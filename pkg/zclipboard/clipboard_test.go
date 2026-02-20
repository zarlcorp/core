package zclipboard

import (
	"runtime"
	"testing"
)

func TestCopyAndClear(t *testing.T) {
	if runtime.GOOS != "darwin" && runtime.GOOS != "linux" {
		t.Skip("clipboard tests require darwin or linux")
	}

	if err := Copy("zclipboard-test"); err != nil {
		t.Fatalf("Copy: %v", err)
	}

	if err := Clear(); err != nil {
		t.Fatalf("Clear: %v", err)
	}
}
