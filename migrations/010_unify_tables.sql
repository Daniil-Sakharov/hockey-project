-- +goose Up
-- +goose StatementBegin

-- ============================================
-- PLAYERS: делаем поля nullable для SPB
-- ============================================
ALTER TABLE players ALTER COLUMN profile_url DROP NOT NULL;
ALTER TABLE players ALTER COLUMN position DROP NOT NULL;
ALTER TABLE players ADD COLUMN IF NOT EXISTS region TEXT;

-- Уникальный индекс для external_id + source (для SPB upsert)
CREATE UNIQUE INDEX IF NOT EXISTS idx_players_external_id_source 
    ON players(external_id, source) WHERE external_id IS NOT NULL;

-- ============================================
-- TOURNAMENTS: добавляем поля для SPB
-- ============================================
ALTER TABLE tournaments ALTER COLUMN url DROP NOT NULL;
ALTER TABLE tournaments ALTER COLUMN domain DROP NOT NULL;
ALTER TABLE tournaments ADD COLUMN IF NOT EXISTS external_id TEXT;
ALTER TABLE tournaments ADD COLUMN IF NOT EXISTS birth_year INT;
ALTER TABLE tournaments ADD COLUMN IF NOT EXISTS group_name TEXT;
ALTER TABLE tournaments ADD COLUMN IF NOT EXISTS region TEXT;

CREATE UNIQUE INDEX IF NOT EXISTS idx_tournaments_external_id_region 
    ON tournaments(external_id, region) WHERE external_id IS NOT NULL;

-- ============================================
-- TEAMS: добавляем поля для SPB
-- ============================================
ALTER TABLE teams ALTER COLUMN url DROP NOT NULL;
ALTER TABLE teams ADD COLUMN IF NOT EXISTS external_id TEXT;
ALTER TABLE teams ADD COLUMN IF NOT EXISTS tournament_id TEXT REFERENCES tournaments(id) ON DELETE CASCADE;
ALTER TABLE teams ADD COLUMN IF NOT EXISTS region TEXT;

CREATE UNIQUE INDEX IF NOT EXISTS idx_teams_external_id_tournament 
    ON teams(external_id, tournament_id) WHERE external_id IS NOT NULL;

-- ============================================
-- PLAYER_TEAMS: добавляем position для SPB
-- ============================================
ALTER TABLE player_teams ADD COLUMN IF NOT EXISTS position TEXT;

-- ============================================
-- PLAYER_STATISTICS: делаем поля nullable, добавляем avg
-- ============================================
ALTER TABLE player_statistics ALTER COLUMN group_name DROP NOT NULL;
ALTER TABLE player_statistics ALTER COLUMN birth_year DROP NOT NULL;
ALTER TABLE player_statistics ADD COLUMN IF NOT EXISTS points_avg DECIMAL(5,2);
ALTER TABLE player_statistics ADD COLUMN IF NOT EXISTS penalty_avg DECIMAL(5,2);

-- ============================================
-- GOALIE_STATISTICS: новая таблица
-- ============================================
CREATE TABLE IF NOT EXISTS goalie_statistics (
    id SERIAL PRIMARY KEY,
    player_id TEXT NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    team_id TEXT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    tournament_id TEXT NOT NULL REFERENCES tournaments(id) ON DELETE CASCADE,
    
    games INT DEFAULT 0,
    minutes INT DEFAULT 0,
    goals_against INT DEFAULT 0,
    shots_against INT DEFAULT 0,
    save_percentage DECIMAL(5,2),
    goals_against_avg DECIMAL(5,2),
    wins INT DEFAULT 0,
    shutouts INT DEFAULT 0,
    assists INT DEFAULT 0,
    penalty_minutes INT DEFAULT 0,
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    
    UNIQUE(player_id, team_id, tournament_id)
);

CREATE INDEX IF NOT EXISTS idx_goalie_statistics_player ON goalie_statistics(player_id);
CREATE INDEX IF NOT EXISTS idx_goalie_statistics_team ON goalie_statistics(team_id);
CREATE INDEX IF NOT EXISTS idx_goalie_statistics_tournament ON goalie_statistics(tournament_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS goalie_statistics;

ALTER TABLE player_statistics DROP COLUMN IF EXISTS penalty_avg;
ALTER TABLE player_statistics DROP COLUMN IF EXISTS points_avg;

ALTER TABLE player_teams DROP COLUMN IF EXISTS position;

DROP INDEX IF EXISTS idx_teams_external_id_tournament;
ALTER TABLE teams DROP COLUMN IF EXISTS region;
ALTER TABLE teams DROP COLUMN IF EXISTS tournament_id;
ALTER TABLE teams DROP COLUMN IF EXISTS external_id;

DROP INDEX IF EXISTS idx_tournaments_external_id_region;
ALTER TABLE tournaments DROP COLUMN IF EXISTS region;
ALTER TABLE tournaments DROP COLUMN IF EXISTS group_name;
ALTER TABLE tournaments DROP COLUMN IF EXISTS birth_year;
ALTER TABLE tournaments DROP COLUMN IF EXISTS external_id;

DROP INDEX IF EXISTS idx_players_external_id_source;
ALTER TABLE players DROP COLUMN IF EXISTS region;

-- +goose StatementEnd
