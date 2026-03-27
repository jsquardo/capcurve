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
	Year               int
	TeamID             int
	TeamName           string
	Age                int
	ValueScore         float64
	GamesPlayed        int
	GamesStarted       int
	PlateAppearances   int
	AtBats             int
	Hits               int
	Doubles            int
	Triples            int
	HomeRuns           int
	Runs               int
	RBI                int
	Walks              int
	Strikeouts         int
	StolenBases        int
	BattingAvg         float64
	OBP                float64
	SLG                float64
	OPS                float64
	BABIP              float64
	Wins               int
	Losses             int
	ERA                float64
	WHIP               float64
	InningsPitched     float64
	HitsAllowed        int
	WalksAllowed       int
	HomeRunsAllowed    int
	StrikeoutsPer9     float64
	WalksPer9          float64
	HitsPer9           float64
	HomeRunsPer9       float64
	StrikeoutWalkRatio float64
	StrikePercentage   float64
	ExpectedBattingAvg *float64
	ExpectedSlugging   *float64
	ExpectedWOBA       *float64
	ExpectedERA        *float64
	BarrelPct          *float64
	HardHitPct         *float64
	AvgExitVelocity    *float64
	AvgLaunchAngle     *float64
	SweetSpotPct       *float64
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
		PlayerID:           int(playerID),
		Year:               fixture.Year,
		TeamID:             fixture.TeamID,
		TeamName:           fixture.TeamName,
		Age:                fixture.Age,
		ValueScore:         fixture.ValueScore,
		GamesPlayed:        fixture.GamesPlayed,
		GamesStarted:       fixture.GamesStarted,
		PlateAppearances:   fixture.PlateAppearances,
		AtBats:             fixture.AtBats,
		Hits:               fixture.Hits,
		Doubles:            fixture.Doubles,
		Triples:            fixture.Triples,
		HomeRuns:           fixture.HomeRuns,
		Runs:               fixture.Runs,
		RBI:                fixture.RBI,
		Walks:              fixture.Walks,
		Strikeouts:         fixture.Strikeouts,
		StolenBases:        fixture.StolenBases,
		BattingAvg:         fixture.BattingAvg,
		OBP:                fixture.OBP,
		SLG:                fixture.SLG,
		OPS:                fixture.OPS,
		BABIP:              fixture.BABIP,
		Wins:               fixture.Wins,
		Losses:             fixture.Losses,
		ERA:                fixture.ERA,
		WHIP:               fixture.WHIP,
		InningsPitched:     fixture.InningsPitched,
		HitsAllowed:        fixture.HitsAllowed,
		WalksAllowed:       fixture.WalksAllowed,
		HomeRunsAllowed:    fixture.HomeRunsAllowed,
		StrikeoutsPer9:     fixture.StrikeoutsPer9,
		WalksPer9:          fixture.WalksPer9,
		HitsPer9:           fixture.HitsPer9,
		HomeRunsPer9:       fixture.HomeRunsPer9,
		StrikeoutWalkRatio: fixture.StrikeoutWalkRatio,
		StrikePercentage:   fixture.StrikePercentage,
		ExpectedBattingAvg: fixture.ExpectedBattingAvg,
		ExpectedSlugging:   fixture.ExpectedSlugging,
		ExpectedWOBA:       fixture.ExpectedWOBA,
		ExpectedERA:        fixture.ExpectedERA,
		BarrelPct:          fixture.BarrelPct,
		HardHitPct:         fixture.HardHitPct,
		AvgExitVelocity:    fixture.AvgExitVelocity,
		AvgLaunchAngle:     fixture.AvgLaunchAngle,
		SweetSpotPct:       fixture.SweetSpotPct,
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

type playerDetailTestResponse struct {
	Data playerDetailTestItem `json:"data"`
}

type playerDetailTestItem struct {
	ID           uint                         `json:"id"`
	MLBID        int                          `json:"mlb_id"`
	FirstName    string                       `json:"first_name"`
	LastName     string                       `json:"last_name"`
	FullName     string                       `json:"full_name"`
	Position     string                       `json:"position"`
	Bats         string                       `json:"bats"`
	Throws       string                       `json:"throws"`
	DateOfBirth  *time.Time                   `json:"date_of_birth"`
	Active       bool                         `json:"active"`
	ImageURL     string                       `json:"image_url"`
	LatestSeason *playerSeasonListItem        `json:"latest_season"`
	CareerStats  []playerDetailCareerStatItem `json:"career_stats"`
}

type playerDetailCareerStatItem struct {
	Year       int                        `json:"year"`
	TeamID     int                        `json:"team_id"`
	TeamName   string                     `json:"team_name"`
	Age        int                        `json:"age"`
	ValueScore float64                    `json:"value_score"`
	Hitting    *playerDetailHittingStats  `json:"hitting"`
	Pitching   *playerDetailPitchingStats `json:"pitching"`
}

type playerDetailHittingStats struct {
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

type playerDetailPitchingStats struct {
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

func TestGetPlayerEndpoint(t *testing.T) {
	db := testHandlersDB(t)
	nowSuffix := time.Now().UnixNano()

	twoWayExpectedBatting := 0.301
	twoWayExpectedSlugging := 0.577
	twoWayExpectedWOBA := 0.412
	twoWayExpectedERA := 3.18
	twoWayBarrelPct := 18.4
	twoWayHardHitPct := 52.1
	twoWayExitVelocity := 95.2
	twoWayLaunchAngle := 13.6
	twoWaySweetSpotPct := 38.7

	withStats := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     911100000 + int(nowSuffix%100000),
		FirstName: fmt.Sprintf("CodexDetail%d", nowSuffix),
		LastName:  "History",
		Position:  "OF",
		Active:    true,
	})
	createTestSeasonStat(t, db, withStats.ID, testSeasonFixture{
		Year:             2023,
		TeamID:           121,
		TeamName:         "New York Mets",
		Age:              25,
		ValueScore:       44.4,
		GamesPlayed:      120,
		PlateAppearances: 480,
		AtBats:           430,
		Hits:             129,
		Doubles:          24,
		Triples:          3,
		HomeRuns:         18,
		Runs:             66,
		RBI:              72,
		Walks:            42,
		Strikeouts:       101,
		StolenBases:      12,
		BattingAvg:       0.300,
		OBP:              0.362,
		SLG:              0.488,
		OPS:              0.850,
		BABIP:            0.321,
	})
	createTestSeasonStat(t, db, withStats.ID, testSeasonFixture{
		Year:             2024,
		TeamID:           147,
		TeamName:         "New York Yankees",
		Age:              26,
		ValueScore:       57.9,
		GamesPlayed:      145,
		PlateAppearances: 610,
		AtBats:           550,
		Hits:             170,
		Doubles:          32,
		Triples:          4,
		HomeRuns:         27,
		Runs:             95,
		RBI:              88,
		Walks:            51,
		Strikeouts:       118,
		StolenBases:      17,
		BattingAvg:       0.309,
		OBP:              0.371,
		SLG:              0.531,
		OPS:              0.902,
		BABIP:            0.333,
	})

	noStats := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     911200000 + int(nowSuffix%100000),
		FirstName: fmt.Sprintf("CodexDetail%d", nowSuffix),
		LastName:  "Empty",
		Position:  "SS",
		Active:    true,
	})

	twoWay := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     911300000 + int(nowSuffix%100000),
		FirstName: fmt.Sprintf("CodexDetail%d", nowSuffix),
		LastName:  "TwoWay",
		Position:  "DH",
		Active:    true,
	})
	createTestSeasonStat(t, db, twoWay.ID, testSeasonFixture{
		Year:               2025,
		TeamID:             119,
		TeamName:           "Los Angeles Dodgers",
		Age:                30,
		ValueScore:         91.3,
		GamesPlayed:        150,
		GamesStarted:       22,
		PlateAppearances:   625,
		AtBats:             560,
		Hits:               168,
		Doubles:            26,
		Triples:            6,
		HomeRuns:           44,
		Runs:               112,
		RBI:                108,
		Walks:              58,
		Strikeouts:         142,
		StolenBases:        21,
		BattingAvg:         0.300,
		OBP:                0.377,
		SLG:                0.590,
		OPS:                0.967,
		BABIP:              0.314,
		Wins:               10,
		Losses:             3,
		ERA:                3.05,
		WHIP:               1.07,
		InningsPitched:     132.1,
		HitsAllowed:        102,
		WalksAllowed:       41,
		HomeRunsAllowed:    12,
		StrikeoutsPer9:     11.8,
		WalksPer9:          2.8,
		HitsPer9:           6.9,
		HomeRunsPer9:       0.8,
		StrikeoutWalkRatio: 4.21,
		StrikePercentage:   67.4,
		ExpectedBattingAvg: &twoWayExpectedBatting,
		ExpectedSlugging:   &twoWayExpectedSlugging,
		ExpectedWOBA:       &twoWayExpectedWOBA,
		ExpectedERA:        &twoWayExpectedERA,
		BarrelPct:          &twoWayBarrelPct,
		HardHitPct:         &twoWayHardHitPct,
		AvgExitVelocity:    &twoWayExitVelocity,
		AvgLaunchAngle:     &twoWayLaunchAngle,
		SweetSpotPct:       &twoWaySweetSpotPct,
	})

	t.Cleanup(func() {
		cleanupTestPlayers(t, db, []int{
			withStats.MLBID,
			noStats.MLBID,
			twoWay.MLBID,
		})
	})

	t.Run("returns player detail with full career stats and latest season derived from the last sorted row", func(t *testing.T) {
		response := hitGetPlayerEndpoint(t, db, fmt.Sprintf("/api/v1/players/%d", withStats.ID), http.StatusOK)

		require.Equal(t, withStats.ID, response.Data.ID)
		require.Equal(t, withStats.MLBID, response.Data.MLBID)
		require.Equal(t, withStats.FirstName+" "+withStats.LastName, response.Data.FullName)
		require.NotNil(t, response.Data.LatestSeason)
		require.Equal(t, 2024, response.Data.LatestSeason.Year)
		require.Equal(t, "New York Yankees", response.Data.LatestSeason.TeamName)
		require.Len(t, response.Data.CareerStats, 2)
		require.Equal(t, 2023, response.Data.CareerStats[0].Year)
		require.Equal(t, 2024, response.Data.CareerStats[1].Year)
		require.NotNil(t, response.Data.CareerStats[0].Hitting)
		require.Nil(t, response.Data.CareerStats[0].Pitching)
		require.Equal(t, 610, response.Data.CareerStats[1].Hitting.PlateAppearances)
	})

	t.Run("returns empty career stats and null latest season when the player has no season rows", func(t *testing.T) {
		response := hitGetPlayerEndpoint(t, db, fmt.Sprintf("/api/v1/players/%d", noStats.ID), http.StatusOK)

		require.Nil(t, response.Data.LatestSeason)
		require.Empty(t, response.Data.CareerStats)
	})

	t.Run("populates both hitting and pitching for a two-way season row", func(t *testing.T) {
		response := hitGetPlayerEndpoint(t, db, fmt.Sprintf("/api/v1/players/%d", twoWay.ID), http.StatusOK)

		require.Len(t, response.Data.CareerStats, 1)
		season := response.Data.CareerStats[0]
		require.NotNil(t, season.Hitting)
		require.NotNil(t, season.Pitching)
		require.Equal(t, 625, season.Hitting.PlateAppearances)
		require.Equal(t, 132.1, season.Pitching.InningsPitched)
		require.NotNil(t, season.Hitting.ExpectedWOBA)
		require.Equal(t, twoWayExpectedWOBA, *season.Hitting.ExpectedWOBA)
		require.NotNil(t, season.Pitching.ExpectedERA)
		require.Equal(t, twoWayExpectedERA, *season.Pitching.ExpectedERA)
	})

	t.Run("returns 404 when the player does not exist", func(t *testing.T) {
		rec := hitPlayerEndpointRaw(t, db, "/api/v1/players/999999999")

		require.Equal(t, http.StatusNotFound, rec.Code)
		require.JSONEq(t, `{"error":"player not found"}`, rec.Body.String())
	})
}

func hitGetPlayerEndpoint(t *testing.T, db *gorm.DB, path string, expectedStatus int) playerDetailTestResponse {
	t.Helper()

	rec := hitPlayerEndpointRaw(t, db, path)
	require.Equal(t, expectedStatus, rec.Code)

	var response playerDetailTestResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))

	return response
}

func hitPlayerEndpointRaw(t *testing.T, db *gorm.DB, path string) *httptest.ResponseRecorder {
	t.Helper()

	e := echo.New()
	RegisterRoutes(e, db, syncjob.NewStatusStore(false), "super-secret")

	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	return rec
}
