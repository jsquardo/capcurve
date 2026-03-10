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
		GamesPlayed:      158,
		PlateAppearances: 639,
		HomeRuns:         46,
		OPS:              0.964,
	}

	pitching := SeasonStatRecord{
		Year:               2021,
		TeamID:             108,
		TeamName:           "Los Angeles Angels",
		GamesStarted:       23,
		InningsPitched:     130.1,
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
