ALTER TABLE conversations
    ADD COLUMN IF NOT EXISTS jd_id uuid REFERENCES job_descriptions(id),
    ADD COLUMN IF NOT EXISTS match_id uuid REFERENCES match_results(id);

CREATE INDEX IF NOT EXISTS idx_conversations_jd_id ON conversations(jd_id);
CREATE INDEX IF NOT EXISTS idx_conversations_match_id ON conversations(match_id);
