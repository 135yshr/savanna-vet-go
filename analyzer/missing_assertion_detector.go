package analyzer

import (
	"go/ast"
	"go/token"
)

// MissingAssertionDetector はアサーションのないテスト関数を検出する。
type MissingAssertionDetector struct{}

func (d *MissingAssertionDetector) Type() SmellType { return SmellMissingAssertion }

func (d *MissingAssertionDetector) Detect(fset *token.FileSet, file *ast.File, fn *ast.FuncDecl) []TestSmell {
	if !isTestFunc(fn) {
		return nil
	}
	if fn.Body == nil || len(fn.Body.List) == 0 {
		return nil // EmptyTestDetector が処理する
	}
	if hasAssertion(fn.Body) {
		return nil
	}
	pos := fset.Position(fn.Pos())
	return []TestSmell{{
		Type:     SmellMissingAssertion,
		File:     pos.Filename,
		FuncName: fn.Name.Name,
		Line:     pos.Line,
		Message:  "テスト関数にアサーションがありません",
	}}
}

// assertionMethods は t.Error, t.Fatal 等のアサーション系メソッド名。
var assertionMethods = map[string]bool{
	"Error":  true,
	"Errorf": true,
	"Fatal":  true,
	"Fatalf": true,
	"Fail":   true,
	"FailNow": true,
}

// testifyMethods は testify の assert/require 関数名。
var testifyMethods = map[string]bool{
	"Equal":         true,
	"NotEqual":      true,
	"True":          true,
	"False":         true,
	"Nil":           true,
	"NotNil":        true,
	"NoError":       true,
	"Error":         true,
	"Contains":      true,
	"NotContains":   true,
	"Len":           true,
	"Empty":         true,
	"NotEmpty":      true,
	"Zero":          true,
	"NotZero":       true,
	"Greater":       true,
	"Less":          true,
	"Panics":        true,
	"JSONEq":        true,
	"ElementsMatch": true,
}

func hasAssertion(body *ast.BlockStmt) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		if found {
			return false
		}
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		// t.Error(), t.Fatal() 等
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if ok {
			if assertionMethods[sel.Sel.Name] {
				found = true
				return false
			}
			// testify: assert.Equal(), require.NoError() 等
			if ident, ok := sel.X.(*ast.Ident); ok {
				if (ident.Name == "assert" || ident.Name == "require") && testifyMethods[sel.Sel.Name] {
					found = true
					return false
				}
			}
		}
		return true
	})
	return found
}
