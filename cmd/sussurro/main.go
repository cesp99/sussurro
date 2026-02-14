package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cesp99/sussurro/internal/asr"
	"github.com/cesp99/sussurro/internal/audio"
	"github.com/cesp99/sussurro/internal/config"
	"github.com/cesp99/sussurro/internal/context"
	"github.com/cesp99/sussurro/internal/hotkey"
	"github.com/cesp99/sussurro/internal/injection"
	"github.com/cesp99/sussurro/internal/llm"
	"github.com/cesp99/sussurro/internal/logger"
	"github.com/cesp99/sussurro/internal/pipeline"
	"github.com/cesp99/sussurro/internal/setup"
	"github.com/cesp99/sussurro/internal/trigger"
	"github.com/cesp99/sussurro/internal/version"

	"golang.design/x/hotkey/mainthread"
)

func main() {
	mainthread.Init(run)
}

func run() {
	// Parse command line flags
	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	// Ensure Setup (First Run Experience)
	if err := setup.EnsureSetup(); err != nil {
		fmt.Printf("Setup failed: %v\n", err)
		os.Exit(1)
	}

	// Load Configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize Logger
	log := logger.Init(cfg.App.LogLevel)
	log.Info("Starting Sussurro CLI", "version", version.Version)

	// Check if models exist
	if _, err := os.Stat(cfg.Models.ASR.Path); os.IsNotExist(err) {
		log.Error("ASR model missing", "path", cfg.Models.ASR.Path)
		fmt.Printf("Error: ASR model not found at %s. Please ensure models are downloaded.\n", cfg.Models.ASR.Path)
		os.Exit(1)
	}
	if _, err := os.Stat(cfg.Models.LLM.Path); os.IsNotExist(err) {
		log.Error("LLM model missing", "path", cfg.Models.LLM.Path)
		fmt.Printf("Error: LLM model not found at %s. Please ensure models are downloaded.\n", cfg.Models.LLM.Path)
		os.Exit(1)
	}

	// Initialize Context Provider
	ctxProvider := context.NewProvider()
	defer ctxProvider.Close()

	// Initialize Audio Capture
	audioEngine, err := audio.NewCaptureEngine(cfg.Audio.SampleRate, cfg.Audio.Channels)
	if err != nil {
		log.Error("Failed to initialize audio engine", "error", err)
		os.Exit(1)
	}
	defer audioEngine.Close()

	// Initialize ASR Engine
	asrEngine, err := asr.NewEngine(cfg.Models.ASR.Path, cfg.Models.ASR.Threads, cfg.App.Debug)
	if err != nil {
		log.Error("Failed to initialize ASR engine", "error", err)
		os.Exit(1)
	}
	defer asrEngine.Close()

	// Initialize LLM Engine
	llmEngine, err := llm.NewEngine(cfg.Models.LLM.Path, cfg.Models.LLM.Threads, cfg.Models.LLM.ContextSize, cfg.Models.LLM.GpuLayers, cfg.App.Debug)
	if err != nil {
		log.Error("Failed to initialize LLM engine", "error", err)
		os.Exit(1)
	}
	defer llmEngine.Close()

	// Initialize Injector
	injector, err := injection.NewInjector()
	if err != nil {
		log.Error("Failed to initialize injector", "error", err)
		// We might want to continue even if injector fails, depending on requirements,
		// but typically it's essential for this app.
	}

	// Initialize and Start Pipeline
	pipe := pipeline.NewPipeline(audioEngine, asrEngine, llmEngine, ctxProvider, injector, log, cfg.Audio.SampleRate, cfg.Audio.MaxDuration)

	// No UI to update on completion, so we can pass nil or a logging callback
	pipe.SetOnCompletion(func() {
		log.Debug("Pipeline processing completed")
	})

	err = pipe.Start()
	if err != nil {
		log.Error("Failed to start pipeline", "error", err)
		os.Exit(1)
	}
	defer pipe.Stop()

	// Initialize input handler (hotkey or trigger server depending on environment)
	if hotkey.IsWayland() {
		log.Debug("Wayland detected - using trigger server")

		triggerServer, err := trigger.NewServer(log)
		if err != nil {
			log.Error("Failed to initialize trigger server", "error", err)
			os.Exit(1)
		}
		defer triggerServer.Stop()

		err = triggerServer.Start(
			func() { // On trigger start
				log.Debug("Trigger: Starting recording")
				pipe.StartRecording()
			},
			func() { // On trigger stop
				log.Debug("Trigger: Stopping recording")
				if !pipe.StopRecording() {
					log.Debug("Recording was not active or already stopped")
				}
			},
		)
		if err != nil {
			log.Error("Failed to start trigger server", "error", err)
			os.Exit(1)
		}

		log.Warn("Wayland detected: Configure keyboard shortcut (see docs/WAYLAND.md)")
	} else {
		log.Info("X11 detected - using global hotkeys")

		hkHandler, err := hotkey.NewHandler(cfg.Hotkey.Trigger, log)
		if err != nil {
			log.Error("Failed to initialize hotkey handler", "error", err)
			os.Exit(1)
		}
		defer hkHandler.Unregister()

		err = hkHandler.Register(
			func() { // On Key Down
				log.Debug("Hotkey pressed: Starting recording")
				pipe.StartRecording()
			},
			func() { // On Key Up
				log.Debug("Hotkey released: Stopping recording")
				if !pipe.StopRecording() {
					log.Debug("Recording was not active or already stopped")
				}
			},
		)
		if err != nil {
			log.Error("Failed to register hotkey", "error", err)
			os.Exit(1)
		}
	}

	log.Info("Sussurro running. Press Ctrl+C to exit.")

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until signal received
	sig := <-sigChan
	log.Info("Received signal, shutting down...", "signal", sig)

	// Defer statements will handle cleanup in reverse order
}
