-- Admin adapters (PEFT/LoRA)
CREATE TABLE IF NOT EXISTS admin_adapters (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    file_path TEXT NOT NULL,
    active BOOLEAN DEFAULT false,
    model_name VARCHAR(255) NOT NULL DEFAULT 'gemma:2b',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- System prompts
CREATE TABLE IF NOT EXISTS system_prompts (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    active BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- LLM settings
CREATE TABLE IF NOT EXISTS llm_settings (
    id SERIAL PRIMARY KEY,
    max_tokens INTEGER NOT NULL DEFAULT 2048,
    temperature REAL NOT NULL DEFAULT 0.7,
    top_p REAL NOT NULL DEFAULT 0.9,
    context_length INTEGER NOT NULL DEFAULT 4096,
    model_name VARCHAR(255) NOT NULL DEFAULT 'gemma:2b',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Query logs
CREATE TABLE IF NOT EXISTS query_logs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    query TEXT NOT NULL,
    response TEXT,
    model VARCHAR(255) NOT NULL DEFAULT 'gemma:2b',
    adapter VARCHAR(255),
    duration_ms BIGINT NOT NULL DEFAULT 0,
    token_count INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'success',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_query_logs_user_id ON query_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_query_logs_created_at ON query_logs(created_at DESC);

-- Knowledge entries (DeepKwiki)
CREATE TABLE IF NOT EXISTS knowledge_entries (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    content TEXT NOT NULL,
    tags TEXT[] DEFAULT '{}',
    category VARCHAR(255),
    source VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_knowledge_entries_user_id ON knowledge_entries(user_id);
CREATE INDEX IF NOT EXISTS idx_knowledge_entries_category ON knowledge_entries(category);
CREATE INDEX IF NOT EXISTS idx_knowledge_entries_tags ON knowledge_entries USING GIN(tags);
CREATE INDEX IF NOT EXISTS idx_knowledge_entries_search ON knowledge_entries USING GIN(to_tsvector('english', title || ' ' || content));

-- Seed default settings
INSERT INTO llm_settings (max_tokens, temperature, top_p, context_length, model_name)
SELECT 2048, 0.7, 0.9, 4096, 'gemma:2b'
WHERE NOT EXISTS (SELECT 1 FROM llm_settings LIMIT 1);
