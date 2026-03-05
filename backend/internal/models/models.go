package models

import (
	"time"

	"gorm.io/gorm"
)

type Player struct {
	gorm.Model
	MLBID       int        `gorm:"uniqueIndex;not null" json:"mlb_id"`
	FirstName   string     `gorm:"not null" json:"first_name"`
	LastName    string     `gorm:"not null" json:"last_name"`
	Position    string     `gorm:"not null" json:"position"`
	Bats        string     `json:"bats"`
	Throws      string     `json:"throws"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Active      bool       `gorm:"default:true" json:"active"`
	ImageURL    string     `json:"image_url"`

	SeasonStats []SeasonStat `gorm:"foreignKey:PlayerID" json:"season_stats,omitempty"`
	Contracts   []Contract   `gorm:"foreignKey:PlayerID" json:"contracts,omitempty"`
	CareerArc   *CareerArc   `gorm:"foreignKey:PlayerID" json:"career_arc,omitempty"`
}

type SeasonStat struct {
	gorm.Model
	PlayerID int    `gorm:"index;not null" json:"player_id"`
	Year     int    `gorm:"not null" json:"year"`
	TeamID   int    `json:"team_id"`
	TeamName string `json:"team_name"`

	GamesPlayed    int     `json:"games_played"`
	AtBats         int     `json:"at_bats"`
	Hits           int     `json:"hits"`
	HomeRuns       int     `json:"home_runs"`
	RBI            int     `json:"rbi"`
	StolenBases    int     `json:"stolen_bases"`
	BattingAvg     float64 `json:"batting_avg"`
	OBP            float64 `json:"obp"`
	SLG            float64 `json:"slg"`
	OPS            float64 `json:"ops"`
	OPSPlus        int     `json:"ops_plus"`
	WRCPlus        int     `json:"wrc_plus"`

	Wins           int     `json:"wins"`
	Losses         int     `json:"losses"`
	ERA            float64 `json:"era"`
	ERAPlus        int     `json:"era_plus"`
	WHIP           float64 `json:"whip"`
	Strikeouts     int     `json:"strikeouts"`
	InningsPitched float64 `json:"innings_pitched"`
	FIP            float64 `json:"fip"`
	WAR            float64 `json:"war"`

	ValueScore   float64 `json:"value_score"`
	WARPerDollar float64 `json:"war_per_dollar"`
}

type Contract struct {
	gorm.Model
	PlayerID     int     `gorm:"index;not null" json:"player_id"`
	TeamID       int     `json:"team_id"`
	TeamName     string  `json:"team_name"`
	TotalValue   float64 `json:"total_value"`
	AAV          float64 `json:"aav"`
	Years        int     `json:"years"`
	StartYear    int     `json:"start_year"`
	EndYear      int     `json:"end_year"`
	SigningAge   int     `json:"signing_age"`
	ContractType string  `json:"contract_type"`

	OverallValueScore float64 `json:"overall_value_score"`
	IsActive          bool    `gorm:"default:false" json:"is_active"`

	ContractSeasons []ContractSeason `gorm:"foreignKey:ContractID" json:"contract_seasons,omitempty"`
}

type ContractSeason struct {
	gorm.Model
	ContractID   int     `gorm:"index;not null" json:"contract_id"`
	PlayerID     int     `gorm:"index;not null" json:"player_id"`
	Year         int     `json:"year"`
	Salary       float64 `json:"salary"`
	WAR          float64 `json:"war"`
	WARValue     float64 `json:"war_value"`
	ValueScore   float64 `json:"value_score"`
	VerdictLabel string  `json:"verdict_label"`
}

type CareerArc struct {
	gorm.Model
	PlayerID         int       `gorm:"uniqueIndex;not null" json:"player_id"`
	PeakYearStart    int       `json:"peak_year_start"`
	PeakYearEnd      int       `json:"peak_year_end"`
	PeakWAR          float64   `json:"peak_war"`
	CareerWAR        float64   `json:"career_war"`
	DeclineOnsetYear int       `json:"decline_onset_year"`
	ArcShape         string    `json:"arc_shape"`
	LastComputedAt   time.Time `json:"last_computed_at"`
}
