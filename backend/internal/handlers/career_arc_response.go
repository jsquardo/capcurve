package handlers

import (
	"strings"
	"time"

	"github.com/jsquardo/capcurve/internal/models"
)

type careerArcResponse struct {
	Data careerArcData `json:"data"`
}

type projectionResponse struct {
	Data projectionData `json:"data"`
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
