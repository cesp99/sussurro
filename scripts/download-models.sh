#!/bin/bash

# Directory to store models
MODEL_DIR="models"
mkdir -p $MODEL_DIR

# Check if curl is installed
if ! command -v curl &> /dev/null; then
    echo "Error: curl is required to download models."
    exit 1
fi

# Download Whisper small model (recommended)
# Source: HuggingFace (ggerganov/whisper.cpp)
# We need the ggml-small.bin file
# The official download script usually uses: https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-small.bin

MODEL_URL="https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-small.bin"
MODEL_FILE="$MODEL_DIR/ggml-small.bin"

echo "Downloading Whisper Small model..."
if [ -f "$MODEL_FILE" ]; then
    echo "Model already exists at $MODEL_FILE"
else
    curl -L -o "$MODEL_FILE" "$MODEL_URL"
    if [ $? -eq 0 ]; then
        echo "Download complete: $MODEL_FILE"
    else
        echo "Download failed."
        exit 1
    fi
fi

# Note: LLM download will be added in Phase 4
echo "Done."
