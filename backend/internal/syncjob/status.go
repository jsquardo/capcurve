package syncjob

import (
	"sync"
	"time"
)

type StatusSnapshot struct {
	Enabled              bool       `json:"enabled"`
	Running              bool       `json:"running"`
	InSeason             bool       `json:"in_season"`
	TargetSeasonYear     int        `json:"target_season_year"`
	LastSyncStartedAt    *time.Time `json:"last_sync_started_at"`
	LastSyncCompletedAt  *time.Time `json:"last_sync_completed_at"`
	LastSuccessfulSyncAt *time.Time `json:"last_successful_sync_at"`
	NextScheduledSyncAt  *time.Time `json:"next_scheduled_sync_at"`
	LastError            string     `json:"last_error"`
}

type StatusStore struct {
	mu       sync.RWMutex
	snapshot StatusSnapshot
}

func NewStatusStore(enabled bool) *StatusStore {
	return &StatusStore{
		snapshot: StatusSnapshot{
			Enabled: enabled,
		},
	}
}

func (s *StatusStore) Snapshot() StatusSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.snapshot
}

func (s *StatusStore) SetNextRun(nextRun time.Time, inSeason bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	nextRunCopy := nextRun
	s.snapshot.NextScheduledSyncAt = &nextRunCopy
	s.snapshot.InSeason = inSeason
}

func (s *StatusStore) MarkRunStarted(startedAt time.Time, seasonYear int, inSeason bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	startedAtCopy := startedAt
	s.snapshot.Running = true
	s.snapshot.InSeason = inSeason
	s.snapshot.TargetSeasonYear = seasonYear
	s.snapshot.LastSyncStartedAt = &startedAtCopy
}

func (s *StatusStore) MarkRunCompleted(completedAt time.Time, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	completedAtCopy := completedAt
	s.snapshot.Running = false
	s.snapshot.LastSyncCompletedAt = &completedAtCopy
	if err != nil {
		s.snapshot.LastError = err.Error()
		return
	}

	s.snapshot.LastSuccessfulSyncAt = &completedAtCopy
	s.snapshot.LastError = ""
}
