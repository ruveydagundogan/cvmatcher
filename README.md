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

# CV Matcher

**AI-powered CV & Job Description matching platform using local LLM (Ollama/Gemma) and Go backend**

<p>
⚡ Go + Chi • 🧠 Ollama (Gemma-2B) • ▲ Next.js 16 • 🐘 PostgreSQL • 🐳 Docker
</p>

[![Deployed on Render](https://img.shields.io/badge/Render-46E3B7?logo=render&logoColor=fff)](https://cvmatcher-api.onrender.com)
[![Deployed on Vercel](https://img.shields.io/badge/Vercel-000?logo=vercel)](https://cvmatcherapp.vercel.app)
[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go)](https://go.dev)
[![Next.js](https://img.shields.io/badge/Next.js-16-000?logo=next.js)](https://nextjs.org)

</div>

---

## Overview

CV Matcher is a full-stack application that helps recruiters and hiring managers match **CVs (Resumes)** with **Job Descriptions** using AI. The platform:

1. **Parses** CVs to extract skills, experience, and education
2. **Analyzes** job descriptions to identify required/preferred skills
3. **Matches** CVs against JDs and produces a detailed compatibility score (0-100%)
4. Provides an **admin panel** for managing AI model settings, system prompts, and query logs
5. Includes a **knowledge base** (DeepKwiki) for storing and searching reference information

All AI processing runs locally via **Ollama** (Gemma-2B) — no external API calls, no data leaves your machine.

---

## Architecture

```
┌─────────────────────────┐      HTTP/JWT      ┌──────────────────────┐
│   Frontend (Next.js)    │ ──────────────────> │   Backend (Go+Chi)  │
│                         │ <────────────────── │                      │
│  /dashboard             │      JSON API       │  - CV/JD CRUD        │
│  /dashboard/resumes     │                     │  - Match Engine      │
│  /dashboard/jds         │                     │  - MCP Protocol      │
│  /dashboard/matches     │                     │  - Knowledge Base    │
│  /dashboard/knowledge   │                     │  - Admin Panel       │
│  /admin/*               │                     │  - Auth (JWT)        │
│                         │                     │  - Rate Limiting     │
└─────────────────────────┘                     └───────┬──────────────┘
                                                        │
                                                        │ OpenAI-compatible API
                                                        ▼
                                               ┌──────────────────┐
                                               │  Ollama (Gemma)  │
                                               │  Local LLM       │
                                               └──────────────────┘
```

### Tech Stack

| Component | Technology | Purpose |
|-----------|-----------|---------|
| **Frontend** | Next.js 16, React 19, Tailwind CSS v4 | User interface |
| **Backend** | Go 1.26, Chi router, pgx | REST API, business logic |
| **Database** | PostgreSQL 16 (primary), In-memory (fallback) | Data persistence |
| **Cache** | Redis 7 | Rate limiting, caching |
| **AI Engine** | Ollama (Gemma-2B) | CV parsing, JD analysis, matching |
| **Protocol** | MCP (Model Context Protocol) | Standardized LLM communication |
| **Monitoring** | Prometheus + Grafana | Metrics, dashboards |
| **Deployment** | Render (backend), Vercel (frontend) | Cloud hosting |

---

## Features

### Core Features
- **CV Management** — Upload, view, parse CVs with AI-powered skill extraction
- **JD Management** — Create, edit, analyze job descriptions
- **AI Matching** — Match CVs vs JDs with detailed score breakdown (overall, skill, experience, education)
- **Dashboard** — Overview stats (total CVs, JDs, matches, average scores)
- **Fallback Mode** — Keyword-based parsing/scoring when LLM is unavailable

### Advanced Features
- **MCP Protocol** — Standardized LLM query interface (`POST /api/v1/mcp/query`)
- **DeepKwiki** — Knowledge base with full-text search (`/dashboard/knowledge`)
- **Admin Panel** (`/admin`) — Adapter management, system prompts, LLM settings, query logs
- **Rich Results** — AI responses rendered as sections, tables, code blocks
- **Dark Mode** — Theme toggle with persistent preference

---

## Quick Start

### Prerequisites
- Go 1.26+
- Node.js 20+
- Docker & Docker Compose
- [Ollama](https://ollama.ai) (for local AI)

### 1. Start AI Model

```bash
# Install and run Ollama with Qwen
ollama pull qwen2.5:1.5b-instruct
ollama run qwen2.5:1.5b-instruct
```

### 2. Start Backend

```bash
cd backend

# Option A: Full stack with Docker
docker compose -f ../docker-compose.yml up postgres redis

# Option B: With Ollama in Docker
docker compose -f ../docker-compose.yml up postgres redis ollama

# Run the backend
MLC_LLM_BASE_URL=http://localhost:11434 go run ./cmd/server
```

### 3. Start Frontend

```bash
cd frontend
npm install
NEXT_PUBLIC_API_URL=http://localhost:8080 npm run dev
```

### 4. Open Browser
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Grafana: http://localhost:3001 (admin/admin)
- Prometheus: http://localhost:9090

---

## API Endpoints

### Public
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/auth/register` | Register new user |
| POST | `/api/v1/auth/login` | Login |
| POST | `/api/v1/llm/chat` | Chat with AI |
| POST | `/api/v1/mcp/query` | MCP protocol query |

### Protected (JWT required)
| Method | Path | Description |
|--------|------|-------------|
| GET/PUT | `/api/v1/me` | User profile |
| CRUD | `/api/v1/cvs` | CV management |
| POST | `/api/v1/cvs/{id}/parse` | Parse CV with AI |
| CRUD | `/api/v1/jds` | JD management |
| POST | `/api/v1/jds/{id}/analyze` | Analyze JD with AI |
| POST | `/api/v1/matches` | Run CV-JD match |
| GET | `/api/v1/dashboard/stats` | Dashboard stats |
| CRUD | `/api/v1/knowledge` | Knowledge base |
| GET | `/api/v1/knowledge/search` | Full-text search |
| | `/api/v1/admin/*` | Admin panel API |

---

## Deployment

### Backend (Render)
The backend auto-deploys from the `main` branch via `render.yaml`. It connects to the local Ollama instance through a **WebSocket tunnel** — the `start-tunnel.sh` script on your machine establishes a tunnel so Render can reach your local AI model.

```bash
# Start the tunnel (connects Render to local Ollama)
cd backend && bash start-tunnel.sh
```

### Frontend (Vercel)
Auto-deploys from the `main` branch. Configure `NEXT_PUBLIC_API_URL` environment variable in the Vercel dashboard to point to your Render backend URL.

---

## Project Structure

```
├── backend/                  # Go backend
│   ├── cmd/server/           # Entry point
│   ├── internal/
│   │   ├── application/      # Use cases
│   │   │   ├── admin/        # Admin panel logic
│   │   │   ├── cv/           # CV use cases
│   │   │   ├── jobdescription/
│   │   │   ├── knowledge/    # DeepKwiki use cases
│   │   │   ├── matching/     # Match scoring
│   │   │   └── iam/          # Auth & users
│   │   ├── domain/           # Domain models & interfaces
│   │   ├── infrastructure/   # Implementations
│   │   │   ├── http/         # Handlers, router
│   │   │   ├── llm/          # LLM client
│   │   │   ├── mcp/          # MCP engine
│   │   │   ├── postgres/     # Postgres repos
│   │   │   ├── tunnel/       # WebSocket tunnel
│   │   │   └── memory/       # In-memory fallbacks
│   │   └── shared/           # Config, middleware, errors
│   ├── migrations/           # SQL migrations
│   └── deployments/          # Docker Compose
├── frontend/                 # Next.js frontend
│   ├── src/app/
│   │   ├── dashboard/        # Main app pages
│   │   ├── admin/            # Admin panel pages
│   │   └── register/         # Registration
│   └── src/components/       # Shared components
├── docker-compose.yml        # Full stack
├── render.yaml               # Render config
└── grafana/                  # Monitoring dashboards
```