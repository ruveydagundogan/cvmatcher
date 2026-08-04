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
    echo "Run: git clone --depth 1 https://github.com/ggerganov/llama.cpp.git /tmp/llama.cpp"
    exit 1
fi

# Patch llama.cpp converter for compatibility with newer torch versions
# (LoraTorchTensor is missing dim()/numel()/item(), used by conversion/base.py heuristics)
if ! grep -q "def dim(self)" "$CONVERT_SCRIPT"; then
    echo "Patching $CONVERT_SCRIPT for torch compatibility..."
    python3 - <<'PATCH'
path = "/tmp/llama.cpp/convert_lora_to_gguf.py"
src = open(path).read()
anchor = "    def contiguous(self) -> LoraTorchTensor:"
addition = """    def dim(self) -> int:
        return len(self._lora_A.shape)

    def numel(self) -> int:
        return self._lora_A.numel() + self._lora_B.numel()

    def item(self):
        raise NotImplementedError("LoraTorchTensor has no scalar value")

"""
if "def dim(self)" not in src:
    src = src.replace(anchor, addition + anchor, 1)
    open(path, "w").write(src)
    print("  -> patched")
PATCH
fi

echo "Extracting adapters to $WORK_DIR..."
mkdir -p "$WORK_DIR"
unzip -o "$ZIP_PATH" -d "$WORK_DIR"

echo ""
echo "Fixing macOS extended attributes (com.apple.provenance)..."
xattr -cr "$WORK_DIR" 2>/dev/null || true

BASE_MODEL_PATH="$HOME/.cache/huggingface/hub/models--Qwen--Qwen2.5-1.5B-Instruct"
# Resolve the actual snapshot directory (HF hub layout: .../snapshots/<hash>/config.json)
SNAPSHOT_DIR="$(ls -d "$BASE_MODEL_PATH"/snapshots/*/ 2>/dev/null | head -n1)"
if [ -n "$SNAPSHOT_DIR" ]; then
    BASE_MODEL_PATH="${SNAPSHOT_DIR%/}"
fi
if [ ! -f "$BASE_MODEL_PATH/config.json" ]; then
    echo "ERROR: base model config.json not found at $BASE_MODEL_PATH"
    echo "Download it first: huggingface-cli download Qwen/Qwen2.5-1.5B-Instruct"
    exit 1
fi
echo "Using base model at: $BASE_MODEL_PATH"

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
        if [ -f "$adapter_dir/adapter_model.safetensors" ] || [ -f "$adapter_dir/model.safetensors" ]; then
            python3 "$CONVERT_SCRIPT" \
                --base "$BASE_MODEL_PATH" \
                --outfile "$gguf_out" \
                --outtype auto \
                "$adapter_dir"
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
PARAMETER temperature 0.6
PARAMETER top_p 0.9
PARAMETER repeat_penalty 1.3
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
