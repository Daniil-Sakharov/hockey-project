-- +goose Up
ALTER TABLE tournaments ADD COLUMN birth_year_groups jsonb;

-- Заполнить из существующих данных player_teams
UPDATE tournaments t SET birth_year_groups = sub.groups
FROM (
    SELECT pt.tournament_id,
        jsonb_object_agg(
            pt.birth_year::text,
            pt.group_names
        ) as groups
    FROM (
        SELECT tournament_id, birth_year,
            jsonb_agg(DISTINCT group_name ORDER BY group_name) as group_names
        FROM player_teams
        WHERE birth_year IS NOT NULL AND birth_year > 0 AND group_name IS NOT NULL
        GROUP BY tournament_id, birth_year
    ) pt
    GROUP BY pt.tournament_id
) sub
WHERE t.id = sub.tournament_id;

-- +goose Down
ALTER TABLE tournaments DROP COLUMN IF EXISTS birth_year_groups;
