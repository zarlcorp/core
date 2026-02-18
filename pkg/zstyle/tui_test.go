package zstyle_test

import (
	"strings"
	"testing"

	"github.com/zarlcorp/core/pkg/zstyle"
)

func TestRenderFooter(t *testing.T) {
	t.Run("multiple pairs", func(t *testing.T) {
		pairs := []zstyle.HelpPair{
			{Key: "q", Desc: "quit"},
			{Key: "?", Desc: "help"},
			{Key: "enter", Desc: "confirm"},
		}
		got := zstyle.RenderFooter(pairs)

		if !strings.Contains(got, "q") {
			t.Error("output missing key 'q'")
		}
		if !strings.Contains(got, "quit") {
			t.Error("output missing desc 'quit'")
		}
		if !strings.Contains(got, "?") {
			t.Error("output missing key '?'")
		}
		if !strings.Contains(got, "help") {
			t.Error("output missing desc 'help'")
		}
		if !strings.Contains(got, "enter") {
			t.Error("output missing key 'enter'")
		}
		if !strings.Contains(got, "confirm") {
			t.Error("output missing desc 'confirm'")
		}
		if !strings.Contains(got, "|") {
			t.Error("output missing separator '|'")
		}
	})

	t.Run("single pair", func(t *testing.T) {
		pairs := []zstyle.HelpPair{{Key: "q", Desc: "quit"}}
		got := zstyle.RenderFooter(pairs)

		if !strings.Contains(got, "q") {
			t.Error("output missing key 'q'")
		}
		if !strings.Contains(got, "quit") {
			t.Error("output missing desc 'quit'")
		}
		// single pair should not contain a separator
		if strings.Contains(got, "|") {
			t.Error("single pair should not have separator")
		}
	})

	t.Run("empty pairs", func(t *testing.T) {
		got := zstyle.RenderFooter(nil)
		if got != "" {
			t.Errorf("expected empty string for nil pairs, got %q", got)
		}

		got = zstyle.RenderFooter([]zstyle.HelpPair{})
		if got != "" {
			t.Errorf("expected empty string for empty pairs, got %q", got)
		}
	})
}

func TestRenderMenuItem(t *testing.T) {
	t.Run("active with count", func(t *testing.T) {
		item := zstyle.MenuItem{Label: "Secrets", Count: "(3)", Active: true}
		got := zstyle.RenderMenuItem(item, zstyle.Mauve)

		if !strings.Contains(got, "▸") {
			t.Error("active item missing cursor ▸")
		}
		if !strings.Contains(got, "Secrets") {
			t.Error("active item missing label")
		}
		if !strings.Contains(got, "(3)") {
			t.Error("active item missing count")
		}
	})

	t.Run("inactive with count", func(t *testing.T) {
		item := zstyle.MenuItem{Label: "Tasks", Count: "(2 pending)", Active: false}
		got := zstyle.RenderMenuItem(item, zstyle.Mauve)

		if strings.Contains(got, "▸") {
			t.Error("inactive item should not have cursor ▸")
		}
		if !strings.Contains(got, "Tasks") {
			t.Error("inactive item missing label")
		}
		if !strings.Contains(got, "(2 pending)") {
			t.Error("inactive item missing count")
		}
	})

	t.Run("active without count", func(t *testing.T) {
		item := zstyle.MenuItem{Label: "Settings", Active: true}
		got := zstyle.RenderMenuItem(item, zstyle.Peach)

		if !strings.Contains(got, "▸") {
			t.Error("active item missing cursor ▸")
		}
		if !strings.Contains(got, "Settings") {
			t.Error("active item missing label")
		}
	})

	t.Run("inactive without count", func(t *testing.T) {
		item := zstyle.MenuItem{Label: "About", Active: false}
		got := zstyle.RenderMenuItem(item, zstyle.Teal)

		if strings.Contains(got, "▸") {
			t.Error("inactive item should not have cursor")
		}
		if !strings.Contains(got, "About") {
			t.Error("inactive item missing label")
		}
	})
}

func TestRenderHeader(t *testing.T) {
	t.Run("with view title", func(t *testing.T) {
		got := zstyle.RenderHeader("zvault", "Secrets", zstyle.Mauve)

		if !strings.Contains(got, "zvault") {
			t.Error("header missing app name")
		}
		if !strings.Contains(got, "/") {
			t.Error("header missing separator '/'")
		}
		if !strings.Contains(got, "Secrets") {
			t.Error("header missing view title")
		}
	})

	t.Run("without view title", func(t *testing.T) {
		got := zstyle.RenderHeader("zburn", "", zstyle.Peach)

		if !strings.Contains(got, "zburn") {
			t.Error("header missing app name")
		}
		if strings.Contains(got, "/") {
			t.Error("header without view title should not have separator")
		}
	})
}

func TestRenderSeparator(t *testing.T) {
	t.Run("positive width", func(t *testing.T) {
		got := zstyle.RenderSeparator(40)

		if !strings.Contains(got, "─") {
			t.Error("separator missing ─ character")
		}
		// count ─ runes in the output (ignoring ANSI escapes)
		count := strings.Count(got, "─")
		if count != 40 {
			t.Errorf("expected 40 ─ characters, got %d", count)
		}
	})

	t.Run("zero width defaults to 60", func(t *testing.T) {
		got := zstyle.RenderSeparator(0)

		count := strings.Count(got, "─")
		if count != 60 {
			t.Errorf("expected 60 ─ characters for zero width, got %d", count)
		}
	})

	t.Run("negative width defaults to 60", func(t *testing.T) {
		got := zstyle.RenderSeparator(-5)

		count := strings.Count(got, "─")
		if count != 60 {
			t.Errorf("expected 60 ─ characters for negative width, got %d", count)
		}
	})
}
