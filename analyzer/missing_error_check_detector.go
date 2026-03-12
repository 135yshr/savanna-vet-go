package analyzer

import (
	"go/ast"
	"go/token"
)

// MissingErrorCheckDetector は error 戻り値を _ で無視しているケースを検出する。
type MissingErrorCheckDetector struct{}

func (d *MissingErrorCheckDetector) Type() SmellType { return SmellMissingErrorCheck }

func (d *MissingErrorCheckDetector) Detect(fset *token.FileSet, file *ast.File, fn *ast.FuncDecl) []TestSmell {
	if !isTestFunc(fn) {
		return nil
	}
	if fn.Body == nil {
		return nil
	}

	var smells []TestSmell
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}
		// 最後の戻り値が _ で、変数名が err パターンのものを検出
		if len(assign.Lhs) < 2 {
			return true
		}
		last := assign.Lhs[len(assign.Lhs)-1]
		ident, ok := last.(*ast.Ident)
		if !ok {
			return true
		}
		if ident.Name == "_" {
			// 関数呼び出しの戻り値で最後が _ の場合、error を無視している可能性
			if _, ok := assign.Rhs[0].(*ast.CallExpr); ok {
				pos := fset.Position(assign.Pos())
				smells = append(smells, TestSmell{
					Type:     SmellMissingErrorCheck,
					File:     pos.Filename,
					FuncName: fn.Name.Name,
					Line:     pos.Line,
					Message:  "error の戻り値が無視されています。テストでもエラーチェックを行ってください",
				})
			}
		}
		return true
	})
	return smells
}
