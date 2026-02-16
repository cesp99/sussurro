# Changelog

All notable changes to Sussurro will be documented in this file.

## [1.3] - 2026-02-16

### Changed
- **Upgraded LLM model** from Qwen 3 base to fine-tuned **Qwen 3 Sussurro**
- Model now hosted at https://huggingface.co/cesp99/qwen3-sussurro
- Improved transcription cleanup and accuracy with domain-specific training
- Automatic detection and migration for users upgrading from versions < v1.3
- Setup now displays file sizes for model downloads (Whisper: 488 MB, LLM: 1.28 GB)

## [1.2] - 2026-02-14

### Added
- **Full Linux support** with automatic platform detection
- **Wayland support** via trigger server and UNIX socket
- **Pure-Go clipboard** implementation (no external dependencies on X11)
- Platform-specific hotkey handlers (X11 vs Wayland)
- Trigger server for Wayland with desktop notifications
- Helper script (`scripts/trigger.sh`) for Wayland keyboard shortcuts
- Comprehensive documentation:
  - Quick Start Guide
  - Dependencies guide with distro-specific commands
  - Wayland setup guide for all major DEs
  - Platform-specific README sections
- Graceful shutdown handling (Ctrl+C now works properly)
- Parallel compilation support (multi-core builds)

### Changed
- Refactored hotkey system with platform-specific implementations
- Improved log verbosity (moved technical details to DEBUG level)
- Updated clipboard to use `github.com/atotto/clipboard` on Linux
- Build system now detects CPU cores for faster compilation
- Context providers now use build tags for platform selection

### Fixed
- macOS-specific code now properly excluded on Linux builds
- Build errors on Linux due to missing build tags
- Clipboard failures on Wayland (now requires `wl-clipboard`)
- Application hanging on shutdown
- sed syntax incompatibility in patch script (macOS vs Linux)
- Metal GPU framework attempted on Linux builds

### Documentation
- Reorganized README with platform-specific quick start sections
- Added system dependency requirements for each platform
- Clear Wayland vs X11 usage instructions
- Desktop environment-specific setup guides (GNOME, KDE, Sway, Hyprland)

## [1.1] - 2025-02-13

### Added
- Initial release
- macOS support with native hotkeys
- Whisper.cpp integration for ASR
- LLM-based text cleanup with Qwen 3
- Configuration system
- First-run setup flow

## [1.0] - 2025-02-13

- Initial development version
