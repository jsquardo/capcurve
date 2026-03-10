package ingestion

import (
	"fmt"
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
		record.Age = intFromMap(split.Stat, "age")
		record.GamesPlayed = intFromMap(split.Stat, "gamesPlayed")
		record.GamesStarted = intFromMap(split.Stat, "gamesStarted")
		record.Wins = intFromMap(split.Stat, "wins")
		record.Losses = intFromMap(split.Stat, "losses")
		record.ERA = floatFromMap(split.Stat, "era")
		record.WHIP = floatFromMap(split.Stat, "whip")
		record.InningsPitched = floatFromMap(split.Stat, "inningsPitched")
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

// MergeSeasonStats overlays hitting, pitching, and enrichment data into a single
// row keyed by player, season, and team.
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
	mergeFloat(&merged.InningsPitched, incoming.InningsPitched)
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
