# Playground scripts

# List all available commands
_default:
    @just --list --unsorted

# Setup NVIDIA Container Toolkit for GPU access
setup-gpu:
    @echo "Setting up NVIDIA Container Toolkit..."
    @bash setup-gpu.sh

# Launch the Dockerized environment
up:
    @echo "Starting Docker environment..."
    @docker compose -f compose.yaml up \
        --abort-on-container-failure \
        --build

# Stop and remove all containers
down:
    @echo "Stopping Docker environment..."
    @docker compose -f compose.yaml down --remove-orphans

# View container logs
logs service="":
    #!/usr/bin/env bash
    if [ -z "{{ service }}" ]; then
        docker compose -f compose.yaml logs -f
    else
        docker compose -f compose.yaml logs -f {{ service }}
    fi

# Pull Ollama model manually
pull-model model="llama3":
    @echo "Pulling Ollama model: {{ model }}"
    @docker exec ollama ollama pull {{ model }}

# Initialize Ollama with custom configuration
init:
    @echo "Running Ollama initialization..."
    @docker exec ollama /usr/local/bin/init-ollama.sh

# List currently loaded models
list-models:
    @docker exec ollama ollama list

# Check GPU status and availability
gpu-check:
    @echo "Host GPU status:"
    @nvidia-smi 2>/dev/null || echo "No NVIDIA GPU found on host"
    @echo ""
    @echo "Container GPU access:"
    @docker exec ollama printenv | grep -i nvidia || echo "NVIDIA env vars not set"
    @docker exec ollama ls -la /dev | grep -i nvidia 2>/dev/null || echo "No NVIDIA devices in container"

# Run the custom personalized model
run model="llama-dev":  pull-model init
    @docker exec -it ollama ollama run {{ model }}

# Show the custom model configuration
show-custom:
    @docker exec ollama ollama show llama-dev

# Run unit tests
pytest *args:
    @uv run pytest {{ args }}

# Lint the workspace
lint:
    @uv run black .
    @uv run isort .

# List all available Python tools
tools-list:
    @echo "Available Python tools:"
    @docker exec ollama /toolbox/.venv/bin/python3 /toolbox/cli.py list

# Show detailed info about a specific tool
tools-info tool:
    @docker exec ollama /toolbox/.venv/bin/python3 /toolbox/cli.py info {{ tool }}

# Execute a Python tool
tools-exec tool *args:
    @docker exec ollama /toolbox/.venv/bin/python3 /toolbox/cli.py execute {{ tool }} {{ args }}

# Export tool schemas in Ollama format
tools-schema:
    @docker exec ollama /toolbox/.venv/bin/python3 /toolbox/cli.py schema

# Run Python tools tests
tools-test:
    @docker exec ollama sh -c "cd /toolbox && .venv/bin/python3 -m pytest"
