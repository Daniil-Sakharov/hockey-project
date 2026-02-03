-- +goose Up
-- Migration: Restore citizenship column for MIHF parser
-- The MIHF parser (stats.mihf.ru) parses player citizenship from profile pages

ALTER TABLE players ADD COLUMN IF NOT EXISTS citizenship VARCHAR(100);
CREATE INDEX IF NOT EXISTS idx_players_citizenship ON players(citizenship);
COMMENT ON COLUMN players.citizenship IS 'Гражданство игрока (из профиля MIHF)';

-- +goose Down
DROP INDEX IF EXISTS idx_players_citizenship;
ALTER TABLE players DROP COLUMN IF EXISTS citizenship;
