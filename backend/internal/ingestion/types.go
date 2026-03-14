package ingestion

import "time"

type PlayerRecord struct {
	MLBID       int
	FirstName   string
	LastName    string
	Position    string
	Bats        string
	Throws      string
	DateOfBirth *time.Time
	Active      bool
	ImageURL    string
}

type SeasonStatRecord struct {
	Year        int
	TeamID      int
	TeamName    string
	Age         int
	HasHitting  bool
	HasPitching bool

	GamesPlayed      int
	GamesStarted     int
	PlateAppearances int
	AtBats           int
	Hits             int
	Doubles          int
	Triples          int
	HomeRuns         int
	Runs             int
	RBI              int
	Walks            int
	Strikeouts       int
	StolenBases      int
	BattingAvg       float64
	OBP              float64
	SLG              float64
	OPS              float64
	BABIP            float64

	Wins               int
	Losses             int
	ERA                float64
	WHIP               float64
	InningsPitched     float64
	InningsPitchedOuts int
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

type SavantType string

const (
	SavantTypeBatter  SavantType = "batter"
	SavantTypePitcher SavantType = "pitcher"
)

type SavantEnrichment struct {
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
