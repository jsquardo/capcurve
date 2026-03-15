package scoring

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/jsquardo/capcurve/internal/models"
)

func TestScoreSeasonStatsScopesPercentilesByYear(t *testing.T) {
	stats := []models.SeasonStat{
		newHitterStat(1, 1968, 220, 0.300, 0.380, 0.510, 24, 80, 10, 0.310),
		newHitterStat(2, 1968, 610, 0.270, 0.320, 0.410, 12, 58, 6, 0.280),
		newHitterStat(3, 2000, 220, 0.300, 0.380, 0.510, 24, 80, 10, 0.310),
		newHitterStat(4, 2000, 640, 0.340, 0.440, 0.650, 46, 128, 12, 0.345),
	}

	scores := ScoreSeasonStats(stats)

	require.Greater(t, scores[1].FinalScore, scores[3].FinalScore)
	require.Greater(t, scores[4].FinalScore, scores[3].FinalScore)
}

func TestPercentileHandlesPositiveInverseAndTies(t *testing.T) {
	require.Equal(t, 100.0, percentile(0.400, []float64{0.320, 0.360, 0.400}, false))
	require.Equal(t, 100.0, percentile(2.1, []float64{4.2, 3.0, 2.1}, true))
	require.Equal(t, 75.0, percentile(0.360, []float64{0.320, 0.360, 0.360}, false))
}

func TestHitterScoreDoesNotRequireSavantData(t *testing.T) {
	stats := []models.SeasonStat{
		newHitterStat(1, 2024, 600, 0.310, 0.400, 0.600, 35, 110, 18, 0.320),
		newHitterStat(2, 2024, 580, 0.280, 0.350, 0.500, 24, 88, 9, 0.290),
	}

	scores := ScoreSeasonStats(stats)

	require.Greater(t, scores[1].HitterScore, 0.0)
	require.Greater(t, scores[1].FinalScore, scores[2].FinalScore)
}

func TestPitcherScoreDoesNotRequireExpectedERA(t *testing.T) {
	stats := []models.SeasonStat{
		newPitcherStat(1, 2024, 180, 2.80, 1.02, 11.2, 2.1, 6.8, 0.7, 5.3, 68.0, nil),
		newPitcherStat(2, 2024, 175, 4.10, 1.28, 8.9, 3.0, 8.5, 1.2, 2.9, 63.0, nil),
	}

	scores := ScoreSeasonStats(stats)

	require.Greater(t, scores[1].PitcherScore, 0.0)
	require.Greater(t, scores[1].FinalScore, scores[2].FinalScore)
}

func TestSampleDampenerPenalizesTinySamples(t *testing.T) {
	stats := []models.SeasonStat{
		newHitterStat(1, 2024, 40, 0.350, 0.460, 0.700, 12, 30, 4, 0.360),
		newHitterStat(2, 2024, 600, 0.320, 0.420, 0.620, 30, 95, 12, 0.330),
		newHitterStat(3, 2024, 580, 0.280, 0.350, 0.500, 20, 75, 8, 0.300),
	}

	scores := ScoreSeasonStats(stats)

	require.Less(t, scores[1].FinalScore, scores[2].FinalScore)
}

func TestHitterPercentileCohortExcludesTinySamples(t *testing.T) {
	stats := []models.SeasonStat{
		newHitterStat(1, 2024, 90, 0.410, 0.520, 0.780, 10, 28, 5, 0.390),
		newHitterStat(2, 2024, 610, 0.315, 0.405, 0.615, 38, 112, 11, 0.325),
		newHitterStat(3, 2024, 590, 0.275, 0.345, 0.470, 18, 74, 7, 0.295),
	}

	scores := ScoreSeasonStats(stats)

	require.Equal(t, 100.0, scores[2].HitterScore)
	require.Equal(t, 0.0, scores[3].HitterScore)
	require.Greater(t, scores[1].HitterScore, 0.0)
}

func TestPitcherPercentileCohortExcludesTinySamples(t *testing.T) {
	stats := []models.SeasonStat{
		newPitcherStat(1, 2024, 24, 1.10, 0.82, 13.4, 1.4, 5.0, 0.3, 9.6, 71.0, nil),
		newPitcherStat(2, 2024, 182, 2.95, 1.04, 10.9, 2.1, 6.4, 0.7, 5.2, 67.0, nil),
		newPitcherStat(3, 2024, 168, 4.05, 1.27, 8.4, 3.3, 8.1, 1.1, 2.6, 62.0, nil),
	}

	scores := ScoreSeasonStats(stats)

	require.Equal(t, 100.0, scores[2].PitcherScore)
	require.Equal(t, 0.0, scores[3].PitcherScore)
	require.Greater(t, scores[1].PitcherScore, 0.0)
}

func TestTrueInningsFromBaseballNotation(t *testing.T) {
	require.InDelta(t, 29.6666667, trueInningsFromBaseballNotation(29.2), 0.000001)
	require.InDelta(t, 10.3333333, trueInningsFromBaseballNotation(10.1), 0.000001)
	require.Equal(t, 0.0, trueInningsFromBaseballNotation(0))
	require.Equal(t, 10.4, trueInningsFromBaseballNotation(10.4))
}

func TestPitcherScoreUsesOutsBasedWorkloadDampener(t *testing.T) {
	stats := []models.SeasonStat{
		newPitcherStat(1, 2024, 29.2, 3.20, 1.08, 10.0, 2.4, 6.7, 0.8, 4.2, 67.0, nil),
		newPitcherStat(2, 2024, 29.0, 3.20, 1.08, 10.0, 2.4, 6.7, 0.8, 4.2, 67.0, nil),
		newPitcherStat(3, 2024, 180, 3.20, 1.08, 10.0, 2.4, 6.7, 0.8, 4.2, 67.0, nil),
	}

	scores := ScoreSeasonStats(stats)

	require.Equal(t, 49.44, scores[1].PitcherScore)
	require.Equal(t, 48.33, scores[2].PitcherScore)
}

func TestPartialSavantCoverageOnlyUsesAvailableMetrics(t *testing.T) {
	expectedWOBAHigh := 0.420
	hardHitHigh := 55.0
	expectedWOBALow := 0.330
	hardHitLow := 31.0

	stats := []models.SeasonStat{
		{
			Model:            gorm.Model{ID: 1},
			Year:             2024,
			PlateAppearances: 560,
			HomeRuns:         30,
			RBI:              90,
			StolenBases:      12,
			OBP:              0.370,
			SLG:              0.540,
			OPS:              0.910,
			BABIP:            0.305,
			ExpectedWOBA:     &expectedWOBAHigh,
			HardHitPct:       &hardHitHigh,
		},
		{
			Model:            gorm.Model{ID: 2},
			Year:             2024,
			PlateAppearances: 560,
			HomeRuns:         30,
			RBI:              90,
			StolenBases:      12,
			OBP:              0.370,
			SLG:              0.540,
			OPS:              0.910,
			BABIP:            0.305,
			ExpectedWOBA:     &expectedWOBALow,
			HardHitPct:       &hardHitLow,
		},
	}

	scores := ScoreSeasonStats(stats)

	require.Greater(t, scores[1].HitterScore, scores[2].HitterScore)
}

func TestTwoWayBlendProducesBothSubscoresAndWorkloadWeightedFinal(t *testing.T) {
	expectedWOBA := 0.410
	expectedERA := 3.05
	stats := []models.SeasonStat{
		{
			Model:              gorm.Model{ID: 1},
			Year:               2024,
			PlateAppearances:   520,
			HomeRuns:           38,
			RBI:                95,
			StolenBases:        24,
			OBP:                0.390,
			SLG:                0.610,
			OPS:                1.000,
			BABIP:              0.320,
			ExpectedWOBA:       &expectedWOBA,
			InningsPitched:     135,
			GamesStarted:       24,
			ERA:                3.20,
			WHIP:               1.10,
			StrikeoutsPer9:     11.1,
			WalksPer9:          2.8,
			HitsPer9:           6.9,
			HomeRunsPer9:       0.9,
			StrikeoutWalkRatio: 4.0,
			StrikePercentage:   66.0,
			ExpectedERA:        &expectedERA,
		},
		{
			Model:              gorm.Model{ID: 2},
			Year:               2024,
			PlateAppearances:   610,
			HomeRuns:           42,
			RBI:                118,
			StolenBases:        8,
			OBP:                0.410,
			SLG:                0.640,
			OPS:                1.050,
			BABIP:              0.330,
			InningsPitched:     0,
			ERA:                0,
			WHIP:               0,
			StrikeoutsPer9:     0,
			WalksPer9:          0,
			HitsPer9:           0,
			HomeRunsPer9:       0,
			StrikeoutWalkRatio: 0,
			StrikePercentage:   0,
		},
		{
			Model:              gorm.Model{ID: 3},
			Year:               2024,
			InningsPitched:     185,
			GamesStarted:       31,
			ERA:                2.95,
			WHIP:               1.02,
			StrikeoutsPer9:     10.5,
			WalksPer9:          2.0,
			HitsPer9:           6.3,
			HomeRunsPer9:       0.8,
			StrikeoutWalkRatio: 5.1,
			StrikePercentage:   67.2,
		},
	}

	scores := ScoreSeasonStats(stats)
	twoWay := scores[1]

	require.Greater(t, twoWay.HitterScore, 0.0)
	require.Greater(t, twoWay.PitcherScore, 0.0)
	require.GreaterOrEqual(t, twoWay.FinalScore, math.Min(twoWay.HitterScore, twoWay.PitcherScore))
	require.LessOrEqual(t, twoWay.FinalScore, math.Max(twoWay.HitterScore, twoWay.PitcherScore))
}

func TestTwoWayBlendLetsEligibleSideDominate(t *testing.T) {
	stats := []models.SeasonStat{
		{
			Model:              gorm.Model{ID: 1},
			Year:               2024,
			PlateAppearances:   540,
			HomeRuns:           34,
			RBI:                101,
			StolenBases:        18,
			OBP:                0.385,
			SLG:                0.585,
			OPS:                0.970,
			BABIP:              0.318,
			InningsPitched:     12,
			GamesStarted:       3,
			ERA:                3.05,
			WHIP:               1.08,
			StrikeoutsPer9:     10.8,
			WalksPer9:          2.7,
			HitsPer9:           7.0,
			HomeRunsPer9:       0.8,
			StrikeoutWalkRatio: 4.0,
			StrikePercentage:   65.0,
		},
		newHitterStat(2, 2024, 610, 0.315, 0.405, 0.615, 41, 120, 9, 0.326),
		newHitterStat(4, 2024, 520, 0.265, 0.330, 0.430, 16, 66, 5, 0.288),
		newPitcherStat(3, 2024, 180, 2.90, 1.04, 10.9, 2.1, 6.4, 0.7, 5.2, 67.0, nil),
		newPitcherStat(5, 2024, 165, 4.05, 1.27, 8.4, 3.3, 8.1, 1.1, 2.6, 62.0, nil),
	}

	scores := ScoreSeasonStats(stats)
	twoWay := scores[1]

	require.Greater(t, twoWay.HitterScore, 0.0)
	require.Greater(t, math.Abs(twoWay.FinalScore-twoWay.PitcherScore), math.Abs(twoWay.FinalScore-twoWay.HitterScore))
}

func TestTwoWayBlendUsesOutsBasedPitchingWorkload(t *testing.T) {
	stat := models.SeasonStat{
		PlateAppearances: 100,
		InningsPitched:   29.2,
		GamesStarted:     6,
	}

	score := finalScore(stat, Breakdown{
		HitterScore:  80,
		PitcherScore: 20,
	})

	require.Equal(t, 60.15, score)
}

func newHitterStat(id uint, year int, pa int, avg float64, obp float64, slg float64, homeRuns int, rbi int, steals int, babip float64) models.SeasonStat {
	return models.SeasonStat{
		Model:            gorm.Model{ID: id},
		Year:             year,
		PlateAppearances: pa,
		BattingAvg:       avg,
		OBP:              obp,
		SLG:              slg,
		OPS:              obp + slg,
		HomeRuns:         homeRuns,
		RBI:              rbi,
		StolenBases:      steals,
		BABIP:            babip,
	}
}

func newPitcherStat(id uint, year int, innings float64, era float64, whip float64, k9 float64, bb9 float64, h9 float64, hr9 float64, kbb float64, strikePct float64, expectedERA *float64) models.SeasonStat {
	return models.SeasonStat{
		Model:              gorm.Model{ID: id},
		Year:               year,
		InningsPitched:     innings,
		GamesStarted:       28,
		ERA:                era,
		WHIP:               whip,
		StrikeoutsPer9:     k9,
		WalksPer9:          bb9,
		HitsPer9:           h9,
		HomeRunsPer9:       hr9,
		StrikeoutWalkRatio: kbb,
		StrikePercentage:   strikePct,
		ExpectedERA:        expectedERA,
	}
}
