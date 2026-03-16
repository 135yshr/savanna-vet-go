package main

import (
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/135yshr/savanna/analyzer"
)

func main() {
	multichecker.Main(analyzer.AllAnalyzers()...)
}
