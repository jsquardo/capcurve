package ingestion

import (
	"context"
	"fmt"
	"sort"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/jsquardo/capcurve/internal/models"
	"github.com/jsquardo/capcurve/internal/scoring"
)

type Service struct {
	db           *gorm.DB
	mlbClient    *MLBClient
	savantClient *SavantClient
	scorer       *scoring.Service
}

// NewService wires the MLB fetcher, Savant fetcher, and database persistence into
// one ingestion entry point.
func NewService(db *gorm.DB, mlbClient *MLBClient, savantClient *SavantClient) *Service {
	if mlbClient == nil {
		mlbClient = NewMLBClient(nil)
	}
	if savantClient == nil {
		savantClient = NewSavantClient(nil)
	}

	return &Service{
		db:           db,
		mlbClient:    mlbClient,
		savantClient: savantClient,
		scorer:       scoring.NewService(),
	}
}

// SyncPlayer fetches the player bio, season splits, and optional Savant data, then
// upserts the merged result into players and season_stats.
func (s *Service) SyncPlayer(ctx context.Context, playerID int) (*models.Player, error) {
	playerData, err := s.mlbClient.FetchPlayer(ctx, playerID)
	if err != nil {
		return nil, err
	}

	playerRecord, err := NormalizePlayer(playerData)
	if err != nil {
		return nil, err
	}

	hittingSplits, err := s.mlbClient.FetchYearByYearStats(ctx, playerID, "hitting")
	if err != nil {
		return nil, err
	}

	pitchingSplits, err := s.mlbClient.FetchYearByYearStats(ctx, playerID, "pitching")
	if err != nil {
		return nil, err
	}

	seasonRecords := make(map[string]SeasonStatRecord)
	if err := s.mergeSplits(seasonRecords, hittingSplits, "hitting"); err != nil {
		return nil, err
	}
	if err := s.mergeSplits(seasonRecords, pitchingSplits, "pitching"); err != nil {
		return nil, err
	}

	if err := s.applySavant(ctx, playerID, seasonRecords); err != nil {
		return nil, err
	}

	var player models.Player
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		player = models.Player{
			MLBID:       playerRecord.MLBID,
			FirstName:   playerRecord.FirstName,
			LastName:    playerRecord.LastName,
			Position:    playerRecord.Position,
			Bats:        playerRecord.Bats,
			Throws:      playerRecord.Throws,
			DateOfBirth: playerRecord.DateOfBirth,
			Active:      playerRecord.Active,
			ImageURL:    playerRecord.ImageURL,
		}

		if err := tx.Where("mlb_id = ?", playerRecord.MLBID).Assign(player).FirstOrCreate(&player).Error; err != nil {
			return err
		}

		if err := tx.Where("player_id = ? AND team_id = 0", player.ID).Delete(&models.SeasonStat{}).Error; err != nil {
			return err
		}

		for _, record := range orderedSeasonRecords(seasonRecords) {
			if err := tx.Where("player_id = ? AND year = ? AND team_id <> ?", player.ID, record.Year, record.TeamID).Delete(&models.SeasonStat{}).Error; err != nil {
				return err
			}

			season := modelFromSeasonRecord(player.ID, record)
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{Name: "player_id"},
					{Name: "year"},
					{Name: "team_id"},
				},
				DoUpdates: clause.AssignmentColumns([]string{
					"team_name",
					"age",
					"games_played",
					"games_started",
					"plate_appearances",
					"at_bats",
					"hits",
					"doubles",
					"triples",
					"home_runs",
					"runs",
					"rbi",
					"walks",
					"strikeouts",
					"stolen_bases",
					"batting_avg",
					"obp",
					"slg",
					"ops",
					"babip",
					"wins",
					"losses",
					"era",
					"whip",
					"innings_pitched",
					"hits_allowed",
					"walks_allowed",
					"home_runs_allowed",
					"strikeouts_per_9",
					"walks_per_9",
					"hits_per_9",
					"home_runs_per_9",
					"strikeout_walk_ratio",
					"strike_percentage",
					"expected_batting_avg",
					"expected_slugging",
					"expected_woba",
					"expected_era",
					"barrel_pct",
					"hard_hit_pct",
					"avg_exit_velocity",
					"avg_launch_angle",
					"sweet_spot_pct",
				}),
			}).Create(&season).Error; err != nil {
				return err
			}
		}

		if err := s.scorer.RecalculateYears(ctx, tx, seasonYears(seasonRecords)); err != nil {
			return err
		}

		return tx.Preload("SeasonStats").First(&player, player.ID).Error
	})
	if err != nil {
		return nil, err
	}

	return &player, nil
}

// mergeSplits folds one MLB stat group into the in-memory season map before persistence.
func (s *Service) mergeSplits(records map[string]SeasonStatRecord, splits []MLBSeasonSplit, group string) error {
	for _, split := range splits {
		if IsAggregateSeasonSplit(split) {
			continue
		}

		record, err := NormalizeSeasonSplit(split, group)
		if err != nil {
			return err
		}

		key := seasonKey(record.Year)
		records[key] = MergeSeasonGroup(records[key], record, group)
	}

	return nil
}

// applySavant enriches only the seasons that have hitting or pitching activity.
func (s *Service) applySavant(ctx context.Context, playerID int, records map[string]SeasonStatRecord) error {
	for key, record := range records {
		if record.PlateAppearances > 0 {
			enrichment, err := s.savantClient.FetchSeasonEnrichment(ctx, record.Year, playerID, SavantTypeBatter)
			if err != nil {
				return err
			}
			if enrichment != nil {
				records[key] = ApplySavantEnrichment(record, enrichment)
				record = records[key]
			}
		}

		if record.GamesStarted > 0 || record.InningsPitched > 0 {
			enrichment, err := s.savantClient.FetchSeasonEnrichment(ctx, record.Year, playerID, SavantTypePitcher)
			if err != nil {
				return err
			}
			if enrichment != nil {
				records[key] = ApplySavantEnrichment(record, enrichment)
			}
		}
	}

	return nil
}

// orderedSeasonRecords keeps inserts deterministic and easier to read in logs.
func orderedSeasonRecords(records map[string]SeasonStatRecord) []SeasonStatRecord {
	keys := make([]string, 0, len(records))
	for key := range records {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	ordered := make([]SeasonStatRecord, 0, len(keys))
	for _, key := range keys {
		ordered = append(ordered, records[key])
	}

	return ordered
}

func seasonYears(records map[string]SeasonStatRecord) []int {
	years := make([]int, 0, len(records))
	for _, record := range records {
		years = append(years, record.Year)
	}
	return years
}

// seasonKey matches the uniqueness rule used by the season_stats table.
func seasonKey(year int) string {
	return fmt.Sprintf("%d", year)
}

// modelFromSeasonRecord converts the normalized ingestion shape into the GORM model.
func modelFromSeasonRecord(playerID uint, record SeasonStatRecord) models.SeasonStat {
	return models.SeasonStat{
		PlayerID:           int(playerID),
		Year:               record.Year,
		TeamID:             record.TeamID,
		TeamName:           record.TeamName,
		Age:                record.Age,
		GamesPlayed:        record.GamesPlayed,
		GamesStarted:       record.GamesStarted,
		PlateAppearances:   record.PlateAppearances,
		AtBats:             record.AtBats,
		Hits:               record.Hits,
		Doubles:            record.Doubles,
		Triples:            record.Triples,
		HomeRuns:           record.HomeRuns,
		Runs:               record.Runs,
		RBI:                record.RBI,
		Walks:              record.Walks,
		Strikeouts:         record.Strikeouts,
		StolenBases:        record.StolenBases,
		BattingAvg:         record.BattingAvg,
		OBP:                record.OBP,
		SLG:                record.SLG,
		OPS:                record.OPS,
		BABIP:              record.BABIP,
		Wins:               record.Wins,
		Losses:             record.Losses,
		ERA:                record.ERA,
		WHIP:               record.WHIP,
		InningsPitched:     record.InningsPitched,
		HitsAllowed:        record.HitsAllowed,
		WalksAllowed:       record.WalksAllowed,
		HomeRunsAllowed:    record.HomeRunsAllowed,
		StrikeoutsPer9:     record.StrikeoutsPer9,
		WalksPer9:          record.WalksPer9,
		HitsPer9:           record.HitsPer9,
		HomeRunsPer9:       record.HomeRunsPer9,
		StrikeoutWalkRatio: record.StrikeoutWalkRatio,
		StrikePercentage:   record.StrikePercentage,
		ExpectedBattingAvg: record.ExpectedBattingAvg,
		ExpectedSlugging:   record.ExpectedSlugging,
		ExpectedWOBA:       record.ExpectedWOBA,
		ExpectedERA:        record.ExpectedERA,
		BarrelPct:          record.BarrelPct,
		HardHitPct:         record.HardHitPct,
		AvgExitVelocity:    record.AvgExitVelocity,
		AvgLaunchAngle:     record.AvgLaunchAngle,
		SweetSpotPct:       record.SweetSpotPct,
	}
}
