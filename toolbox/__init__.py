"""Python tools for Ollama agent function calling."""

# Import tools to trigger registration
import tools  # noqa: F401
from registry import tool_registry

__version__ = "0.1.0"

__all__ = ["tool_registry"]
