//go:build linux

package hotkey

import (
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/godbus/dbus/v5"
)

// WaylandHandler manages global hotkeys on Wayland via Desktop Portal
type WaylandHandler struct {
	conn      *dbus.Conn
	log       *slog.Logger
	done      chan struct{}
	sessionID string
	shortcutID string

	onKeyDown func()
	onKeyUp   func()

	mu sync.Mutex
	pressed bool
}

// NewWaylandHandler creates a new Wayland hotkey handler
func NewWaylandHandler(trigger string, log *slog.Logger) (*WaylandHandler, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to session bus: %w", err)
	}

	return &WaylandHandler{
		conn: conn,
		log:  log,
		done: make(chan struct{}),
	}, nil
}

// Register registers the hotkey using Desktop Portal
func (h *WaylandHandler) Register(onKeyDown, onKeyUp func()) error {
	h.onKeyDown = onKeyDown
	h.onKeyUp = onKeyUp

	// For now, log that we're on Wayland and suggest alternatives
	h.log.Warn("Wayland detected: Global hotkeys have limited support")
	h.log.Info("Alternative: Use your desktop environment's keyboard settings to bind a command")
	h.log.Info("Alternative: The application will work if you switch to an X11 session")

	// Try to use Desktop Portal if available
	obj := h.conn.Object("org.freedesktop.portal.Desktop", "/org/freedesktop/portal/desktop")

	// Check if GlobalShortcuts portal is available
	var version uint32
	err := obj.Call("org.freedesktop.DBus.Properties.Get", 0,
		"org.freedesktop.portal.GlobalShortcuts", "version").Store(&version)

	if err != nil {
		h.log.Warn("GlobalShortcuts portal not available on this system")
		h.log.Info("You can still use the app by configuring your DE to trigger it")
		return fmt.Errorf("global shortcuts not available on Wayland: use X11 or configure DE shortcuts")
	}

	h.log.Info("GlobalShortcuts portal available", "version", version)
	// TODO: Implement full portal integration
	return fmt.Errorf("GlobalShortcuts portal integration coming soon - please use X11 for now")
}

// Unregister unregisters the hotkey
func (h *WaylandHandler) Unregister() {
	close(h.done)
	if h.conn != nil {
		h.conn.Close()
	}
}

// IsWayland checks if we're running on Wayland
func IsWayland() bool {
	// Check common Wayland environment variables
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		return true
	}
	if os.Getenv("XDG_SESSION_TYPE") == "wayland" {
		return true
	}
	return false
}
