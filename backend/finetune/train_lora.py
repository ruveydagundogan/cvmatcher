#!/usr/bin/env python3
"""
PEFT/LoRA Fine-Tuning Pipeline for CV Matcher

Fine-tunes Qwen2.5-1.5B-Instruct using LoRA (Low-Rank Adaptation) on a CV parsing dataset.
Outputs a LoRA adapter that can be loaded into Ollama via Modelfile.

Usage:
    python train_lora.py --base-model Qwen/Qwen2.5-1.5B-Instruct --data data/cv_parse_dataset.json --output-dir ./adapters/cv-parser-v1
"""

import argparse
import json
import os
import sys

import torch
from datasets import Dataset
from transformers import (
    AutoModelForCausalLM,
    AutoTokenizer,
    BitsAndBytesConfig,
    TrainingArguments,
    Trainer,
    DataCollatorForLanguageModeling,
    pipeline,
)
from peft import (
    LoraConfig,
    get_peft_model,
    prepare_model_for_kbit_training,
    PeftModel,
    PeftConfig,
)


def parse_args():
    parser = argparse.ArgumentParser(description="LoRA fine-tuning for CV Matcher")
    parser.add_argument("--base-model", type=str, default="Qwen/Qwen2.5-1.5B-Instruct",
                        help="Base model name or path")
    parser.add_argument("--data", type=str, required=True,
                        help="Path to training data JSON")
    parser.add_argument("--output-dir", type=str, default="./adapters/cv-parser-v1",
                        help="Output directory for LoRA adapter")
    parser.add_argument("--lora-r", type=int, default=8,
                        help="LoRA rank (default: 8)")
    parser.add_argument("--lora-alpha", type=int, default=16,
                        help="LoRA alpha (default: 16)")
    parser.add_argument("--lora-dropout", type=float, default=0.05,
                        help="LoRA dropout (default: 0.05)")
    parser.add_argument("--epochs", type=int, default=3,
                        help="Number of training epochs (default: 3)")
    parser.add_argument("--batch-size", type=int, default=4,
                        help="Training batch size (default: 4)")
    parser.add_argument("--lr", type=float, default=2e-4,
                        help="Learning rate (default: 2e-4)")
    parser.add_argument("--max-length", type=int, default=512,
                        help="Max sequence length (default: 512)")
    parser.add_argument("--quantize", action="store_true", default=False,
                        help="Use 4-bit quantization to reduce memory")
    parser.add_argument("--mode", type=str, default="cv-parse", choices=["cv-parse", "cv-jd-match", "cv-coach"],
                        help="Training mode: cv-parse (extract structured CV data), cv-jd-match (CV-JD match scoring), or cv-coach (CV improvement chat)")
    parser.add_argument("--test-only", action="store_true", default=False,
                        help="Load existing adapter and test inference")
    return parser.parse_args()


def get_system_prompt(mode):
    if mode == "cv-jd-match":
        return "You are a CV-JD matching assistant. Analyze CV and job description pairs, provide match scores (0.0-1.0) for each category, highlight matched and missing skills, and give a detailed analysis."
    if mode == "cv-coach":
        return "You are a CV Coach. Help job seekers improve their CVs with concrete, actionable advice. Be encouraging but honest, and always give specific examples and numbers."
    return "You are a CV parsing assistant. Extract structured information from CV texts."


def load_dataset(data_path, mode, tokenizer):
    """Load and prepare the training dataset using the model's chat template."""
    with open(data_path, "r") as f:
        raw_data = json.load(f)

    formatted = []
    for item in raw_data:
        messages = [
            {"role": "user", "content": f"{item['instruction']}\n\n{item['input']}"},
            {"role": "assistant", "content": item["output"]},
        ]
        text = tokenizer.apply_chat_template(messages, tokenize=False)
        formatted.append({"text": text})

    return Dataset.from_list(formatted)


def create_lora_model(base_model_name, tokenizer, quantize=False):
    """Load base model and apply LoRA configuration."""
    print(f"[INFO] Loading base model: {base_model_name}")

    bnb_config = None
    if quantize:
        bnb_config = BitsAndBytesConfig(
            load_in_4bit=True,
            bnb_4bit_quant_type="nf4",
            bnb_4bit_compute_dtype=torch.bfloat16,
            bnb_4bit_use_double_quant=True,
        )

    model = AutoModelForCausalLM.from_pretrained(
        base_model_name,
        quantization_config=bnb_config,
        device_map="auto",
        torch_dtype=torch.bfloat16 if not quantize else None,
        trust_remote_code=True,
    )

    if quantize:
        model = prepare_model_for_kbit_training(model)

    lora_config = LoraConfig(
        r=args.lora_r,
        lora_alpha=args.lora_alpha,
        lora_dropout=args.lora_dropout,
        bias="none",
        task_type="CAUSAL_LM",
        target_modules=["q_proj", "k_proj", "v_proj", "o_proj",
                        "gate_proj", "up_proj", "down_proj"],
    )

    model = get_peft_model(model, lora_config)
    model.print_trainable_parameters()

    return model


def train(model, tokenizer, dataset, output_dir, args):
    """Run LoRA fine-tuning."""
    print(f"[INFO] Starting training for {args.epochs} epochs")

    training_args = TrainingArguments(
        output_dir=output_dir,
        num_train_epochs=args.epochs,
        per_device_train_batch_size=args.batch_size,
        gradient_accumulation_steps=4,
        gradient_checkpointing=True,
        optim="paged_adamw_8bit",
        logging_steps=10,
        save_strategy="epoch",
        learning_rate=args.lr,
        bf16=torch.cuda.is_bf16_supported(),
        fp16=not torch.cuda.is_bf16_supported(),
        max_grad_norm=0.3,
        warmup_ratio=0.03,
        lr_scheduler_type="cosine",
        report_to="none",
    )

    def tokenize_fn(examples):
        return tokenizer(examples["text"], truncation=True, max_length=args.max_length, padding=False)

    tokenized = dataset.map(tokenize_fn, batched=True, remove_columns=["text"])

    collator = DataCollatorForLanguageModeling(tokenizer=tokenizer, mlm=False)

    trainer = Trainer(
        model=model,
        args=training_args,
        train_dataset=tokenized,
        data_collator=collator,
    )

    trainer.train()
    trainer.save_model(output_dir)
    tokenizer.save_pretrained(output_dir)

    print(f"[INFO] Model saved to {output_dir}")
    return trainer


def create_ollama_modelfile(adapter_path, base_model, output_name="cv-parser"):
    """Create a Modelfile for Ollama to load the LoRA adapter."""
    ollama_model = base_model.replace("google/", "").replace("-it", "")
    if "/" in ollama_model:
        ollama_model = ollama_model.split("/")[-1].lower()
    adapter_dirname = os.path.basename(adapter_path)

    # Qwen2.5 requires GGUF adapter format (Safetensors ADAPTER not supported)
    model_id = base_model.replace("/", "--")
    gguf_path = os.path.join(adapter_path, f"{output_name}.gguf")
    modelfile_content = f"""FROM {ollama_model}
ADAPTER {gguf_path}
"""
    modelfile_path = os.path.join(adapter_path, "Modelfile")
    with open(modelfile_path, "w") as f:
        f.write(modelfile_content)

    print(f"[INFO] Modelfile created at {modelfile_path}")
    print(f"[INFO] Pull base model if needed: ollama pull {ollama_model}")
    print(f"[INFO] NOTE: Qwen2.5 does not support Safetensors ADAPTER in Ollama.")
    print(f"[INFO] Convert to GGUF first:")
    print(f"  python3 /tmp/llama.cpp/convert_lora_to_gguf.py \\")
    print(f"    --base-model ~/.cache/huggingface/hub/models--{model_id}/snapshots/*/ \\")
    print(f"    --lora-model {adapter_path} \\")
    print(f"    --output {gguf_path}")
    print(f"[INFO] Then load into Ollama:")
    print(f"  cd {os.path.dirname(adapter_path)}")
    print(f"  ollama create {output_name} -f {os.path.join(adapter_dirname, 'Modelfile')}")
    print(f"  ollama run {output_name}")

    return modelfile_path


def test_inference(model, tokenizer, prompt_text):
    """Test inference with the fine-tuned model."""
    pipe = pipeline(
        "text-generation",
        model=model,
        tokenizer=tokenizer,
        max_new_tokens=256,
        temperature=0.1,
        top_p=0.95,
        do_sample=True,
    )

    result = pipe(prompt_text)
    return result[0]["generated_text"]


def main():
    global args
    args = parse_args()

    os.makedirs(args.output_dir, exist_ok=True)

    tokenizer = AutoTokenizer.from_pretrained(args.base_model)
    tokenizer.pad_token = tokenizer.eos_token
    tokenizer.padding_side = "right"

    if args.test_only:
        print("[INFO] Test mode: loading existing adapter")
        peft_config = PeftConfig.from_pretrained(args.output_dir)
        base_model_name = peft_config.base_model_name_or_path

        model = AutoModelForCausalLM.from_pretrained(
            base_model_name,
            device_map="auto",
            torch_dtype=torch.bfloat16,
        )
        model = PeftModel.from_pretrained(model, args.output_dir)
        model.eval()

        test_messages = [
            {"role": "user", "content": "Parse the following CV text and extract structured information: skills, experience, education, and a brief summary.\n\nPython developer with 6 years of backend experience. Skilled in Django, FastAPI, Redis, and Celery. Led team of 5 engineers at TechCo. Master's in Software Engineering."}
        ]
        test_prompt = tokenizer.apply_chat_template(test_messages, tokenize=False, add_generation_prompt=True)
        result = test_inference(model, tokenizer, test_prompt)
        print(f"[RESULT]\n{result}")
        return

    dataset = load_dataset(args.data, args.mode, tokenizer)
    print(f"[INFO] Loaded {len(dataset)} training examples")

    model = create_lora_model(args.base_model, tokenizer, args.quantize)

    trainer = train(model, tokenizer, dataset, args.output_dir, args)

    output_name = "cv-parser"
    if args.mode == "cv-jd-match":
        output_name = "cv-jd-matcher"
    elif args.mode == "cv-coach":
        output_name = "cv-coach"
    create_ollama_modelfile(
        adapter_path=os.path.abspath(args.output_dir),
        base_model=args.base_model,
        output_name=output_name
    )

    print("[INFO] Pipeline complete!")
    print(f"[INFO] Adapter saved to: {os.path.abspath(args.output_dir)}")
    print(f"[INFO] To test: python train_lora.py --test-only --output-dir {args.output_dir}")


if __name__ == "__main__":
    main()
