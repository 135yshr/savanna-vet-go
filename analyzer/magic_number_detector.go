package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
)

// MagicNumberTestDetector はアサーション内のマジックナンバーを検出する。
type MagicNumberTestDetector struct{}

func (d *MagicNumberTestDetector) Type() SmellType { return SmellMagicNumberTest }

// 許容される数値リテラル
var allowedNumbers = map[string]bool{
	"0":    true,
	"1":    true,
	"2":    true,
	"-1":   true,
	"0.0":  true,
	"1.0":  true,
	"true": true,
	"false": true,
}

func (d *MagicNumberTestDetector) Detect(fset *token.FileSet, file *ast.File, fn *ast.FuncDecl) []TestSmell {
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
		if !isAssertionCall(call) {
			return true
		}
		for _, arg := range call.Args {
			if lit, ok := arg.(*ast.BasicLit); ok {
				if lit.Kind == token.INT || lit.Kind == token.FLOAT {
					if !allowedNumbers[lit.Value] {
						pos := fset.Position(lit.Pos())
						smells = append(smells, TestSmell{
							Type:     SmellMagicNumberTest,
							File:     pos.Filename,
							FuncName: fn.Name.Name,
							Line:     pos.Line,
							Message:  "マジックナンバー " + lit.Value + " には名前付き定数を使用してください",
						})
					}
				}
			}
			// 負の数値: -42 のようなケース
			if unary, ok := arg.(*ast.UnaryExpr); ok && unary.Op == token.SUB {
				if lit, ok := unary.X.(*ast.BasicLit); ok {
					if lit.Kind == token.INT || lit.Kind == token.FLOAT {
						val := "-" + lit.Value
						if !allowedNumbers[val] {
							if _, err := strconv.ParseFloat(lit.Value, 64); err == nil {
								pos := fset.Position(unary.Pos())
								smells = append(smells, TestSmell{
									Type:     SmellMagicNumberTest,
									File:     pos.Filename,
									FuncName: fn.Name.Name,
									Line:     pos.Line,
									Message:  "マジックナンバー " + val + " には名前付き定数を使用してください",
								})
							}
						}
					}
				}
			}
		}
		return true
	})
	return smells
}

func isAssertionCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	// t.Error, t.Fatal 等
	if assertionMethods[sel.Sel.Name] {
		return true
	}
	// assert.Equal, require.Equal 等
	if ident, ok := sel.X.(*ast.Ident); ok {
		if (ident.Name == "assert" || ident.Name == "require") && testifyMethods[sel.Sel.Name] {
			return true
		}
	}
	return false
}
