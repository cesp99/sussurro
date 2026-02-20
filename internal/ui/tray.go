package ui

import (
	_ "embed"

	"github.com/getlantern/systray"
)

//go:embed assets/tray.png
var trayIcon []byte

//go:embed assets/tray_rec.png
var trayIconRec []byte

// runTray starts the system tray in the calling goroutine (blocks).
// It must be started with go m.runTray() so it doesn't block the UI thread.
func (m *Manager) runTray() {
	systray.Run(m.onTrayReady, m.onTrayExit)
}

func (m *Manager) onTrayReady() {
	systray.SetIcon(trayIcon)
	systray.SetTooltip("Sussurro")

	mSettings := systray.AddMenuItem("Open Settings", "Open the settings window")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Exit Sussurro")

	go func() {
		for {
			select {
			case <-mSettings.ClickedCh:
				m.settings.Show()

			case <-mQuit.ClickedCh:
				m.Quit()
				return
			}
		}
	}()
}

func (m *Manager) onTrayExit() {}

// updateTrayIcon swaps the tray icon based on recording state.
func (m *Manager) updateTrayIcon(state AppState) {
	if state == StateRecording {
		systray.SetIcon(trayIconRec)
	} else {
		systray.SetIcon(trayIcon)
	}
}
