package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cesp99/sussurro/internal/asr"
	"github.com/cesp99/sussurro/internal/audio"
	"github.com/cesp99/sussurro/internal/config"
	"github.com/cesp99/sussurro/internal/context"
	"github.com/cesp99/sussurro/internal/llm"
	"github.com/cesp99/sussurro/internal/logger"
	"github.com/cesp99/sussurro/internal/pipeline"
)

func main() {
	// Load Configuration
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		// Fallback logging if config fails
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize Logger
	log := logger.Init(cfg.App.LogLevel)
	log.Info("Starting Sussurro", "version", cfg.App.Version)

	// Initialize Context Provider (Phase 5)
	// TODO: Add config for provider selection (default to macOS for now)
	ctxProvider := context.NewMacOSProvider()
	defer ctxProvider.Close()
	
	// Test Context Detection
	ctxInfo, err := ctxProvider.GetContext()
	if err != nil {
		log.Warn("Failed to get initial context", "error", err)
	} else {
		log.Info("Initial Context", "app", ctxInfo.AppName, "window", ctxInfo.WindowTitle)
	}

	// Initialize Audio Capture (Phase 2)
	audioEngine, err := audio.NewCaptureEngine(cfg.Audio.SampleRate, cfg.Audio.Channels)
	if err != nil {
		log.Error("Failed to initialize audio engine", "error", err)
		os.Exit(1)
	}
	defer audioEngine.Close()
	log.Info("Audio engine initialized", "sample_rate", cfg.Audio.SampleRate, "channels", cfg.Audio.Channels)

	// Initialize ASR Engine (Phase 3)
	// Check if model exists first
	if _, err := os.Stat(cfg.Models.ASR.Path); os.IsNotExist(err) {
		log.Error("ASR model not found. Please run scripts/download-models.sh", "path", cfg.Models.ASR.Path)
		os.Exit(1)
	}

	asrEngine, err := asr.NewEngine(cfg.Models.ASR.Path, cfg.Models.ASR.Threads)
	if err != nil {
		log.Error("Failed to initialize ASR engine", "error", err)
		os.Exit(1)
	}
	defer asrEngine.Close()
	log.Info("ASR engine initialized", "model", cfg.Models.ASR.Path)

	// Initialize LLM Engine (Phase 4)
	if _, err := os.Stat(cfg.Models.LLM.Path); os.IsNotExist(err) {
		log.Error("LLM model not found. Please run scripts/download-models.sh", "path", cfg.Models.LLM.Path)
		os.Exit(1)
	}

	llmEngine, err := llm.NewEngine(cfg.Models.LLM.Path, cfg.Models.LLM.Threads, cfg.Models.LLM.ContextSize, cfg.Models.LLM.GpuLayers)
	if err != nil {
		log.Error("Failed to initialize LLM engine", "error", err)
		os.Exit(1)
	}
	defer llmEngine.Close()
	log.Info("LLM engine initialized", "model", cfg.Models.LLM.Path)

	// Initialize and Start Pipeline (Phase 6)
	pipe := pipeline.NewPipeline(audioEngine, asrEngine, llmEngine, ctxProvider, log)
	err = pipe.Start()
	if err != nil {
		log.Error("Failed to start pipeline", "error", err)
		os.Exit(1)
	}
	defer pipe.Stop()

	// Setup Signal Handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Info("Sussurro initialized and waiting for signals")

	// Wait for shutdown signal
	sig := <-sigChan
	log.Info("Received shutdown signal", "signal", sig.String())
	log.Info("Shutting down Sussurro...")
}
