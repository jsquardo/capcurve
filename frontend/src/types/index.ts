export interface Player {
  id: number
  mlb_id: number
  first_name: string
  last_name: string
  position: string
  bats: string
  throws: string
  date_of_birth: string | null
  active: boolean
  image_url: string
  season_stats?: SeasonStat[]
  contracts?: Contract[]
  career_arc?: CareerArc
}

export interface PlayerListSeason {
  year: number
  team_id: number
  team_name: string
  age: number
  value_score: number
}

export interface PlayerListItem {
  id: number
  mlb_id: number
  first_name: string
  last_name: string
  full_name: string
  position: string
  bats: string
  throws: string
  date_of_birth: string | null
  active: boolean
  image_url: string
  latest_season: PlayerListSeason | null
}

export interface PlayerListMeta {
  limit: number
  offset: number
  count: number
  total: number
}

export interface PlayerListResponse {
  data: PlayerListItem[]
  meta: PlayerListMeta
}

export interface SeasonStat {
  id: number
  player_id: number
  year: number
  age: number
  team_id: number
  team_name: string
  games_played: number
  games_started: number
  plate_appearances: number
  at_bats: number
  hits: number
  doubles: number
  triples: number
  home_runs: number
  runs: number
  rbi: number
  walks: number
  strikeouts: number
  stolen_bases: number
  batting_avg: number
  obp: number
  slg: number
  ops: number
  babip: number
  wins: number
  losses: number
  era: number
  whip: number
  innings_pitched: number
  hits_allowed: number
  walks_allowed: number
  home_runs_allowed: number
  strikeouts_per_9: number
  walks_per_9: number
  hits_per_9: number
  home_runs_per_9: number
  strikeout_walk_ratio: number
  strike_percentage: number
  expected_batting_avg: number | null
  expected_slugging: number | null
  expected_woba: number | null
  expected_era: number | null
  barrel_pct: number | null
  hard_hit_pct: number | null
  avg_exit_velocity: number | null
  avg_launch_angle: number | null
  sweet_spot_pct: number | null
  value_score: number
}

export interface Contract {
  id: number
  player_id: number
  team_name: string
  total_value: number
  aav: number
  years: number
  start_year: number
  end_year: number
  signing_age: number
  contract_type: string
  overall_value_score: number
  is_active: boolean
  contract_seasons?: ContractSeason[]
}

export interface ContractSeason {
  id: number
  contract_id: number
  player_id: number
  year: number
  salary: number
  war: number
  war_value: number
  value_score: number
  verdict_label: 'Bargain' | 'Fair' | 'Overpaid' | 'Albatross'
}

export interface CareerArc {
  id: number
  player_id: number
  peak_year_start: number
  peak_year_end: number
  peak_war: number
  career_war: number
  decline_onset_year: number
  arc_shape: 'early_peak' | 'late_bloomer' | 'sustained' | 'flash' | 'declining'
  last_computed_at: string
}

export interface CareerArcResponse {
  arc: CareerArc
  season_stats: SeasonStat[]
}

export interface SyncStatus {
  enabled: boolean
  running: boolean
  in_season: boolean
  target_season_year: number
  last_sync_started_at: string | null
  last_sync_completed_at: string | null
  last_successful_sync_at: string | null
  next_scheduled_sync_at: string | null
  last_error: string
}

export interface AdminDashboard {
  total_players: number
  active_players: number
  sync_status: SyncStatus
}
