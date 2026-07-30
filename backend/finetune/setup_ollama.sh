#!/bin/bash
# CV Matcher - Ollama LoRA Adapter Setup
# Run this after downloading the LoRA adapter zip from Colab
# Usage: ./setup_ollama.sh /path/to/cvmatcher-lora.zip

set -e

ZIP_PATH="${1:-$HOME/Downloads/cvmatcher-lora.zip}"
ADAPTER_DIR="backend/finetune/adapters"

if [ ! -f "$ZIP_PATH" ]; then
    echo "Zip file not found at: $ZIP_PATH"
    echo "Usage: $0 /path/to/cvmatcher-lora.zip"
    exit 1
fi

echo "Extracting adapters..."
mkdir -p "$ADAPTER_DIR"
unzip -o "$ZIP_PATH" -d "$ADAPTER_DIR"

echo ""
echo "Creating Ollama models..."

if [ -d "$ADAPTER_DIR/cv-parser-v1" ]; then
    echo "  -> cv-parser"
    ollama create cv-parser -f "$ADAPTER_DIR/cv-parser-v1/Modelfile" 2>/dev/null || {
        echo "Creating Modelfile for cv-parser..."
        echo "FROM gemma:2b" > "$ADAPTER_DIR/cv-parser-v1/Modelfile"
        echo "ADAPTER $ADAPTER_DIR/cv-parser-v1" >> "$ADAPTER_DIR/cv-parser-v1/Modelfile"
        ollama create cv-parser -f "$ADAPTER_DIR/cv-parser-v1/Modelfile"
    }
fi

if [ -d "$ADAPTER_DIR/cv-jd-matcher-v1" ]; then
    echo "  -> cv-jd-matcher"
    ollama create cv-jd-matcher -f "$ADAPTER_DIR/cv-jd-matcher-v1/Modelfile" 2>/dev/null || {
        echo "Creating Modelfile for cv-jd-matcher..."
        echo "FROM gemma:2b" > "$ADAPTER_DIR/cv-jd-matcher-v1/Modelfile"
        echo "ADAPTER $ADAPTER_DIR/cv-jd-matcher-v1" >> "$ADAPTER_DIR/cv-jd-matcher-v1/Modelfile"
        ollama create cv-jd-matcher -f "$ADAPTER_DIR/cv-jd-matcher-v1/Modelfile"
    }
fi

echo ""
echo "Done! Models created:"
ollama list | grep -E "cv-parser|cv-jd-matcher"

echo ""
echo "To test:"
echo "  ollama run cv-parser"
echo ""
echo "To use in backend:"
echo "  export OLLAMA_MODEL=cv-parser"
echo "  go run cmd/server/main.go"
