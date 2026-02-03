-- +goose Up

-- Добавляем флаг is_home для событий команды (вратарь, пустые ворота, тайм-аут)
ALTER TABLE match_events ADD COLUMN IF NOT EXISTS is_home BOOLEAN;

COMMENT ON COLUMN match_events.is_home IS 'true = домашняя команда, false = гости (для событий goalie_change, empty_net, timeout)';

-- +goose Down
ALTER TABLE match_events DROP COLUMN IF EXISTS is_home;
