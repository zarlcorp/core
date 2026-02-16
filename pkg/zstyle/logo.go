package zstyle

import "github.com/charmbracelet/lipgloss"

// Logo is the zarlcorp ASCII art wordmark for TUI splash screens.
const Logo = `_  _  _| _ _  _ _ 
/_(_|| |(_(_)| |_)
               |  `

// StyledLogo returns the logo rendered with the given lipgloss style.
func StyledLogo(s lipgloss.Style) string { //nolint:gocritic // value type avoids heap allocation
	return s.Render(Logo)
}
