package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jsquardo/capcurve/internal/models"
	"github.com/jsquardo/capcurve/internal/syncjob"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestGetLeaderboardsEndpoint(t *testing.T) {
	db := testHandlersDB(t)
	nowSuffix := time.Now().UnixNano()
	defaultSeason := defaultLeaderboardSeason(time.Now())
	priorSeason := defaultSeason - 1
	customSeason := defaultSeason + 50

	peakLeader := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     914100000 + int(nowSuffix%100000),
		FirstName: fmt.Sprintf("CodexBoard%d", nowSuffix),
		LastName:  "PeakLeader",
		Position:  "CF",
		Active:    false,
	})
	createTestSeasonStat(t, db, peakLeader.ID, testSeasonFixture{
		Year:             priorSeason - 1,
		TeamID:           141,
		TeamName:         "Toronto Blue Jays",
		Age:              25,
		ValueScore:       72.3,
		PlateAppearances: 590,
		AtBats:           532,
		Hits:             161,
		HomeRuns:         27,
	})
	createTestSeasonStat(t, db, peakLeader.ID, testSeasonFixture{
		Year:             priorSeason,
		TeamID:           141,
		TeamName:         "Toronto Blue Jays",
		Age:              26,
		ValueScore:       94.2,
		PlateAppearances: 608,
		AtBats:           547,
		Hits:             176,
		HomeRuns:         34,
	})
	require.NoError(t, db.Create(&models.CareerArc{
		PlayerID:         int(peakLeader.ID),
		PeakYearStart:    priorSeason - 1,
		PeakYearEnd:      priorSeason,
		PeakWAR:          5.1,
		DeclineOnsetYear: defaultSeason + 1,
		ArcShape:         "Sustained Peak",
		LastComputedAt:   time.Now().UTC(),
	}).Error)

	peakWarStubLoser := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     914200000 + int(nowSuffix%100000),
		FirstName: fmt.Sprintf("CodexBoard%d", nowSuffix),
		LastName:  "PeakWar",
		Position:  "RF",
		Active:    false,
	})
	createTestSeasonStat(t, db, peakWarStubLoser.ID, testSeasonFixture{
		Year:             priorSeason - 1,
		TeamID:           111,
		TeamName:         "Boston Red Sox",
		Age:              28,
		ValueScore:       75.4,
		PlateAppearances: 601,
		AtBats:           542,
		Hits:             169,
		HomeRuns:         29,
	})
	createTestSeasonStat(t, db, peakWarStubLoser.ID, testSeasonFixture{
		Year:             defaultSeason,
		TeamID:           111,
		TeamName:         "Boston Red Sox",
		Age:              29,
		ValueScore:       68.0,
		PlateAppearances: 584,
		AtBats:           524,
		Hits:             153,
		HomeRuns:         22,
	})
	require.NoError(t, db.Create(&models.CareerArc{
		PlayerID:         int(peakWarStubLoser.ID),
		PeakYearStart:    priorSeason - 1,
		PeakYearEnd:      priorSeason - 1,
		PeakWAR:          9.9,
		DeclineOnsetYear: defaultSeason,
		ArcShape:         "Sharp Peak",
		LastComputedAt:   time.Now().UTC(),
	}).Error)

	hrLeader := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     914300000 + int(nowSuffix%100000),
		FirstName: fmt.Sprintf("CodexBoard%d", nowSuffix),
		LastName:  "Slugger",
		Position:  "LF",
		Active:    true,
	})
	createTestSeasonStat(t, db, hrLeader.ID, testSeasonFixture{
		Year:             defaultSeason,
		TeamID:           147,
		TeamName:         "New York Yankees",
		Age:              30,
		ValueScore:       81.0,
		PlateAppearances: 640,
		AtBats:           570,
		Hits:             172,
		HomeRuns:         88,
		BattingAvg:       0.302,
	})
	createTestSeasonStat(t, db, hrLeader.ID, testSeasonFixture{
		Year:             customSeason,
		TeamID:           147,
		TeamName:         "New York Yankees",
		Age:              32,
		ValueScore:       83.0,
		PlateAppearances: 641,
		AtBats:           573,
		Hits:             228,
		HomeRuns:         36,
		BattingAvg:       0.398,
	})

	hrRunnerUp := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     914400000 + int(nowSuffix%100000),
		FirstName: fmt.Sprintf("CodexBoard%d", nowSuffix),
		LastName:  "RunnerUp",
		Position:  "1B",
		Active:    true,
	})
	createTestSeasonStat(t, db, hrRunnerUp.ID, testSeasonFixture{
		Year:             defaultSeason,
		TeamID:           119,
		TeamName:         "Los Angeles Dodgers",
		Age:              31,
		ValueScore:       78.2,
		PlateAppearances: 632,
		AtBats:           564,
		Hits:             167,
		HomeRuns:         77,
		BattingAvg:       0.296,
	})
	createTestSeasonStat(t, db, hrRunnerUp.ID, testSeasonFixture{
		Year:             customSeason,
		TeamID:           119,
		TeamName:         "Los Angeles Dodgers",
		Age:              30,
		ValueScore:       84.0,
		PlateAppearances: 627,
		AtBats:           558,
		Hits:             190,
		HomeRuns:         28,
		BattingAvg:       0.411,
	})

	eraAce := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     914500000 + int(nowSuffix%100000),
		FirstName: fmt.Sprintf("CodexBoard%d", nowSuffix),
		LastName:  "EraAce",
		Position:  "SP",
		Active:    true,
	})
	createTestSeasonStat(t, db, eraAce.ID, testSeasonFixture{
		Year:           defaultSeason,
		TeamID:         138,
		TeamName:       "St. Louis Cardinals",
		Age:            28,
		ValueScore:     74.0,
		InningsPitched: 188.1,
		ERA:            0.91,
		StrikeoutsPer9: 10.1,
	})

	k9Ace := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     914600000 + int(nowSuffix%100000),
		FirstName: fmt.Sprintf("CodexBoard%d", nowSuffix),
		LastName:  "K9Ace",
		Position:  "SP",
		Active:    true,
	})
	createTestSeasonStat(t, db, k9Ace.ID, testSeasonFixture{
		Year:           defaultSeason,
		TeamID:         121,
		TeamName:       "New York Mets",
		Age:            27,
		ValueScore:     70.0,
		InningsPitched: 176.0,
		ERA:            3.22,
		StrikeoutsPer9: 19.7,
	})

	t.Cleanup(func() {
		cleanupTestPlayers(t, db, []int{
			peakLeader.MLBID,
			peakWarStubLoser.MLBID,
			hrLeader.MLBID,
			hrRunnerUp.MLBID,
			eraAce.MLBID,
			k9Ace.MLBID,
		})
	})

	t.Run("peak_arc ranks by derived peak value score and omits season", func(t *testing.T) {
		response := hitGetLeaderboardsEndpoint(t, db, "/api/v1/leaderboards?category=peak_arc", http.StatusOK)

		require.Equal(t, "peak_arc", response.Data.Category)
		require.EqualValues(t, 2, response.Data.Meta.Total)
		require.Equal(t, 1, response.Data.Meta.Page)
		require.Equal(t, 25, response.Data.Meta.PageSize)
		require.Equal(t, 1, response.Data.Meta.TotalPages)
		require.Len(t, response.Data.Leaders, 2)
		require.Equal(t, 1, response.Data.Leaders[0].Rank)
		require.Equal(t, peakLeader.ID, response.Data.Leaders[0].PlayerID)
		require.Equal(t, peakLeader.FirstName+" "+peakLeader.LastName, response.Data.Leaders[0].PlayerName)
		require.Equal(t, "Toronto Blue Jays", response.Data.Leaders[0].Team)
		require.Equal(t, 94.2, response.Data.Leaders[0].Value)
		require.Nil(t, response.Data.Leaders[0].Season)
		require.Equal(t, peakWarStubLoser.ID, response.Data.Leaders[1].PlayerID)
		require.Equal(t, 75.4, response.Data.Leaders[1].Value)
	})

	t.Run("season stat categories default to the most recent completed season and paginate ranks", func(t *testing.T) {
		response := hitGetLeaderboardsEndpoint(t, db, "/api/v1/leaderboards?category=hr&page=2&page_size=1", http.StatusOK)

		require.Equal(t, "hr", response.Data.Category)
		require.GreaterOrEqual(t, response.Data.Meta.Total, int64(2))
		require.Equal(t, 2, response.Data.Meta.Page)
		require.Equal(t, 1, response.Data.Meta.PageSize)
		require.GreaterOrEqual(t, response.Data.Meta.TotalPages, 2)
		require.Len(t, response.Data.Leaders, 1)
		require.Equal(t, 2, response.Data.Leaders[0].Rank)
		require.Equal(t, hrRunnerUp.ID, response.Data.Leaders[0].PlayerID)
		require.NotNil(t, response.Data.Leaders[0].Season)
		require.Equal(t, defaultSeason, *response.Data.Leaders[0].Season)
		require.Equal(t, 77.0, response.Data.Leaders[0].Value)
	})

	t.Run("explicit season works for avg", func(t *testing.T) {
		response := hitGetLeaderboardsEndpoint(t, db, fmt.Sprintf("/api/v1/leaderboards?category=avg&season=%d", customSeason), http.StatusOK)

		require.Equal(t, "avg", response.Data.Category)
		require.Len(t, response.Data.Leaders, 2)
		require.Equal(t, hrRunnerUp.ID, response.Data.Leaders[0].PlayerID)
		require.NotNil(t, response.Data.Leaders[0].Season)
		require.Equal(t, customSeason, *response.Data.Leaders[0].Season)
		require.Equal(t, 0.411, response.Data.Leaders[0].Value)
		require.Equal(t, hrLeader.ID, response.Data.Leaders[1].PlayerID)
		require.Equal(t, 0.398, response.Data.Leaders[1].Value)
	})

	t.Run("era sorts ascending and k9 sorts descending", func(t *testing.T) {
		eraResponse := hitGetLeaderboardsEndpoint(t, db, "/api/v1/leaderboards?category=era", http.StatusOK)
		k9Response := hitGetLeaderboardsEndpoint(t, db, "/api/v1/leaderboards?category=k9", http.StatusOK)

		require.Equal(t, eraAce.ID, eraResponse.Data.Leaders[0].PlayerID)
		require.Equal(t, 0.91, eraResponse.Data.Leaders[0].Value)
		require.Equal(t, k9Ace.ID, k9Response.Data.Leaders[0].PlayerID)
		require.Equal(t, 19.7, k9Response.Data.Leaders[0].Value)
	})

	t.Run("returns 400 for unsupported categories", func(t *testing.T) {
		rec := hitLeaderboardsEndpointRaw(t, db, "/api/v1/leaderboards?category=most_overpaid")

		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.JSONEq(t, `{"error":"invalid category value"}`, rec.Body.String())
	})
}

func TestRemovedDeferredContractRoutes(t *testing.T) {
	db := testHandlersDB(t)
	e := echo.New()
	RegisterRoutes(e, db, syncjob.NewStatusStore(false), "super-secret")

	for _, path := range []string{
		"/api/v1/players/1/contracts",
		"/api/v1/contracts/1",
		"/api/v1/leaderboards/most-overpaid",
		"/api/v1/leaderboards/best-value",
		"/api/v1/leaderboards/peak-arcs",
	} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		require.Equalf(t, http.StatusNotFound, rec.Code, "expected %s to be removed", path)
	}
}

type leaderboardTestResponse struct {
	Data leaderboardTestData `json:"data"`
}

type leaderboardTestData struct {
	Category string                `json:"category"`
	Leaders  []leaderboardTestItem `json:"leaders"`
	Meta     leaderboardTestMeta   `json:"meta"`
}

type leaderboardTestMeta struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

type leaderboardTestItem struct {
	Rank       int     `json:"rank"`
	PlayerID   uint    `json:"player_id"`
	PlayerName string  `json:"player_name"`
	Position   string  `json:"position"`
	Team       string  `json:"team"`
	Value      float64 `json:"value"`
	Season     *int    `json:"season,omitempty"`
}

func hitGetLeaderboardsEndpoint(t *testing.T, db *gorm.DB, path string, expectedStatus int) leaderboardTestResponse {
	t.Helper()

	rec := hitLeaderboardsEndpointRaw(t, db, path)
	require.Equal(t, expectedStatus, rec.Code)

	var response leaderboardTestResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))

	return response
}

func hitLeaderboardsEndpointRaw(t *testing.T, db *gorm.DB, path string) *httptest.ResponseRecorder {
	t.Helper()

	e := echo.New()
	RegisterRoutes(e, db, syncjob.NewStatusStore(false), "super-secret")

	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	return rec
}
