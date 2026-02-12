package llm

import (
	"fmt"
	"os"

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
	// TinyLlama Chat template
	prompt := fmt.Sprintf(`<|system|>
You are a text cleanup assistant. Your task is to clean up transcribed speech while preserving the original meaning.
Rules:
1. Remove filler words (um, uh, like, you know)
2. Fix grammar and punctuation
3. Remove speech artifacts and repetitions
4. Maintain the speaker's intent and tone
5. Do NOT add new information
6. Do NOT change the meaning
7. Output ONLY the cleaned text, nothing else</s>
<|user|>
Input: %s</s>
<|assistant|>`, rawText)

	// We use Predict with empty options for now
	// SetThreads is a PredictOption, not ModelOption
	cleaned, err := e.model.Predict(prompt, llama.SetTokens(0), llama.SetThreads(e.threads))
	if err != nil {
		return "", fmt.Errorf("prediction failed: %w", err)
	}

	return cleaned, nil
}

// Close releases resources
func (e *Engine) Close() {
	if e.model != nil {
		e.model.Free()
	}
}
