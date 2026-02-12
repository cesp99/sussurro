package main

import (
	"fmt"
	"os"

	"github.com/cesp99/sussurro/internal/asr"
	"github.com/cesp99/sussurro/internal/audio"
	"github.com/cesp99/sussurro/internal/config"
	"github.com/cesp99/sussurro/internal/context"
	"github.com/cesp99/sussurro/internal/hotkey"
	"github.com/cesp99/sussurro/internal/injection"
	"github.com/cesp99/sussurro/internal/llm"
	"github.com/cesp99/sussurro/internal/logger"
	"github.com/cesp99/sussurro/internal/pipeline"
	"github.com/getlantern/systray"
)

func main() {
	// Systray must run on the main thread
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("Sussurro")
	systray.SetTooltip("Sussurro AI Assistant")
	
	mStatus := systray.AddMenuItem("Status: Idle", "Current status")
	mStatus.Disable()
	
	mQuit := systray.AddMenuItem("Quit", "Quit Sussurro")

	// Load Configuration
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		systray.Quit()
		return
	}

	// Initialize Logger
	log := logger.Init(cfg.App.LogLevel)
	log.Info("Starting Sussurro", "version", cfg.App.Version)

	go func() {
		// Initialize Context Provider (Phase 5)
		ctxProvider := context.NewMacOSProvider()
		defer ctxProvider.Close()

		// Initialize Audio Capture (Phase 2)
		audioEngine, err := audio.NewCaptureEngine(cfg.Audio.SampleRate, cfg.Audio.Channels)
		if err != nil {
			log.Error("Failed to initialize audio engine", "error", err)
			return
		}
		defer audioEngine.Close()
		log.Info("Audio engine initialized", "sample_rate", cfg.Audio.SampleRate, "channels", cfg.Audio.Channels)

		// Initialize ASR Engine (Phase 3)
		if _, err := os.Stat(cfg.Models.ASR.Path); os.IsNotExist(err) {
			log.Error("ASR model not found. Please run scripts/download-models.sh", "path", cfg.Models.ASR.Path)
			return
		}

		asrEngine, err := asr.NewEngine(cfg.Models.ASR.Path, cfg.Models.ASR.Threads)
		if err != nil {
			log.Error("Failed to initialize ASR engine", "error", err)
			return
		}
		defer asrEngine.Close()
		log.Info("ASR engine initialized", "model", cfg.Models.ASR.Path)

		// Initialize LLM Engine (Phase 4)
		if _, err := os.Stat(cfg.Models.LLM.Path); os.IsNotExist(err) {
			log.Error("LLM model not found. Please run scripts/download-models.sh", "path", cfg.Models.LLM.Path)
			return
		}

		llmEngine, err := llm.NewEngine(cfg.Models.LLM.Path, cfg.Models.LLM.Threads, cfg.Models.LLM.ContextSize, cfg.Models.LLM.GpuLayers)
		if err != nil {
			log.Error("Failed to initialize LLM engine", "error", err)
			return
		}
		defer llmEngine.Close()
		log.Info("LLM engine initialized", "model", cfg.Models.LLM.Path)

		// Initialize Injector (Phase 6)
		injector, err := injection.NewInjector()
		if err != nil {
			log.Error("Failed to initialize injector", "error", err)
			// Proceed without injector (will fall back to clipboard only)
		}
		log.Info("Injector initialized")

		// Initialize and Start Pipeline (Phase 7 Integration)
		pipe := pipeline.NewPipeline(audioEngine, asrEngine, llmEngine, ctxProvider, injector, log)
		err = pipe.Start()
		if err != nil {
			log.Error("Failed to start pipeline", "error", err)
			return
		}
		defer pipe.Stop()

		// Initialize Hotkey Handler (Phase 7)
		hkHandler, err := hotkey.NewHandler(cfg.Hotkey.Trigger, log)
		if err != nil {
			log.Error("Failed to initialize hotkey handler", "error", err)
			return
		}

		// Register Hotkey Callbacks
		err = hkHandler.Register(
			func() { // On Key Down
				mStatus.SetTitle("Status: Recording...")
				pipe.StartRecording()
			},
			func() { // On Key Up
				mStatus.SetTitle("Status: Processing...")
				pipe.StopRecording()
				// After processing (which is async), we ideally want to set status back to Idle.
				// But we don't have a callback for "Processing Done" here yet.
				// For now, we can just leave it or set a timer.
				// Or pipeline could accept a status callback.
				// For simplicity, we just leave "Processing..." or set it to "Idle" after a delay.
				// The user will see the text appear.
				
				// Quick hack: Reset status after 1 second (might be too fast)
				// Better: pipeline emits events.
				mStatus.SetTitle("Status: Idle") 
			},
		)
		if err != nil {
			log.Error("Failed to register hotkey", "error", err)
			return
		}
		defer hkHandler.Unregister()

		log.Info("Sussurro initialized and running", "hotkey", cfg.Hotkey.Trigger)
		
		// Wait for quit signal from tray
		<-mQuit.ClickedCh
		log.Info("Quit requested from tray")
		systray.Quit()
	}()
}

func onExit() {
	// Cleanup happens via defers in the goroutine when it returns/exits, 
	// but since systray.Quit() kills the app, we might want to ensure graceful shutdown.
	// For now, the OS cleanup is sufficient for this stage.
	fmt.Println("Sussurro exiting...")
}
