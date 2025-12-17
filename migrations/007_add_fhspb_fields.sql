-- +goose Up
-- +goose StatementBegin

-- Добавить внешний ID источника (PlayerID из fhspb.ru)
ALTER TABLE players ADD COLUMN IF NOT EXISTS external_id VARCHAR(100);

-- Добавить гражданство
ALTER TABLE players ADD COLUMN IF NOT EXISTS citizenship VARCHAR(100);

-- Добавить роль в команде (К - капитан, А - ассистент)
ALTER TABLE players ADD COLUMN IF NOT EXISTS role VARCHAR(10);

-- Добавить место рождения
ALTER TABLE players ADD COLUMN IF NOT EXISTS birth_place VARCHAR(255);

-- Индекс для поиска по внешнему ID
CREATE INDEX IF NOT EXISTS idx_players_external_id ON players(external_id);

-- Уникальный индекс для upsert по external_id + source (только для не-null external_id)
CREATE UNIQUE INDEX IF NOT EXISTS idx_players_external_source_unique 
    ON players(external_id, source) WHERE external_id IS NOT NULL;

-- Комментарии
COMMENT ON COLUMN players.source IS 'Источник данных: junior.fhr.ru, fhspb.ru или both';
COMMENT ON COLUMN players.external_id IS 'Внешний ID из источника (PlayerID UUID для fhspb.ru)';
COMMENT ON COLUMN players.citizenship IS 'Гражданство игрока';
COMMENT ON COLUMN players.role IS 'Роль в команде: К - капитан, А - ассистент';
COMMENT ON COLUMN players.birth_place IS 'Место рождения';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_players_external_source_unique;
DROP INDEX IF EXISTS idx_players_external_id;
ALTER TABLE players DROP COLUMN IF EXISTS birth_place;
ALTER TABLE players DROP COLUMN IF EXISTS role;
ALTER TABLE players DROP COLUMN IF EXISTS citizenship;
ALTER TABLE players DROP COLUMN IF EXISTS external_id;

-- +goose StatementEnd
