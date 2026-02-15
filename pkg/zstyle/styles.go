package zstyle

import "github.com/charmbracelet/lipgloss"

// text styles
var (
	Title     = lipgloss.NewStyle().Bold(true).Foreground(Cyan)
	Subtitle  = lipgloss.NewStyle().Bold(true).Foreground(Muted)
	Highlight = lipgloss.NewStyle().Foreground(Orange)
	MutedText = lipgloss.NewStyle().Foreground(Muted)
)

// status indicators
var (
	StatusOK   = lipgloss.NewStyle().Foreground(Success)
	StatusErr  = lipgloss.NewStyle().Foreground(Error)
	StatusWarn = lipgloss.NewStyle().Foreground(Warning)
)

// structural styles
var (
	Border       = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(Subtle)
	ActiveBorder = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(Cyan)
)
