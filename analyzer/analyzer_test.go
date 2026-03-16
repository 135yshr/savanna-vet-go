package analyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/135yshr/savanna/analyzer"
)

var analyzerPackages = map[string]string{
	"emptytest":            "emptytest",
	"missingassertion":     "missingassertion",
	"sleepytest":           "sleepytest",
	"redundantprint":       "redundantprint",
	"conditionaltestlogic": "conditionaltestlogic",
	"magicnumbertest":      "magicnumbertest",
	"missingerrorcheck":    "missingerrorcheck",
	"missinghelper":        "missinghelper",
}

func TestAllAnalyzers(t *testing.T) {
	testdata := analysistest.TestData()
	for _, a := range analyzer.AllAnalyzers() {
		pkg, ok := analyzerPackages[a.Name]
		if !ok {
			t.Fatalf("analyzer %q のテストデータパッケージが analyzerPackages に未登録です", a.Name)
		}
		t.Run(a.Name, func(t *testing.T) {
			analysistest.Run(t, testdata, a, pkg)
		})
	}
}
