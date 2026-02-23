package ui

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/cesp99/sussurro/internal/config"
	"github.com/cesp99/sussurro/internal/setup"
	"github.com/cesp99/sussurro/internal/version"
)

// modelInfo describes a model for the settings UI.
type modelInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"desc"`
	Size        string `json:"size"`
	Installed   bool   `json:"installed"`
	Active      bool   `json:"active"`
	Type        string `json:"type"` // "whisper" or "llm"
}

// initialData is returned by getInitialData().
type initialData struct {
	Platform  string      `json:"platform"`
	Version   string      `json:"version"`
	Models    []modelInfo `json:"models"`
	Hotkey    string      `json:"hotkey"`
	IsWayland bool        `json:"isWayland"`
}

// bindBridge attaches all Go↔JS bridge functions to the webview.
func bindBridge(sw *settingsWindow) {
	mgr := sw.mgr

	sw.w.Bind("getInitialData", func() (result string) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("panic in getInitialData", "error", r)
				result = `{"error":"internal error"}`
			}
		}()
		data := buildInitialData(mgr)
		b, _ := json.Marshal(data)
		return string(b)
	})

	sw.w.Bind("saveHotkey", func(trigger string) (result string) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("panic in saveHotkey", "error", r)
				result = fmt.Sprintf("error: panic: %v", r)
			}
		}()
		if err := config.SaveHotkey(mgr.cfg, trigger); err != nil {
			return fmt.Sprintf("error: %v", err)
		}
		mgr.cfg.Hotkey.Trigger = trigger
		// Re-register the OS-level hotkey with the new trigger so it takes
		// effect immediately without requiring a restart.
		go mgr.reinstallHotkey(trigger)
		return "ok"
	})

	sw.w.Bind("downloadModel", func(modelID string) {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("panic in downloadModel goroutine", "error", r)
				}
			}()
			url, dest, name := resolveModelDownload(modelID)
			if url == "" {
				return
			}
			setup.SetProgressCallback(func(n string, pct float64, _, _ int64) {
				sw.pushDownloadProgress(n, pct)
			})
			defer setup.SetProgressCallback(nil)
			if err := setup.DownloadModel(url, dest, name); err != nil {
				sw.w.Dispatch(func() {
					sw.w.Eval(fmt.Sprintf("onDownloadError('%s', '%v')", modelID, err))
				})
				return
			}
			sw.w.Dispatch(func() {
				sw.w.Eval(fmt.Sprintf("onDownloadComplete('%s')", modelID))
			})
		}()
	})

	sw.w.Bind("setActiveModel", func(modelID string) (result string) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("panic in setActiveModel", "error", r)
				result = fmt.Sprintf("error: panic: %v", r)
			}
		}()
		if err := setup.SetActiveModel(modelID); err != nil {
			return fmt.Sprintf("error: %v", err)
		}
		// Config written — restart the process to load the new model.
		go func() {
			time.Sleep(300 * time.Millisecond)
			exe, err := os.Executable()
			if err != nil {
				slog.Error("restart: cannot resolve executable", "error", err)
				os.Exit(0)
			}
			cmd := exec.Command(exe, os.Args[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Start(); err != nil {
				slog.Error("restart: failed to start new process", "error", err)
			}
			os.Exit(0)
		}()
		return "ok"
	})

	sw.w.Bind("openURL", func(url string) {
		go func() {
			var cmd *exec.Cmd
			if runtime.GOOS == "darwin" {
				cmd = exec.Command("open", url)
			} else {
				cmd = exec.Command("xdg-open", url)
			}
			if err := cmd.Start(); err != nil {
				slog.Error("openURL failed", "url", url, "error", err)
			}
		}()
	})

	sw.w.Bind("closeSettings", func() {
		sw.Hide()
	})
}

func buildInitialData(mgr *Manager) initialData {
	homeDir, _ := os.UserHomeDir()
	modelsDir := homeDir + "/.sussurro/models"

	whisperSmallPath := modelsDir + "/ggml-small.bin"
	whisperLargePath := modelsDir + "/ggml-large-v3-turbo.bin"
	llmPath := modelsDir + "/qwen3-sussurro-q4_k_m.gguf"

	currentASR := mgr.cfg.Models.ASR.Path
	currentLLM := mgr.cfg.Models.LLM.Path

	models := []modelInfo{
		{
			ID:          "whisper-small",
			Name:        "Whisper Small",
			Description: "Faster, lower memory usage",
			Size:        "~488 MB",
			Installed:   fileExists(whisperSmallPath),
			Active:      currentASR == whisperSmallPath,
			Type:        "whisper",
		},
		{
			ID:          "whisper-large-v3-turbo",
			Name:        "Whisper Large v3 Turbo",
			Description: "Higher accuracy, more memory",
			Size:        "~1.62 GB",
			Installed:   fileExists(whisperLargePath),
			Active:      currentASR == whisperLargePath,
			Type:        "whisper",
		},
		{
			ID:          "qwen3-sussurro",
			Name:        "Qwen 3 Sussurro",
			Description: "Fine-tuned for transcription cleanup",
			Size:        "~1.28 GB",
			Installed:   fileExists(llmPath),
			Active:      currentLLM == llmPath,
			Type:        "llm",
		},
	}

	platform := "LINUX"
	if runtime.GOOS == "darwin" {
		platform = "MACOS"
	}

	isWayland := os.Getenv("WAYLAND_DISPLAY") != "" ||
		os.Getenv("XDG_SESSION_TYPE") == "wayland"
	if isWayland {
		platform += " (WAYLAND)"
	} else if runtime.GOOS == "linux" {
		platform += " (X11)"
	}

	return initialData{
		Platform:  platform,
		Version:   version.Version,
		Models:    models,
		Hotkey:    mgr.cfg.Hotkey.Trigger,
		IsWayland: isWayland,
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// resolveModelDownload maps a model ID to its download URL and local path.
func resolveModelDownload(modelID string) (url, dest, name string) {
	homeDir, _ := os.UserHomeDir()
	modelsDir := homeDir + "/.sussurro/models"

	switch modelID {
	case "whisper-small":
		return "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-small.bin",
			modelsDir + "/ggml-small.bin",
			"Whisper Small"
	case "whisper-large-v3-turbo":
		return "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-large-v3-turbo.bin",
			modelsDir + "/ggml-large-v3-turbo.bin",
			"Whisper Large v3 Turbo"
	case "qwen3-sussurro":
		return "https://huggingface.co/cesp99/qwen3-sussurro/resolve/main/qwen3-sussurro-q4_k_m.gguf",
			modelsDir + "/qwen3-sussurro-q4_k_m.gguf",
			"Qwen 3 Sussurro"
	}
	return "", "", ""
}
