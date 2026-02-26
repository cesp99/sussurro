# Known Issues

## Platform Support

### Windows not supported
The project currently only runs on Linux and macOS. Windows support is not yet implemented. The main blockers are:

- The overlay window uses GTK3 + Cairo (Linux) and NSPanel/CoreGraphics (macOS) — neither works on Windows.
- The settings window relies on `webview_go` with `webkit2gtk` (Linux) and `WKWebView` (macOS).
- The system tray uses `getlantern/systray`, which does support Windows, but has not been tested.
- Global hotkey capture on X11 uses a custom GDK `XGrabKey` filter in `overlay_linux.c`; a Windows equivalent (e.g., `RegisterHotKey`) would need to be written.
- Audio capture and text injection would also need Windows-specific implementations or a cross-platform abstraction.
