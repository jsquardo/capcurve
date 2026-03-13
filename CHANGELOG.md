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
