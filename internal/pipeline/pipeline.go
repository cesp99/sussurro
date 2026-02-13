package pipeline

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/cesp99/sussurro/internal/asr"
	"github.com/cesp99/sussurro/internal/audio"
	"github.com/cesp99/sussurro/internal/clipboard"
	ctxProvider "github.com/cesp99/sussurro/internal/context"
	"github.com/cesp99/sussurro/internal/injection"
	"github.com/cesp99/sussurro/internal/llm"
)

// Pipeline orchestrates the flow of data from audio capture to text output
type Pipeline struct {
	audioEngine *audio.CaptureEngine
	asrEngine   *asr.Engine
	llmEngine   *llm.Engine
	ctxProvider ctxProvider.Provider
	injector    *injection.Injector
	log         *slog.Logger
	vadParams   audio.VADParams

	onCompletion func() // Callback for when processing finishes

	// Channels for data flow
	audioChan chan []float32
	stopChan  chan struct{}
	wg        sync.WaitGroup

	// State
	isRecording bool
	audioBuffer []float32
	mu          sync.Mutex // Protects isRecording and audioBuffer
}

// NewPipeline creates a new processing pipeline
func NewPipeline(
	audioEngine *audio.CaptureEngine,
	asrEngine *asr.Engine,
	llmEngine *llm.Engine,
	ctxProvider ctxProvider.Provider,
	injector *injection.Injector,
	log *slog.Logger,
	sampleRate int,
) *Pipeline {
	vadParams := audio.DefaultVADParams()
	vadParams.SampleRate = sampleRate // Override with actual sample rate

	return &Pipeline{
		audioEngine: audioEngine,
		asrEngine:   asrEngine,
		llmEngine:   llmEngine,
		ctxProvider: ctxProvider,
		injector:    injector,
		log:         log,
		vadParams:   vadParams,
		audioChan:   make(chan []float32, 100), // Buffer audio chunks
		stopChan:    make(chan struct{}),
	}
}

// SetOnCompletion sets a callback to be called when processing is done
func (p *Pipeline) SetOnCompletion(callback func()) {
	p.onCompletion = callback
}

// Start begins the pipeline processing
func (p *Pipeline) Start() error {
	p.log.Info("Starting pipeline...")

	// Start Audio Capture Loop (runs continuously to keep device ready)
	p.wg.Add(1)
	go p.captureLoop()

	return nil
}

// Stop gracefully shuts down the pipeline
func (p *Pipeline) Stop() {
	p.log.Info("Stopping pipeline...")
	close(p.stopChan)
	p.wg.Wait()
	p.log.Info("Pipeline stopped")
}

// StartRecording begins accumulating audio data
func (p *Pipeline) StartRecording() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isRecording {
		return
	}

	// Drain channel to ensure no stale audio is included
	for len(p.audioChan) > 0 {
		<-p.audioChan
	}

	p.isRecording = true
	p.audioBuffer = nil // Clear buffer
	p.log.Info("Recording started")
}

// StopRecording stops accumulating and triggers processing
// Returns true if recording was stopped and processing started, false if not recording
func (p *Pipeline) StopRecording() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isRecording {
		return false
	}

	p.isRecording = false
	p.log.Info("Recording stopped", "buffer_size", len(p.audioBuffer))

	// Process the captured audio in a separate goroutine to not block
	// Make a copy of the buffer
	bufferCopy := make([]float32, len(p.audioBuffer))
	copy(bufferCopy, p.audioBuffer)

	go p.processSegment(bufferCopy)
	return true
}

func (p *Pipeline) captureLoop() {
	defer p.wg.Done()

	// Start audio capture
	err := p.audioEngine.StartRecording(p.audioChan)
	if err != nil {
		p.log.Error("Failed to start recording", "error", err)
		return
	}

	defer p.audioEngine.Stop()

	// Max recording duration (30 seconds at 16kHz)
	maxSamples := 16000 * 30

	for {
		select {
		case chunk := <-p.audioChan:
			p.mu.Lock()
			if p.isRecording {
				// Safety check: Auto-stop if recording gets too long (prevents OOM/Stuck state)
				if len(p.audioBuffer) >= maxSamples {
					p.log.Warn("Max recording duration reached (30s), forcing stop")
					p.isRecording = false

					// Copy and process immediately
					bufferCopy := make([]float32, len(p.audioBuffer))
					copy(bufferCopy, p.audioBuffer)

					// Launch processing in background
					go p.processSegment(bufferCopy)
				} else {
					p.audioBuffer = append(p.audioBuffer, chunk...)
				}
			}
			p.mu.Unlock()

		case <-p.stopChan:
			return
		}
	}
}

func (p *Pipeline) processSegment(samples []float32) {
	defer func() {
		if r := recover(); r != nil {
			p.log.Error("Recovered from panic in processSegment", "error", r)
		}
		if p.onCompletion != nil {
			p.onCompletion()
		}
	}()

	if len(samples) == 0 {
		p.log.Warn("Empty audio buffer, skipping processing")
		return
	}

	// Check duration (SampleRate is typically 16000)
	// If recording is less than 2 seconds, skip transcription
	durationSeconds := float64(len(samples)) / float64(p.vadParams.SampleRate)
	p.log.Info("Processing segment", "samples", len(samples), "rate", p.vadParams.SampleRate, "duration", durationSeconds)

	if durationSeconds < 2.0 {
		p.log.Info("Recording too short (< 2s), skipping transcription", "duration", durationSeconds)
		return
	}

	start := time.Now()

	// 1. ASR: Transcribe Audio
	text, err := p.asrEngine.Transcribe(samples)
	if err != nil {
		p.log.Error("ASR failed", "error", err)
		return
	}

	// Check word count
	// If detected less than 4 words, avoid transcribing completely (treat as false positive)
	// We do this after transcription as we need the text to count words
	words := strings.Fields(text)
	if len(words) < 4 {
		p.log.Info("Transcription too short (< 4 words), ignoring", "text", text, "word_count", len(words))
		return
	}

	if strings.TrimSpace(text) == "" {
		p.log.Info("No speech detected")
		return
	}

	p.log.Debug("ASR Output", "text", text, "duration", time.Since(start))

	// 2. Context: Get Current Window Info
	ctxInfo, err := p.ctxProvider.GetContext()
	if err != nil {
		p.log.Warn("Failed to get context", "error", err)
		// Proceed without context
	}

	// 3. LLM: Cleanup and Contextualize
	// TODO: Pass context info to LLM if supported
	cleanedText, err := p.llmEngine.CleanupText(text)
	if err != nil {
		p.log.Error("LLM cleanup failed", "error", err)
		// Fallback to raw text
		cleanedText = text
	}

	p.log.Info("Final Output",
		"raw", text,
		"cleaned", cleanedText,
		"app", ctxInfo.AppName,
		"window", ctxInfo.WindowTitle,
		"total_duration", time.Since(start),
	)

	// 4. Output: Print to Stdout
	fmt.Println(cleanedText)

	// 5. Output: Inject Text
	// First write to clipboard as backup/mechanism
	if err := clipboard.Write(cleanedText); err != nil {
		p.log.Error("Failed to write to clipboard", "error", err)
	}

	// Then inject via keyboard
	if p.injector != nil {
		if err := p.injector.Inject(cleanedText); err != nil {
			p.log.Error("Failed to inject text", "error", err)
		}
	}
}
