package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jsquardo/capcurve/internal/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type leaderboardParams struct {
	Category string
	Season   *int
	Page     int
	PageSize int
}

type leaderboardRow struct {
	PlayerID   uint
	PlayerName string
	Position   string
	Team       string
	Value      float64
	Season     *int
}

func (h *Handler) GetLeaderboards(c echo.Context) error {
	params, err := parseLeaderboardParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": leaderboardParamError(err)})
	}

	rowsQuery := h.buildLeaderboardRowsQuery(params)

	var total int64
	if err := h.db.Table("(?) AS leaderboard_rows", rowsQuery).Count(&total).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var rows []leaderboardRow
	if err := h.db.Table("(?) AS leaderboard_rows", rowsQuery).
		Order(leaderboardValueOrder(params.Category)).
		Order("player_name ASC").
		Limit(params.PageSize).
		Offset(leaderboardOffset(params)).
		Scan(&rows).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := leaderboardsResponse{
		Data: leaderboardData{
			Category: params.Category,
			Leaders:  make([]leaderboardItem, 0, len(rows)),
			Meta: leaderboardMeta{
				Total:      total,
				Page:       params.Page,
				PageSize:   params.PageSize,
				TotalPages: totalPages(total, params.PageSize),
			},
		},
	}

	startRank := leaderboardOffset(params) + 1
	for i, row := range rows {
		response.Data.Leaders = append(response.Data.Leaders, row.toResponse(startRank+i))
	}

	return c.JSON(http.StatusOK, response)
}

func parseLeaderboardParams(c echo.Context) (leaderboardParams, error) {
	params := leaderboardParams{
		Category: strings.ToLower(strings.TrimSpace(c.QueryParam("category"))),
		Page:     1,
		PageSize: 25,
	}

	if !isSupportedLeaderboardCategory(params.Category) {
		return leaderboardParams{}, echo.NewHTTPError(http.StatusBadRequest, "invalid category value")
	}

	if season := strings.TrimSpace(c.QueryParam("season")); season != "" {
		parsed, err := strconv.Atoi(season)
		if err != nil || parsed <= 0 {
			return leaderboardParams{}, echo.NewHTTPError(http.StatusBadRequest, "invalid season value")
		}
		params.Season = &parsed
	}

	if page := strings.TrimSpace(c.QueryParam("page")); page != "" {
		parsed, err := strconv.Atoi(page)
		if err != nil || parsed <= 0 {
			return leaderboardParams{}, echo.NewHTTPError(http.StatusBadRequest, "invalid page value")
		}
		params.Page = parsed
	}

	if pageSize := strings.TrimSpace(c.QueryParam("page_size")); pageSize != "" {
		parsed, err := strconv.Atoi(pageSize)
		if err != nil || parsed <= 0 {
			return leaderboardParams{}, echo.NewHTTPError(http.StatusBadRequest, "invalid page_size value")
		}
		if parsed > 100 {
			parsed = 100
		}
		params.PageSize = parsed
	}

	return params, nil
}

func leaderboardParamError(err error) string {
	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		if message, ok := httpErr.Message.(string); ok {
			return message
		}
	}

	return err.Error()
}

func isSupportedLeaderboardCategory(category string) bool {
	switch category {
	case "peak_arc", "hr", "avg", "era", "k9":
		return true
	default:
		return false
	}
}

func (h *Handler) buildLeaderboardRowsQuery(params leaderboardParams) *gorm.DB {
	if params.Category == "peak_arc" {
		return h.buildPeakArcLeaderboardQuery()
	}

	return h.buildSeasonStatLeaderboardQuery(params)
}

func (h *Handler) buildPeakArcLeaderboardQuery() *gorm.DB {
	latestSeason := latestSeasonSnapshotSubquery(h.db)

	// Peak arc remains value-score-first even though the current career_arcs table
	// still stores WAR-era summary columns. We derive the leaderboard value from
	// season_stats inside the stored peak window and only fall back to the overall
	// max season score when the peak window has no matching rows.
	return h.db.Model(&models.CareerArc{}).
		Joins("JOIN players ON players.id = career_arcs.player_id AND players.deleted_at IS NULL").
		Joins("JOIN season_stats ON season_stats.player_id = players.id AND season_stats.deleted_at IS NULL").
		Joins("LEFT JOIN (?) AS latest_season ON latest_season.player_id = players.id", latestSeason).
		Select(`
			players.id AS player_id,
			TRIM(players.first_name || ' ' || players.last_name) AS player_name,
			players.position AS position,
			COALESCE(latest_season.team_name, '') AS team,
			COALESCE(
				MAX(CASE
					WHEN season_stats.year BETWEEN career_arcs.peak_year_start AND career_arcs.peak_year_end
					THEN season_stats.value_score
				END),
				MAX(season_stats.value_score)
			) AS value
		`).
		Group("players.id, players.first_name, players.last_name, players.position, latest_season.team_name")
}

func (h *Handler) buildSeasonStatLeaderboardQuery(params leaderboardParams) *gorm.DB {
	season := leaderboardSeason(params, time.Now())

	query := h.db.Model(&models.SeasonStat{}).
		Joins("JOIN players ON players.id = season_stats.player_id AND players.deleted_at IS NULL").
		Where("season_stats.year = ?", season).
		Select(`
			players.id AS player_id,
			TRIM(players.first_name || ' ' || players.last_name) AS player_name,
			players.position AS position,
			season_stats.team_name AS team,
			season_stats.year AS season
		`)

	switch params.Category {
	case "hr":
		return query.
			Where("season_stats.plate_appearances > 0").
			Select(`
				players.id AS player_id,
				TRIM(players.first_name || ' ' || players.last_name) AS player_name,
				players.position AS position,
				season_stats.team_name AS team,
				season_stats.home_runs AS value,
				season_stats.year AS season
			`)
	case "avg":
		return query.
			Where("season_stats.at_bats > 0").
			Select(`
				players.id AS player_id,
				TRIM(players.first_name || ' ' || players.last_name) AS player_name,
				players.position AS position,
				season_stats.team_name AS team,
				season_stats.batting_avg AS value,
				season_stats.year AS season
			`)
	case "era":
		return query.
			Where("season_stats.innings_pitched > 0").
			Select(`
				players.id AS player_id,
				TRIM(players.first_name || ' ' || players.last_name) AS player_name,
				players.position AS position,
				season_stats.team_name AS team,
				season_stats.era AS value,
				season_stats.year AS season
			`)
	case "k9":
		return query.
			Where("season_stats.innings_pitched > 0").
			Select(`
				players.id AS player_id,
				TRIM(players.first_name || ' ' || players.last_name) AS player_name,
				players.position AS position,
				season_stats.team_name AS team,
				season_stats.strikeouts_per_9 AS value,
				season_stats.year AS season
			`)
	default:
		return query
	}
}

func leaderboardSeason(params leaderboardParams, now time.Time) int {
	if params.Season != nil {
		return *params.Season
	}

	return defaultLeaderboardSeason(now)
}

func defaultLeaderboardSeason(now time.Time) int {
	// Season-based leaderboards should expose the just-finished regular season as
	// soon as the calendar reaches October 1, which is the project's explicit
	// postseason cutoff for default leaderboard queries.
	if now.Month() > time.September {
		return now.Year()
	}

	return now.Year() - 1
}

func leaderboardValueOrder(category string) string {
	if category == "era" {
		return "value ASC"
	}

	return "value DESC"
}

func leaderboardOffset(params leaderboardParams) int {
	return (params.Page - 1) * params.PageSize
}
