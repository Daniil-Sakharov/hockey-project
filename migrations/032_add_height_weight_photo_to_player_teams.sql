-- +goose Up
-- +goose StatementBegin

ALTER TABLE player_teams ADD COLUMN IF NOT EXISTS height INTEGER;
ALTER TABLE player_teams ADD COLUMN IF NOT EXISTS weight INTEGER;
ALTER TABLE player_teams ADD COLUMN IF NOT EXISTS photo_url TEXT;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE player_teams DROP COLUMN IF EXISTS photo_url;
ALTER TABLE player_teams DROP COLUMN IF EXISTS weight;
ALTER TABLE player_teams DROP COLUMN IF EXISTS height;

-- +goose StatementEnd
