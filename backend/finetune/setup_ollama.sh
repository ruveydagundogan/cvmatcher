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

create_from_modelfile() {
    local name="$1"
    local dir="$ADAPTER_DIR/$2"
    if [ ! -d "$dir" ]; then
        echo "  SKIP: $dir not found"
        return
    fi
    echo "  -> $name"

    if [ -f "$dir/Modelfile" ]; then
        ollama create "$name" -f "$dir/Modelfile"
    else
        BASE_MODEL="${3:-gemma:2b}"
        echo "FROM $BASE_MODEL" > "$dir/Modelfile"
        echo "ADAPTER $dir" >> "$dir/Modelfile"
        echo "  (Modelfile created with base: $BASE_MODEL)"
        ollama create "$name" -f "$dir/Modelfile"
    fi
}

create_from_modelfile "cv-parser" "cv-parser-v1" "gemma:2b"
create_from_modelfile "cv-jd-matcher" "cv-jd-matcher-v1" "gemma:2b"

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
