-- +goose Up

-- Добавляем роль капитана/ассистента в составы
ALTER TABLE match_lineups ADD COLUMN IF NOT EXISTS captain_role VARCHAR(1);
-- К = капитан, А = ассистент, NULL = обычный игрок

-- Добавляем список игроков на льду при голе (JSONB массив player_id)
ALTER TABLE match_events ADD COLUMN IF NOT EXISTS home_players_on_ice JSONB;
ALTER TABLE match_events ADD COLUMN IF NOT EXISTS away_players_on_ice JSONB;

-- Добавляем вратаря в события (для голов)
ALTER TABLE match_events ADD COLUMN IF NOT EXISTS goalie_player_id TEXT REFERENCES players(id);

-- Комментарии для понимания структуры
COMMENT ON COLUMN match_lineups.captain_role IS 'К = капитан, А = ассистент, NULL = обычный игрок';
COMMENT ON COLUMN match_events.home_players_on_ice IS 'JSON массив player_id домашней команды на льду при событии';
COMMENT ON COLUMN match_events.away_players_on_ice IS 'JSON массив player_id гостевой команды на льду при событии';
COMMENT ON COLUMN match_events.goalie_player_id IS 'Вратарь при голе (NULL если пустые ворота)';

-- +goose Down
ALTER TABLE match_lineups DROP COLUMN IF EXISTS captain_role;
ALTER TABLE match_events DROP COLUMN IF EXISTS home_players_on_ice;
ALTER TABLE match_events DROP COLUMN IF EXISTS away_players_on_ice;
ALTER TABLE match_events DROP COLUMN IF EXISTS goalie_player_id;
