//go:build darwin

package hotkey

// IsWayland checks if we're running on Wayland (always false on macOS)
func IsWayland() bool {
	return false
}
