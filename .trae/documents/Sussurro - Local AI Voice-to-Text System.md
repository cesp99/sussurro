# Implementation Plan - Sussurro

A local-first, AI-powered voice-to-text system for macOS, Windows, and Linux.

## 1. Project Initialization & Structure
- **Initialize Go Module**: `github.com/cesp99/sussurro`
- **Directory Structure**:
  ```
  /cmd
    /daemon      # Main background process
    /ui          # Settings/Tray UI
  /internal
    /audio       # Microphone capture
    /whisper     # ASR implementation
    /llm         # Text refinement (Qwen 1.5B)
    /input       # Keyboard injection
    /hotkey      # Global hotkey management
    /config      # Configuration management
    /utils       # Helper functions
  /models        # Directory for storing GGUF models
  ```

## 2. Configuration & Model Management
- **Config System**: YAML/JSON config to store preferences.
- **Model Selection**: Support the following Whisper models (GGUF format):
  - **`medium.en` (Default)**: Best balance for English.
  - `large-v3-turbo`: Fast & accurate (Multilingual).
  - `large-v3`: Maximum accuracy (Slow).
  - `small.en`: Faster, lower resource usage.
  - `base.en` / `tiny.en`: Extremely fast, lower accuracy.
- **LLM Model**: Default to `Qwen2.5-1.5B-Instruct-GGUF` for text cleanup.
- **Model Downloader**: Utility to automatically fetch selected models from HuggingFace if missing.

## 3. Audio Capture Layer
- **Library**: `github.com/gordonklaus/portaudio` (via cgo).
- **Functionality**: 
  - Detect default input device.
  - Record 16kHz mono audio (Whisper requirement).
  - VAD (Voice Activity Detection) - simple energy-based or WebRTC VAD to detect speech end.

## 4. AI Pipeline Implementation
### Stage 1: ASR (Whisper)
- **Library**: `github.com/ggerganov/whisper.cpp/bindings/go`
- **Features**:
  - Load model based on config.
  - Stream audio chunks to `whisper_full`.
  - Return raw text segments.

### Stage 2: Text Refinement (LLM)
- **Library**: `github.com/go-skynet/go-llama.cpp` or `github.com/ollama/ollama/api` (Embedding `llama.cpp` directly is preferred for standalone).
- **Model**: `Qwen2.5-1.5B-Instruct` (Low VRAM/CPU usage).
- **Prompting**: System prompt to remove filler words, fix punctuation, and format based on context (if detectable) without hallucinating.

## 5. System Integration
- **Global Hotkey**: `golang.design/x/hotkey`
  - Default: `Ctrl + Alt + Space` (Configurable).
  - Behavior: Push-to-talk or Toggle (Configurable).
- **Text Injection**: `github.com/micmonay/keybd_event`
  - Type simulated keystrokes into the active window.
  - Fallback to Clipboard paste if typing is too slow/problematic.

## 6. User Interface (Wails/Systray)
- **Tray Icon**: Simple menu to Quit, Open Settings, Toggle Status.
- **Settings Window** (Wails):
  - Model selector (Dropdown: Tiny -> Large).
  - Hotkey configuration.
  - Download progress bars for models.

## 7. Execution Steps
1.  **Setup**: Create project structure and install dependencies.
2.  **Core**: Implement Audio Capture + Hotkey trigger.
3.  **ASR**: Integrate Whisper.cpp and test raw transcription.
4.  **LLM**: Integrate Qwen 1.5B and test cleanup.
5.  **Output**: Implement keyboard injection.
6.  **UI**: Build simple settings UI to switch models.
7.  **Packaging**: Ensure binaries are portable (handle `.dylib` / `.dll` embedding).
