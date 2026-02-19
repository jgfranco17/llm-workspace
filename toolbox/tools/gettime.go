package tools

import (
	"context"
	"time"
)

func GetCurrentTime(_ context.Context, _ map[string]string) (string, error) {
	return time.Now().Format(time.RFC3339), nil
}
