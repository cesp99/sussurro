package hotkey

import (
	"fmt"
	"log/slog"
	"strings"

	"golang.design/x/hotkey"
)

// Handler manages global hotkeys
type Handler struct {
	hk  *hotkey.Hotkey
	log *slog.Logger
	
	onKeyDown func()
	onKeyUp   func()
}

// NewHandler creates a new hotkey handler
func NewHandler(trigger string, log *slog.Logger) (*Handler, error) {
	mods, key, err := parseTrigger(trigger)
	if err != nil {
		return nil, err
	}

	hk := hotkey.New(mods, key)
	
	return &Handler{
		hk:  hk,
		log: log,
	}, nil
}

// Register registers the hotkey and starts listening
func (h *Handler) Register(onKeyDown, onKeyUp func()) error {
	h.onKeyDown = onKeyDown
	h.onKeyUp = onKeyUp

	if err := h.hk.Register(); err != nil {
		return fmt.Errorf("failed to register hotkey: %w", err)
	}

	// Start listening loop in a goroutine
	go h.listen()

	return nil
}

// Unregister unregisters the hotkey
func (h *Handler) Unregister() {
	h.hk.Unregister()
}

func (h *Handler) listen() {
	for {
		select {
		case <-h.hk.Keydown():
			h.log.Debug("Hotkey pressed")
			if h.onKeyDown != nil {
				h.onKeyDown()
			}
		case <-h.hk.Keyup():
			h.log.Debug("Hotkey released")
			if h.onKeyUp != nil {
				h.onKeyUp()
			}
		}
	}
}

// parseTrigger parses a string like "ctrl+shift+space" into modifiers and key
func parseTrigger(trigger string) ([]hotkey.Modifier, hotkey.Key, error) {
	parts := strings.Split(strings.ToLower(trigger), "+")
	if len(parts) == 0 {
		return nil, 0, fmt.Errorf("empty hotkey trigger")
	}

	var mods []hotkey.Modifier
	var key hotkey.Key

	// Map strings to hotkey constants
	// Note: specific mapping depends on golang.design/x/hotkey definitions
	// We'll implement a basic mapping here
	
	for i, part := range parts {
		// Last part is the key
		if i == len(parts)-1 {
			k, ok := keyMap[part]
			if !ok {
				return nil, 0, fmt.Errorf("unknown key: %s", part)
			}
			key = k
			continue
		}

		// Modifiers
		switch part {
		case "ctrl", "control":
			mods = append(mods, hotkey.ModCtrl)
		case "shift":
			mods = append(mods, hotkey.ModShift)
		case "alt", "option":
			mods = append(mods, hotkey.ModOption)
		case "cmd", "command":
			mods = append(mods, hotkey.ModCmd)
		default:
			return nil, 0, fmt.Errorf("unknown modifier: %s", part)
		}
	}

	return mods, key, nil
}

// Basic key map - expand as needed
var keyMap = map[string]hotkey.Key{
	"space": hotkey.KeySpace,
	"enter": hotkey.KeyReturn,
	"f1":    hotkey.KeyF1,
	"f2":    hotkey.KeyF2,
	"f3":    hotkey.KeyF3,
	"f4":    hotkey.KeyF4,
	"f5":    hotkey.KeyF5,
	"f6":    hotkey.KeyF6,
	"f7":    hotkey.KeyF7,
	"f8":    hotkey.KeyF8,
	"f9":    hotkey.KeyF9,
	"f10":   hotkey.KeyF10,
	"f11":   hotkey.KeyF11,
	"f12":   hotkey.KeyF12,
	"a":     hotkey.KeyA,
	"b":     hotkey.KeyB,
	"c":     hotkey.KeyC,
	"d":     hotkey.KeyD,
	"e":     hotkey.KeyE,
	"f":     hotkey.KeyF,
	"g":     hotkey.KeyG,
	"h":     hotkey.KeyH,
	"i":     hotkey.KeyI,
	"j":     hotkey.KeyJ,
	"k":     hotkey.KeyK,
	"l":     hotkey.KeyL,
	"m":     hotkey.KeyM,
	"n":     hotkey.KeyN,
	"o":     hotkey.KeyO,
	"p":     hotkey.KeyP,
	"q":     hotkey.KeyQ,
	"r":     hotkey.KeyR,
	"s":     hotkey.KeyS,
	"t":     hotkey.KeyT,
	"u":     hotkey.KeyU,
	"v":     hotkey.KeyV,
	"w":     hotkey.KeyW,
	"x":     hotkey.KeyX,
	"y":     hotkey.KeyY,
	"z":     hotkey.KeyZ,
}
