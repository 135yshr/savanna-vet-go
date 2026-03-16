package analyzer

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var MissingErrorCheckAnalyzer = &analysis.Analyzer{
	Name: "missingerrorcheck",
	Doc:  "error 戻り値を無視しているケースを検出する",
	Run:  runMissingErrorCheck,
}

func runMissingErrorCheck(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		if !strings.HasSuffix(filename, "_test.go") {
			continue
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if !isTestFunc(fn) {
				continue
			}
			if fn.Body == nil {
				continue
			}
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				assign, ok := n.(*ast.AssignStmt)
				if !ok {
					return true
				}
				if len(assign.Lhs) < 2 || len(assign.Rhs) == 0 {
					return true
				}
				last := assign.Lhs[len(assign.Lhs)-1]
				ident, ok := last.(*ast.Ident)
				if !ok {
					return true
				}
				if ident.Name != "_" {
					return true
				}
				call, ok := assign.Rhs[0].(*ast.CallExpr)
				if !ok {
					return true
				}
				// 型情報を使って最後の戻り値が error 型か検証
				if isLastReturnError(pass, call) {
					pass.Reportf(assign.Pos(), "error の戻り値が無視されています。テストでもエラーチェックを行ってください")
				}
				return true
			})
		}
	}
	return nil, nil
}

// isLastReturnError は関数呼び出しの最後の戻り値が error 型かどうかを判定する。
func isLastReturnError(pass *analysis.Pass, call *ast.CallExpr) bool {
	t := pass.TypesInfo.TypeOf(call)
	if t == nil {
		return false
	}
	// 複数戻り値の場合
	if tuple, ok := t.(*types.Tuple); ok {
		if tuple.Len() == 0 {
			return false
		}
		lastType := tuple.At(tuple.Len() - 1).Type()
		return isErrorType(lastType)
	}
	// 単一戻り値の場合
	return isErrorType(t)
}

// isErrorType は型が error インターフェースかどうかを判定する。
func isErrorType(t types.Type) bool {
	return types.Implements(t, errorInterface())
}

func errorInterface() *types.Interface {
	return types.Universe.Lookup("error").Type().Underlying().(*types.Interface)
}
