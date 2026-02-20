package tools

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFile(t *testing.T) {
	cases := []struct {
		name         string
		setup        func(t *testing.T) Parameters
		wantErr      bool
		wantErrMatch string
		wantResult   string
	}{
		{
			name: "missing filepath param",
			setup: func(t *testing.T) Parameters {
				return Parameters{}
			},
			wantErr:      true,
			wantErrMatch: "missing required parameter: filepath",
		},
		{
			name: "empty filepath param",
			setup: func(t *testing.T) Parameters {
				return Parameters{"filepath": ""}
			},
			wantErr:      true,
			wantErrMatch: "missing required parameter: filepath",
		},
		{
			name: "file does not exist",
			setup: func(t *testing.T) Parameters {
				return Parameters{"filepath": filepath.Join(t.TempDir(), "nonexistent.txt")}
			},
			wantErr:      true,
			wantErrMatch: "file does not exist",
		},
		{
			name: "path is a directory",
			setup: func(t *testing.T) Parameters {
				return Parameters{"filepath": t.TempDir()}
			},
			wantErr:      true,
			wantErrMatch: "path is not a file",
		},
		{
			name: "file too large",
			setup: func(t *testing.T) Parameters {
				path := filepath.Join(t.TempDir(), "big.txt")
				data := strings.Repeat("x", MaxFileSize+1)
				assert.NoError(t, os.WriteFile(path, []byte(data), 0o644))
				return Parameters{"filepath": path}
			},
			wantErr:      true,
			wantErrMatch: "file too large",
		},
		{
			name: "reads file content successfully",
			setup: func(t *testing.T) Parameters {
				path := filepath.Join(t.TempDir(), "hello.txt")
				assert.NoError(t, os.WriteFile(path, []byte("hello world"), 0o644))
				return Parameters{"filepath": path}
			},
			wantErr:    false,
			wantResult: "hello world",
		},
		{
			name: "reads empty file",
			setup: func(t *testing.T) Parameters {
				path := filepath.Join(t.TempDir(), "empty.txt")
				assert.NoError(t, os.WriteFile(path, []byte(""), 0o644))
				return Parameters{"filepath": path}
			},
			wantErr:    false,
			wantResult: "",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			params := testCase.setup(t)
			result, err := ReadFile(context.Background(), params)
			if testCase.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, testCase.wantErrMatch)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, testCase.wantResult, result)
		})
	}
}
