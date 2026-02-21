# Go Tools for Ollama Agent

This directory contains Go-based tools that can be used by the Ollama agent for function
calling and task execution.

## Architecture

- **toolbox-go/**: Go tool runner and schema loader
- **toolbox-go/config/tools.json**: Tool schema definitions loaded at runtime
- **cli.py**: Legacy Python CLI (no longer used by default)

## Available Tools

| Tool                 | Description                               | Parameters                      |
| -------------------- | ----------------------------------------- | ------------------------------- |
| `get_current_time`   | Get current date/time in ISO 8601 format  | None                            |
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
just tools-info fetch_url

# Execute a tool
just tools-exec fetch_url url="https://example.com"

# Export tool schemas
just tools-schema
```

### From Container

```bash
# Using the tool-runner binary (recommended)
docker exec ollama tool-runner list
docker exec ollama tool-runner execute fetch_url url="https://example.com"
```

## Adding New Tools

1. Define your tool handler in `toolbox-go/internal/toolbox/tools.go`
2. Add its schema to `toolbox-go/config/tools.json`
3. Register the handler in `toolbox-go/internal/toolbox/handlers.go`
4. Rebuild the container: `just down && just up`

Example:

```json
{
    "name": "my_new_tool",
    "description": "Description of what this tool does",
    "parameters": {
        "param1": {
            "type": "string",
            "description": "Description of param1",
            "required": true
        },
        "param2": {
            "type": "integer",
            "description": "Description of param2",
            "required": false
        }
    }
}
```

## Security Notes

- File operations are limited to 100KB files
- Shell commands have a 30-second timeout
- URL fetches are limited to 1000 characters
- All tools run in the container environment
- Extend safety measures for production use
