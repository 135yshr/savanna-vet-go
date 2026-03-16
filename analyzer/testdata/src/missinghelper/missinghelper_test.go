package missinghelper

import "testing"

func setupTest(t *testing.T) { // want "テストヘルパー関数に t\\.Helper\\(\\) がありません"
	// no t.Helper() call
}

func setupTestWithHelper(t *testing.T) {
	t.Helper()
}
