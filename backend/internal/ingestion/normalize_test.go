package ingestion

import (
	"testing"
	"time"
)

func TestNormalizePlayer(t *testing.T) {
	t.Parallel()

	player, err := NormalizePlayer(&MLBPlayer{
		ID:        660271,
		FirstName: "Shohei",
		LastName:  "Ohtani",
		BirthDate: "1994-07-05",
		Active:    true,
		PrimaryPosition: mlbPosition{
			Name: "Two-Way Player",
		},
		BatSide: mlbHandDescriptor{Code: "L"},
		PitchHand: mlbHandDescriptor{
			Code: "R",
		},
	})
	if err != nil {
		t.Fatalf("NormalizePlayer returned error: %v", err)
	}

	if player.Position != "Two-Way Player" {
		t.Fatalf("Position = %q, want %q", player.Position, "Two-Way Player")
	}
	if player.ImageURL == "" {
		t.Fatal("ImageURL is empty")
	}
	if player.DateOfBirth == nil || !player.DateOfBirth.Equal(time.Date(1994, 7, 5, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("DateOfBirth = %v, want 1994-07-05", player.DateOfBirth)
	}
}

func TestMergeSeasonStats(t *testing.T) {
	t.Parallel()

	hitting := SeasonStatRecord{
		Year:             2021,
		TeamID:           108,
		TeamName:         "Los Angeles Angels",
		HasHitting:       true,
		GamesPlayed:      158,
		PlateAppearances: 639,
		HomeRuns:         46,
		OPS:              0.964,
	}

	pitching := SeasonStatRecord{
		Year:               2021,
		TeamID:             108,
		TeamName:           "Los Angeles Angels",
		HasPitching:        true,
		GamesStarted:       23,
		InningsPitched:     130.1,
		InningsPitchedOuts: 391,
		ERA:                3.18,
		StrikeoutsPer9:     10.77,
		StrikeoutWalkRatio: 3.55,
		StrikePercentage:   0.64,
	}

	merged := MergeSeasonStats(hitting, pitching)
	if merged.HomeRuns != 46 {
		t.Fatalf("HomeRuns = %d, want 46", merged.HomeRuns)
	}
	if merged.GamesStarted != 23 {
		t.Fatalf("GamesStarted = %d, want 23", merged.GamesStarted)
	}
	if merged.InningsPitched != 130.1 {
		t.Fatalf("InningsPitched = %v, want 130.1", merged.InningsPitched)
	}
	if !merged.HasHitting || !merged.HasPitching {
		t.Fatalf("expected merged record to retain both role flags: %+v", merged)
	}
}

func TestNormalizeSeasonSplit(t *testing.T) {
	t.Parallel()

	split := MLBSeasonSplit{
		Season: "2024",
		Stat: map[string]any{
			"age":              float64(29),
			"gamesPlayed":      float64(159),
			"plateAppearances": float64(731),
			"homeRuns":         float64(54),
			"ops":              "1.036",
			"babip":            ".336",
		},
		Team: mlbStatTeam{
			ID:   119,
			Name: "Los Angeles Dodgers",
		},
	}

	record, err := NormalizeSeasonSplit(split, "hitting")
	if err != nil {
		t.Fatalf("NormalizeSeasonSplit returned error: %v", err)
	}

	if record.Year != 2024 || record.TeamID != 119 {
		t.Fatalf("unexpected record identity: %+v", record)
	}
	if record.OPS != 1.036 {
		t.Fatalf("OPS = %v, want 1.036", record.OPS)
	}
	if record.BABIP != 0.336 {
		t.Fatalf("BABIP = %v, want 0.336", record.BABIP)
	}
}

func TestNormalizeSeasonSplitPitchingConvertsBaseballInningsToOuts(t *testing.T) {
	t.Parallel()

	split := MLBSeasonSplit{
		Season: "2024",
		Stat: map[string]any{
			"inningsPitched": "130.2",
			"era":            "3.47",
		},
		Team: mlbStatTeam{
			ID:   119,
			Name: "Los Angeles Dodgers",
		},
	}

	record, err := NormalizeSeasonSplit(split, "pitching")
	if err != nil {
		t.Fatalf("NormalizeSeasonSplit returned error: %v", err)
	}

	if record.InningsPitchedOuts != 392 {
		t.Fatalf("InningsPitchedOuts = %d, want 392", record.InningsPitchedOuts)
	}
	if record.InningsPitched != 130.2 {
		t.Fatalf("InningsPitched = %v, want 130.2", record.InningsPitched)
	}
}

func TestIsAggregateSeasonSplit(t *testing.T) {
	t.Parallel()

	if !IsAggregateSeasonSplit(MLBSeasonSplit{
		Season: "2022",
		Team: mlbStatTeam{
			ID:   0,
			Name: "TOT",
		},
	}) {
		t.Fatal("expected team_id 0 split to be treated as aggregate")
	}

	if IsAggregateSeasonSplit(MLBSeasonSplit{
		Season: "2022",
		Team: mlbStatTeam{
			ID:   120,
			Name: "Washington Nationals",
		},
	}) {
		t.Fatal("expected real team split to be retained")
	}
}

func TestMergeSplitsSkipsAggregateRows(t *testing.T) {
	t.Parallel()

	service := &Service{}
	records := make(map[string]SeasonStatRecord)
	splits := []MLBSeasonSplit{
		{
			Season: "2022",
			Stat: map[string]any{
				"plateAppearances": float64(52),
				"homeRuns":         float64(6),
			},
			Team: mlbStatTeam{
				ID:   120,
				Name: "Washington Nationals",
			},
		},
		{
			Season: "2022",
			Stat: map[string]any{
				"plateAppearances": float64(427),
				"homeRuns":         float64(21),
			},
			Team: mlbStatTeam{
				ID:   109,
				Name: "San Diego Padres",
			},
		},
		{
			Season: "2022",
			Stat: map[string]any{
				"plateAppearances": float64(479),
				"homeRuns":         float64(27),
			},
			Team: mlbStatTeam{
				ID:   0,
				Name: "TOT",
			},
		},
	}

	if err := service.mergeSplits(records, splits, "hitting"); err != nil {
		t.Fatalf("mergeSplits returned error: %v", err)
	}

	if len(records) != 1 {
		t.Fatalf("expected 1 merged season record, got %d", len(records))
	}

	record, ok := records[seasonKey(2022)]
	if !ok {
		t.Fatalf("expected merged 2022 season record, got keys: %+v", records)
	}
	if _, ok := records["2022:0"]; ok {
		t.Fatal("aggregate team_id 0 row should not be present")
	}
	if record.TeamID != 109 || record.TeamName != "San Diego Padres" {
		t.Fatalf("expected canonical team metadata to reflect final split, got team_id=%d team_name=%q", record.TeamID, record.TeamName)
	}
	if record.PlateAppearances != 479 {
		t.Fatalf("PlateAppearances = %d, want 479", record.PlateAppearances)
	}
	if record.HomeRuns != 27 {
		t.Fatalf("HomeRuns = %d, want 27", record.HomeRuns)
	}
}

func TestMergeSeasonGroupAggregatesPitchingRatesAcrossTeams(t *testing.T) {
	t.Parallel()

	first := SeasonStatRecord{
		Year:               2017,
		TeamID:             117,
		TeamName:           "Detroit Tigers",
		HasPitching:        true,
		GamesPlayed:        28,
		GamesStarted:       28,
		Wins:               10,
		Losses:             8,
		InningsPitched:     172.0,
		InningsPitchedOuts: 516,
		HitsAllowed:        138,
		WalksAllowed:       42,
		HomeRunsAllowed:    21,
		Strikeouts:         176,
		ERA:                3.82,
		WHIP:               1.047,
		StrikeoutsPer9:     9.209,
		WalksPer9:          2.198,
		HitsPer9:           7.221,
		HomeRunsPer9:       1.099,
		StrikeoutWalkRatio: 4.19,
		StrikePercentage:   0.668,
	}
	second := SeasonStatRecord{
		Year:               2017,
		TeamID:             146,
		TeamName:           "Houston Astros",
		HasPitching:        true,
		GamesPlayed:        5,
		GamesStarted:       5,
		Wins:               5,
		Losses:             0,
		InningsPitched:     34.0,
		InningsPitchedOuts: 102,
		HitsAllowed:        26,
		WalksAllowed:       5,
		HomeRunsAllowed:    2,
		Strikeouts:         43,
		ERA:                1.06,
		WHIP:               0.912,
		StrikeoutsPer9:     11.382,
		WalksPer9:          1.324,
		HitsPer9:           6.882,
		HomeRunsPer9:       0.529,
		StrikeoutWalkRatio: 8.6,
		StrikePercentage:   0.69,
	}

	merged := MergeSeasonGroup(first, second, "pitching")

	if merged.TeamID != 146 || merged.TeamName != "Houston Astros" {
		t.Fatalf("expected final team metadata, got team_id=%d team_name=%q", merged.TeamID, merged.TeamName)
	}
	if merged.InningsPitched != 206 {
		t.Fatalf("InningsPitched = %v, want 206", merged.InningsPitched)
	}
	if merged.InningsPitchedOuts != 618 {
		t.Fatalf("InningsPitchedOuts = %d, want 618", merged.InningsPitchedOuts)
	}
	if merged.ERA != 3.364 {
		t.Fatalf("ERA = %v, want 3.364", merged.ERA)
	}
	if merged.WHIP != 1.024 {
		t.Fatalf("WHIP = %v, want 1.024", merged.WHIP)
	}
	if merged.StrikeoutsPer9 != 9.568 {
		t.Fatalf("StrikeoutsPer9 = %v, want 9.568", merged.StrikeoutsPer9)
	}
}

func TestMergeSeasonGroupAggregatesBaseballPartialInningsAcrossTeams(t *testing.T) {
	t.Parallel()

	first := SeasonStatRecord{
		Year:               2025,
		TeamID:             111,
		TeamName:           "Boston Red Sox",
		HasPitching:        true,
		GamesPlayed:        10,
		GamesStarted:       10,
		InningsPitched:     10.1,
		InningsPitchedOuts: 31,
		HitsAllowed:        10,
		WalksAllowed:       3,
		HomeRunsAllowed:    1,
		Strikeouts:         15,
		ERA:                2.613,
		StrikePercentage:   0.64,
	}
	second := SeasonStatRecord{
		Year:               2025,
		TeamID:             112,
		TeamName:           "Chicago Cubs",
		HasPitching:        true,
		GamesPlayed:        4,
		GamesStarted:       4,
		InningsPitched:     2.2,
		InningsPitchedOuts: 8,
		HitsAllowed:        4,
		WalksAllowed:       1,
		HomeRunsAllowed:    0,
		Strikeouts:         5,
		ERA:                3.375,
		StrikePercentage:   0.61,
	}

	merged := MergeSeasonGroup(first, second, "pitching")

	if merged.InningsPitchedOuts != 39 {
		t.Fatalf("InningsPitchedOuts = %d, want 39", merged.InningsPitchedOuts)
	}
	if merged.InningsPitched != 13 {
		t.Fatalf("InningsPitched = %v, want 13", merged.InningsPitched)
	}
	if merged.ERA != 2.769 {
		t.Fatalf("ERA = %v, want 2.769", merged.ERA)
	}
}
