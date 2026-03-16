package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var EmptyTestAnalyzer = &analysis.Analyzer{
	Name: "emptytest",
	Doc:  "テスト関数の本体が空のケースを検出する",
	Run:  runEmptyTest,
}

func runEmptyTest(pass *analysis.Pass) (any, error) {
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
			if fn.Body == nil || len(fn.Body.List) == 0 {
				pass.Reportf(fn.Pos(), "テスト関数の本体が空です")
			}
		}
	}
	return nil, nil
}
