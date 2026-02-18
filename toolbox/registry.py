"""Tool registry for Ollama function-calling.

This module maintains a registry of available tools and their schemas.
"""

from typing import Any, Callable

from pydantic import BaseModel


class ToolParameter(BaseModel):
    """Schema for a tool parameter."""

    type: str
    description: str
    enum: list[str] | None = None
    required: bool = True


class ToolSchema(BaseModel):
    """Schema for a tool definition."""

    name: str
    description: str
    parameters: dict[str, ToolParameter]


class _ToolRegistry:
    """Registry for managing available tools."""

    def __init__(self):
        self._tools: dict[str, tuple[ToolSchema, Callable]] = {}

    def register(self, schema: ToolSchema, handler: Callable) -> None:
        """Register a tool with its schema and handler function."""
        self._tools[schema.name] = (schema, handler)

    def get_tool(self, name: str) -> tuple[ToolSchema, Callable] | None:
        """Get a tool by name."""
        return self._tools.get(name)

    def list_tools(self) -> list[ToolSchema]:
        """List all registered tools."""
        return [schema for schema, _ in self._tools.values()]

    def to_ollama_format(self) -> list[dict[str, Any]]:
        """Convert tool registry to Ollama function calling format."""
        tools = []
        for schema, _ in self._tools.values():
            tool = {
                "type": "function",
                "function": {
                    "name": schema.name,
                    "description": schema.description,
                    "parameters": {
                        "type": "object",
                        "properties": {
                            name: {
                                "type": param.type,
                                "description": param.description,
                                **({"enum": param.enum} if param.enum else {}),
                            }
                            for name, param in schema.parameters.items()
                        },
                        "required": [
                            name for name, param in schema.parameters.items() if param.required
                        ],
                    },
                },
            }
            tools.append(tool)
        return tools

    def execute(self, name: str, **kwargs) -> Any:
        """Execute a tool by name with given parameters."""
        tool_data = self.get_tool(name)
        if not tool_data:
            raise ValueError(f"Tool '{name}' not found in registry")

        _, handler = tool_data
        return handler(**kwargs)


# Global registry instance
tool_registry = _ToolRegistry()
