package tools

import (
	"context"
	"fmt"
)

type Handler func(ctx context.Context, params Parameters) (string, error)

type HandlerCollection map[string]Handler

func (c HandlerCollection) GetHandler(name string) (*Handler, error) {
	if handler, ok := c[name]; ok {
		return &handler, nil
	}
	return nil, fmt.Errorf("handler not found: %s", name)
}

func DefaultHandlers() HandlerCollection {
	return HandlerCollection{
		"get_current_time":  GetCurrentTime,
		"fetch_url":         FetchURL,
		"run_shell_command": RunShellCommand,
		"read_file":         ReadFile,
		"write_file":        WriteFile,
	}
}
