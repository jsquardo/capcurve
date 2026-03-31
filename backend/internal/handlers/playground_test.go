package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jsquardo/capcurve/internal/syncjob"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type playgroundQueryTestResponse struct {
	Data []playgroundQueryItem `json:"data"`
	Meta playerListMeta        `json:"meta"`
}

func TestGetPlaygroundQueryEndpoint(t *testing.T) {
	db := testHandlersDB(t)
	nowSuffix := time.Now().UnixNano()

	hitterPrefix := fmt.Sprintf("CodexPlayH%d", nowSuffix)
	hitterOne := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     914100000 + int(nowSuffix%100000),
		FirstName: hitterPrefix,
		LastName:  "Slugger",
		Position:  "OF",
		Active:    true,
	})
	createTestSeasonStat(t, db, hitterOne.ID, testSeasonFixture{
		Year:             2024,
		TeamID:           147,
		TeamName:         "New York Yankees",
		Age:              29,
		ValueScore:       82.4,
		GamesPlayed:      151,
		PlateAppearances: 665,
		AtBats:           590,
		Hits:             182,
		Doubles:          34,
		HomeRuns:         41,
		Runs:             109,
		RBI:              117,
		Walks:            66,
		Strikeouts:       149,
		StolenBases:      13,
		BattingAvg:       0.308,
		OBP:              0.394,
		SLG:              0.582,
		OPS:              0.976,
		BABIP:            0.336,
	})

	hitterTwo := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     914200000 + int(nowSuffix%100000),
		FirstName: hitterPrefix,
		LastName:  "Runner",
		Position:  "SS",
		Active:    true,
	})
	createTestSeasonStat(t, db, hitterTwo.ID, testSeasonFixture{
		Year:             2023,
		TeamID:           147,
		TeamName:         "New York Yankees",
		Age:              27,
		ValueScore:       71.1,
		GamesPlayed:      149,
		PlateAppearances: 641,
		AtBats:           575,
		Hits:             171,
		Doubles:          28,
		HomeRuns:         23,
		Runs:             97,
		RBI:              74,
		Walks:            49,
		Strikeouts:       131,
		StolenBases:      29,
		BattingAvg:       0.297,
		OBP:              0.361,
		SLG:              0.482,
		OPS:              0.843,
		BABIP:            0.329,
	})

	pitcherPrefix := fmt.Sprintf("CodexPlayP%d", nowSuffix)
	pitcherOne := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     914300000 + int(nowSuffix%100000),
		FirstName: pitcherPrefix,
		LastName:  "Ace",
		Position:  "SP",
		Active:    false,
	})
	createTestSeasonStat(t, db, pitcherOne.ID, testSeasonFixture{
		Year:               2022,
		TeamID:             121,
		TeamName:           "New York Mets",
		Age:                31,
		ValueScore:         78.9,
		GamesPlayed:        32,
		GamesStarted:       32,
		Wins:               18,
		Losses:             7,
		ERA:                2.88,
		WHIP:               1.03,
		InningsPitched:     201.2,
		HitsAllowed:        156,
		WalksAllowed:       41,
		HomeRunsAllowed:    17,
		StrikeoutsPer9:     10.7,
		WalksPer9:          1.8,
		HitsPer9:           7.0,
		HomeRunsPer9:       0.8,
		StrikeoutWalkRatio: 5.96,
		StrikePercentage:   68.8,
	})

	pitcherTwo := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     914400000 + int(nowSuffix%100000),
		FirstName: pitcherPrefix,
		LastName:  "Bulk",
		Position:  "RP",
		Active:    true,
	})
	createTestSeasonStat(t, db, pitcherTwo.ID, testSeasonFixture{
		Year:               2024,
		TeamID:             121,
		TeamName:           "New York Mets",
		Age:                28,
		ValueScore:         49.2,
		GamesPlayed:        62,
		GamesStarted:       0,
		ERA:                4.21,
		WHIP:               1.31,
		InningsPitched:     68.0,
		HitsAllowed:        57,
		WalksAllowed:       23,
		HomeRunsAllowed:    9,
		StrikeoutsPer9:     8.9,
		WalksPer9:          3.0,
		HitsPer9:           7.5,
		HomeRunsPer9:       1.2,
		StrikeoutWalkRatio: 2.96,
		StrikePercentage:   64.2,
	})

	pitcherThree := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     914500000 + int(nowSuffix%100000),
		FirstName: pitcherPrefix,
		LastName:  "Control",
		Position:  "SP",
		Active:    false,
	})
	createTestSeasonStat(t, db, pitcherThree.ID, testSeasonFixture{
		Year:               2021,
		TeamID:             121,
		TeamName:           "New York Mets",
		Age:                30,
		ValueScore:         76.3,
		GamesPlayed:        31,
		GamesStarted:       31,
		Wins:               16,
		Losses:             8,
		ERA:                2.95,
		WHIP:               1.05,
		InningsPitched:     189.1,
		HitsAllowed:        149,
		WalksAllowed:       36,
		HomeRunsAllowed:    15,
		StrikeoutsPer9:     10.2,
		WalksPer9:          1.7,
		HitsPer9:           7.1,
		HomeRunsPer9:       0.7,
		StrikeoutWalkRatio: 6.00,
		StrikePercentage:   69.1,
	})

	t.Cleanup(func() {
		cleanupTestPlayers(t, db, []int{
			hitterOne.MLBID,
			hitterTwo.MLBID,
			pitcherOne.MLBID,
			pitcherTwo.MLBID,
			pitcherThree.MLBID,
		})
	})

	t.Run("returns filtered hitter season rows with approved response shape and pagination", func(t *testing.T) {
		path := fmt.Sprintf(
			"/api/v1/playground/query?q=%s&group=hitting&team=147&era_start=2023&era_end=2024&age_min=28&age_max=30&min_pa=600&min_hr=30&min_obp=0.390&min_slg=0.580&position=OF,SS&sort=-value_score&page=1&page_size=1",
			hitterPrefix,
		)
		response := hitPlaygroundQueryEndpoint(t, db, path, http.StatusOK)

		require.EqualValues(t, 1, response.Meta.Total)
		require.Equal(t, 1, response.Meta.Page)
		require.Equal(t, 1, response.Meta.PageSize)
		require.Equal(t, 1, response.Meta.TotalPages)
		require.Len(t, response.Data, 1)

		row := response.Data[0]
		require.Equal(t, hitterOne.ID, row.Player.ID)
		require.Equal(t, hitterOne.MLBID, row.Player.MLBID)
		require.Equal(t, hitterOne.FirstName+" "+hitterOne.LastName, row.Player.FullName)
		require.Equal(t, 2024, row.Season.Year)
		require.Equal(t, "New York Yankees", row.Season.TeamName)
		require.Equal(t, 82.4, row.Season.ValueScore)
		require.NotNil(t, row.Hitting)
		require.Nil(t, row.Pitching)
		require.Equal(t, 41, row.Hitting.HomeRuns)
		require.Equal(t, 665, row.Hitting.PlateAppearances)
	})

	t.Run("returns filtered pitcher rows and offset pagination compatibility", func(t *testing.T) {
		path := fmt.Sprintf(
			"/api/v1/playground/query?q=%s&group=pitching&team=Mets&active=false&season=2022&min_ip=150&max_era=3.00&min_k9=10.0&sort=era&limit=1&offset=0",
			pitcherPrefix,
		)
		response := hitPlaygroundQueryEndpoint(t, db, path, http.StatusOK)

		require.EqualValues(t, 1, response.Meta.Total)
		require.Equal(t, 1, response.Meta.Page)
		require.Equal(t, 1, response.Meta.PageSize)
		require.Equal(t, 1, response.Meta.TotalPages)
		require.Len(t, response.Data, 1)

		row := response.Data[0]
		require.Equal(t, pitcherOne.ID, row.Player.ID)
		require.Equal(t, 2022, row.Season.Year)
		require.Nil(t, row.Hitting)
		require.NotNil(t, row.Pitching)
		require.Equal(t, 2.88, row.Pitching.ERA)
		require.Equal(t, 10.7, row.Pitching.StrikeoutsPer9)
	})

	t.Run("uses the caller offset without snapping to page boundaries", func(t *testing.T) {
		path := fmt.Sprintf(
			"/api/v1/playground/query?q=%s&group=pitching&team=Mets&active=false&min_ip=150&max_era=3.00&min_k9=10.0&sort=era&limit=2&offset=1",
			pitcherPrefix,
		)
		response := hitPlaygroundQueryEndpoint(t, db, path, http.StatusOK)

		require.EqualValues(t, 2, response.Meta.Total)
		require.Equal(t, 1, response.Meta.Page)
		require.Equal(t, 2, response.Meta.PageSize)
		require.Equal(t, 1, response.Meta.TotalPages)
		require.Len(t, response.Data, 1)

		row := response.Data[0]
		require.Equal(t, pitcherThree.ID, row.Player.ID)
		require.Equal(t, 2021, row.Season.Year)
		require.Nil(t, row.Hitting)
		require.NotNil(t, row.Pitching)
		require.Equal(t, 2.95, row.Pitching.ERA)
		require.Equal(t, 10.2, row.Pitching.StrikeoutsPer9)
	})

	t.Run("group all pitching thresholds ignore zero workload hitter rows", func(t *testing.T) {
		path := fmt.Sprintf(
			"/api/v1/playground/query?q=%s&group=all&team=147&max_era=1.00",
			hitterPrefix,
		)
		response := hitPlaygroundQueryEndpoint(t, db, path, http.StatusOK)

		require.EqualValues(t, 0, response.Meta.Total)
		require.Equal(t, 1, response.Meta.Page)
		require.Equal(t, 25, response.Meta.PageSize)
		require.Equal(t, 0, response.Meta.TotalPages)
		require.Empty(t, response.Data)
	})

	t.Run("group all hitting thresholds ignore zero workload pitcher rows", func(t *testing.T) {
		path := fmt.Sprintf(
			"/api/v1/playground/query?q=%s&group=all&team=Mets&active=false&max_obp=0.050",
			pitcherPrefix,
		)
		response := hitPlaygroundQueryEndpoint(t, db, path, http.StatusOK)

		require.EqualValues(t, 0, response.Meta.Total)
		require.Equal(t, 1, response.Meta.Page)
		require.Equal(t, 25, response.Meta.PageSize)
		require.Equal(t, 0, response.Meta.TotalPages)
		require.Empty(t, response.Data)
	})

	t.Run("group all pitching workload filters ignore zero workload hitter rows", func(t *testing.T) {
		path := fmt.Sprintf(
			"/api/v1/playground/query?q=%s&group=all&team=147&max_ip=1.00",
			hitterPrefix,
		)
		response := hitPlaygroundQueryEndpoint(t, db, path, http.StatusOK)

		require.EqualValues(t, 0, response.Meta.Total)
		require.Equal(t, 1, response.Meta.Page)
		require.Equal(t, 25, response.Meta.PageSize)
		require.Equal(t, 0, response.Meta.TotalPages)
		require.Empty(t, response.Data)
	})

	t.Run("group all hitting workload filters ignore zero workload pitcher rows", func(t *testing.T) {
		path := fmt.Sprintf(
			"/api/v1/playground/query?q=%s&group=all&team=Mets&active=false&max_pa=5",
			pitcherPrefix,
		)
		response := hitPlaygroundQueryEndpoint(t, db, path, http.StatusOK)

		require.EqualValues(t, 0, response.Meta.Total)
		require.Equal(t, 1, response.Meta.Page)
		require.Equal(t, 25, response.Meta.PageSize)
		require.Equal(t, 0, response.Meta.TotalPages)
		require.Empty(t, response.Data)
	})

	t.Run("rejects pitching thresholds for group hitting", func(t *testing.T) {
		rec := hitPlaygroundQueryEndpointRaw(t, db, "/api/v1/playground/query?group=hitting&min_era=3.00")

		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.JSONEq(t, `{"error":"pitching thresholds are not supported for group=hitting"}`, rec.Body.String())
	})

	t.Run("rejects season combined with era range", func(t *testing.T) {
		rec := hitPlaygroundQueryEndpointRaw(t, db, "/api/v1/playground/query?season=2024&era_start=2023")

		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.JSONEq(t, `{"error":"season cannot be combined with era_start or era_end"}`, rec.Body.String())
	})
}

func hitPlaygroundQueryEndpoint(t *testing.T, db *gorm.DB, path string, expectedStatus int) playgroundQueryTestResponse {
	t.Helper()

	e := echo.New()
	RegisterRoutes(e, db, syncjob.NewStatusStore(false), "super-secret")

	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	require.Equal(t, expectedStatus, rec.Code)

	var response playgroundQueryTestResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))

	return response
}

func hitPlaygroundQueryEndpointRaw(t *testing.T, db *gorm.DB, path string) *httptest.ResponseRecorder {
	t.Helper()

	e := echo.New()
	RegisterRoutes(e, db, syncjob.NewStatusStore(false), "super-secret")

	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	return rec
}
