package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jsquardo/capcurve/internal/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type playgroundQueryParams struct {
	Query         string
	Group         string
	Positions     []string
	Team          string
	Active        *bool
	Season        *int
	EraStart      *int
	EraEnd        *int
	AgeMin        *int
	AgeMax        *int
	MinPA         *int
	MaxPA         *int
	MinIP         *float64
	MaxIP         *float64
	MinValueScore *float64
	MaxValueScore *float64
	MinHR         *int
	MaxHR         *int
	MinAvg        *float64
	MaxAvg        *float64
	MinOBP        *float64
	MaxOBP        *float64
	MinSLG        *float64
	MaxSLG        *float64
	MinSB         *int
	MaxSB         *int
	MinERA        *float64
	MaxERA        *float64
	MinWHIP       *float64
	MaxWHIP       *float64
	MinK9         *float64
	MaxK9         *float64
	Page          int
	PageSize      int
	Offset        int
	Sort          string
}

type playgroundQueryRow struct {
	PlayerID           uint     `gorm:"column:player_id"`
	MLBID              int      `gorm:"column:mlb_id"`
	FirstName          string   `gorm:"column:first_name"`
	LastName           string   `gorm:"column:last_name"`
	Position           string   `gorm:"column:position"`
	Bats               string   `gorm:"column:bats"`
	Throws             string   `gorm:"column:throws"`
	Active             bool     `gorm:"column:active"`
	ImageURL           string   `gorm:"column:image_url"`
	Year               int      `gorm:"column:year"`
	TeamID             int      `gorm:"column:team_id"`
	TeamName           string   `gorm:"column:team_name"`
	Age                int      `gorm:"column:age"`
	ValueScore         float64  `gorm:"column:value_score"`
	GamesPlayed        int      `gorm:"column:games_played"`
	GamesStarted       int      `gorm:"column:games_started"`
	PlateAppearances   int      `gorm:"column:plate_appearances"`
	AtBats             int      `gorm:"column:at_bats"`
	Hits               int      `gorm:"column:hits"`
	Doubles            int      `gorm:"column:doubles"`
	Triples            int      `gorm:"column:triples"`
	HomeRuns           int      `gorm:"column:home_runs"`
	Runs               int      `gorm:"column:runs"`
	RBI                int      `gorm:"column:rbi"`
	Walks              int      `gorm:"column:walks"`
	Strikeouts         int      `gorm:"column:strikeouts"`
	StolenBases        int      `gorm:"column:stolen_bases"`
	BattingAvg         float64  `gorm:"column:batting_avg"`
	OBP                float64  `gorm:"column:obp"`
	SLG                float64  `gorm:"column:slg"`
	OPS                float64  `gorm:"column:ops"`
	BABIP              float64  `gorm:"column:babip"`
	Wins               int      `gorm:"column:wins"`
	Losses             int      `gorm:"column:losses"`
	ERA                float64  `gorm:"column:era"`
	WHIP               float64  `gorm:"column:whip"`
	InningsPitched     float64  `gorm:"column:innings_pitched"`
	HitsAllowed        int      `gorm:"column:hits_allowed"`
	WalksAllowed       int      `gorm:"column:walks_allowed"`
	HomeRunsAllowed    int      `gorm:"column:home_runs_allowed"`
	StrikeoutsPer9     float64  `gorm:"column:strikeouts_per_9"`
	WalksPer9          float64  `gorm:"column:walks_per_9"`
	HitsPer9           float64  `gorm:"column:hits_per_9"`
	HomeRunsPer9       float64  `gorm:"column:home_runs_per_9"`
	StrikeoutWalkRatio float64  `gorm:"column:strikeout_walk_ratio"`
	StrikePercentage   float64  `gorm:"column:strike_percentage"`
	ExpectedBattingAvg *float64 `gorm:"column:expected_batting_avg"`
	ExpectedSlugging   *float64 `gorm:"column:expected_slugging"`
	ExpectedWOBA       *float64 `gorm:"column:expected_woba"`
	ExpectedERA        *float64 `gorm:"column:expected_era"`
	BarrelPct          *float64 `gorm:"column:barrel_pct"`
	HardHitPct         *float64 `gorm:"column:hard_hit_pct"`
	AvgExitVelocity    *float64 `gorm:"column:avg_exit_velocity"`
	AvgLaunchAngle     *float64 `gorm:"column:avg_launch_angle"`
	SweetSpotPct       *float64 `gorm:"column:sweet_spot_pct"`
}

func (h *Handler) GetPlaygroundQuery(c echo.Context) error {
	params, err := parsePlaygroundQueryParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": playgroundQueryParamError(err)})
	}

	query := h.buildPlaygroundQuery(params)

	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var rows []playgroundQueryRow
	if err := h.buildPlaygroundOrderedQuery(params).
		Select(playgroundQuerySelectColumns()).
		Limit(params.PageSize).
		Offset(params.paginationOffset()).
		Scan(&rows).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := playgroundQueryResponse{
		Data: make([]playgroundQueryItem, 0, len(rows)),
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

func parsePlaygroundQueryParams(c echo.Context) (playgroundQueryParams, error) {
	params := playgroundQueryParams{
		Query:    strings.TrimSpace(c.QueryParam("q")),
		Group:    strings.ToLower(strings.TrimSpace(c.QueryParam("group"))),
		Team:     strings.TrimSpace(c.QueryParam("team")),
		Page:     1,
		PageSize: 25,
		Sort:     strings.TrimSpace(c.QueryParam("sort")),
	}

	if params.Group == "" {
		params.Group = "all"
	}
	if params.Sort == "" {
		params.Sort = "-value_score"
	}

	switch params.Group {
	case "all", "hitting", "pitching":
	default:
		return playgroundQueryParams{}, echo.NewHTTPError(http.StatusBadRequest, "invalid group value")
	}

	if positions := strings.TrimSpace(c.QueryParam("position")); positions != "" {
		for _, position := range strings.Split(positions, ",") {
			trimmed := strings.TrimSpace(position)
			if trimmed != "" {
				params.Positions = append(params.Positions, trimmed)
			}
		}
	}

	var err error
	if params.Active, err = parseOptionalBoolQuery(c, "active"); err != nil {
		return playgroundQueryParams{}, err
	}
	if params.Season, err = parseOptionalIntQuery(c, "season"); err != nil {
		return playgroundQueryParams{}, err
	}
	if params.EraStart, err = parseOptionalIntQuery(c, "era_start"); err != nil {
		return playgroundQueryParams{}, err
	}
	if params.EraEnd, err = parseOptionalIntQuery(c, "era_end"); err != nil {
		return playgroundQueryParams{}, err
	}
	if params.AgeMin, err = parseOptionalIntQuery(c, "age_min"); err != nil {
		return playgroundQueryParams{}, err
	}
	if params.AgeMax, err = parseOptionalIntQuery(c, "age_max"); err != nil {
		return playgroundQueryParams{}, err
	}
	if params.MinPA, err = parseOptionalIntQuery(c, "min_pa"); err != nil {
		return playgroundQueryParams{}, err
	}
	if params.MaxPA, err = parseOptionalIntQuery(c, "max_pa"); err != nil {
		return playgroundQueryParams{}, err
	}
	if params.MinIP, err = parseOptionalFloatQuery(c, "min_ip"); err != nil {
		return playgroundQueryParams{}, err
	}
	if params.MaxIP, err = parseOptionalFloatQuery(c, "max_ip"); err != nil {
		return playgroundQueryParams{}, err
	}
	if params.MinValueScore, err = parseOptionalFloatQuery(c, "min_value_score"); err != nil {
		return playgroundQueryParams{}, err
	}
	if params.MaxValueScore, err = parseOptionalFloatQuery(c, "max_value_score"); err != nil {
		return playgroundQueryParams{}, err
	}
	if params.MinHR, err = parseOptionalIntQuery(c, "min_hr"); err != nil {
		return playgroundQueryParams{}, err
	}
	if params.MaxHR, err = parseOptionalIntQuery(c, "max_hr"); err != nil {
		return playgroundQueryParams{}, err
	}
	if params.MinAvg, err = parseOptionalFloatQuery(c, "min_avg"); err != nil {
		return playgroundQueryParams{}, err
	}
	if params.MaxAvg, err = parseOptionalFloatQuery(c, "max_avg"); err != nil {
		return playgroundQueryParams{}, err
	}
	if params.MinOBP, err = parseOptionalFloatQuery(c, "min_obp"); err != nil {
		return playgroundQueryParams{}, err
	}
	if params.MaxOBP, err = parseOptionalFloatQuery(c, "max_obp"); err != nil {
		return playgroundQueryParams{}, err
	}
	if params.MinSLG, err = parseOptionalFloatQuery(c, "min_slg"); err != nil {
		return playgroundQueryParams{}, err
	}
	if params.MaxSLG, err = parseOptionalFloatQuery(c, "max_slg"); err != nil {
		return playgroundQueryParams{}, err
	}
	if params.MinSB, err = parseOptionalIntQuery(c, "min_sb"); err != nil {
		return playgroundQueryParams{}, err
	}
	if params.MaxSB, err = parseOptionalIntQuery(c, "max_sb"); err != nil {
		return playgroundQueryParams{}, err
	}
	if params.MinERA, err = parseOptionalFloatQuery(c, "min_era"); err != nil {
		return playgroundQueryParams{}, err
	}
	if params.MaxERA, err = parseOptionalFloatQuery(c, "max_era"); err != nil {
		return playgroundQueryParams{}, err
	}
	if params.MinWHIP, err = parseOptionalFloatQuery(c, "min_whip"); err != nil {
		return playgroundQueryParams{}, err
	}
	if params.MaxWHIP, err = parseOptionalFloatQuery(c, "max_whip"); err != nil {
		return playgroundQueryParams{}, err
	}
	if params.MinK9, err = parseOptionalFloatQuery(c, "min_k9"); err != nil {
		return playgroundQueryParams{}, err
	}
	if params.MaxK9, err = parseOptionalFloatQuery(c, "max_k9"); err != nil {
		return playgroundQueryParams{}, err
	}

	if page := strings.TrimSpace(c.QueryParam("page")); page != "" {
		parsed, err := strconv.Atoi(page)
		if err != nil || parsed <= 0 {
			return playgroundQueryParams{}, echo.NewHTTPError(http.StatusBadRequest, "invalid page value")
		}
		params.Page = parsed
	}

	if pageSize := strings.TrimSpace(c.QueryParam("page_size")); pageSize != "" {
		parsed, err := strconv.Atoi(pageSize)
		if err != nil || parsed <= 0 {
			return playgroundQueryParams{}, echo.NewHTTPError(http.StatusBadRequest, "invalid page_size value")
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
				return playgroundQueryParams{}, echo.NewHTTPError(http.StatusBadRequest, "invalid limit value")
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
				return playgroundQueryParams{}, echo.NewHTTPError(http.StatusBadRequest, "invalid offset value")
			}
			params.Offset = parsed
			params.Page = (parsed / params.PageSize) + 1
		}
	}

	if params.Season != nil && (params.EraStart != nil || params.EraEnd != nil) {
		return playgroundQueryParams{}, echo.NewHTTPError(http.StatusBadRequest, "season cannot be combined with era_start or era_end")
	}
	if params.EraStart != nil && params.EraEnd != nil && *params.EraStart > *params.EraEnd {
		return playgroundQueryParams{}, echo.NewHTTPError(http.StatusBadRequest, "era_start must be less than or equal to era_end")
	}
	if params.AgeMin != nil && params.AgeMax != nil && *params.AgeMin > *params.AgeMax {
		return playgroundQueryParams{}, echo.NewHTTPError(http.StatusBadRequest, "age_min must be less than or equal to age_max")
	}

	if params.Group == "hitting" && hasPitchingThresholds(params) {
		return playgroundQueryParams{}, echo.NewHTTPError(http.StatusBadRequest, "pitching thresholds are not supported for group=hitting")
	}
	if params.Group == "pitching" && hasHittingThresholds(params) {
		return playgroundQueryParams{}, echo.NewHTTPError(http.StatusBadRequest, "hitting thresholds are not supported for group=pitching")
	}

	if !isSupportedPlaygroundSort(params.Sort) {
		return playgroundQueryParams{}, echo.NewHTTPError(http.StatusBadRequest, "invalid sort value")
	}

	return params, nil
}

func (h *Handler) buildPlaygroundQuery(params playgroundQueryParams) *gorm.DB {
	query := h.db.Model(&models.SeasonStat{}).
		Joins("JOIN players ON players.id = season_stats.player_id AND players.deleted_at IS NULL")

	if params.Query != "" {
		searchTerm := "%" + params.Query + "%"
		query = query.Where(
			"players.first_name ILIKE ? OR players.last_name ILIKE ? OR TRIM(players.first_name || ' ' || players.last_name) ILIKE ?",
			searchTerm, searchTerm, searchTerm,
		)
	}

	if len(params.Positions) > 0 {
		query = query.Where("players.position IN ?", params.Positions)
	}
	if params.Team != "" {
		if teamID, err := strconv.Atoi(params.Team); err == nil {
			query = query.Where("season_stats.team_id = ?", teamID)
		} else {
			query = query.Where("season_stats.team_name ILIKE ?", "%"+params.Team+"%")
		}
	}
	if params.Active != nil {
		query = query.Where("players.active = ?", *params.Active)
	}
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
		if hasHittingThresholds(params) || hasHittingWorkloadFilters(params) {
			query = query.Where("season_stats.plate_appearances > 0")
		}
		if hasPitchingThresholds(params) || hasPitchingWorkloadFilters(params) {
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
	if params.MinValueScore != nil {
		query = query.Where("season_stats.value_score >= ?", *params.MinValueScore)
	}
	if params.MaxValueScore != nil {
		query = query.Where("season_stats.value_score <= ?", *params.MaxValueScore)
	}
	if params.MinHR != nil {
		query = query.Where("season_stats.home_runs >= ?", *params.MinHR)
	}
	if params.MaxHR != nil {
		query = query.Where("season_stats.home_runs <= ?", *params.MaxHR)
	}
	if params.MinAvg != nil {
		query = query.Where("season_stats.batting_avg >= ?", *params.MinAvg)
	}
	if params.MaxAvg != nil {
		query = query.Where("season_stats.batting_avg <= ?", *params.MaxAvg)
	}
	if params.MinOBP != nil {
		query = query.Where("season_stats.obp >= ?", *params.MinOBP)
	}
	if params.MaxOBP != nil {
		query = query.Where("season_stats.obp <= ?", *params.MaxOBP)
	}
	if params.MinSLG != nil {
		query = query.Where("season_stats.slg >= ?", *params.MinSLG)
	}
	if params.MaxSLG != nil {
		query = query.Where("season_stats.slg <= ?", *params.MaxSLG)
	}
	if params.MinSB != nil {
		query = query.Where("season_stats.stolen_bases >= ?", *params.MinSB)
	}
	if params.MaxSB != nil {
		query = query.Where("season_stats.stolen_bases <= ?", *params.MaxSB)
	}
	if params.MinERA != nil {
		query = query.Where("season_stats.era >= ?", *params.MinERA)
	}
	if params.MaxERA != nil {
		query = query.Where("season_stats.era <= ?", *params.MaxERA)
	}
	if params.MinWHIP != nil {
		query = query.Where("season_stats.whip >= ?", *params.MinWHIP)
	}
	if params.MaxWHIP != nil {
		query = query.Where("season_stats.whip <= ?", *params.MaxWHIP)
	}
	if params.MinK9 != nil {
		query = query.Where("season_stats.strikeouts_per_9 >= ?", *params.MinK9)
	}
	if params.MaxK9 != nil {
		query = query.Where("season_stats.strikeouts_per_9 <= ?", *params.MaxK9)
	}

	switch params.Group {
	case "hitting":
		query = query.Where("season_stats.plate_appearances > 0")
	case "pitching":
		query = query.Where("season_stats.innings_pitched > 0")
	}

	return query
}

func (h *Handler) buildPlaygroundOrderedQuery(params playgroundQueryParams) *gorm.DB {
	query := h.buildPlaygroundQuery(params)
	for _, order := range playgroundQueryOrderBy(params.Sort) {
		query = query.Order(order)
	}
	return query
}

func playgroundQuerySelectColumns() string {
	return `
		players.id AS player_id,
		players.mlb_id AS mlb_id,
		players.first_name AS first_name,
		players.last_name AS last_name,
		players.position AS position,
		players.bats AS bats,
		players.throws AS throws,
		players.active AS active,
		players.image_url AS image_url,
		season_stats.year AS year,
		season_stats.team_id AS team_id,
		season_stats.team_name AS team_name,
		season_stats.age AS age,
		season_stats.value_score AS value_score,
		season_stats.games_played AS games_played,
		season_stats.games_started AS games_started,
		season_stats.plate_appearances AS plate_appearances,
		season_stats.at_bats AS at_bats,
		season_stats.hits AS hits,
		season_stats.doubles AS doubles,
		season_stats.triples AS triples,
		season_stats.home_runs AS home_runs,
		season_stats.runs AS runs,
		season_stats.rbi AS rbi,
		season_stats.walks AS walks,
		season_stats.strikeouts AS strikeouts,
		season_stats.stolen_bases AS stolen_bases,
		season_stats.batting_avg AS batting_avg,
		season_stats.obp AS obp,
		season_stats.slg AS slg,
		season_stats.ops AS ops,
		season_stats.babip AS babip,
		season_stats.wins AS wins,
		season_stats.losses AS losses,
		season_stats.era AS era,
		season_stats.whip AS whip,
		season_stats.innings_pitched AS innings_pitched,
		season_stats.hits_allowed AS hits_allowed,
		season_stats.walks_allowed AS walks_allowed,
		season_stats.home_runs_allowed AS home_runs_allowed,
		season_stats.strikeouts_per_9 AS strikeouts_per_9,
		season_stats.walks_per_9 AS walks_per_9,
		season_stats.hits_per_9 AS hits_per_9,
		season_stats.home_runs_per_9 AS home_runs_per_9,
		season_stats.strikeout_walk_ratio AS strikeout_walk_ratio,
		season_stats.strike_percentage AS strike_percentage,
		season_stats.expected_batting_avg AS expected_batting_avg,
		season_stats.expected_slugging AS expected_slugging,
		season_stats.expected_woba AS expected_woba,
		season_stats.expected_era AS expected_era,
		season_stats.barrel_pct AS barrel_pct,
		season_stats.hard_hit_pct AS hard_hit_pct,
		season_stats.avg_exit_velocity AS avg_exit_velocity,
		season_stats.avg_launch_angle AS avg_launch_angle,
		season_stats.sweet_spot_pct AS sweet_spot_pct
	`
}

func (r playgroundQueryRow) toResponse() playgroundQueryItem {
	stat := models.SeasonStat{
		Year:               r.Year,
		TeamID:             r.TeamID,
		TeamName:           r.TeamName,
		Age:                r.Age,
		ValueScore:         r.ValueScore,
		GamesPlayed:        r.GamesPlayed,
		GamesStarted:       r.GamesStarted,
		PlateAppearances:   r.PlateAppearances,
		AtBats:             r.AtBats,
		Hits:               r.Hits,
		Doubles:            r.Doubles,
		Triples:            r.Triples,
		HomeRuns:           r.HomeRuns,
		Runs:               r.Runs,
		RBI:                r.RBI,
		Walks:              r.Walks,
		Strikeouts:         r.Strikeouts,
		StolenBases:        r.StolenBases,
		BattingAvg:         r.BattingAvg,
		OBP:                r.OBP,
		SLG:                r.SLG,
		OPS:                r.OPS,
		BABIP:              r.BABIP,
		Wins:               r.Wins,
		Losses:             r.Losses,
		ERA:                r.ERA,
		WHIP:               r.WHIP,
		InningsPitched:     r.InningsPitched,
		HitsAllowed:        r.HitsAllowed,
		WalksAllowed:       r.WalksAllowed,
		HomeRunsAllowed:    r.HomeRunsAllowed,
		StrikeoutsPer9:     r.StrikeoutsPer9,
		WalksPer9:          r.WalksPer9,
		HitsPer9:           r.HitsPer9,
		HomeRunsPer9:       r.HomeRunsPer9,
		StrikeoutWalkRatio: r.StrikeoutWalkRatio,
		StrikePercentage:   r.StrikePercentage,
		ExpectedBattingAvg: r.ExpectedBattingAvg,
		ExpectedSlugging:   r.ExpectedSlugging,
		ExpectedWOBA:       r.ExpectedWOBA,
		ExpectedERA:        r.ExpectedERA,
		BarrelPct:          r.BarrelPct,
		HardHitPct:         r.HardHitPct,
		AvgExitVelocity:    r.AvgExitVelocity,
		AvgLaunchAngle:     r.AvgLaunchAngle,
		SweetSpotPct:       r.SweetSpotPct,
	}

	// Workload presence, not listed position, decides whether a row behaves like
	// a hitter, pitcher, or two-way season in the playground response.
	return playgroundQueryItem{
		Player: playgroundQueryPlayerItem{
			ID:        r.PlayerID,
			MLBID:     r.MLBID,
			FirstName: r.FirstName,
			LastName:  r.LastName,
			FullName:  strings.TrimSpace(r.FirstName + " " + r.LastName),
			Position:  r.Position,
			Bats:      r.Bats,
			Throws:    r.Throws,
			Active:    r.Active,
			ImageURL:  r.ImageURL,
		},
		Season: playgroundQuerySeasonItem{
			Year:       r.Year,
			TeamID:     r.TeamID,
			TeamName:   r.TeamName,
			Age:        r.Age,
			ValueScore: r.ValueScore,
		},
		Hitting:  newPlayerHittingStats(stat),
		Pitching: newPlayerPitchingStats(stat),
	}
}

func (p playgroundQueryParams) paginationOffset() int {
	if p.Offset > 0 {
		return p.Offset
	}

	return playerListOffset(playerListParams{Page: p.Page, PageSize: p.PageSize})
}

func parseOptionalIntQuery(c echo.Context, key string) (*int, error) {
	value := strings.TrimSpace(c.QueryParam(key))
	if value == "" {
		return nil, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid %s value", key))
	}
	return &parsed, nil
}

func parseOptionalFloatQuery(c echo.Context, key string) (*float64, error) {
	value := strings.TrimSpace(c.QueryParam(key))
	if value == "" {
		return nil, nil
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid %s value", key))
	}
	return &parsed, nil
}

func parseOptionalBoolQuery(c echo.Context, key string) (*bool, error) {
	value := strings.TrimSpace(c.QueryParam(key))
	if value == "" {
		return nil, nil
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid %s value", key))
	}
	return &parsed, nil
}

func hasPitchingThresholds(params playgroundQueryParams) bool {
	return params.MinERA != nil || params.MaxERA != nil ||
		params.MinWHIP != nil || params.MaxWHIP != nil ||
		params.MinK9 != nil || params.MaxK9 != nil
}

func hasPitchingWorkloadFilters(params playgroundQueryParams) bool {
	return params.MinIP != nil || params.MaxIP != nil
}

func hasHittingThresholds(params playgroundQueryParams) bool {
	return params.MinHR != nil || params.MaxHR != nil ||
		params.MinAvg != nil || params.MaxAvg != nil ||
		params.MinOBP != nil || params.MaxOBP != nil ||
		params.MinSLG != nil || params.MaxSLG != nil ||
		params.MinSB != nil || params.MaxSB != nil
}

func hasHittingWorkloadFilters(params playgroundQueryParams) bool {
	return params.MinPA != nil || params.MaxPA != nil
}

func playgroundQueryParamError(err error) string {
	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		if message, ok := httpErr.Message.(string); ok {
			return message
		}
	}
	return err.Error()
}

func isSupportedPlaygroundSort(sort string) bool {
	switch sort {
	case "name", "-name",
		"year", "-year",
		"age", "-age",
		"value_score", "-value_score",
		"home_runs", "-home_runs",
		"batting_avg", "-batting_avg",
		"obp", "-obp",
		"slg", "-slg",
		"stolen_bases", "-stolen_bases",
		"era", "-era",
		"whip", "-whip",
		"innings_pitched", "-innings_pitched",
		"strikeouts_per_9", "-strikeouts_per_9":
		return true
	default:
		return false
	}
}

func playgroundQueryOrderBy(sort string) []string {
	switch sort {
	case "name":
		return []string{"players.last_name ASC", "players.first_name ASC", "season_stats.year DESC", "season_stats.id ASC"}
	case "-name":
		return []string{"players.last_name DESC", "players.first_name DESC", "season_stats.year DESC", "season_stats.id ASC"}
	case "year":
		return []string{"season_stats.year ASC", "players.last_name ASC", "players.first_name ASC", "season_stats.id ASC"}
	case "-year":
		return []string{"season_stats.year DESC", "players.last_name ASC", "players.first_name ASC", "season_stats.id ASC"}
	case "age":
		return []string{"season_stats.age ASC", "players.last_name ASC", "players.first_name ASC", "season_stats.id ASC"}
	case "-age":
		return []string{"season_stats.age DESC", "players.last_name ASC", "players.first_name ASC", "season_stats.id ASC"}
	case "value_score":
		return []string{"season_stats.value_score ASC", "players.last_name ASC", "players.first_name ASC", "season_stats.id ASC"}
	case "-value_score":
		return []string{"season_stats.value_score DESC", "players.last_name ASC", "players.first_name ASC", "season_stats.id ASC"}
	case "home_runs":
		return []string{"season_stats.home_runs ASC", "players.last_name ASC", "players.first_name ASC", "season_stats.id ASC"}
	case "-home_runs":
		return []string{"season_stats.home_runs DESC", "players.last_name ASC", "players.first_name ASC", "season_stats.id ASC"}
	case "batting_avg":
		return []string{"season_stats.batting_avg ASC", "players.last_name ASC", "players.first_name ASC", "season_stats.id ASC"}
	case "-batting_avg":
		return []string{"season_stats.batting_avg DESC", "players.last_name ASC", "players.first_name ASC", "season_stats.id ASC"}
	case "obp":
		return []string{"season_stats.obp ASC", "players.last_name ASC", "players.first_name ASC", "season_stats.id ASC"}
	case "-obp":
		return []string{"season_stats.obp DESC", "players.last_name ASC", "players.first_name ASC", "season_stats.id ASC"}
	case "slg":
		return []string{"season_stats.slg ASC", "players.last_name ASC", "players.first_name ASC", "season_stats.id ASC"}
	case "-slg":
		return []string{"season_stats.slg DESC", "players.last_name ASC", "players.first_name ASC", "season_stats.id ASC"}
	case "stolen_bases":
		return []string{"season_stats.stolen_bases ASC", "players.last_name ASC", "players.first_name ASC", "season_stats.id ASC"}
	case "-stolen_bases":
		return []string{"season_stats.stolen_bases DESC", "players.last_name ASC", "players.first_name ASC", "season_stats.id ASC"}
	case "era":
		return []string{"season_stats.era ASC", "players.last_name ASC", "players.first_name ASC", "season_stats.id ASC"}
	case "-era":
		return []string{"season_stats.era DESC", "players.last_name ASC", "players.first_name ASC", "season_stats.id ASC"}
	case "whip":
		return []string{"season_stats.whip ASC", "players.last_name ASC", "players.first_name ASC", "season_stats.id ASC"}
	case "-whip":
		return []string{"season_stats.whip DESC", "players.last_name ASC", "players.first_name ASC", "season_stats.id ASC"}
	case "innings_pitched":
		return []string{"season_stats.innings_pitched ASC", "players.last_name ASC", "players.first_name ASC", "season_stats.id ASC"}
	case "-innings_pitched":
		return []string{"season_stats.innings_pitched DESC", "players.last_name ASC", "players.first_name ASC", "season_stats.id ASC"}
	case "strikeouts_per_9":
		return []string{"season_stats.strikeouts_per_9 ASC", "players.last_name ASC", "players.first_name ASC", "season_stats.id ASC"}
	case "-strikeouts_per_9":
		return []string{"season_stats.strikeouts_per_9 DESC", "players.last_name ASC", "players.first_name ASC", "season_stats.id ASC"}
	default:
		return []string{"season_stats.value_score DESC", "players.last_name ASC", "players.first_name ASC", "season_stats.id ASC"}
	}
}
