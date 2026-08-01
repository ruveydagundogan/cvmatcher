#!/bin/bash
# CV Matcher - Ollama LoRA Adapter Setup (Qwen2.5 + GGUF)
# Usage: ./setup_ollama.sh /path/to/cvmatcher-lora.zip

set -e

ZIP_PATH="${1:-$HOME/Downloads/cvmatcher-lora.zip}"
WORK_DIR="backend/finetune/adapters-clean"
LLAMA_CPP="/tmp/llama.cpp"
CONVERT_SCRIPT="$LLAMA_CPP/convert_lora_to_gguf.py"

if [ ! -f "$ZIP_PATH" ]; then
    echo "Zip file not found at: $ZIP_PATH"
    echo "Usage: $0 /path/to/cvmatcher-lora.zip"
    exit 1
fi

if [ ! -f "$CONVERT_SCRIPT" ]; then
    echo "convert_lora_to_gguf.py not found at $CONVERT_SCRIPT"
    echo "Run: git clone https://github.com/ggerganov/llama.cpp.git /tmp/llama.cpp"
    exit 1
fi

echo "Extracting adapters to $WORK_DIR..."
mkdir -p "$WORK_DIR"
unzip -o "$ZIP_PATH" -d "$WORK_DIR"

echo ""
echo "Fixing macOS extended attributes (com.apple.provenance)..."
xattr -cr "$WORK_DIR" 2>/dev/null || true

BASE_MODEL_PATH="$HOME/.cache/huggingface/hub/Qwen/Qwen2.5-1.5B-Instruct"
# Also search in common locations
if [ ! -d "$BASE_MODEL_PATH" ]; then
    BASE_MODEL_PATH="$HOME/.cache/huggingface/hub/models--Qwen--Qwen2.5-1.5B-Instruct/snapshots/"*
fi

echo ""
echo "Creating Ollama models..."
echo "------------------------"

process_adapter() {
    local name="$1"
    local subdir="$2"

    local adapter_dir="$WORK_DIR/$subdir"
    if [ ! -d "$adapter_dir" ]; then
        echo "  SKIP: $adapter_dir not found"
        return
    fi

    echo "Processing: $name ($adapter_dir)"

    # Fix xattr on all files
    xattr -cr "$adapter_dir" 2>/dev/null || true

    local gguf_out="$adapter_dir/$name.gguf"
    if [ ! -f "$gguf_out" ]; then
        echo "  -> Converting LoRA to GGUF..."
        if [ -f "$adapter_dir/adapter_model.safetensors" ]; then
            python3 "$CONVERT_SCRIPT" \
                --base-model "$BASE_MODEL_PATH" \
                --lora-model "$adapter_dir" \
                --output "$gguf_out"
        elif [ -f "$adapter_dir/model.safetensors" ]; then
            python3 "$CONVERT_SCRIPT" \
                --base-model "$BASE_MODEL_PATH" \
                --lora-model "$adapter_dir" \
                --output "$gguf_out"
        else
            echo "  ERROR: No safetensors found in $adapter_dir!"
            ls -la "$adapter_dir/"
            return 1
        fi
    else
        echo "  -> GGUF already exists, skipping conversion"
    fi

    echo "  -> Creating Ollama model: $name"
    local abs_gguf_path
    abs_gguf_path="$(cd "$(dirname "$gguf_out")" && pwd)/$(basename "$gguf_out")"
    cat > "$adapter_dir/Modelfile" <<EOF
FROM qwen2.5:1.5b-instruct
ADAPTER $abs_gguf_path
EOF

    ollama create "$name" -f "$adapter_dir/Modelfile"
    echo "  -> Done: $name"
}

process_adapter "cv-parser" "cv-parser-v1"
process_adapter "cv-jd-matcher" "cv-jd-matcher-v1"
process_adapter "cv-coach" "cv-coach-v1"

echo ""
echo "Done! Models created:"
ollama list | grep -E "cv-parser|cv-jd-matcher|cv-coach"

echo ""
echo "To test:"
echo "  ollama run cv-parser"
echo "  ollama run cv-jd-matcher"
echo "  ollama run cv-coach"
echo ""
echo "To start the backend with fine-tuned models:"
echo "  cd backend && go run ./cmd/server"
