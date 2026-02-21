package config

import (
	"sort"
)

// ToolParameter describes a single parameter accepted by a tool.
type ToolParameter struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
	Required    bool     `json:"required"`
}

// ToolSchema describes a tool's name, purpose, and parameter contract.
type ToolSchema struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Parameters  map[string]ToolParameter `json:"parameters"`
}

// Toolset maps tool names to their ToolSchema definitions.
type Toolset map[string]ToolSchema

// OllamaFormat wraps an OllamaFunction in the envelope expected by the
// Ollama API.
type OllamaFormat struct {
	Type     string         `json:"type"`
	Function OllamaFunction `json:"function"`
}

// OllamaFunction holds the name, description, and parameter schema
// of a tool in the Ollama function-calling format.
type OllamaFunction struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Parameters  OllamaParameters `json:"parameters"`
}

// OllamaParameters describes the object schema for an OllamaFunction's
// arguments, including required field names.
type OllamaParameters struct {
	Type       string                    `json:"type"`
	Properties map[string]OllamaProperty `json:"properties"`
	Required   []string                  `json:"required,omitempty"`
}

// OllamaProperty describes a single property within OllamaParameters.
type OllamaProperty struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

// ToOllama converts the ToolSchema to the OllamaFormat used by the
// Ollama API. Required parameter names are sorted alphabetically.
func (schema ToolSchema) ToOllama() OllamaFormat {
	properties := make(map[string]OllamaProperty, len(schema.Parameters))
	required := make([]string, 0, len(schema.Parameters))

	for name, param := range schema.Parameters {
		properties[name] = OllamaProperty{
			Type:        param.Type,
			Description: param.Description,
			Enum:        param.Enum,
		}

		if param.Required {
			required = append(required, name)
		}
	}

	sort.Strings(required)

	return OllamaFormat{
		Type: "function",
		Function: OllamaFunction{
			Name:        schema.Name,
			Description: schema.Description,
			Parameters: OllamaParameters{
				Type:       "object",
				Properties: properties,
				Required:   required,
			},
		},
	}
}
