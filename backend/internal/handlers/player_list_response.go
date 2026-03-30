package handlers

import (
	"strings"
	"time"
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
