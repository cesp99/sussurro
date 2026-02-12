package llm

import (
	"fmt"
	"os"
	"strings"

	llama "github.com/go-skynet/go-llama.cpp"
)

// Engine handles the LLM model and text generation
type Engine struct {
	model   *llama.LLama
	threads int
}

// NewEngine initializes the LLM model from a file path
func NewEngine(modelPath string, threads int, contextSize int, gpuLayers int) (*Engine, error) {
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("model file not found at %s: %w", modelPath, err)
	}

	model, err := llama.New(
		modelPath,
		llama.SetContext(contextSize),
		llama.SetGPULayers(gpuLayers),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load llm model: %w", err)
	}

	return &Engine{
		model:   model,
		threads: threads,
	}, nil
}

// CleanupText processes the raw transcription to remove artifacts and fix grammar
func (e *Engine) CleanupText(rawText string) (string, error) {
	// TinyLlama Chat template - Simplified for stability
	prompt := fmt.Sprintf(`<|system|>
You are a text cleanup assistant. Rewrite the user's text to remove filler words and fix grammar. Output ONLY the corrected text. Do not provide examples or notes.</s>
<|user|>
%s</s>
<|assistant|>`, rawText)

	// We use Predict with strict options
	cleaned, err := e.model.Predict(prompt, 
		llama.SetTokens(0), 
		llama.SetThreads(e.threads),
		llama.SetTemperature(0.1), // Low temperature for deterministic output
		llama.SetTopP(0.9),
	)
	if err != nil {
		return "", fmt.Errorf("prediction failed: %w", err)
	}

	// Post-processing cleanup
	cleaned = strings.TrimSpace(cleaned)
	
	// Cut off at common hallucination markers if stop strings didn't catch them
	if idx := strings.Index(cleaned, "Input:"); idx != -1 {
		cleaned = cleaned[:idx]
	}
	if idx := strings.Index(cleaned, "Example:"); idx != -1 {
		cleaned = cleaned[:idx]
	}
	if idx := strings.Index(cleaned, "<|user|>"); idx != -1 {
		cleaned = cleaned[:idx]
	}
	
	return strings.TrimSpace(cleaned), nil
}

// Close releases resources
func (e *Engine) Close() {
	if e.model != nil {
		e.model.Free()
	}
}
