package zstyle

import "github.com/charmbracelet/lipgloss"

// text styles
var (
	Title     = lipgloss.NewStyle().Bold(true).Foreground(Lavender)
	Subtitle  = lipgloss.NewStyle().Bold(true).Foreground(Subtext1)
	Highlight = lipgloss.NewStyle().Foreground(Peach)
	MutedText = lipgloss.NewStyle().Foreground(Overlay1)
)

// status indicators
var (
	StatusOK   = lipgloss.NewStyle().Foreground(Success)
	StatusErr  = lipgloss.NewStyle().Foreground(Error)
	StatusWarn = lipgloss.NewStyle().Foreground(Warning)
)

// structural styles
var (
	Border       = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(Surface1)
	ActiveBorder = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(Lavender)
)
