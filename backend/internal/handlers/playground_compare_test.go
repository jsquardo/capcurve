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

type playgroundCompareTestResponse struct {
	Data []playgroundCompareItem `json:"data"`
}

func TestGetPlaygroundCompareEndpoint(t *testing.T) {
	db := testHandlersDB(t)
	nowSuffix := time.Now().UnixNano()

	first := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     915100000 + int(nowSuffix%100000),
		FirstName: fmt.Sprintf("Compare%d", nowSuffix),
		LastName:  "First",
		Position:  "OF",
		Active:    true,
	})
	createTestSeasonStat(t, db, first.ID, testSeasonFixture{
		Year:             2023,
		TeamID:           147,
		TeamName:         "New York Yankees",
		Age:              28,
		ValueScore:       70.5,
		GamesPlayed:      145,
		PlateAppearances: 640,
		AtBats:           570,
		Hits:             177,
		HomeRuns:         28,
		BattingAvg:       0.311,
		OBP:              0.389,
		SLG:              0.541,
		OPS:              0.930,
	})
	createTestSeasonStat(t, db, first.ID, testSeasonFixture{
		Year:             2024,
		TeamID:           147,
		TeamName:         "New York Yankees",
		Age:              29,
		ValueScore:       82.1,
		GamesPlayed:      150,
		PlateAppearances: 666,
		AtBats:           594,
		Hits:             184,
		HomeRuns:         39,
		BattingAvg:       0.310,
		OBP:              0.401,
		SLG:              0.590,
		OPS:              0.991,
	})

	second := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     915200000 + int(nowSuffix%100000),
		FirstName: fmt.Sprintf("Compare%d", nowSuffix),
		LastName:  "Second",
		Position:  "SP",
		Active:    false,
	})
	createTestSeasonStat(t, db, second.ID, testSeasonFixture{
		Year:               2022,
		TeamID:             121,
		TeamName:           "New York Mets",
		Age:                31,
		ValueScore:         76.4,
		GamesPlayed:        32,
		GamesStarted:       32,
		Wins:               17,
		Losses:             8,
		ERA:                2.94,
		WHIP:               1.06,
		InningsPitched:     196.1,
		HitsAllowed:        151,
		WalksAllowed:       39,
		HomeRunsAllowed:    18,
		StrikeoutsPer9:     10.1,
		WalksPer9:          1.8,
		HitsPer9:           6.9,
		HomeRunsPer9:       0.8,
		StrikeoutWalkRatio: 5.61,
		StrikePercentage:   68.3,
	})

	third := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     915300000 + int(nowSuffix%100000),
		FirstName: fmt.Sprintf("Compare%d", nowSuffix),
		LastName:  "Third",
		Position:  "SS",
		Active:    true,
	})
	createTestSeasonStat(t, db, third.ID, testSeasonFixture{
		Year:             2021,
		TeamID:           144,
		TeamName:         "Atlanta Braves",
		Age:              24,
		ValueScore:       55.0,
		GamesPlayed:      80,
		PlateAppearances: 280,
		AtBats:           250,
		Hits:             70,
		HomeRuns:         8,
		BattingAvg:       0.280,
		OBP:              0.340,
		SLG:              0.430,
		OPS:              0.770,
	})

	t.Cleanup(func() {
		cleanupTestPlayers(t, db, []int{first.MLBID, second.MLBID, third.MLBID})
	})

	t.Run("returns requested players in request order with filtered seasons", func(t *testing.T) {
		path := fmt.Sprintf(
			"/api/v1/playground/compare?player_ids=%d,%d&era_start=2023",
			second.ID,
			first.ID,
		)

		response := hitPlaygroundCompareEndpoint(t, db, path, http.StatusOK)

		require.Len(t, response.Data, 2)
		require.Equal(t, second.ID, response.Data[0].Player.ID)
		require.Empty(t, response.Data[0].Seasons)

		require.Equal(t, first.ID, response.Data[1].Player.ID)
		require.Len(t, response.Data[1].Seasons, 2)
		require.Equal(t, 2023, response.Data[1].Seasons[0].Season.Year)
		require.Equal(t, 2024, response.Data[1].Seasons[1].Season.Year)
		require.NotNil(t, response.Data[1].Seasons[0].Hitting)
		require.Nil(t, response.Data[1].Seasons[0].Pitching)
	})

	t.Run("filters pitching group rows by workload", func(t *testing.T) {
		path := fmt.Sprintf(
			"/api/v1/playground/compare?player_ids=%d,%d&group=pitching&min_ip=190",
			first.ID,
			second.ID,
		)

		response := hitPlaygroundCompareEndpoint(t, db, path, http.StatusOK)

		require.Len(t, response.Data, 2)
		require.Equal(t, first.ID, response.Data[0].Player.ID)
		require.Empty(t, response.Data[0].Seasons)
		require.Equal(t, second.ID, response.Data[1].Player.ID)
		require.Len(t, response.Data[1].Seasons, 1)
		require.NotNil(t, response.Data[1].Seasons[0].Pitching)
		require.Nil(t, response.Data[1].Seasons[0].Hitting)
		require.Equal(t, 196.1, response.Data[1].Seasons[0].Pitching.InningsPitched)
	})

	t.Run("rejects too few player ids", func(t *testing.T) {
		rec := hitPlaygroundCompareEndpointRaw(t, db, fmt.Sprintf("/api/v1/playground/compare?player_ids=%d", first.ID))

		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.JSONEq(t, `{"error":"player_ids must include between 2 and 4 players"}`, rec.Body.String())
	})

	t.Run("rejects missing players without partial results", func(t *testing.T) {
		rec := hitPlaygroundCompareEndpointRaw(
			t,
			db,
			fmt.Sprintf("/api/v1/playground/compare?player_ids=%d,999999999", first.ID),
		)

		require.Equal(t, http.StatusNotFound, rec.Code)
		require.JSONEq(t, `{"error":"player not found: 999999999"}`, rec.Body.String())
	})

	t.Run("rejects hitting workload filters for group pitching", func(t *testing.T) {
		rec := hitPlaygroundCompareEndpointRaw(t, db, fmt.Sprintf("/api/v1/playground/compare?player_ids=%d,%d&group=pitching&min_pa=50", first.ID, second.ID))

		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.JSONEq(t, `{"error":"hitting thresholds are not supported for group=pitching"}`, rec.Body.String())
	})
}

func hitPlaygroundCompareEndpoint(t *testing.T, db *gorm.DB, path string, expectedStatus int) playgroundCompareTestResponse {
	t.Helper()

	e := echo.New()
	RegisterRoutes(e, db, syncjob.NewStatusStore(false), "super-secret")

	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	require.Equal(t, expectedStatus, rec.Code)

	var response playgroundCompareTestResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))

	return response
}

func hitPlaygroundCompareEndpointRaw(t *testing.T, db *gorm.DB, path string) *httptest.ResponseRecorder {
	t.Helper()

	e := echo.New()
	RegisterRoutes(e, db, syncjob.NewStatusStore(false), "super-secret")

	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	return rec
}
