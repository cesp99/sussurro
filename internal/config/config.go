package config

import (
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
	Version  string `mapstructure:"version"`
	Debug    bool   `mapstructure:"debug"`
	LogLevel string `mapstructure:"log_level"`
}

type AudioConfig struct {
	SampleRate int `mapstructure:"sample_rate"`
	Channels   int `mapstructure:"channels"`
	BitDepth   int `mapstructure:"bit_depth"`
	BufferSize int `mapstructure:"buffer_size"`
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

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigName("default")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("SUSSURRO")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
