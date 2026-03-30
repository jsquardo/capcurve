package projection

import (
	"testing"

	"github.com/jsquardo/capcurve/internal/models"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestServiceBuildReturnsOnlyQualityComparableMatches(t *testing.T) {
	service := NewService()

	target := models.Player{Model: testModel(1), MLBID: 1001, FirstName: "Target", LastName: "Player", Position: "OF", Active: true}
	targetHistory := []models.SeasonStat{
		testSeason(1, 1, 2023, 27, 54, 580, 0, 22, 16, 0.285, 0.350, 0.488),
		testSeason(2, 1, 2024, 28, 60, 610, 0, 26, 18, 0.295, 0.362, 0.515),
		testSeason(3, 1, 2025, 29, 66, 625, 0, 29, 17, 0.301, 0.371, 0.534),
	}

	goodComp := models.Player{Model: testModel(2), MLBID: 1002, FirstName: "Good", LastName: "Comp", Position: "OF", Active: false}
	badComp := models.Player{Model: testModel(3), MLBID: 1003, FirstName: "Bad", LastName: "Comp", Position: "C", Active: false}

	allPlayers := []models.Player{goodComp, badComp}
	allStats := []models.SeasonStat{
		testSeason(4, 2, 2021, 27, 52, 575, 0, 21, 15, 0.284, 0.349, 0.482),
		testSeason(5, 2, 2022, 28, 59, 605, 0, 25, 17, 0.293, 0.360, 0.508),
		testSeason(6, 2, 2023, 29, 64, 620, 0, 27, 16, 0.299, 0.368, 0.526),
		testSeason(7, 2, 2024, 30, 61, 598, 0, 24, 13, 0.289, 0.354, 0.492),
		testSeason(8, 2, 2025, 31, 57, 566, 0, 19, 11, 0.281, 0.342, 0.458),
		testSeason(9, 3, 2021, 27, 28, 410, 0, 8, 1, 0.236, 0.298, 0.371),
		testSeason(10, 3, 2022, 28, 24, 390, 0, 6, 1, 0.227, 0.285, 0.348),
		testSeason(11, 3, 2023, 29, 19, 365, 0, 5, 0, 0.219, 0.277, 0.332),
		testSeason(12, 3, 2024, 30, 14, 330, 0, 3, 0, 0.205, 0.260, 0.301),
	}

	result := service.Build(target, targetHistory, allPlayers, allStats)

	require.Equal(t, "ready", result.Status)
	require.Len(t, result.Comparables, 1)
	require.Equal(t, goodComp.ID, result.Comparables[0].PlayerID)
	require.Len(t, result.Points, len(result.ConfidenceBand))
}

func TestServiceBuildSkipsActiveComparableCandidates(t *testing.T) {
	service := NewService()

	target := models.Player{Model: testModel(1), MLBID: 1001, FirstName: "Target", LastName: "Player", Position: "OF", Active: true}
	targetHistory := []models.SeasonStat{
		testSeason(1, 1, 2023, 27, 54, 580, 0, 22, 16, 0.285, 0.350, 0.488),
		testSeason(2, 1, 2024, 28, 60, 610, 0, 26, 18, 0.295, 0.362, 0.515),
		testSeason(3, 1, 2025, 29, 66, 625, 0, 29, 17, 0.301, 0.371, 0.534),
	}

	activeComp := models.Player{Model: testModel(2), MLBID: 1002, FirstName: "Active", LastName: "Comp", Position: "OF", Active: true}
	retiredComp := models.Player{Model: testModel(3), MLBID: 1003, FirstName: "Retired", LastName: "Comp", Position: "OF", Active: false}

	allPlayers := []models.Player{activeComp, retiredComp}
	allStats := []models.SeasonStat{
		testSeason(4, 2, 2021, 27, 52, 575, 0, 21, 15, 0.284, 0.349, 0.482),
		testSeason(5, 2, 2022, 28, 59, 605, 0, 25, 17, 0.293, 0.360, 0.508),
		testSeason(6, 2, 2023, 29, 64, 620, 0, 27, 16, 0.299, 0.368, 0.526),
		testSeason(7, 2, 2024, 30, 61, 598, 0, 24, 13, 0.289, 0.354, 0.492),
		testSeason(8, 3, 2021, 27, 51, 570, 0, 20, 14, 0.282, 0.347, 0.478),
		testSeason(9, 3, 2022, 28, 58, 600, 0, 24, 16, 0.291, 0.358, 0.503),
		testSeason(10, 3, 2023, 29, 63, 617, 0, 28, 15, 0.298, 0.366, 0.521),
		testSeason(11, 3, 2024, 30, 60, 592, 0, 23, 12, 0.287, 0.351, 0.487),
	}

	result := service.Build(target, targetHistory, allPlayers, allStats)

	require.Equal(t, "ready", result.Status)
	require.Len(t, result.Comparables, 1)
	require.Equal(t, retiredComp.ID, result.Comparables[0].PlayerID)
}

func TestServiceBuildFallsBackWhenComparablePoolIsEmpty(t *testing.T) {
	service := NewService()

	target := models.Player{Model: testModel(1), MLBID: 1001, FirstName: "Solo", LastName: "Player", Position: "SS", Active: true}
	targetHistory := []models.SeasonStat{
		testSeason(1, 1, 2023, 24, 48, 550, 0, 12, 24, 0.281, 0.341, 0.434),
		testSeason(2, 1, 2024, 25, 55, 590, 0, 15, 27, 0.289, 0.352, 0.451),
		testSeason(3, 1, 2025, 26, 58, 608, 0, 16, 29, 0.293, 0.357, 0.459),
	}

	result := service.Build(target, targetHistory, nil, nil)

	require.Equal(t, "ready", result.Status)
	require.Empty(t, result.Comparables)
	require.NotEmpty(t, result.Points)
	require.Len(t, result.Points, len(result.ConfidenceBand))
	for _, band := range result.ConfidenceBand {
		require.Greater(t, band.Upper, band.Lower)
	}
}

func TestServiceBuildReturnsIneligibleForRetiredPlayers(t *testing.T) {
	service := NewService()

	player := models.Player{Model: testModel(1), MLBID: 1001, FirstName: "Retired", LastName: "Player", Position: "1B", Active: false}
	result := service.Build(player, nil, nil, nil)

	require.Equal(t, "ineligible", result.Status)
	require.False(t, result.Eligible)
	require.Equal(t, "player is not active", result.Reason)
	require.Empty(t, result.Points)
}

func testSeason(id uint, playerID int, year, age int, valueScore float64, plateAppearances int, inningsPitched float64, homeRuns, stolenBases int, battingAvg, obp, slg float64) models.SeasonStat {
	return models.SeasonStat{
		Model:            testModel(id),
		PlayerID:         playerID,
		Year:             year,
		Age:              age,
		ValueScore:       valueScore,
		PlateAppearances: plateAppearances,
		InningsPitched:   inningsPitched,
		HomeRuns:         homeRuns,
		StolenBases:      stolenBases,
		BattingAvg:       battingAvg,
		OBP:              obp,
		SLG:              slg,
	}
}

func testModel(id uint) gorm.Model {
	return gorm.Model{ID: id}
}
