package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteFile(t *testing.T) {
	cases := []struct {
		name         string
		setup        func(t *testing.T) Parameters
		wantErr      bool
		wantErrMatch string
		verify       func(t *testing.T, result string, params Parameters)
	}{
		{
			name: "missing filepath param",
			setup: func(t *testing.T) Parameters {
				return Parameters{"content": "hello"}
			},
			wantErr:      true,
			wantErrMatch: "missing required parameter: filepath",
		},
		{
			name: "missing content param",
			setup: func(t *testing.T) Parameters {
				return Parameters{"filepath": filepath.Join(t.TempDir(), "out.txt")}
			},
			wantErr:      true,
			wantErrMatch: "missing required parameter: content",
		},
		{
			name: "writes content to file",
			setup: func(t *testing.T) Parameters {
				return Parameters{
					"filepath": filepath.Join(t.TempDir(), "out.txt"),
					"content":  "hello world",
				}
			},
			verify: func(t *testing.T, result string, params Parameters) {
				data, err := os.ReadFile(params["filepath"])
				assert.NoError(t, err)
				assert.Equal(t, "hello world", string(data))
				assert.Contains(t, result, "11") // byte count
			},
		},
		{
			name: "creates intermediate directories",
			setup: func(t *testing.T) Parameters {
				return Parameters{
					"filepath": filepath.Join(t.TempDir(), "a", "b", "c", "out.txt"),
					"content":  "nested",
				}
			},
			verify: func(t *testing.T, _ string, params Parameters) {
				data, err := os.ReadFile(params["filepath"])
				assert.NoError(t, err)
				assert.Equal(t, "nested", string(data))
			},
		},
		{
			name: "overwrites existing file",
			setup: func(t *testing.T) Parameters {
				path := filepath.Join(t.TempDir(), "out.txt")
				assert.NoError(t, os.WriteFile(path, []byte("old content"), 0o644))
				return Parameters{"filepath": path, "content": "new content"}
			},
			verify: func(t *testing.T, _ string, params Parameters) {
				data, err := os.ReadFile(params["filepath"])
				assert.NoError(t, err)
				assert.Equal(t, "new content", string(data))
			},
		},
		{
			name: "result message includes byte count and filepath",
			setup: func(t *testing.T) Parameters {
				path := filepath.Join(t.TempDir(), "out.txt")
				return Parameters{"filepath": path, "content": "abc"}
			},
			verify: func(t *testing.T, result string, params Parameters) {
				assert.Contains(t, result, fmt.Sprintf("%d", len("abc")))
				assert.Contains(t, result, params["filepath"])
			},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			params := testCase.setup(t)
			result, err := WriteFile(context.Background(), params)
			if testCase.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, testCase.wantErrMatch)
				return
			}
			assert.NoError(t, err)
			if testCase.verify != nil {
				testCase.verify(t, result, params)
			}
		})
	}
}
