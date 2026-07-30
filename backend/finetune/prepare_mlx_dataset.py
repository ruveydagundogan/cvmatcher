import json
import os
import random

DATA_DIR = "backend/finetune/data"
OUTPUT_DIR = "backend/finetune/mlx_data"
os.makedirs(OUTPUT_DIR, exist_ok=True)

def format_gemma(instruction, input_text, output_text):
    if input_text:
        user_content = f"{instruction}\n\n{input_text}"
    else:
        user_content = instruction
    return (
        f"<start_of_turn>user\n{user_content}<end_of_turn>\n"
        f"<start_of_turn>model\n{output_text}<end_of_turn>"
    )

all_samples = []

for fname in ["cv_parse_dataset.json", "cv_jd_match_dataset.json"]:
    samples = json.load(open(os.path.join(DATA_DIR, fname)))
    for s in samples:
        text = format_gemma(s["instruction"], s.get("input", ""), s["output"])
        all_samples.append({"text": text})

random.seed(42)
random.shuffle(all_samples)

n = len(all_samples)
train = all_samples[:int(n * 0.8)]
valid = all_samples[int(n * 0.8):int(n * 0.9)]
test = all_samples[int(n * 0.9):]

for name, data in [("train", train), ("valid", valid), ("test", test)]:
    path = os.path.join(OUTPUT_DIR, f"{name}.jsonl")
    with open(path, "w") as f:
        for item in data:
            f.write(json.dumps(item) + "\n")
    print(f"{name}: {len(data)} samples -> {path}")

print("Done!")
