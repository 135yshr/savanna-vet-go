package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var SleepyTestAnalyzer = &analysis.Analyzer{
	Name: "sleepytest",
	Doc:  "テストでの time.Sleep() の使用を検出する",
	Run:  runSleepyTest,
}

func runSleepyTest(pass *analysis.Pass) (any, error) {
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
				ident, ok := sel.X.(*ast.Ident)
				if !ok {
					return true
				}
				if ident.Name == "time" && sel.Sel.Name == "Sleep" {
					pass.Reportf(call.Pos(), "time.Sleep() はテストを不安定にします。チャネルやsync.WaitGroupの使用を検討してください")
				}
				return true
			})
		}
	}
	return nil, nil
}
