package zstyle_test

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	"github.com/zarlcorp/core/pkg/zstyle"
)

func TestBaseColors(t *testing.T) {
	colors := []struct {
		name string
		color lipgloss.Color
		hex  string
	}{
		{"Base", zstyle.Base, "#1e1e2e"},
		{"Mantle", zstyle.Mantle, "#181825"},
		{"Crust", zstyle.Crust, "#11111b"},
		{"Surface0", zstyle.Surface0, "#313244"},
		{"Surface1", zstyle.Surface1, "#45475a"},
		{"Surface2", zstyle.Surface2, "#585b70"},
		{"Overlay0", zstyle.Overlay0, "#6c7086"},
		{"Overlay1", zstyle.Overlay1, "#7f849c"},
		{"Overlay2", zstyle.Overlay2, "#9399b2"},
		{"Subtext0", zstyle.Subtext0, "#a6adc8"},
		{"Subtext1", zstyle.Subtext1, "#bac2de"},
		{"Text", zstyle.Text, "#cdd6f4"},
	}

	for _, c := range colors {
		t.Run(c.name, func(t *testing.T) {
			if string(c.color) != c.hex {
				t.Errorf("got %s, want %s", c.color, c.hex)
			}
		})
	}
}

func TestAccentColors(t *testing.T) {
	colors := []struct {
		name  string
		color lipgloss.Color
		hex   string
	}{
		{"Rosewater", zstyle.Rosewater, "#f5e0dc"},
		{"Flamingo", zstyle.Flamingo, "#f2cdcd"},
		{"Pink", zstyle.Pink, "#f5c2e7"},
		{"Mauve", zstyle.Mauve, "#cba6f7"},
		{"Red", zstyle.Red, "#f38ba8"},
		{"Maroon", zstyle.Maroon, "#eba0ac"},
		{"Peach", zstyle.Peach, "#fab387"},
		{"Yellow", zstyle.Yellow, "#f9e2af"},
		{"Green", zstyle.Green, "#a6e3a1"},
		{"Teal", zstyle.Teal, "#94e2d5"},
		{"Sky", zstyle.Sky, "#89dceb"},
		{"Sapphire", zstyle.Sapphire, "#74c7ec"},
		{"Blue", zstyle.Blue, "#89b4fa"},
		{"Lavender", zstyle.Lavender, "#b4befe"},
	}

	for _, c := range colors {
		t.Run(c.name, func(t *testing.T) {
			if string(c.color) != c.hex {
				t.Errorf("got %s, want %s", c.color, c.hex)
			}
		})
	}
}

func TestSemanticColors(t *testing.T) {
	tests := []struct {
		name     string
		semantic lipgloss.Color
		accent   lipgloss.Color
	}{
		{"Success=Green", zstyle.Success, zstyle.Green},
		{"Error=Red", zstyle.Error, zstyle.Red},
		{"Warning=Yellow", zstyle.Warning, zstyle.Yellow},
		{"Info=Blue", zstyle.Info, zstyle.Blue},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.semantic != tt.accent {
				t.Errorf("got %s, want %s", tt.semantic, tt.accent)
			}
		})
	}
}

func TestToolAccents(t *testing.T) {
	tests := []struct {
		name   string
		accent lipgloss.Color
		want   lipgloss.Color
	}{
		{"ZburnAccent=Peach", zstyle.ZburnAccent, zstyle.Peach},
		{"ZvaultAccent=Mauve", zstyle.ZvaultAccent, zstyle.Mauve},
		{"ZshieldAccent=Teal", zstyle.ZshieldAccent, zstyle.Teal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.accent != tt.want {
				t.Errorf("got %s, want %s", tt.accent, tt.want)
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
			got := s.style.Render("test")
			if got == "" {
				t.Errorf("style %s rendered empty string", s.name)
			}
		})
	}
}

func TestCSSVariables(t *testing.T) {
	required := []string{
		"--ctp-base: #1e1e2e",
		"--ctp-mantle: #181825",
		"--ctp-crust: #11111b",
		"--ctp-text: #cdd6f4",
		"--ctp-rosewater: #f5e0dc",
		"--ctp-lavender: #b4befe",
		"--ctp-peach: #fab387",
		"--ctp-green: #a6e3a1",
	}

	for _, want := range required {
		if !strings.Contains(zstyle.CSSVariables, want) {
			t.Errorf("CSSVariables missing %q", want)
		}
	}

	if !strings.HasPrefix(zstyle.CSSVariables, ":root {") {
		t.Error("CSSVariables should start with :root {")
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
