package analyzer

import (
	"go/ast"
	"go/token"
)

// EmptyTestDetector は本体が空のテスト関数を検出する。
type EmptyTestDetector struct{}

func (d *EmptyTestDetector) Type() SmellType { return SmellEmptyTest }

func (d *EmptyTestDetector) Detect(fset *token.FileSet, file *ast.File, fn *ast.FuncDecl) []TestSmell {
	if !isTestFunc(fn) {
		return nil
	}
	if fn.Body == nil || len(fn.Body.List) == 0 {
		pos := fset.Position(fn.Pos())
		return []TestSmell{{
			Type:     SmellEmptyTest,
			File:     pos.Filename,
			FuncName: fn.Name.Name,
			Line:     pos.Line,
			Message:  "テスト関数の本体が空です",
		}}
	}
	return nil
}
