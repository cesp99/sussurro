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
	
	cleaned = strings.TrimSpace(cleaned)

	// Anti-Hallucination Check
	if !validateOutput(rawText, cleaned) {
		return rawText, nil // Fallback to raw text
	}
	
	return cleaned, nil
}

func validateOutput(raw, cleaned string) bool {
	// 1. Length Check
	// If cleaned is significantly longer than raw (more than 2x), it's likely a hallucination
	// unless raw is very short.
	if len(raw) > 10 && len(cleaned) > len(raw)*2 {
		return false
	}

	// 2. Pattern Check for Common Hallucinations
	lowerCleaned := strings.ToLower(cleaned)
	invalidPrefixes := []string{
		"the user", "input:", "output:", "rewrite", "corrected text:", 
		"here is", "sure, i can", "i'm sorry", "assistant:",
	}
	for _, prefix := range invalidPrefixes {
		if strings.HasPrefix(lowerCleaned, prefix) {
			return false
		}
	}

	// 3. Semantic Content Check (Bag of Words)
	// Ensure significant words from raw text are present in cleaned text.
	// We ignore common filler words.
	rawWords := strings.Fields(strings.ToLower(raw))
	cleanedLower := strings.ToLower(cleaned)
	
	missingCount := 0
	totalSignificant := 0
	
	// Basic stop words to ignore
	stopWords := map[string]bool{
		"umm": true, "ah": true, "uh": true, "like": true, "so": true, 
		"just": true, "a": true, "an": true, "the": true,
	} 
	
	for _, w := range rawWords {
		// Clean punctuation
		w = strings.Trim(w, ".,!?-")
		if w == "" || stopWords[w] { continue }
		
		totalSignificant++
		// Check if word exists in cleaned text
		if !strings.Contains(cleanedLower, w) {
			missingCount++
		}
	}
	
	// If we are missing more than 50% of significant words, it's likely a hallucination
	// (or a complete rewrite which we don't want)
	if totalSignificant > 0 && float64(missingCount)/float64(totalSignificant) > 0.5 {
		return false
	}
	
	return true
}

// Close releases resources
func (e *Engine) Close() {
	if e.model != nil {
		e.model.Free()
	}
}
