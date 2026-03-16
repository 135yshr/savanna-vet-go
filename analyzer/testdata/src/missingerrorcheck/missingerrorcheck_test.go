package missingerrorcheck

import (
	"os"
	"testing"
)

func TestIgnoreErr(t *testing.T) {
	f, _ := os.Open("test.txt") // want "error の戻り値が無視されています"
	_ = f
	t.Error("dummy")
}

func TestNoIgnore(t *testing.T) {
	f, err := os.Open("test.txt")
	if err != nil {
		t.Fatal(err)
	}
	_ = f
}

func TestMapLookup(t *testing.T) {
	m := map[string]int{"a": 1}
	v, _ := m["a"] // map lookup: last value is bool, not error — should NOT be reported
	_ = v
	t.Error("dummy")
}
