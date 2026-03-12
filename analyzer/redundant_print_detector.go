package analyzer

import (
	"go/ast"
	"go/token"
)

// RedundantPrintDetector は fmt.Print 系の呼び出しを検出する。
type RedundantPrintDetector struct{}

func (d *RedundantPrintDetector) Type() SmellType { return SmellRedundantPrint }

var printFuncs = map[string]bool{
	"Print":   true,
	"Printf":  true,
	"Println": true,
}

func (d *RedundantPrintDetector) Detect(fset *token.FileSet, file *ast.File, fn *ast.FuncDecl) []TestSmell {
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
		if ident.Name == "fmt" && printFuncs[sel.Sel.Name] {
			pos := fset.Position(call.Pos())
			smells = append(smells, TestSmell{
				Type:     SmellRedundantPrint,
				File:     pos.Filename,
				FuncName: fn.Name.Name,
				Line:     pos.Line,
				Message:  "fmt." + sel.Sel.Name + "() の代わりに t.Log() を使用してください",
			})
		}
		return true
	})
	return smells
}
