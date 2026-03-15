## [2026-03-15] — Session Summary

### Added
- Added migration `000004_enforce_active_season_stat_uniqueness` to enforce one active `season_stats` row per `player_id + year` with a partial unique index on non-deleted rows
- Added unit coverage proving the ingestion upsert targets active `player_id + year` conflicts and keeps the expected mutable season columns

### Changed
- Updated ingestion upserts to target the new active-row partial unique index instead of the old `player_id + year + team_id` uniqueness rule
- Updated ingestion conflict updates so `team_id` can change when a re-ingested traded season ends with a different canonical final-team split

### Fixed
- Prevented soft-deleted `season_stats` tombstones from blocking re-ingestion under the new one-row-per-player-year rule

### Notes
- Applied the migration successfully in Docker, verified migration version `4`, and confirmed Postgres created `idx_season_stats_player_year_active_unique` with `WHERE deleted_at IS NULL`
- Ran `make test` successfully

## [2026-03-15] — Session Summary

### Added
- Added a regression test proving hitter base scoring ignores derived `OPS` changes when the non-derivative hitter inputs stay the same
- Added a regression test proving higher `BarrelPct` and `HardHitPct` outrank a SweetSpot-only improvement when the rest of a hitter season is unchanged

### Changed
- Reweighted hitter base scoring to use `OBP 0.28`, `SLG 0.24`, `BattingAvg 0.08`, `HomeRuns 0.10`, `RBI 0.07`, `StolenBases 0.05`, and `BABIP 0.05`
- Reweighted hitter Savant scoring to emphasize `BarrelPct` and `HardHitPct` above `SweetSpotPct` while keeping expected outcome metrics and exit velocity in the blend
- Updated `AGENTS.md` Phase 1 scoring guidance to record the permanent hitter base-weight rule and note that WAR should eventually be added at `0.05`-`0.10` once the data model supports it
- Updated `AGENTS.md` Phase 1 scoring guidance to record the permanent differentiated hitter Savant weights

### Fixed
- Removed hitter-score double counting caused by scoring `OBP`, `SLG`, and derived `OPS` together

### Notes
- Ran `make test` successfully

## [2026-03-14] — Session Summary

### Added
- Added regression tests for baseball-notation innings parsing, pitching split normalization into outs, and partial-inning traded-season merges

### Changed
- Updated ingestion normalization to convert MLB innings strings like `130.1` and `130.2` into integer outs before any merge or scoring math
- Updated merged pitching-season recomputation to use outs-based innings internally while still persisting `innings_pitched` in MLB-style baseball notation for schema compatibility
- Updated `AGENTS.md` Current State to mark the innings-pitched blocker complete and leave the schema-enforcement migration as the next Phase 1 task

### Fixed
- Corrected partial-inning aggregation so merged pitching rates no longer treat baseball notation tenths as decimal fractions

### Notes
- Internal ingestion records now keep `InningsPitchedOuts` as the authoritative pitching-workload value; the existing model/schema field remains unchanged for now
- Ran `make test` successfully

## [2026-03-11] — Session Summary

### Added
- Added the first `CHANGELOG.md` session record for cross-session tracking

### Changed
- Seeded the development database with a representative 20-player Phase 1 validation set using the existing ingestion CLI
- Updated `AGENTS.md` Current State to reflect the completed seed-and-verify pass and the exact next follow-up task

### Fixed
- Confirmed migrations `000002` and `000003` are already applied by verifying the database is at migration version `3`

### Notes
- Verified `players.position` accepts `"Two-Way Player"` and Ohtani persists as a merged two-way player
- Verified seeded rows have non-zero/non-null `value_score`, current pure pitchers have Savant pitcher enrichment, and older pre-Statcast seasons persist without Savant data
- Identified a blocker: multi-team seasons persist an extra aggregate row with `team_id = 0`; there are currently `7` such rows and this needs to be fixed before Phase 2

## [2026-03-12] — Session Summary

### Added
- Added regression tests covering aggregate `team_id = 0` split detection and exclusion during ingestion merge

### Changed
- Updated ingestion normalization to skip MLB aggregate `TOT` season rows before persistence
- Updated player sync persistence to soft-delete any legacy `season_stats` rows with `team_id = 0` for the player being re-ingested
- Re-ran ingestion for Shohei Ohtani, Juan Soto, Albert Pujols, Justin Verlander, Ken Griffey Jr., and Greg Maddux to clean up existing aggregate rows
- Updated `AGENTS.md` Current State to reflect the completed aggregate-row fix and the next scoring-design task

### Fixed
- Removed all active aggregate multi-team rows from `season_stats`; Postgres verification now shows `0` active `team_id = 0` rows and no duplicate `player + year + team` records

### Notes
- Active `season_stats` row count is now `305`, with `7` legacy aggregate rows preserved only as soft-deleted tombstones
- Ohtani's `2025` merged row still has `expected_era = NULL`, and the live `2025` Savant pitcher leaderboard does not include `entity_id = 660271`, so this currently looks like an upstream coverage gap rather than an ingestion bug
- Ran `GOCACHE=/tmp/capcurve-gocache go test -mod=mod ./...` successfully from `backend/`

## [2026-03-13] — Session Summary

### Added
- Added regression tests covering traded multi-team season aggregation and weighted pitching-rate recomputation across team splits

### Changed
- Changed ingestion season merging to collapse real traded splits into one persisted `season_stats` row per `player + year`
- Recomputed merged hitting and pitching rate stats from the combined split totals instead of carrying a single team split's rates forward
- Kept the final real team split as the canonical `team_id` / `team_name` on merged traded seasons and soft-deleted superseded split rows during re-ingest
- Re-ran ingestion for Juan Soto, Albert Pujols, Justin Verlander, Ken Griffey Jr., and Greg Maddux to rewrite the affected traded seasons
- Updated `AGENTS.md` and ingestion docs to record the permanent traded-season merge rule and the next remaining Phase 1 scoring tasks

### Fixed
- Eliminated active duplicate `season_stats` rows at the `player + year` grain; traded seasons now carry exactly one active `value_score` row per player-season

### Notes
- Postgres verification now shows `298` active `season_stats` rows, `14` soft-deleted historical rows, and `0` active duplicate player-year rows
- Canonical team metadata on merged traded seasons intentionally reflects the player's final real team split, not MLB's synthetic `TOT` row
- Ran `GOCACHE=/tmp/capcurve-gocache go test ./internal/ingestion/...` and `make test` successfully

## [2026-03-13] — Session Summary

### Added
- Added regression tests proving sub-threshold hitter and pitcher seasons are excluded from percentile cohorts while eligible full-season peers keep stable percentile anchors

### Changed
- Updated season scoring so hitter percentile cohorts now require at least `100` plate appearances and pitcher percentile cohorts require at least `30.0` innings pitched
- Updated `AGENTS.md` Phase 1 scoring design guidance to record the percentile threshold policy as a permanent rule instead of a temporary Current State note

### Fixed
- Prevented tiny-sample seasons from distorting season-scoped percentile rankings used for `value_score`

### Notes
- Sub-threshold seasons still receive a dampened score; they no longer define the percentile baseline for fuller workloads
- Ran `make test` successfully

## [2026-03-15] — Session Summary

### Added
- Added regression tests covering baseball-notation innings conversion inside scoring, partial-inning pitcher dampening, and two-way final-score weighting

### Changed
- Updated scoring workload math to convert persisted MLB baseball-notation innings into true outs-based innings before applying pitcher thresholds, dampeners, and two-way role weighting
- Updated `AGENTS.md` Current State to mark the innings-pitched scoring blocker complete and set the active `season_stats` uniqueness migration as the next task

### Fixed
- Corrected scoring-side workload handling so seasons like `29.2` innings are no longer underweighted as decimal tenths during `value_score` calculation

### Notes
- Persisted `season_stats.innings_pitched` still remains in MLB baseball notation for schema compatibility; scoring now normalizes it back to true innings internally
- Ran `make test` successfully
