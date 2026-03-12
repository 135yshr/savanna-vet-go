package analyzer

import (
	"go/ast"
	"go/token"
)

// SleepyTestDetector は time.Sleep() を使用しているテストを検出する。
type SleepyTestDetector struct{}

func (d *SleepyTestDetector) Type() SmellType { return SmellSleepyTest }

func (d *SleepyTestDetector) Detect(fset *token.FileSet, file *ast.File, fn *ast.FuncDecl) []TestSmell {
	if !isTestFunc(fn) {
		return nil
	}
	if fn.Body == nil {
		return nil
	}

	var smells []TestSmell
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
			pos := fset.Position(call.Pos())
			smells = append(smells, TestSmell{
				Type:     SmellSleepyTest,
				File:     pos.Filename,
				FuncName: fn.Name.Name,
				Line:     pos.Line,
				Message:  "time.Sleep() はテストを不安定にします。チャネルやsync.WaitGroupの使用を検討してください",
			})
		}
		return true
	})
	return smells
}
