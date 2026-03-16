package analyzer

import (
	"go/ast"
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
				if len(assign.Lhs) < 2 {
					return true
				}
				last := assign.Lhs[len(assign.Lhs)-1]
				ident, ok := last.(*ast.Ident)
				if !ok {
					return true
				}
				if ident.Name == "_" {
					if _, ok := assign.Rhs[0].(*ast.CallExpr); ok {
						pass.Reportf(assign.Pos(), "error の戻り値が無視されています。テストでもエラーチェックを行ってください")
					}
				}
				return true
			})
		}
	}
	return nil, nil
}
