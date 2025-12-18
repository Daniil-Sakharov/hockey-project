-- Migration: Create failed_parsing_jobs table for retry system
-- This table stores failed parsing tasks that need to be retried

CREATE TABLE IF NOT EXISTS failed_parsing_jobs (
    id SERIAL PRIMARY KEY,
    job_type VARCHAR(50) NOT NULL,     -- 'tournament', 'team', 'player'
    source VARCHAR(20) NOT NULL,       -- 'fhspb', 'junior'
    external_id VARCHAR(255) NOT NULL,
    url TEXT,
    error_message TEXT,
    retry_count INTEGER DEFAULT 0,
    max_retries INTEGER DEFAULT 3,
    next_retry_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Index for efficient retry processing
CREATE INDEX IF NOT EXISTS idx_failed_jobs_retry ON failed_parsing_jobs(source, next_retry_at, retry_count);

-- Index for cleanup
CREATE INDEX IF NOT EXISTS idx_failed_jobs_created ON failed_parsing_jobs(created_at);
