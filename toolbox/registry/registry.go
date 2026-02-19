package registry

import (
	"context"
	"fmt"

	"github.com/jgfranco17/llm-workspace/internal/config"
)

type Handler func(ctx context.Context, params map[string]string) (string, error)

type Manager struct {
	tools    map[string]config.ToolSchema
	order    []string
	handlers map[string]Handler
}

func New(handlers map[string]Handler) *Manager {
	if handlers == nil {
		handlers = DefaultHandlers()
	}
	return &Manager{
		tools:    make(map[string]config.ToolSchema),
		order:    []string{},
		handlers: handlers,
	}
}

func (registry *Manager) Add(schemas ...config.ToolSchema) error {
	for _, schema := range schemas {
		if _, ok := registry.handlers[schema.Name]; !ok {
			return fmt.Errorf("no handler registered for tool: %s", schema.Name)
		}
		if _, exists := registry.tools[schema.Name]; exists {
			return fmt.Errorf("duplicate tool schema: %s", schema.Name)
		}
		registry.tools[schema.Name] = schema
		registry.order = append(registry.order, schema.Name)
	}
	return nil
}

func (registry *Manager) ListTools() []config.ToolSchema {
	tools := make([]config.ToolSchema, 0, len(registry.order))
	for _, name := range registry.order {
		tools = append(tools, registry.tools[name])
	}
	return tools
}

func (registry *Manager) GetTool(name string) (config.ToolSchema, Handler, bool) {
	schema, ok := registry.tools[name]
	if !ok {
		return config.ToolSchema{}, nil, false
	}
	handler := registry.handlers[name]
	return schema, handler, true
}

func (registry *Manager) OllamaSchemas() []config.OllamaFormat {
	schemas := make([]config.OllamaFormat, 0, len(registry.order))
	for _, name := range registry.order {
		schemas = append(schemas, registry.tools[name].ToOllama())
	}
	return schemas
}

func (registry *Manager) Execute(ctx context.Context, name string, params map[string]string) (string, error) {
	schema, handler, ok := registry.GetTool(name)
	if !ok {
		return "", fmt.Errorf("tool '%s' not found in registry", name)
	}

	_ = schema
	return handler(ctx, params)
}
