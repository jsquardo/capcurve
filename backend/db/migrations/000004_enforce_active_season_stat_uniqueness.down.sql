DROP INDEX IF EXISTS idx_season_stats_player_year_active_unique;

-- Rolling back to the old (player_id, year, team_id) uniqueness rule must
-- first move colliding soft-deleted tombstones out of that keyspace so active
-- rows and preserved history can coexist during the rollback.
WITH ranked_team_keys AS (
    SELECT
        id,
        ROW_NUMBER() OVER (
            PARTITION BY player_id, year, team_id
            ORDER BY
                CASE WHEN deleted_at IS NULL THEN 0 ELSE 1 END,
                deleted_at DESC NULLS LAST,
                id DESC
        ) AS team_key_rank
    FROM season_stats
    WHERE team_id IS NOT NULL
),
conflicting_tombstones AS (
    SELECT id
    FROM ranked_team_keys
    WHERE team_key_rank > 1
)
UPDATE season_stats
SET
    team_id = NULL,
    team_name = CONCAT(COALESCE(team_name, 'Unknown Team'), ' (rollback tombstone)'),
    updated_at = NOW()
WHERE id IN (SELECT id FROM conflicting_tombstones)
  AND deleted_at IS NOT NULL;

ALTER TABLE season_stats
    ADD CONSTRAINT season_stats_player_id_year_team_id_key UNIQUE (player_id, year, team_id);
