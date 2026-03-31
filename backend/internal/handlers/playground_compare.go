package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/jsquardo/capcurve/internal/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type playgroundCompareParams struct {
	PlayerIDs []uint
	Group     string
	Season    *int
	EraStart  *int
	EraEnd    *int
	AgeMin    *int
	AgeMax    *int
	MinPA     *int
	MaxPA     *int
	MinIP     *float64
	MaxIP     *float64
}

func (h *Handler) GetPlaygroundCompare(c echo.Context) error {
	params, err := parsePlaygroundCompareParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": playgroundCompareParamError(err)})
	}

	var players []models.Player
	if err := h.db.
		Where("id IN ?", params.PlayerIDs).
		Order("id ASC").
		Find(&players).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if len(players) != len(params.PlayerIDs) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": formatPlaygroundCompareMissingPlayers(params.PlayerIDs, players)})
	}

	rowsQuery := h.buildPlaygroundCompareQuery(params)

	var rows []playgroundQueryRow
	if err := rowsQuery.
		Select(playgroundQuerySelectColumns()).
		Order("season_stats.player_id ASC").
		Order("season_stats.year ASC").
		Order("season_stats.id ASC").
		Scan(&rows).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	rowsByPlayerID := make(map[uint][]playgroundCompareSeasonItem, len(params.PlayerIDs))
	for _, row := range rows {
		rowsByPlayerID[row.PlayerID] = append(rowsByPlayerID[row.PlayerID], row.toCompareSeasonResponse())
	}

	playerByID := make(map[uint]models.Player, len(players))
	for _, player := range players {
		playerByID[player.ID] = player
	}

	response := playgroundCompareResponse{
		Data: make([]playgroundCompareItem, 0, len(params.PlayerIDs)),
	}

	for _, playerID := range params.PlayerIDs {
		player := playerByID[playerID]
		response.Data = append(response.Data, playgroundCompareItem{
			Player:  newPlaygroundQueryPlayerItem(player),
			Seasons: rowsByPlayerID[playerID],
		})
	}

	return c.JSON(http.StatusOK, response)
}

func parsePlaygroundCompareParams(c echo.Context) (playgroundCompareParams, error) {
	params := playgroundCompareParams{
		Group: strings.ToLower(strings.TrimSpace(c.QueryParam("group"))),
	}

	if params.Group == "" {
		params.Group = "all"
	}

	switch params.Group {
	case "all", "hitting", "pitching":
	default:
		return playgroundCompareParams{}, echo.NewHTTPError(http.StatusBadRequest, "invalid group value")
	}

	playerIDs, err := parsePlaygroundComparePlayerIDs(c.QueryParam("player_ids"))
	if err != nil {
		return playgroundCompareParams{}, err
	}
	params.PlayerIDs = playerIDs

	if params.Season, err = parseOptionalIntQuery(c, "season"); err != nil {
		return playgroundCompareParams{}, err
	}
	if params.EraStart, err = parseOptionalIntQuery(c, "era_start"); err != nil {
		return playgroundCompareParams{}, err
	}
	if params.EraEnd, err = parseOptionalIntQuery(c, "era_end"); err != nil {
		return playgroundCompareParams{}, err
	}
	if params.AgeMin, err = parseOptionalIntQuery(c, "age_min"); err != nil {
		return playgroundCompareParams{}, err
	}
	if params.AgeMax, err = parseOptionalIntQuery(c, "age_max"); err != nil {
		return playgroundCompareParams{}, err
	}
	if params.MinPA, err = parseOptionalIntQuery(c, "min_pa"); err != nil {
		return playgroundCompareParams{}, err
	}
	if params.MaxPA, err = parseOptionalIntQuery(c, "max_pa"); err != nil {
		return playgroundCompareParams{}, err
	}
	if params.MinIP, err = parseOptionalFloatQuery(c, "min_ip"); err != nil {
		return playgroundCompareParams{}, err
	}
	if params.MaxIP, err = parseOptionalFloatQuery(c, "max_ip"); err != nil {
		return playgroundCompareParams{}, err
	}

	if err := validatePlaygroundSeasonAndAgeFilters(
		params.Season,
		params.EraStart,
		params.EraEnd,
		params.AgeMin,
		params.AgeMax,
	); err != nil {
		return playgroundCompareParams{}, err
	}

	if err := validatePlaygroundGroupFilters(
		params.Group,
		playgroundThresholdSupport{
			HasHittingWorkloadFilters:  hasHittingWorkloadFiltersFromCompare(params),
			HasPitchingWorkloadFilters: hasPitchingWorkloadFiltersFromCompare(params),
		},
	); err != nil {
		return playgroundCompareParams{}, err
	}

	return params, nil
}

func parsePlaygroundComparePlayerIDs(raw string) ([]uint, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "player_ids is required")
	}

	parts := strings.Split(raw, ",")
	if len(parts) < 2 || len(parts) > 4 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "player_ids must include between 2 and 4 players")
	}

	playerIDs := make([]uint, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid player_ids value")
		}

		parsed, err := strconv.ParseUint(trimmed, 10, 64)
		if err != nil || parsed == 0 {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid player_ids value")
		}

		playerID := uint(parsed)
		if slices.Contains(playerIDs, playerID) {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "player_ids must be unique")
		}

		playerIDs = append(playerIDs, playerID)
	}

	return playerIDs, nil
}

func (h *Handler) buildPlaygroundCompareQuery(params playgroundCompareParams) *gorm.DB {
	query := h.db.Model(&models.SeasonStat{}).
		Joins("JOIN players ON players.id = season_stats.player_id AND players.deleted_at IS NULL").
		Where("season_stats.player_id IN ?", params.PlayerIDs)

	if params.Season != nil {
		query = query.Where("season_stats.year = ?", *params.Season)
	}
	if params.EraStart != nil {
		query = query.Where("season_stats.year >= ?", *params.EraStart)
	}
	if params.EraEnd != nil {
		query = query.Where("season_stats.year <= ?", *params.EraEnd)
	}
	if params.AgeMin != nil {
		query = query.Where("season_stats.age >= ?", *params.AgeMin)
	}
	if params.AgeMax != nil {
		query = query.Where("season_stats.age <= ?", *params.AgeMax)
	}

	if params.Group == "all" {
		if hasHittingWorkloadFiltersFromCompare(params) {
			query = query.Where("season_stats.plate_appearances > 0")
		}
		if hasPitchingWorkloadFiltersFromCompare(params) {
			query = query.Where("season_stats.innings_pitched > 0")
		}
	}

	if params.MinPA != nil {
		query = query.Where("season_stats.plate_appearances >= ?", *params.MinPA)
	}
	if params.MaxPA != nil {
		query = query.Where("season_stats.plate_appearances <= ?", *params.MaxPA)
	}
	if params.MinIP != nil {
		query = query.Where("season_stats.innings_pitched >= ?", *params.MinIP)
	}
	if params.MaxIP != nil {
		query = query.Where("season_stats.innings_pitched <= ?", *params.MaxIP)
	}

	switch params.Group {
	case "hitting":
		query = query.Where("season_stats.plate_appearances > 0")
	case "pitching":
		query = query.Where("season_stats.innings_pitched > 0")
	}

	return query
}

func (r playgroundQueryRow) toCompareSeasonResponse() playgroundCompareSeasonItem {
	queryItem := r.toResponse()

	return playgroundCompareSeasonItem{
		Season:   queryItem.Season,
		Hitting:  queryItem.Hitting,
		Pitching: queryItem.Pitching,
	}
}

func newPlaygroundQueryPlayerItem(player models.Player) playgroundQueryPlayerItem {
	return playgroundQueryPlayerItem{
		ID:        uint(player.ID),
		MLBID:     player.MLBID,
		FirstName: player.FirstName,
		LastName:  player.LastName,
		FullName:  strings.TrimSpace(player.FirstName + " " + player.LastName),
		Position:  player.Position,
		Bats:      player.Bats,
		Throws:    player.Throws,
		Active:    player.Active,
		ImageURL:  player.ImageURL,
	}
}

func hasHittingWorkloadFiltersFromCompare(params playgroundCompareParams) bool {
	return params.MinPA != nil || params.MaxPA != nil
}

func hasPitchingWorkloadFiltersFromCompare(params playgroundCompareParams) bool {
	return params.MinIP != nil || params.MaxIP != nil
}

func playgroundCompareParamError(err error) string {
	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		if message, ok := httpErr.Message.(string); ok {
			return message
		}
	}

	return err.Error()
}

type playgroundThresholdSupport struct {
	HasHittingThresholds       bool
	HasPitchingThresholds      bool
	HasHittingWorkloadFilters  bool
	HasPitchingWorkloadFilters bool
}

func validatePlaygroundSeasonAndAgeFilters(
	season *int,
	eraStart *int,
	eraEnd *int,
	ageMin *int,
	ageMax *int,
) error {
	if season != nil && (eraStart != nil || eraEnd != nil) {
		return echo.NewHTTPError(http.StatusBadRequest, "season cannot be combined with era_start or era_end")
	}
	if eraStart != nil && eraEnd != nil && *eraStart > *eraEnd {
		return echo.NewHTTPError(http.StatusBadRequest, "era_start must be less than or equal to era_end")
	}
	if ageMin != nil && ageMax != nil && *ageMin > *ageMax {
		return echo.NewHTTPError(http.StatusBadRequest, "age_min must be less than or equal to age_max")
	}

	return nil
}

func validatePlaygroundGroupFilters(group string, support playgroundThresholdSupport) error {
	if group == "hitting" && (support.HasPitchingThresholds || support.HasPitchingWorkloadFilters) {
		return echo.NewHTTPError(http.StatusBadRequest, "pitching thresholds are not supported for group=hitting")
	}
	if group == "pitching" && (support.HasHittingThresholds || support.HasHittingWorkloadFilters) {
		return echo.NewHTTPError(http.StatusBadRequest, "hitting thresholds are not supported for group=pitching")
	}

	return nil
}

func formatPlaygroundCompareMissingPlayers(playerIDs []uint, players []models.Player) string {
	found := make(map[uint]struct{}, len(players))
	for _, player := range players {
		found[player.ID] = struct{}{}
	}

	missing := make([]string, 0)
	for _, playerID := range playerIDs {
		if _, ok := found[playerID]; !ok {
			missing = append(missing, strconv.FormatUint(uint64(playerID), 10))
		}
	}

	return fmt.Sprintf("player not found: %s", strings.Join(missing, ","))
}
