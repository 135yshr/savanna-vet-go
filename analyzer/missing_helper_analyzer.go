package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var MissingHelperAnalyzer = &analysis.Analyzer{
	Name: "missinghelper",
	Doc:  "テストヘルパー関数で t.Helper() が呼ばれていないケースを検出する",
	Run:  runMissingHelper,
}

func runMissingHelper(pass *analysis.Pass) (any, error) {
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
			if !isHelperFunc(fn) {
				continue
			}
			if fn.Body == nil {
				continue
			}
			if !hasHelperCall(fn.Body) {
				pass.Reportf(fn.Pos(), "テストヘルパー関数に t.Helper() がありません。エラー時の行番号が不正確になります")
			}
		}
	}
	return nil, nil
}
