package reporter

import (
	"fmt"
	"io"
	"os"

	"github.com/135yshr/savanna/analyzer"
	"github.com/135yshr/savanna/banner"
)

// ConsoleReporter はコンソールにレポートを出力する。
type ConsoleReporter struct {
	Writer io.Writer
}

// NewConsoleReporter は標準出力に書き込むレポーターを返す。
func NewConsoleReporter() *ConsoleReporter {
	return &ConsoleReporter{Writer: os.Stdout}
}

func (r *ConsoleReporter) Report(smells []analyzer.TestSmell) error {
	if len(smells) == 0 {
		fmt.Fprintln(r.Writer, "テストスメルは検出されませんでした。")
		return nil
	}

	// スメルタイプごとに集計
	byType := make(map[analyzer.SmellType][]analyzer.TestSmell)
	for _, s := range smells {
		byType[s.Type] = append(byType[s.Type], s)
	}

	// ライオンの吠え（最初に見つかったスメルのメッセージ）
	firstType := smells[0].Type
	if msg, ok := analyzer.SmellMessage[firstType]; ok {
		fmt.Fprint(r.Writer, banner.Roar(msg))
	}

	fmt.Fprintf(r.Writer, "\n%d 件のテストスメルが検出されました:\n\n", len(smells))

	for smellType, typeSmells := range byType {
		displayName := analyzer.SmellDisplayName[smellType]
		fmt.Fprintf(r.Writer, "── %s (%d件) ──\n", displayName, len(typeSmells))
		for _, s := range typeSmells {
			fmt.Fprintf(r.Writer, "  %s:%d %s() - %s\n", s.File, s.Line, s.FuncName, s.Message)
		}
		fmt.Fprintln(r.Writer)
	}

	return nil
}
