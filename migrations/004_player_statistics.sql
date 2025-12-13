-- +goose Up
-- +goose StatementBegin
-- ==============================================================================
-- Миграция 004: Таблица статистики игроков
-- ==============================================================================
-- Создаёт таблицу для хранения детальной статистики игроков по турнирам
-- с разбивкой по группам и годам рождения.
-- ==============================================================================

CREATE TABLE IF NOT EXISTS player_statistics (
    id SERIAL PRIMARY KEY,
    
    -- ==============================================================================
    -- Ключевые связи
    -- ==============================================================================
    tournament_id TEXT NOT NULL REFERENCES tournaments(id) ON DELETE CASCADE,
    player_id TEXT NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    team_id TEXT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    
    -- ==============================================================================
    -- Контекст статистики
    -- ==============================================================================
    group_name TEXT NOT NULL,           -- "Общая статистика", "Группа А", "Группа Б"
    birth_year INT NOT NULL,            -- 2008, 2009, 2010, 2011...
    
    -- ==============================================================================
    -- Основная статистика
    -- ==============================================================================
    games INT NOT NULL DEFAULT 0,                   -- И   - сыгранные матчи
    goals INT NOT NULL DEFAULT 0,                   -- Ш   - заброшенные шайбы
    assists INT NOT NULL DEFAULT 0,                 -- А   - передачи
    points INT NOT NULL DEFAULT 0,                  -- О   - очки (гол + пас)
    plus INT NOT NULL DEFAULT 0,                    -- +   - плюс
    minus INT NOT NULL DEFAULT 0,                   -- -   - минус
    plus_minus INT NOT NULL DEFAULT 0,              -- +/- - коэффициент полезности
    penalty_minutes INT NOT NULL DEFAULT 0,         -- ШТР - штрафное время
    
    -- ==============================================================================
    -- Детальная статистика голов
    -- ==============================================================================
    goals_even_strength INT NOT NULL DEFAULT 0,     -- ШР  - шайб в равенстве
    goals_power_play INT NOT NULL DEFAULT 0,        -- ШБ  - шайб в большинстве
    goals_short_handed INT NOT NULL DEFAULT 0,      -- ШМ  - шайб в меньшинстве
    goals_period_1 INT NOT NULL DEFAULT 0,          -- Ш1п - шайб в 1 периоде
    goals_period_2 INT NOT NULL DEFAULT 0,          -- Ш2п - шайб в 2 периоде
    goals_period_3 INT NOT NULL DEFAULT 0,          -- Ш3п - шайб в 3 периоде
    goals_overtime INT NOT NULL DEFAULT 0,          -- ШОт - шайб в овертайме
    hat_tricks INT NOT NULL DEFAULT 0,              -- ХТ  - хет-трики
    game_winning_goals INT NOT NULL DEFAULT 0,      -- РБ  - решающие буллиты
    
    -- ==============================================================================
    -- Средние показатели
    -- ==============================================================================
    goals_per_game DECIMAL(5,2) NOT NULL DEFAULT 0.00,          -- Ш/И
    points_per_game DECIMAL(5,2) NOT NULL DEFAULT 0.00,         -- О/И
    penalty_minutes_per_game DECIMAL(5,2) NOT NULL DEFAULT 0.00, -- ШТР/И
    
    -- ==============================================================================
    -- Метаданные
    -- ==============================================================================
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ==============================================================================
-- Индексы
-- ==============================================================================

-- Основной индекс для поиска статистики игрока
CREATE INDEX idx_player_statistics_player_id 
    ON player_statistics(player_id);

-- Индекс для поиска по турниру
CREATE INDEX idx_player_statistics_tournament_id 
    ON player_statistics(tournament_id);

-- Индекс для поиска по команде
CREATE INDEX idx_player_statistics_team_id 
    ON player_statistics(team_id);

-- Composite индекс для дедупликации и быстрого поиска
CREATE INDEX idx_player_statistics_composite 
    ON player_statistics(tournament_id, player_id, group_name, birth_year);

-- Индекс для поиска по году рождения
CREATE INDEX idx_player_statistics_birth_year 
    ON player_statistics(birth_year);

-- ==============================================================================
-- Комментарии
-- ==============================================================================

COMMENT ON TABLE player_statistics IS 
    'Детальная статистика игроков по турнирам с разбивкой по группам и годам';

COMMENT ON COLUMN player_statistics.tournament_id IS 
    'ID турнира из таблицы tournaments';

COMMENT ON COLUMN player_statistics.player_id IS 
    'ID игрока из таблицы players';

COMMENT ON COLUMN player_statistics.team_id IS 
    'ID команды из таблицы teams';

COMMENT ON COLUMN player_statistics.group_name IS 
    'Название группы: "Общая статистика", "Группа А", "Группа Б" и т.д.';

COMMENT ON COLUMN player_statistics.birth_year IS 
    'Группа года рождения (может не совпадать с реальным годом рождения игрока)';

COMMENT ON COLUMN player_statistics.games IS 
    'Количество сыгранных матчей';

COMMENT ON COLUMN player_statistics.goals IS 
    'Количество заброшенных шайб';

COMMENT ON COLUMN player_statistics.assists IS 
    'Количество передач';

COMMENT ON COLUMN player_statistics.points IS 
    'Очки (голы + передачи)';

COMMENT ON COLUMN player_statistics.penalty_minutes IS 
    'Штрафное время в минутах';

COMMENT ON COLUMN player_statistics.goals_even_strength IS 
    'Шайбы забитые в равенстве';

COMMENT ON COLUMN player_statistics.goals_power_play IS 
    'Шайбы забитые в большинстве';

COMMENT ON COLUMN player_statistics.goals_short_handed IS 
    'Шайбы забитые в меньшинстве';

COMMENT ON COLUMN player_statistics.hat_tricks IS 
    'Количество хет-триков (3+ гола за матч)';

COMMENT ON COLUMN player_statistics.game_winning_goals IS 
    'Решающие голы/буллиты';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS player_statistics;
-- +goose StatementEnd
