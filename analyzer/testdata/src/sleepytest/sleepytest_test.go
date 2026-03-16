package sleepytest

import (
	"testing"
	"time"
)

func TestWithSleep(t *testing.T) {
	time.Sleep(time.Second) // want "time\\.Sleep\\(\\) はテストを不安定にします"
	t.Error("dummy")
}

func TestNoSleep(t *testing.T) {
	t.Error("ok")
}
