-- +goose Up
-- +goose StatementBegin
-- ==============================================================================
-- Миграция 005: Добавление UNIQUE constraint для дедупликации статистики
-- ==============================================================================
-- Добавляет уникальный индекс чтобы предотвратить дубликаты записей статистики
-- для одного игрока в одном турнире/группе/году.
-- 
-- Логика дедупликации:
-- - Один игрок МОЖЕТ быть в разных группах одного турнира/года (Группа А, Группа Б)
-- - Один игрок МОЖЕТ быть в одной группе разных годов (2015, 2016)
-- - Один игрок НЕ МОЖЕТ дублироваться в одной комбинации турнир+группа+год
-- ==============================================================================

-- Удаляем возможные дубликаты перед созданием индекса
-- Оставляем только последнюю запись для каждой комбинации (по id)
DELETE FROM player_statistics
WHERE id IN (
    SELECT id
    FROM (
        SELECT id,
               ROW_NUMBER() OVER (
                   PARTITION BY tournament_id, player_id, group_name, birth_year
                   ORDER BY created_at DESC, id DESC
               ) as rn
        FROM player_statistics
    ) t
    WHERE t.rn > 1
);

-- Создаем уникальный индекс
CREATE UNIQUE INDEX idx_unique_player_stat 
ON player_statistics(tournament_id, player_id, group_name, birth_year);

COMMENT ON INDEX idx_unique_player_stat IS 
    'Уникальный индекс для предотвращения дубликатов статистики: один игрок не может иметь несколько записей для одной комбинации турнир+группа+год';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_unique_player_stat;
-- +goose StatementEnd
