-- +goose Up

-- Турниры SPB
CREATE TABLE spb_tournaments (
    id SERIAL PRIMARY KEY,
    external_id INT UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    birth_year INT,
    group_name VARCHAR(50),
    season VARCHAR(20),
    start_date DATE,
    end_date DATE,
    is_ended BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_spb_tournaments_season ON spb_tournaments(season);
CREATE INDEX idx_spb_tournaments_birth_year ON spb_tournaments(birth_year);

-- Команды в турнире
CREATE TABLE spb_teams (
    id SERIAL PRIMARY KEY,
    external_id UUID NOT NULL,
    tournament_id INT NOT NULL REFERENCES spb_tournaments(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(external_id, tournament_id)
);

CREATE INDEX idx_spb_teams_tournament ON spb_teams(tournament_id);

-- Профили игроков (уникальные люди)
CREATE TABLE spb_players (
    id SERIAL PRIMARY KEY,
    external_id UUID UNIQUE NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    birth_date DATE,
    birth_place VARCHAR(255),
    current_team_id INT REFERENCES spb_teams(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_spb_players_birth_date ON spb_players(birth_date);
CREATE INDEX idx_spb_players_current_team ON spb_players(current_team_id);

-- Связь: игрок <-> команда в турнире
CREATE TABLE spb_player_teams (
    id SERIAL PRIMARY KEY,
    player_id INT NOT NULL REFERENCES spb_players(id) ON DELETE CASCADE,
    team_id INT NOT NULL REFERENCES spb_teams(id) ON DELETE CASCADE,
    number INT,
    role VARCHAR(5),
    position VARCHAR(20),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(player_id, team_id)
);

CREATE INDEX idx_spb_player_teams_player ON spb_player_teams(player_id);
CREATE INDEX idx_spb_player_teams_team ON spb_player_teams(team_id);

-- +goose Down

DROP TABLE IF EXISTS spb_player_teams;
DROP TABLE IF EXISTS spb_players;
DROP TABLE IF EXISTS spb_teams;
DROP TABLE IF EXISTS spb_tournaments;
