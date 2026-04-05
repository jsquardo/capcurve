package projection

import (
	"math"
	"sort"
	"strings"

	"github.com/jsquardo/capcurve/internal/models"
)

const (
	defaultComparableLimit     = 5
	defaultComparableThreshold = 18.0
	defaultProjectionHorizon   = 5
	defaultProjectionFloor     = 5.0
)

type Result struct {
	Status         string
	Eligible       bool
	Reason         string
	Points         []Point
	ConfidenceBand []ConfidenceBand
	Comparables    []Comparable
}

type Point struct {
	Year       int
	Age        int
	ValueScore float64
}

type ConfidenceBand struct {
	Year  int
	Lower float64
	Upper float64
}

type Comparable struct {
	PlayerID uint
	MLBID    int
	FullName string
	Position string
}

type Service struct {
	comparableLimit     int
	comparableThreshold float64
	projectionHorizon   int
	projectionFloor     float64
}

type roleProfile struct {
	broadRole      string
	profile        string
	peakAge        int
	prePeakRise    float64
	peakDrift      float64
	declineEarly   float64
	declineLate    float64
	positionBucket string
}

type comparableMatch struct {
	player  models.Player
	score   float64
	futures []models.SeasonStat
}

func NewService() *Service {
	return &Service{
		comparableLimit:     defaultComparableLimit,
		comparableThreshold: defaultComparableThreshold,
		projectionHorizon:   defaultProjectionHorizon,
		projectionFloor:     defaultProjectionFloor,
	}
}

func (s *Service) Build(player models.Player, history []models.SeasonStat, players []models.Player, stats []models.SeasonStat) Result {
	if !player.Active {
		return Result{
			Status:         "ineligible",
			Eligible:       false,
			Reason:         "player is not active",
			Points:         []Point{},
			ConfidenceBand: []ConfidenceBand{},
			Comparables:    []Comparable{},
		}
	}

	if len(history) == 0 {
		return Result{
			Status:         "insufficient_data",
			Eligible:       true,
			Reason:         "player does not have any season history to project from",
			Points:         []Point{},
			ConfidenceBand: []ConfidenceBand{},
			Comparables:    []Comparable{},
		}
	}

	// Use only qualified seasons (value_score > 0) as the baseline for role
	// inference, comparable matching, and confidence-band volatility. Sub-threshold
	// or in-progress seasons score 0 and would otherwise distort all three paths
	// (e.g. a trailing in-progress season with no pitching stats would wrongly
	// classify a starter as a hitter, misalign the comparable anchor age by 1+
	// years, and inflate fallback volatility against real-season deltas).
	qualifiedHistory := filterQualifiedSeasons(history)
	if len(qualifiedHistory) == 0 {
		return Result{
			Status:         "ready",
			Eligible:       true,
			Reason:         "",
			Points:         []Point{},
			ConfidenceBand: []ConfidenceBand{},
			Comparables:    []Comparable{},
		}
	}

	latest := qualifiedHistory[len(qualifiedHistory)-1]
	profile := inferRoleProfile(player, latest)
	points := s.projectPoints(history, profile)
	matches := s.findComparableMatches(player, qualifiedHistory, profile, players, stats)
	bands := s.buildConfidenceBands(points, qualifiedHistory, matches)

	comparables := make([]Comparable, 0, len(matches))
	for _, match := range matches {
		comparables = append(comparables, Comparable{
			PlayerID: match.player.ID,
			MLBID:    match.player.MLBID,
			FullName: strings.TrimSpace(match.player.FirstName + " " + match.player.LastName),
			Position: match.player.Position,
		})
	}

	return Result{
		Status:         "ready",
		Eligible:       true,
		Reason:         "",
		Points:         points,
		ConfidenceBand: bands,
		Comparables:    comparables,
	}
}

func (s *Service) projectPoints(history []models.SeasonStat, profile roleProfile) []Point {
	// Only seasons with a real score are useful baselines. Sub-threshold or
	// incomplete seasons (e.g. an in-progress season below 100 PA) have a
	// value_score of 0, which would anchor the projection at zero.
	qualifiedHistory := filterQualifiedSeasons(history)
	if len(qualifiedHistory) == 0 {
		return []Point{}
	}
	latest := qualifiedHistory[len(qualifiedHistory)-1]
	weightedRecent := weightedRecentAverage(qualifiedHistory)
	trend := weightedTrend(qualifiedHistory)

	points := make([]Point, 0, s.projectionHorizon)
	for horizon := 1; horizon <= s.projectionHorizon; horizon++ {
		targetYear := latest.Year + horizon
		targetAge := latest.Age + horizon

		ageCurveValue := latest.ValueScore + cumulativeAgeAdjustment(profile, latest.Age, targetAge)
		recentComponent := weightedRecent + (trend * trendDecayForYear(horizon))
		blendRecent, blendAge := blendWeightsForYear(horizon)
		value := clampScore((recentComponent * blendRecent) + (ageCurveValue * blendAge))

		points = append(points, Point{
			Year:       targetYear,
			Age:        targetAge,
			ValueScore: roundToTenth(value),
		})

		if value <= s.projectionFloor {
			break
		}
	}

	return points
}

func (s *Service) findComparableMatches(target models.Player, targetHistory []models.SeasonStat, targetProfile roleProfile, players []models.Player, stats []models.SeasonStat) []comparableMatch {
	if len(targetHistory) == 0 {
		return nil
	}

	historyByPlayer := make(map[int][]models.SeasonStat)
	for _, stat := range stats {
		historyByPlayer[stat.PlayerID] = append(historyByPlayer[stat.PlayerID], stat)
	}

	for playerID := range historyByPlayer {
		sort.Slice(historyByPlayer[playerID], func(i, j int) bool {
			if historyByPlayer[playerID][i].Year == historyByPlayer[playerID][j].Year {
				return historyByPlayer[playerID][i].ID < historyByPlayer[playerID][j].ID
			}
			return historyByPlayer[playerID][i].Year < historyByPlayer[playerID][j].Year
		})
	}

	targetLatest := targetHistory[len(targetHistory)-1]
	matches := make([]comparableMatch, 0, s.comparableLimit)
	for _, candidate := range players {
		if candidate.ID == target.ID {
			continue
		}

		// Confidence bands need completed future outcomes, so only retired or
		// otherwise historical players are eligible as comparables.
		if candidate.Active {
			continue
		}

		candidateHistory := historyByPlayer[int(candidate.ID)]
		if len(candidateHistory) < 2 {
			continue
		}

		anchorIndex, ok := findComparableAnchor(targetLatest.Age, candidateHistory)
		if !ok || anchorIndex == len(candidateHistory)-1 {
			continue
		}

		candidateProfile := inferRoleProfile(candidate, candidateHistory[anchorIndex])
		if candidateProfile.broadRole != targetProfile.broadRole {
			continue
		}

		score := comparableDistance(target, targetHistory, targetProfile, candidate, candidateHistory, candidateProfile, anchorIndex)
		if score > s.comparableThreshold {
			continue
		}

		matches = append(matches, comparableMatch{
			player:  candidate,
			score:   score,
			futures: candidateHistory[anchorIndex+1:],
		})
	}

	sort.Slice(matches, func(i, j int) bool {
		if matches[i].score == matches[j].score {
			return matches[i].player.ID < matches[j].player.ID
		}
		return matches[i].score < matches[j].score
	})

	if len(matches) > s.comparableLimit {
		matches = matches[:s.comparableLimit]
	}

	return matches
}

func (s *Service) buildConfidenceBands(points []Point, history []models.SeasonStat, matches []comparableMatch) []ConfidenceBand {
	bands := make([]ConfidenceBand, 0, len(points))
	previousWidth := 0.0
	volatilityFallback := maxFloat(recentVolatility(history), 4.0)

	for index, point := range points {
		scores := futureComparableScores(matches, index)
		width := 0.0

		switch {
		case len(scores) >= 2:
			width = stddev(scores) * horizonFactor(index+1)
		case previousWidth > 0:
			width = previousWidth * horizonFactor(index+1)
		default:
			width = volatilityFallback * horizonFactor(index+1)
		}

		width = maxFloat(width, 2.5)
		previousWidth = width

		bands = append(bands, ConfidenceBand{
			Year:  point.Year,
			Lower: roundToTenth(clampScore(point.ValueScore - width)),
			Upper: roundToTenth(clampScore(point.ValueScore + width)),
		})
	}

	return bands
}

func inferRoleProfile(player models.Player, latest models.SeasonStat) roleProfile {
	hasHitting := latest.PlateAppearances > 0
	hasPitching := latest.InningsPitched > 0

	switch {
	case hasHitting && hasPitching:
		return roleProfile{
			broadRole:      "two_way",
			profile:        "two_way",
			peakAge:        29,
			prePeakRise:    1.5,
			peakDrift:      -0.2,
			declineEarly:   2.1,
			declineLate:    3.3,
			positionBucket: strings.ToUpper(strings.TrimSpace(player.Position)),
		}
	case hasPitching:
		if isStarter(latest) {
			return roleProfile{
				broadRole:      "pitcher",
				profile:        "starter",
				peakAge:        30,
				prePeakRise:    1.0,
				peakDrift:      -0.3,
				declineEarly:   1.8,
				declineLate:    2.9,
				positionBucket: "SP",
			}
		}
		return roleProfile{
			broadRole:      "pitcher",
			profile:        "reliever",
			peakAge:        27,
			prePeakRise:    0.8,
			peakDrift:      -0.4,
			declineEarly:   2.0,
			declineLate:    3.2,
			positionBucket: "RP",
		}
	default:
		position := strings.ToUpper(strings.TrimSpace(player.Position))
		switch {
		case position == "C":
			return roleProfile{
				broadRole:      "hitter",
				profile:        "catcher",
				peakAge:        26,
				prePeakRise:    1.4,
				peakDrift:      -0.5,
				declineEarly:   2.8,
				declineLate:    4.0,
				positionBucket: "C",
			}
		case position == "DH" || position == "1B":
			return roleProfile{
				broadRole:      "hitter",
				profile:        "corner_bat",
				peakAge:        32,
				prePeakRise:    0.9,
				peakDrift:      -0.1,
				declineEarly:   1.4,
				declineLate:    2.4,
				positionBucket: position,
			}
		case latest.HomeRuns >= 30 || latest.SLG >= 0.500:
			return roleProfile{
				broadRole:      "hitter",
				profile:        "power",
				peakAge:        29,
				prePeakRise:    1.2,
				peakDrift:      -0.2,
				declineEarly:   1.8,
				declineLate:    2.7,
				positionBucket: position,
			}
		case latest.StolenBases >= 20 || (latest.BattingAvg >= 0.285 && latest.HomeRuns < 20):
			return roleProfile{
				broadRole:      "hitter",
				profile:        "speed_contact",
				peakAge:        26,
				prePeakRise:    1.8,
				peakDrift:      -0.4,
				declineEarly:   2.3,
				declineLate:    3.5,
				positionBucket: position,
			}
		default:
			return roleProfile{
				broadRole:      "hitter",
				profile:        "default_hitter",
				peakAge:        27,
				prePeakRise:    1.3,
				peakDrift:      -0.3,
				declineEarly:   2.0,
				declineLate:    3.0,
				positionBucket: position,
			}
		}
	}
}

func isStarter(stat models.SeasonStat) bool {
	if stat.GamesStarted >= 10 {
		return true
	}
	if stat.GamesPlayed == 0 {
		return false
	}
	return float64(stat.GamesStarted)/float64(stat.GamesPlayed) >= 0.35
}

func weightedRecentAverage(history []models.SeasonStat) float64 {
	weights := []float64{1.0, 0.6, 0.36}
	used := lastNSeasons(history, len(weights))

	var weightedSum float64
	var totalWeight float64
	for index := range used {
		weight := weights[index]
		season := used[len(used)-1-index]
		weightedSum += season.ValueScore * weight
		totalWeight += weight
	}

	if totalWeight == 0 {
		return 0
	}

	return weightedSum / totalWeight
}

func weightedTrend(history []models.SeasonStat) float64 {
	weights := []float64{1.0, 0.6}
	used := lastNSeasons(history, 3)
	if len(used) < 2 {
		return 0
	}

	var weightedSum float64
	var totalWeight float64
	for index := len(used) - 1; index > 0; index-- {
		weightIndex := len(used) - 1 - index
		weight := weights[minInt(weightIndex, len(weights)-1)]
		delta := used[index].ValueScore - used[index-1].ValueScore
		weightedSum += delta * weight
		totalWeight += weight
	}

	if totalWeight == 0 {
		return 0
	}

	return weightedSum / totalWeight
}

func trendDecayForYear(year int) float64 {
	switch year {
	case 1:
		return 1.0
	case 2:
		return 0.6
	case 3:
		return 0.35
	default:
		return 0.2
	}
}

func blendWeightsForYear(year int) (recent float64, age float64) {
	switch year {
	case 1:
		return 0.65, 0.35
	case 2:
		return 0.50, 0.50
	case 3:
		return 0.35, 0.65
	default:
		return 0.25, 0.75
	}
}

func cumulativeAgeAdjustment(profile roleProfile, currentAge, targetAge int) float64 {
	if targetAge <= currentAge {
		return 0
	}

	var total float64
	for age := currentAge + 1; age <= targetAge; age++ {
		total += yearlyAgeAdjustment(profile, age)
	}

	return total
}

func yearlyAgeAdjustment(profile roleProfile, age int) float64 {
	switch {
	case age < profile.peakAge:
		return profile.prePeakRise
	case age <= profile.peakAge+1:
		return profile.peakDrift
	case age <= profile.peakAge+3:
		return -profile.declineEarly
	default:
		return -profile.declineLate
	}
}

func findComparableAnchor(targetAge int, history []models.SeasonStat) (int, bool) {
	bestIndex := -1
	bestDistance := math.MaxInt

	for index, season := range history {
		if season.Age == 0 {
			continue
		}

		distance := absInt(season.Age - targetAge)
		if distance > 2 {
			continue
		}

		if distance < bestDistance || (distance == bestDistance && season.Year > history[bestIndex].Year) {
			bestIndex = index
			bestDistance = distance
		}
	}

	return bestIndex, bestIndex >= 0
}

func comparableDistance(target models.Player, targetHistory []models.SeasonStat, targetProfile roleProfile, candidate models.Player, candidateHistory []models.SeasonStat, candidateProfile roleProfile, anchorIndex int) float64 {
	targetLatest := targetHistory[len(targetHistory)-1]
	anchor := candidateHistory[anchorIndex]

	positionPenalty := 0.0
	if strings.EqualFold(target.Position, candidate.Position) {
		positionPenalty = 0
	} else if targetProfile.positionBucket == candidateProfile.positionBucket {
		positionPenalty = 1.0
	} else {
		positionPenalty = 5.0
	}

	profilePenalty := 0.0
	if targetProfile.profile != candidateProfile.profile {
		profilePenalty = 4.0
	}

	agePenalty := float64(absInt(targetLatest.Age-anchor.Age)) * 2.0
	scorePenalty := math.Abs(targetLatest.ValueScore-anchor.ValueScore) / 5.0
	shapePenalty := trajectoryPenalty(targetHistory, candidateHistory[:anchorIndex+1])

	return positionPenalty + profilePenalty + agePenalty + scorePenalty + shapePenalty
}

func trajectoryPenalty(targetHistory []models.SeasonStat, candidateHistory []models.SeasonStat) float64 {
	targetWindow := lastNSeasons(targetHistory, 3)
	candidateWindow := lastNSeasons(candidateHistory, 3)
	length := minInt(len(targetWindow), len(candidateWindow))
	if length == 0 {
		return 6.0
	}

	var total float64
	for i := 0; i < length; i++ {
		targetSeason := targetWindow[len(targetWindow)-length+i]
		candidateSeason := candidateWindow[len(candidateWindow)-length+i]
		total += math.Abs(targetSeason.ValueScore-candidateSeason.ValueScore) / 8.0
	}

	if length >= 2 {
		targetDelta := targetWindow[len(targetWindow)-1].ValueScore - targetWindow[len(targetWindow)-2].ValueScore
		candidateDelta := candidateWindow[len(candidateWindow)-1].ValueScore - candidateWindow[len(candidateWindow)-2].ValueScore
		total += math.Abs(targetDelta-candidateDelta) / 4.0
	}

	return total / float64(length)
}

func futureComparableScores(matches []comparableMatch, horizonIndex int) []float64 {
	scores := make([]float64, 0, len(matches))
	for _, match := range matches {
		if horizonIndex >= len(match.futures) {
			continue
		}
		scores = append(scores, match.futures[horizonIndex].ValueScore)
	}
	return scores
}

func recentVolatility(history []models.SeasonStat) float64 {
	used := lastNSeasons(history, 3)
	if len(used) < 2 {
		return 4.0
	}

	deltas := make([]float64, 0, len(used)-1)
	for i := 1; i < len(used); i++ {
		deltas = append(deltas, math.Abs(used[i].ValueScore-used[i-1].ValueScore))
	}

	return mean(deltas)
}

func horizonFactor(year int) float64 {
	return 1.0 + (float64(year-1) * 0.10)
}

func stddev(values []float64) float64 {
	if len(values) < 2 {
		return 0
	}

	avg := mean(values)
	var sum float64
	for _, value := range values {
		sum += math.Pow(value-avg, 2)
	}

	return math.Sqrt(sum / float64(len(values)-1))
}

func mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	var total float64
	for _, value := range values {
		total += value
	}

	return total / float64(len(values))
}

func filterQualifiedSeasons(history []models.SeasonStat) []models.SeasonStat {
	qualified := make([]models.SeasonStat, 0, len(history))
	for _, s := range history {
		if s.ValueScore > 0 {
			qualified = append(qualified, s)
		}
	}
	return qualified
}

func lastNSeasons(history []models.SeasonStat, count int) []models.SeasonStat {
	if len(history) <= count {
		return history
	}
	return history[len(history)-count:]
}

func clampScore(value float64) float64 {
	return math.Max(0, math.Min(100, value))
}

func roundToTenth(value float64) float64 {
	return math.Round(value*10) / 10
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
