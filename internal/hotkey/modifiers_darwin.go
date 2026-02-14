//go:build darwin

package hotkey

import (
	"fmt"
	"golang.design/x/hotkey"
)

// parseModifier parses a modifier string to hotkey.Modifier for macOS
func parseModifier(part string) (hotkey.Modifier, error) {
	switch part {
	case "ctrl", "control":
		return hotkey.ModCtrl, nil
	case "shift":
		return hotkey.ModShift, nil
	case "alt", "option":
		return hotkey.ModOption, nil
	case "cmd", "command":
		return hotkey.ModCmd, nil
	default:
		return 0, fmt.Errorf("unknown modifier: %s", part)
	}
}
