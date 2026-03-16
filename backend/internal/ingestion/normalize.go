package ingestion

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

const mlbHeadshotURLTemplate = "https://img.mlbstatic.com/mlb-photos/image/upload/w_360,q_auto:best/v1/people/%d/headshot/67/current"

// NormalizePlayer converts the MLB person payload into the shape CapCurve stores.
func NormalizePlayer(player *MLBPlayer) (*PlayerRecord, error) {
	if player == nil {
		return nil, fmt.Errorf("player is nil")
	}

	var dateOfBirth *time.Time
	if player.BirthDate != "" {
		parsed, err := time.Parse("2006-01-02", player.BirthDate)
		if err != nil {
			return nil, fmt.Errorf("parse birth date: %w", err)
		}
		dateOfBirth = &parsed
	}

	return &PlayerRecord{
		MLBID:       player.ID,
		FirstName:   player.FirstName,
		LastName:    player.LastName,
		Position:    player.PrimaryPosition.Name,
		Bats:        player.BatSide.Code,
		Throws:      player.PitchHand.Code,
		DateOfBirth: dateOfBirth,
		Active:      player.Active,
		ImageURL:    fmt.Sprintf(mlbHeadshotURLTemplate, player.ID),
	}, nil
}

// NormalizeSeasonSplit maps either a hitting or pitching split into the unified
// season record used by the database.
func NormalizeSeasonSplit(split MLBSeasonSplit, group string) (SeasonStatRecord, error) {
	record := SeasonStatRecord{
		Year:     parseStringInt(split.Season),
		TeamID:   split.Team.ID,
		TeamName: split.Team.Name,
	}

	switch group {
	case "hitting":
		record.HasHitting = true
		record.Age = intFromMap(split.Stat, "age")
		record.GamesPlayed = intFromMap(split.Stat, "gamesPlayed")
		record.PlateAppearances = intFromMap(split.Stat, "plateAppearances")
		record.AtBats = intFromMap(split.Stat, "atBats")
		record.Hits = intFromMap(split.Stat, "hits")
		record.Doubles = intFromMap(split.Stat, "doubles")
		record.Triples = intFromMap(split.Stat, "triples")
		record.HomeRuns = intFromMap(split.Stat, "homeRuns")
		record.Runs = intFromMap(split.Stat, "runs")
		record.RBI = intFromMap(split.Stat, "rbi")
		record.Walks = intFromMap(split.Stat, "baseOnBalls")
		record.Strikeouts = intFromMap(split.Stat, "strikeOuts")
		record.StolenBases = intFromMap(split.Stat, "stolenBases")
		record.BattingAvg = floatFromMap(split.Stat, "avg")
		record.OBP = floatFromMap(split.Stat, "obp")
		record.SLG = floatFromMap(split.Stat, "slg")
		record.OPS = floatFromMap(split.Stat, "ops")
		record.BABIP = floatFromMap(split.Stat, "babip")
	case "pitching":
		record.HasPitching = true
		record.Age = intFromMap(split.Stat, "age")
		record.GamesPlayed = intFromMap(split.Stat, "gamesPlayed")
		record.GamesStarted = intFromMap(split.Stat, "gamesStarted")
		record.Wins = intFromMap(split.Stat, "wins")
		record.Losses = intFromMap(split.Stat, "losses")
		record.ERA = floatFromMap(split.Stat, "era")
		record.WHIP = floatFromMap(split.Stat, "whip")
		inningsOuts, err := baseballInningsOutsFromMap(split.Stat, "inningsPitched")
		if err != nil {
			return SeasonStatRecord{}, err
		}
		record.InningsPitchedOuts = inningsOuts
		record.InningsPitched = baseballInningsFromOuts(inningsOuts)
		record.HitsAllowed = intFromMap(split.Stat, "hits")
		record.WalksAllowed = intFromMap(split.Stat, "baseOnBalls")
		record.HomeRunsAllowed = intFromMap(split.Stat, "homeRuns")
		record.Strikeouts = intFromMap(split.Stat, "strikeOuts")
		record.StrikeoutsPer9 = floatFromMap(split.Stat, "strikeoutsPer9Inn")
		record.WalksPer9 = floatFromMap(split.Stat, "walksPer9Inn")
		record.HitsPer9 = floatFromMap(split.Stat, "hitsPer9Inn")
		record.HomeRunsPer9 = floatFromMap(split.Stat, "homeRunsPer9")
		record.StrikeoutWalkRatio = floatFromMap(split.Stat, "strikeoutWalkRatio")
		record.StrikePercentage = floatFromMap(split.Stat, "strikePercentage")
	default:
		return SeasonStatRecord{}, fmt.Errorf("unsupported stat group %q", group)
	}

	return record, nil
}

// IsAggregateSeasonSplit identifies MLB's synthetic "TOT" row so ingestion can
// keep only real team splits in storage.
func IsAggregateSeasonSplit(split MLBSeasonSplit) bool {
	return split.Team.ID == 0
}

// MergeSeasonStats overlays hitting, pitching, and enrichment data into a single
// row keyed by player and season.
func MergeSeasonStats(existing SeasonStatRecord, incoming SeasonStatRecord) SeasonStatRecord {
	merged := existing

	if merged.Year == 0 {
		merged.Year = incoming.Year
	}
	if merged.TeamID == 0 {
		merged.TeamID = incoming.TeamID
	}
	if merged.TeamName == "" {
		merged.TeamName = incoming.TeamName
	}
	if merged.Age == 0 {
		merged.Age = incoming.Age
	}
	if incoming.HasHitting {
		merged.HasHitting = true
	}
	if incoming.HasPitching {
		merged.HasPitching = true
	}

	mergeInt := func(current *int, next int) {
		if next != 0 {
			*current = next
		}
	}
	mergeFloat := func(current *float64, next float64) {
		if next != 0 {
			*current = next
		}
	}
	mergeOptional := func(current **float64, next *float64) {
		if next != nil {
			*current = next
		}
	}

	mergeInt(&merged.GamesPlayed, incoming.GamesPlayed)
	mergeInt(&merged.GamesStarted, incoming.GamesStarted)
	mergeInt(&merged.PlateAppearances, incoming.PlateAppearances)
	mergeInt(&merged.AtBats, incoming.AtBats)
	mergeInt(&merged.Hits, incoming.Hits)
	mergeInt(&merged.Doubles, incoming.Doubles)
	mergeInt(&merged.Triples, incoming.Triples)
	mergeInt(&merged.HomeRuns, incoming.HomeRuns)
	mergeInt(&merged.Runs, incoming.Runs)
	mergeInt(&merged.RBI, incoming.RBI)
	mergeInt(&merged.Walks, incoming.Walks)
	mergeInt(&merged.Strikeouts, incoming.Strikeouts)
	mergeInt(&merged.StolenBases, incoming.StolenBases)
	mergeFloat(&merged.BattingAvg, incoming.BattingAvg)
	mergeFloat(&merged.OBP, incoming.OBP)
	mergeFloat(&merged.SLG, incoming.SLG)
	mergeFloat(&merged.OPS, incoming.OPS)
	mergeFloat(&merged.BABIP, incoming.BABIP)
	mergeInt(&merged.Wins, incoming.Wins)
	mergeInt(&merged.Losses, incoming.Losses)
	mergeFloat(&merged.ERA, incoming.ERA)
	mergeFloat(&merged.WHIP, incoming.WHIP)
	if incoming.InningsPitchedOuts != 0 {
		merged.InningsPitchedOuts = incoming.InningsPitchedOuts
		merged.InningsPitched = baseballInningsFromOuts(merged.InningsPitchedOuts)
	}
	mergeInt(&merged.HitsAllowed, incoming.HitsAllowed)
	mergeInt(&merged.WalksAllowed, incoming.WalksAllowed)
	mergeInt(&merged.HomeRunsAllowed, incoming.HomeRunsAllowed)
	mergeFloat(&merged.StrikeoutsPer9, incoming.StrikeoutsPer9)
	mergeFloat(&merged.WalksPer9, incoming.WalksPer9)
	mergeFloat(&merged.HitsPer9, incoming.HitsPer9)
	mergeFloat(&merged.HomeRunsPer9, incoming.HomeRunsPer9)
	mergeFloat(&merged.StrikeoutWalkRatio, incoming.StrikeoutWalkRatio)
	mergeFloat(&merged.StrikePercentage, incoming.StrikePercentage)
	mergeOptional(&merged.ExpectedBattingAvg, incoming.ExpectedBattingAvg)
	mergeOptional(&merged.ExpectedSlugging, incoming.ExpectedSlugging)
	mergeOptional(&merged.ExpectedWOBA, incoming.ExpectedWOBA)
	mergeOptional(&merged.ExpectedERA, incoming.ExpectedERA)
	mergeOptional(&merged.BarrelPct, incoming.BarrelPct)
	mergeOptional(&merged.HardHitPct, incoming.HardHitPct)
	mergeOptional(&merged.AvgExitVelocity, incoming.AvgExitVelocity)
	mergeOptional(&merged.AvgLaunchAngle, incoming.AvgLaunchAngle)
	mergeOptional(&merged.SweetSpotPct, incoming.SweetSpotPct)

	return merged
}

// MergeSeasonGroup combines multiple team splits from the same stat group into
// one season-level row. When a player is traded mid-season, the latest real-team
// split becomes the canonical team reference for the merged season.
func MergeSeasonGroup(existing SeasonStatRecord, incoming SeasonStatRecord, group string) SeasonStatRecord {
	if existing.Year == 0 {
		return incoming
	}

	switch group {
	case "hitting":
		if !existing.HasHitting {
			return MergeSeasonStats(existing, incoming)
		}
		return mergeHittingSeason(existing, incoming)
	case "pitching":
		if !existing.HasPitching {
			return MergeSeasonStats(existing, incoming)
		}
		return mergePitchingSeason(existing, incoming)
	default:
		return MergeSeasonStats(existing, incoming)
	}
}

func mergeHittingSeason(existing SeasonStatRecord, incoming SeasonStatRecord) SeasonStatRecord {
	merged := existing
	mergeSeasonIdentity(&merged, incoming)
	merged.HasHitting = true
	merged.GamesPlayed += incoming.GamesPlayed
	merged.PlateAppearances += incoming.PlateAppearances
	merged.AtBats += incoming.AtBats
	merged.Hits += incoming.Hits
	merged.Doubles += incoming.Doubles
	merged.Triples += incoming.Triples
	merged.HomeRuns += incoming.HomeRuns
	merged.Runs += incoming.Runs
	merged.RBI += incoming.RBI
	merged.Walks += incoming.Walks
	merged.Strikeouts += incoming.Strikeouts
	merged.StolenBases += incoming.StolenBases
	recomputeHittingRates(&merged)

	return merged
}

func mergePitchingSeason(existing SeasonStatRecord, incoming SeasonStatRecord) SeasonStatRecord {
	merged := existing
	combinedEarnedRuns := earnedRuns(existing) + earnedRuns(incoming)
	combinedStrikePct := weightedAverage(existing.StrikePercentage, pitchingWeight(existing), incoming.StrikePercentage, pitchingWeight(incoming))
	mergeSeasonIdentity(&merged, incoming)
	merged.HasPitching = true
	merged.GamesPlayed += incoming.GamesPlayed
	merged.GamesStarted += incoming.GamesStarted
	merged.Wins += incoming.Wins
	merged.Losses += incoming.Losses
	merged.InningsPitchedOuts += incoming.InningsPitchedOuts
	merged.InningsPitched = baseballInningsFromOuts(merged.InningsPitchedOuts)
	merged.HitsAllowed += incoming.HitsAllowed
	merged.WalksAllowed += incoming.WalksAllowed
	merged.HomeRunsAllowed += incoming.HomeRunsAllowed
	merged.Strikeouts += incoming.Strikeouts
	recomputePitchingRates(&merged, combinedEarnedRuns, combinedStrikePct)

	return merged
}

func mergeSeasonIdentity(current *SeasonStatRecord, incoming SeasonStatRecord) {
	if incoming.sourceOrder >= current.sourceOrder && incoming.TeamID != 0 {
		current.TeamID = incoming.TeamID
		current.TeamName = incoming.TeamName
		current.sourceOrder = incoming.sourceOrder
	}
	if current.TeamName == "" {
		current.TeamName = incoming.TeamName
	}
	// MLB split ages can differ within one season if a player is traded after a
	// birthday. Keeping the higher age preserves the end-of-season season_stats
	// age instead of making it depend on split merge order.
	if incoming.Age > current.Age {
		current.Age = incoming.Age
	}
}

func recomputeHittingRates(record *SeasonStatRecord) {
	if record.AtBats > 0 {
		record.BattingAvg = roundRate(float64(record.Hits) / float64(record.AtBats))
		record.SLG = roundRate(float64(totalBases(*record)) / float64(record.AtBats))
	}

	obpDenominator := record.AtBats + record.Walks
	if obpDenominator > 0 {
		record.OBP = roundRate(float64(record.Hits+record.Walks) / float64(obpDenominator))
	}

	ballsInPlay := record.AtBats - record.Strikeouts - record.HomeRuns
	if ballsInPlay > 0 {
		record.BABIP = roundRate(float64(record.Hits-record.HomeRuns) / float64(ballsInPlay))
	}

	if record.OBP > 0 || record.SLG > 0 {
		record.OPS = roundRate(record.OBP + record.SLG)
	}
}

func recomputePitchingRates(record *SeasonStatRecord, combinedEarnedRuns float64, combinedStrikePct float64) {
	innings := inningsFromOuts(record.InningsPitchedOuts)
	if innings <= 0 {
		return
	}

	record.InningsPitched = baseballInningsFromOuts(record.InningsPitchedOuts)
	record.ERA = roundRate((combinedEarnedRuns / innings) * 9)
	record.WHIP = roundRate(float64(record.HitsAllowed+record.WalksAllowed) / innings)
	record.StrikeoutsPer9 = roundRate(float64(record.Strikeouts) * 9 / innings)
	record.WalksPer9 = roundRate(float64(record.WalksAllowed) * 9 / innings)
	record.HitsPer9 = roundRate(float64(record.HitsAllowed) * 9 / innings)
	record.HomeRunsPer9 = roundRate(float64(record.HomeRunsAllowed) * 9 / innings)
	if record.WalksAllowed > 0 {
		record.StrikeoutWalkRatio = roundRate(float64(record.Strikeouts) / float64(record.WalksAllowed))
	}
	record.StrikePercentage = roundRate(combinedStrikePct)
}

func pitchingWeight(record SeasonStatRecord) float64 {
	if record.InningsPitchedOuts > 0 {
		return float64(record.InningsPitchedOuts)
	}
	return float64(record.GamesPlayed)
}

func earnedRuns(record SeasonStatRecord) float64 {
	innings := inningsFromOuts(record.InningsPitchedOuts)
	if innings <= 0 || record.ERA <= 0 {
		return 0
	}
	return (record.ERA * innings) / 9
}

func weightedAverage(left float64, leftWeight float64, right float64, rightWeight float64) float64 {
	totalWeight := leftWeight + rightWeight
	if totalWeight == 0 {
		return 0
	}
	return ((left * leftWeight) + (right * rightWeight)) / totalWeight
}

func totalBases(record SeasonStatRecord) int {
	singles := record.Hits - record.Doubles - record.Triples - record.HomeRuns
	return singles + (record.Doubles * 2) + (record.Triples * 3) + (record.HomeRuns * 4)
}

func roundRate(value float64) float64 {
	return math.Round(value*1000) / 1000
}

// ApplySavantEnrichment attaches optional Statcast-derived fields after the MLB
// baseline row has already been built.
func ApplySavantEnrichment(record SeasonStatRecord, enrichment *SavantEnrichment) SeasonStatRecord {
	if enrichment == nil {
		return record
	}

	record.ExpectedBattingAvg = enrichment.ExpectedBattingAvg
	record.ExpectedSlugging = enrichment.ExpectedSlugging
	record.ExpectedWOBA = enrichment.ExpectedWOBA
	record.ExpectedERA = enrichment.ExpectedERA
	record.BarrelPct = enrichment.BarrelPct
	record.HardHitPct = enrichment.HardHitPct
	record.AvgExitVelocity = enrichment.AvgExitVelocity
	record.AvgLaunchAngle = enrichment.AvgLaunchAngle
	record.SweetSpotPct = enrichment.SweetSpotPct

	return record
}

// intFromMap handles the mixed float/string shapes returned by the MLB Stats API.
func intFromMap(values map[string]any, key string) int {
	raw, ok := values[key]
	if !ok || raw == nil {
		return 0
	}

	switch value := raw.(type) {
	case float64:
		return int(value)
	case string:
		return parseStringInt(value)
	default:
		return 0
	}
}

// floatFromMap handles the mixed float/string shapes returned by the MLB Stats API.
func floatFromMap(values map[string]any, key string) float64 {
	raw, ok := values[key]
	if !ok || raw == nil {
		return 0
	}

	switch value := raw.(type) {
	case float64:
		return value
	case string:
		return parseStringFloat(value)
	default:
		return 0
	}
}

func baseballInningsOutsFromMap(values map[string]any, key string) (int, error) {
	raw, ok := values[key]
	if !ok || raw == nil {
		return 0, nil
	}

	switch value := raw.(type) {
	case float64:
		return parseBaseballInnings(strconv.FormatFloat(value, 'f', -1, 64))
	case string:
		return parseBaseballInnings(value)
	default:
		return 0, nil
	}
}

func inningsFromOuts(outs int) float64 {
	if outs <= 0 {
		return 0
	}

	return float64(outs) / 3
}
