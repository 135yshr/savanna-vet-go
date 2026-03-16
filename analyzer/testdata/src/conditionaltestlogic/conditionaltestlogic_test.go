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

func TestWithNestedIf(t *testing.T) {
	x := 1
	if x == 1 { // want "テスト内に if 文があります"
		if x > 0 { // want "テスト内に if 文があります"
			t.Error("nested")
		}
	}
}

func TestTableDriven(t *testing.T) {
	tests := []struct {
		name string
		val  int
	}{
		{"one", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.val != 1 { // want "テスト内に if 文があります"
				t.Error("fail")
			}
		})
	}
}
