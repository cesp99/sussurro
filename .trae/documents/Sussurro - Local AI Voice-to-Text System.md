## Overview

**Sussurro** is a fully local, open-source, cross-platform desktop voice-to-text system that acts as a system-wide AI dictation layer. It transforms speech into clean, formatted, context-aware text injected into any application.

---

## Technology Stack

| Component | Technology | Rationale |
|-----------|------------|-----------|
| **Language** | Go | Core system language as specified |
| **ASR Model** | whisper.cpp (small or base model) | Fast, accurate, multilingual support via Go bindings |
| **LLM Model** | Qwen2.5-1.5B-Instruct (GGUF) | Best 1.5B model for text cleanup, llama.cpp compatible |
| **Audio Capture** | gordonklaus/portaudio | Cross-platform microphone capture |
| **Global Hotkey** | golang.design/x/hotkey | Cross-platform hotkey registration |
| **Text Injection** | micmonay/keybd_event | Cross-platform keystroke simulation |
| **AI Bindings** | ggerganov/whisper.cpp/bindings/go + go-skynet/go-llama.cpp | Go bindings for local inference |

---

## Project Structure

```
sussurro/
├── cmd/
│   └── sussurro/
│       └── main.go                 # Entry point
├── internal/
│   ├── audio/
│   │   ├── capture.go              # Microphone capture via PortAudio
│   │   └── buffer.go               # Audio buffering for streaming
│   ├── asr/
│   │   ├── whisper.go              # Whisper.cpp integration
│   │   └── transcriber.go          # Transcription pipeline
│   ├── llm/
│   │   ├── llama.go                # llama.cpp integration
│   │   └── cleanup.go              # Text cleanup prompts
│   ├── context/
│   │   ├── detector.go             # Active window detection
│   │   └── profiles.go             # Context-aware formatting profiles
│   ├── injection/
│   │   ├── injector.go             # Text injection interface
│   │   └── keyboard.go             # Keystroke emulation
│   ├── hotkey/
│   │   └── handler.go              # Global hotkey management
│   └── config/
│       └── config.go               # Configuration management
├── pkg/
│   └── pipeline/
│       └── pipeline.go             # Main processing pipeline orchestration
├── models/                         # Model files directory (gitignored)
├── configs/
│   └── default.yaml                # Default configuration
├── scripts/
│   ├── download-models.sh          # Model download script
│   └── build.sh                    # Build script
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## Implementation Phases

### Phase 1: Project Foundation
1. Initialize Go module and project structure
2. Set up configuration system (YAML-based)
3. Implement logging infrastructure
4. Create build system (Makefile)

### Phase 2: Audio Capture Layer
1. Integrate PortAudio via `gordonklaus/portaudio`
2. Implement microphone stream capture
3. Create audio buffer management (16kHz, 16-bit mono for Whisper)
4. Handle voice activity detection (VAD) basics

### Phase 3: ASR Integration (Whisper.cpp)
1. Integrate `ggerganov/whisper.cpp/bindings/go`
2. Implement model loading (whisper-small or whisper-base)
3. Create transcription pipeline with streaming support
4. Handle partial results for low latency

### Phase 4: LLM Text Cleanup Layer
1. Integrate `go-skynet/go-llama.cpp`
2. Implement Qwen2.5-1.5B-Instruct model loading
3. Create text cleanup prompt engineering:
   - Remove filler words ("um", "uh", "like")
   - Fix grammar and punctuation
   - Maintain semantic meaning
   - No content invention
4. Implement streaming inference for responsiveness

### Phase 5: Context Detection
1. Implement active window detection (platform-specific):
   - macOS: CGWindowListCopyWindowInfo
   - Windows: GetForegroundWindow
   - Linux: X11/Wayland APIs
2. Extract application name, window title
3. Create context-aware formatting profiles

### Phase 6: Text Injection Layer
1. Integrate `micmonay/keybd_event`
2. Implement keystroke simulation
3. Add clipboard fallback for complex text
4. Handle Unicode text properly

### Phase 7: Global Hotkey & Pipeline Integration
1. Integrate `golang.design/x/hotkey`
2. Implement hold-to-record pattern:
   - Hotkey down: Start recording
   - Hotkey up: Process and inject
3. Create main pipeline orchestration
4. Add status indicators (console-based initially)

### Phase 8: Configuration & Polish
1. Implement configurable hotkeys
2. Add model selection options
3. Create model download scripts
4. Write comprehensive README

---

## Recommended Whisper Models

| Model | Size | Speed | Accuracy | Multilingual |
|-------|------|-------|----------|--------------|
| **whisper-small** (recommended) | ~460MB | ~4x | Good | Yes |
| whisper-base | ~140MB | ~7x | Medium | Yes |
| whisper-small.en | ~460MB | ~4x | Better for EN | No |

**Recommendation**: Start with `whisper-small` (ggml-small.bin) for balance of speed and accuracy with multilingual support.

---

## LLM Model: Qwen2.5-1.5B-Instruct

- **Size**: ~1.5GB (Q4_K_M quantized: ~900MB)
- **Context**: 32K tokens
- **Format**: GGUF (llama.cpp compatible)
- **Strengths**: Excellent instruction following, text editing, multilingual
- **Source**: huggingface.co/Qwen/Qwen2.5-1.5B-Instruct-GGUF

---

## Text Cleanup Prompt Template

```
You are a text cleanup assistant. Your task is to clean up transcribed speech while preserving the original meaning. Rules:
1. Remove filler words (um, uh, like, you know)
2. Fix grammar and punctuation
3. Remove speech artifacts and repetitions
4. Maintain the speaker's intent and tone
5. Do NOT add new information
6. Do NOT change the meaning
7. Output ONLY the cleaned text, nothing else

Input: {raw_transcription}
Output:
```

---

## Key Dependencies

```go
// go.mod dependencies
require (
    github.com/ggerganov/whisper.cpp/bindings/go v0.0.0
    github.com/go-skynet/go-llama.cpp v0.0.0
    github.com/gordonklaus/portaudio v0.0.0
    github.com/micmonay/keybd_event v0.0.0
    github.com/micmonay/keybd_event v0.0.0
    golang.design/x/hotkey v0.0.0
)