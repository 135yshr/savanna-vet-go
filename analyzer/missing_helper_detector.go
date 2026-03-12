package analyzer

import (
	"go/ast"
	"go/token"
)

// MissingHelperDetector は *testing.T を受け取るヘルパー関数で t.Helper() が呼ばれていないケースを検出する。
type MissingHelperDetector struct{}

func (d *MissingHelperDetector) Type() SmellType { return SmellMissingHelper }

func (d *MissingHelperDetector) Detect(fset *token.FileSet, file *ast.File, fn *ast.FuncDecl) []TestSmell {
	if !isHelperFunc(fn) {
		return nil
	}
	if fn.Body == nil {
		return nil
	}

	if hasHelperCall(fn.Body) {
		return nil
	}

	pos := fset.Position(fn.Pos())
	return []TestSmell{{
		Type:     SmellMissingHelper,
		File:     pos.Filename,
		FuncName: fn.Name.Name,
		Line:     pos.Line,
		Message:  "テストヘルパー関数に t.Helper() がありません。エラー時の行番号が不正確になります",
	}}
}

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
