package audio

import (
	"fmt"
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

// Start initiates the audio stream
// onData is called with raw PCM bytes
func (e *CaptureEngine) Start(onData func([]byte)) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if e.isRecording {
		return fmt.Errorf("already recording")
	}

	e.dataCallback = onData

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = uint32(e.channels)
	deviceConfig.SampleRate = uint32(e.sampleRate)
	deviceConfig.Alsa.NoMMap = 1 // Common fix for Linux ALSA

	var err error
	// Callback to handle incoming audio data
	e.device, err = malgo.InitDevice(e.ctx.Context, deviceConfig, malgo.DeviceCallbacks{
		Data: func(pOutputSample, pInputSamples []byte, framecount uint32) {
			if e.dataCallback != nil {
				// Create a copy of the data to be safe, or pass it directly if the consumer handles it quickly.
				// For now, let's pass it. The slice is valid only for this callback.
				// To be safe for async processing, we should copy it.
				// Whisper processing is slow, so we definitely need to buffer this elsewhere.
				// The CaptureEngine just hands off the raw buffer.
				
				// Make a copy because pInputSamples is reused by miniaudio
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
