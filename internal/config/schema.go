package config

import (
	"sort"
)

type ToolParameter struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
	Required    bool     `json:"required"`
}

type ToolSchema struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Parameters  map[string]ToolParameter `json:"parameters"`
}

type OllamaFormat struct {
	Type     string         `json:"type"`
	Function OllamaFunction `json:"function"`
}

type OllamaFunction struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Parameters  OllamaParameters `json:"parameters"`
}

type OllamaParameters struct {
	Type       string                    `json:"type"`
	Properties map[string]OllamaProperty `json:"properties"`
	Required   []string                  `json:"required,omitempty"`
}

type OllamaProperty struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

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
