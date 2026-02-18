# Python Tools for Ollama Agent

This directory contains Python-based tools that can be used by the Ollama agent for function
calling and task execution.

## Architecture

- **registry.py**: Core tool registry system with Pydantic schemas
- **tools.py**: Example tool implementations (calculator, file ops, shell commands, etc.)
- **cli.py**: Command-line interface for managing and executing tools
- **pyproject.toml**: UV configuration for dependency management

## Available Tools

| Tool                 | Description                               | Parameters                      |
| -------------------- | ----------------------------------------- | ------------------------------- |
| `get_current_time`   | Get current date/time in ISO 8601 format  | None                            |
| `calculate`          | Evaluate mathematical expressions         | `expression`                    |
| `fetch_url`          | Fetch content from a URL                  | `url`                           |
| `run_shell_command`  | Execute shell commands safely             | `command`, `workdir` (optional) |
| `read_file`          | Read file contents (max 100KB)            | `filepath`                      |
| `write_file`         | Write content to a file                   | `filepath`, `content`           |

## Usage

### From Host

```bash
# List all available tools
just tools-list

# Get info about a specific tool
just tools-info calculate

# Execute a tool
just tools-exec calculate expression="2 + 2"

# Export tool schemas
just tools-schema
```

### From Container

```bash
# Using the tool-runner symlink (recommended)
docker exec ollama tool-runner list
docker exec ollama tool-runner execute calculate expression="sqrt(16)"

# Or using the venv directly
docker exec ollama /toolbox/.venv/bin/python3 /toolbox/cli.py list
```

## Adding New Tools

1. Define your tool function in `tools.py`
2. Create a `ToolSchema` for it
3. Register it with the global `registry`
4. Rebuild the container: `just down && just up`

Example:

```python
def my_new_tool(param1: str, param2: int = 10) -> str:
    """Tool description."""
    return f"Processed {param1} with {param2}"


registry.register(
    ToolSchema(
        name="my_new_tool",
        description="Description of what this tool does",
        parameters={
            "param1": ToolParameter(
                type="string",
                description="Description of param1",
                required=True,
            ),
            "param2": ToolParameter(
                type="integer",
                description="Description of param2",
                required=False,
            ),
        },
    ),
    my_new_tool,
)
```

## Security Notes

- File operations are limited to 100KB files
- Shell commands have a 30-second timeout
- URL fetches are limited to 1000 characters
- All tools run in the container environment
- Extend safety measures for production use
