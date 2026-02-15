// Package zstyle provides the zarlcorp visual identity for TUI applications.
//
// It defines a consistent color palette, lipgloss style presets, and standard
// keybinding constants. Every zarlcorp TUI tool imports zstyle so they share
// the same look and feel.
//
// The palette is Catppuccin Mocha — warm pastels on dark backgrounds, designed
// for the riced/unixporn aesthetic: rounded borders, floating containers,
// monospaced everything.
//
// # Usage
//
//	import "github.com/zarlcorp/core/pkg/zstyle"
//
//	fmt.Println(zstyle.Title.Render("My Tool"))
//	fmt.Println(zstyle.StatusOK.Render("done"))
package zstyle

import "github.com/charmbracelet/lipgloss"

// catppuccin mocha base colors
var (
	Base     = lipgloss.Color("#1e1e2e") // primary background
	Mantle   = lipgloss.Color("#181825") // darker background
	Crust    = lipgloss.Color("#11111b") // darkest background
	Surface0 = lipgloss.Color("#313244") // elevated surface
	Surface1 = lipgloss.Color("#45475a") // borders, separators
	Surface2 = lipgloss.Color("#585b70") // inactive elements
	Overlay0 = lipgloss.Color("#6c7086") // muted text
	Overlay1 = lipgloss.Color("#7f849c") // secondary text
	Overlay2 = lipgloss.Color("#9399b2")
	Subtext0 = lipgloss.Color("#a6adc8")
	Subtext1 = lipgloss.Color("#bac2de")
	Text     = lipgloss.Color("#cdd6f4") // primary text
)

// catppuccin mocha accent colors
var (
	Rosewater = lipgloss.Color("#f5e0dc")
	Flamingo  = lipgloss.Color("#f2cdcd")
	Pink      = lipgloss.Color("#f5c2e7")
	Mauve     = lipgloss.Color("#cba6f7")
	Red       = lipgloss.Color("#f38ba8")
	Maroon    = lipgloss.Color("#eba0ac")
	Peach     = lipgloss.Color("#fab387")
	Yellow    = lipgloss.Color("#f9e2af")
	Green     = lipgloss.Color("#a6e3a1")
	Teal      = lipgloss.Color("#94e2d5")
	Sky       = lipgloss.Color("#89dceb")
	Sapphire  = lipgloss.Color("#74c7ec")
	Blue      = lipgloss.Color("#89b4fa")
	Lavender  = lipgloss.Color("#b4befe")
)

// semantic colors mapped to catppuccin accents
var (
	Success = Green
	Error   = Red
	Warning = Yellow
	Info    = Blue
)

// per-tool accent colors
var (
	ZburnAccent   = Peach    // warm — fire, burning
	ZvaultAccent  = Mauve    // regal — vaults, secrets
	ZshieldAccent = Teal     // protective — shields, safety
)

// CSSVariables exports the full catppuccin mocha palette as CSS custom properties.
const CSSVariables = `:root {
  --ctp-base: #1e1e2e;
  --ctp-mantle: #181825;
  --ctp-crust: #11111b;
  --ctp-surface0: #313244;
  --ctp-surface1: #45475a;
  --ctp-surface2: #585b70;
  --ctp-overlay0: #6c7086;
  --ctp-overlay1: #7f849c;
  --ctp-overlay2: #9399b2;
  --ctp-subtext0: #a6adc8;
  --ctp-subtext1: #bac2de;
  --ctp-text: #cdd6f4;
  --ctp-rosewater: #f5e0dc;
  --ctp-flamingo: #f2cdcd;
  --ctp-pink: #f5c2e7;
  --ctp-mauve: #cba6f7;
  --ctp-red: #f38ba8;
  --ctp-maroon: #eba0ac;
  --ctp-peach: #fab387;
  --ctp-yellow: #f9e2af;
  --ctp-green: #a6e3a1;
  --ctp-teal: #94e2d5;
  --ctp-sky: #89dceb;
  --ctp-sapphire: #74c7ec;
  --ctp-blue: #89b4fa;
  --ctp-lavender: #b4befe;
}
`
