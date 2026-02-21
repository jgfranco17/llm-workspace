package config

import (
	"encoding/json"
	"fmt"
	"io"
)

const (
	DefaultConfigFile string = "tools.json"
)

type Config struct {
	Tools []ToolSchema `json:"tools"`
}

func Read(configData io.Reader) (Config, error) {
	if configData == nil {
		return Config{}, fmt.Errorf("config data is required")
	}

	data, err := io.ReadAll(configData)
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}

	if len(cfg.Tools) == 0 {
		return Config{}, fmt.Errorf("config has no tools")
	}

	return cfg, nil
}
