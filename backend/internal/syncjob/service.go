package syncjob

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"sync"
	"time"

	"github.com/jsquardo/capcurve/internal/models"
	"gorm.io/gorm"
)

type ActivePlayerSource interface {
	ActivePlayerMLBIDs(ctx context.Context) ([]int, error)
}

type PlayerSyncer interface {
	RefreshPlayerSeason(ctx context.Context, playerID int, seasonYear int) (*models.Player, error)
}

type Service struct {
	logger     *slog.Logger
	schedule   Schedule
	location   *time.Location
	now        func() time.Time
	status     *StatusStore
	statusMu   sync.Mutex
	players    ActivePlayerSource
	syncer     PlayerSyncer
	timerAfter func(time.Duration) <-chan time.Time
}

type Options struct {
	Logger   *slog.Logger
	Schedule Schedule
	Location *time.Location
	Now      func() time.Time
	Status   *StatusStore
}

func NewService(db *gorm.DB, syncer PlayerSyncer, opts Options) *Service {
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}

	location := opts.Location
	if location == nil {
		location = time.Local
	}

	now := opts.Now
	if now == nil {
		now = time.Now
	}

	status := opts.Status
	if status == nil {
		status = NewStatusStore(true)
	}

	return &Service{
		logger:     logger,
		schedule:   opts.Schedule,
		location:   location,
		now:        now,
		status:     status,
		players:    &gormPlayerStore{db: db},
		syncer:     syncer,
		timerAfter: time.After,
	}
}

func (s *Service) Start(ctx context.Context) {
	s.logger.Info("season-aware sync scheduler started", "timezone", s.location.String(), "daily_hour", s.schedule.Hour, "daily_minute", s.schedule.Minute, "offseason_weekday", s.schedule.Weekday.String())

	for {
		now := s.now().In(s.location)
		nextRun := s.schedule.NextRun(now)
		wait := time.Until(nextRun)
		if wait < 0 {
			wait = 0
		}

		s.statusStore().SetNextRun(nextRun, IsInSeason(now))

		s.logger.Info("next sync scheduled", "now", now.Format(time.RFC3339), "next_run", nextRun.Format(time.RFC3339), "in_season", IsInSeason(now))

		select {
		case <-ctx.Done():
			s.logger.Info("sync scheduler stopped")
			return
		case <-s.timerAfter(wait):
			if err := s.RunOnce(ctx); err != nil && !errors.Is(err, context.Canceled) {
				s.logger.Error("scheduled sync failed", "err", err)
			}
		}
	}
}

func (s *Service) RunOnce(ctx context.Context) error {
	now := s.now().In(s.location)
	seasonYear := TargetSeasonYear(now)
	s.statusStore().MarkRunStarted(now, seasonYear, IsInSeason(now))

	mlbIDs, err := s.players.ActivePlayerMLBIDs(ctx)
	if err != nil {
		s.statusStore().MarkRunCompleted(s.now().In(s.location), err)
		return fmt.Errorf("load active players: %w", err)
	}

	s.logger.Info("starting active-player sync", "player_count", len(mlbIDs), "season_year", seasonYear, "in_season", IsInSeason(now))

	var failed []int
	for _, playerID := range mlbIDs {
		if err := ctx.Err(); err != nil {
			return err
		}

		if _, err := s.syncer.RefreshPlayerSeason(ctx, playerID, seasonYear); err != nil {
			failed = append(failed, playerID)
			s.logger.Error("player sync failed", "mlb_id", playerID, "err", err)
			continue
		}

		s.logger.Info("player synced", "mlb_id", playerID)
	}

	if len(failed) > 0 {
		sort.Ints(failed)
		err := fmt.Errorf("sync completed with %d player errors: %v", len(failed), failed)
		s.statusStore().MarkRunCompleted(s.now().In(s.location), err)
		return err
	}

	s.logger.Info("active-player sync completed", "player_count", len(mlbIDs), "season_year", seasonYear)
	s.statusStore().MarkRunCompleted(s.now().In(s.location), nil)
	return nil
}

func (s *Service) Snapshot() StatusSnapshot {
	return s.statusStore().Snapshot()
}

func (s *Service) statusStore() *StatusStore {
	s.statusMu.Lock()
	defer s.statusMu.Unlock()

	if s.status == nil {
		s.status = NewStatusStore(true)
	}

	return s.status
}

type gormPlayerStore struct {
	db *gorm.DB
}

func (s *gormPlayerStore) ActivePlayerMLBIDs(ctx context.Context) ([]int, error) {
	var ids []int
	err := s.db.WithContext(ctx).
		Model(&models.Player{}).
		Where("active = ?", true).
		Order("mlb_id ASC").
		Pluck("mlb_id", &ids).
		Error
	if err != nil {
		return nil, err
	}

	return ids, nil
}
