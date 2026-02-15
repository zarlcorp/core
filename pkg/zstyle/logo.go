package zstyle

import "github.com/charmbracelet/lipgloss"

// Logo is the zarlcorp ASCII art wordmark for TUI splash screens.
// ribbon/origami letterforms with diagonal box-drawing characters.
const Logo = "" +
	"   ________  ________  ________  ______    ________  ________  ________  ________ \n" +
	"  ╱        ╲╱        ╲╱        ╲╱      ╲  ╱        ╲╱        ╲╱        ╲╱        ╲\n" +
	" ╱-        ╱    /    ╱     /   ╱       ╱ ╱         ╱    /    ╱    /    ╱    /    ╱\n" +
	"╱        _╱         ╱        _╱       ╱_╱       --╱         ╱        _╱       __╱ \n" +
	"╲________╱╲___╱____╱╲____╱___╱╲________╱╲________╱╲________╱╲____╱___╱╲______╱    "

// StyledLogo returns the logo rendered with the given lipgloss style.
func StyledLogo(s *lipgloss.Style) string {
	return s.Render(Logo)
}
