package tools

import (
	"context"
	"fmt"
	"os"
)

func ReadFile(ctx context.Context, params Parameters) (string, error) {
	filepathValue, err := getRequiredParam(ctx, params, "filepath")
	if err != nil {
		return "", err
	}
	if err := validateFile(filepathValue); err != nil {
		return "", err
	}

	data, err := os.ReadFile(filepathValue)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(data), nil
}

func validateFile(filepath string) error {
	stat, err := os.Stat(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", filepath)
		}
		return fmt.Errorf("failed to stat file: %w", err)
	}
	if stat.IsDir() {
		return fmt.Errorf("path is not a file: %s", filepath)
	}
	if stat.Size() > MaxFileSize {
		return fmt.Errorf("file too large (max %d bytes)", MaxFileSize)
	}
	return nil
}
