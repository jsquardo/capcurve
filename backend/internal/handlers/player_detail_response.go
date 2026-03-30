package handlers

import (
	"strings"
	"time"

	"github.com/jsquardo/capcurve/internal/models"
)

type playerDetailResponse struct {
	Data playerDetailItem `json:"data"`
}

type playerDetailItem struct {
	ID           uint                   `json:"id"`
	MLBID        int                    `json:"mlb_id"`
	FirstName    string                 `json:"first_name"`
	LastName     string                 `json:"last_name"`
	FullName     string                 `json:"full_name"`
	Position     string                 `json:"position"`
	Bats         string                 `json:"bats"`
	Throws       string                 `json:"throws"`
	DateOfBirth  *time.Time             `json:"date_of_birth"`
	Active       bool                   `json:"active"`
	ImageURL     string                 `json:"image_url"`
	LatestSeason *playerSeasonListItem  `json:"latest_season"`
	CareerStats  []playerCareerStatItem `json:"career_stats"`
}

type playerCareerStatItem struct {
	Year       int                  `json:"year"`
	TeamID     int                  `json:"team_id"`
	TeamName   string               `json:"team_name"`
	Age        int                  `json:"age"`
	ValueScore float64              `json:"value_score"`
	Hitting    *playerHittingStats  `json:"hitting"`
	Pitching   *playerPitchingStats `json:"pitching"`
}

type playerHittingStats struct {
	GamesPlayed        int      `json:"games_played"`
	PlateAppearances   int      `json:"plate_appearances"`
	AtBats             int      `json:"at_bats"`
	Hits               int      `json:"hits"`
	Doubles            int      `json:"doubles"`
	Triples            int      `json:"triples"`
	HomeRuns           int      `json:"home_runs"`
	Runs               int      `json:"runs"`
	RBI                int      `json:"rbi"`
	Walks              int      `json:"walks"`
	Strikeouts         int      `json:"strikeouts"`
	StolenBases        int      `json:"stolen_bases"`
	BattingAvg         float64  `json:"batting_avg"`
	OBP                float64  `json:"obp"`
	SLG                float64  `json:"slg"`
	OPS                float64  `json:"ops"`
	BABIP              float64  `json:"babip"`
	ExpectedBattingAvg *float64 `json:"expected_batting_avg"`
	ExpectedSlugging   *float64 `json:"expected_slugging"`
	ExpectedWOBA       *float64 `json:"expected_woba"`
	BarrelPct          *float64 `json:"barrel_pct"`
	HardHitPct         *float64 `json:"hard_hit_pct"`
	AvgExitVelocity    *float64 `json:"avg_exit_velocity"`
	AvgLaunchAngle     *float64 `json:"avg_launch_angle"`
	SweetSpotPct       *float64 `json:"sweet_spot_pct"`
}

type playerPitchingStats struct {
	GamesPlayed        int      `json:"games_played"`
	GamesStarted       int      `json:"games_started"`
	Wins               int      `json:"wins"`
	Losses             int      `json:"losses"`
	ERA                float64  `json:"era"`
	WHIP               float64  `json:"whip"`
	InningsPitched     float64  `json:"innings_pitched"`
	HitsAllowed        int      `json:"hits_allowed"`
	WalksAllowed       int      `json:"walks_allowed"`
	HomeRunsAllowed    int      `json:"home_runs_allowed"`
	StrikeoutsPer9     float64  `json:"strikeouts_per_9"`
	WalksPer9          float64  `json:"walks_per_9"`
	HitsPer9           float64  `json:"hits_per_9"`
	HomeRunsPer9       float64  `json:"home_runs_per_9"`
	StrikeoutWalkRatio float64  `json:"strikeout_walk_ratio"`
	StrikePercentage   float64  `json:"strike_percentage"`
	ExpectedERA        *float64 `json:"expected_era"`
}

func newPlayerDetailItem(player models.Player, seasonStats []models.SeasonStat) playerDetailItem {
	item := playerDetailItem{
		ID:          player.ID,
		MLBID:       player.MLBID,
		FirstName:   player.FirstName,
		LastName:    player.LastName,
		FullName:    strings.TrimSpace(player.FirstName + " " + player.LastName),
		Position:    player.Position,
		Bats:        player.Bats,
		Throws:      player.Throws,
		DateOfBirth: player.DateOfBirth,
		Active:      player.Active,
		ImageURL:    player.ImageURL,
		CareerStats: make([]playerCareerStatItem, 0, len(seasonStats)),
	}

	for _, stat := range seasonStats {
		item.CareerStats = append(item.CareerStats, newPlayerCareerStatItem(stat))
	}

	// The detail endpoint loads season rows in year-ascending order, so the final
	// element is the latest available season snapshot for the profile header.
	if len(item.CareerStats) > 0 {
		latest := item.CareerStats[len(item.CareerStats)-1]
		item.LatestSeason = &playerSeasonListItem{
			Year:       latest.Year,
			TeamID:     latest.TeamID,
			TeamName:   latest.TeamName,
			Age:        latest.Age,
			ValueScore: latest.ValueScore,
		}
	}

	return item
}

func newPlayerCareerStatItem(stat models.SeasonStat) playerCareerStatItem {
	return playerCareerStatItem{
		Year:       stat.Year,
		TeamID:     stat.TeamID,
		TeamName:   stat.TeamName,
		Age:        stat.Age,
		ValueScore: stat.ValueScore,
		Hitting:    newPlayerHittingStats(stat),
		Pitching:   newPlayerPitchingStats(stat),
	}
}

func newPlayerHittingStats(stat models.SeasonStat) *playerHittingStats {
	if stat.PlateAppearances <= 0 {
		return nil
	}

	return &playerHittingStats{
		GamesPlayed:        stat.GamesPlayed,
		PlateAppearances:   stat.PlateAppearances,
		AtBats:             stat.AtBats,
		Hits:               stat.Hits,
		Doubles:            stat.Doubles,
		Triples:            stat.Triples,
		HomeRuns:           stat.HomeRuns,
		Runs:               stat.Runs,
		RBI:                stat.RBI,
		Walks:              stat.Walks,
		Strikeouts:         stat.Strikeouts,
		StolenBases:        stat.StolenBases,
		BattingAvg:         stat.BattingAvg,
		OBP:                stat.OBP,
		SLG:                stat.SLG,
		OPS:                stat.OPS,
		BABIP:              stat.BABIP,
		ExpectedBattingAvg: stat.ExpectedBattingAvg,
		ExpectedSlugging:   stat.ExpectedSlugging,
		ExpectedWOBA:       stat.ExpectedWOBA,
		BarrelPct:          stat.BarrelPct,
		HardHitPct:         stat.HardHitPct,
		AvgExitVelocity:    stat.AvgExitVelocity,
		AvgLaunchAngle:     stat.AvgLaunchAngle,
		SweetSpotPct:       stat.SweetSpotPct,
	}
}

func newPlayerPitchingStats(stat models.SeasonStat) *playerPitchingStats {
	if stat.InningsPitched <= 0 {
		return nil
	}

	return &playerPitchingStats{
		GamesPlayed:        stat.GamesPlayed,
		GamesStarted:       stat.GamesStarted,
		Wins:               stat.Wins,
		Losses:             stat.Losses,
		ERA:                stat.ERA,
		WHIP:               stat.WHIP,
		InningsPitched:     stat.InningsPitched,
		HitsAllowed:        stat.HitsAllowed,
		WalksAllowed:       stat.WalksAllowed,
		HomeRunsAllowed:    stat.HomeRunsAllowed,
		StrikeoutsPer9:     stat.StrikeoutsPer9,
		WalksPer9:          stat.WalksPer9,
		HitsPer9:           stat.HitsPer9,
		HomeRunsPer9:       stat.HomeRunsPer9,
		StrikeoutWalkRatio: stat.StrikeoutWalkRatio,
		StrikePercentage:   stat.StrikePercentage,
		ExpectedERA:        stat.ExpectedERA,
	}
}
