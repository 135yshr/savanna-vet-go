package analyzer

import "fmt"

// SmellType はテストスメルの種類を表す。
type SmellType string

const (
	SmellEmptyTest            SmellType = "EMPTY_TEST"
	SmellMissingAssertion     SmellType = "MISSING_ASSERTION"
	SmellSleepyTest           SmellType = "SLEEPY_TEST"
	SmellRedundantPrint       SmellType = "REDUNDANT_PRINT"
	SmellConditionalTestLogic SmellType = "CONDITIONAL_TEST_LOGIC"
	SmellMagicNumberTest      SmellType = "MAGIC_NUMBER_TEST"
	SmellMissingErrorCheck    SmellType = "MISSING_ERROR_CHECK"
	SmellMissingHelper        SmellType = "MISSING_HELPER"
)

// SmellDisplayName はスメルの表示名を返す。
var SmellDisplayName = map[SmellType]string{
	SmellEmptyTest:            "Empty Test",
	SmellMissingAssertion:     "Missing Assertion",
	SmellSleepyTest:           "Sleepy Test",
	SmellRedundantPrint:       "Redundant Print",
	SmellConditionalTestLogic: "Conditional Test Logic",
	SmellMagicNumberTest:      "Magic Number Test",
	SmellMissingErrorCheck:    "Missing Error Check",
	SmellMissingHelper:        "Missing t.Helper()",
}

// SmellMessage はスメル検出時のライオンのメッセージを返す。
var SmellMessage = map[SmellType]string{
	SmellEmptyTest:            "空のテストとかお前それ@t_wadaの前でも同じ事言えんの？",
	SmellMissingAssertion:     "アサーションのないテストとかお前それ@t_wadaの前でも同じ事言えんの？",
	SmellSleepyTest:           "テストでsleepとかお前それ@t_wadaの前でも同じ事言えんの？",
	SmellRedundantPrint:       "テストでfmt.Printlnとかお前それ@t_wadaの前でも同じ事言えんの？",
	SmellConditionalTestLogic: "テストで条件分岐とかお前それ@t_wadaの前でも同じ事言えんの？",
	SmellMagicNumberTest:      "テストでマジックナンバーとかお前それ@t_wadaの前でも同じ事言えんの？",
	SmellMissingErrorCheck:    "errorを無視するテストとかお前それ@t_wadaの前でも同じ事言えんの？",
	SmellMissingHelper:        "t.Helper()を書かないヘルパーとかお前それ@t_wadaの前でも同じ事言えんの？",
}

// TestSmell は検出されたテストスメルを表す。
type TestSmell struct {
	Type     SmellType
	File     string
	FuncName string
	Line     int
	Message  string
}

func (s TestSmell) String() string {
	return fmt.Sprintf("[%s] %s:%d in %s() - %s", s.Type, s.File, s.Line, s.FuncName, s.Message)
}

// AllSmellTypes は全スメルタイプを返す。
func AllSmellTypes() []SmellType {
	return []SmellType{
		SmellEmptyTest,
		SmellMissingAssertion,
		SmellSleepyTest,
		SmellRedundantPrint,
		SmellConditionalTestLogic,
		SmellMagicNumberTest,
		SmellMissingErrorCheck,
		SmellMissingHelper,
	}
}
