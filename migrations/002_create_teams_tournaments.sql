-- +goose Up
-- +goose StatementBegin

-- Таблица турниров
CREATE TABLE tournaments (
    id TEXT PRIMARY KEY,  -- ID из URL (например "16756891")
    url TEXT UNIQUE NOT NULL,  -- /tournaments/pervenstvo-tsfo-18171615-let-16756891/
    name TEXT NOT NULL,
    domain TEXT NOT NULL,  -- https://cfo.fhr.ru
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tournaments_url ON tournaments(url);
CREATE INDEX idx_tournaments_domain ON tournaments(domain);

COMMENT ON TABLE tournaments IS 'Турниры из junior.fhr.ru';
COMMENT ON COLUMN tournaments.id IS 'ID извлеченный из URL турнира';
COMMENT ON COLUMN tournaments.domain IS 'Домен откуда был спарсен турнир';

-- Таблица команд
CREATE TABLE teams (
    id TEXT PRIMARY KEY,  -- ID из URL или hash
    url TEXT UNIQUE NOT NULL,  -- /tournaments/.../buran_5136295/
    name TEXT NOT NULL,
    city TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_teams_url ON teams(url);
CREATE INDEX idx_teams_name ON teams(name);

COMMENT ON TABLE teams IS 'Команды из junior.fhr.ru';

-- Связь многие-ко-многим: игрок-команда-турнир
CREATE TABLE player_teams (
    player_id TEXT NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    team_id TEXT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    tournament_id TEXT NOT NULL REFERENCES tournaments(id) ON DELETE CASCADE,
    number TEXT,  -- Номер игрока в команде (может меняться)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    PRIMARY KEY (player_id, team_id, tournament_id)
);

CREATE INDEX idx_player_teams_player ON player_teams(player_id);
CREATE INDEX idx_player_teams_team ON player_teams(team_id);
CREATE INDEX idx_player_teams_tournament ON player_teams(tournament_id);

COMMENT ON TABLE player_teams IS 'Связь игрок-команда-турнир (многие-ко-многим)';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS player_teams;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS tournaments;
-- +goose StatementEnd
