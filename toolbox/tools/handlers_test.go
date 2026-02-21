package tools

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerCollection_Get(t *testing.T) {
	cases := []struct {
		name         string
		collection   HandlerCollection
		key          string
		wantFound    bool
		wantErrMatch string
	}{
		{
			name:       "existing handler",
			collection: HandlerCollection{"tool": noopHandler("ok")},
			key:        "tool",
			wantFound:  true,
		},
		{
			name:         "missing handler",
			collection:   HandlerCollection{},
			key:          "missing",
			wantFound:    false,
			wantErrMatch: "handler not found: missing",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			h, err := testCase.collection.Get(testCase.key)
			if testCase.wantFound {
				assert.NoError(t, err)
				assert.NotNil(t, h)
			} else {
				assert.Error(t, err)
				assert.ErrorContains(t, err, testCase.wantErrMatch)
				assert.Nil(t, h)
			}
		})
	}
}

func TestHandlerCollection_Has(t *testing.T) {
	cases := []struct {
		name       string
		collection HandlerCollection
		key        string
		want       bool
	}{
		{
			name:       "key present",
			collection: HandlerCollection{"tool": noopHandler("ok")},
			key:        "tool",
			want:       true,
		},
		{
			name:       "key absent",
			collection: HandlerCollection{},
			key:        "tool",
			want:       false,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.want, testCase.collection.Has(testCase.key))
		})
	}
}

func TestHandlerCollection_Add(t *testing.T) {
	t.Run("adds handler and makes it retrievable", func(t *testing.T) {
		c := HandlerCollection{}
		c.Add("tool", noopHandler("ok"))

		assert.True(t, c.Has("tool"))
		h, err := c.Get("tool")
		assert.NoError(t, err)
		assert.NotNil(t, h)
	})

	t.Run("overwrites existing handler", func(t *testing.T) {
		c := HandlerCollection{"tool": noopHandler("first")}
		c.Add("tool", noopHandler("second"))

		h, err := c.Get("tool")
		assert.NoError(t, err)
		result, err := (*h)(context.Background(), Parameters{})
		assert.NoError(t, err)
		assert.Equal(t, "second", result)
	})
}

func TestDefaultHandlers(t *testing.T) {
	c := DefaultHandlers()
	assert.NotNil(t, c)

	expectedKeys := []string{
		"get_current_time",
		"fetch_url",
		"run_shell_command",
		"read_file",
		"write_file",
	}

	for _, key := range expectedKeys {
		t.Run("contains "+key, func(t *testing.T) {
			assert.True(t, c.Has(key))
		})
	}
}

// noopHandler returns a Handler that always returns the given result.
func noopHandler(result string) Handler {
	return func(_ context.Context, _ Parameters) (string, error) {
		return result, nil
	}
}
