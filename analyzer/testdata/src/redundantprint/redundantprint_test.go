package redundantprint

import (
	"fmt"
	"testing"
)

func TestWithPrint(t *testing.T) {
	fmt.Println("debug") // want "fmt\\.Println\\(\\) の代わりに t\\.Log\\(\\) を使用してください"
	t.Error("dummy")
}

func TestNoPrint(t *testing.T) {
	t.Error("ok")
}
