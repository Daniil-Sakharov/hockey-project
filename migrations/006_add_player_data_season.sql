-- +goose Up
-- +goose StatementBegin

-- Добавляем поле data_season для хранения сезона из которого взяты данные игрока (рост, вес, хват)
-- Это позволяет обновлять данные только из более свежих турниров
ALTER TABLE players ADD COLUMN data_season TEXT;

-- Комментарий к полю
COMMENT ON COLUMN players.data_season IS 'Сезон турнира из которого взяты актуальные данные игрока (рост, вес, хват). Формат: 2024/2025';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE players DROP COLUMN IF EXISTS data_season;
-- +goose StatementEnd
