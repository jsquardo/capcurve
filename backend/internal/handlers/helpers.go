package handlers

import (
	"strings"
	"time"

	"github.com/jsquardo/capcurve/internal/models"
)

type playerListItem struct {
	ID           uint                  `json:"id"`
	MLBID        int                   `json:"mlb_id"`
	FirstName    string                `json:"first_name"`
	LastName     string                `json:"last_name"`
	FullName     string                `json:"full_name"`
	Position     string                `json:"position"`
	Bats         string                `json:"bats"`
	Throws       string                `json:"throws"`
	DateOfBirth  *time.Time            `json:"date_of_birth"`
	Active       bool                  `json:"active"`
	ImageURL     string                `json:"image_url"`
	LatestSeason *playerSeasonListItem `json:"latest_season"`
}

type playerSeasonListItem struct {
	Year       int     `json:"year"`
	TeamID     int     `json:"team_id"`
	TeamName   string  `json:"team_name"`
	Age        int     `json:"age"`
	ValueScore float64 `json:"value_score"`
}

type playerListMeta struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

type playerListResponse struct {
	Data []playerListItem `json:"data"`
	Meta playerListMeta   `json:"meta"`
}

type playerListRow struct {
	ID               uint
	MLBID            int
	FirstName        string
	LastName         string
	Position         string
	Bats             string
	Throws           string
	DateOfBirth      *time.Time
	Active           bool
	ImageURL         string
	LatestSeasonYear *int
	LatestTeamID     *int
	LatestTeamName   *string
	LatestAge        *int
	LatestValueScore *float64
}

type playerDetailResponse struct {
	Data playerDetailItem `json:"data"`
}

type careerArcResponse struct {
	Data careerArcData `json:"data"`
}

type projectionResponse struct {
	Data projectionData `json:"data"`
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

type careerArcData struct {
	Player     careerArcPlayerItem     `json:"player"`
	Arc        *careerArcMetadata      `json:"arc"`
	Timeline   []careerArcTimelineItem `json:"timeline"`
	Projection careerArcProjection     `json:"projection"`
}

type projectionData struct {
	Player     careerArcPlayerItem `json:"player"`
	Projection careerArcProjection `json:"projection"`
}

type careerArcPlayerItem struct {
	ID          uint       `json:"id"`
	MLBID       int        `json:"mlb_id"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	FullName    string     `json:"full_name"`
	Position    string     `json:"position"`
	Bats        string     `json:"bats"`
	Throws      string     `json:"throws"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Active      bool       `json:"active"`
	ImageURL    string     `json:"image_url"`
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

type careerArcTimelineItem struct {
	Year         int                  `json:"year"`
	TeamID       int                  `json:"team_id"`
	TeamName     string               `json:"team_name"`
	Age          int                  `json:"age"`
	ValueScore   float64              `json:"value_score"`
	IsPeak       bool                 `json:"is_peak"`
	IsProjection bool                 `json:"is_projection"`
	Hitting      *playerHittingStats  `json:"hitting"`
	Pitching     *playerPitchingStats `json:"pitching"`
}

type careerArcMetadata struct {
	PeakYearStart         int       `json:"peak_year_start"`
	PeakYearEnd           int       `json:"peak_year_end"`
	DeclineOnsetYear      int       `json:"decline_onset_year"`
	ArcShape              string    `json:"arc_shape"`
	PeakValueScore        float64   `json:"peak_value_score"`
	CareerValueScoreTotal float64   `json:"career_value_score_total"`
	LastComputedAt        time.Time `json:"last_computed_at"`
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

type careerArcProjection struct {
	Status         string                      `json:"status"`
	Eligible       bool                        `json:"eligible"`
	Reason         string                      `json:"reason"`
	Points         []careerArcProjectionPoint  `json:"points"`
	ConfidenceBand []careerArcConfidenceBand   `json:"confidence_band"`
	Comparables    []careerArcComparablePlayer `json:"comparables"`
}

type careerArcProjectionPoint struct {
	Year         int     `json:"year"`
	Age          int     `json:"age"`
	ValueScore   float64 `json:"value_score"`
	IsProjection bool    `json:"is_projection"`
}

type careerArcConfidenceBand struct {
	Year  int     `json:"year"`
	Lower float64 `json:"lower"`
	Upper float64 `json:"upper"`
}

type careerArcComparablePlayer struct {
	PlayerID uint   `json:"player_id"`
	MLBID    int    `json:"mlb_id"`
	FullName string `json:"full_name"`
	Position string `json:"position"`
}

func (r playerListRow) toResponse() playerListItem {
	item := playerListItem{
		ID:          r.ID,
		MLBID:       r.MLBID,
		FirstName:   r.FirstName,
		LastName:    r.LastName,
		FullName:    strings.TrimSpace(r.FirstName + " " + r.LastName),
		Position:    r.Position,
		Bats:        r.Bats,
		Throws:      r.Throws,
		DateOfBirth: r.DateOfBirth,
		Active:      r.Active,
		ImageURL:    r.ImageURL,
	}

	if r.LatestSeasonYear != nil {
		item.LatestSeason = &playerSeasonListItem{
			Year:       *r.LatestSeasonYear,
			TeamID:     derefInt(r.LatestTeamID),
			TeamName:   derefString(r.LatestTeamName),
			Age:        derefInt(r.LatestAge),
			ValueScore: derefFloat64(r.LatestValueScore),
		}
	}

	return item
}

func derefInt(value *int) int {
	if value == nil {
		return 0
	}

	return *value
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}

func derefFloat64(value *float64) float64 {
	if value == nil {
		return 0
	}

	return *value
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

func newCareerArcData(player models.Player, seasonStats []models.SeasonStat, arc *models.CareerArc, projectionPayload careerArcProjection) careerArcData {
	timeline := make([]careerArcTimelineItem, 0, len(seasonStats))
	for _, stat := range seasonStats {
		timeline = append(timeline, newCareerArcTimelineItem(stat, arc))
	}

	return careerArcData{
		Player:     newCareerArcPlayerItem(player),
		Arc:        newCareerArcMetadata(arc, timeline),
		Timeline:   timeline,
		Projection: projectionPayload,
	}
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

func newCareerArcPlayerItem(player models.Player) careerArcPlayerItem {
	return careerArcPlayerItem{
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
	}
}

func newCareerArcTimelineItem(stat models.SeasonStat, arc *models.CareerArc) careerArcTimelineItem {
	detail := newPlayerCareerStatItem(stat)

	return careerArcTimelineItem{
		Year:         detail.Year,
		TeamID:       detail.TeamID,
		TeamName:     detail.TeamName,
		Age:          detail.Age,
		ValueScore:   detail.ValueScore,
		IsPeak:       isPeakSeason(detail.Year, arc),
		IsProjection: false,
		Hitting:      detail.Hitting,
		Pitching:     detail.Pitching,
	}
}

func newCareerArcMetadata(arc *models.CareerArc, timeline []careerArcTimelineItem) *careerArcMetadata {
	if arc == nil {
		return nil
	}

	return &careerArcMetadata{
		PeakYearStart:         arc.PeakYearStart,
		PeakYearEnd:           arc.PeakYearEnd,
		DeclineOnsetYear:      arcDeclineOnsetYear(arc),
		ArcShape:              arc.ArcShape,
		PeakValueScore:        peakTimelineValueScore(timeline, arc),
		CareerValueScoreTotal: totalTimelineValueScore(timeline),
		LastComputedAt:        arc.LastComputedAt,
	}
}

func isPeakSeason(year int, arc *models.CareerArc) bool {
	if arc == nil || arc.PeakYearStart == 0 || arc.PeakYearEnd == 0 {
		return false
	}

	return year >= arc.PeakYearStart && year <= arc.PeakYearEnd
}

func arcDeclineOnsetYear(arc *models.CareerArc) int {
	if arc == nil {
		return 0
	}

	return arc.DeclineOnsetYear
}

func peakTimelineValueScore(timeline []careerArcTimelineItem, arc *models.CareerArc) float64 {
	if len(timeline) == 0 {
		return 0
	}

	peak, ok := peakTimelineValueScoreInWindow(timeline, arc)
	if ok {
		return peak
	}

	return peakTimelineValueScoreAcrossTimeline(timeline)
}

func peakTimelineValueScoreInWindow(timeline []careerArcTimelineItem, arc *models.CareerArc) (float64, bool) {
	if arc == nil || arc.PeakYearStart == 0 || arc.PeakYearEnd == 0 {
		return 0, false
	}

	var (
		peak  float64
		found bool
	)

	for _, point := range timeline {
		if point.Year < arc.PeakYearStart || point.Year > arc.PeakYearEnd {
			continue
		}

		if !found || point.ValueScore > peak {
			peak = point.ValueScore
			found = true
		}
	}

	return peak, found
}

func peakTimelineValueScoreAcrossTimeline(timeline []careerArcTimelineItem) float64 {
	peak := timeline[0].ValueScore
	for _, point := range timeline[1:] {
		if point.ValueScore > peak {
			peak = point.ValueScore
		}
	}

	return peak
}

func totalTimelineValueScore(timeline []careerArcTimelineItem) float64 {
	var total float64
	for _, point := range timeline {
		total += point.ValueScore
	}

	return total
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
