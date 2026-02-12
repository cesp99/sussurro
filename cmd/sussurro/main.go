package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

	// Setup Signal Handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// TODO: Initialize other components (Audio, ASR, LLM, etc.)
	log.Info("Sussurro initialized and waiting for signals (Phase 1 complete)")

	// Wait for shutdown signal
	sig := <-sigChan
	log.Info("Received shutdown signal", "signal", sig.String())
	log.Info("Shutting down Sussurro...")
}
