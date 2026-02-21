package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

func WriteFile(ctx context.Context, params Parameters) (string, error) {
	filepathValue, err := getRequiredParam(ctx, params, "filepath")
	if err != nil {
		return "", err
	}

	content, err := getRequiredParam(ctx, params, "content")
	if err != nil {
		return "", err
	}

	dir := filepath.Dir(filepathValue)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}

	if err := os.WriteFile(filepathValue, []byte(content), 0o644); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}

	return fmt.Sprintf("Successfully wrote %d bytes to %s", len(content), filepathValue), nil
}
