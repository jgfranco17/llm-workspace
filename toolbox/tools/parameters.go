package tools

import (
	"fmt"
	"strings"
)

// Parameters holds the named string arguments passed to a Handler.
type Parameters map[string]string

// Get retrieves the value of the parameter named key. It returns an
// error if the parameter is missing or empty.
func (p Parameters) Get(key string) (string, error) {
	value, ok := p[key]
	if !ok || strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("missing required parameter: %s", key)
	}
	return value, nil
}
