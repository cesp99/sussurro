# Sussurro

[![Version 1.3](https://img.shields.io/badge/Version-1.3-black?style=flat)](https://github.com/cesp99/sussurro/releases)
[![GPL-3.0](https://img.shields.io/badge/License-GPL--3.0-black?style=flat)](LICENSE)
[![Go 1.24+](https://img.shields.io/badge/Go-1.24+-black?style=flat&logo=go&logoColor=white)](https://golang.org)
[![Linux](https://img.shields.io/badge/Linux-black?style=flat&logo=linux&logoColor=white)](https://github.com/cesp99/sussurro)
[![macOS](https://img.shields.io/badge/macOS-black?style=flat&logo=apple&logoColor=white)](https://github.com/cesp99/sussurro)

Sussurro is a fully local, open-source, cross-platform CLI voice-to-text system that acts as a system-wide AI dictation layer. It transforms speech into clean, formatted, context-aware text injected into any application.

**New to Sussurro?** Start with the [Quick Start Guide](docs/quickstart.md) to get running in under 5 minutes.

## Overview

Sussurro uses local AI models to ensure privacy and low latency. It combines:
- **Whisper.cpp** for automatic speech recognition (ASR).
- **LLMs (Qwen 3 Sussurro)** for intelligent text cleanup, removing filler words, and fixing grammar errors.

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

*   [**Quick Start**](docs/quickstart.md): Get up and running in under 5 minutes (recommended for first-time users).
*   [**Dependencies**](docs/dependencies.md): System requirements and package installation for your platform.
*   [**Wayland Setup**](docs/wayland.md): Setup guide for Wayland users (one-time configuration required).
*   [**Configuration**](docs/configuration.md): Detailed guide on `default.yaml` and environment variables.
*   [**Architecture**](docs/architecture.md): Learn how the audio pipeline, ASR, and LLM engines work together.
*   [**Compilation**](docs/compilation.md): Instructions for building from source.

## Getting Started

Choose your platform below for specific instructions:

### macOS

**Using Prebuilt Binary:**
1. Download `sussurro-macos-<arch>.tar.gz` from the [GitHub Releases](https://github.com/cesp99/sussurro/releases) page
2. Extract and prepare:
   ```bash
   tar -xzf sussurro-macos-*.tar.gz
   cd sussurro-macos-<arch>
   chmod +x sussurro trigger.sh
   xattr -d com.apple.quarantine sussurro  # Remove macOS quarantine
   ```
3. Run:
   ```bash
   ./sussurro
   ```
   On first run, Sussurro will guide you through model download

**Usage:** Hold `Cmd+Shift+Space` to talk, release to transcribe. Works immediately, no configuration needed.

**Building from Source:** See [Compilation Guide](docs/compilation.md)

---

### Linux (Arch/Manjaro)

**Step 1: Install Dependencies**
```bash
# For Wayland users (GNOME Wayland, KDE Wayland, Sway, Hyprland)
sudo pacman -S wl-clipboard

# For X11 users (optional, for window context detection)
sudo pacman -S xdotool xorg-xprop
```

**Step 2: Get Sussurro**

**Option A: Prebuilt Binary**
1. Download `sussurro-linux-<arch>.tar.gz` from [GitHub Releases](https://github.com/cesp99/sussurro/releases)
2. Extract and prepare:
   ```bash
   tar -xzf sussurro-linux-*.tar.gz
   cd sussurro-linux-<arch>
   chmod +x sussurro trigger.sh
   ```
   Package includes: `sussurro` binary, `trigger.sh` script (for Wayland), example config
3. Run:
   ```bash
   ./sussurro
   ```

**Option B: Build from Source**
```bash
git clone https://github.com/cesp99/sussurro.git
cd sussurro
make build
./bin/sussurro
```

**Step 3: First Run Setup**
- Sussurro will create `~/.sussurro/config.yaml`
- Follow prompts to download AI models

**Step 4: Configure Keyboard Shortcut (Wayland Only)**

If using **Wayland** (check with `echo $XDG_SESSION_TYPE`):

**GNOME:**
1. Settings → Keyboard → Keyboard Shortcuts → Custom Shortcuts
2. Click "+" to add new
3. Name: "Sussurro Voice Input"
4. Command: `/full/path/to/extracted/folder/trigger.sh`
5. Set shortcut: `Ctrl+Shift+Space`

**KDE Plasma:**
1. System Settings → Shortcuts → Custom Shortcuts
2. Right-click → New → Global Shortcut → Command/URL
3. Trigger tab: Set `Ctrl+Shift+Space`
4. Action tab: `/full/path/to/extracted/folder/trigger.sh`

**Sway/Hyprland:** See [Wayland Setup Guide](docs/wayland.md)

**Usage:**
- **X11**: Hold `Ctrl+Shift+Space`, talk, release (works immediately)
- **Wayland**: Press once to start, talk, press again to stop

---

### Linux (Ubuntu/Debian)

**Step 1: Install Dependencies**
```bash
# For Wayland users
sudo apt install wl-clipboard

# For X11 users (optional)
sudo apt install xdotool x11-utils
```

**Step 2-4:** Same as Arch instructions above

---

### Linux (Other Distributions)

See [dependencies.md](docs/dependencies.md) for your distribution's package manager commands, then follow the Arch instructions above.

---

### Windows

**Status:** Not yet tested. Contributions welcome.

---

## Quick Reference

| Platform | Hotkey Behavior | Setup Required |
|----------|----------------|----------------|
| macOS | Hold to talk | None |
| Linux X11 | Hold to talk | None |
| Linux Wayland | Toggle (press twice) | One-time DE shortcut |

**Troubleshooting:** See [dependencies.md](docs/dependencies.md)

## License

This project is licensed under the GNU General Public License v3.0 (GPLv3) - see the [LICENSE](LICENSE) file for details.
