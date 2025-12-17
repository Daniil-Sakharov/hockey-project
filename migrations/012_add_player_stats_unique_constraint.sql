-- +goose Up
-- +goose StatementBegin

-- Добавляем уникальный индекс для FHSPB статистики (без group_name и birth_year)
CREATE UNIQUE INDEX IF NOT EXISTS idx_player_statistics_unique_simple 
    ON player_statistics(player_id, team_id, tournament_id) 
    WHERE group_name IS NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_player_statistics_unique_simple;

-- +goose StatementEnd
