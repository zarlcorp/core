// Package zclipboard provides cross-platform system clipboard access.
package zclipboard

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// Copy writes text to the system clipboard.
func Copy(text string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		if _, err := exec.LookPath("xclip"); err == nil {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		} else if _, err := exec.LookPath("xsel"); err == nil {
			cmd = exec.Command("xsel", "--clipboard", "--input")
		} else {
			return fmt.Errorf("no clipboard tool: install xclip or xsel")
		}
	case "windows":
		cmd = exec.Command("cmd", "/c", "clip")
	default:
		return fmt.Errorf("clipboard not supported on %s", runtime.GOOS)
	}

	cmd.Stdin = strings.NewReader(text)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clipboard: %w", err)
	}

	return nil
}

// Clear empties the system clipboard.
func Clear() error {
	return Copy("")
}
