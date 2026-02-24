package ui

// Overlay is the platform-independent interface for the capsule overlay window.
type Overlay interface {
	Show()
	Hide()
	SetState(state AppState)
	PushRMS(rms float32)
	Close()
}
