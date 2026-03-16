package main

import (
	"golang.org/x/tools/go/analysis/unitchecker"

	"github.com/135yshr/savanna/analyzer"
)

func main() {
	unitchecker.Main(analyzer.AllAnalyzers()...)
}
