# syntax=docker/dockerfile:1

# Versioning
ARG BASE_VERSION=latest
ARG GO_VERSION=1.24.3

FROM golang:${GO_VERSION} AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY internal/ ./internal/
COPY toolbox/ ./toolbox/
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/tool-runner .

# This Dockerfile customizes the official Ollama image with
# context configuration and initialization capabilities.
FROM ollama/ollama:${BASE_VERSION} AS model

LABEL org.opencontainers.image.source="https://github.com/jgfranco17/llm-workspace"
LABEL org.opencontainers.image.description="Sandbox environment for testing customized LLMs."

# Set default environment variables for Ollama configuration
# Optimized for RTX 2060 (6GB VRAM)
# These can be overridden in compose.yaml or at runtime
ENV OLLAMA_NUM_CTX=4096
ENV OLLAMA_NUM_PARALLEL=1
ENV OLLAMA_MAX_LOADED_MODELS=1
ENV OLLAMA_KEEP_ALIVE=10m
ENV OLLAMA_NUM_GPU=1

# Copy initialization script and custom Modelfile
COPY --chmod=755 scripts/init-ollama.sh /usr/local/bin/init-ollama.sh
COPY Modelfile.custom /etc/ollama/Modelfile.custom

# Install compiled tool runner and config
ENV TOOL_RUNNER_PATH="/usr/local/bin/tool-runner"
COPY --from=build /out/tool-runner "${TOOL_RUNNER_PATH}"
ENV TOOL_CONFIG_PATH="/config/tools.json"
COPY tools.json "${TOOL_CONFIG_PATH}"

# Reset workdir and preserve entrypoint
WORKDIR /workspace
ENV TOOL_LOG_LEVEL="INFO"

# The default entrypoint from the base image will be preserved
# To run initialization after Ollama starts, use the healthcheck
# or docker-compose depends_on to trigger init-ollama.sh
FROM model AS app
