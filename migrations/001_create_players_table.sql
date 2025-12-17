-- +goose Up
-- +goose StatementBegin
CREATE TABLE players (
    id TEXT PRIMARY KEY,  -- ID извлеченный из URL (например "924040")
    profile_url TEXT UNIQUE NOT NULL,  -- /player/...-924040/
    name TEXT NOT NULL,
    birth_date DATE NOT NULL,
    position TEXT NOT NULL,  -- Защитник/Нападающий/Вратарь
    height INT,
    weight INT,
    handedness TEXT,  -- Левый/Правый
    
    -- Данные из registrynew.fhr.ru (nullable, будут добавлены позже)
    registry_id TEXT,
    school TEXT,
    rank TEXT,
    
    source TEXT NOT NULL DEFAULT 'junior.fhr.ru',  -- junior.fhr.ru или both
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Индексы для быстрого поиска
CREATE UNIQUE INDEX idx_players_profile_url ON players(profile_url);
CREATE INDEX idx_players_name ON players USING gin(to_tsvector('russian', name));  -- Full-text search
CREATE INDEX idx_players_birth_date ON players(birth_date);
CREATE INDEX idx_players_position ON players(position);

-- Комментарии для документации
COMMENT ON TABLE players IS 'Игроки из junior.fhr.ru и registrynew.fhr.ru';
COMMENT ON COLUMN players.id IS 'ID из URL профиля (924040)';
COMMENT ON COLUMN players.profile_url IS 'Уникальный URL профиля игрока';
COMMENT ON COLUMN players.source IS 'Источник данных: junior.fhr.ru или both';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS players;
-- +goose StatementEnd
