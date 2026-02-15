package zstyle

import "github.com/charmbracelet/lipgloss"

// Logo is the zarlcorp ASCII art wordmark for TUI splash screens.
// lowercase letterforms with box-drawing characters.
const Logo = "" +
	"──┐ ┌─┐ ┌─┐ │   ┌── ┌─┐ ┌─┐ ┌─┐\n" +
	"┌─┘ ├─┤ ├─┘ │   │   │ │ ├─┘ ├─┘\n" +
	"└── ┴ ┴ ┴   └── └── └─┘ ┴   │  "

// StyledLogo returns the logo rendered with the given lipgloss style.
func StyledLogo(s lipgloss.Style) string {
	return s.Render(Logo)
}
