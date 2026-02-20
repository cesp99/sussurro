package audio

import (
	"encoding/binary"
	"fmt"
	"math"
	"sync"

	"github.com/gen2brain/malgo"
)

// CaptureEngine handles audio recording using malgo (miniaudio)
type CaptureEngine struct {
	ctx          *malgo.AllocatedContext
	device       *malgo.Device
	sampleRate   int
	channels     int
	bitDepth     int // Should be 16 for Whisper
	isRecording  bool
	mutex        sync.Mutex
	dataCallback func([]byte)
	rmsCB        func(float32) // optional RMS callback, set via SetRMSCallback
}

// SetRMSCallback installs a callback that receives the RMS level of each
// incoming audio chunk.  The callback is invoked from the audio thread —
// implementations must be non-blocking.
func (e *CaptureEngine) SetRMSCallback(cb func(float32)) {
	e.mutex.Lock()
	e.rmsCB = cb
	e.mutex.Unlock()
}

// computeRMS returns the root-mean-square of a float32 sample slice.
func computeRMS(samples []float32) float32 {
	if len(samples) == 0 {
		return 0
	}
	var sum float64
	for _, s := range samples {
		sum += float64(s) * float64(s)
	}
	return float32(math.Sqrt(sum / float64(len(samples))))
}

// NewCaptureEngine creates a new engine instance
// Whisper typically expects 16kHz, 1 channel, 16-bit PCM
func NewCaptureEngine(sampleRate, channels int) (*CaptureEngine, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to init audio context: %w", err)
	}

	return &CaptureEngine{
		ctx:        ctx,
		sampleRate: sampleRate,
		channels:   channels,
		bitDepth:   16, // Fixed to 16-bit for now
	}, nil
}

// StartRecording starts capturing audio and sends data to the provided channel
func (e *CaptureEngine) StartRecording(dataChan chan<- []float32) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if e.isRecording {
		return nil
	}

	// Define the callback that writes to the channel
	onData := func(data []byte) {
		// Convert byte slice to float32 slice
		// data contains raw F32 samples (4 bytes each)
		numSamples := len(data) / 4
		floats := make([]float32, numSamples)

		for i := 0; i < numSamples; i++ {
			bits := binary.LittleEndian.Uint32(data[i*4 : (i+1)*4])
			floats[i] = math.Float32frombits(bits)
		}

		// Invoke RMS callback (non-blocking) if installed
		e.mutex.Lock()
		cb := e.rmsCB
		e.mutex.Unlock()
		if cb != nil {
			rms := computeRMS(floats)
			cb(rms)
		}

		// Non-blocking send
		select {
		case dataChan <- floats:
		default:
			// Drop frame if buffer is full
		}
	}

	// Start the internal device
	return e.startDevice(onData)
}

// startDevice initiates the low-level audio stream
func (e *CaptureEngine) startDevice(onData func([]byte)) error {
	e.dataCallback = onData

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = malgo.FormatF32
	deviceConfig.Capture.Channels = uint32(e.channels)
	deviceConfig.SampleRate = uint32(e.sampleRate)
	deviceConfig.Alsa.NoMMap = 1 // Common fix for Linux ALSA

	var err error
	// Callback to handle incoming audio data
	e.device, err = malgo.InitDevice(e.ctx.Context, deviceConfig, malgo.DeviceCallbacks{
		Data: func(pOutputSample, pInputSamples []byte, framecount uint32) {
			if e.dataCallback != nil {
				// We received F32 samples as bytes
				// Copy them to ensure memory safety
				if len(pInputSamples) == 0 {
					return
				}

				dataCopy := make([]byte, len(pInputSamples))
				copy(dataCopy, pInputSamples)

				e.dataCallback(dataCopy)
			}
		},
	})
	if err != nil {
		return fmt.Errorf("failed to init device: %w", err)
	}

	err = e.device.Start()
	if err != nil {
		return fmt.Errorf("failed to start device: %w", err)
	}

	e.isRecording = true
	return nil
}

// Stop halts the stream
func (e *CaptureEngine) Stop() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if !e.isRecording {
		return nil
	}

	if e.device != nil {
		e.device.Uninit()
		e.device = nil
	}
	e.isRecording = false
	return nil
}

// Close releases resources
func (e *CaptureEngine) Close() {
	e.Stop()
	if e.ctx != nil {
		e.ctx.Free()
		e.ctx = nil
	}
}
