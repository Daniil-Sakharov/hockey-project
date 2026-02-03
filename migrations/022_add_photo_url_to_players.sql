-- +goose Up
-- Migration: Add photo_url column to players table for storing player photos

ALTER TABLE players ADD COLUMN IF NOT EXISTS photo_url TEXT;
CREATE INDEX IF NOT EXISTS idx_players_photo_url ON players(photo_url) WHERE photo_url IS NOT NULL;
COMMENT ON COLUMN players.photo_url IS 'URL фотографии игрока';

-- +goose Down
DROP INDEX IF EXISTS idx_players_photo_url;
ALTER TABLE players DROP COLUMN IF EXISTS photo_url;
