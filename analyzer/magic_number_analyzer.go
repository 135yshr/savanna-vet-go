package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var MagicNumberTestAnalyzer = &analysis.Analyzer{
	Name: "magicnumbertest",
	Doc:  "アサーション内のマジックナンバーを検出する",
	Run:  runMagicNumberTest,
}

func runMagicNumberTest(pass *analysis.Pass) (any, error) {
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
			tName := testingParamName(fn)
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				if !isAssertionCallWithParam(call, tName) {
					return true
				}
				for _, arg := range call.Args {
					if lit, ok := arg.(*ast.BasicLit); ok {
						if lit.Kind == token.INT || lit.Kind == token.FLOAT {
							if !allowedNumbers[lit.Value] {
								pass.Reportf(lit.Pos(), "マジックナンバー %s には名前付き定数を使用してください", lit.Value)
							}
						}
					}
					if unary, ok := arg.(*ast.UnaryExpr); ok && unary.Op == token.SUB {
						if lit, ok := unary.X.(*ast.BasicLit); ok {
							if lit.Kind == token.INT || lit.Kind == token.FLOAT {
								val := "-" + lit.Value
								if !allowedNumbers[val] {
									if _, err := strconv.ParseFloat(lit.Value, 64); err == nil {
										pass.Reportf(unary.Pos(), "マジックナンバー %s には名前付き定数を使用してください", val)
									}
								}
							}
						}
					}
				}
				return true
			})
		}
	}
	return nil, nil
}
