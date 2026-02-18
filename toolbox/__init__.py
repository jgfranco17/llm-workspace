"""Python tools for Ollama agent function calling."""

# Import tools to trigger registration
from registry import tool_registry

import toolbox  # noqa: F401

__version__ = "0.1.0"

__all__ = ["tool_registry"]
