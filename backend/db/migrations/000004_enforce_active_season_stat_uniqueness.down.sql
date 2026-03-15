DROP INDEX IF EXISTS idx_season_stats_player_year_active_unique;

ALTER TABLE season_stats
    ADD CONSTRAINT season_stats_player_id_year_team_id_key UNIQUE (player_id, year, team_id);
