// ── Phase 5 — Contract Value Tracker (DEFERRED) ──────────────────────────────
//
// These types are explicitly deferred until a clean salary/contract data source
// for 2017–present is identified. Do NOT import or use these types anywhere in
// the active codebase until Phase 5 begins and the AGENTS.md note is removed.
//
// See the Phase 5 section of AGENTS.md for the full list of evaluated sources
// and the reasons they were ruled out.

export interface Contract {
  id: number;
  player_id: number;
  team_name: string;
  total_value: number;
  aav: number;
  years: number;
  start_year: number;
  end_year: number;
  signing_age: number;
  contract_type: string;
  overall_value_score: number;
  is_active: boolean;
  contract_seasons?: ContractSeason[];
}

export interface ContractSeason {
  id: number;
  contract_id: number;
  player_id: number;
  year: number;
  salary: number;
  war: number;
  war_value: number;
  value_score: number;
  verdict_label: "Bargain" | "Fair" | "Overpaid" | "Albatross";
}
