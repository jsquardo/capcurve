CREATE TABLE players (
    id SERIAL PRIMARY KEY,
    mlb_id INTEGER NOT NULL UNIQUE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    position VARCHAR(10) NOT NULL,
    bats VARCHAR(1),
    throws VARCHAR(1),
    date_of_birth DATE,
    active BOOLEAN DEFAULT TRUE,
    image_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_players_active ON players(active);
CREATE INDEX idx_players_position ON players(position);
CREATE INDEX idx_players_last_name ON players(last_name);

CREATE TABLE season_stats (
    id SERIAL PRIMARY KEY,
    player_id INTEGER NOT NULL REFERENCES players(id),
    year INTEGER NOT NULL,
    team_id INTEGER,
    team_name VARCHAR(100),
    games_played INTEGER DEFAULT 0,
    at_bats INTEGER DEFAULT 0,
    hits INTEGER DEFAULT 0,
    home_runs INTEGER DEFAULT 0,
    rbi INTEGER DEFAULT 0,
    stolen_bases INTEGER DEFAULT 0,
    batting_avg DECIMAL(5,3) DEFAULT 0,
    obp DECIMAL(5,3) DEFAULT 0,
    slg DECIMAL(5,3) DEFAULT 0,
    ops DECIMAL(5,3) DEFAULT 0,
    ops_plus INTEGER DEFAULT 0,
    wrc_plus INTEGER DEFAULT 0,
    wins INTEGER DEFAULT 0,
    losses INTEGER DEFAULT 0,
    era DECIMAL(5,2) DEFAULT 0,
    era_plus INTEGER DEFAULT 0,
    whip DECIMAL(5,3) DEFAULT 0,
    strikeouts INTEGER DEFAULT 0,
    innings_pitched DECIMAL(6,1) DEFAULT 0,
    fip DECIMAL(5,2) DEFAULT 0,
    war DECIMAL(5,1) DEFAULT 0,
    value_score DECIMAL(8,2) DEFAULT 0,
    war_per_dollar DECIMAL(8,4) DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE(player_id, year, team_id)
);

CREATE INDEX idx_season_stats_player_id ON season_stats(player_id);
CREATE INDEX idx_season_stats_year ON season_stats(year);
CREATE INDEX idx_season_stats_war ON season_stats(war DESC);

CREATE TABLE contracts (
    id SERIAL PRIMARY KEY,
    player_id INTEGER NOT NULL REFERENCES players(id),
    team_id INTEGER,
    team_name VARCHAR(100),
    total_value DECIMAL(10,2) NOT NULL,
    aav DECIMAL(10,2) NOT NULL,
    years INTEGER NOT NULL,
    start_year INTEGER NOT NULL,
    end_year INTEGER NOT NULL,
    signing_age INTEGER,
    contract_type VARCHAR(20) DEFAULT 'free_agent',
    overall_value_score DECIMAL(8,2) DEFAULT 0,
    is_active BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_contracts_player_id ON contracts(player_id);
CREATE INDEX idx_contracts_is_active ON contracts(is_active);
CREATE INDEX idx_contracts_value_score ON contracts(overall_value_score);

CREATE TABLE contract_seasons (
    id SERIAL PRIMARY KEY,
    contract_id INTEGER NOT NULL REFERENCES contracts(id),
    player_id INTEGER NOT NULL REFERENCES players(id),
    year INTEGER NOT NULL,
    salary DECIMAL(10,2) NOT NULL,
    war DECIMAL(5,1) DEFAULT 0,
    war_value DECIMAL(10,2) DEFAULT 0,
    value_score DECIMAL(8,2) DEFAULT 0,
    verdict_label VARCHAR(20) DEFAULT 'Fair',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_contract_seasons_contract_id ON contract_seasons(contract_id);
CREATE INDEX idx_contract_seasons_player_id ON contract_seasons(player_id);

CREATE TABLE career_arcs (
    id SERIAL PRIMARY KEY,
    player_id INTEGER NOT NULL UNIQUE REFERENCES players(id),
    peak_year_start INTEGER,
    peak_year_end INTEGER,
    peak_war DECIMAL(5,1) DEFAULT 0,
    career_war DECIMAL(6,1) DEFAULT 0,
    decline_onset_year INTEGER,
    arc_shape VARCHAR(30) DEFAULT 'sustained',
    last_computed_at TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_career_arcs_player_id ON career_arcs(player_id);
CREATE INDEX idx_career_arcs_peak_war ON career_arcs(peak_war DESC);
