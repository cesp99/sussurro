//go:build darwin

package ui

// installOverlayHotkey is a no-op on macOS — hotkeys use golang.design/x/hotkey
// via the standard handler in main.go with mainthread.Init.
func installOverlayHotkey(overlay Overlay, trigger string, onDown, onUp func()) {
	_ = overlay
	_ = trigger
	_ = onDown
	_ = onUp
}

// installOverlayContextMenu is a no-op on macOS (NSStatusItem menu handles this).
func installOverlayContextMenu(overlay Overlay, openSettings, quit func()) {
	_ = overlay
	_ = openSettings
	_ = quit
}
