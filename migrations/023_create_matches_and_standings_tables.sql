-- +goose Up

-- ============================================================================
-- Турнирные таблицы (standings)
-- ============================================================================
CREATE TABLE team_standings (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    tournament_id TEXT NOT NULL REFERENCES tournaments(id) ON DELETE CASCADE,
    team_id TEXT NOT NULL REFERENCES teams(id),

    -- Позиция и очки
    position INTEGER,
    points INTEGER DEFAULT 0,

    -- Матчи
    games INTEGER DEFAULT 0,
    wins INTEGER DEFAULT 0,
    wins_ot INTEGER DEFAULT 0,
    wins_so INTEGER DEFAULT 0,
    losses_so INTEGER DEFAULT 0,
    losses_ot INTEGER DEFAULT 0,
    losses INTEGER DEFAULT 0,
    draws INTEGER DEFAULT 0,

    -- Шайбы
    goals_for INTEGER DEFAULT 0,
    goals_against INTEGER DEFAULT 0,
    goal_difference INTEGER DEFAULT 0,

    -- Метаданные
    group_name VARCHAR(50),
    birth_year INTEGER,
    source VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    UNIQUE(tournament_id, team_id, source)
);

CREATE INDEX idx_team_standings_tournament ON team_standings(tournament_id);
CREATE INDEX idx_team_standings_team ON team_standings(team_id);
CREATE INDEX idx_team_standings_source ON team_standings(source);

-- ============================================================================
-- Матчи
-- ============================================================================
CREATE TABLE matches (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    external_id VARCHAR(50) NOT NULL,
    tournament_id TEXT REFERENCES tournaments(id),

    -- Команды
    home_team_id TEXT REFERENCES teams(id),
    away_team_id TEXT REFERENCES teams(id),

    -- Результат
    home_score INTEGER,
    away_score INTEGER,
    home_score_p1 INTEGER,
    away_score_p1 INTEGER,
    home_score_p2 INTEGER,
    away_score_p2 INTEGER,
    home_score_p3 INTEGER,
    away_score_p3 INTEGER,
    home_score_ot INTEGER,
    away_score_ot INTEGER,

    -- Метаданные
    match_number INTEGER,
    scheduled_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) DEFAULT 'scheduled',
    result_type VARCHAR(10),
    venue VARCHAR(255),

    -- Фильтры
    group_name VARCHAR(50),
    birth_year INTEGER,

    -- Видео
    video_url TEXT,

    -- Служебные
    source VARCHAR(50) NOT NULL,
    domain VARCHAR(100),
    details_parsed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    UNIQUE(source, external_id)
);

CREATE INDEX idx_matches_tournament ON matches(tournament_id);
CREATE INDEX idx_matches_home_team ON matches(home_team_id);
CREATE INDEX idx_matches_away_team ON matches(away_team_id);
CREATE INDEX idx_matches_scheduled ON matches(scheduled_at);
CREATE INDEX idx_matches_status ON matches(status);
CREATE INDEX idx_matches_source ON matches(source);

-- ============================================================================
-- События матча (голы, штрафы)
-- ============================================================================
CREATE TABLE match_events (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    match_id TEXT NOT NULL REFERENCES matches(id) ON DELETE CASCADE,

    event_type VARCHAR(20) NOT NULL,
    period INTEGER,
    time_minutes INTEGER,
    time_seconds INTEGER,

    -- Для голов
    scorer_player_id TEXT REFERENCES players(id),
    assist1_player_id TEXT REFERENCES players(id),
    assist2_player_id TEXT REFERENCES players(id),
    team_id TEXT REFERENCES teams(id),
    goal_type VARCHAR(20),

    -- Для штрафов
    penalty_player_id TEXT REFERENCES players(id),
    penalty_minutes INTEGER,
    penalty_reason VARCHAR(100),

    source VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_match_events_match ON match_events(match_id);
CREATE INDEX idx_match_events_scorer ON match_events(scorer_player_id);

-- ============================================================================
-- Составы матча
-- ============================================================================
CREATE TABLE match_lineups (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    match_id TEXT NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    player_id TEXT NOT NULL REFERENCES players(id),
    team_id TEXT NOT NULL REFERENCES teams(id),

    jersey_number INTEGER,
    position VARCHAR(5),

    -- Статистика за матч
    goals INTEGER DEFAULT 0,
    assists INTEGER DEFAULT 0,
    penalty_minutes INTEGER DEFAULT 0,
    plus_minus INTEGER DEFAULT 0,

    -- Для вратарей
    saves INTEGER,
    goals_against INTEGER,
    time_on_ice INTEGER,

    source VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),

    UNIQUE(match_id, player_id)
);

CREATE INDEX idx_match_lineups_match ON match_lineups(match_id);
CREATE INDEX idx_match_lineups_player ON match_lineups(player_id);

-- +goose Down
DROP TABLE IF EXISTS match_lineups;
DROP TABLE IF EXISTS match_events;
DROP TABLE IF EXISTS matches;
DROP TABLE IF EXISTS team_standings;
