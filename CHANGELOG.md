## [2026-03-30] — Session Summary

### Added
- Added backend integration coverage proving `GET /api/v1/players` returns explicit page-based metadata and the correct second-page slice for `page=2&page_size=1`

### Changed
- Updated `GET /api/v1/players` to return `{ data, meta: { total, page, page_size, total_pages } }` so the frontend can build a full paginator from a single response
- Updated the player-list handler to accept `page` and `page_size` as first-class pagination inputs while still deriving equivalent paging from legacy `limit` and `offset` values when needed
- Updated the shared frontend `PlayerListResponse` type to match the explicit page-based pagination metadata
- Updated `AGENTS.md` Current State so the next session resumes with tightening projection comparable-player candidate queries

### Fixed
- Removed the old list-response assumption that clients needed `limit`, `offset`, and `count` metadata instead of explicit paginator-friendly page totals

### Notes
- Verified with `make test`

## [2026-03-30] — Session Summary

### Added
- Added handler-level regression coverage proving projection comparable-candidate loading excludes active players and only returns season rows for the retired candidate set

### Changed
- Updated `buildPlayerProjectionPayload` to load comparable-player candidates through a dedicated helper that queries only retired non-target players
- Updated projection candidate stat loading so request-time projection work only fetches `season_stats` rows for the filtered retired candidate IDs
- Updated `AGENTS.md` Current State so the next session resumes with splitting `handlers/players.go` into smaller focused files

### Fixed
- Removed the projection request-path inefficiency where active-player candidate rows were loaded from the database even though the projection service would always discard them before comparable matching

### Notes
- Verified with `make test`

## [2026-03-28] — Session Summary

### Added
- Added a dedicated `GET /api/v1/players/:id/projection` endpoint that returns a `player` header plus a populated projection payload with forecast points, confidence bands, and comparable-player metadata
- Added a new backend projection engine that blends a role-aware age curve, exponentially weighted recent performance, comparable-player matching, and confidence-band fallbacks when comparable futures are sparse
- Added backend integration coverage for active-player projection responses, retired-player ineligibility responses, invalid ids, missing players, and forced database failures
- Added projection-service unit coverage for quality-threshold comparable filtering, no-comparable confidence-band fallback behavior, and retired-player ineligibility handling

### Changed
- Updated the player routes to expose `GET /api/v1/players/:id/projection` under the versioned `/api/v1` group
- Updated `AGENTS.md` Current State so the next session wires the `/career-arc` projection block to the real projection engine after the dedicated endpoint is stable

### Fixed
- Fixed comparable-player selection so the projection engine returns only matches that pass the distance threshold and never pads the response with low-quality comparables just to hit a target count

### Notes
- Verified with `make test`

## [2026-03-28] — Session Summary

### Added
- Added backend integration coverage proving `GET /api/v1/players/:id/career-arc` uses the in-window peak season score when the overall timeline max lives outside `peak_year_start..peak_year_end`
- Added backend integration coverage for the closed-database `500` branch on `GET /api/v1/players/:id/career-arc`

### Changed
- Updated career-arc metadata shaping so `peak_value_score` now prefers the inclusive stored peak window and only falls back to the overall historical max when the peak window has no matching timeline seasons
- Updated `AGENTS.md` Current State so the next session resumes with `GET /api/v1/players/:id/projection`

### Fixed
- Fixed `GET /api/v1/players/:id/career-arc` so peak-score reporting no longer leaks a higher `value_score` from seasons outside the stored peak window

### Notes
- Verified with `make test`

## [2026-03-27] — Session Summary

### Added
- Added a deliberate `GET /api/v1/players/:id/career-arc` response envelope with `player`, optional `arc`, chart-ready `timeline`, and `projection` placeholder sections
- Added backend integration coverage for the career-arc endpoint covering arc-backed responses, players with history but no `career_arcs` row, two-way timeline rows, invalid ids, and missing players
- Added backend integration coverage for `GET /api/v1/players/:id` invalid-id and forced database-failure paths

### Changed
- Updated `GET /api/v1/players/:id/career-arc` to reuse the player-detail season shaping helpers so chart timeline rows return the same `hitting` / `pitching` structure and workload-based nullability
- Updated career-arc metadata shaping to derive `peak_value_score` and `career_value_score_total` from the loaded historical timeline instead of requiring schema changes
- Updated `AGENTS.md` Current State so the next session moves to implementing `GET /api/v1/players/:id/projection`

### Fixed
- Fixed `GET /api/v1/players/:id` so only `gorm.ErrRecordNotFound` returns `404`, while real database/query failures now return `500`
- Fixed `GET /api/v1/players/:id/career-arc` so a missing `career_arcs` row no longer incorrectly returns `404`; existing players now get `200` with `"arc": null`

### Notes
- Verified with `make test`

## [2026-03-27] — Session Summary

### Added
- Added a typed `GET /api/v1/players/:id` response envelope with full `career_stats` history and explicit `hitting` / `pitching` sub-objects per season
- Added backend integration coverage for player detail with season history, no-history players, two-way seasons, and unknown-player `404` responses

### Changed
- Updated the player detail handler to derive `latest_season` from the final `year ASC` season row already loaded for the response instead of issuing a separate lookup
- Updated season-detail shaping to treat workload presence (`plate_appearances`, `innings_pitched`) as the source of truth for whether hitting or pitching data should be returned
- Updated `AGENTS.md` Current State so the next session moves to reshaping `GET /api/v1/players/:id/career-arc`

### Fixed
- Removed the raw GORM preload response from `GET /api/v1/players/:id` so the endpoint no longer leaks contracts and unrelated model structure into the player-detail API surface

### Notes
- Verified with `make test`

## [2026-03-26] — Session Summary

### Added
- Added frontend TypeScript types for the player-list response envelope, including list rows, joined latest-season snapshot data, and pagination metadata

### Changed
- Updated the frontend player-search helper to call `GET /api/v1/players?q=...` and unwrap the new `{ data, meta }` response shape
- Updated the frontend list helper to unwrap the same typed list-response envelope instead of assuming a raw array
- Updated `AGENTS.md` Current State so the removed legacy search route and the successful frontend build are recorded for the next session

### Fixed
- Removed the legacy `GET /api/v1/players/search` backend route and handler now that the frontend no longer depends on them

### Notes
- Verified with `make test`
- Verified frontend build with `docker compose exec frontend npm run build`

## [2026-03-26] — Session Summary

### Added
- Added backend integration coverage for `GET /api/v1/players?team=...`, including both joined `team_name` substring matching and numeric `team_id` matching
- Added backend integration coverage for `GET /api/v1/players?season=<year>` so the fixed-year season snapshot branch is explicitly exercised

### Changed
- Refactored handler test fixtures to support creating multiple `season_stats` rows for one player when list-endpoint snapshot behavior needs to be verified
- Updated `AGENTS.md` Current State so the completed test coverage and the remaining legacy `/players/search` decision are both recorded for the next session

### Fixed
- Locked in regression coverage for the player-list snapshot contract so future query changes are less likely to break `team` filtering or year-scoped snapshot shaping silently

### Notes
- Verified with `make test`

## [2026-03-26] — Session Summary

### Added
- Added a compact `GET /api/v1/players` response envelope with `data` and `meta` sections for list consumers
- Added backend endpoint coverage for `q`, `active`, and `sort=-value_score` on the real `/api/v1/players` route

### Changed
- Updated `GET /api/v1/players` to support `q`, `active`, `position`, `team`, `season`, `limit`, `offset`, and `sort`
- Updated the player list query to join a derived season snapshot so list rows can include latest-season context without preloading full season histories
- Updated `AGENTS.md` Current State so the next session moves to reshaping `GET /api/players/:id`

### Fixed
- Fixed the player list endpoint so joined-field sorting and team filtering operate on a one-row-per-player derived season snapshot instead of raw `season_stats` joins
- Fixed the list response for players without season data so `latest_season` is `null` instead of an empty object

### Notes
- Verified with `make test`

## [2026-03-26] — Session Summary

### Added
- Added frontend build verification for the admin dashboard API path fix via `docker compose exec frontend npm run build`

### Changed
- Updated the frontend API client to normalize `VITE_API_URL` to the versioned `/api/v1` base before issuing requests
- Updated `AGENTS.md` Current State so the next session still resumes with Phase 2 player endpoint work after this frontend fix

### Fixed
- Fixed the admin dashboard frontend request so it targets `GET /api/v1/admin/dashboard` instead of falling back to the stale unversioned `/api` base

### Notes
- Verified with `make test`
- Verified frontend build with `docker compose exec frontend npm run build`

## [2026-03-26] — Session Summary

### Added
- Added backend regression coverage for the admin dashboard `OPTIONS` preflight so `/api/v1/admin/dashboard` now has explicit CORS-header test coverage

### Changed
- Updated backend route registration to attach the shared CORS middleware to the `/api/v1` Echo group that owns the admin dashboard endpoint
- Updated `AGENTS.md` Current State so the next session still resumes with Phase 2 player endpoint work after this regression fix

### Fixed
- Fixed the admin dashboard CORS regression where the browser preflight returned `204` without any `Access-Control-Allow-*` headers

### Notes
- Verified with `make test`

## [2026-03-26] — Session Summary

### Added
- Added backend regression tests covering admin bearer-token validation and unauthorized dashboard requests
- Added root `.env` entries for `ADMIN_SECRET` and `VITE_ADMIN_SECRET` to support the new local admin dashboard access flow

### Changed
- Updated the backend admin dashboard endpoint to require `Authorization: Bearer <ADMIN_SECRET>`
- Updated the frontend admin dashboard request to send `VITE_ADMIN_SECRET` only on the admin dashboard API call
- Pointed Vite at the repo-root `.env` so frontend env vars load from the same local file as backend env vars
- Updated `.env.example` to document the new backend/frontend admin secret variables
- Updated `AGENTS.md` Current State so the next session resumes with Phase 2 player endpoint work

### Fixed
- Fixed the admin dashboard path/auth mismatch so the page can call the versioned backend admin endpoint with the expected bearer token
- Removed the Admin link from the public navbar while keeping the `/admin` route directly accessible

### Notes
- Verified with `make test`
- Frontend build was not re-run in this shell because `npm` was unavailable locally

## [2026-03-25] — Session Summary

### Added
- Added regression coverage proving a canceled scheduler run now clears the in-memory `running` flag, records completion, and exposes the friendly `sync interrupted` admin status
- Added status-store coverage for friendly cancellation and timeout messages

### Changed
- Updated scheduler run lifecycle handling to complete through one deferred path so early exits no longer skip status cleanup
- Updated admin-facing scheduler error text to show `sync interrupted` for `context.Canceled` and `sync timed out` for `context.DeadlineExceeded`
- Updated `AGENTS.md` Current State so the next session starts with Phase 2 `GET /api/players` work

### Fixed
- Fixed the false `running` dashboard state that could persist after shutdown or context cancellation interrupted a scheduled sync

### Notes
- Verified with `make test`
- Suggested commit message: `fix(syncjob): clear status on canceled runs`

## [2026-03-23] — Session Summary

### Added
- Added review-driven documentation follow-ups to `AGENTS.md` so the scheduler cancellation bug is now the top-priority next task

### Changed
- Updated `AGENTS.md` Current State ordering to put the in-memory scheduler cancellation bug ahead of Phase 2 API work
- Reconciled `app-structure.md` with `AGENTS.md` by marking Phase 1 complete in both places
- Recorded the manually verified frontend build command in `AGENTS.md`

### Fixed
- Corrected the project-state docs so the repo no longer gives conflicting guidance about whether Phase 1 is complete

### Notes
- Docs-only session; no code changes were made
- Frontend build was manually verified passing via `docker compose exec frontend npm run build`

## [2026-03-23] — Session Summary

### Added
- Added a dedicated scheduled-refresh ingestion path that refreshes only one relevant season per active player instead of replaying full history on each scheduled run
- Added an in-memory scheduler status store plus a minimal read-only admin dashboard endpoint and frontend page at `/admin`
- Added regression coverage for target-season selection, season-scoped scheduled refresh behavior, and scheduler status lifecycle tracking

### Changed
- Updated the background sync job to target the active season during April through October and the most recently completed season during the off-season
- Updated backend startup and routing so the dashboard can read live scheduler status from memory without introducing persistence or auth
- Updated `AGENTS.md` and `app-structure.md` to mark Phase 1 complete and point the next session at Phase 2 API work

### Fixed
- Stopped scheduled syncs from re-running full year-by-year ingestion history for active players on every pass

### Notes
- Verified backend tests with `cd backend && GOCACHE=/tmp/capcurve-gocache go test -mod=mod ./...`
- Could not run frontend type-check/build in this shell because `npm` is unavailable

## [2026-03-19] — Session Summary

### Added
- Added a new `internal/syncjob` package that computes season-aware sync timing and runs background active-player re-ingestion through the existing ingestion service
- Added unit tests covering in-season vs. off-season scheduling boundaries and the active-player-only sync execution path

### Changed
- Updated backend startup to launch the scheduler automatically with graceful shutdown handling
- Extended backend config with sync env settings for enablement, timezone, run time, and off-season weekday
- Updated `AGENTS.md` and `app-structure.md` to reflect the new daily in-season / weekly off-season cadence and to move the admin dashboard to the front of the queue

### Fixed
- Ensured scheduled syncs only target players whose `active` flag is true, avoiding unnecessary re-ingestion for retired players

### Notes
- Default scheduler behavior is `05:00` in `America/New_York`, daily from April through October and weekly on Monday from November through March
- Manually verified tests passing via `cd backend && GOCACHE=/tmp/capcurve-gocache go test -mod=mod ./...`

## [2026-03-16] — Session Summary

### Added
- Added documented source-research findings for Cot's Opening Day Salaries and Lahman `Salaries.csv` to `AGENTS.md` so the contract-data source decision can be made from verified inputs instead of assumptions

### Changed
- Updated `AGENTS.md` Current State to record the confirmed Cot's workbook URL, the verified `2000-2025` tab coverage, the observed Cot's schema drift, and the official Lahman salary-table field list

### Fixed
- Corrected the project's contract-source decision context by verifying that Cot's is a name-first sheet with no stable player IDs or team field on salary rows, and that Lahman `Salaries.csv` does not contain columns beyond `yearID`, `teamID`, `lgID`, `playerID`, and `salary`

### Notes
- Research-only session; no importer code, schema changes, or tests were run
- Verified that the current SABR-linked Lahman documentation still works, but the linked public CSV archive URL resolves to a removed Box share page as of March 16, 2026

## [2026-03-16] — Session Summary

### Added
- Added regression coverage proving ingestion prefers MLB display-name fields for players whose legal first name differs from the public baseball name

### Changed
- Updated MLB player normalization to prefer `useName` / `useLastName` over `firstName` / `lastName` when the Stats API provides both
- Re-ran the ingestion CLI for Albert Pujols (`405395`) so the local `players` row now reflects the corrected canonical name
- Updated `AGENTS.md` Current State so the next task is the remaining broader Phase 1 scoring/data cleanup

### Fixed
- Corrected the persisted local player record for `mlb_id = 405395` from `Jose Pujols` to `Albert Pujols`

### Notes
- Verified in Postgres that `players.mlb_id = 405395` now stores `first_name = Albert` and `last_name = Pujols`
- Ran `make test` successfully

## [2026-03-16] — Session Summary

### Added
- Added a Current State follow-up to fix the local player-name mismatch for `mlb_id = 405395` during the next ingestion/data cleanup pass

### Changed
- Updated `AGENTS.md` notes to clarify that `405395` is Albert Pujols' correct MLB ID and that the mismatch is in local persisted data, not project documentation

### Fixed
- Corrected the review follow-up scope so the next session targets the database/ingestion name discrepancy instead of changing `AGENTS.md` player-ID documentation

### Notes
- No code changes; this was a project-state correction after review

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
## [2026-03-16] — Session Summary

### Added
- Added a Database Schema section to `README.md`

### Changed
- Documented in `README.md` that active `season_stats` uniqueness is enforced in PostgreSQL with a partial unique index on `(player_id, year) WHERE deleted_at IS NULL`
- Updated `AGENTS.md` Current State so the Albert Pujols `mlb_id = 405395` player-name mismatch is now the next Phase 1 cleanup task

### Fixed
- Closed the outstanding repo-doc follow-up for the database-enforced one-active-row-per-player-season rule

### Notes
- No tests were run because this session only changed repository documentation/state tracking

## [2026-03-16] — Session Summary

### Added
- Added contract-data exploration findings to `AGENTS.md` after probing MLB Stats API endpoints for player finances

### Changed
- Documented that MLB Stats API exposes contract-adjacent transaction events such as signings and free agency, but not salary, total value, AAV, guarantees, options, or contract years
- Updated `AGENTS.md` Current State so the next Phase 1 task is choosing a non-Stats-API source and ingestion strategy for financial contract data before building the contract importer

### Fixed
- Closed the open uncertainty about whether MLB Stats API alone can power CapCurve's contract importer strategy

### Notes
- Tested directly against Shohei Ohtani (`660271`), Aaron Judge (`592450`), and Albert Pujols (`405395`) using `people`, `stats`, `transactions`, and `hydrate=transactions` endpoints
- Albert Pujols transaction coverage only reached back to 2009 despite his MLB debut on April 2, 2001, so older transaction history should not be treated as complete
- No tests were run because this session only changed repository documentation/state tracking

## [2026-03-30] — Session Summary

### Added
- Added projection-service regression coverage proving active-player candidates are excluded from comparable matching even when their trajectory would otherwise qualify

### Changed
- Restricted projection comparable-player selection to historical or retired players by filtering on `players.active = false` before distance scoring
- Updated `AGENTS.md` Current State so wiring the real projection engine into `GET /api/v1/players/:id/career-arc` is now the next Phase 2 task

### Fixed
- Prevented active-player futures from influencing projection confidence bands that are supposed to be anchored by completed career outcomes

### Notes
- The projection engine still requires comparable candidates to have a valid anchor season plus at least one later season, so retired-only filtering works alongside the existing future-outcome guardrails
- Ran `GOCACHE=/tmp/capcurve-gocache go test ./internal/projection` and `make test` successfully

## [2026-03-30] — Session Summary

### Added
- Added backend integration coverage proving `GET /api/v1/players/:id/career-arc` returns the same projection payload as `GET /api/v1/players/:id/projection` for the same player

### Changed
- Wired the career-arc endpoint to the real shared projection engine instead of the temporary placeholder payload
- Consolidated projection loading into a shared handler helper so career-arc and dedicated projection responses stay in sync

### Fixed
- Removed the stale `"projection engine not implemented yet"` placeholder behavior from the career-arc response for active players

### Notes
- `/career-arc` now returns the full shared projection payload shape because the player page needs forecast points, confidence bands, and comparable-player context together
- Ran `GOCACHE=/tmp/capcurve-gocache go test ./internal/handlers` and `make test` successfully
