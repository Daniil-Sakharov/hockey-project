-- +goose Up
-- Расширяем поле group_name для матчей (некоторые названия групп длиннее 50 символов)
ALTER TABLE matches ALTER COLUMN group_name TYPE VARCHAR(255);

-- +goose Down
ALTER TABLE matches ALTER COLUMN group_name TYPE VARCHAR(50);
