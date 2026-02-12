package pipeline

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/cesp99/sussurro/internal/asr"
	"github.com/cesp99/sussurro/internal/audio"
	ctxProvider "github.com/cesp99/sussurro/internal/context"
	"github.com/cesp99/sussurro/internal/llm"
)

// Pipeline orchestrates the flow of data from audio capture to text output
type Pipeline struct {
	audioEngine *audio.CaptureEngine
	asrEngine   *asr.Engine
	llmEngine   *llm.Engine
	ctxProvider ctxProvider.Provider
	log         *slog.Logger

	// Channels for data flow
	audioChan chan []float32
	textChan  chan string
	stopChan  chan struct{}
	wg        sync.WaitGroup
}

// NewPipeline creates a new processing pipeline
func NewPipeline(
	audioEngine *audio.CaptureEngine,
	asrEngine *asr.Engine,
	llmEngine *llm.Engine,
	ctxProvider ctxProvider.Provider,
	log *slog.Logger,
) *Pipeline {
	return &Pipeline{
		audioEngine: audioEngine,
		asrEngine:   asrEngine,
		llmEngine:   llmEngine,
		ctxProvider: ctxProvider,
		log:         log,
		audioChan:   make(chan []float32, 100), // Buffer audio chunks
		textChan:    make(chan string, 10),     // Buffer text segments
		stopChan:    make(chan struct{}),
	}
}

// Start begins the pipeline processing
func (p *Pipeline) Start() error {
	p.log.Info("Starting pipeline...")

	// 1. Start Audio Capture Loop
	p.wg.Add(1)
	go p.captureLoop()

	// 2. Start Processing Loop (ASR + LLM)
	p.wg.Add(1)
	go p.processLoop()

	return nil
}

// Stop gracefully shuts down the pipeline
func (p *Pipeline) Stop() {
	p.log.Info("Stopping pipeline...")
	close(p.stopChan)
	p.wg.Wait()
	p.log.Info("Pipeline stopped")
}

func (p *Pipeline) captureLoop() {
	defer p.wg.Done()
	
	// Start audio capture
	err := p.audioEngine.StartRecording(p.audioChan)
	if err != nil {
		p.log.Error("Failed to start recording", "error", err)
		return
	}
	
	<-p.stopChan
	p.audioEngine.Stop()
}

func (p *Pipeline) processLoop() {
	defer p.wg.Done()

	var audioBuffer []float32
	// Configurable silence threshold and duration for VAD-like behavior
	// const silenceThreshold = 0.01
	// const silenceDuration = 500 * time.Millisecond 

	for {
		select {
		case chunk := <-p.audioChan:
			audioBuffer = append(audioBuffer, chunk...)
			
			// Simple segmentation logic (placeholder)
			// If buffer is long enough (e.g., 2 seconds), process it
			if len(audioBuffer) > 16000*2 { // 2 seconds at 16kHz
				p.processSegment(audioBuffer)
				audioBuffer = nil // Clear buffer
			}
			
		case <-p.stopChan:
			return
		}
	}
}

func (p *Pipeline) processSegment(samples []float32) {
	// 1. ASR: Transcribe Audio
	text, err := p.asrEngine.Transcribe(samples)
	if err != nil {
		p.log.Error("ASR failed", "error", err)
		return
	}
	
	if strings.TrimSpace(text) == "" {
		return
	}
	
	p.log.Debug("ASR Output", "text", text)

	// 2. Context: Get Current Window Info
	ctxInfo, err := p.ctxProvider.GetContext()
	if err != nil {
		p.log.Warn("Failed to get context", "error", err)
		// Proceed without context
	}

	// 3. LLM: Cleanup and Contextualize
	cleanedText, err := p.llmEngine.CleanupText(text) // We might want to add context here later
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
	)
	
	// 4. Output: Print to Stdout (or clipboard later)
	fmt.Println(cleanedText)
}
