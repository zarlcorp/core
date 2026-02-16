package zstyle

import "github.com/charmbracelet/bubbles/key"

// standard keybindings for zarlcorp TUI applications
var (
	KeyQuit   = key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit"))
	KeyHelp   = key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help"))
	KeyUp     = key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("↑/k", "up"))
	KeyDown   = key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("↓/j", "down"))
	KeyEnter  = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm"))
	KeyBack   = key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back"))
	KeyTab    = key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next"))
	KeyFilter = key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter"))
)
