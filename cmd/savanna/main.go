package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/135yshr/savanna/analyzer"
	"github.com/135yshr/savanna/reporter"
)

var (
	version = "0.1.0"
)

func main() {
	var (
		reportFormat string
		outputDir    string
		failOnSmell  bool
		enabledStr   string
		disabledStr  string
		showVersion  bool
		listSmells   bool
	)

	flag.StringVar(&reportFormat, "format", "console", "レポート形式 (console, json)")
	flag.StringVar(&outputDir, "output", "savanna-reports", "レポート出力ディレクトリ (json用)")
	flag.BoolVar(&failOnSmell, "fail", false, "スメル検出時に終了コード1で終了")
	flag.StringVar(&enabledStr, "enable", "", "有効にするスメル (カンマ区切り)")
	flag.StringVar(&disabledStr, "disable", "", "無効にするスメル (カンマ区切り)")
	flag.BoolVar(&showVersion, "version", false, "バージョンを表示")
	flag.BoolVar(&listSmells, "list", false, "検出可能なスメル一覧を表示")
	flag.Parse()

	if showVersion {
		fmt.Printf("savanna %s\n", version)
		return
	}

	if listSmells {
		fmt.Println("検出可能なテストスメル:")
		for _, t := range analyzer.AllSmellTypes() {
			fmt.Printf("  %-30s %s\n", t, analyzer.SmellDisplayName[t])
		}
		return
	}

	dirs := flag.Args()
	if len(dirs) == 0 {
		dirs = []string{"."}
	}

	scanner := analyzer.NewScanner()

	if enabledStr != "" {
		scanner.EnabledSmells = parseSmellList(enabledStr)
	}
	if disabledStr != "" {
		scanner.DisabledSmells = parseSmellList(disabledStr)
	}

	var allSmells []analyzer.TestSmell
	for _, dir := range dirs {
		smells, err := scanner.ScanDir(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %s のスキャンに失敗: %v\n", dir, err)
			os.Exit(1)
		}
		allSmells = append(allSmells, smells...)
	}

	var rep reporter.Reporter
	switch reportFormat {
	case "json":
		rep = &reporter.JSONReporter{OutputDir: outputDir}
	default:
		rep = reporter.NewConsoleReporter()
	}

	if err := rep.Report(allSmells); err != nil {
		fmt.Fprintf(os.Stderr, "レポート出力エラー: %v\n", err)
		os.Exit(1)
	}

	if failOnSmell && len(allSmells) > 0 {
		os.Exit(1)
	}
}

func parseSmellList(s string) map[analyzer.SmellType]bool {
	m := make(map[analyzer.SmellType]bool)
	for _, name := range strings.Split(s, ",") {
		name = strings.TrimSpace(strings.ToUpper(name))
		if name != "" {
			m[analyzer.SmellType(name)] = true
		}
	}
	return m
}
