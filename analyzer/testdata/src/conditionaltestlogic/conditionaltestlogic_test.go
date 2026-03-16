package conditionaltestlogic

import "testing"

func TestWithIf(t *testing.T) {
	x := 1
	if x == 1 { // want "テスト内に if 文があります"
		t.Error("one")
	}
}

func TestWithSwitch(t *testing.T) {
	x := 1
	switch x { // want "テスト内に switch 文があります"
	case 1:
		t.Error("one")
	}
}
