-- +goose Up
-- Очистка старых данных команд FHSPB (со старым форматом ID spb:<team_id>)
-- Новый формат: spb:<tournament_id>:<team_id>

-- Удаляем статистику игроков для команд FHSPB (будет перепарсена)
DELETE FROM player_statistics WHERE tournament_id LIKE 'spb:%';

-- Удаляем статистику вратарей для команд FHSPB (будет перепарсена)
DELETE FROM goalie_statistics WHERE tournament_id LIKE 'spb:%';

-- Удаляем связи игрок-команда для FHSPB (будут перепарсены)
DELETE FROM player_teams WHERE team_id LIKE 'spb:%';

-- Удаляем команды FHSPB со старым форматом ID
DELETE FROM teams WHERE id LIKE 'spb:%';

-- +goose Down
-- Откат невозможен - данные будут перепарсены
