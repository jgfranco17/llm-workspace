package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	cases := []struct {
		name         string
		input        string
		nilReader    bool
		wantErr      bool
		wantErrMatch string
		wantCount    int
	}{
		{
			name:         "nil reader",
			nilReader:    true,
			wantErr:      true,
			wantErrMatch: "config data is required",
		},
		{
			name:         "empty input",
			input:        "",
			wantErr:      true,
			wantErrMatch: "parse config",
		},
		{
			name:         "invalid JSON",
			input:        `{not valid json}`,
			wantErr:      true,
			wantErrMatch: "parse config",
		},
		{
			name:         "valid JSON but no tools",
			input:        `{"tools": []}`,
			wantErr:      true,
			wantErrMatch: "config has no tools",
		},
		{
			name: "single tool",
			input: `{
"tools": [
{"name": "get_time", "description": "Returns current time", "parameters": {}}
]
}`,
			wantErr:   false,
			wantCount: 1,
		},
		{
			name: "multiple tools",
			input: `{
"tools": [
{"name": "get_time", "description": "Returns current time", "parameters": {}},
{"name": "fetch_url", "description": "Fetches a URL", "parameters": {}}
]
}`,
			wantErr:   false,
			wantCount: 2,
		},
		{
			name: "tool with parameters",
			input: `{
"tools": [
{
"name": "run_shell",
"description": "Runs a shell command",
"parameters": {
"command": {"type": "string", "description": "The command", "required": true}
}
}
]
}`,
			wantErr:   false,
			wantCount: 1,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			var cfg Config
			var err error

			if testCase.nilReader {
				cfg, err = Read(nil)
			} else {
				cfg, err = Read(strings.NewReader(testCase.input))
			}

			if testCase.wantErr {
				assert.Error(t, err)
				if testCase.wantErrMatch != "" {
					assert.ErrorContains(t, err, testCase.wantErrMatch)
				}
				return
			}

			assert.NoError(t, err)
			assert.Len(t, cfg.Tools, testCase.wantCount)
		})
	}
}

func TestToolSchema_ToOllama(t *testing.T) {
	cases := []struct {
		name             string
		schema           ToolSchema
		wantType         string
		wantFuncName     string
		wantFuncDesc     string
		wantRequired     []string
		wantPropertyKeys []string
	}{
		{
			name: "no parameters",
			schema: ToolSchema{
				Name:        "get_time",
				Description: "Returns current time",
				Parameters:  map[string]ToolParameter{},
			},
			wantType:         "function",
			wantFuncName:     "get_time",
			wantFuncDesc:     "Returns current time",
			wantRequired:     []string{},
			wantPropertyKeys: []string{},
		},
		{
			name: "all required parameters",
			schema: ToolSchema{
				Name:        "run_shell",
				Description: "Runs a shell command",
				Parameters: map[string]ToolParameter{
					"command": {Type: "string", Description: "The command to run", Required: true},
				},
			},
			wantType:         "function",
			wantFuncName:     "run_shell",
			wantRequired:     []string{"command"},
			wantPropertyKeys: []string{"command"},
		},
		{
			name: "optional parameter excluded from required",
			schema: ToolSchema{
				Name:        "fetch_url",
				Description: "Fetches a URL",
				Parameters: map[string]ToolParameter{
					"url":     {Type: "string", Description: "The URL", Required: true},
					"timeout": {Type: "integer", Description: "Timeout in seconds", Required: false},
				},
			},
			wantType:         "function",
			wantFuncName:     "fetch_url",
			wantRequired:     []string{"url"},
			wantPropertyKeys: []string{"url", "timeout"},
		},
		{
			name: "required list is sorted alphabetically",
			schema: ToolSchema{
				Name:        "write_file",
				Description: "Writes content to a file",
				Parameters: map[string]ToolParameter{
					"path":    {Type: "string", Description: "File path", Required: true},
					"content": {Type: "string", Description: "File content", Required: true},
					"mode":    {Type: "string", Description: "Write mode", Required: true},
				},
			},
			wantType:     "function",
			wantFuncName: "write_file",
			wantRequired: []string{"content", "mode", "path"},
		},
		{
			name: "enum values preserved",
			schema: ToolSchema{
				Name:        "set_level",
				Description: "Sets a level",
				Parameters: map[string]ToolParameter{
					"level": {Type: "string", Description: "The level", Enum: []string{"low", "medium", "high"}, Required: true},
				},
			},
			wantType:         "function",
			wantFuncName:     "set_level",
			wantRequired:     []string{"level"},
			wantPropertyKeys: []string{"level"},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			result := testCase.schema.ToOllama()

			assert.Equal(t, testCase.wantType, result.Type)
			assert.Equal(t, testCase.wantFuncName, result.Function.Name)

			if testCase.wantFuncDesc != "" {
				assert.Equal(t, testCase.wantFuncDesc, result.Function.Description)
			}

			assert.Equal(t, "object", result.Function.Parameters.Type)

			if testCase.wantRequired != nil {
				assert.Equal(t, testCase.wantRequired, result.Function.Parameters.Required)
			}

			for _, key := range testCase.wantPropertyKeys {
				assert.Contains(t, result.Function.Parameters.Properties, key)
			}
		})
	}
}

func TestToolSchema_ToOllama_EnumPreserved(t *testing.T) {
	schema := ToolSchema{
		Name:        "set_level",
		Description: "Sets a level",
		Parameters: map[string]ToolParameter{
			"level": {Type: "string", Description: "The level", Enum: []string{"low", "medium", "high"}, Required: true},
		},
	}

	result := schema.ToOllama()

	prop, ok := result.Function.Parameters.Properties["level"]
	assert.True(t, ok)
	assert.Equal(t, []string{"low", "medium", "high"}, prop.Enum)
}

func TestToolSchema_ToOllama_PropertyFields(t *testing.T) {
	schema := ToolSchema{
		Name:        "run_shell",
		Description: "Runs a shell command",
		Parameters: map[string]ToolParameter{
			"command": {Type: "string", Description: "The command to run", Required: true},
		},
	}

	result := schema.ToOllama()

	prop, ok := result.Function.Parameters.Properties["command"]
	assert.True(t, ok)
	assert.Equal(t, "string", prop.Type)
	assert.Equal(t, "The command to run", prop.Description)
}
