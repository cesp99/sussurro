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

# Download Qwen 3 1.7B GGUF
# Source: HuggingFace (enacimie/Qwen3-1.7B-Q4_K_M-GGUF)
# Using Q4_K_M quantization
LLM_URL="https://huggingface.co/enacimie/Qwen3-1.7B-Q4_K_M-GGUF/resolve/main/qwen3-1.7b-q4_k_m.gguf"
LLM_FILE="$MODEL_DIR/qwen3-1.7b-q4_k_m.gguf"

echo "Downloading Qwen 3 1.7B model..."
if [ -f "$LLM_FILE" ]; then
    echo "Model already exists at $LLM_FILE"
else
    curl -L -o "$LLM_FILE" "$LLM_URL"
    if [ $? -eq 0 ]; then
        echo "Download complete: $LLM_FILE"
    else
        echo "Download failed."
        exit 1
    fi
fi

echo "Done."
