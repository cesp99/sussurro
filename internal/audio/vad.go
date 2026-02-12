package audio

import (
	"math"
)

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

// ComputeRMS calculates the Root Mean Square of the audio signal
func ComputeRMS(pcm []float32) float32 {
	if len(pcm) == 0 {
		return 0
	}
	var sum float64
	for _, sample := range pcm {
		sum += float64(sample * sample)
	}
	return float32(math.Sqrt(sum / float64(len(pcm))))
}

// IsSpeechSimple checks if the buffer RMS is above a threshold
func IsSpeechSimple(pcm []float32, threshold float32) bool {
	return ComputeRMS(pcm) > threshold
}
