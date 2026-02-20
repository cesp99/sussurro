//go:build linux

package ui

// installOverlayHotkey registers an X11 global hotkey via GDK XGrabKey.
// On Wayland, the overlay is a *linuxOverlay but IsWayland() returns true,
// so the caller should skip this and use the trigger server instead.
func installOverlayHotkey(overlay Overlay, trigger string, onDown, onUp func()) {
	if lo, ok := overlay.(*linuxOverlay); ok {
		lo.installHotkey(trigger, onDown, onUp)
	}
}

// installOverlayContextMenu wires the right-click menu on the GTK3 overlay.
func installOverlayContextMenu(overlay Overlay, openSettings, quit func()) {
	if lo, ok := overlay.(*linuxOverlay); ok {
		lo.installContextMenu(openSettings, quit)
	}
}
