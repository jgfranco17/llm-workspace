package tools

import (
	"context"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// requireShell skips the test on platforms where sh is not available.
func requireShell(t *testing.T) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("sh not available on Windows")
	}
}

func TestRunShellCommand_MissingParam(t *testing.T) {
	cases := []struct {
		name         string
		params       Parameters
		wantErrMatch string
	}{
		{
			name:         "missing command param",
			params:       Parameters{},
			wantErrMatch: "missing required parameter: command",
		},
		{
			name:         "empty command param",
			params:       Parameters{"command": ""},
			wantErrMatch: "missing required parameter: command",
		},
		{
			name:         "whitespace command param",
			params:       Parameters{"command": "   "},
			wantErrMatch: "missing required parameter: command",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := RunShellCommand(context.Background(), testCase.params)
			assert.Error(t, err)
			assert.ErrorContains(t, err, testCase.wantErrMatch)
		})
	}
}

func TestRunShellCommand_Success(t *testing.T) {
	requireShell(t)

	cases := []struct {
		name       string
		command    string
		wantResult string
	}{
		{
			name:       "echo command",
			command:    "echo hello",
			wantResult: "hello",
		},
		{
			name:       "multi-word echo",
			command:    "echo foo bar",
			wantResult: "foo bar",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := RunShellCommand(context.Background(), Parameters{
				"command": testCase.command,
				"workdir": t.TempDir(),
			})
			assert.NoError(t, err)
			assert.Equal(t, testCase.wantResult, strings.TrimRight(result, "\n"))
		})
	}
}

func TestRunShellCommand_Failure(t *testing.T) {
	requireShell(t)

	_, err := RunShellCommand(context.Background(), Parameters{
		"command": "exit 1",
		"workdir": t.TempDir(),
	})
	assert.Error(t, err)
	assert.ErrorContains(t, err, "run shell command")
}

func TestRunShellCommand_OutputTruncation(t *testing.T) {
	requireShell(t)

	// Generate output larger than MaxOutputSize using printf and a loop.
	command := "printf '%0.sa' {1..2000}"
	result, err := RunShellCommand(context.Background(), Parameters{
		"command": command,
		"workdir": t.TempDir(),
	})
	assert.NoError(t, err)
	assert.LessOrEqual(t, len(result), MaxOutputSize)
}

func TestRunShellCommand_ContextCancelled(t *testing.T) {
	requireShell(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := RunShellCommand(ctx, Parameters{
		"command": "echo hello",
		"workdir": t.TempDir(),
	})
	assert.Error(t, err)
}
