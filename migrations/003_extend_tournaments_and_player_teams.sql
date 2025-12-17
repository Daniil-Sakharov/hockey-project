-- +goose Up
-- +goose StatementBegin

-- Расширяем таблицу tournaments
ALTER TABLE tournaments ADD COLUMN IF NOT EXISTS season TEXT;
ALTER TABLE tournaments ADD COLUMN IF NOT EXISTS start_date TIMESTAMP;
ALTER TABLE tournaments ADD COLUMN IF NOT EXISTS end_date TIMESTAMP;
ALTER TABLE tournaments ADD COLUMN IF NOT EXISTS is_ended BOOLEAN DEFAULT false;

-- Расширяем таблицу player_teams
ALTER TABLE player_teams ADD COLUMN IF NOT EXISTS season TEXT;
ALTER TABLE player_teams ADD COLUMN IF NOT EXISTS started_at TIMESTAMP;
ALTER TABLE player_teams ADD COLUMN IF NOT EXISTS ended_at TIMESTAMP;
ALTER TABLE player_teams ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT true;
ALTER TABLE player_teams ADD COLUMN IF NOT EXISTS jersey_number INTEGER;
ALTER TABLE player_teams ADD COLUMN IF NOT EXISTS role TEXT;
ALTER TABLE player_teams ADD COLUMN IF NOT EXISTS source TEXT NOT NULL DEFAULT 'junior';
ALTER TABLE player_teams ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP NOT NULL DEFAULT NOW();

-- Создаем индексы для быстрых запросов
CREATE INDEX IF NOT EXISTS idx_player_teams_player ON player_teams(player_id);
CREATE INDEX IF NOT EXISTS idx_player_teams_team ON player_teams(team_id);
CREATE INDEX IF NOT EXISTS idx_player_teams_tournament ON player_teams(tournament_id);
CREATE INDEX IF NOT EXISTS idx_player_teams_active ON player_teams(player_id) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_player_teams_season ON player_teams(season);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Удаляем индексы
DROP INDEX IF EXISTS idx_player_teams_season;
DROP INDEX IF EXISTS idx_player_teams_active;
DROP INDEX IF EXISTS idx_player_teams_tournament;
DROP INDEX IF EXISTS idx_player_teams_team;
DROP INDEX IF EXISTS idx_player_teams_player;

-- Откатываем изменения player_teams
ALTER TABLE player_teams DROP COLUMN IF EXISTS updated_at;
ALTER TABLE player_teams DROP COLUMN IF EXISTS source;
ALTER TABLE player_teams DROP COLUMN IF EXISTS role;
ALTER TABLE player_teams DROP COLUMN IF EXISTS jersey_number;
ALTER TABLE player_teams DROP COLUMN IF EXISTS is_active;
ALTER TABLE player_teams DROP COLUMN IF EXISTS ended_at;
ALTER TABLE player_teams DROP COLUMN IF EXISTS started_at;
ALTER TABLE player_teams DROP COLUMN IF EXISTS season;

-- Откатываем изменения tournaments
ALTER TABLE tournaments DROP COLUMN IF EXISTS is_ended;
ALTER TABLE tournaments DROP COLUMN IF EXISTS end_date;
ALTER TABLE tournaments DROP COLUMN IF EXISTS start_date;
ALTER TABLE tournaments DROP COLUMN IF EXISTS season;

-- +goose StatementEnd
