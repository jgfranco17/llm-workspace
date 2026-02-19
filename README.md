# AI Playground

A customized Ollama environment with personalized context, configurable settings, and GPU support.

## Features

- **Personalized AI Assistant**: Custom model configured with your preferences
- **Go Tools Integration**: Built-in tool support with a compiled tool runner
- **Extended Context Window**: Default 4096 tokens (configurable, optimized for RTX 2060)
- **GPU Support**: NVIDIA GPU acceleration enabled
- **Auto-initialization**: Automatic model creation and configuration on startup
- **Easy Management**: justfile commands for common operations

## Personalized Model

The setup creates a custom model called `llama-dev` that knows your preferences and work context.
This context is automatically applied, so the AI provides more relevant assistance without
repeated explanations on startup.

## Configuration

Environment variables can be set in a `.env` file (see `.env.example`):

```bash
# Custom Model Settings
CREATE_CUSTOM_MODEL=true
CUSTOM_MODEL_NAME=llama-dev
BASE_MODEL=llama3

# Ollama Runtime Configuration
OLLAMA_NUM_CTX=4096
OLLAMA_NUM_PARALLEL=1
OLLAMA_MAX_LOADED_MODELS=1
OLLAMA_KEEP_ALIVE=5m
```

To customize the AI's personality and context, edit [Modelfile.custom](Modelfile.custom).

## Quick Start

```bash
# Start the environment (auto-creates custom model)
just up

# Run your personalized model
just run

# Or run with a different model
just run llama3

# List loaded models
just list-models

# Show custom model configuration
just show-custom

# Pull additional models
just pull-model codellama

# Stop the environment
just down
```

## Go Tools

The agent has access to Go-based tools for enhanced capabilities.
Available tools include:

- **get_current_time**: Get current date and time
- **fetch_url**: Retrieve content from URLs
- **run_shell_command**: Execute shell commands safely
- **read_file**: Read file contents
- **write_file**: Write content to files

See [toolbox/README.md](toolbox/README.md) for detailed documentation and how to add new tools.

## Troubleshooting

### Slow Text Generation (1 word/minute)

This usually means Ollama is running on CPU instead of GPU.

**Check GPU access:**

```bash
just gpu-check
```

**Fix GPU access:**

```bash
# Install NVIDIA Container Toolkit
just setup-gpu

# Then restart
just down && just up
```

**Verify it's working:**

- GPU generation should be ~20-50+ tokens/second
- CPU generation is ~1-5 tokens/second

### Out of Memory

If you see OOM errors with the RTX 2060 (6GB VRAM):

1. Reduce context window in `.envrc`: `OLLAMA_NUM_CTX=2048`
2. Use a smaller model: `just pull-model llama3:8b-instruct-q4_0`
3. Ensure only 1 model loaded: `OLLAMA_MAX_LOADED_MODELS=1`
