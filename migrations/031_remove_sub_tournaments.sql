-- +goose Up
-- +goose StatementBegin

-- 1. Добавляем birth_year и group_name в player_teams
ALTER TABLE player_teams ADD COLUMN IF NOT EXISTS birth_year INTEGER;
ALTER TABLE player_teams ADD COLUMN IF NOT EXISTS group_name VARCHAR(50);

-- 2. Обновляем UNIQUE constraint на team_standings
--    Без под-турниров одна команда может иметь standings для разных year/group
--    в одном базовом турнире
ALTER TABLE team_standings DROP CONSTRAINT IF EXISTS team_standings_tournament_id_team_id_source_key;
CREATE UNIQUE INDEX team_standings_uniq_tournament_team_year_group
  ON team_standings (tournament_id, team_id, COALESCE(birth_year, 0), COALESCE(group_name, ''), source);

-- 3. Мигрируем данные: ссылки с под-турниров на базовые турниры
UPDATE matches SET tournament_id = SPLIT_PART(tournament_id, '_y', 1)
  WHERE tournament_id LIKE '%\_y%' ESCAPE '\';

UPDATE team_standings SET tournament_id = SPLIT_PART(tournament_id, '_y', 1)
  WHERE tournament_id LIKE '%\_y%' ESCAPE '\';

UPDATE player_teams SET tournament_id = SPLIT_PART(tournament_id, '_y', 1)
  WHERE tournament_id LIKE '%\_y%' ESCAPE '\';

UPDATE player_statistics SET tournament_id = SPLIT_PART(tournament_id, '_y', 1)
  WHERE tournament_id LIKE '%\_y%' ESCAPE '\';

UPDATE teams SET tournament_id = SPLIT_PART(tournament_id, '_y', 1)
  WHERE tournament_id LIKE '%\_y%' ESCAPE '\';

-- 4. Удаляем под-турниры
DELETE FROM tournaments WHERE parent_tournament_id IS NOT NULL;

-- 5. Индексы
CREATE INDEX IF NOT EXISTS idx_player_teams_birth_year ON player_teams(birth_year) WHERE birth_year IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_player_teams_group_name ON player_teams(group_name) WHERE group_name IS NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE player_teams DROP COLUMN IF EXISTS group_name;
ALTER TABLE player_teams DROP COLUMN IF EXISTS birth_year;
DROP INDEX IF EXISTS idx_player_teams_birth_year;
DROP INDEX IF EXISTS idx_player_teams_group_name;
DROP INDEX IF EXISTS team_standings_uniq_tournament_team_year_group;
ALTER TABLE team_standings ADD CONSTRAINT team_standings_tournament_id_team_id_source_key
  UNIQUE(tournament_id, team_id, source);

-- +goose StatementEnd
