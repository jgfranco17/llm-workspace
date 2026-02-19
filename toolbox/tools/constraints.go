package tools

import (
	"context"
	"fmt"
	"strings"
)

const (
	MaxOutputSize = 1000
	MaxFileSize   = 100_000
)

func limitOutput(text string, limit int) string {
	if len(text) <= limit {
		return text
	}
	return text[:limit]
}

func getRequiredParam(_ context.Context, params map[string]string, key string) (string, error) {
	value, ok := params[key]
	if !ok || strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("missing required parameter: %s", key)
	}
	return value, nil
}
