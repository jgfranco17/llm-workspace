package tools

import (
	"context"
	"fmt"
	"strings"
)

const (
	// MaxOutputSize is the maximum number of bytes returned by a tool.
	MaxOutputSize = 1000

	// MaxFileSize is the maximum file size in bytes that ReadFile will accept.
	MaxFileSize = 100_000
)

func limitOutput(text string, limit int) string {
	if len(text) <= limit {
		return text
	}
	return text[:limit]
}

func getRequiredParam(_ context.Context, params Parameters, key string) (string, error) {
	value, ok := params[key]
	if !ok || strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("missing required parameter: %s", key)
	}
	return value, nil
}
