package tools

import (
	"context"
	"time"
)

// GetCurrentTime returns the current time formatted as RFC3339.
func GetCurrentTime(_ context.Context, _ Parameters) (string, error) {
	return time.Now().Format(time.RFC3339), nil
}
