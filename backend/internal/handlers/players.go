package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/jsquardo/capcurve/internal/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

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
		Limit(params.PageSize).
		Offset(playerListOffset(params)).
		Scan(&rows).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := playerListResponse{
		Data: make([]playerListItem, 0, len(rows)),
		Meta: playerListMeta{
			Total:      total,
			Page:       params.Page,
			PageSize:   params.PageSize,
			TotalPages: totalPages(total, params.PageSize),
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
	if err := h.db.First(&player, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "player not found"})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var seasonStats []models.SeasonStat
	if err := h.db.Where("player_id = ?", player.ID).Order("year ASC, id ASC").Find(&seasonStats).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := playerDetailResponse{
		Data: newPlayerDetailItem(player, seasonStats),
	}

	return c.JSON(http.StatusOK, response)
}

type playerListParams struct {
	Query    string
	Active   *bool
	Position string
	Team     string
	Season   *int
	Page     int
	PageSize int
	Sort     string
}

func parsePlayerListParams(c echo.Context) (playerListParams, error) {
	params := playerListParams{
		Query:    strings.TrimSpace(c.QueryParam("q")),
		Position: strings.TrimSpace(c.QueryParam("position")),
		Team:     strings.TrimSpace(c.QueryParam("team")),
		Page:     1,
		PageSize: 25,
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

	if page := strings.TrimSpace(c.QueryParam("page")); page != "" {
		parsed, err := strconv.Atoi(page)
		if err != nil || parsed <= 0 {
			return playerListParams{}, echo.NewHTTPError(http.StatusBadRequest, "invalid page value")
		}
		params.Page = parsed
	}

	if pageSize := strings.TrimSpace(c.QueryParam("page_size")); pageSize != "" {
		parsed, err := strconv.Atoi(pageSize)
		if err != nil || parsed <= 0 {
			return playerListParams{}, echo.NewHTTPError(http.StatusBadRequest, "invalid page_size value")
		}
		if parsed > 100 {
			parsed = 100
		}
		params.PageSize = parsed
	}

	if pageSize := strings.TrimSpace(c.QueryParam("page_size")); pageSize == "" {
		if limit := strings.TrimSpace(c.QueryParam("limit")); limit != "" {
			parsed, err := strconv.Atoi(limit)
			if err != nil || parsed <= 0 {
				return playerListParams{}, echo.NewHTTPError(http.StatusBadRequest, "invalid limit value")
			}
			if parsed > 100 {
				parsed = 100
			}
			params.PageSize = parsed
		}
	}

	if page := strings.TrimSpace(c.QueryParam("page")); page == "" {
		if offset := strings.TrimSpace(c.QueryParam("offset")); offset != "" {
			parsed, err := strconv.Atoi(offset)
			if err != nil || parsed < 0 {
				return playerListParams{}, echo.NewHTTPError(http.StatusBadRequest, "invalid offset value")
			}
			params.Page = (parsed / params.PageSize) + 1
		}
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

func playerListOffset(params playerListParams) int {
	return (params.Page - 1) * params.PageSize
}

func totalPages(total int64, pageSize int) int {
	if total == 0 || pageSize <= 0 {
		return 0
	}

	return int((total + int64(pageSize) - 1) / int64(pageSize))
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
