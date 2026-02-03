-- +goose Up

-- Таблица для статистики команды за матч (броски по периодам)
CREATE TABLE IF NOT EXISTS match_team_stats (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    match_id TEXT NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    team_id TEXT NOT NULL REFERENCES teams(id),

    -- Броски по периодам
    shots_p1 INTEGER DEFAULT 0,
    shots_p2 INTEGER DEFAULT 0,
    shots_p3 INTEGER DEFAULT 0,
    shots_ot INTEGER DEFAULT 0,
    shots_total INTEGER DEFAULT 0,

    -- Метаданные
    source VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),

    UNIQUE(match_id, team_id)
);

CREATE INDEX idx_match_team_stats_match ON match_team_stats(match_id);
CREATE INDEX idx_match_team_stats_team ON match_team_stats(team_id);

COMMENT ON TABLE match_team_stats IS 'Статистика команды за матч (броски по периодам)';

-- Добавить поля в match_events для счёта в момент гола
ALTER TABLE match_events ADD COLUMN IF NOT EXISTS score_home INTEGER;
ALTER TABLE match_events ADD COLUMN IF NOT EXISTS score_away INTEGER;

COMMENT ON COLUMN match_events.score_home IS 'Счёт домашней команды после события';
COMMENT ON COLUMN match_events.score_away IS 'Счёт гостевой команды после события';

-- +goose Down
ALTER TABLE match_events DROP COLUMN IF EXISTS score_away;
ALTER TABLE match_events DROP COLUMN IF EXISTS score_home;
DROP TABLE IF EXISTS match_team_stats;
