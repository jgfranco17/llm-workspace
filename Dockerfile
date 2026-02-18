# syntax=docker/dockerfile:1

# This Dockerfile customizes the official Ollama image with
# context configuration and initialization capabilities.
ARG BASE_VERSION=latest
FROM ollama/ollama:${BASE_VERSION} AS setup

# Install Python, UV, and system dependencies
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    curl \
    && rm -rf /var/lib/apt/lists/*
COPY --from=ghcr.io/astral-sh/uv:0.9.28 /uv /uvx /bin/

# Set default environment variables for Ollama configuration
# Optimized for RTX 2060 (6GB VRAM)
# These can be overridden in compose.yaml or at runtime
ENV OLLAMA_NUM_CTX=4096
ENV OLLAMA_NUM_PARALLEL=1
ENV OLLAMA_MAX_LOADED_MODELS=1
ENV OLLAMA_KEEP_ALIVE=5m
ENV OLLAMA_NUM_GPU=1

# Copy Python tools and dependencies
WORKDIR /toolbox
COPY toolbox/ .
COPY pyproject.toml uv.lock ./
RUN uv sync --frozen

# Copy initialization script and custom Modelfile
COPY --chmod=755 scripts/init-ollama.sh /usr/local/bin/init-ollama.sh
COPY Modelfile.custom /etc/ollama/Modelfile.custom

# Create symlink for CLI access
RUN ln -s /toolbox/.venv/bin/python3 /usr/local/bin/tool-python && \
    chmod +x /toolbox/cli.py && \
    ln -s /toolbox/cli.py /usr/local/bin/tool-runner

# Reset workdir and preserve entrypoint
WORKDIR /workspace

# The default entrypoint from the base image will be preserved
# To run initialization after Ollama starts, use the healthcheck
# or docker-compose depends_on to trigger init-ollama.sh
FROM setup AS app
