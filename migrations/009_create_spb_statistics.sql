-- +goose Up

-- Статистика полевых игроков в турнире
CREATE TABLE spb_player_statistics (
    id SERIAL PRIMARY KEY,
    player_id INT NOT NULL REFERENCES spb_players(id) ON DELETE CASCADE,
    team_id INT NOT NULL REFERENCES spb_teams(id) ON DELETE CASCADE,
    tournament_id INT NOT NULL REFERENCES spb_tournaments(id) ON DELETE CASCADE,
    
    -- Статистика
    games INT DEFAULT 0,              -- И (игры)
    points INT DEFAULT 0,             -- О (очки)
    points_avg DECIMAL(5,2),          -- О ср.
    goals INT DEFAULT 0,              -- Г (голы)
    assists INT DEFAULT 0,            -- П (передачи)
    plus_minus INT DEFAULT 0,         -- +/-
    penalty_minutes INT DEFAULT 0,    -- Шт (штрафные минуты)
    penalty_avg DECIMAL(5,2),         -- Шт ср.
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(player_id, team_id, tournament_id)
);

CREATE INDEX idx_spb_player_stats_player ON spb_player_statistics(player_id);
CREATE INDEX idx_spb_player_stats_team ON spb_player_statistics(team_id);
CREATE INDEX idx_spb_player_stats_tournament ON spb_player_statistics(tournament_id);
CREATE INDEX idx_spb_player_stats_points ON spb_player_statistics(points DESC);

-- Статистика вратарей в турнире
CREATE TABLE spb_goalie_statistics (
    id SERIAL PRIMARY KEY,
    player_id INT NOT NULL REFERENCES spb_players(id) ON DELETE CASCADE,
    team_id INT NOT NULL REFERENCES spb_teams(id) ON DELETE CASCADE,
    tournament_id INT NOT NULL REFERENCES spb_tournaments(id) ON DELETE CASCADE,
    
    -- Статистика
    games INT DEFAULT 0,              -- И (игры)
    minutes INT DEFAULT 0,            -- Мин. (время на льду)
    goals_against INT DEFAULT 0,      -- Г (пропущенные голы)
    shots_against INT DEFAULT 0,      -- Бр. (броски по воротам)
    save_percentage DECIMAL(5,2),     -- % (процент отражённых)
    goals_against_avg DECIMAL(5,2),   -- Ср. (среднее пропущенных за игру)
    wins INT DEFAULT 0,               -- В (победы)
    shutouts INT DEFAULT 0,           -- На 0 (сухие матчи)
    assists INT DEFAULT 0,            -- П (передачи)
    penalty_minutes INT DEFAULT 0,    -- Шт (штрафные минуты)
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(player_id, team_id, tournament_id)
);

CREATE INDEX idx_spb_goalie_stats_player ON spb_goalie_statistics(player_id);
CREATE INDEX idx_spb_goalie_stats_team ON spb_goalie_statistics(team_id);
CREATE INDEX idx_spb_goalie_stats_tournament ON spb_goalie_statistics(tournament_id);
CREATE INDEX idx_spb_goalie_stats_save_pct ON spb_goalie_statistics(save_percentage DESC);

-- +goose Down

DROP TABLE IF EXISTS spb_goalie_statistics;
DROP TABLE IF EXISTS spb_player_statistics;
