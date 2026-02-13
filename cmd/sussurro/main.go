package main

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"

	"github.com/cesp99/sussurro/internal/asr"
	"github.com/cesp99/sussurro/internal/audio"
	"github.com/cesp99/sussurro/internal/config"
	"github.com/cesp99/sussurro/internal/context"
	"github.com/cesp99/sussurro/internal/hotkey"
	"github.com/cesp99/sussurro/internal/injection"
	"github.com/cesp99/sussurro/internal/llm"
	"github.com/cesp99/sussurro/internal/logger"
	"github.com/cesp99/sussurro/internal/pipeline"
	"github.com/cesp99/sussurro/internal/ui"
	"github.com/cesp99/sussurro/internal/ui/theme"
)

func main() {
	// Initialize Fyne App
	a := app.New()
	a.Settings().SetTheme(&theme.SussurroTheme{})

	// Load Configuration
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	// Initialize Logger
	log := logger.Init(cfg.App.LogLevel)
	log.Info("Starting Sussurro", "version", cfg.App.Version)

	// Check if models exist
	modelsExist := true
	if _, err := os.Stat(cfg.Models.ASR.Path); os.IsNotExist(err) {
		modelsExist = false
	}
	if _, err := os.Stat(cfg.Models.LLM.Path); os.IsNotExist(err) {
		modelsExist = false
	}

	// Create Main Window (Model Manager)
	w := a.NewWindow("Sussurro Models")
	manager := ui.NewModelManager(w)
	w.SetContent(manager.GetContent())
	w.Resize(fyne.NewSize(600, 500))

	// Create Overlay Window
	overlay := ui.NewOverlayWindow(a)
	overlay.SetState(ui.StateLoading)
	overlay.Show()

	// If models are missing, show window immediately
	if !modelsExist {
		w.Show()
	}

	// Setup System Tray
	if desk, ok := a.(desktop.App); ok {
		m := fyne.NewMenu("Sussurro",
			fyne.NewMenuItem("Show Models", func() {
				w.Show()
			}),
			fyne.NewMenuItem("Toggle Overlay", func() {
				// Simple toggle for now
				overlay.Show()
			}),
			fyne.NewMenuItem("Quit", func() {
				a.Quit()
			}),
		)
		desk.SetSystemTrayMenu(m)
	}

	// Start Backend Services in Background
	go func() {
		// Initialize Context Provider
		ctxProvider := context.NewMacOSProvider()
		defer ctxProvider.Close()

		// Initialize Audio Capture
		audioEngine, err := audio.NewCaptureEngine(cfg.Audio.SampleRate, cfg.Audio.Channels)
		if err != nil {
			log.Error("Failed to initialize audio engine", "error", err)
			return
		}
		defer audioEngine.Close()

		// Initialize ASR Engine
		if _, err := os.Stat(cfg.Models.ASR.Path); os.IsNotExist(err) {
			log.Warn("ASR model missing. Please download via UI.")
			return
		}

		asrEngine, err := asr.NewEngine(cfg.Models.ASR.Path, cfg.Models.ASR.Threads)
		if err != nil {
			log.Error("Failed to initialize ASR engine", "error", err)
			return
		}
		defer asrEngine.Close()

		// Initialize LLM Engine
		if _, err := os.Stat(cfg.Models.LLM.Path); os.IsNotExist(err) {
			log.Warn("LLM model missing. Please download via UI.")
			return
		}

		llmEngine, err := llm.NewEngine(cfg.Models.LLM.Path, cfg.Models.LLM.Threads, cfg.Models.LLM.ContextSize, cfg.Models.LLM.GpuLayers)
		if err != nil {
			log.Error("Failed to initialize LLM engine", "error", err)
			return
		}
		defer llmEngine.Close()

		// Initialize Injector
		injector, err := injection.NewInjector()
		if err != nil {
			log.Error("Failed to initialize injector", "error", err)
		}

		// Initialize and Start Pipeline
		pipe := pipeline.NewPipeline(audioEngine, asrEngine, llmEngine, ctxProvider, injector, log)
		pipe.SetOnCompletion(func() {
			overlay.SetState(ui.StateIdle)
		})
		err = pipe.Start()
		if err != nil {
			log.Error("Failed to start pipeline", "error", err)
			return
		}
		defer pipe.Stop()

		// Backend ready
		overlay.SetState(ui.StateIdle)

		// Initialize Hotkey Handler
		hkHandler, err := hotkey.NewHandler(cfg.Hotkey.Trigger, log)
		if err != nil {
			log.Error("Failed to initialize hotkey handler", "error", err)
			return
		}

		// Register Hotkey Callbacks
		err = hkHandler.Register(
			func() { // On Key Down
				overlay.SetState(ui.StateListening)
				pipe.StartRecording()
			},
			func() { // On Key Up
				if pipe.StopRecording() {
					overlay.SetState(ui.StateTranscribing)
				} else {
					// Was not recording (likely auto-stopped due to duration limit)
					// Ensure UI is back to idle
					overlay.SetState(ui.StateIdle)
				}
			},
		)
		if err != nil {
			log.Error("Failed to register hotkey", "error", err)
			return
		}
		defer hkHandler.Unregister()

		log.Info("Sussurro backend running")

		// Block forever (until app quit)
		select {}
	}()

	a.Run()
}
