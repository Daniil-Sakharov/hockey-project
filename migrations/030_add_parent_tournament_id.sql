-- +goose Up
-- Добавляем колонку parent_tournament_id для связи под-турниров с базовым турниром
ALTER TABLE tournaments ADD COLUMN parent_tournament_id VARCHAR(50);

-- Индекс для быстрого поиска под-турниров
CREATE INDEX idx_tournaments_parent_id ON tournaments(parent_tournament_id) WHERE parent_tournament_id IS NOT NULL;

-- Обновляем существующие под-турниры (извлекаем базовый ID из строки)
-- Формат ID: "16731969_y2012_gГруппа_А" -> parent = "16731969"
UPDATE tournaments
SET parent_tournament_id = SPLIT_PART(id, '_y', 1)
WHERE id LIKE '%\_y%' ESCAPE '\'
  AND parent_tournament_id IS NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_tournaments_parent_id;
ALTER TABLE tournaments DROP COLUMN IF EXISTS parent_tournament_id;
