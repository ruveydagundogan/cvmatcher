---
license: mit
language:
- en
tags:
- cv-parsing
- cv-jd-matching
- resume
- job-description
---

# CV Matcher Dataset

Fine-tuning dataset for CV parsing and CV-JD matching tasks.

## Files

- `cv_parse_dataset.json` — 50 CV examples with structured extraction (skills, experience, education, summary)
- `cv_jd_match_dataset.json` — 50 CV-JD pair examples with match scores (0.0-1.0) per category

## Usage

```python
from datasets import load_dataset

dataset = load_dataset("RuveydaGundogan/cvmatcher-dataset")
```

## Format

### CV Parse
```json
{
  "instruction": "Parse the following CV text...",
  "input": "Senior Full Stack Developer with...",
  "output": "{\"skills\": [...], \"experience\": [...], ...}"
}
```

### CV-JD Match
```json
{
  "instruction": "Analyze how well this CV matches...",
  "input": "CV: ...\n\nJD: ...",
  "output": "{\"overall_score\": 0.85, ...}"
}
```
