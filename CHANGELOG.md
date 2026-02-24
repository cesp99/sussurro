# Changelog

All notable changes to Sussurro will be documented in this file.

## [1.6] - 2026-02-24

### Added
- **macOS overlay blur + border**: the capsule overlay now uses `NSVisualEffectView` (material `HUDWindow`, `NSAppearanceNameVibrantDark`) as a frosted-glass backdrop clipped to the pill silhouette, making it legible over any background. A 1.5 px semi-transparent white border is drawn as an inset stroke around the pill on both macOS and Linux.
- **Model-switch restart banner**: switching the active Whisper model in Settings no longer force-quits and relaunches the process. Instead, the config is saved silently and a blue info banner — *"Restart Sussurro to load the new model into memory"* — appears at the bottom of the settings window. The running pipeline is not disrupted.
- **In-memory config sync after model switch**: after `setup.SetActiveModel` writes the new ASR path to disk, `mgr.cfg.Models.ASR.Path` is updated in memory immediately. This fixes a race where `reloadSettings()` would read stale data and snap the UI back to the previously active model for one frame.

### Fixed
- **`onDownloadProgress` fragile name match**: download progress updates now target `#prog-<modelId>` / `#pct-<modelId>` directly by element ID instead of scanning all `.model-name` spans for a matching first word — removes a latent bug if two models share a first word.
- **`onTrayExit` no-op**: the systray exit callback now calls `m.Quit()` so the `quitCh` is closed and `processUpdates` goroutine drains cleanly when the OS removes the tray icon.
- **`sussurroModelsDir()` helper**: the `~/.sussurro/models` path was duplicated in `buildInitialData` and `resolveModelDownload`; both now call a single `sussurroModelsDir()` helper.
- **Removed stale `time` import** from `settings_bridge.go` after the auto-restart goroutine was deleted.

## [1.5] - 2026-02-23

### Added
- **macOS full overlay UI**: NSPanel overlay (Cocoa + CoreVideo CVDisplayLink), settings window, system tray, and right-click context menu now all work on macOS (previously macOS was headless-only)
- **Live hotkey reconfiguration**: changing the global hotkey in Settings takes effect immediately — no restart required (`reinstallOverlayHotkey` on both Linux and macOS)
- **Linux X11 modifier support**: `alt`/`option` (X11 Mod1) and `super`/`meta`/`cmd` (X11 Mod4) hotkey modifiers now work on Linux (previously returned an error)
- **macOS modifier aliases**: `super` and `meta` are now accepted as aliases for `cmd`/`command` on macOS
- **Hotkey recording modal**: live key-combo preview as keys are held; finalises on full key release; requires at least one non-modifier key; supports up to 3 simultaneous keys
- **Metal-safe exit on macOS**: `platformExit()` calls `overlay_terminate_macos()` which stops `CVDisplayLink` and calls `_exit(0)` to bypass C++ global destructors, preventing a Metal render-encoder assertion from `ggml-metal` on quit
- **macOS settings window close fix**: `NSWindowDelegate` now hides the window instead of destroying it, preserving the WKWebView backing store across open/close cycles
- **`ParseTrigger` exported** from `internal/hotkey` package so platform-specific UI code can reuse the modifier/key mapping without duplication

### Changed
- macOS overlay panel: window level raised to `NSStatusWindowLevel`, `hidesOnDeactivate=NO`, `NSWindowCollectionBehaviorFullScreenAuxiliary` (stays visible above full-screen apps), uses `orderFrontRegardless` instead of `makeKeyAndOrderFront` to avoid stealing keyboard focus
- macOS hotkey now registered via CGEventTap in a goroutine after `[NSApp run]` is live (300 ms defer), replacing the previous no-op stub
- `Manager.Quit()` uses `platformExit()` instead of `os.Exit(0)` directly
- Log message: `"X11/macOS detected - using overlay hotkey"` → `"Using overlay hotkey"`
- Log message: `"X11 detected - using global hotkeys"` → `"Using global hotkeys (X11 / macOS)"`

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
