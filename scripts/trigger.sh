#!/bin/bash
# Trigger script for Sussurro on Wayland
# Bind this script to your keyboard shortcut in your DE settings

SOCKET="${XDG_RUNTIME_DIR:-/tmp}/sussurro.sock"

if [ ! -S "$SOCKET" ]; then
    notify-send "Sussurro" "Sussurro is not running"
    exit 1
fi

echo "toggle" | nc -U "$SOCKET" 2>/dev/null || {
    # Fallback if nc is not available
    echo "toggle" | socat - UNIX-CONNECT:"$SOCKET" 2>/dev/null
}
