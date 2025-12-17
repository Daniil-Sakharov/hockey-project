-- +goose Up
-- +goose StatementBegin

DROP TABLE IF EXISTS spb_goalie_statistics;
DROP TABLE IF EXISTS spb_player_statistics;
DROP TABLE IF EXISTS spb_player_teams;
DROP TABLE IF EXISTS spb_players;
DROP TABLE IF EXISTS spb_teams;
DROP TABLE IF EXISTS spb_tournaments;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Восстановление SPB таблиц (если нужен откат)
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

CREATE TABLE spb_teams (
    id SERIAL PRIMARY KEY,
    external_id UUID NOT NULL,
    tournament_id INT NOT NULL REFERENCES spb_tournaments(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(external_id, tournament_id)
);

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

CREATE TABLE spb_player_statistics (
    id SERIAL PRIMARY KEY,
    player_id INT NOT NULL REFERENCES spb_players(id) ON DELETE CASCADE,
    team_id INT NOT NULL REFERENCES spb_teams(id) ON DELETE CASCADE,
    tournament_id INT NOT NULL REFERENCES spb_tournaments(id) ON DELETE CASCADE,
    games INT DEFAULT 0,
    points INT DEFAULT 0,
    points_avg DECIMAL(5,2),
    goals INT DEFAULT 0,
    assists INT DEFAULT 0,
    plus_minus INT DEFAULT 0,
    penalty_minutes INT DEFAULT 0,
    penalty_avg DECIMAL(5,2),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(player_id, team_id, tournament_id)
);

CREATE TABLE spb_goalie_statistics (
    id SERIAL PRIMARY KEY,
    player_id INT NOT NULL REFERENCES spb_players(id) ON DELETE CASCADE,
    team_id INT NOT NULL REFERENCES spb_teams(id) ON DELETE CASCADE,
    tournament_id INT NOT NULL REFERENCES spb_tournaments(id) ON DELETE CASCADE,
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
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(player_id, team_id, tournament_id)
);

-- +goose StatementEnd
