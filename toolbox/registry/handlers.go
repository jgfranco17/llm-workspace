package registry

import "github.com/jgfranco17/llm-workspace/toolbox/tools"

func DefaultHandlers() map[string]Handler {
	return map[string]Handler{
		"get_current_time":  tools.GetCurrentTime,
		"fetch_url":         tools.FetchURL,
		"run_shell_command": tools.RunShellCommand,
		"read_file":         tools.ReadFile,
		"write_file":        tools.WriteFile,
	}
}
