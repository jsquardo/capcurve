package syncjob

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/jsquardo/capcurve/internal/models"
	"github.com/stretchr/testify/require"
)

type fakeActivePlayerSource struct {
	ids []int
	err error
}

func (f fakeActivePlayerSource) ActivePlayerMLBIDs(context.Context) ([]int, error) {
	if f.err != nil {
		return nil, f.err
	}

	return f.ids, nil
}

type fakePlayerSyncer struct {
	syncedBySeason map[int][]int
	fail           map[int]error
}

func (f *fakePlayerSyncer) RefreshPlayerSeason(_ context.Context, playerID int, seasonYear int) (*models.Player, error) {
	if f.syncedBySeason == nil {
		f.syncedBySeason = make(map[int][]int)
	}
	f.syncedBySeason[seasonYear] = append(f.syncedBySeason[seasonYear], playerID)
	if err, ok := f.fail[playerID]; ok {
		return nil, err
	}

	return &models.Player{MLBID: playerID}, nil
}

func TestRunOnceSyncsOnlyActivePlayersReturnedBySource(t *testing.T) {
	t.Parallel()

	syncer := &fakePlayerSyncer{}
	now := time.Date(2026, time.July, 10, 5, 0, 0, 0, time.UTC)
	service := &Service{
		logger:   slog.Default(),
		schedule: Schedule{Hour: 5, Weekday: time.Monday},
		location: time.UTC,
		now:      func() time.Time { return now },
		players: fakeActivePlayerSource{
			ids: []int{592450, 660271, 665742},
		},
		syncer: syncer,
	}

	err := service.RunOnce(context.Background())

	require.NoError(t, err)
	require.Equal(t, []int{592450, 660271, 665742}, syncer.syncedBySeason[2026])
}

func TestRunOnceReturnsFailedPlayerIDs(t *testing.T) {
	t.Parallel()

	syncer := &fakePlayerSyncer{
		fail: map[int]error{
			660271: errors.New("mlb unavailable"),
		},
	}
	now := time.Date(2026, time.July, 10, 5, 0, 0, 0, time.UTC)
	service := &Service{
		logger:   slog.Default(),
		schedule: Schedule{Hour: 5, Weekday: time.Monday},
		location: time.UTC,
		now:      func() time.Time { return now },
		players: fakeActivePlayerSource{
			ids: []int{592450, 660271},
		},
		syncer: syncer,
	}

	err := service.RunOnce(context.Background())

	require.EqualError(t, err, "sync completed with 1 player errors: [660271]")
	require.Equal(t, []int{592450, 660271}, syncer.syncedBySeason[2026])
}

func TestRunOnceStopsWhenContextIsCanceled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	syncer := &fakePlayerSyncer{}
	service := &Service{
		logger:   slog.Default(),
		schedule: Schedule{Hour: 5, Weekday: time.Monday},
		location: time.UTC,
		now:      func() time.Time { return time.Date(2026, time.July, 10, 5, 0, 0, 0, time.UTC) },
		players: fakeActivePlayerSource{
			ids: []int{592450, 660271},
		},
		syncer: syncer,
	}

	err := service.RunOnce(ctx)

	require.ErrorIs(t, err, context.Canceled)
	require.Empty(t, syncer.syncedBySeason)
}

func TestRunOnceUsesMostRecentCompletedSeasonDuringOffseason(t *testing.T) {
	t.Parallel()

	syncer := &fakePlayerSyncer{}
	now := time.Date(2026, time.January, 12, 5, 0, 0, 0, time.UTC)
	service := &Service{
		logger:   slog.Default(),
		schedule: Schedule{Hour: 5, Weekday: time.Monday},
		location: time.UTC,
		now:      func() time.Time { return now },
		players: fakeActivePlayerSource{
			ids: []int{592450, 660271},
		},
		syncer: syncer,
	}

	err := service.RunOnce(context.Background())

	require.NoError(t, err)
	require.Equal(t, []int{592450, 660271}, syncer.syncedBySeason[2025])
}
