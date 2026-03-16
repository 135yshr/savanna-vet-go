package analyzer

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var SleepyTestAnalyzer = &analysis.Analyzer{
	Name: "sleepytest",
	Doc:  "テストでの time.Sleep() の使用を検出する",
	Run:  runSleepyTest,
}

func runSleepyTest(pass *analysis.Pass) (any, error) {
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
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				if sel.Sel.Name != "Sleep" {
					return true
				}
				if isPkgCall(pass, sel, "time") {
					pass.Reportf(call.Pos(), "time.Sleep() はテストを不安定にします。チャネルやsync.WaitGroupの使用を検討してください")
				}
				return true
			})
		}
	}
	return nil, nil
}

// isPkgCall は SelectorExpr のレシーバが指定パッケージパスのパッケージ名か判定する。
// エイリアスインポートやシャドウされたローカル変数にも正しく対応する。
func isPkgCall(pass *analysis.Pass, sel *ast.SelectorExpr, pkgPath string) bool {
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	obj := pass.TypesInfo.Uses[ident]
	if obj == nil {
		return false
	}
	pkgName, ok := obj.(*types.PkgName)
	if !ok {
		return false
	}
	return pkgName.Imported().Path() == pkgPath
}
