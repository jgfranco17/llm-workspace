#!/usr/bin/env python3
"""Example tools for Ollama agent."""

import json
import subprocess
from datetime import datetime
from pathlib import Path

import httpx
from registry import ToolParameter, ToolSchema, tool_registry


def get_current_time() -> str:
    """Get the current time in ISO 8601 format."""
    return datetime.now().isoformat()


def calculate(expression: str) -> str:
    """Safely evaluate a mathematical expression.

    Args:
        expression: Mathematical expression to evaluate (e.g., "2 + 2", "sqrt(16)")

    Returns:
        Result of the calculation as a string
    """
    try:
        # Use Python's eval with restricted builtins for safety
        import math

        allowed_names = {k: v for k, v in math.__dict__.items() if not k.startswith("__")}
        allowed_names.update({"abs": abs, "round": round})

        result = eval(expression, {"__builtins__": {}}, allowed_names)
        return str(result)
    except Exception as e:
        return f"Error: {str(e)}"


def fetch_url(url: str) -> str:
    """Fetch content from a URL.

    Args:
        url: URL to fetch

    Returns:
        Response text or error message
    """
    try:
        with httpx.Client(timeout=10.0) as client:
            response = client.get(url)
            response.raise_for_status()
            return response.text[:1000]  # Limit response size
    except Exception as e:
        return f"Error fetching URL: {str(e)}"


def run_shell_command(command: str, workdir: str = "/workspace") -> str:
    """Run a shell command in a safe environment.

    Args:
        command: Shell command to execute
        workdir: Working directory for the command

    Returns:
        Command output or error message
    """
    try:
        result = subprocess.run(
            command,
            shell=True,
            cwd=workdir,
            capture_output=True,
            text=True,
            timeout=30,
        )
        output = result.stdout if result.returncode == 0 else result.stderr
        return output[:1000]  # Limit output size
    except Exception as e:
        return f"Error executing command: {str(e)}"


def read_file(filepath: str) -> str:
    """Read contents of a file.

    Args:
        filepath: Path to the file to read

    Returns:
        File contents or error message
    """
    try:
        path = Path(filepath)
        if not path.exists():
            return f"Error: File '{filepath}' does not exist"

        # Limit file size for safety
        if path.stat().st_size > 100_000:
            return "Error: File too large (max 100KB)"

        return path.read_text()
    except Exception as e:
        return f"Error reading file: {str(e)}"


def write_file(filepath: str, content: str) -> str:
    """Write content to a file.

    Args:
        filepath: Path to the file to write
        content: Content to write to the file

    Returns:
        Success message or error message
    """
    try:
        path = Path(filepath)
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(content)
        return f"Successfully wrote {len(content)} bytes to {filepath}"
    except Exception as e:
        return f"Error writing file: {str(e)}"


# Register all tools
tool_registry.register(
    ToolSchema(
        name="get_current_time",
        description="Get the current date and time in ISO 8601 format",
        parameters={},
    ),
    get_current_time,
)

tool_registry.register(
    ToolSchema(
        name="calculate",
        description="Evaluate a mathematical expression safely",
        parameters={
            "expression": ToolParameter(
                type="string",
                description="Mathematical expression to evaluate (e.g., '2 + 2', 'sqrt(16)')",
            )
        },
    ),
    calculate,
)

tool_registry.register(
    ToolSchema(
        name="fetch_url",
        description="Fetch content from a URL (first 1000 characters)",
        parameters={
            "url": ToolParameter(
                type="string",
                description="URL to fetch",
            )
        },
    ),
    fetch_url,
)

tool_registry.register(
    ToolSchema(
        name="run_shell_command",
        description="Execute a shell command in a safe environment",
        parameters={
            "command": ToolParameter(
                type="string",
                description="Shell command to execute",
            ),
            "workdir": ToolParameter(
                type="string",
                description="Working directory for the command",
                required=False,
            ),
        },
    ),
    run_shell_command,
)

tool_registry.register(
    ToolSchema(
        name="read_file",
        description="Read contents of a file (max 100KB)",
        parameters={
            "filepath": ToolParameter(
                type="string",
                description="Path to the file to read",
            )
        },
    ),
    read_file,
)

tool_registry.register(
    ToolSchema(
        name="write_file",
        description="Write content to a file",
        parameters={
            "filepath": ToolParameter(
                type="string",
                description="Path to the file to write",
            ),
            "content": ToolParameter(
                type="string",
                description="Content to write to the file",
            ),
        },
    ),
    write_file,
)


if __name__ == "__main__":
    # Export tools schema for reference
    tools = tool_registry.to_ollama_format()
    print(json.dumps(tools, indent=2))
