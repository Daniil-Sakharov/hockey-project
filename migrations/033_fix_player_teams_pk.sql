-- +goose Up
-- +goose StatementBegin

-- Текущий PK: (player_id, team_id, tournament_id)
-- Проблема: игрок в одной команде/турнире может участвовать в разных группах (А2, А4)
-- UPSERT перезаписывает group_name последним значением, теряя данные

-- 1. Удаляем текущий PK
ALTER TABLE player_teams DROP CONSTRAINT player_teams_pkey;

-- 2. Добавляем суррогатный PK
ALTER TABLE player_teams ADD COLUMN id BIGSERIAL;
ALTER TABLE player_teams ADD PRIMARY KEY (id);

-- 3. Создаём UNIQUE constraint с учётом birth_year и group_name
CREATE UNIQUE INDEX player_teams_uniq_player_team_tournament_year_group
  ON player_teams (player_id, team_id, tournament_id, COALESCE(birth_year, 0), COALESCE(group_name, ''));

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS player_teams_uniq_player_team_tournament_year_group;
ALTER TABLE player_teams DROP CONSTRAINT player_teams_pkey;
ALTER TABLE player_teams DROP COLUMN id;
ALTER TABLE player_teams ADD PRIMARY KEY (player_id, team_id, tournament_id);

-- +goose StatementEnd
