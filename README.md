# Sussurro

[![Version 1.5](https://img.shields.io/badge/Version-1.5-black?style=flat)](https://github.com/cesp99/sussurro/releases)
[![GPL-3.0](https://img.shields.io/badge/License-GPL--3.0-black?style=flat)](LICENSE)
[![Go 1.24+](https://img.shields.io/badge/Go-1.24+-black?style=flat&logo=go&logoColor=white)](https://golang.org)
[![Linux](https://img.shields.io/badge/Linux-black?style=flat&logo=linux&logoColor=white)](https://github.com/cesp99/sussurro)
[![macOS](https://img.shields.io/badge/macOS-black?style=flat&logo=apple&logoColor=white)](https://github.com/cesp99/sussurro)

Sussurro is a fully local, open-source voice-to-text system with a built-in native overlay UI. It transforms speech into clean, formatted, context-aware text and injects it into any application — entirely on your machine.

**New to Sussurro?** Start with the [Quick Start Guide](docs/quickstart.md) to get running in under 5 minutes.

> **Platform note:** The native overlay UI is currently available on **Linux only**. macOS support is planned. The headless `--no-ui` mode works on both platforms.

## Overview

Sussurro uses local AI models to ensure privacy and low latency:
- **Whisper.cpp** for automatic speech recognition (ASR)
- **Qwen 3 Sussurro** (fine-tuned LLM) for intelligent text cleanup

## Features

- **Built-in Native Overlay** *(Linux)*: A minimal, aesthetically clean floating capsule shows recording/transcribing state — always on top, no taskbar entry
- **Settings UI** *(Linux)*: Dark-themed settings window accessible via system tray or right-click on the overlay
- **Smart Cleanup**: Removes filler words, handles self-corrections, prevents hallucinations
- **Local Processing**: No data leaves your machine
- **System-Wide**: Works in any application where you can type
- **Flexible ASR**: Whisper Small (fast) or Large v3 Turbo (accurate), switchable from the UI
- **Headless Mode**: `--no-ui` flag for CLI/scripting use on any platform

## Documentation

- [**Quick Start**](docs/quickstart.md): Get up and running in under 5 minutes
- [**Dependencies**](docs/dependencies.md): System requirements and package installation
- [**Wayland Setup**](docs/wayland.md): One-time configuration for Wayland users
- [**Configuration**](docs/configuration.md): Detailed guide on `config.yaml` and environment variables
- [**Architecture**](docs/architecture.md): How the audio pipeline, ASR, and LLM engines work
- [**Compilation**](docs/compilation.md): Building from source (CLI and UI builds)

## Getting Started

### Linux (Arch/Manjaro) — UI Mode

**Step 1: Install UI dependencies**
```bash
# Core UI libraries (GTK3, WebKit, AppIndicator)
sudo pacman -S gtk3 webkit2gtk-4.1 libappindicator-gtk3

# Optional: wlr-layer-shell for true Wayland overlay
sudo pacman -S gtk-layer-shell

# Wayland clipboard support
sudo pacman -S wl-clipboard

# X11 optional helpers
sudo pacman -S xdotool xorg-xprop
```

**Step 2: Get Sussurro**

Option A — prebuilt binary:
```bash
tar -xzf sussurro-linux-*.tar.gz
cd sussurro-linux-*
chmod +x sussurro
```

Option B — build from source:
```bash
git clone https://github.com/cesp99/sussurro.git
cd sussurro
make build-ui        # builds with overlay + settings UI
```

**Step 3: First run** — open a terminal and run:
```bash
./sussurro        # prebuilt
# or
./bin/sussurro    # built from source
```
Follow the prompts to download the AI models.

**Step 4 (Wayland only):** Configure a keyboard shortcut — see [Wayland Setup](docs/wayland.md).

---

### Linux (Ubuntu/Debian) — UI Mode

**Step 1: Install UI dependencies**
```bash
sudo apt install libgtk-3-0 libwebkit2gtk-4.1-0 libayatana-appindicator3-1

# Optional: wlr-layer-shell overlay (Ubuntu 22.04+)
sudo apt install libgtk-layer-shell0

# Wayland clipboard
sudo apt install wl-clipboard
```

**Step 2–4:** Same as Arch above (use `make build-ui`).

---

### macOS — Headless Mode Only

The UI overlay is **not yet available on macOS**. Run Sussurro in headless mode:

```bash
tar -xzf sussurro-macos-*.tar.gz
cd sussurro-macos-*
chmod +x sussurro
xattr -d com.apple.quarantine sussurro   # remove quarantine
./sussurro --no-ui
```

**Usage:** Hold `Cmd+Shift+Space` to talk, release to transcribe. Cleaned text is injected into the active application.

---

## UI: The Overlay Capsule *(Linux)*

When Sussurro runs on Linux, a sleek pill-shaped capsule appears at the bottom-center of your screen:

| State | Appearance |
|-------|-----------|
| **Idle** | 7 softly pulsing white dots |
| **Recording** | 7 waveform bars animated by your voice |
| **Transcribing** | "transcribing" text with a shimmer effect |

**Accessing Settings:**

| Method | How |
|--------|-----|
| System tray | Click the Sussurro icon → **Open Settings** |
| Right-click overlay | Right-click the capsule → **Open Settings** |

The settings window lets you switch Whisper models, download models with a progress bar, change the global hotkey, and toggle auto-start.

---

## Headless / CLI Mode

If you don't want the overlay (e.g. for scripting or low-resource environments):

```bash
./sussurro --no-ui
```

This runs Sussurro exactly as before — terminal output only, no overlay, no tray.

---

## Known Limitations

### "Start at Login" toggle

The "Start at Login" toggle in Settings is present in the UI but is not yet implemented. It will be addressed in a future release.

---

## Quick Reference

| Platform | Hotkey | Access Settings |
|----------|--------|----------------|
| Linux X11 | Hold `Ctrl+Shift+Space` | System tray or right-click capsule |
| Linux Wayland | Toggle (press twice) | System tray or right-click capsule |
| macOS *(headless)* | Hold `Cmd+Shift+Space` | — (no UI yet) |

## Switching Whisper Models

Via the Settings UI (recommended) — or from the command line:

```bash
./sussurro --whisper   # (or --wsp)
```

| Model | Size | Best for |
|-------|------|----------|
| Whisper Small | ~488 MB | Faster, lower RAM |
| Whisper Large v3 Turbo | ~1.62 GB | Higher accuracy |

## License

GNU General Public License v3.0 — see [LICENSE](LICENSE).
