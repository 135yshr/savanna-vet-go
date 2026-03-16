package analyzer

import "golang.org/x/tools/go/analysis"

// AllAnalyzers は全 Analyzer を返す。
func AllAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		EmptyTestAnalyzer,
		MissingAssertionAnalyzer,
		SleepyTestAnalyzer,
		RedundantPrintAnalyzer,
		ConditionalTestLogicAnalyzer,
		MagicNumberTestAnalyzer,
		MissingErrorCheckAnalyzer,
		MissingHelperAnalyzer,
	}
}
