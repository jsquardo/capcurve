ALTER TABLE season_stats
    DROP CONSTRAINT IF EXISTS season_stats_player_id_year_team_id_key;

DROP INDEX IF EXISTS idx_season_stats_player_year_active_unique;

CREATE UNIQUE INDEX idx_season_stats_player_year_active_unique
    ON season_stats (player_id, year)
    WHERE deleted_at IS NULL;
