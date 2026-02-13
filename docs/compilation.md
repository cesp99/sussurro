# Compilation and Build Guide

## Prerequisites

To build Sussurro from source, you need the following tools installed on your system:

1.  **Go 1.24+**: The project uses the latest Go toolchain.
2.  **C/C++ Compiler**:
    *   **macOS**: Xcode Command Line Tools (`xcode-select --install`).
    *   **Linux**: GCC (`build-essential`).
3.  **Make**: For running build scripts.
4.  **Git**: For cloning dependencies.

## Building Dependencies

Sussurro relies on C++ libraries (`whisper.cpp` and `llama.cpp`) which must be built and linked statically.

Run the following command to download and build these dependencies:

```bash
make deps
```

This command will:
1.  Clone `whisper.cpp` into `third_party/`.
2.  Clone `go-llama.cpp` into `third_party/`.
3.  Compile the static libraries (`.a` files) optimized for your architecture (including Metal support on macOS).

## Building the Application

Once dependencies are ready, build the main binary:

```bash
make build
```

The binary will be output to `bin/sussurro`.

## Downloading Models

Before running, you must download the AI models. We provide a script for this:

```bash
./scripts/download-models.sh
```

This script downloads:
*   **Whisper Small**: For ASR.
*   **Qwen 3 1.7B (Q4_K_M)**: For text cleanup.

## Running

```bash
# Run via Make
make run

# Or run the binary directly
./bin/sussurro
```
