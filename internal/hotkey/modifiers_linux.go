//go:build linux

package hotkey

import (
	"fmt"

	"golang.design/x/hotkey"
)

// parseModifier parses a modifier string to hotkey.Modifier for Linux (X11).
// On X11: Mod1 is conventionally Alt, Mod4 is conventionally Super/Windows key.
func parseModifier(part string) (hotkey.Modifier, error) {
	switch part {
	case "ctrl", "control":
		return hotkey.ModCtrl, nil
	case "shift":
		return hotkey.ModShift, nil
	case "alt", "option":
		// Mod1 = Alt_L on virtually all X11 systems
		return hotkey.Mod1, nil
	case "cmd", "command", "super", "meta":
		// Mod4 = Super_L (Windows/Meta key) on virtually all X11 systems
		return hotkey.Mod4, nil
	default:
		return 0, fmt.Errorf("unknown modifier: %s", part)
	}
}
