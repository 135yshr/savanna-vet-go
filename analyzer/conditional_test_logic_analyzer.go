package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var ConditionalTestLogicAnalyzer = &analysis.Analyzer{
	Name: "conditionaltestlogic",
	Doc:  "テスト内の条件分岐を検出する",
	Run:  runConditionalTestLogic,
}

func runConditionalTestLogic(pass *analysis.Pass) (any, error) {
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
				switch s := n.(type) {
				case *ast.IfStmt:
					pass.Reportf(s.Pos(), "テスト内に if 文があります。テーブル駆動テストやサブテストの使用を検討してください")
				case *ast.SwitchStmt:
					pass.Reportf(s.Pos(), "テスト内に switch 文があります。テーブル駆動テストの使用を検討してください")
				case *ast.ForStmt:
					if !isTableDrivenLoop(s) {
						pass.Reportf(s.Pos(), "テスト内に for 文があります。テーブル駆動テストでない場合は構造を見直してください")
					}
				case *ast.RangeStmt:
					if !isTableDrivenRangeLoop(s) {
						pass.Reportf(s.Pos(), "テスト内に range 文があります。テーブル駆動テストでない場合は構造を見直してください")
					}
				}
				return true
			})
		}
	}
	return nil, nil
}
