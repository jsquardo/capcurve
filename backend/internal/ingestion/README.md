# Phase 1 Data Shape Notes

## MLB Stats API

Canonical endpoints for v1 ingestion:

- `GET /people/{id}` for player bio
- `GET /people/{id}/stats?stats=yearByYear&group=hitting` for season hitting
- `GET /people/{id}/stats?stats=yearByYear&group=pitching` for season pitching
- `GET /teams?sportId=1` for team reference data

Important response-shape notes:

- Many numeric rate fields arrive as strings, for example `.285`, `1.066`, and `130.1`.
- Placeholder values like `.---` appear and must be treated as missing.
- Two-way players can have hitting and pitching splits for the same `player + season + team`.
- The MLB Stats API season endpoints do not provide `WAR`, `wRC+`, `FIP`, or `ERA+`.

Phase 1 model fields sourced from MLB directly:

- Player bio: MLB ID, names, active status, DOB, bats, throws, primary position
- Hitting: age, games, plate appearances, at-bats, hits, doubles, triples, home runs, runs, RBI, walks, strikeouts, stolen bases, AVG, OBP, SLG, OPS, BABIP
- Pitching: age, games, games started, wins, losses, ERA, WHIP, innings pitched, hits allowed, walks allowed, home runs allowed, strikeouts, K/9, BB/9, H/9, HR/9, K/BB, strike %

## Baseball Savant

Phase 1 enrichment source:

- `GET /leaderboard/statcast?type={batter|pitcher}&year={season}&position=&team=&min=q`

Important response-shape notes:

- The useful payload is embedded in the HTML as `var leaderboard_data = [...]`.
- The embedded data includes both quality-of-contact metrics and expected metrics in the same row set.
- The hitter and pitcher payloads share the same shape, with `xera` populated for pitchers.

Phase 1 enrichment fields sourced from Savant:

- `est_ba` -> `expected_batting_avg`
- `est_slg` -> `expected_slugging`
- `est_woba` -> `expected_woba`
- `xera` -> `expected_era`
- `barrels_per_bip` -> `barrel_pct`
- `hard_hit_percent` -> `hard_hit_pct`
- `exit_velocity_avg` -> `avg_exit_velocity`
- `launch_angle_avg` -> `avg_launch_angle`
- `sweet_spot_percent` -> `sweet_spot_pct`

## Merge Rules

- Persist one `season_stats` row per `player + season + team`.
- Merge hitting and pitching splits into that row when both exist.
- Apply Savant enrichment after the MLB baseline merge.
- Seasons without Savant coverage still persist normally from MLB data alone.
