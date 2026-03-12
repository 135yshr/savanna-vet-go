package reporter

import "github.com/135yshr/savanna/analyzer"

// Reporter はスメルレポートを出力するインターフェース。
type Reporter interface {
	Report(smells []analyzer.TestSmell) error
}
