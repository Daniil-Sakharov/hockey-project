-- +goose Up
-- Migration: Create scheduler_locks table for distributed locking

CREATE TABLE IF NOT EXISTS scheduler_locks (
    job_name VARCHAR(50) PRIMARY KEY,
    locked_at TIMESTAMP NOT NULL,
    locked_until TIMESTAMP NOT NULL,
    instance_id VARCHAR(100)
);

-- +goose Down
DROP TABLE IF EXISTS scheduler_locks;
