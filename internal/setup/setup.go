package setup

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultConfigTemplate = `app:
  name: "Sussurro"
  debug: false
  log_level: "info" # debug, info, warn, error

audio:
  sample_rate: 16000
  channels: 1
  bit_depth: 16
  buffer_size: 1024
  max_duration: "60s"

models:
  asr:
    path: "{{ASR_PATH}}"
    type: "whisper"
    threads: 4
  llm:
    path: "{{LLM_PATH}}"
    context_size: 32768
    gpu_layers: 0
    threads: 4

hotkey:
  trigger: "ctrl+shift+space"

injection:
  method: "keyboard"
`
	// Whisper Small model (approx 500MB)
	urlASR = "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-small.bin"
	// Qwen 3 1.7B GGUF (approx 1GB+)
	urlLLM = "https://huggingface.co/enacimie/Qwen3-1.7B-Q4_K_M-GGUF/resolve/main/qwen3-1.7b-q4_k_m.gguf"
)

// EnsureSetup checks for the necessary configuration and models,
// and prompts the user to set them up if missing.
func EnsureSetup() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	sussurroDir := filepath.Join(homeDir, ".sussurro")
	modelsDir := filepath.Join(sussurroDir, "models")
	configFile := filepath.Join(sussurroDir, "config.yaml")

	// 1. Create .sussurro directory if it doesn't exist
	if _, err := os.Stat(sussurroDir); os.IsNotExist(err) {
		fmt.Println("Welcome to Sussurro! It looks like this is your first run.")
		fmt.Printf("Creating configuration directory at %s...\n", sussurroDir)
		if err := os.MkdirAll(modelsDir, 0755); err != nil {
			return fmt.Errorf("failed to create directories: %w", err)
		}
	} else {
		// Ensure models dir exists even if sussurro dir exists
		if err := os.MkdirAll(modelsDir, 0755); err != nil {
			return fmt.Errorf("failed to create models directory: %w", err)
		}
	}

	// 2. Create config.yaml if it doesn't exist
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Println("Creating default configuration file...")

		asrPath := filepath.Join(modelsDir, "ggml-small.bin")
		llmPath := filepath.Join(modelsDir, "qwen3-1.7b-q4_k_m.gguf")

		configContent := strings.ReplaceAll(defaultConfigTemplate, "{{ASR_PATH}}", asrPath)
		configContent = strings.ReplaceAll(configContent, "{{LLM_PATH}}", llmPath)

		if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
			return fmt.Errorf("failed to write config file: %w", err)
		}
		fmt.Printf("Configuration saved to %s\n", configFile)
	}

	// 3. Check for models and prompt to download
	asrPath := filepath.Join(modelsDir, "ggml-small.bin")
	llmPath := filepath.Join(modelsDir, "qwen3-1.7b-q4_k_m.gguf")

	missingASR := false
	missingLLM := false

	if _, err := os.Stat(asrPath); os.IsNotExist(err) {
		missingASR = true
	}
	if _, err := os.Stat(llmPath); os.IsNotExist(err) {
		missingLLM = true
	}

	if missingASR || missingLLM {
		fmt.Println("\nMissing model files:")
		if missingASR {
			fmt.Printf(" - Whisper Model (ASR): %s\n", asrPath)
		}
		if missingLLM {
			fmt.Printf(" - LLM Model: %s\n", llmPath)
		}

		fmt.Print("\nWould you like to download them now? (Y/n): ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response == "" || response == "y" || response == "yes" {
			if missingASR {
				if err := downloadFile(urlASR, asrPath, "Whisper Model"); err != nil {
					return fmt.Errorf("failed to download ASR model: %w", err)
				}
			}
			if missingLLM {
				if err := downloadFile(urlLLM, llmPath, "LLM Model"); err != nil {
					return fmt.Errorf("failed to download LLM model: %w", err)
				}
			}
			fmt.Println("\nAll models downloaded successfully!")
		} else {
			fmt.Println("Skipping download. Note: Sussurro may not function correctly without these models.")
		}
	}

	return nil
}

// downloadFile downloads a file from url to filepath with a simple progress indicator
func downloadFile(url, filepath, name string) error {
	fmt.Printf("Downloading %s...\n", name)

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create a proxy reader to track progress
	contentLength := resp.ContentLength
	reader := &progressReader{
		Reader: resp.Body,
		Total:  contentLength,
		Name:   name,
	}

	_, err = io.Copy(out, reader)
	fmt.Println() // Newline after progress
	return err
}

type progressReader struct {
	io.Reader
	Total   int64
	Current int64
	Name    string
	Last    int64
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.Current += int64(n)

	// Update progress every 1MB or so to avoid spamming stdout
	if pr.Current-pr.Last > 1024*1024 || pr.Current == pr.Total {
		pr.Last = pr.Current
		if pr.Total > 0 {
			percent := float64(pr.Current) / float64(pr.Total) * 100
			fmt.Printf("\rDownloading %s: %.1f%% (%.1f/%.1f MB)", pr.Name, percent, float64(pr.Current)/1024/1024, float64(pr.Total)/1024/1024)
		} else {
			fmt.Printf("\rDownloading %s: %.1f MB", pr.Name, float64(pr.Current)/1024/1024)
		}
	}

	return n, err
}
