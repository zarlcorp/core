// Package zstyle provides the zarlcorp visual identity for TUI applications.
//
// It defines a consistent color palette, lipgloss style presets, and standard
// keybinding constants. Every zarlcorp TUI tool imports zstyle so they share
// the same look and feel.
//
// zstyle assumes a dark terminal background and never sets background colors
// on main content areas.
//
// # Usage
//
//	import "github.com/zarlcorp/core/pkg/zstyle"
//
//	fmt.Println(zstyle.Title.Render("My Tool"))
//	fmt.Println(zstyle.StatusOK.Render("✓ done"))
package zstyle

import "github.com/charmbracelet/lipgloss"

// core brand colors
var (
	Cyan   = lipgloss.Color("#00E5FF") // primary — headings, active elements, highlights
	Orange = lipgloss.Color("#FF6E40") // accent — calls to action, selected items, emphasis
)

// semantic colors
var (
	Success = lipgloss.Color("#69F0AE") // green — completed, passed, ok
	Error   = lipgloss.Color("#FF5252") // red — errors, failures, destructive actions
	Warning = lipgloss.Color("#FFD740") // amber — warnings, caution
	Info    = lipgloss.Color("#40C4FF") // light blue — informational
)

// neutral tones
var (
	Muted  = lipgloss.Color("#78909C") // grey-blue — secondary text, timestamps, metadata
	Subtle = lipgloss.Color("#37474F") // dark grey — borders, separators, inactive elements
	Bright = lipgloss.Color("#ECEFF1") // near-white — primary text when emphasis needed
)
