package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var MissingAssertionAnalyzer = &analysis.Analyzer{
	Name: "missingassertion",
	Doc:  "アサーションのないテスト関数を検出する",
	Run:  runMissingAssertion,
}

func runMissingAssertion(pass *analysis.Pass) (any, error) {
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
				continue
			}
			tName := testingParamName(fn)
			if !hasAssertionWithParam(fn.Body, tName) {
				pass.Reportf(fn.Pos(), "テスト関数にアサーションがありません")
			}
		}
	}
	return nil, nil
}
