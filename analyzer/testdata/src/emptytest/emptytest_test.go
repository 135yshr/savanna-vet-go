package emptytest

import "testing"

func TestEmpty(t *testing.T) { // want "テスト関数の本体が空です"
}

func TestNotEmpty(t *testing.T) {
	t.Error("ok")
}
