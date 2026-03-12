package analyzer

import (
	"go/ast"
	"go/token"
)

// ConditionalTestLogicDetector はテスト内の条件分岐を検出する。
type ConditionalTestLogicDetector struct{}

func (d *ConditionalTestLogicDetector) Type() SmellType { return SmellConditionalTestLogic }

func (d *ConditionalTestLogicDetector) Detect(fset *token.FileSet, file *ast.File, fn *ast.FuncDecl) []TestSmell {
	if !isTestFunc(fn) {
		return nil
	}
	if fn.Body == nil {
		return nil
	}

	var smells []TestSmell
	for _, stmt := range fn.Body.List {
		// t.Run 内のサブテストは除外するため、直接の子ステートメントのみ検査
		switch s := stmt.(type) {
		case *ast.IfStmt:
			pos := fset.Position(s.Pos())
			smells = append(smells, TestSmell{
				Type:     SmellConditionalTestLogic,
				File:     pos.Filename,
				FuncName: fn.Name.Name,
				Line:     pos.Line,
				Message:  "テスト内に if 文があります。テーブル駆動テストやサブテストの使用を検討してください",
			})
		case *ast.SwitchStmt:
			pos := fset.Position(s.Pos())
			smells = append(smells, TestSmell{
				Type:     SmellConditionalTestLogic,
				File:     pos.Filename,
				FuncName: fn.Name.Name,
				Line:     pos.Line,
				Message:  "テスト内に switch 文があります。テーブル駆動テストの使用を検討してください",
			})
		case *ast.ForStmt:
			if !isTableDrivenLoop(s) {
				pos := fset.Position(s.Pos())
				smells = append(smells, TestSmell{
					Type:     SmellConditionalTestLogic,
					File:     pos.Filename,
					FuncName: fn.Name.Name,
					Line:     pos.Line,
					Message:  "テスト内に for 文があります。テーブル駆動テストでない場合は構造を見直してください",
				})
			}
		}
	}
	return smells
}

// isTableDrivenLoop はテーブル駆動テストの for-range ループかどうかを簡易判定する。
// ast.ForStmt は for i := 0; i < N; i++ 形式で、range は ast.RangeStmt なので
// 通常の for 文はテーブル駆動ではないと判断する。
func isTableDrivenLoop(_ *ast.ForStmt) bool {
	return false
}
