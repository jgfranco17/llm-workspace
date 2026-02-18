# syntax=docker/dockerfile:1

# This Dockerfile customizes the official Ollama image with
# context configuration and initialization capabilities.
ARG BASE_VERSION=latest
FROM ollama/ollama:${BASE_VERSION} AS setup

# Set default environment variables for Ollama configuration
# Optimized for RTX 2060 (6GB VRAM)
# These can be overridden in compose.yaml or at runtime
ENV OLLAMA_NUM_CTX=4096
ENV OLLAMA_NUM_PARALLEL=1
ENV OLLAMA_MAX_LOADED_MODELS=1
ENV OLLAMA_KEEP_ALIVE=5m
ENV OLLAMA_NUM_GPU=1

# Copy initialization script and custom Modelfile
COPY --chmod=755 init-ollama.sh /usr/local/bin/init-ollama.sh
COPY Modelfile.custom /etc/ollama/Modelfile.custom

# The default entrypoint from the base image will be preserved
# To run initialization after Ollama starts, use the healthcheck
# or docker-compose depends_on to trigger init-ollama.sh
FROM setup AS app
