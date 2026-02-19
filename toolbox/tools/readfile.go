package tools

import (
	"context"
	"fmt"
	"os"
)

func ReadFile(ctx context.Context, params map[string]string) (string, error) {
	filepathValue, err := getRequiredParam(ctx, params, "filepath")
	if err != nil {
		return "", err
	}

	stat, err := os.Stat(filepathValue)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("file does not exist: %s", filepathValue)
		}
		return "", fmt.Errorf("read file: %w", err)
	}

	if stat.IsDir() {
		return "", fmt.Errorf("path is not a file: %s", filepathValue)
	}

	if stat.Size() > MaxFileSize {
		return "", fmt.Errorf("file too large (max %d bytes)", MaxFileSize)
	}

	data, err := os.ReadFile(filepathValue)
	if err != nil {
		return "", fmt.Errorf("read file: %w", err)
	}

	return string(data), nil
}
