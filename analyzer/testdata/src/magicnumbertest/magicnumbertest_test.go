package magicnumbertest

import "testing"

func TestMagic(t *testing.T) {
	t.Errorf("expected %d", 42) // want "マジックナンバー 42 には名前付き定数を使用してください"
}

func TestAllowed(t *testing.T) {
	t.Errorf("expected %d", 0) // allowed number, no diagnostic
}
