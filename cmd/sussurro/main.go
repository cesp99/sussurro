package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cesp99/sussurro/internal/asr"
	"github.com/cesp99/sussurro/internal/audio"
	"github.com/cesp99/sussurro/internal/config"
	"github.com/cesp99/sussurro/internal/logger"
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

	// Setup Signal Handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// TODO: Initialize LLM and Pipeline
	log.Info("Sussurro initialized and waiting for signals")

	// Wait for shutdown signal
	sig := <-sigChan
	log.Info("Received shutdown signal", "signal", sig.String())
	log.Info("Shutting down Sussurro...")
}
