package zstyle_test

import (
	"strings"
	"testing"

	"github.com/zarlcorp/core/pkg/zstyle"
)

func TestLogo(t *testing.T) {
	t.Run("not empty", func(t *testing.T) {
		if zstyle.Logo == "" {
			t.Error("Logo is empty")
		}
	})

	t.Run("three lines", func(t *testing.T) {
		lines := strings.Split(zstyle.Logo, "\n")
		if len(lines) != 3 {
			t.Errorf("Logo has %d lines, want 3", len(lines))
		}
	})

	t.Run("fits 80 columns", func(t *testing.T) {
		for i, line := range strings.Split(zstyle.Logo, "\n") {
			// count runes, not bytes — logo uses multibyte box-drawing chars
			n := 0
			for range line {
				n++
			}
			if n > 80 {
				t.Errorf("line %d is %d runes wide, want <= 80", i, n)
			}
		}
	})
}

func TestStyledLogo(t *testing.T) {
	got := zstyle.StyledLogo(zstyle.Title)
	if got == "" {
		t.Error("StyledLogo returned empty string")
	}
	// lipgloss may or may not emit ANSI in test environments,
	// so we just verify the logo text survives styling
	if !strings.Contains(got, "┌") {
		t.Error("StyledLogo output missing logo content")
	}
}
