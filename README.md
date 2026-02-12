# Sussurro

Sussurro is a fully local, open-source, cross-platform desktop voice-to-text system that acts as a system-wide AI dictation layer. It transforms speech into clean, formatted, context-aware text injected into any application.

## Overview

Sussurro uses local AI models to ensure privacy and low latency. It combines:
- **Whisper.cpp** for automatic speech recognition (ASR).
- **Qwen2.5-1.5B** for intelligent text cleanup and formatting.

## Features

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
├── pkg/            # Library code that's ok to use by external applications
├── models/         # Model files (gitignored)
└── configs/        # Configuration files
```

## Getting Started

### Prerequisites

- Go (latest version)
- Make

### Building

```bash
make build
```

### Running

```bash
make run
```

## Development Status

### Phase 1: Project Foundation (Completed)
- [x] Initialize Go module and project structure
- [x] Set up configuration system (YAML-based)
- [x] Implement logging infrastructure
- [x] Create build system (Makefile)

