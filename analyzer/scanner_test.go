package analyzer_test

import (
	"testing"

	"github.com/135yshr/savanna/analyzer"
)

func TestEmptyTest(t *testing.T) {
	src := `package foo

import "testing"

func TestEmpty(t *testing.T) {
}
`
	scanner := analyzer.NewScanner()
	smells, err := scanner.ScanSource(src)
	if err != nil {
		t.Fatalf("ScanSource failed: %v", err)
	}
	if len(smells) == 0 {
		t.Fatal("expected at least one smell for empty test")
	}
	found := false
	for _, s := range smells {
		if s.Type == analyzer.SmellEmptyTest {
			found = true
		}
	}
	if !found {
		t.Errorf("expected SmellEmptyTest, got %v", smells)
	}
}

func TestMissingAssertion(t *testing.T) {
	src := `package foo

import "testing"

func TestNoAssert(t *testing.T) {
	x := 1 + 1
	_ = x
}
`
	scanner := analyzer.NewScanner()
	smells, err := scanner.ScanSource(src)
	if err != nil {
		t.Fatalf("ScanSource failed: %v", err)
	}
	found := false
	for _, s := range smells {
		if s.Type == analyzer.SmellMissingAssertion {
			found = true
		}
	}
	if !found {
		t.Errorf("expected SmellMissingAssertion, got %v", smells)
	}
}

func TestMissingAssertionWithAssert(t *testing.T) {
	src := `package foo

import "testing"

func TestWithAssert(t *testing.T) {
	x := 1 + 1
	if x != 2 {
		t.Errorf("expected 2, got %d", x)
	}
}
`
	scanner := analyzer.NewScanner()
	smells, err := scanner.ScanSource(src)
	if err != nil {
		t.Fatalf("ScanSource failed: %v", err)
	}
	for _, s := range smells {
		if s.Type == analyzer.SmellMissingAssertion {
			t.Error("should not detect MissingAssertion when t.Errorf is used")
		}
	}
}

func TestSleepyTest(t *testing.T) {
	src := `package foo

import (
	"testing"
	"time"
)

func TestWithSleep(t *testing.T) {
	time.Sleep(time.Second)
	t.Error("dummy")
}
`
	scanner := analyzer.NewScanner()
	smells, err := scanner.ScanSource(src)
	if err != nil {
		t.Fatalf("ScanSource failed: %v", err)
	}
	found := false
	for _, s := range smells {
		if s.Type == analyzer.SmellSleepyTest {
			found = true
		}
	}
	if !found {
		t.Errorf("expected SmellSleepyTest, got %v", smells)
	}
}

func TestRedundantPrint(t *testing.T) {
	src := `package foo

import (
	"fmt"
	"testing"
)

func TestWithPrint(t *testing.T) {
	fmt.Println("debug")
	t.Error("dummy")
}
`
	scanner := analyzer.NewScanner()
	smells, err := scanner.ScanSource(src)
	if err != nil {
		t.Fatalf("ScanSource failed: %v", err)
	}
	found := false
	for _, s := range smells {
		if s.Type == analyzer.SmellRedundantPrint {
			found = true
		}
	}
	if !found {
		t.Errorf("expected SmellRedundantPrint, got %v", smells)
	}
}

func TestConditionalTestLogic(t *testing.T) {
	src := `package foo

import "testing"

func TestWithIf(t *testing.T) {
	x := 1
	if x == 1 {
		t.Error("one")
	}
}
`
	scanner := analyzer.NewScanner()
	smells, err := scanner.ScanSource(src)
	if err != nil {
		t.Fatalf("ScanSource failed: %v", err)
	}
	found := false
	for _, s := range smells {
		if s.Type == analyzer.SmellConditionalTestLogic {
			found = true
		}
	}
	if !found {
		t.Errorf("expected SmellConditionalTestLogic, got %v", smells)
	}
}

func TestMagicNumber(t *testing.T) {
	src := `package foo

import "testing"

func TestMagic(t *testing.T) {
	t.Errorf("expected %d", 42)
}
`
	scanner := analyzer.NewScanner()
	smells, err := scanner.ScanSource(src)
	if err != nil {
		t.Fatalf("ScanSource failed: %v", err)
	}
	found := false
	for _, s := range smells {
		if s.Type == analyzer.SmellMagicNumberTest {
			found = true
		}
	}
	if !found {
		t.Errorf("expected SmellMagicNumberTest, got %v", smells)
	}
}

func TestMissingErrorCheck(t *testing.T) {
	src := `package foo

import (
	"os"
	"testing"
)

func TestIgnoreErr(t *testing.T) {
	f, _ := os.Open("test.txt")
	_ = f
	t.Error("dummy")
}
`
	scanner := analyzer.NewScanner()
	smells, err := scanner.ScanSource(src)
	if err != nil {
		t.Fatalf("ScanSource failed: %v", err)
	}
	found := false
	for _, s := range smells {
		if s.Type == analyzer.SmellMissingErrorCheck {
			found = true
		}
	}
	if !found {
		t.Errorf("expected SmellMissingErrorCheck, got %v", smells)
	}
}

func TestMissingHelper(t *testing.T) {
	src := `package foo

import "testing"

func setupTest(t *testing.T) {
	// no t.Helper() call
}
`
	scanner := analyzer.NewScanner()
	smells, err := scanner.ScanSource(src)
	if err != nil {
		t.Fatalf("ScanSource failed: %v", err)
	}
	found := false
	for _, s := range smells {
		if s.Type == analyzer.SmellMissingHelper {
			found = true
		}
	}
	if !found {
		t.Errorf("expected SmellMissingHelper, got %v", smells)
	}
}

func TestMissingHelperWithHelper(t *testing.T) {
	src := `package foo

import "testing"

func setupTest(t *testing.T) {
	t.Helper()
}
`
	scanner := analyzer.NewScanner()
	smells, err := scanner.ScanSource(src)
	if err != nil {
		t.Fatalf("ScanSource failed: %v", err)
	}
	for _, s := range smells {
		if s.Type == analyzer.SmellMissingHelper {
			t.Error("should not detect MissingHelper when t.Helper() is present")
		}
	}
}

func TestCleanTest(t *testing.T) {
	src := `package foo

import "testing"

func TestGood(t *testing.T) {
	x := 1 + 1
	if x != 2 {
		t.Errorf("expected 2, got %d", x)
	}
}
`
	scanner := analyzer.NewScanner()
	smells, err := scanner.ScanSource(src)
	if err != nil {
		t.Fatalf("ScanSource failed: %v", err)
	}
	// ConditionalTestLogic は if 文を検出するので、それ以外のスメルがないことを確認
	for _, s := range smells {
		if s.Type != analyzer.SmellConditionalTestLogic {
			t.Errorf("unexpected smell in clean test: %v", s)
		}
	}
}

func TestEnabledSmells(t *testing.T) {
	src := `package foo

import "testing"

func TestEmpty(t *testing.T) {
}
`
	scanner := analyzer.NewScanner()
	scanner.EnabledSmells = map[analyzer.SmellType]bool{
		analyzer.SmellSleepyTest: true,
	}
	smells, err := scanner.ScanSource(src)
	if err != nil {
		t.Fatalf("ScanSource failed: %v", err)
	}
	if len(smells) != 0 {
		t.Errorf("expected no smells when EmptyTest is not enabled, got %v", smells)
	}
}

func TestDisabledSmells(t *testing.T) {
	src := `package foo

import "testing"

func TestEmpty(t *testing.T) {
}
`
	scanner := analyzer.NewScanner()
	scanner.DisabledSmells = map[analyzer.SmellType]bool{
		analyzer.SmellEmptyTest:        true,
		analyzer.SmellMissingAssertion: true,
	}
	smells, err := scanner.ScanSource(src)
	if err != nil {
		t.Fatalf("ScanSource failed: %v", err)
	}
	if len(smells) != 0 {
		t.Errorf("expected no smells when EmptyTest is disabled, got %v", smells)
	}
}
