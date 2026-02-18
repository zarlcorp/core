package zstyle

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// HelpPair is a keybinding label and its description, used for footer help text.
type HelpPair struct {
	Key  string
	Desc string
}

// MenuItem describes a single entry in a TUI menu list.
type MenuItem struct {
	Label  string
	Count  string // e.g. "(3)", "(2 pending)", or ""
	Active bool
}

// RenderFooter renders key/desc pairs as a single-line help bar.
// Output format: "  key desc | key desc | key desc"
func RenderFooter(pairs []HelpPair) string {
	if len(pairs) == 0 {
		return ""
	}

	keyStyle := lipgloss.NewStyle().Foreground(Lavender).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(Overlay1)
	sepStyle := lipgloss.NewStyle().Foreground(Surface2)

	sep := sepStyle.Render(" | ")

	parts := make([]string, len(pairs))
	for i, p := range pairs {
		parts[i] = keyStyle.Render(p.Key) + " " + descStyle.Render(p.Desc)
	}

	return "  " + strings.Join(parts, sep)
}

// RenderMenuItem renders a single menu item with an optional active cursor.
// Active items show "  ▸ Label" in the accent color; inactive show "    Label" in Text.
func RenderMenuItem(item MenuItem, accent lipgloss.Color) string {
	accentStyle := lipgloss.NewStyle().Foreground(accent).Bold(true)
	textStyle := lipgloss.NewStyle().Foreground(Text)
	countStyle := lipgloss.NewStyle().Foreground(Overlay1)

	var b strings.Builder

	if item.Active {
		b.WriteString("  ")
		b.WriteString(accentStyle.Render("▸"))
		b.WriteString(" ")
		b.WriteString(accentStyle.Render(item.Label))
	} else {
		b.WriteString("    ")
		b.WriteString(textStyle.Render(item.Label))
	}

	if item.Count != "" {
		b.WriteString(" ")
		b.WriteString(countStyle.Render(item.Count))
	}

	return b.String()
}

// RenderHeader renders a breadcrumb-style header: "  appName / viewTitle".
// If viewTitle is empty, only the app name is rendered.
func RenderHeader(appName, viewTitle string, accent lipgloss.Color) string {
	nameStyle := lipgloss.NewStyle().Foreground(accent).Bold(true)
	sepStyle := lipgloss.NewStyle().Foreground(Overlay1)
	viewStyle := lipgloss.NewStyle().Foreground(Subtext1)

	var b strings.Builder
	b.WriteString("  ")
	b.WriteString(nameStyle.Render(appName))

	if viewTitle != "" {
		b.WriteString(sepStyle.Render(" / "))
		b.WriteString(viewStyle.Render(viewTitle))
	}

	return b.String()
}

// RenderSeparator renders a horizontal line of ─ characters styled in Surface1.
// If width is zero or negative, it defaults to 60.
func RenderSeparator(width int) string {
	if width <= 0 {
		width = 60
	}

	line := strings.Repeat("─", width)
	return lipgloss.NewStyle().Foreground(Surface1).Render(line)
}
