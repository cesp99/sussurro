//go:build linux

package clipboard

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
)

// Write copies the provided string to the system clipboard
func Write(text string) error {
	err := clipboard.WriteAll(text)
	if err != nil {
		// Provide helpful error message if clipboard tools are missing
		if os.Getenv("WAYLAND_DISPLAY") != "" || os.Getenv("XDG_SESSION_TYPE") == "wayland" {
			return fmt.Errorf("clipboard failed (Wayland requires wl-clipboard): %w - Install with: sudo pacman -S wl-clipboard", err)
		}
		return fmt.Errorf("clipboard failed: %w", err)
	}
	return nil
}

// Read returns the current string content of the system clipboard
func Read() (string, error) {
	text, err := clipboard.ReadAll()
	if err != nil {
		if os.Getenv("WAYLAND_DISPLAY") != "" || os.Getenv("XDG_SESSION_TYPE") == "wayland" {
			return "", fmt.Errorf("clipboard failed (Wayland requires wl-clipboard): %w - Install with: sudo pacman -S wl-clipboard", err)
		}
		return "", fmt.Errorf("clipboard failed: %w", err)
	}
	return text, nil
}
