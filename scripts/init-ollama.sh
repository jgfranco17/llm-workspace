#!/usr/bin/env bash
# Initialization script for Ollama custom configuration
# This script can be used to pre-pull models, create custom Modelfiles,
# or apply any other initialization logic after Ollama starts.

set -eux

echo "Ollama initialization script started..."

# Wait for Ollama service to be ready
while ! ollama list >/dev/null 2>&1; do
    echo "Waiting for Ollama service to start..."
    sleep 2
done

echo "Ollama service is ready."

# Create custom personalized model if configured
if [ "${CREATE_CUSTOM_MODEL:-true}" = "true" ] && [ -f /etc/ollama/Modelfile.custom ]; then
    echo "Creating custom personalized model: ${CUSTOM_MODEL_NAME:-llama-dev}"

    # First ensure base model is available
    BASE_MODEL="${BASE_MODEL:-llama3}"
    if ! ollama list | grep -q "$BASE_MODEL"; then
        echo "Pulling base model: $BASE_MODEL"
        ollama pull "$BASE_MODEL" || {
            echo "Failed to pull base model $BASE_MODEL"
            exit 1
        }
    fi

    # Create custom model from Modelfile
    ollama create "${CUSTOM_MODEL_NAME:-llama-dev}" -f /etc/ollama/Modelfile.custom || {
        echo "Failed to create custom model"
        exit 1
    }

    echo "Custom model ${CUSTOM_MODEL_NAME:-llama-dev} created successfully"
fi

# Example: Pull default models if specified
if [ -n "${OLLAMA_DEFAULT_MODELS:-}" ]; then
    IFS=',' read -ra MODELS <<< "${OLLAMA_DEFAULT_MODELS}"
    for model in "${MODELS[@]}"; do
        model=$(echo "$model" | xargs) # trim whitespace
        if [ -n "$model" ]; then
            echo "Pulling model: $model"
            ollama pull "$model" || echo "Failed to pull $model"
        fi
    done
fi

# Example: Create custom Modelfile with extended context
if [ "${CREATE_EXTENDED_MODELS:-false}" = "true" ]; then
    echo "Creating custom models with extended context..."

    # You can create custom Modelfiles here
    # Example for llama3 with custom context:
    cat > /tmp/Modelfile.llama3-extended <<EOF
FROM llama3
PARAMETER num_ctx 8192
PARAMETER temperature 0.7
EOF
    # ollama create llama3-extended -f /tmp/Modelfile.llama3-extended || echo "Failed to create custom model"
fi

echo "Ollama initialization complete."
