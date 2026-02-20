package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App       AppConfig       `mapstructure:"app"`
	Audio     AudioConfig     `mapstructure:"audio"`
	Models    ModelsConfig    `mapstructure:"models"`
	Hotkey    HotkeyConfig    `mapstructure:"hotkey"`
	Injection InjectionConfig `mapstructure:"injection"`
}

type AppConfig struct {
	Name     string `mapstructure:"name"`
	Debug    bool   `mapstructure:"debug"`
	LogLevel string `mapstructure:"log_level"`
}

type AudioConfig struct {
	SampleRate  int    `mapstructure:"sample_rate"`
	Channels    int    `mapstructure:"channels"`
	BitDepth    int    `mapstructure:"bit_depth"`
	BufferSize  int    `mapstructure:"buffer_size"`
	MaxDuration string `mapstructure:"max_duration"`
}

type ModelsConfig struct {
	ASR ASRConfig `mapstructure:"asr"`
	LLM LLMConfig `mapstructure:"llm"`
}

type ASRConfig struct {
	Path    string `mapstructure:"path"`
	Type    string `mapstructure:"type"`
	Threads int    `mapstructure:"threads"`
}

type LLMConfig struct {
	Path        string `mapstructure:"path"`
	ContextSize int    `mapstructure:"context_size"`
	GpuLayers   int    `mapstructure:"gpu_layers"`
	Threads     int    `mapstructure:"threads"`
}

type HotkeyConfig struct {
	Trigger string `mapstructure:"trigger"`
}

type InjectionConfig struct {
	Method string `mapstructure:"method"`
}

// SaveHotkey rewrites only the hotkey.trigger field in the YAML config file.
func SaveHotkey(cfg *Config, trigger string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot find home directory: %w", err)
	}
	configFile := filepath.Join(homeDir, ".sussurro", "config.yaml")

	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("cannot read config file: %w", err)
	}

	// Simple line-by-line replacement of the trigger value.
	lines := strings.Split(string(data), "\n")
	replaced := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "trigger:") {
			indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
			lines[i] = indent + "trigger: \"" + trigger + "\""
			replaced = true
			break
		}
	}
	if !replaced {
		return fmt.Errorf("trigger key not found in config file")
	}

	return os.WriteFile(configFile, []byte(strings.Join(lines, "\n")), 0644)
}

func LoadConfig(path string) (*Config, error) {
	if path != "" {
		// If a specific file path is provided, use it directly
		viper.SetConfigFile(path)
	} else {
		// Otherwise search in default locations
		viper.SetConfigName("config") // Look for config.yaml (or .json, .toml)
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.sussurro")
		viper.AddConfigPath("./configs")
	}

	viper.SetEnvPrefix("SUSSURRO")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Try fallback to "default" (old behavior)
			viper.SetConfigName("default")
			if err := viper.ReadInConfig(); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
