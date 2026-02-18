# Quick Start Guide

Get Sussurro running in under 5 minutes.

## Step 1: Check Your System

```bash
# Check if you're on Wayland
echo $XDG_SESSION_TYPE
```

## Step 2: Install Dependencies (Linux Only)

### If you got "wayland" above:
```bash
# Arch/Manjaro
sudo pacman -S wl-clipboard

# Ubuntu/Debian
sudo apt install wl-clipboard

# Fedora
sudo dnf install wl-clipboard
```

### If you got "x11" or you're on macOS:
Skip to Step 3. No dependencies needed.

## Step 3: Download Sussurro

Go to [GitHub Releases](https://github.com/cesp99/sussurro/releases) and download the binary for your platform:
- macOS: `sussurro-macos-<arch>.tar.gz`
- Linux: `sussurro-linux-<arch>.tar.gz`

Extract and prepare:
```bash
tar -xzf sussurro-*.tar.gz
cd sussurro-*
chmod +x sussurro trigger.sh

# macOS only: Remove quarantine
xattr -d com.apple.quarantine sussurro 2>/dev/null || true
```

## Step 4: Run for the First Time

```bash
./sussurro
```

Follow the prompts to:
1. Choose your Whisper ASR model:
   - **Whisper Small** (488 MB) — faster, good accuracy
   - **Whisper Large v3 Turbo** (1.62 GB) — slower, best accuracy
2. Download the AI models (LLM is always ~1.28 GB)
3. Wait for models to download (this takes a few minutes)

> **Tip:** You can switch the Whisper model later with `./sussurro --whisper`

## Step 5: Configure Hotkey (Wayland Only)

**Skip this if you're on X11 or macOS - hotkeys work automatically.**

If you're on Wayland, you need to set up a keyboard shortcut once:

### GNOME (Wayland)
1. Open Settings → Keyboard → Keyboard Shortcuts
2. Scroll to bottom, click "Custom Shortcuts"
3. Click "+" to add new
4. Name: `Sussurro`
5. Command: `/full/path/to/sussurro/scripts/trigger.sh`
   - Replace `/full/path/to/` with actual path (use `pwd` to find it)
6. Click "Set Shortcut" and press `Ctrl+Shift+Space`
7. Click "Add"

### KDE Plasma (Wayland)
1. System Settings → Shortcuts
2. Click "Custom Shortcuts"
3. Right-click in empty area → New → Global Shortcut → Command/URL
4. Trigger tab: Click button and press `Ctrl+Shift+Space`
5. Action tab: Enter `/full/path/to/sussurro/scripts/trigger.sh`
6. Click "Apply"

### Other DEs
See [wayland.md](wayland.md) for Sway, Hyprland, etc.

## Step 6: Test It

1. Open any text editor (gedit, kate, notepad, etc.)
2. Click inside to focus
3. **X11/macOS:** Hold `Ctrl+Shift+Space` (or `Cmd+Shift+Space` on Mac), say something, release
4. **Wayland:** Press `Ctrl+Shift+Space`, say something, press again
5. Text appears!

## Troubleshooting

### "clipboard failed" error
- **Wayland:** You forgot Step 2. Install `wl-clipboard`
- **X11:** This shouldn't happen. Check your X11 setup.

### Hotkey doesn't work (Wayland)
- Did you complete Step 5?
- Is Sussurro running? Check with `ps aux | grep sussurro`
- Test the trigger manually: `echo toggle | nc -U /run/user/$(id -u)/sussurro.sock`

### Hotkey doesn't work (X11/macOS)
- Is another app using `Ctrl+Shift+Space`? Try a different hotkey in config.
- Check logs for error messages

### No text appears
- Check that Sussurro is running
- Look at terminal for error messages
- Make sure you spoke for at least 2 seconds

## Next Steps

- Customize hotkey: Edit `~/.sussurro/config.yaml`
- Switch Whisper model: Run `./sussurro --whisper` (or `--wsp`)
- Use different models: See [Configuration Guide](configuration.md)
- Build from source: See [Compilation Guide](compilation.md)

## Daily Usage

**X11/macOS:**
1. Keep Sussurro running in background
2. Hold hotkey anywhere you can type
3. Speak
4. Release
5. Text appears

**Wayland:**
1. Keep Sussurro running in background
2. Press hotkey once to start
3. Speak
4. Press hotkey again to stop
5. Text appears

To stop Sussurro: Press `Ctrl+C` in the terminal where it's running.
