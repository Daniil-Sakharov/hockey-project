-- +goose Up
-- Migration: Sync indexes that were created manually in production

-- Indexes for failed_parsing_jobs (created manually in prod)
CREATE INDEX IF NOT EXISTS idx_failed_parsing_jobs_job_type ON failed_parsing_jobs(job_type);
CREATE INDEX IF NOT EXISTS idx_failed_parsing_jobs_next_retry ON failed_parsing_jobs(next_retry_at) WHERE next_retry_at IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_failed_parsing_jobs_next_retry;
DROP INDEX IF EXISTS idx_failed_parsing_jobs_job_type;
