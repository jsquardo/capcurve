export interface Player {
  id: number;
  mlb_id: number;
  first_name: string;
  last_name: string;
  position: string;
  bats: string;
  throws: string;
  date_of_birth: string | null;
  active: boolean;
  image_url: string;
}

export interface PlayerListSeason {
  year: number;
  team_id: number;
  team_name: string;
  age: number;
  value_score: number;
}

export interface PlayerListItem {
  id: number;
  mlb_id: number;
  first_name: string;
  last_name: string;
  full_name: string;
  position: string;
  bats: string;
  throws: string;
  date_of_birth: string | null;
  active: boolean;
  image_url: string;
  latest_season: PlayerListSeason | null;
}

export interface PlayerListMeta {
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface PlayerListResponse {
  data: PlayerListItem[];
  meta: PlayerListMeta;
}

export interface SeasonStat {
  id: number;
  player_id: number;
  year: number;
  age: number;
  team_id: number;
  team_name: string;
  games_played: number;
  games_started: number;
  plate_appearances: number;
  at_bats: number;
  hits: number;
  doubles: number;
  triples: number;
  home_runs: number;
  runs: number;
  rbi: number;
  walks: number;
  strikeouts: number;
  stolen_bases: number;
  batting_avg: number;
  obp: number;
  slg: number;
  ops: number;
  babip: number;
  wins: number;
  losses: number;
  era: number;
  whip: number;
  innings_pitched: number;
  hits_allowed: number;
  walks_allowed: number;
  home_runs_allowed: number;
  strikeouts_per_9: number;
  walks_per_9: number;
  hits_per_9: number;
  home_runs_per_9: number;
  strikeout_walk_ratio: number;
  strike_percentage: number;
  expected_batting_avg: number | null;
  expected_slugging: number | null;
  expected_woba: number | null;
  expected_era: number | null;
  barrel_pct: number | null;
  hard_hit_pct: number | null;
  avg_exit_velocity: number | null;
  avg_launch_angle: number | null;
  sweet_spot_pct: number | null;
  value_score: number;
}

// ── Player Detail (GET /api/v1/players/:id) ───────────────────────────────────

export interface HittingStats {
  games_played: number;
  plate_appearances: number;
  at_bats: number;
  hits: number;
  doubles: number;
  triples: number;
  home_runs: number;
  runs: number;
  rbi: number;
  walks: number;
  strikeouts: number;
  stolen_bases: number;
  batting_avg: number;
  obp: number;
  slg: number;
  ops: number;
  babip: number;
  expected_batting_avg: number | null;
  expected_slugging: number | null;
  expected_woba: number | null;
  barrel_pct: number | null;
  hard_hit_pct: number | null;
  avg_exit_velocity: number | null;
  avg_launch_angle: number | null;
  sweet_spot_pct: number | null;
}

export interface PitchingStats {
  games_played: number;
  games_started: number;
  wins: number;
  losses: number;
  era: number;
  whip: number;
  innings_pitched: number;
  hits_allowed: number;
  walks_allowed: number;
  home_runs_allowed: number;
  strikeouts_per_9: number;
  walks_per_9: number;
  hits_per_9: number;
  home_runs_per_9: number;
  strikeout_walk_ratio: number;
  strike_percentage: number;
  expected_era: number | null;
}

export interface CareerStatItem {
  year: number;
  team_id: number;
  team_name: string;
  age: number;
  value_score: number;
  hitting: HittingStats | null;
  pitching: PitchingStats | null;
}

// PlayerDetail is what GET /api/v1/players/:id returns (inside the { data: ... } envelope)
export interface PlayerDetail {
  id: number;
  mlb_id: number;
  first_name: string;
  last_name: string;
  full_name: string;
  position: string;
  bats: string;
  throws: string;
  date_of_birth: string | null;
  active: boolean;
  image_url: string;
  latest_season: PlayerListSeason | null;
  career_stats: CareerStatItem[];
}

// ── Career Arc (GET /api/v1/players/:id/career-arc) ──────────────────────────

// CareerArcPlayer is the player fields embedded inside the arc response envelope.
export interface CareerArcPlayer {
  id: number;
  mlb_id: number;
  first_name: string;
  last_name: string;
  full_name: string;
  position: string;
  bats: string;
  throws: string;
  date_of_birth: string | null;
  active: boolean;
  image_url: string;
}

// CareerArcMeta is the arc summary row. null when no career_arcs row exists yet.
export interface CareerArcMeta {
  peak_year_start: number;
  peak_year_end: number;
  decline_onset_year: number;
  arc_shape:
    | "early_peak"
    | "late_bloomer"
    | "sustained"
    | "flash"
    | "declining";
  peak_value_score: number;
  career_value_score_total: number;
  last_computed_at: string;
}

export interface CareerArcTimelineItem {
  year: number;
  team_id: number;
  team_name: string;
  age: number;
  value_score: number;
  is_peak: boolean;
  is_projection: boolean;
  hitting: HittingStats | null;
  pitching: PitchingStats | null;
}

export interface ProjectionPoint {
  year: number;
  age: number;
  value_score: number;
}

export interface ConfidenceBandPoint {
  year: number;
  lower: number;
  upper: number;
}

export interface ComparablePlayer {
  player_id: number;
  mlb_id: number;
  full_name: string;
  position: string;
}

export interface CareerArcProjection {
  status: "ready" | "insufficient_data" | "ineligible";
  eligible: boolean;
  reason: string;
  points: ProjectionPoint[];
  confidence_band: ConfidenceBandPoint[];
  comparables: ComparablePlayer[];
}

export interface CareerArcData {
  player: CareerArcPlayer;
  arc: CareerArcMeta | null;
  timeline: CareerArcTimelineItem[];
  projection: CareerArcProjection;
}

// CareerArcResponse is the full HTTP response envelope from GET /api/v1/players/:id/career-arc
export interface CareerArcResponse {
  data: CareerArcData;
}

// ── Admin ─────────────────────────────────────────────────────────────────────

export interface SyncStatus {
  enabled: boolean;
  running: boolean;
  in_season: boolean;
  target_season_year: number;
  last_sync_started_at: string | null;
  last_sync_completed_at: string | null;
  last_successful_sync_at: string | null;
  next_scheduled_sync_at: string | null;
  last_error: string;
}

export interface AdminDashboard {
  total_players: number;
  active_players: number;
  sync_status: SyncStatus;
}

// ── Leaderboards ─────────────────────────────────────────────────────────────

export type LeaderboardCategory = "peak_arc" | "hr" | "avg" | "era" | "k9";

export interface LeaderboardEntry {
  rank: number;
  player_id: number;
  player_name: string;
  position: string;
  team: string;
  value: number;
  season: number | null; // null for peak_arc (career metric, not season-scoped)
}

export interface LeaderboardMeta {
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface LeaderboardsResponse {
  data: {
    category: LeaderboardCategory;
    leaders: LeaderboardEntry[];
    meta: LeaderboardMeta;
  };
}
