# Playground scripts

# List all available commands
_default:
    @just --list --unsorted

# Build the tools binary for local use
build-tools:
    #!/usr/bin/env bash
    echo "Building tools binary..."
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
        go build -o tool-runner .

# Run tests for the toolbox
test:
    #!/usr/bin/env bash
    echo "Running tests..."
    go clean -testcache
    go test -cover ./...

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
    @docker exec llm-workspace ollama pull {{ model }}

# Initialize Ollama with custom configuration
init:
    @echo "Running Ollama initialization..."
    @docker exec llm-workspace /usr/local/bin/init-ollama.sh

# List currently loaded models
list-models:
    @docker exec llm-workspace ollama list

# Check GPU status and availability
gpu-check:
    @echo "Host GPU status:"
    @nvidia-smi 2>/dev/null || echo "No NVIDIA GPU found on host"
    @echo ""
    @echo "Container GPU access:"
    @docker exec llm-workspace printenv | grep -i nvidia || echo "NVIDIA env vars not set"
    @docker exec llm-workspace ls -la /dev | grep -i nvidia 2>/dev/null || echo "No NVIDIA devices in container"

# Run the custom personalized model
run model="llama-dev":  pull-model init
    @docker exec -it llm-workspace ollama run {{ model }}

# Show the custom model configuration
show-custom:
    @docker exec llm-workspace ollama show llama-dev
