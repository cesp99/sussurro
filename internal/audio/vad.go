package audio

// VADParams holds configuration for Voice Activity Detection
type VADParams struct {
	SampleRate    int
	EnergyThresh  float32 // Threshold for RMS energy (e.g., 0.005 for quiet environments)
	SilenceThresh float32 // Threshold to consider silence (lower than EnergyThresh)
}

// DefaultVADParams returns sensible defaults
func DefaultVADParams() VADParams {
	return VADParams{
		SampleRate:    16000,
		EnergyThresh:  0.01,
		SilenceThresh: 0.002,
	}
}
