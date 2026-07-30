#!/usr/bin/env python3
import subprocess
import sys
import os

DATA_DIR = os.path.join(os.path.dirname(__file__), "mlx_data")
MODEL = "mlx-community/quantized-gemma-2b-it"
ADAPTER_DIR = os.path.join(os.path.dirname(__file__), "adapters")

os.makedirs(ADAPTER_DIR, exist_ok=True)

def run_lora(output_name, train=True):
    adapter_path = os.path.join(ADAPTER_DIR, output_name)
    cmd = [
        sys.executable, "-m", "mlx_lm.lora",
        "--model", MODEL,
        "--data", DATA_DIR,
        "--adapter-path", adapter_path,
        "--num-layers", "16",
        "--batch-size", "2",
        "--iters", "200",
        "--steps-per-report", "10",
        "--steps-per-eval", "20",
        "--val-batches", "5",
        "--learning-rate", "1e-4",
        "--max-seq-length", "512",
        "--fine-tune-type", "lora",
    ]
    if train:
        cmd += ["--train"]
    else:
        cmd += ["--test", "--test-batches", "5"]

    print(f"Running: {' '.join(cmd)}")
    subprocess.run(cmd, check=True)
    return adapter_path

if __name__ == "__main__":
    mode = sys.argv[1] if len(sys.argv) > 1 else "train"
    
    if mode == "train":
        adapter = run_lora("cv-parser-v1", train=True)
        print(f"\nAdapter saved to: {adapter}")
        print(f"\nTo create Ollama model:")
        print(f"  echo 'FROM gemma:2b' > {adapter}/Modelfile")
        print(f"  echo 'ADAPTER {adapter}' >> {adapter}/Modelfile")
        print(f"  ollama create cv-parser -f {adapter}/Modelfile")
