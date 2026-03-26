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
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
	Count  int   `json:"count"`
	Total  int64 `json:"total"`
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

func (h *Handler) ListPlayers(c echo.Context) error {
	params, err := parsePlayerListParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": playerListParamError(err)})
	}

	query := h.buildPlayerListQuery(params)

	var total int64
	if err := query.Session(&gorm.Session{}).
		Distinct("players.id").
		Count(&total).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var rows []playerListRow
	if err := h.buildPlayerListOrderedQuery(params).
		Select(playerListSelectColumns()).
		Order("players.id ASC").
		Limit(params.Limit).
		Offset(params.Offset).
		Scan(&rows).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := playerListResponse{
		Data: make([]playerListItem, 0, len(rows)),
		Meta: playerListMeta{
			Limit:  params.Limit,
			Offset: params.Offset,
			Count:  len(rows),
			Total:  total,
		},
	}

	for _, row := range rows {
		response.Data = append(response.Data, row.toResponse())
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) GetPlayer(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid player id"})
	}

	var player models.Player
	if err := h.db.Preload("SeasonStats").Preload("Contracts").Preload("CareerArc").First(&player, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "player not found"})
	}
	return c.JSON(http.StatusOK, player)
}

type playerListParams struct {
	Query    string
	Active   *bool
	Position string
	Team     string
	Season   *int
	Limit    int
	Offset   int
	Sort     string
}

func parsePlayerListParams(c echo.Context) (playerListParams, error) {
	params := playerListParams{
		Query:    strings.TrimSpace(c.QueryParam("q")),
		Position: strings.TrimSpace(c.QueryParam("position")),
		Team:     strings.TrimSpace(c.QueryParam("team")),
		Limit:    25,
		Offset:   0,
		Sort:     strings.TrimSpace(c.QueryParam("sort")),
	}

	if params.Sort == "" {
		params.Sort = "name"
	}

	if active := strings.TrimSpace(c.QueryParam("active")); active != "" {
		parsed, err := strconv.ParseBool(active)
		if err != nil {
			return playerListParams{}, echo.NewHTTPError(http.StatusBadRequest, "invalid active value")
		}
		params.Active = &parsed
	}

	if season := strings.TrimSpace(c.QueryParam("season")); season != "" {
		parsed, err := strconv.Atoi(season)
		if err != nil {
			return playerListParams{}, echo.NewHTTPError(http.StatusBadRequest, "invalid season value")
		}
		params.Season = &parsed
	}

	if limit := strings.TrimSpace(c.QueryParam("limit")); limit != "" {
		parsed, err := strconv.Atoi(limit)
		if err != nil || parsed <= 0 {
			return playerListParams{}, echo.NewHTTPError(http.StatusBadRequest, "invalid limit value")
		}
		if parsed > 100 {
			parsed = 100
		}
		params.Limit = parsed
	}

	if offset := strings.TrimSpace(c.QueryParam("offset")); offset != "" {
		parsed, err := strconv.Atoi(offset)
		if err != nil || parsed < 0 {
			return playerListParams{}, echo.NewHTTPError(http.StatusBadRequest, "invalid offset value")
		}
		params.Offset = parsed
	}

	switch params.Sort {
	case "name", "-name", "value_score", "-value_score", "recent_year", "-recent_year":
		return params, nil
	default:
		return playerListParams{}, echo.NewHTTPError(http.StatusBadRequest, "invalid sort value")
	}
}

func (h *Handler) buildPlayerListQuery(params playerListParams) *gorm.DB {
	snapshot := seasonSnapshotSubquery(h.db, params.Season)

	query := h.db.Model(&models.Player{}).
		Joins("LEFT JOIN (?) AS season_snapshot ON season_snapshot.player_id = players.id", snapshot)

	if params.Query != "" {
		searchTerm := "%" + params.Query + "%"
		query = query.Where(
			"players.first_name ILIKE ? OR players.last_name ILIKE ? OR (players.first_name || ' ' || players.last_name) ILIKE ?",
			searchTerm, searchTerm, searchTerm,
		)
	}

	if params.Active != nil {
		query = query.Where("players.active = ?", *params.Active)
	}

	if params.Position != "" {
		query = query.Where("players.position = ?", params.Position)
	}

	if params.Team != "" {
		if teamID, err := strconv.Atoi(params.Team); err == nil {
			query = query.Where("season_snapshot.team_id = ?", teamID)
		} else {
			query = query.Where("season_snapshot.team_name ILIKE ?", "%"+params.Team+"%")
		}
	}

	return query
}

func (h *Handler) buildPlayerListOrderedQuery(params playerListParams) *gorm.DB {
	query := h.buildPlayerListQuery(params)

	for _, order := range playerListOrderBy(params.Sort) {
		query = query.Order(order)
	}

	return query
}

func playerListParamError(err error) string {
	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		if message, ok := httpErr.Message.(string); ok {
			return message
		}
	}

	return err.Error()
}

func seasonSnapshotSubquery(db *gorm.DB, season *int) *gorm.DB {
	if season != nil {
		return seasonSnapshotForYearSubquery(db, *season)
	}

	return latestSeasonSnapshotSubquery(db)
}

func latestSeasonSnapshotSubquery(db *gorm.DB) *gorm.DB {
	// DISTINCT ON keeps one deterministic season row per player so the outer list
	// can sort/filter by the same "latest season" snapshot without duplicating players.
	return db.Model(&models.SeasonStat{}).
		Select(`
			DISTINCT ON (player_id)
			player_id,
			year,
			team_id,
			team_name,
			age,
			value_score
		`).
		Order("player_id, year DESC, id DESC")
}

func seasonSnapshotForYearSubquery(db *gorm.DB, season int) *gorm.DB {
	// When the client scopes the list to a season, the joined snapshot should be
	// that exact year rather than the latest available season.
	return db.Model(&models.SeasonStat{}).
		Select(`
			player_id,
			year,
			team_id,
			team_name,
			age,
			value_score
		`).
		Where("year = ?", season)
}

func playerListSelectColumns() []string {
	return []string{
		"players.id",
		"players.mlb_id",
		"players.first_name",
		"players.last_name",
		"players.position",
		"players.bats",
		"players.throws",
		"players.date_of_birth",
		"players.active",
		"players.image_url",
		"season_snapshot.year AS latest_season_year",
		"season_snapshot.team_id AS latest_team_id",
		"season_snapshot.team_name AS latest_team_name",
		"season_snapshot.age AS latest_age",
		"season_snapshot.value_score AS latest_value_score",
	}
}

func playerListOrderBy(sort string) []string {
	switch sort {
	case "-name":
		return []string{"players.last_name DESC", "players.first_name DESC"}
	case "value_score":
		return []string{
			"season_snapshot.value_score ASC NULLS LAST",
			"players.last_name ASC",
			"players.first_name ASC",
		}
	case "-value_score":
		return []string{
			"season_snapshot.value_score DESC NULLS LAST",
			"players.last_name ASC",
			"players.first_name ASC",
		}
	case "recent_year":
		return []string{
			"season_snapshot.year ASC NULLS LAST",
			"players.last_name ASC",
			"players.first_name ASC",
		}
	case "-recent_year":
		return []string{
			"season_snapshot.year DESC NULLS LAST",
			"players.last_name ASC",
			"players.first_name ASC",
		}
	default:
		return []string{"players.last_name ASC", "players.first_name ASC"}
	}
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

func (h *Handler) GetCareerArc(c echo.Context) error {
	playerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid player id"})
	}

	var arc models.CareerArc
	if err := h.db.Where("player_id = ?", playerID).First(&arc).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "career arc not found"})
	}

	var stats []models.SeasonStat
	h.db.Where("player_id = ?", playerID).Order("year ASC").Find(&stats)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"arc":          arc,
		"season_stats": stats,
	})
}

func (h *Handler) GetPlayerContracts(c echo.Context) error {
	playerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid player id"})
	}

	var contracts []models.Contract
	if err := h.db.Preload("ContractSeasons").Where("player_id = ?", playerID).Find(&contracts).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, contracts)
}

func (h *Handler) GetContract(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid contract id"})
	}

	var contract models.Contract
	if err := h.db.Preload("ContractSeasons").First(&contract, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "contract not found"})
	}
	return c.JSON(http.StatusOK, contract)
}

func (h *Handler) MostOverpaid(c echo.Context) error {
	var contracts []models.Contract
	h.db.Where("is_active = ? AND overall_value_score < ?", true, -20).
		Order("overall_value_score ASC").Limit(25).Find(&contracts)
	return c.JSON(http.StatusOK, contracts)
}

func (h *Handler) BestValue(c echo.Context) error {
	var contracts []models.Contract
	h.db.Where("is_active = ? AND overall_value_score > ?", true, 20).
		Order("overall_value_score DESC").Limit(25).Find(&contracts)
	return c.JSON(http.StatusOK, contracts)
}

func (h *Handler) PeakArcs(c echo.Context) error {
	var arcs []models.CareerArc
	h.db.Order("peak_war DESC").Limit(25).Find(&arcs)
	return c.JSON(http.StatusOK, arcs)
}
