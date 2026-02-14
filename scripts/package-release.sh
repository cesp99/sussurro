#!/bin/bash
# Package Sussurro for release
# Usage: ./scripts/package-release.sh <version> <platform> <arch>
# Example: ./scripts/package-release.sh 1.2 linux amd64

set -e

VERSION=${1:-"1.2"}
PLATFORM=${2:-"linux"}
ARCH=${3:-"amd64"}

RELEASE_NAME="sussurro-${PLATFORM}-${ARCH}"
RELEASE_DIR="release/${RELEASE_NAME}"

echo "Packaging Sussurro v${VERSION} for ${PLATFORM}-${ARCH}..."

# Clean and create release directory
rm -rf release
mkdir -p "${RELEASE_DIR}"

# Check if binary exists
if [ ! -f "bin/sussurro" ]; then
    echo "Error: bin/sussurro not found. Run 'make build' first."
    exit 1
fi

# Copy binary
echo "Copying binary..."
cp bin/sussurro "${RELEASE_DIR}/sussurro"
chmod +x "${RELEASE_DIR}/sussurro"

# Copy scripts
echo "Copying scripts..."
cp scripts/trigger.sh "${RELEASE_DIR}/trigger.sh"
chmod +x "${RELEASE_DIR}/trigger.sh"

# Copy example config
echo "Copying example config..."
cp configs/default.yaml "${RELEASE_DIR}/config.example.yaml"

# Create a quick install guide
cat > "${RELEASE_DIR}/INSTALL.txt" << 'EOF'
Sussurro v${VERSION} Installation
================================

Quick Start:
1. Make executable: chmod +x sussurro trigger.sh
2. macOS only: xattr -d com.apple.quarantine sussurro
3. Run: ./sussurro
4. Follow the prompts to download AI models
5. Set up keyboard shortcut (Wayland users only, see below)

For Wayland Users:
-----------------
If you're on Wayland (check with: echo $XDG_SESSION_TYPE):

1. Make sure you have wl-clipboard installed:
   Arch: sudo pacman -S wl-clipboard
   Ubuntu: sudo apt install wl-clipboard

2. Set up keyboard shortcut in your desktop environment:
   - Open keyboard settings
   - Add custom shortcut: Ctrl+Shift+Space
   - Command: /full/path/to/trigger.sh
   - See full guide: https://github.com/cesp99/sussurro/blob/master/docs/WAYLAND.md

For X11 Users:
-------------
Just run ./sussurro - hotkeys work automatically!
Hold Ctrl+Shift+Space to talk, release to transcribe.

Documentation:
-------------
Full documentation: https://github.com/cesp99/sussurro
Quick Start Guide: https://github.com/cesp99/sussurro/blob/master/docs/QUICKSTART.md
EOF

# Replace version placeholder (compatible with both GNU and BSD sed)
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "s/\${VERSION}/${VERSION}/g" "${RELEASE_DIR}/INSTALL.txt"
else
    sed -i "s/\${VERSION}/${VERSION}/g" "${RELEASE_DIR}/INSTALL.txt"
fi

# Create tarball
echo "Creating tarball..."
cd release
tar -czf "${RELEASE_NAME}.tar.gz" "${RELEASE_NAME}/"
cd ..

# Create checksum
echo "Generating checksum..."
cd release
if command -v sha256sum &> /dev/null; then
    sha256sum "${RELEASE_NAME}.tar.gz" > "${RELEASE_NAME}.tar.gz.sha256"
elif command -v shasum &> /dev/null; then
    shasum -a 256 "${RELEASE_NAME}.tar.gz" > "${RELEASE_NAME}.tar.gz.sha256"
else
    echo "Warning: sha256sum or shasum not found. Skipping checksum generation."
fi
cd ..

echo ""
echo "✓ Release package created successfully!"
echo ""
echo "Package: release/${RELEASE_NAME}.tar.gz"
echo "SHA256: release/${RELEASE_NAME}.tar.gz.sha256"
echo ""
echo "Contents:"
ls -lh "release/${RELEASE_NAME}/"
echo ""
echo "Upload these files to GitHub Releases:"
echo "  - release/${RELEASE_NAME}.tar.gz"
echo "  - release/${RELEASE_NAME}.tar.gz.sha256"
