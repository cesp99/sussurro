//go:build linux

package hotkey

import (
	"fmt"
	"golang.design/x/hotkey"
)

// parseModifier parses a modifier string to hotkey.Modifier for Linux
func parseModifier(part string) (hotkey.Modifier, error) {
	switch part {
	case "ctrl", "control":
		return hotkey.ModCtrl, nil
	case "shift":
		return hotkey.ModShift, nil
	case "alt", "option":
		// Note: Alt modifier support may vary on Linux
		return 0, fmt.Errorf("alt modifier not fully supported on Linux, use ctrl+shift instead")
	case "cmd", "command", "super", "meta":
		// Note: Cmd/Super modifier support may vary on Linux
		return 0, fmt.Errorf("cmd/super modifier not fully supported on Linux, use ctrl+shift instead")
	default:
		return 0, fmt.Errorf("unknown modifier: %s", part)
	}
}
