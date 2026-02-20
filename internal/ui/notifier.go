package ui

// AppState represents the current state of the application overlay.
type AppState int

const (
	StateIdle         AppState = iota // 7 animated dots
	StateRecording                    // waveform bars
	StateTranscribing                 // shimmer text
)

// StateNotifier is the interface called by the pipeline to update UI state.
type StateNotifier interface {
	OnStateChange(state AppState)
	OnRMSData(rms float32)
}
