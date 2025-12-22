-- +goose Up
-- Migration: Add composite key support and Event Sourcing
-- Implements architectural decisions #3 (Composite Key) and #5 (Event Sourcing)

-- ============================================
-- PLAYERS: Add missing external_id and domain for composite key
-- ============================================
ALTER TABLE players ADD COLUMN IF NOT EXISTS external_id TEXT;
ALTER TABLE players ADD COLUMN IF NOT EXISTS domain TEXT;

-- Update existing data: extract external_id from existing id field
UPDATE players SET external_id = id WHERE external_id IS NULL;

-- Create composite unique constraint (source, external_id, domain)
CREATE UNIQUE INDEX IF NOT EXISTS idx_players_composite_key 
    ON players(source, external_id, domain);

-- ============================================
-- EVENT SOURCING: player_events table
-- ============================================
CREATE TABLE IF NOT EXISTS player_events (
    id SERIAL PRIMARY KEY,
    player_id TEXT NOT NULL,
    event_type VARCHAR(50) NOT NULL,    -- 'created', 'updated', 'position_changed', 'stats_updated'
    event_data JSONB NOT NULL,          -- Full event payload
    source VARCHAR(20) NOT NULL,        -- 'fhspb', 'junior'
    domain VARCHAR(50),                 -- 'cfo.fhr.ru', 'fhspb.ru', etc.
    created_at TIMESTAMPTZ DEFAULT NOW(),
    
    -- Metadata for debugging and ML
    parsing_session_id UUID,
    user_agent TEXT,
    ip_address INET
);

-- Indexes for Event Sourcing queries
CREATE INDEX IF NOT EXISTS idx_player_events_player ON player_events(player_id);
CREATE INDEX IF NOT EXISTS idx_player_events_type ON player_events(event_type);
CREATE INDEX IF NOT EXISTS idx_player_events_source ON player_events(source);
CREATE INDEX IF NOT EXISTS idx_player_events_created ON player_events(created_at);
CREATE INDEX IF NOT EXISTS idx_player_events_session ON player_events(parsing_session_id);

-- GIN index for JSONB queries
CREATE INDEX IF NOT EXISTS idx_player_events_data ON player_events USING gin(event_data);

-- Comments
COMMENT ON TABLE player_events IS 'Event Sourcing: All player changes for audit trail and ML';
COMMENT ON COLUMN player_events.event_data IS 'JSONB payload with before/after values';
COMMENT ON COLUMN player_events.parsing_session_id IS 'Links events to parsing session for debugging';

-- +goose Down
DROP TABLE IF EXISTS player_events;
DROP INDEX IF EXISTS idx_players_composite_key;
ALTER TABLE players DROP COLUMN IF EXISTS domain;
ALTER TABLE players DROP COLUMN IF EXISTS external_id;
