-- +goose Up
-- Migration: Add source column to tournaments and teams for unified architecture

-- ============================================
-- TOURNAMENTS: Add source column
-- ============================================
ALTER TABLE tournaments ADD COLUMN IF NOT EXISTS source TEXT DEFAULT 'junior';

-- Create unique index for (external_id, source)
CREATE UNIQUE INDEX IF NOT EXISTS idx_tournaments_external_id_source 
    ON tournaments(external_id, source) WHERE external_id IS NOT NULL;

-- ============================================
-- TEAMS: Add source column
-- ============================================
ALTER TABLE teams ADD COLUMN IF NOT EXISTS source TEXT DEFAULT 'junior';

-- Create unique index for (external_id, source)
CREATE UNIQUE INDEX IF NOT EXISTS idx_teams_external_id_source 
    ON teams(external_id, source) WHERE external_id IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_teams_external_id_source;
ALTER TABLE teams DROP COLUMN IF EXISTS source;

DROP INDEX IF EXISTS idx_tournaments_external_id_source;
ALTER TABLE tournaments DROP COLUMN IF EXISTS source;
