package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	PrintWidth int `json:"printWidth"`
	TabWidth   int `json:"tabWidth"`
}

var defaults = Config{
	PrintWidth: 80,
	TabWidth:   4,
}

func LoadConfig(path string) (Config, error) {
	if path == "" {
		return defaults, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse configuration: %w", err)
	}

	mergeDefaults(&config)

	return config, nil
}

func mergeDefaults(config *Config) {
	if config.PrintWidth == 0 {
		config.PrintWidth = defaults.PrintWidth
	}

	if config.TabWidth == 0 {
		config.TabWidth = defaults.TabWidth
	}
}
