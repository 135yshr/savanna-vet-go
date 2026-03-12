package analyzer

import (
	"go/ast"
	"go/token"
	"strings"
)

// Detector はテストスメルを検出するインターフェース。
type Detector interface {
	Type() SmellType
	Detect(fset *token.FileSet, file *ast.File, fn *ast.FuncDecl) []TestSmell
}

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
