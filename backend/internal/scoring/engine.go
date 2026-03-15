package scoring

import (
	"math"
	"sort"

	"github.com/jsquardo/capcurve/internal/models"
)

const (
	hitterThreshold  = 100.0
	pitcherThreshold = 30.0
)

type Breakdown struct {
	HitterScore  float64
	PitcherScore float64
	FinalScore   float64
}

type metricConfig struct {
	weight  float64
	inverse bool
	value   func(models.SeasonStat) (float64, bool)
}

var hitterMetricConfigs = []metricConfig{
	{weight: 0.28, value: func(stat models.SeasonStat) (float64, bool) { return stat.OBP, true }},
	{weight: 0.24, value: func(stat models.SeasonStat) (float64, bool) { return stat.SLG, true }},
	{weight: 0.08, value: func(stat models.SeasonStat) (float64, bool) { return stat.BattingAvg, true }},
	{weight: 0.10, value: func(stat models.SeasonStat) (float64, bool) { return float64(stat.HomeRuns), true }},
	{weight: 0.07, value: func(stat models.SeasonStat) (float64, bool) { return float64(stat.RBI), true }},
	{weight: 0.05, value: func(stat models.SeasonStat) (float64, bool) { return float64(stat.StolenBases), true }},
	{weight: 0.05, value: func(stat models.SeasonStat) (float64, bool) { return stat.BABIP, true }},
}

var hitterSavantConfigs = []metricConfig{
	{weight: 1.20, value: func(stat models.SeasonStat) (float64, bool) { return optionalValue(stat.ExpectedWOBA) }},
	{weight: 0.80, value: func(stat models.SeasonStat) (float64, bool) { return optionalValue(stat.ExpectedBattingAvg) }},
	{weight: 1.00, value: func(stat models.SeasonStat) (float64, bool) { return optionalValue(stat.ExpectedSlugging) }},
	{weight: 1.40, value: func(stat models.SeasonStat) (float64, bool) { return optionalValue(stat.BarrelPct) }},
	{weight: 1.20, value: func(stat models.SeasonStat) (float64, bool) { return optionalValue(stat.HardHitPct) }},
	{weight: 0.90, value: func(stat models.SeasonStat) (float64, bool) { return optionalValue(stat.AvgExitVelocity) }},
	{weight: 0.60, value: func(stat models.SeasonStat) (float64, bool) { return optionalValue(stat.SweetSpotPct) }},
}

var pitcherMetricConfigs = []metricConfig{
	{weight: 0.22, inverse: true, value: func(stat models.SeasonStat) (float64, bool) { return stat.ERA, true }},
	{weight: 0.18, inverse: true, value: func(stat models.SeasonStat) (float64, bool) { return stat.WHIP, true }},
	{weight: 0.16, value: func(stat models.SeasonStat) (float64, bool) { return stat.StrikeoutsPer9, true }},
	{weight: 0.10, inverse: true, value: func(stat models.SeasonStat) (float64, bool) { return stat.WalksPer9, true }},
	{weight: 0.08, inverse: true, value: func(stat models.SeasonStat) (float64, bool) { return stat.HitsPer9, true }},
	{weight: 0.08, inverse: true, value: func(stat models.SeasonStat) (float64, bool) { return stat.HomeRunsPer9, true }},
	{weight: 0.10, value: func(stat models.SeasonStat) (float64, bool) { return stat.StrikeoutWalkRatio, true }},
	{weight: 0.04, value: func(stat models.SeasonStat) (float64, bool) { return stat.StrikePercentage, true }},
	{weight: 0.04, inverse: true, value: func(stat models.SeasonStat) (float64, bool) { return optionalValue(stat.ExpectedERA) }},
}

// ScoreSeasonStats computes 0-100 season-scoped scores grouped by MLB season so
// a stat line is only compared to peers from the same run environment.
func ScoreSeasonStats(stats []models.SeasonStat) map[uint]Breakdown {
	byYear := make(map[int][]models.SeasonStat)
	for _, stat := range stats {
		byYear[stat.Year] = append(byYear[stat.Year], stat)
	}

	scores := make(map[uint]Breakdown, len(stats))
	for _, yearStats := range byYear {
		scoreYear(yearStats, scores)
	}

	return scores
}

func scoreYear(stats []models.SeasonStat, scores map[uint]Breakdown) {
	hitterCohort := filterStats(stats, isEligibleHitterPercentileSeason)
	pitcherCohort := filterStats(stats, isEligiblePitcherPercentileSeason)

	for _, stat := range stats {
		breakdown := Breakdown{}

		if isHitterSeason(stat) {
			breakdown.HitterScore = hitterScore(stat, hitterCohort)
		}
		if isPitcherSeason(stat) {
			breakdown.PitcherScore = pitcherScore(stat, pitcherCohort)
		}

		breakdown.FinalScore = finalScore(stat, breakdown)
		scores[stat.ID] = breakdown
	}
}

func hitterScore(stat models.SeasonStat, cohort []models.SeasonStat) float64 {
	score, _ := weightedPercentileScore(stat, cohort, hitterMetricConfigs)
	if savantScore, ok := savantBlendScore(stat, cohort, hitterSavantConfigs); ok {
		score = mergeWeightedScores(score, 0.87, savantScore, 0.13)
	}

	return roundScore(score * sampleDampener(float64(stat.PlateAppearances), hitterThreshold))
}

func pitcherScore(stat models.SeasonStat, cohort []models.SeasonStat) float64 {
	score, _ := weightedPercentileScore(stat, cohort, pitcherMetricConfigs)
	return roundScore(score * sampleDampener(pitchingWorkloadInnings(stat), pitcherThreshold))
}

func finalScore(stat models.SeasonStat, breakdown Breakdown) float64 {
	hitterActive := isHitterSeason(stat)
	pitcherActive := isPitcherSeason(stat)

	switch {
	case hitterActive && pitcherActive:
		hitterWeight := roleWorkload(float64(stat.PlateAppearances), hitterThreshold)
		pitcherWeight := roleWorkload(pitchingWorkloadInnings(stat), pitcherThreshold)

		hitterEligible := float64(stat.PlateAppearances) >= hitterThreshold
		pitcherEligible := pitchingWorkloadInnings(stat) >= pitcherThreshold
		if hitterEligible != pitcherEligible {
			if hitterEligible {
				hitterWeight *= 2
			}
			if pitcherEligible {
				pitcherWeight *= 2
			}
		}

		totalWeight := hitterWeight + pitcherWeight
		if totalWeight == 0 {
			return roundScore((breakdown.HitterScore + breakdown.PitcherScore) / 2)
		}

		return roundScore(((breakdown.HitterScore * hitterWeight) + (breakdown.PitcherScore * pitcherWeight)) / totalWeight)
	case hitterActive:
		return breakdown.HitterScore
	case pitcherActive:
		return breakdown.PitcherScore
	default:
		return 0
	}
}

func weightedPercentileScore(stat models.SeasonStat, cohort []models.SeasonStat, configs []metricConfig) (float64, bool) {
	var weightedTotal float64
	var totalWeight float64

	for _, config := range configs {
		value, ok := config.value(stat)
		if !ok {
			continue
		}

		peerValues := peerValues(cohort, config.value)
		if len(peerValues) == 0 {
			continue
		}

		weightedTotal += percentile(value, peerValues, config.inverse) * config.weight
		totalWeight += config.weight
	}

	if totalWeight == 0 {
		return 0, false
	}

	return weightedTotal / totalWeight, true
}

func savantBlendScore(stat models.SeasonStat, cohort []models.SeasonStat, configs []metricConfig) (float64, bool) {
	return weightedPercentileScore(stat, cohort, configs)
}

func peerValues(cohort []models.SeasonStat, extractor func(models.SeasonStat) (float64, bool)) []float64 {
	values := make([]float64, 0, len(cohort))
	for _, candidate := range cohort {
		value, ok := extractor(candidate)
		if !ok {
			continue
		}
		values = append(values, value)
	}
	return values
}

func percentile(value float64, peers []float64, inverse bool) float64 {
	if len(peers) == 0 {
		return 0
	}
	if len(peers) == 1 {
		return 50
	}

	normalizedTarget := normalizeForOrdering(value, inverse)
	normalizedPeers := make([]float64, len(peers))
	for i, peer := range peers {
		normalizedPeers[i] = normalizeForOrdering(peer, inverse)
	}
	sort.Float64s(normalizedPeers)

	lower := sort.Search(len(normalizedPeers), func(i int) bool {
		return normalizedPeers[i] >= normalizedTarget
	})
	upper := sort.Search(len(normalizedPeers), func(i int) bool {
		return normalizedPeers[i] > normalizedTarget
	})

	rank := float64(lower) + (float64(upper-lower)-1)/2
	return (rank / float64(len(normalizedPeers)-1)) * 100
}

func filterStats(stats []models.SeasonStat, predicate func(models.SeasonStat) bool) []models.SeasonStat {
	filtered := make([]models.SeasonStat, 0, len(stats))
	for _, stat := range stats {
		if predicate(stat) {
			filtered = append(filtered, stat)
		}
	}
	return filtered
}

func isHitterSeason(stat models.SeasonStat) bool {
	return stat.PlateAppearances > 0
}

func isPitcherSeason(stat models.SeasonStat) bool {
	return stat.InningsPitched > 0 || stat.GamesStarted > 0
}

func isEligibleHitterPercentileSeason(stat models.SeasonStat) bool {
	return float64(stat.PlateAppearances) >= hitterThreshold
}

func isEligiblePitcherPercentileSeason(stat models.SeasonStat) bool {
	return pitchingWorkloadInnings(stat) >= pitcherThreshold
}

func sampleDampener(sample float64, threshold float64) float64 {
	if sample <= 0 {
		return 0
	}
	return math.Min(sample/threshold, 1)
}

func roleWorkload(sample float64, threshold float64) float64 {
	if sample <= 0 {
		return 0
	}
	return math.Min(sample/threshold, 1)
}

// season_stats.innings_pitched is still stored in MLB baseball notation, so
// scoring must convert values like 29.2 back into 29 2/3 true innings before
// applying workload thresholds, dampeners, or two-way role weighting.
func pitchingWorkloadInnings(stat models.SeasonStat) float64 {
	return trueInningsFromBaseballNotation(stat.InningsPitched)
}

func trueInningsFromBaseballNotation(innings float64) float64 {
	if innings <= 0 {
		return 0
	}

	tenths := int(math.Round(innings * 10))
	whole := tenths / 10
	partial := tenths % 10

	switch partial {
	case 0:
		return float64(whole)
	case 1, 2:
		return float64((whole*3)+partial) / 3
	default:
		return innings
	}
}

func mergeWeightedScores(primary float64, primaryWeight float64, secondary float64, secondaryWeight float64) float64 {
	return ((primary * primaryWeight) + (secondary * secondaryWeight)) / (primaryWeight + secondaryWeight)
}

func optionalValue(value *float64) (float64, bool) {
	if value == nil {
		return 0, false
	}
	return *value, true
}

func normalizeForOrdering(value float64, inverse bool) float64 {
	if inverse {
		return -value
	}
	return value
}

func roundScore(value float64) float64 {
	return math.Round(value*100) / 100
}
