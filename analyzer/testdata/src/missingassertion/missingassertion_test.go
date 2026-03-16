package missingassertion

import "testing"

func TestNoAssert(t *testing.T) { // want "テスト関数にアサーションがありません"
	x := 1 + 1
	_ = x
}

func TestWithAssert(t *testing.T) {
	t.Error("ok")
}
