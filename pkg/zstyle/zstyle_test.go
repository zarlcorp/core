package zstyle_test

import (
	"testing"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	"github.com/zarlcorp/core/pkg/zstyle"
)

func TestColors(t *testing.T) {
	colors := []struct {
		name  string
		color lipgloss.Color
	}{
		{"Cyan", zstyle.Cyan},
		{"Orange", zstyle.Orange},
		{"Success", zstyle.Success},
		{"Error", zstyle.Error},
		{"Warning", zstyle.Warning},
		{"Info", zstyle.Info},
		{"Muted", zstyle.Muted},
		{"Subtle", zstyle.Subtle},
		{"Bright", zstyle.Bright},
	}

	for _, c := range colors {
		t.Run(c.name, func(t *testing.T) {
			if c.color == "" {
				t.Errorf("color %s is empty", c.name)
			}
		})
	}
}

func TestStyles(t *testing.T) {
	styles := []struct {
		name  string
		style lipgloss.Style
	}{
		{"Title", zstyle.Title},
		{"Subtitle", zstyle.Subtitle},
		{"Highlight", zstyle.Highlight},
		{"MutedText", zstyle.MutedText},
		{"StatusOK", zstyle.StatusOK},
		{"StatusErr", zstyle.StatusErr},
		{"StatusWarn", zstyle.StatusWarn},
		{"Border", zstyle.Border},
		{"ActiveBorder", zstyle.ActiveBorder},
	}

	for _, s := range styles {
		t.Run(s.name, func(t *testing.T) {
			// verify rendering does not panic
			got := s.style.Render("test")
			if got == "" {
				t.Errorf("style %s rendered empty string", s.name)
			}
		})
	}
}

func TestKeys(t *testing.T) {
	bindings := []struct {
		name    string
		binding key.Binding
	}{
		{"KeyQuit", zstyle.KeyQuit},
		{"KeyHelp", zstyle.KeyHelp},
		{"KeyUp", zstyle.KeyUp},
		{"KeyDown", zstyle.KeyDown},
		{"KeyEnter", zstyle.KeyEnter},
		{"KeyBack", zstyle.KeyBack},
		{"KeyTab", zstyle.KeyTab},
		{"KeyFilter", zstyle.KeyFilter},
	}

	for _, b := range bindings {
		t.Run(b.name, func(t *testing.T) {
			keys := b.binding.Keys()
			if len(keys) == 0 {
				t.Errorf("binding %s has no keys", b.name)
			}

			h := b.binding.Help()
			if h.Key == "" {
				t.Errorf("binding %s has empty help key", b.name)
			}
			if h.Desc == "" {
				t.Errorf("binding %s has empty help description", b.name)
			}
		})
	}
}
