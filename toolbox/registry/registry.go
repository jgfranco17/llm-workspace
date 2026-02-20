package registry

import (
	"context"
	"fmt"

	"github.com/jgfranco17/llm-workspace/internal/config"
	"github.com/jgfranco17/llm-workspace/toolbox/tools"
)

type Manager struct {
	tools    config.Toolset
	order    []string
	handlers tools.HandlerCollection
}

func New(handlers tools.HandlerCollection) *Manager {
	if handlers == nil {
		handlers = tools.DefaultHandlers()
	}
	return &Manager{
		tools:    make(config.Toolset),
		order:    []string{},
		handlers: handlers,
	}
}

func (m *Manager) Add(schemas ...config.ToolSchema) error {
	for _, schema := range schemas {
		if !m.handlers.Has(schema.Name) {
			return fmt.Errorf("no handler registered for tool: %s", schema.Name)
		}
		if _, exists := m.tools[schema.Name]; exists {
			return fmt.Errorf("duplicate tool schema: %s", schema.Name)
		}
		m.tools[schema.Name] = schema
		m.order = append(m.order, schema.Name)
	}
	return nil
}

func (m *Manager) ListTools() []config.ToolSchema {
	tools := make([]config.ToolSchema, 0, len(m.order))
	for _, name := range m.order {
		tools = append(tools, m.tools[name])
	}
	return tools
}

func (m *Manager) GetHandler(name string) (config.ToolSchema, tools.Handler, bool) {
	schema, ok := m.tools[name]
	if !ok {
		return config.ToolSchema{}, nil, false
	}
	handler := m.handlers[name]
	return schema, handler, true
}

func (m *Manager) AsOllamaSchema() []config.OllamaFormat {
	schemas := make([]config.OllamaFormat, 0, len(m.order))
	for _, name := range m.order {
		schemas = append(schemas, m.tools[name].ToOllama())
	}
	return schemas
}

func (m *Manager) Execute(ctx context.Context, name string, params tools.Parameters) (string, error) {
	schema, handler, ok := m.GetHandler(name)
	if !ok {
		return "", fmt.Errorf("no handler registered for tool: %s", name)
	}

	_ = schema
	return handler(ctx, params)
}
