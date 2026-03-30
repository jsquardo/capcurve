package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jsquardo/capcurve/internal/models"
	"github.com/jsquardo/capcurve/internal/projection"
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

func newCareerArcData(player models.Player, seasonStats []models.SeasonStat, arc *models.CareerArc) careerArcData {
	timeline := make([]careerArcTimelineItem, 0, len(seasonStats))
	for _, stat := range seasonStats {
		timeline = append(timeline, newCareerArcTimelineItem(stat, arc))
	}

	return careerArcData{
		Player:     newCareerArcPlayerItem(player),
		Arc:        newCareerArcMetadata(arc, timeline),
		Timeline:   timeline,
		Projection: newCareerArcProjection(player.Active),
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

func newCareerArcProjection(active bool) careerArcProjection {
	projection := careerArcProjection{
		Status:         "pending",
		Eligible:       active,
		Points:         []careerArcProjectionPoint{},
		ConfidenceBand: []careerArcConfidenceBand{},
		Comparables:    []careerArcComparablePlayer{},
	}

	if active {
		projection.Reason = "projection engine not implemented yet"
		return projection
	}

	projection.Reason = "projections are only available for active players"
	return projection
}

func newProjectionData(player models.Player, projection careerArcProjection) projectionData {
	return projectionData{
		Player:     newCareerArcPlayerItem(player),
		Projection: projection,
	}
}

func newProjectionPayload(result projection.Result) careerArcProjection {
	payload := careerArcProjection{
		Status:         result.Status,
		Eligible:       result.Eligible,
		Reason:         result.Reason,
		Points:         make([]careerArcProjectionPoint, 0, len(result.Points)),
		ConfidenceBand: make([]careerArcConfidenceBand, 0, len(result.ConfidenceBand)),
		Comparables:    make([]careerArcComparablePlayer, 0, len(result.Comparables)),
	}

	for _, point := range result.Points {
		payload.Points = append(payload.Points, careerArcProjectionPoint{
			Year:         point.Year,
			Age:          point.Age,
			ValueScore:   point.ValueScore,
			IsProjection: true,
		})
	}

	for _, band := range result.ConfidenceBand {
		payload.ConfidenceBand = append(payload.ConfidenceBand, careerArcConfidenceBand{
			Year:  band.Year,
			Lower: band.Lower,
			Upper: band.Upper,
		})
	}

	for _, comparable := range result.Comparables {
		payload.Comparables = append(payload.Comparables, careerArcComparablePlayer{
			PlayerID: comparable.PlayerID,
			MLBID:    comparable.MLBID,
			FullName: comparable.FullName,
			Position: comparable.Position,
		})
	}

	return payload
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

func (h *Handler) GetCareerArc(c echo.Context) error {
	playerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid player id"})
	}

	var player models.Player
	if err := h.db.First(&player, playerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "player not found"})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var stats []models.SeasonStat
	if err := h.db.Where("player_id = ?", player.ID).Order("year ASC, id ASC").Find(&stats).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var arc *models.CareerArc
	var arcRecord models.CareerArc
	if err := h.db.Where("player_id = ?", player.ID).First(&arcRecord).Error; err != nil {
		// Missing arc metadata should not block the chart endpoint from returning
		// the player's historical timeline.
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	} else {
		arc = &arcRecord
	}

	return c.JSON(http.StatusOK, careerArcResponse{
		Data: newCareerArcData(player, stats, arc),
	})
}

func (h *Handler) GetProjection(c echo.Context) error {
	playerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid player id"})
	}

	var player models.Player
	if err := h.db.First(&player, playerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "player not found"})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var history []models.SeasonStat
	if err := h.db.Where("player_id = ?", player.ID).Order("year ASC, id ASC").Find(&history).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	service := projection.NewService()
	result := projection.Result{}

	if !player.Active || len(history) == 0 {
		result = service.Build(player, history, nil, nil)
	} else {
		var candidates []models.Player
		if err := h.db.Where("id <> ?", player.ID).Find(&candidates).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		var candidateStats []models.SeasonStat
		if err := h.db.Where("player_id <> ?", player.ID).Order("player_id ASC, year ASC, id ASC").Find(&candidateStats).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		result = service.Build(player, history, candidates, candidateStats)
	}

	return c.JSON(http.StatusOK, projectionResponse{
		Data: newProjectionData(player, newProjectionPayload(result)),
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
