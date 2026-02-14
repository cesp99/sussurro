# Sussurro

[![License: GPLv3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Language: Go](https://img.shields.io/badge/Language-Go-00ADD8?style=flat&logo=go)](https://golang.org)
[![Platform: macOS | Windows | Linux](https://img.shields.io/badge/Platform-macOS%20%7C%20Windows%20%7C%20Linux-blue?style=flat&logo=apple)](https://github.com/cesp99/sussurro)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen?style=flat)](https://github.com/cesp99/sussurro/actions)
[![Version](https://img.shields.io/badge/Version-1.1-blue?style=flat)](https://github.com/cesp99/sussurro)

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
- **Cross-Platform**: Designed for macOS, Windows and Linux.

## Documentation

*   [**Architecture**](docs/architecture.md): Learn how the audio pipeline, ASR, and LLM engines work together.
*   [**Configuration**](docs/configuration.md): Detailed guide on `default.yaml` and environment variables.
*   [**Compilation**](docs/compilation.md): Instructions for building from source.

## Getting Started

### Prerequisites

- Go 1.24+
- Make
- C/C++ Compiler

### Prebuilt Binaries

1.  Download the latest release for your OS from the GitHub Releases page.
2.  Unzip the archive and run:
    ```bash
    ./sussurro
    ```
    On first run Sussurro creates `~/.sussurro/config.yaml` and asks to download the models into `~/.sussurro/models`.

### Quick Install

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/cesp99/sussurro.git
    cd sussurro
    ```

2.  **Build**:
    ```bash
    make build
    ```

3.  **Run**:
    ```bash
    ./bin/sussurro
    ```
    On first run Sussurro creates `~/.sussurro/config.yaml` and asks to download the models into `~/.sussurro/models`.
    Or with a specific config:
    ```bash
    ./bin/sussurro -config /path/to/config.yaml
    ```

The application runs in the background. Hold `Ctrl+Shift+Space` to talk, release to transcribe.

## License

This project is licensed under the GNU General Public License v3.0 (GPLv3) - see the [LICENSE](LICENSE) file for details.
