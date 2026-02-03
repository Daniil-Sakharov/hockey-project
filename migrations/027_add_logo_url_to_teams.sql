-- +goose Up
-- +goose StatementBegin

-- Добавляем поле logo_url в таблицу teams для хранения URL логотипа команды
ALTER TABLE teams ADD COLUMN IF NOT EXISTS logo_url TEXT;

COMMENT ON COLUMN teams.logo_url IS 'URL логотипа команды';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE teams DROP COLUMN IF EXISTS logo_url;

-- +goose StatementEnd
