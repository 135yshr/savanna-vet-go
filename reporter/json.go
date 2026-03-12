package reporter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/135yshr/savanna/analyzer"
)

// JSONReporter はJSON形式でレポートを出力する。
type JSONReporter struct {
	OutputDir string
}

type jsonReport struct {
	Timestamp   string       `json:"timestamp"`
	TotalSmells int          `json:"totalSmells"`
	Smells      []jsonSmell  `json:"smells"`
	Summary     []jsonCount  `json:"summary"`
}

type jsonSmell struct {
	Type     string `json:"type"`
	File     string `json:"file"`
	Function string `json:"function"`
	Line     int    `json:"line"`
	Message  string `json:"message"`
}

type jsonCount struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func (r *JSONReporter) Report(smells []analyzer.TestSmell) error {
	if err := os.MkdirAll(r.OutputDir, 0o755); err != nil {
		return fmt.Errorf("ディレクトリの作成に失敗: %w", err)
	}

	byType := make(map[analyzer.SmellType]int)
	var items []jsonSmell
	for _, s := range smells {
		byType[s.Type]++
		items = append(items, jsonSmell{
			Type:     string(s.Type),
			File:     s.File,
			Function: s.FuncName,
			Line:     s.Line,
			Message:  s.Message,
		})
	}

	var summary []jsonCount
	for t, c := range byType {
		summary = append(summary, jsonCount{
			Type:  string(t),
			Name:  analyzer.SmellDisplayName[t],
			Count: c,
		})
	}

	report := jsonReport{
		Timestamp:   time.Now().Format(time.RFC3339),
		TotalSmells: len(smells),
		Smells:      items,
		Summary:     summary,
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("JSONの生成に失敗: %w", err)
	}

	path := filepath.Join(r.OutputDir, "test-smells.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗: %w", err)
	}

	fmt.Printf("レポートを出力しました: %s\n", path)
	return nil
}
