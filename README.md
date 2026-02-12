# Sussurro

Sussurro is a fully local, open-source, cross-platform desktop voice-to-text system that acts as a system-wide AI dictation layer. It transforms speech into clean, formatted, context-aware text injected into any application.

## Overview

Sussurro uses local AI models to ensure privacy and low latency. It combines:
- **Whisper.cpp** for automatic speech recognition (ASR).
- **LLMs (TinyLlama, Qwen, etc.)** for intelligent text cleanup and formatting.

## Features

- **Beautiful UI**: High-contrast black and white theme with a model manager.
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
│   ├── ui/         # Fyne-based GUI code
│   ├── theme/      # Custom UI theme
│   └── ...
├── pkg/            # Library code that's ok to use by external applications
├── models/         # Model files (gitignored)
└── configs/        # Configuration files
```

## Getting Started

### Prerequisites

- Go (latest version)
- Make
- C/C++ Compiler (for building dependencies)

### Building

```bash
make build
```

### Running

```bash
make run
```

On first launch, the **Model Manager** will appear, allowing you to download the required Whisper and LLM models.

## Installation

### Releases (Recommended)
You can download pre-compiled binaries for macOS, Windows, and Linux from the [Releases](https://github.com/cesp99/sussurro/releases) page. No compilation required.

### Building from Source
If you prefer to compile from scratch:

```bash
make build
```

## License

This project is licensed under the GNU General Public License v3.0 (GPLv3) - see the [LICENSE](LICENSE) file for details.

