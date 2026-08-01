-- 008: HR role, qwen defaults, and chat tables

-- Add hr role (IK specialist)
INSERT INTO roles (name, permissions) VALUES
    ('hr', '{score:write,history:read,history:write,profile:read,profile:write,candidate:read,candidate:write}')
ON CONFLICT (name) DO NOTHING;

-- Fix admin table defaults (gemma:2b -> qwen2.5:1.5b-instruct)
UPDATE admin_adapters SET model_name = 'qwen2.5:1.5b-instruct' WHERE model_name = 'gemma:2b';
UPDATE query_logs SET model = 'qwen2.5:1.5b-instruct' WHERE model = 'gemma:2b';
UPDATE llm_settings SET model_name = 'qwen2.5:1.5b-instruct' WHERE model_name = 'gemma:2b';

-- Chat conversations (CV coach)
CREATE TABLE IF NOT EXISTS conversations (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL DEFAULT 'New Chat',
    cv_id UUID REFERENCES cvs(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY,
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    token_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_conversations_user_id ON conversations(user_id, updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages(conversation_id, created_at ASC);
