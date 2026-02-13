# Sussurro

Sussurro is a fully local, open-source, cross-platform CLI voice-to-text system that acts as a system-wide AI dictation layer. It transforms speech into clean, formatted, context-aware text injected into any application.

## Overview

Sussurro uses local AI models to ensure privacy and low latency. It combines:
- **Whisper.cpp** for automatic speech recognition (ASR).
- **LLMs (TinyLlama, Qwen, etc.)** for intelligent text cleanup and formatting.

## Features

- **CLI-First**: Lightweight command-line interface controlled by configuration files.
- **Smart Cleanup**: Algorithmic anti-hallucination guardrails (minimum duration and word count checks) ensure accurate transcription and reduce false positives.
- **Local Processing**: No data leaves your machine.
- **System-Wide**: Works in any application where you can type.
- **Context-Aware**: Adapts formatting based on the active application.
- **Cross-Platform**: Designed for macOS, Windows, and Linux.

## Project Structure

The project is written in Go and structured as follows:

```
sussurro/
├── cmd/            # Application entry points
├── internal/       # Private application and library code
│   ├── pipeline/   # Core processing pipeline
│   ├── asr/        # ASR engine
│   ├── llm/        # LLM engine
│   └── ...
├── pkg/            # Library code that's ok to use by external applications
├── models/         # Model files (gitignored)
├── scripts/        # Helper scripts
└── configs/        # Configuration files
```

## Getting Started

### Prerequisites

- Go (latest version)
- Make
- C/C++ Compiler (for building dependencies)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/cesp99/sussurro.git
   cd sussurro
   ```

2. Download Models:
   Since Sussurro runs locally, you need to download the required models first.
   ```bash
   chmod +x scripts/download-models.sh
   ./scripts/download-models.sh
   ```
   This will download Whisper (small) and TinyLlama models to the `models/` directory.

3. Build the application:
   ```bash
   make build
   ```

### Running

To run Sussurro:

```bash
./bin/sussurro
```
Or via Make:
```bash
make run
```

The application will start in the background. Use the configured hotkey (default: `Ctrl+Space`) to start recording. Release the hotkey to transcribe and inject text.

Press `Ctrl+C` in the terminal to stop the application gracefully.

## Configuration

Configuration is loaded from `configs/default.yaml`. You can customize models, audio settings, and hotkeys there.

## License

This project is licensed under the GNU General Public License v3.0 (GPLv3) - see the [LICENSE](LICENSE) file for details.
