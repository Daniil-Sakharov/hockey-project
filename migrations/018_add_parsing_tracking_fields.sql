-- +goose Up
-- Migration: Add parsing tracking fields to tournaments

ALTER TABLE tournaments ADD COLUMN IF NOT EXISTS last_players_parsed_at TIMESTAMP;
ALTER TABLE tournaments ADD COLUMN IF NOT EXISTS last_stats_parsed_at TIMESTAMP;

CREATE INDEX IF NOT EXISTS idx_tournaments_players_parsed ON tournaments(last_players_parsed_at);
CREATE INDEX IF NOT EXISTS idx_tournaments_stats_parsed ON tournaments(last_stats_parsed_at);

-- +goose Down
DROP INDEX IF EXISTS idx_tournaments_stats_parsed;
DROP INDEX IF EXISTS idx_tournaments_players_parsed;
ALTER TABLE tournaments DROP COLUMN IF EXISTS last_stats_parsed_at;
ALTER TABLE tournaments DROP COLUMN IF EXISTS last_players_parsed_at;
