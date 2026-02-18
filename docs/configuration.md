# Configuration Guide

Sussurro uses a flexible configuration system powered by [Viper](https://github.com/spf13/viper).

## Loading Mechanism

When Sussurro starts, it looks for a configuration file in the following order:

1.  **Command Line Flag**: If provided via `-config`.
    ```bash
    ./sussurro -config /path/to/my-config.yaml
    ```
2.  **Current Directory**: Checks for `config.yaml` in the directory where the binary is run.
3.  **Home Directory**: Checks for `~/.sussurro/config.yaml`.
4.  **Configs Directory**: Checks for `./configs/config.yaml`.
5.  **Fallback**: If `config.yaml` is not found, the same paths are checked for `default.yaml`.

## Configuration Structure (`config.yaml`)

The repo also includes `configs/default.yaml` with the same keys. It is a fallback if `config.yaml` is missing.

### App Settings
```yaml
app:
  name: "Sussurro"
  debug: true        # Enable verbose logging
  log_level: "info"  # debug, info, warn, error
```

### Audio Settings
```yaml
audio:
  sample_rate: 16000 # Required by Whisper
  channels: 1        # Mono audio
  bit_depth: 16
  buffer_size: 1024
  max_duration: "60s" # Maximum recording time (default: 60s, 0 for no limit)
```

### Model Settings
Sussurro requires two models: one for ASR and one for LLM cleanup.

```yaml
models:
  asr:
    path: "/home/you/.sussurro/models/ggml-small.bin"
    type: "whisper"
    threads: 4
  llm:
    path: "/home/you/.sussurro/models/qwen3-sussurro-q4_k_m.gguf" # Path to Qwen 3 model
    context_size: 32768                   # Qwen 3 supports large context
    gpu_layers: 0                         # Set > 0 if compiled with Metal or CUDA support
    threads: 4
```

Use absolute paths for model files. The first run setup writes a config file with absolute paths based on your home directory.

#### Whisper ASR Models

Two Whisper models are supported. During first-run setup you will be asked which one to download. You can also switch at any time:

```bash
sussurro --whisper   # or: sussurro --wsp
```

| Model | Filename | Size | Notes |
|-------|----------|------|-------|
| Whisper Small | `ggml-small.bin` | 488 MB | Faster, lower RAM |
| Whisper Large v3 Turbo | `ggml-large-v3-turbo.bin` | 1.62 GB | Slower, higher accuracy |

The `--whisper` / `--wsp` flag opens an interactive menu, downloads the chosen model if needed, and updates `~/.sussurro/config.yaml` automatically.

### Hotkey Settings
```yaml
hotkey:
  trigger: "ctrl+shift+space" # The key combination to hold for recording
```

### Injection Settings
```yaml
injection:
  method: "keyboard"
```

### Environment Variables

All configuration values can be overridden using environment variables prefixed with `SUSSURRO_`. Nested keys are separated by underscores.

Example:
```bash
export SUSSURRO_APP_DEBUG=true
export SUSSURRO_MODELS_LLM_THREADS=8
./sussurro
```
