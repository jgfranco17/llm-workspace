#!/usr/bin/env python3
"""Example usage of the tool system.

This demonstrates how to use the tool registry from Python code.
"""

from toolbox import tool_registry


def main():
    """Run example tool invocations."""
    print("=== Tool System Examples ===\n")

    # Example 1: Get current time
    print("1. Getting current time:")
    time_result = tool_registry.execute("get_current_time")
    print(f"   Result: {time_result}\n")

    # Example 2: Calculate expression
    print("2. Calculate 2 + 2:")
    calc_result = tool_registry.execute("calculate", expression="2 + 2")
    print(f"   Result: {calc_result}\n")

    # Example 3: Calculate with math functions
    print("3. Calculate sqrt(16) + pi:")
    math_result = tool_registry.execute("calculate", expression="sqrt(16) + pi")
    print(f"   Result: {math_result}\n")

    # Example 4: List all tools
    print("4. Available tools:")
    tools = tool_registry.list_tools()
    for tool in tools:
        param_count = len(tool.parameters)
        print(f"   - {tool.name}: {tool.description} ({param_count} params)")
    print()

    # Example 5: Export Ollama format
    print("5. Ollama format schema (first tool):")
    schemas = tool_registry.to_ollama_format()
    if schemas:
        import json

        print(json.dumps(schemas[0], indent=2))


if __name__ == "__main__":
    main()
