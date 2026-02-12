package asr

import (
	"fmt"
	"os"
	"sync"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
)

// Engine handles the Whisper model and transcription
type Engine struct {
	model   whisper.Model
	context whisper.Context
	mutex   sync.Mutex
}

// NewEngine initializes the Whisper model from a file path
func NewEngine(modelPath string, threads int) (*Engine, error) {
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("model file not found at %s: %w", modelPath, err)
	}

	model, err := whisper.New(modelPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load whisper model: %w", err)
	}

	ctx, err := model.NewContext()
	if err != nil {
		return nil, fmt.Errorf("failed to create whisper context: %w", err)
	}

	return &Engine{
		model:   model,
		context: ctx,
	}, nil
}

// Transcribe processes the audio samples and returns the text
func (e *Engine) Transcribe(samples []float32) (string, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if len(samples) == 0 {
		return "", nil
	}

	if err := e.context.Process(samples, nil, nil, nil); err != nil {
		return "", fmt.Errorf("transcription failed: %w", err)
	}

	// Iterate through segments to build the full text
	var result string
	for {
		segment, err := e.context.NextSegment()
		if err != nil {
			break // End of segments
		}
		result += segment.Text
	}

	return result, nil
}

// Close releases resources
// Note: context.Close() is not available in the bindings, we rely on GC or explicit C-level cleanup if exposed.
// However, Model has Close().
func (e *Engine) Close() {
	// e.context.Close() // Not available in current bindings
	if e.model != nil {
		e.model.Close()
	}
}
