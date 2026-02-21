package tools

import (
	"context"
	"fmt"
)

// Handler is a function that executes a tool given a context and parameters.
type Handler func(ctx context.Context, params Parameters) (string, error)

// HandlerCollection maps tool names to their Handler implementations.
type HandlerCollection map[string]Handler

// Get returns a pointer to the Handler registered under name.
// It returns an error if no handler with that name exists.
func (c HandlerCollection) Get(name string) (*Handler, error) {
	if handler, ok := c[name]; ok {
		return &handler, nil
	}
	return nil, fmt.Errorf("handler not found: %s", name)
}

// Has reports whether a handler is registered under name.
func (c HandlerCollection) Has(name string) bool {
	_, ok := c[name]
	return ok
}

// Add registers handler under name in the collection.
func (c HandlerCollection) Add(name string, handler Handler) {
	if c == nil {
		c = make(HandlerCollection)
	}
	c[name] = handler
}

// DefaultHandlers returns a HandlerCollection pre-populated with
// all built-in tool handlers.
func DefaultHandlers() HandlerCollection {
	return HandlerCollection{
		"get_current_time":  GetCurrentTime,
		"fetch_url":         FetchURL,
		"run_shell_command": RunShellCommand,
		"read_file":         ReadFile,
		"write_file":        WriteFile,
	}
}
