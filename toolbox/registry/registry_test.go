package registry

import (
	"context"
	"testing"

	"github.com/jgfranco17/llm-workspace/internal/config"
	"github.com/jgfranco17/llm-workspace/toolbox/tools"
	"github.com/stretchr/testify/assert"
)

func TestManager_Add(t *testing.T) {
	cases := []struct {
		name         string
		handlers     tools.HandlerCollection
		schemas      []config.ToolSchema
		wantErr      bool
		wantErrMatch string
	}{
		{
			name:     "adds schema",
			handlers: tools.HandlerCollection{"tool": noopHandler("ok")},
			schemas: []config.ToolSchema{
				{Name: "tool", Description: "desc", Parameters: map[string]config.ToolParameter{}},
			},
			wantErr: false,
		},
		{
			name:         "missing handler",
			handlers:     tools.HandlerCollection{},
			schemas:      []config.ToolSchema{{Name: "tool", Description: "desc", Parameters: map[string]config.ToolParameter{}}},
			wantErr:      true,
			wantErrMatch: "no handler registered",
		},
		{
			name:     "duplicate schema",
			handlers: tools.HandlerCollection{"tool": noopHandler("ok")},
			schemas: []config.ToolSchema{
				{Name: "tool", Description: "desc", Parameters: map[string]config.ToolParameter{}},
				{Name: "tool", Description: "dup", Parameters: map[string]config.ToolParameter{}},
			},
			wantErr:      true,
			wantErrMatch: "duplicate tool schema",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			manager := New(testCase.handlers)
			err := manager.Add(testCase.schemas...)
			if testCase.wantErr {
				assert.Error(t, err)
				if testCase.wantErrMatch != "" {
					assert.ErrorContains(t, err, testCase.wantErrMatch)
				}
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestManager_ListTools_PreservesOrder(t *testing.T) {
	manager := New(tools.HandlerCollection{
		"first":  noopHandler("1"),
		"second": noopHandler("2"),
	})

	err := manager.Add(
		config.ToolSchema{Name: "first", Description: "desc", Parameters: map[string]config.ToolParameter{}},
		config.ToolSchema{Name: "second", Description: "desc", Parameters: map[string]config.ToolParameter{}},
	)
	assert.NoError(t, err)

	schemas := manager.ListTools()
	if assert.Len(t, schemas, 2) {
		assert.Equal(t, "first", schemas[0].Name)
		assert.Equal(t, "second", schemas[1].Name)
	}
}

func TestManagerGetHandler(t *testing.T) {
	handler := noopHandler("ok")
	manager := New(tools.HandlerCollection{"tool": handler})
	assert.NoError(t, manager.Add(config.ToolSchema{Name: "tool", Description: "desc", Parameters: map[string]config.ToolParameter{}}))

	schema, gotHandler, ok := manager.GetHandler("tool")
	assert.True(t, ok)
	assert.Equal(t, "tool", schema.Name)
	assert.NotNil(t, gotHandler)

	_, _, ok = manager.GetHandler("missing")
	assert.False(t, ok)
}

func TestManagerOllamaSchemas(t *testing.T) {
	manager := New(tools.HandlerCollection{"tool": noopHandler("ok")})
	err := manager.Add(config.ToolSchema{
		Name:        "tool",
		Description: "desc",
		Parameters: map[string]config.ToolParameter{
			"value": {Type: "string", Description: "v", Required: true},
		},
	})
	assert.NoError(t, err)

	schemas := manager.AsOllamaSchema()
	if assert.Len(t, schemas, 1) {
		assert.Equal(t, "function", schemas[0].Type)
		assert.Equal(t, "tool", schemas[0].Function.Name)
	}
}

func TestManagerExecute(t *testing.T) {
	cases := []struct {
		name         string
		manager      *Manager
		toolName     string
		params       tools.Parameters
		wantResult   string
		wantErr      bool
		wantErrMatch string
	}{
		{
			name:       "executes handler",
			manager:    managerWithTool("tool", noopHandler("ok")),
			toolName:   "tool",
			params:     tools.Parameters{"value": "1"},
			wantErr:    false,
			wantResult: "ok",
		},
		{
			name:         "missing tool",
			manager:      managerWithTool("tool", noopHandler("ok")),
			toolName:     "missing",
			params:       tools.Parameters{},
			wantErr:      true,
			wantErrMatch: "no handler registered",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := testCase.manager.Execute(context.Background(), testCase.toolName, testCase.params)
			if testCase.wantErr {
				assert.Error(t, err)
				if testCase.wantErrMatch != "" {
					assert.ErrorContains(t, err, testCase.wantErrMatch)
				}
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, testCase.wantResult, result)
		})
	}
}

func noopHandler(result string) tools.Handler {
	return func(_ context.Context, _ tools.Parameters) (string, error) {
		return result, nil
	}
}

func managerWithTool(name string, handler tools.Handler) *Manager {
	manager := New(tools.HandlerCollection{name: handler})
	_ = manager.Add(config.ToolSchema{Name: name, Description: "desc", Parameters: map[string]config.ToolParameter{}})
	return manager
}
