# Sussurro

Sussurro is a fully local, open-source, cross-platform CLI voice-to-text system that acts as a system-wide AI dictation layer. It transforms speech into clean, formatted, context-aware text injected into any application.

## Overview

Sussurro uses local AI models to ensure privacy and low latency. It combines:
- **Whisper.cpp** for automatic speech recognition (ASR).
- **LLMs (Qwen 3)** for intelligent text cleanup, removing filler words, and fixing grammar errors.

## Features

- **CLI-First**: Lightweight command-line interface controlled by configuration files.
- **Smart Cleanup**:
    - **Filler Removal**: Automatically removes "umm", "ah", "like".
    - **Self-Correction**: Handles speech repairs (e.g., "I want blue... no red" -> "I want red").
    - **Guardrails**: Algorithmic checks to ensure accurate transcription and prevent hallucinations.
- **Local Processing**: No data leaves your machine.
- **System-Wide**: Works in any application where you can type.
- **Configurable**: Load custom configs at runtime.
- **Cross-Platform**: Designed for macOS, Windows, and Linux.

## Documentation

*   [**Architecture**](docs/architecture.md): Learn how the audio pipeline, ASR, and LLM engines work together.
*   [**Configuration**](docs/configuration.md): Detailed guide on `default.yaml` and environment variables.
*   [**Compilation**](docs/compilation.md): Instructions for building from source.

## Getting Started

### Prerequisites

- Go 1.24+
- Make
- C/C++ Compiler

### Quick Install

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/cesp99/sussurro.git
    cd sussurro
    ```

2.  **Download Models**:
    ```bash
    chmod +x scripts/download-models.sh
    ./scripts/download-models.sh
    ```
    This downloads the Whisper (small) and Qwen 3 (1.7B) models to `models/`.

3.  **Build**:
    ```bash
    make build
    ```

4.  **Run**:
    ```bash
    ./bin/sussurro
    ```
    Or with a specific config:
    ```bash
    ./bin/sussurro -config my_custom_config.yaml
    ```

The application runs in the background. Hold `Ctrl+Space` to talk, release to transcribe.

## License

This project is licensed under the GNU General Public License v3.0 (GPLv3) - see the [LICENSE](LICENSE) file for details.
