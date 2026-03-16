package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var RedundantPrintAnalyzer = &analysis.Analyzer{
	Name: "redundantprint",
	Doc:  "テストでの fmt.Print 系呼び出しを検出する",
	Run:  runRedundantPrint,
}

func runRedundantPrint(pass *analysis.Pass) (any, error) {
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
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				if !printFuncs[sel.Sel.Name] {
					return true
				}
				if isPkgCall(pass, sel, "fmt") {
					pass.Reportf(call.Pos(), "fmt.%s() の代わりに t.Log() を使用してください", sel.Sel.Name)
				}
				return true
			})
		}
	}
	return nil, nil
}
