package config

import (
	"encoding/json"
	"fmt"
	"io"
)

const (
	// DefaultConfigFile is the default filename for the tool configuration.
	DefaultConfigFile string = "tools.json"
)

// Config holds the parsed tool definitions loaded from a config file.
type Config struct {
	Tools []ToolSchema `json:"tools"`
}

// Read decodes a JSON config from r and returns the resulting Config.
// It returns an error if r is nil, the JSON is malformed, or no tools
// are defined.
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
