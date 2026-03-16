package main

import (
	"golang.org/x/tools/go/analysis/unitchecker"

	"github.com/135yshr/savanna-vet-go/analyzer"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	unitchecker.Main(analyzer.AllAnalyzers()...)
}
