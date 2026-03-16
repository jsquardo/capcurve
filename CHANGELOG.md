## [2026-03-16] — Session Summary

### Added
- Added a post-ordering traded-player ingestion spot check covering Juan Soto (`665742`), Justin Verlander (`434378`), and Albert Pujols (`405395`)

### Changed
- Re-ran the ingestion CLI for the representative traded-player set to confirm the deterministic split-ordering work still preserves canonical final-team metadata on merged seasons

### Fixed
- Confirmed the updated active-row unique index and upsert target still leave exactly one active `season_stats` row per `player_id + year` for the spot-checked traded seasons

### Notes
- Postgres verification showed one active row each for Soto `2022`, Verlander `2017`, and Pujols `2021`, with no active duplicate player-year rows for the spot-checked players

## [2026-03-16] — Session Summary

### Added
- Added stricter ingestion upsert regression coverage that asserts the full mutable `season_stats` update column set

### Changed
- Updated the ingestion upsert clause test to compare the emitted conflict-update columns against the complete declared `seasonStatUpsertColumns` list instead of spot-checking representative fields

### Fixed
- Closed a regression gap where future mutable `season_stats` columns could have been added or removed from the upsert clause without the test noticing

### Notes
- Ran `make test` successfully

## [2026-03-16] — Session Summary

### Added
- Added rollback-time SQL handling in migration `000004` for soft-deleted `season_stats` tombstones that would otherwise collide with the restored legacy `(player_id, year, team_id)` uniqueness rule

### Changed
- Updated the `000004` down migration to null only conflicting soft-deleted tombstone team keys before re-adding the old unique constraint

### Fixed
- Prevented rollback of migration `000004` from failing when preserved tombstones share a `(player_id, year, team_id)` key with active data or newer tombstones

### Notes
- Verified the rollback SQL against a temporary Postgres table containing both active-plus-tombstone and tombstone-only duplicate team keys, then rolled the verification transaction back
- Ran `make test` successfully

## [2026-03-16] — Session Summary

### Added
- Added regression coverage proving traded-season merges keep the higher split `Age` when MLB rows disagree across an in-season birthday

### Changed
- Documented the traded-season age merge rule directly in ingestion so future cleanup work does not have to infer why `Age` uses a max merge

### Fixed
- Removed ambiguity around traded-season age handling by making the end-of-season-age rule explicit in code and tests

### Notes
- MLB split payloads can report different ages within one season when a trade happens before and after a birthday; merged season rows now document that they intentionally keep the higher age
- Ran `make test` successfully

## [2026-03-16] — Session Summary

### Added
- Added regression coverage for deterministic traded-season split ordering and canonical team metadata retention across out-of-order merge calls

### Changed
- Updated ingestion split merging to carry explicit MLB source split order through normalization so canonical traded-season team metadata no longer depends on incidental merge overwrite order
- Forced aggregate `TOT` rows behind real team splits inside ingestion ordering so they cannot compete for canonical team metadata

### Fixed
- Removed the implicit "last merge wins" behavior that previously decided traded-season canonical team metadata without recording why that team won

### Notes
- The MLB year-by-year payload does not expose a dedicated trade timestamp, so ingestion now treats source split order as the explicit chronology signal for intra-season team ordering
- Ran `make test` successfully

## [2026-03-16] — Session Summary

### Added
- Added a shared `internal/baseball` helper package for MLB baseball-notation innings parsing and outs conversion
- Added regression coverage for shared innings parsing, float-based persisted innings validation, and scoring rejection of malformed pitcher workloads

### Changed
- Updated ingestion innings parsing and outs-to-notation conversion to delegate to the shared helper instead of maintaining duplicate rules
- Updated scoring workload normalization to use the shared helper and reject malformed persisted innings values such as `10.4` as zero workload

### Fixed
- Eliminated the scoring-side fallback that previously treated malformed baseball-notation innings as valid decimal innings

### Notes
- Ran `make test` successfully

## [2026-03-16] — Session Summary

### Added
- Added the permanent Phase 1 rule documenting database-enforced active `season_stats` uniqueness at the `player_id + year` grain

### Changed
- Updated `AGENTS.md` Current State priorities to reflect the review findings, keeping malformed baseball-notation innings as the top Phase 1 blocker
- Added the review follow-up tasks for rollback safety and stronger ingestion upsert coverage to the Current State secondary task list

### Fixed
- None

### Notes
- Documentation-only session; no code changes were made

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
