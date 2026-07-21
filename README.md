<div align="center">

<a href="https://academy.masterfabric.co">
  <img src="https://academy.masterfabric.co/academy-badge.png" width="120" alt="MasterFabric Academy">
</a>

<p>
  <sub>
    academy.masterfabric.co is a
    <a href="https://masterfabric.co">MasterFabric</a>
    subsidiary.
  </sub>
</p>

# LLM Decision Score

**AI-powered decision scoring using WebLLM (Gemma) and Go Backend**

<p>
🧠 WebLLM • 💎 Gemma-2B • ⚡ Go (Gin) • ▲ Next.js • 🎯 Decision Scoring
</p>

</div>

---

## Overview

LLM Decision Score is a full-stack AI application that combines **browser-based LLM inference** with **backend evaluation**. Using **WebLLM (Gemma-2B)**, responses are generated locally in the browser, while a **Go (Gin)** backend analyzes each prompt-response pair to calculate a decision score and maintain request history.

---

## Features

- 🚀 Local AI inference with WebLLM (Gemma-2B)
- 📊 AI-powered decision score calculation
- ⚡ Real-time inference metrics
- 📝 Prompt & response history
- 🎨 Responsive dashboard interface
- 🔗 REST API powered by Gin

---

## Tech Stack

| Layer | Technologies |
| :--- | :--- |
| **Frontend** | Next.js, TypeScript, Tailwind CSS |
| **Backend** | Go, Gin |
| **AI Model** | WebLLM (Gemma-2B) |

---

## Architecture

| Component | Responsibility |
| :--- | :--- |
| **Frontend** | Runs Gemma locally with WebLLM and displays prompts, responses, metrics, and history |
| **Backend** | Calculates decision scores and manages request history |
| **WebLLM** | Performs browser-based AI inference |

---

## Solution Workflow

```text
User Prompt
     │
     ▼
WebLLM (Gemma-2B)
     │
     ▼
Generated Response
     │
     ▼
Go Backend (Gin)
     │
     ▼
Decision Score
     │
     ▼
Dashboard & History
```

1. The user submits a prompt through the dashboard.
2. WebLLM generates a response locally in the browser.
3. The generated response is sent to the Go backend.
4. The backend evaluates the response and calculates a decision score.
5. The dashboard displays the response, metrics, score, and stores the request in history.

---

## Decision Score Criteria

Each response is evaluated using multiple quality indicators:

- Response length
- Prompt relevance
- Response structure
- Content richness
- Overall completeness

The final score is calculated on a **0–100** scale.

---

## Dashboard Metrics

| Metric | Description |
| :--- | :--- |
| **Model** | Active language model |
| **Inference Time** | Time required to generate the response |
| **Status** | Current model loading status |
| **Backend** | Backend connection status |
| **Word Count** | Total generated words |
| **Character Count** | Total generated characters |
| **Decision Score** | Backend evaluation score (0–100) |

---

## Project Structure

```text
.
├── frontend/
│   ├── app/
│   ├── components/
│   ├── hooks/
│   └── public/
│
├── backend/
│   ├── handler/
│   ├── model/
│   ├── router/
│   └── main.go
│
└── README.md
```

---

## API Endpoints

| Method | Endpoint | Description |
| :---: | :--- | :--- |
| POST | `/score` | Evaluate an AI response and return its decision score |
| GET | `/llm/history` | Retrieve prompt-response history |
| DELETE | `/llm/history` | Clear stored history |

---

## Getting Started

### 1. Start the Backend

```bash
cd backend
go run .
```

### 2. Start the Frontend

```bash
cd frontend
npm install
npm run dev
```

After starting both services, open the application in your browser.