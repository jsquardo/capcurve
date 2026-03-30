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

	require.NoError(t, db.Unscoped().Where("player_id IN ?", playerIDs).Delete(&models.CareerArc{}).Error)
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

type careerArcTestResponse struct {
	Data careerArcTestData `json:"data"`
}

type projectionTestResponse struct {
	Data projectionTestData `json:"data"`
}

type careerArcTestData struct {
	Player     careerArcTestPlayer     `json:"player"`
	Arc        *careerArcTestMetadata  `json:"arc"`
	Timeline   []careerArcTestTimeline `json:"timeline"`
	Projection careerArcTestProjection `json:"projection"`
}

type projectionTestData struct {
	Player     careerArcTestPlayer     `json:"player"`
	Projection careerArcTestProjection `json:"projection"`
}

type careerArcTestPlayer struct {
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

type careerArcTestMetadata struct {
	PeakYearStart         int       `json:"peak_year_start"`
	PeakYearEnd           int       `json:"peak_year_end"`
	DeclineOnsetYear      int       `json:"decline_onset_year"`
	ArcShape              string    `json:"arc_shape"`
	PeakValueScore        float64   `json:"peak_value_score"`
	CareerValueScoreTotal float64   `json:"career_value_score_total"`
	LastComputedAt        time.Time `json:"last_computed_at"`
}

type careerArcTestTimeline struct {
	Year         int                        `json:"year"`
	TeamID       int                        `json:"team_id"`
	TeamName     string                     `json:"team_name"`
	Age          int                        `json:"age"`
	ValueScore   float64                    `json:"value_score"`
	IsPeak       bool                       `json:"is_peak"`
	IsProjection bool                       `json:"is_projection"`
	Hitting      *playerDetailHittingStats  `json:"hitting"`
	Pitching     *playerDetailPitchingStats `json:"pitching"`
}

type careerArcTestProjection struct {
	Status         string                          `json:"status"`
	Eligible       bool                            `json:"eligible"`
	Reason         string                          `json:"reason"`
	Points         []careerArcTestProjectionPoint  `json:"points"`
	ConfidenceBand []careerArcTestConfidenceBand   `json:"confidence_band"`
	Comparables    []careerArcTestComparablePlayer `json:"comparables"`
}

type careerArcTestProjectionPoint struct {
	Year         int     `json:"year"`
	Age          int     `json:"age"`
	ValueScore   float64 `json:"value_score"`
	IsProjection bool    `json:"is_projection"`
}

type careerArcTestConfidenceBand struct {
	Year  int     `json:"year"`
	Lower float64 `json:"lower"`
	Upper float64 `json:"upper"`
}

type careerArcTestComparablePlayer struct {
	PlayerID uint   `json:"player_id"`
	MLBID    int    `json:"mlb_id"`
	FullName string `json:"full_name"`
	Position string `json:"position"`
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

	t.Run("returns 400 when the player id is not numeric", func(t *testing.T) {
		rec := hitPlayerEndpointRaw(t, db, "/api/v1/players/abc")

		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.JSONEq(t, `{"error":"invalid player id"}`, rec.Body.String())
	})

	t.Run("returns 500 when the player lookup fails for a database error", func(t *testing.T) {
		failingDB := testHandlersDB(t)
		sqlDB, err := failingDB.DB()
		require.NoError(t, err)
		require.NoError(t, sqlDB.Close())

		rec := hitPlayerEndpointRaw(t, failingDB, "/api/v1/players/1")

		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Contains(t, rec.Body.String(), `"error":`)
	})
}

func TestGetCareerArcEndpoint(t *testing.T) {
	db := testHandlersDB(t)
	nowSuffix := time.Now().UnixNano()

	withArc := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     911400000 + int(nowSuffix%100000),
		FirstName: fmt.Sprintf("CodexArc%d", nowSuffix),
		LastName:  "Summary",
		Position:  "OF",
		Active:    true,
	})
	createTestSeasonStat(t, db, withArc.ID, testSeasonFixture{
		Year:             2022,
		TeamID:           141,
		TeamName:         "Toronto Blue Jays",
		Age:              24,
		ValueScore:       48.2,
		GamesPlayed:      130,
		PlateAppearances: 520,
		AtBats:           470,
		Hits:             141,
		HomeRuns:         19,
		OBP:              0.351,
		SLG:              0.472,
	})
	createTestSeasonStat(t, db, withArc.ID, testSeasonFixture{
		Year:             2023,
		TeamID:           141,
		TeamName:         "Toronto Blue Jays",
		Age:              25,
		ValueScore:       67.9,
		GamesPlayed:      145,
		PlateAppearances: 610,
		AtBats:           548,
		Hits:             169,
		HomeRuns:         31,
		OBP:              0.381,
		SLG:              0.544,
	})
	createTestSeasonStat(t, db, withArc.ID, testSeasonFixture{
		Year:             2024,
		TeamID:           141,
		TeamName:         "Toronto Blue Jays",
		Age:              26,
		ValueScore:       58.6,
		GamesPlayed:      138,
		PlateAppearances: 592,
		AtBats:           530,
		Hits:             156,
		HomeRuns:         24,
		OBP:              0.366,
		SLG:              0.501,
	})

	lastComputedAt := time.Date(2026, 3, 27, 12, 0, 0, 0, time.UTC)
	require.NoError(t, db.Create(&models.CareerArc{
		PlayerID:         int(withArc.ID),
		PeakYearStart:    2023,
		PeakYearEnd:      2023,
		PeakWAR:          0,
		CareerWAR:        0,
		DeclineOnsetYear: 2025,
		ArcShape:         "Peak Prime",
		LastComputedAt:   lastComputedAt,
	}).Error)

	noArc := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     911500000 + int(nowSuffix%100000),
		FirstName: fmt.Sprintf("CodexArc%d", nowSuffix),
		LastName:  "Pending",
		Position:  "SS",
		Active:    true,
	})
	createTestSeasonStat(t, db, noArc.ID, testSeasonFixture{
		Year:             2024,
		TeamID:           147,
		TeamName:         "New York Yankees",
		Age:              27,
		ValueScore:       52.4,
		GamesPlayed:      142,
		PlateAppearances: 598,
		AtBats:           540,
		Hits:             158,
		HomeRuns:         17,
		OBP:              0.349,
		SLG:              0.451,
	})

	twoWayExpectedWOBA := 0.401
	twoWayExpectedERA := 3.42
	twoWay := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     911600000 + int(nowSuffix%100000),
		FirstName: fmt.Sprintf("CodexArc%d", nowSuffix),
		LastName:  "TwoWay",
		Position:  "DH",
		Active:    true,
	})
	createTestSeasonStat(t, db, twoWay.ID, testSeasonFixture{
		Year:               2025,
		TeamID:             119,
		TeamName:           "Los Angeles Dodgers",
		Age:                30,
		ValueScore:         89.4,
		GamesPlayed:        148,
		GamesStarted:       21,
		PlateAppearances:   612,
		AtBats:             549,
		Hits:               166,
		HomeRuns:           39,
		OBP:                0.372,
		SLG:                0.581,
		InningsPitched:     128.2,
		Wins:               11,
		Losses:             4,
		ERA:                3.11,
		WHIP:               1.09,
		HitsAllowed:        97,
		WalksAllowed:       39,
		HomeRunsAllowed:    11,
		StrikeoutsPer9:     11.4,
		WalksPer9:          2.7,
		HitsPer9:           6.8,
		HomeRunsPer9:       0.8,
		StrikeoutWalkRatio: 4.18,
		StrikePercentage:   66.8,
		ExpectedWOBA:       &twoWayExpectedWOBA,
		ExpectedERA:        &twoWayExpectedERA,
	})

	t.Cleanup(func() {
		cleanupTestPlayers(t, db, []int{
			withArc.MLBID,
			noArc.MLBID,
			twoWay.MLBID,
		})
	})

	t.Run("returns player header, arc metadata, timeline, and the same projection payload as the dedicated endpoint", func(t *testing.T) {
		response := hitGetCareerArcEndpoint(t, db, fmt.Sprintf("/api/v1/players/%d/career-arc", withArc.ID), http.StatusOK)
		projectionResponse := hitGetProjectionEndpoint(t, db, fmt.Sprintf("/api/v1/players/%d/projection", withArc.ID), http.StatusOK)

		require.Equal(t, withArc.ID, response.Data.Player.ID)
		require.Equal(t, withArc.FirstName+" "+withArc.LastName, response.Data.Player.FullName)
		require.NotNil(t, response.Data.Arc)
		require.Equal(t, 2023, response.Data.Arc.PeakYearStart)
		require.Equal(t, 2023, response.Data.Arc.PeakYearEnd)
		require.Equal(t, 2025, response.Data.Arc.DeclineOnsetYear)
		require.Equal(t, "Peak Prime", response.Data.Arc.ArcShape)
		require.Equal(t, 67.9, response.Data.Arc.PeakValueScore)
		require.InDelta(t, 174.7, response.Data.Arc.CareerValueScoreTotal, 0.0001)
		require.True(t, response.Data.Arc.LastComputedAt.Equal(lastComputedAt))
		require.Len(t, response.Data.Timeline, 3)
		require.False(t, response.Data.Timeline[0].IsPeak)
		require.True(t, response.Data.Timeline[1].IsPeak)
		require.False(t, response.Data.Timeline[2].IsPeak)
		require.False(t, response.Data.Timeline[0].IsProjection)
		require.NotNil(t, response.Data.Timeline[0].Hitting)
		require.Nil(t, response.Data.Timeline[0].Pitching)
		require.Equal(t, projectionResponse.Data.Projection.Status, response.Data.Projection.Status)
		require.Equal(t, projectionResponse.Data.Projection.Eligible, response.Data.Projection.Eligible)
		require.Equal(t, projectionResponse.Data.Projection.Reason, response.Data.Projection.Reason)
		require.Equal(t, projectionResponse.Data.Projection.Points, response.Data.Projection.Points)
		require.Equal(t, projectionResponse.Data.Projection.ConfidenceBand, response.Data.Projection.ConfidenceBand)
		require.Equal(t, projectionResponse.Data.Projection.Comparables, response.Data.Projection.Comparables)
	})

	t.Run("derives peak_value_score from seasons inside the stored peak window before falling back to the overall max", func(t *testing.T) {
		windowedPeak := createTestPlayer(t, db, testPlayerFixture{
			MLBID:     911700000 + int(nowSuffix%100000),
			FirstName: fmt.Sprintf("CodexArc%d", nowSuffix),
			LastName:  "WindowPeak",
			Position:  "3B",
			Active:    true,
		})
		createTestSeasonStat(t, db, windowedPeak.ID, testSeasonFixture{
			Year:             2021,
			TeamID:           109,
			TeamName:         "Arizona Diamondbacks",
			Age:              24,
			ValueScore:       82.7,
			GamesPlayed:      149,
			PlateAppearances: 630,
			AtBats:           574,
			Hits:             173,
			HomeRuns:         34,
			OBP:              0.374,
			SLG:              0.558,
		})
		createTestSeasonStat(t, db, windowedPeak.ID, testSeasonFixture{
			Year:             2022,
			TeamID:           109,
			TeamName:         "Arizona Diamondbacks",
			Age:              25,
			ValueScore:       61.4,
			GamesPlayed:      141,
			PlateAppearances: 598,
			AtBats:           541,
			Hits:             158,
			HomeRuns:         22,
			OBP:              0.349,
			SLG:              0.486,
		})
		createTestSeasonStat(t, db, windowedPeak.ID, testSeasonFixture{
			Year:             2023,
			TeamID:           109,
			TeamName:         "Arizona Diamondbacks",
			Age:              26,
			ValueScore:       67.2,
			GamesPlayed:      146,
			PlateAppearances: 615,
			AtBats:           556,
			Hits:             165,
			HomeRuns:         26,
			OBP:              0.361,
			SLG:              0.509,
		})
		require.NoError(t, db.Create(&models.CareerArc{
			PlayerID:         int(windowedPeak.ID),
			PeakYearStart:    2022,
			PeakYearEnd:      2023,
			DeclineOnsetYear: 2025,
			ArcShape:         "Late Bloomer",
			LastComputedAt:   lastComputedAt,
		}).Error)

		t.Cleanup(func() {
			cleanupTestPlayers(t, db, []int{windowedPeak.MLBID})
		})

		response := hitGetCareerArcEndpoint(t, db, fmt.Sprintf("/api/v1/players/%d/career-arc", windowedPeak.ID), http.StatusOK)

		require.NotNil(t, response.Data.Arc)
		require.Equal(t, 67.2, response.Data.Arc.PeakValueScore)
		require.NotEqual(t, 82.7, response.Data.Arc.PeakValueScore)
	})

	t.Run("returns 200 with arc null when the player has history but no career arc row", func(t *testing.T) {
		response := hitGetCareerArcEndpoint(t, db, fmt.Sprintf("/api/v1/players/%d/career-arc", noArc.ID), http.StatusOK)

		require.Nil(t, response.Data.Arc)
		require.Len(t, response.Data.Timeline, 1)
		require.Equal(t, 2024, response.Data.Timeline[0].Year)
		require.False(t, response.Data.Timeline[0].IsPeak)
		require.True(t, response.Data.Projection.Eligible)
		require.Equal(t, "ready", response.Data.Projection.Status)
		require.NotEmpty(t, response.Data.Projection.Points)
		require.Len(t, response.Data.Projection.Points, len(response.Data.Projection.ConfidenceBand))
	})

	t.Run("keeps both hitting and pitching branches for two-way timeline rows", func(t *testing.T) {
		response := hitGetCareerArcEndpoint(t, db, fmt.Sprintf("/api/v1/players/%d/career-arc", twoWay.ID), http.StatusOK)

		require.Len(t, response.Data.Timeline, 1)
		point := response.Data.Timeline[0]
		require.NotNil(t, point.Hitting)
		require.NotNil(t, point.Pitching)
		require.Equal(t, 612, point.Hitting.PlateAppearances)
		require.Equal(t, 128.2, point.Pitching.InningsPitched)
		require.NotNil(t, point.Hitting.ExpectedWOBA)
		require.Equal(t, twoWayExpectedWOBA, *point.Hitting.ExpectedWOBA)
		require.NotNil(t, point.Pitching.ExpectedERA)
		require.Equal(t, twoWayExpectedERA, *point.Pitching.ExpectedERA)
	})

	t.Run("returns 400 when the player id is not numeric", func(t *testing.T) {
		rec := hitCareerArcEndpointRaw(t, db, "/api/v1/players/abc/career-arc")

		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.JSONEq(t, `{"error":"invalid player id"}`, rec.Body.String())
	})

	t.Run("returns 404 when the player does not exist", func(t *testing.T) {
		rec := hitCareerArcEndpointRaw(t, db, "/api/v1/players/999999999/career-arc")

		require.Equal(t, http.StatusNotFound, rec.Code)
		require.JSONEq(t, `{"error":"player not found"}`, rec.Body.String())
	})

	t.Run("returns 500 when the player lookup fails for a database error", func(t *testing.T) {
		failingDB := testHandlersDB(t)
		sqlDB, err := failingDB.DB()
		require.NoError(t, err)
		require.NoError(t, sqlDB.Close())

		rec := hitCareerArcEndpointRaw(t, failingDB, "/api/v1/players/1/career-arc")

		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Contains(t, rec.Body.String(), `"error":`)
	})
}

func TestGetProjectionEndpoint(t *testing.T) {
	db := testHandlersDB(t)
	nowSuffix := time.Now().UnixNano()

	active := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     912100000 + int(nowSuffix%100000),
		FirstName: fmt.Sprintf("CodexProjection%d", nowSuffix),
		LastName:  "Active",
		Position:  "OF",
		Active:    true,
	})
	createTestSeasonStat(t, db, active.ID, testSeasonFixture{
		Year:             2023,
		TeamID:           147,
		TeamName:         "New York Yankees",
		Age:              27,
		ValueScore:       54.0,
		GamesPlayed:      138,
		PlateAppearances: 590,
		AtBats:           532,
		Hits:             152,
		HomeRuns:         24,
		StolenBases:      18,
		BattingAvg:       0.286,
		OBP:              0.354,
		SLG:              0.497,
	})
	createTestSeasonStat(t, db, active.ID, testSeasonFixture{
		Year:             2024,
		TeamID:           147,
		TeamName:         "New York Yankees",
		Age:              28,
		ValueScore:       61.0,
		GamesPlayed:      145,
		PlateAppearances: 615,
		AtBats:           550,
		Hits:             165,
		HomeRuns:         27,
		StolenBases:      21,
		BattingAvg:       0.300,
		OBP:              0.368,
		SLG:              0.522,
	})
	createTestSeasonStat(t, db, active.ID, testSeasonFixture{
		Year:             2025,
		TeamID:           147,
		TeamName:         "New York Yankees",
		Age:              29,
		ValueScore:       67.0,
		GamesPlayed:      149,
		PlateAppearances: 632,
		AtBats:           566,
		Hits:             171,
		HomeRuns:         30,
		StolenBases:      19,
		BattingAvg:       0.302,
		OBP:              0.375,
		SLG:              0.541,
	})

	compOne := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     912200000 + int(nowSuffix%100000),
		FirstName: fmt.Sprintf("CodexProjection%d", nowSuffix),
		LastName:  "CompOne",
		Position:  "OF",
		Active:    false,
	})
	for _, season := range []testSeasonFixture{
		{Year: 2021, TeamID: 121, TeamName: "New York Mets", Age: 27, ValueScore: 52.0, GamesPlayed: 140, PlateAppearances: 600, AtBats: 544, Hits: 158, HomeRuns: 22, StolenBases: 17, BattingAvg: 0.290, OBP: 0.356, SLG: 0.489},
		{Year: 2022, TeamID: 121, TeamName: "New York Mets", Age: 28, ValueScore: 60.0, GamesPlayed: 144, PlateAppearances: 618, AtBats: 553, Hits: 166, HomeRuns: 26, StolenBases: 20, BattingAvg: 0.300, OBP: 0.366, SLG: 0.518},
		{Year: 2023, TeamID: 121, TeamName: "New York Mets", Age: 29, ValueScore: 65.0, GamesPlayed: 146, PlateAppearances: 624, AtBats: 557, Hits: 168, HomeRuns: 28, StolenBases: 18, BattingAvg: 0.302, OBP: 0.370, SLG: 0.530},
		{Year: 2024, TeamID: 121, TeamName: "New York Mets", Age: 30, ValueScore: 63.0, GamesPlayed: 143, PlateAppearances: 610, AtBats: 545, Hits: 160, HomeRuns: 24, StolenBases: 15, BattingAvg: 0.294, OBP: 0.361, SLG: 0.505},
		{Year: 2025, TeamID: 121, TeamName: "New York Mets", Age: 31, ValueScore: 58.0, GamesPlayed: 137, PlateAppearances: 585, AtBats: 525, Hits: 149, HomeRuns: 20, StolenBases: 13, BattingAvg: 0.284, OBP: 0.347, SLG: 0.472},
	} {
		createTestSeasonStat(t, db, compOne.ID, season)
	}

	compTwo := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     912300000 + int(nowSuffix%100000),
		FirstName: fmt.Sprintf("CodexProjection%d", nowSuffix),
		LastName:  "CompTwo",
		Position:  "RF",
		Active:    false,
	})
	for _, season := range []testSeasonFixture{
		{Year: 2020, TeamID: 111, TeamName: "Boston Red Sox", Age: 26, ValueScore: 50.0, GamesPlayed: 133, PlateAppearances: 575, AtBats: 520, Hits: 150, HomeRuns: 20, StolenBases: 16, BattingAvg: 0.288, OBP: 0.350, SLG: 0.481},
		{Year: 2021, TeamID: 111, TeamName: "Boston Red Sox", Age: 27, ValueScore: 55.0, GamesPlayed: 141, PlateAppearances: 602, AtBats: 540, Hits: 161, HomeRuns: 24, StolenBases: 18, BattingAvg: 0.298, OBP: 0.360, SLG: 0.503},
		{Year: 2022, TeamID: 111, TeamName: "Boston Red Sox", Age: 28, ValueScore: 62.0, GamesPlayed: 145, PlateAppearances: 621, AtBats: 556, Hits: 170, HomeRuns: 27, StolenBases: 20, BattingAvg: 0.306, OBP: 0.371, SLG: 0.527},
		{Year: 2023, TeamID: 111, TeamName: "Boston Red Sox", Age: 29, ValueScore: 66.0, GamesPlayed: 148, PlateAppearances: 629, AtBats: 560, Hits: 172, HomeRuns: 29, StolenBases: 17, BattingAvg: 0.307, OBP: 0.373, SLG: 0.536},
		{Year: 2024, TeamID: 111, TeamName: "Boston Red Sox", Age: 30, ValueScore: 61.0, GamesPlayed: 140, PlateAppearances: 598, AtBats: 535, Hits: 155, HomeRuns: 23, StolenBases: 14, BattingAvg: 0.290, OBP: 0.355, SLG: 0.491},
		{Year: 2025, TeamID: 111, TeamName: "Boston Red Sox", Age: 31, ValueScore: 57.0, GamesPlayed: 134, PlateAppearances: 570, AtBats: 510, Hits: 145, HomeRuns: 18, StolenBases: 11, BattingAvg: 0.284, OBP: 0.344, SLG: 0.460},
	} {
		createTestSeasonStat(t, db, compTwo.ID, season)
	}

	retired := createTestPlayer(t, db, testPlayerFixture{
		MLBID:     912400000 + int(nowSuffix%100000),
		FirstName: fmt.Sprintf("CodexProjection%d", nowSuffix),
		LastName:  "Retired",
		Position:  "1B",
		Active:    false,
	})
	createTestSeasonStat(t, db, retired.ID, testSeasonFixture{
		Year:             2022,
		TeamID:           138,
		TeamName:         "St. Louis Cardinals",
		Age:              42,
		ValueScore:       34.0,
		GamesPlayed:      109,
		PlateAppearances: 412,
		AtBats:           378,
		Hits:             102,
		HomeRuns:         17,
		BattingAvg:       0.270,
		OBP:              0.332,
		SLG:              0.448,
	})

	t.Cleanup(func() {
		cleanupTestPlayers(t, db, []int{
			active.MLBID,
			compOne.MLBID,
			compTwo.MLBID,
			retired.MLBID,
		})
	})

	t.Run("returns populated projections for an active player", func(t *testing.T) {
		response := hitGetProjectionEndpoint(t, db, fmt.Sprintf("/api/v1/players/%d/projection", active.ID), http.StatusOK)

		require.Equal(t, active.ID, response.Data.Player.ID)
		require.True(t, response.Data.Projection.Eligible)
		require.Equal(t, "ready", response.Data.Projection.Status)
		require.Empty(t, response.Data.Projection.Reason)
		require.NotEmpty(t, response.Data.Projection.Points)
		require.Len(t, response.Data.Projection.ConfidenceBand, len(response.Data.Projection.Points))
		require.NotEmpty(t, response.Data.Projection.Comparables)
		require.LessOrEqual(t, len(response.Data.Projection.Comparables), 5)
		require.Equal(t, 2026, response.Data.Projection.Points[0].Year)
		require.Equal(t, 30, response.Data.Projection.Points[0].Age)
		require.True(t, response.Data.Projection.Points[0].IsProjection)
		require.LessOrEqual(t, response.Data.Projection.ConfidenceBand[0].Lower, response.Data.Projection.Points[0].ValueScore)
		require.GreaterOrEqual(t, response.Data.Projection.ConfidenceBand[0].Upper, response.Data.Projection.Points[0].ValueScore)
	})

	t.Run("returns 200 with an ineligible payload for retired players", func(t *testing.T) {
		response := hitGetProjectionEndpoint(t, db, fmt.Sprintf("/api/v1/players/%d/projection", retired.ID), http.StatusOK)

		require.Equal(t, retired.ID, response.Data.Player.ID)
		require.False(t, response.Data.Projection.Eligible)
		require.Equal(t, "ineligible", response.Data.Projection.Status)
		require.Equal(t, "player is not active", response.Data.Projection.Reason)
		require.Empty(t, response.Data.Projection.Points)
		require.Empty(t, response.Data.Projection.ConfidenceBand)
		require.Empty(t, response.Data.Projection.Comparables)
	})

	t.Run("returns 400 when the player id is not numeric", func(t *testing.T) {
		rec := hitProjectionEndpointRaw(t, db, "/api/v1/players/abc/projection")

		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.JSONEq(t, `{"error":"invalid player id"}`, rec.Body.String())
	})

	t.Run("returns 404 when the player does not exist", func(t *testing.T) {
		rec := hitProjectionEndpointRaw(t, db, "/api/v1/players/999999999/projection")

		require.Equal(t, http.StatusNotFound, rec.Code)
		require.JSONEq(t, `{"error":"player not found"}`, rec.Body.String())
	})

	t.Run("returns 500 when the player lookup fails for a database error", func(t *testing.T) {
		failingDB := testHandlersDB(t)
		sqlDB, err := failingDB.DB()
		require.NoError(t, err)
		require.NoError(t, sqlDB.Close())

		rec := hitProjectionEndpointRaw(t, failingDB, "/api/v1/players/1/projection")

		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Contains(t, rec.Body.String(), `"error":`)
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

func hitGetCareerArcEndpoint(t *testing.T, db *gorm.DB, path string, expectedStatus int) careerArcTestResponse {
	t.Helper()

	rec := hitCareerArcEndpointRaw(t, db, path)
	require.Equal(t, expectedStatus, rec.Code)

	var response careerArcTestResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))

	return response
}

func hitGetProjectionEndpoint(t *testing.T, db *gorm.DB, path string, expectedStatus int) projectionTestResponse {
	t.Helper()

	rec := hitProjectionEndpointRaw(t, db, path)
	require.Equal(t, expectedStatus, rec.Code)

	var response projectionTestResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))

	return response
}

func hitCareerArcEndpointRaw(t *testing.T, db *gorm.DB, path string) *httptest.ResponseRecorder {
	t.Helper()

	e := echo.New()
	RegisterRoutes(e, db, syncjob.NewStatusStore(false), "super-secret")

	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	return rec
}

func hitProjectionEndpointRaw(t *testing.T, db *gorm.DB, path string) *httptest.ResponseRecorder {
	t.Helper()

	e := echo.New()
	RegisterRoutes(e, db, syncjob.NewStatusStore(false), "super-secret")

	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	return rec
}
