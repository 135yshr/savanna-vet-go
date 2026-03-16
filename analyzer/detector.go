package analyzer

import (
	"go/ast"
	"strings"
)

// isTestFunc はテスト関数かどうかを判定する。
func isTestFunc(fn *ast.FuncDecl) bool {
	if fn == nil || fn.Type == nil || fn.Type.Params == nil {
		return false
	}
	if !strings.HasPrefix(fn.Name.Name, "Test") {
		return false
	}
	params := fn.Type.Params.List
	if len(params) != 1 {
		return false
	}
	pt, ok := params[0].Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	sel, ok := pt.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == "testing" && sel.Sel.Name == "T"
}

// isHelperFunc はテストヘルパー関数（*testing.T を引数に取るが Test で始まらない）かどうかを判定する。
func isHelperFunc(fn *ast.FuncDecl) bool {
	if fn == nil || fn.Type == nil || fn.Type.Params == nil {
		return false
	}
	if strings.HasPrefix(fn.Name.Name, "Test") {
		return false
	}
	for _, p := range fn.Type.Params.List {
		pt, ok := p.Type.(*ast.StarExpr)
		if !ok {
			continue
		}
		sel, ok := pt.X.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		ident, ok := sel.X.(*ast.Ident)
		if !ok {
			continue
		}
		if ident.Name == "testing" && sel.Sel.Name == "T" {
			return true
		}
	}
	return false
}

// assertionMethods は t.Error, t.Fatal 等のアサーション系メソッド名。
var assertionMethods = map[string]bool{
	"Error":   true,
	"Errorf":  true,
	"Fatal":   true,
	"Fatalf":  true,
	"Fail":    true,
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

// hasAssertion はブロック内にアサーション呼び出しがあるか判定する。
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
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if ok {
			if assertionMethods[sel.Sel.Name] {
				found = true
				return false
			}
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

// isAssertionCall はアサーション呼び出しかどうかを判定する。
func isAssertionCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	if assertionMethods[sel.Sel.Name] {
		return true
	}
	if ident, ok := sel.X.(*ast.Ident); ok {
		if (ident.Name == "assert" || ident.Name == "require") && testifyMethods[sel.Sel.Name] {
			return true
		}
	}
	return false
}

// printFuncs は fmt.Print 系の関数名。
var printFuncs = map[string]bool{
	"Print":   true,
	"Printf":  true,
	"Println": true,
}

// allowedNumbers は許容される数値リテラル。
var allowedNumbers = map[string]bool{
	"0":     true,
	"1":     true,
	"2":     true,
	"-1":    true,
	"0.0":   true,
	"1.0":   true,
	"true":  true,
	"false": true,
}

// isTableDrivenLoop はテーブル駆動テストの for-range ループかどうかを簡易判定する。
func isTableDrivenLoop(_ *ast.ForStmt) bool {
	return false
}

// hasHelperCall はブロック内に t.Helper() 呼び出しがあるか判定する。
func hasHelperCall(body *ast.BlockStmt) bool {
	for _, stmt := range body.List {
		exprStmt, ok := stmt.(*ast.ExprStmt)
		if !ok {
			continue
		}
		call, ok := exprStmt.X.(*ast.CallExpr)
		if !ok {
			continue
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		if sel.Sel.Name == "Helper" {
			return true
		}
	}
	return false
}
