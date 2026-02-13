package audio

import (
	"encoding/binary"
	"math"
	"sync"
)

// StreamBuffer manages a thread-safe buffer for audio samples
// It converts raw bytes to float32 samples (required by Whisper.cpp)
type StreamBuffer struct {
	mu      sync.Mutex
	samples []float32
}

func NewStreamBuffer() *StreamBuffer {
	return &StreamBuffer{
		samples: make([]float32, 0, 16000*2), // Pre-allocate for ~2 seconds
	}
}

// Write appends raw 16-bit PCM bytes to the buffer
func (b *StreamBuffer) Write(data []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Convert 16-bit PCM bytes to float32 [-1.0, 1.0]
	// Whisper expects float32
	numSamples := len(data) / 2
	for i := 0; i < numSamples; i++ {
		// Read int16 (Little Endian)
		sampleInt16 := int16(binary.LittleEndian.Uint16(data[i*2 : i*2+2]))
		// Normalize to float32
		sampleFloat := float32(sampleInt16) / 32768.0
		b.samples = append(b.samples, sampleFloat)
	}
}

// Read returns all available samples and clears the buffer
func (b *StreamBuffer) Read() []float32 {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.samples) == 0 {
		return nil
	}

	// Copy and clear
	result := make([]float32, len(b.samples))
	copy(result, b.samples)
	b.samples = b.samples[:0] // Reset slice but keep capacity

	return result
}

// CalculateRMS calculates the Root Mean Square amplitude of the current chunk
// Useful for VAD (Voice Activity Detection)
func CalculateRMS(samples []float32) float32 {
	if len(samples) == 0 {
		return 0
	}
	var sum float32
	for _, s := range samples {
		sum += s * s
	}
	return float32(math.Sqrt(float64(sum) / float64(len(samples))))
}
