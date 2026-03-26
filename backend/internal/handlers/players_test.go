package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/jsquardo/capcurve/internal/database"
	"github.com/jsquardo/capcurve/internal/models"
	"github.com/jsquardo/capcurve/internal/syncjob"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestListPlayersEndpoint(t *testing.T) {
	db := testHandlersDB(t)
	nowSuffix := time.Now().UnixNano()

	searchA := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     910000000 + int(nowSuffix%100000),
		FirstName: fmt.Sprintf("CodexQuery%d", nowSuffix),
		LastName:  "Alpha",
		Position:  "OF",
		Active:    true,
		Season: &testSeasonFixture{
			Year:       2024,
			TeamID:     147,
			TeamName:   "New York Yankees",
			Age:        28,
			ValueScore: 61.5,
		},
	})
	searchB := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     910100000 + int(nowSuffix%100000),
		FirstName: fmt.Sprintf("CodexQuery%d", nowSuffix),
		LastName:  "Zulu",
		Position:  "SS",
		Active:    true,
	})

	activePrefix := fmt.Sprintf("CodexActive%d", nowSuffix)
	_ = createTestPlayer(t, db, testPlayerFixture{
		MLBID:     910200000 + int(nowSuffix%100000),
		FirstName: activePrefix,
		LastName:  "Yes",
		Position:  "C",
		Active:    true,
		Season: &testSeasonFixture{
			Year:       2023,
			TeamID:     121,
			TeamName:   "New York Mets",
			Age:        27,
			ValueScore: 45.2,
		},
	})
	_ = createTestPlayer(t, db, testPlayerFixture{
		MLBID:     910300000 + int(nowSuffix%100000),
		FirstName: activePrefix,
		LastName:  "No",
		Position:  "1B",
		Active:    false,
	})

	sortPrefix := fmt.Sprintf("CodexSort%d", nowSuffix)
	_ = createTestPlayer(t, db, testPlayerFixture{
		MLBID:     910400000 + int(nowSuffix%100000),
		FirstName: sortPrefix,
		LastName:  "Low",
		Position:  "DH",
		Active:    true,
		Season: &testSeasonFixture{
			Year:       2024,
			TeamID:     119,
			TeamName:   "Los Angeles Dodgers",
			Age:        30,
			ValueScore: 12.5,
		},
	})
	_ = createTestPlayer(t, db, testPlayerFixture{
		MLBID:     910500000 + int(nowSuffix%100000),
		FirstName: sortPrefix,
		LastName:  "High",
		Position:  "DH",
		Active:    true,
		Season: &testSeasonFixture{
			Year:       2025,
			TeamID:     111,
			TeamName:   "Boston Red Sox",
			Age:        31,
			ValueScore: 88.8,
		},
	})

	teamPrefix := fmt.Sprintf("CodexTeam%d", nowSuffix)
	teamByName := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     910600000 + int(nowSuffix%100000),
		FirstName: teamPrefix,
		LastName:  "NameMatch",
		Position:  "OF",
		Active:    true,
		Season: &testSeasonFixture{
			Year:       2025,
			TeamID:     147,
			TeamName:   "New York Yankees",
			Age:        29,
			ValueScore: 54.3,
		},
	})
	teamByID := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     910700000 + int(nowSuffix%100000),
		FirstName: teamPrefix,
		LastName:  "IDMatch",
		Position:  "2B",
		Active:    true,
		Season: &testSeasonFixture{
			Year:       2025,
			TeamID:     147,
			TeamName:   "NY Snapshot Alias",
			Age:        27,
			ValueScore: 49.1,
		},
	})
	_ = createTestPlayer(t, db, testPlayerFixture{
		MLBID:     910800000 + int(nowSuffix%100000),
		FirstName: teamPrefix,
		LastName:  "Miss",
		Position:  "3B",
		Active:    true,
		Season: &testSeasonFixture{
			Year:       2025,
			TeamID:     121,
			TeamName:   "New York Mets",
			Age:        26,
			ValueScore: 47.7,
		},
	})

	seasonPrefix := fmt.Sprintf("CodexSeason%d", nowSuffix)
	seasonScoped := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     910900000 + int(nowSuffix%100000),
		FirstName: seasonPrefix,
		LastName:  "Switch",
		Position:  "SS",
		Active:    true,
		Season: &testSeasonFixture{
			Year:       2025,
			TeamID:     119,
			TeamName:   "Los Angeles Dodgers",
			Age:        31,
			ValueScore: 90.5,
		},
	})
	createTestSeasonStat(t, db, seasonScoped.ID, testSeasonFixture{
		Year:       2024,
		TeamID:     111,
		TeamName:   "Boston Red Sox",
		Age:        30,
		ValueScore: 70.2,
	})
	_ = createTestPlayer(t, db, testPlayerFixture{
		MLBID:     911000000 + int(nowSuffix%100000),
		FirstName: seasonPrefix,
		LastName:  "OtherYear",
		Position:  "1B",
		Active:    true,
		Season: &testSeasonFixture{
			Year:       2025,
			TeamID:     138,
			TeamName:   "St. Louis Cardinals",
			Age:        28,
			ValueScore: 55.5,
		},
	})

	t.Cleanup(func() {
		cleanupTestPlayers(t, db, []int{
			searchA.MLBID,
			searchB.MLBID,
			910200000 + int(nowSuffix%100000),
			910300000 + int(nowSuffix%100000),
			910400000 + int(nowSuffix%100000),
			910500000 + int(nowSuffix%100000),
			teamByName.MLBID,
			teamByID.MLBID,
			910800000 + int(nowSuffix%100000),
			seasonScoped.MLBID,
			911000000 + int(nowSuffix%100000),
		})
	})

	t.Run("q search returns compact list data and null latest season when absent", func(t *testing.T) {
		response := hitListPlayersEndpoint(t, db, "/api/v1/players?q="+searchA.FirstName)

		require.Equal(t, 2, response.Meta.Count)
		require.Equal(t, 2, len(response.Data))
		require.Equal(t, "Alpha", response.Data[0].LastName)
		require.Equal(t, searchA.FirstName+" Alpha", response.Data[0].FullName)
		require.NotNil(t, response.Data[0].LatestSeason)
		require.Equal(t, 2024, response.Data[0].LatestSeason.Year)
		require.Equal(t, "Zulu", response.Data[1].LastName)
		require.Nil(t, response.Data[1].LatestSeason)
	})

	t.Run("active filter limits the list", func(t *testing.T) {
		response := hitListPlayersEndpoint(t, db, "/api/v1/players?q="+activePrefix+"&active=true")

		require.Equal(t, 1, response.Meta.Count)
		require.Len(t, response.Data, 1)
		require.True(t, response.Data[0].Active)
		require.Equal(t, "Yes", response.Data[0].LastName)
	})

	t.Run("sort by joined latest season value score uses the derived snapshot", func(t *testing.T) {
		response := hitListPlayersEndpoint(t, db, "/api/v1/players?q="+sortPrefix+"&sort=-value_score")

		require.Equal(t, 2, response.Meta.Count)
		require.Len(t, response.Data, 2)
		require.Equal(t, "High", response.Data[0].LastName)
		require.NotNil(t, response.Data[0].LatestSeason)
		require.Equal(t, 88.8, response.Data[0].LatestSeason.ValueScore)
		require.Equal(t, "Low", response.Data[1].LastName)
		require.Equal(t, 12.5, response.Data[1].LatestSeason.ValueScore)
	})

	t.Run("team filter matches the joined snapshot by team name and numeric team id", func(t *testing.T) {
		byName := hitListPlayersEndpoint(t, db, "/api/v1/players?q="+teamPrefix+"&team=Yankees")

		require.Equal(t, 1, byName.Meta.Count)
		require.Len(t, byName.Data, 1)
		require.Equal(t, "NameMatch", byName.Data[0].LastName)
		require.NotNil(t, byName.Data[0].LatestSeason)
		require.Equal(t, 147, byName.Data[0].LatestSeason.TeamID)
		require.Equal(t, "New York Yankees", byName.Data[0].LatestSeason.TeamName)

		byID := hitListPlayersEndpoint(t, db, "/api/v1/players?q="+teamPrefix+"&team=147")

		require.Equal(t, 2, byID.Meta.Count)
		require.Len(t, byID.Data, 2)
		require.Equal(t, "IDMatch", byID.Data[0].LastName)
		require.Equal(t, "NameMatch", byID.Data[1].LastName)
		for _, item := range byID.Data {
			require.NotNil(t, item.LatestSeason)
			require.Equal(t, 147, item.LatestSeason.TeamID)
		}
	})

	t.Run("season filter swaps latest season snapshot to the requested year", func(t *testing.T) {
		response := hitListPlayersEndpoint(t, db, "/api/v1/players?q="+seasonPrefix+"&season=2024")

		require.Equal(t, 2, response.Meta.Count)
		require.Len(t, response.Data, 2)
		require.Equal(t, "OtherYear", response.Data[0].LastName)
		require.Nil(t, response.Data[0].LatestSeason)
		require.Equal(t, "Switch", response.Data[1].LastName)
		require.NotNil(t, response.Data[1].LatestSeason)
		require.Equal(t, 2024, response.Data[1].LatestSeason.Year)
		require.Equal(t, 111, response.Data[1].LatestSeason.TeamID)
		require.Equal(t, "Boston Red Sox", response.Data[1].LatestSeason.TeamName)
		require.Equal(t, 70.2, response.Data[1].LatestSeason.ValueScore)
	})
}

type testPlayerFixture struct {
	MLBID     int
	FirstName string
	LastName  string
	Position  string
	Active    bool
	Season    *testSeasonFixture
}

type testSeasonFixture struct {
	Year       int
	TeamID     int
	TeamName   string
	Age        int
	ValueScore float64
}

func testHandlersDB(t *testing.T) *gorm.DB {
	t.Helper()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL is required for handler integration tests")
	}

	db, err := database.Connect(databaseURL)
	require.NoError(t, err)

	return db
}

func createTestPlayer(t *testing.T, db *gorm.DB, fixture testPlayerFixture) models.Player {
	t.Helper()

	player := models.Player{
		MLBID:     fixture.MLBID,
		FirstName: fixture.FirstName,
		LastName:  fixture.LastName,
		Position:  fixture.Position,
		Active:    fixture.Active,
	}
	require.NoError(t, db.Select("MLBID", "FirstName", "LastName", "Position", "Active").Create(&player).Error)
	if !fixture.Active {
		require.NoError(t, db.Model(&player).UpdateColumn("active", false).Error)
		player.Active = false
	}

	if fixture.Season != nil {
		createTestSeasonStat(t, db, player.ID, *fixture.Season)
	}

	return player
}

func createTestSeasonStat(t *testing.T, db *gorm.DB, playerID uint, fixture testSeasonFixture) {
	t.Helper()

	season := models.SeasonStat{
		PlayerID:    int(playerID),
		Year:        fixture.Year,
		TeamID:      fixture.TeamID,
		TeamName:    fixture.TeamName,
		Age:         fixture.Age,
		ValueScore:  fixture.ValueScore,
		GamesPlayed: 1,
	}
	require.NoError(t, db.Create(&season).Error)
}

func cleanupTestPlayers(t *testing.T, db *gorm.DB, mlbIDs []int) {
	t.Helper()

	var players []models.Player
	require.NoError(t, db.Unscoped().Where("mlb_id IN ?", mlbIDs).Find(&players).Error)
	if len(players) == 0 {
		return
	}

	playerIDs := make([]uint, 0, len(players))
	for _, player := range players {
		playerIDs = append(playerIDs, player.ID)
	}

	require.NoError(t, db.Unscoped().Where("player_id IN ?", playerIDs).Delete(&models.SeasonStat{}).Error)
	require.NoError(t, db.Unscoped().Where("id IN ?", playerIDs).Delete(&models.Player{}).Error)
}

func hitListPlayersEndpoint(t *testing.T, db *gorm.DB, path string) playerListResponse {
	t.Helper()

	e := echo.New()
	RegisterRoutes(e, db, syncjob.NewStatusStore(false), "super-secret")

	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var response playerListResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))

	return response
}
