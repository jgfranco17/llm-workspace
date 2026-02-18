"""Tests for Python tools."""

import pytest
from registry import tool_registry

import toolbox  # noqa: F401 - Import to trigger tool registration


def test_get_current_time():
    """Test get_current_time tool."""
    result = tool_registry.execute("get_current_time")
    assert result
    assert "T" in result  # ISO 8601 format includes T


def test_calculate():
    """Test calculate tool."""
    result = tool_registry.execute("calculate", expression="2 + 2")
    assert result == "4"

    result = tool_registry.execute("calculate", expression="sqrt(16)")
    assert result == "4.0"


def test_calculate_error():
    """Test calculate tool with invalid expression."""
    result = tool_registry.execute("calculate", expression="invalid")
    assert "Error" in result


def test_tool_registry():
    """Test tool registry functionality."""
    tools = tool_registry.list_tools()
    assert len(tools) > 0

    tool_names = [t.name for t in tools]
    assert "calculate" in tool_names
    assert "get_current_time" in tool_names


def test_ollama_format():
    """Test Ollama format export."""
    schemas = tool_registry.to_ollama_format()
    assert len(schemas) > 0

    for schema in schemas:
        assert schema["type"] == "function"
        assert "function" in schema
        assert "name" in schema["function"]
        assert "description" in schema["function"]
        assert "parameters" in schema["function"]


def test_tool_not_found():
    """Test executing non-existent tool."""
    with pytest.raises(ValueError, match="not found"):
        tool_registry.execute("nonexistent_tool")
