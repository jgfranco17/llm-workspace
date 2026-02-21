package tools

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// RunShellCommand executes params["command"] via sh within the directory
// given by params["workdir"], with a 30-second timeout. It returns stdout
// on success or an error wrapping stderr on failure.
func RunShellCommand(ctx context.Context, params Parameters) (string, error) {
	command, err := params.Get("command")
	if err != nil {
		return "", err
	}

	workdir := params["workdir"]
	if strings.TrimSpace(workdir) == "" {
		workdir = "/workspace"
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	cmd.Dir = workdir

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	output := stdout.String()
	if err != nil {
		output = stderr.String()
		if strings.TrimSpace(output) == "" {
			output = err.Error()
		}
		return "", fmt.Errorf("run shell command: %s", strings.TrimSpace(output))
	}

	return limitOutput(output, MaxOutputSize), nil
}
