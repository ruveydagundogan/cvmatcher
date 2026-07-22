-- Migration 002: Allow NULL user_id for anonymous scoring

-- Drop foreign key constraint from scoring_requests
ALTER TABLE scoring_requests DROP CONSTRAINT IF EXISTS scoring_requests_user_id_fkey;

-- Drop foreign key constraint from audit_logs  
ALTER TABLE audit_logs DROP CONSTRAINT IF EXISTS audit_logs_user_id_fkey;

-- Make user_id nullable in scoring_requests
ALTER TABLE scoring_requests ALTER COLUMN user_id DROP NOT NULL;

-- Make user_id nullable in audit_logs
ALTER TABLE audit_logs ALTER COLUMN user_id DROP NOT NULL;
