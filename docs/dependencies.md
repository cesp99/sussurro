# System Dependencies

This guide lists all system packages required to run Sussurro on your platform.

## Quick Check: What Do I Need?

**macOS:** Nothing. All dependencies are built-in.

**Linux X11:** Nothing required. Optional packages improve functionality.

**Linux Wayland:** **REQUIRED:** `wl-clipboard` package (install before first run)

---

## Linux (Wayland)

### Required Packages

**CRITICAL:** Wayland users MUST install a clipboard manager before running Sussurro:

#### Arch Linux / Manjaro
```bash
sudo pacman -S wl-clipboard
```

#### Ubuntu / Debian
```bash
sudo apt install wl-clipboard
```

#### Fedora
```bash
sudo dnf install wl-clipboard
```

#### openSUSE
```bash
sudo zypper install wl-clipboard
```

### Optional Packages

For window context detection (to get active application name):
```bash
# Most DEs include these by default, but if missing:
sudo pacman -S xdg-utils          # Arch
sudo apt install xdg-utils         # Ubuntu/Debian
```

For desktop notifications:
```bash
sudo pacman -S libnotify           # Arch
sudo apt install libnotify-bin     # Ubuntu/Debian
```

For the trigger script (usually pre-installed):
```bash
# One of these is needed:
sudo pacman -S openbsd-netcat      # Arch (nc command)
sudo apt install netcat-openbsd    # Ubuntu/Debian
# OR
sudo pacman -S socat               # Alternative
```

## Linux (X11)

### Required Packages

X11 users don't need additional clipboard tools (uses X11 clipboard directly).

### Optional Packages

For window context detection:
```bash
sudo pacman -S xdotool xorg-xprop  # Arch
sudo apt install xdotool x11-utils # Ubuntu/Debian
```

## macOS

No additional dependencies required. Sussurro uses native macOS APIs for:
- Clipboard (NSPasteboard)
- Window detection (CoreGraphics)
- Global hotkeys (Carbon)

## Windows

Not yet tested...

## Desktop Environment Specific Notes

### GNOME (Wayland)
- **Required**: `wl-clipboard`
- **Included by default**: Desktop notifications, xdg-utils

### KDE Plasma (Wayland)
- **Required**: `wl-clipboard`
- Usually includes klipper (clipboard manager) but wl-clipboard is still needed for CLI access

### Sway
- **Required**: `wl-clipboard`
- **Recommended**: Install `mako` or `dunst` for notifications

### Hyprland
- **Required**: `wl-clipboard`
- **Recommended**: Install `mako` or `dunst` for notifications

### GNOME (X11)
- No additional requirements beyond base X11 packages

### KDE Plasma (X11)
- No additional requirements beyond base X11 packages

## Troubleshooting

### "Failed to write to clipboard"

**On Wayland:**
```bash
# Check if wl-clipboard is installed
which wl-copy
# If not found, install it (see above)

# Test clipboard manually
echo "test" | wl-copy
wl-paste
```

**On X11:**
This should work out of the box. If you see this error, your X11 setup might be incomplete.

### "Context provider" errors

**On Wayland:**
Sussurro may not be able to detect the active window on all Wayland compositors due to security restrictions. This is a Wayland limitation, not a bug.

**On X11:**
Install xdotool and xprop:
```bash
sudo pacman -S xdotool xorg-xprop  # Arch
sudo apt install xdotool x11-utils # Ubuntu/Debian
```

### Trigger script not working (Wayland only)

Make sure you have `nc` (netcat) or `socat` installed:
```bash
# Test if netcat works
echo "test" | nc -U /tmp/test.sock 2>&1

# If not, install:
sudo pacman -S openbsd-netcat  # Arch
sudo apt install netcat-openbsd # Ubuntu/Debian
```

## Build Dependencies

These are only needed if building from source:

### All Platforms
- Go 1.24+
- Make
- CMake 3.15+
- C/C++ compiler (gcc/clang/MSVC)

### Linux
- Build essentials:
  ```bash
  sudo pacman -S base-devel cmake    # Arch
  sudo apt install build-essential cmake  # Ubuntu/Debian
  ```

### macOS
- Xcode Command Line Tools:
  ```bash
  xcode-select --install
  ```

### Windows
- MinGW-w64 or Visual Studio Build Tools
- CMake

## Verifying Dependencies

Quick check script:
```bash
# On Wayland
which wl-copy && echo "Clipboard: OK" || echo "Clipboard: MISSING - install wl-clipboard"
which nc && echo "Netcat: OK" || echo "Netcat: MISSING - install netcat"
which notify-send && echo "Notifications: OK" || echo "Notifications: MISSING (optional)"

# On X11
which xdotool && echo "xdotool: OK" || echo "xdotool: MISSING (optional)"
which xprop && echo "xprop: OK" || echo "xprop: MISSING (optional)"
```
