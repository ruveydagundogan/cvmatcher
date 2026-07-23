-- CV & Job Description Matching Schema

-- CVs table
CREATE TABLE IF NOT EXISTS cvs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    title VARCHAR(255) NOT NULL DEFAULT '',
    content TEXT NOT NULL DEFAULT '',
    parsed_skills TEXT[] DEFAULT '{}',
    parsed_experience JSONB DEFAULT '[]',
    parsed_education JSONB DEFAULT '[]',
    parsed_summary TEXT DEFAULT '',
    status VARCHAR(20) DEFAULT 'pending',
    file_name VARCHAR(255) DEFAULT '',
    file_size BIGINT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Job descriptions table
CREATE TABLE IF NOT EXISTS job_descriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL DEFAULT '',
    required_skills TEXT[] DEFAULT '{}',
    preferred_skills TEXT[] DEFAULT '{}',
    experience_level VARCHAR(50) DEFAULT '',
    employment_type VARCHAR(50) DEFAULT '',
    location VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Match results table
CREATE TABLE IF NOT EXISTS match_results (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    cv_id UUID REFERENCES cvs(id) ON DELETE CASCADE NOT NULL,
    jd_id UUID REFERENCES job_descriptions(id) ON DELETE CASCADE NOT NULL,
    overall_score NUMERIC(5,2) DEFAULT 0,
    skill_match_score NUMERIC(5,2) DEFAULT 0,
    experience_score NUMERIC(5,2) DEFAULT 0,
    education_score NUMERIC(5,2) DEFAULT 0,
    llm_analysis TEXT DEFAULT '',
    details JSONB DEFAULT '{}',
    matched_skills TEXT[] DEFAULT '{}',
    missing_skills TEXT[] DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_cvs_user_id ON cvs(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_jds_user_id ON job_descriptions(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_match_results_user_id ON match_results(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_match_results_cv_id ON match_results(cv_id);
CREATE INDEX IF NOT EXISTS idx_match_results_jd_id ON match_results(jd_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_match_results_unique ON match_results(cv_id, jd_id);
