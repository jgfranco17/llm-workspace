package tools

import (
	"context"
	"time"
)

func GetCurrentTime(_ context.Context, _ Parameters) (string, error) {
	return time.Now().Format(time.RFC3339), nil
}
