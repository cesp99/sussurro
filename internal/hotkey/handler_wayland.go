//go:build linux

package hotkey

import "os"

// IsWayland checks if we're running on Wayland
func IsWayland() bool {
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		return true
	}
	if os.Getenv("XDG_SESSION_TYPE") == "wayland" {
		return true
	}
	return false
}
