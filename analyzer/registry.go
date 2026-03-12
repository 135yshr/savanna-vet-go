package analyzer

// AllDetectors は全検出器を返す。
func AllDetectors() []Detector {
	return []Detector{
		&EmptyTestDetector{},
		&MissingAssertionDetector{},
		&SleepyTestDetector{},
		&RedundantPrintDetector{},
		&ConditionalTestLogicDetector{},
		&MagicNumberTestDetector{},
		&MissingErrorCheckDetector{},
		&MissingHelperDetector{},
	}
}
