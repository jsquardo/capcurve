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

export interface SeasonStat {
  id: number
  player_id: number
  year: number
  team_name: string
  games_played: number
  at_bats: number
  hits: number
  home_runs: number
  rbi: number
  stolen_bases: number
  batting_avg: number
  obp: number
  slg: number
  ops: number
  ops_plus: number
  wrc_plus: number
  wins: number
  losses: number
  era: number
  era_plus: number
  whip: number
  strikeouts: number
  innings_pitched: number
  fip: number
  war: number
  value_score: number
  war_per_dollar: number
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
