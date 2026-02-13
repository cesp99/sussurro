# Configuration Guide

Sussurro uses a flexible configuration system powered by [Viper](https://github.com/spf13/viper).

## Loading Mechanism

When Sussurro starts, it looks for a configuration file in the following order:

1.  **Command Line Flag**: If provided via `-config`.
    ```bash
    ./sussurro -config /path/to/my-config.yaml
    ```
2.  **Current Directory**: Checks for `default.yaml` in the directory where the binary is run.
3.  **Configs Directory**: Checks for `./configs/default.yaml`.
4.  **Home Directory**: Checks for `~/.sussurro/default.yaml`.

## Configuration Structure (`default.yaml`)

### App Settings
```yaml
app:
  name: "sussurro"
  version: "0.1.0"
  debug: true        # Enable verbose logging
  log_level: "info"  # debug, info, warn, error
```

### Audio Settings
```yaml
audio:
  sample_rate: 16000 # Required by Whisper
  channels: 1        # Mono audio
  bit_depth: 16
  buffer_size: 4096
  max_duration: "30s" # Maximum recording time
```

### Model Settings
Sussurro requires two models: one for ASR and one for LLM cleanup.

```yaml
models:
  asr:
    path: "models/ggml-small.bin"
    type: "whisper"
    threads: 4
  llm:
    path: "models/qwen3-1.7b-q4_k_m.gguf" # Path to Qwen 3 model
    context_size: 32768                   # Qwen 3 supports large context
    gpu_layers: 0                         # Set > 0 if compiled with Metal/CUDA support
    threads: 4
```

### Hotkey Settings
```yaml
hotkey:
  trigger: "Ctrl+Shift+Space" # The key combination to hold for recording
```

### Environment Variables

All configuration values can be overridden using environment variables prefixed with `SUSSURRO_`. Nested keys are separated by underscores.

Example:
```bash
export SUSSURRO_APP_DEBUG=true
export SUSSURRO_MODELS_LLM_THREADS=8
./sussurro
```
