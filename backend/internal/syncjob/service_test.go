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
	synced []int
	fail   map[int]error
}

func (f *fakePlayerSyncer) SyncPlayer(_ context.Context, playerID int) (*models.Player, error) {
	f.synced = append(f.synced, playerID)
	if err, ok := f.fail[playerID]; ok {
		return nil, err
	}

	return &models.Player{MLBID: playerID}, nil
}

func TestRunOnceSyncsOnlyActivePlayersReturnedBySource(t *testing.T) {
	t.Parallel()

	syncer := &fakePlayerSyncer{}
	service := &Service{
		logger:   slog.Default(),
		schedule: Schedule{Hour: 5, Weekday: time.Monday},
		location: time.UTC,
		now:      time.Now,
		players: fakeActivePlayerSource{
			ids: []int{592450, 660271, 665742},
		},
		syncer: syncer,
	}

	err := service.RunOnce(context.Background())

	require.NoError(t, err)
	require.Equal(t, []int{592450, 660271, 665742}, syncer.synced)
}

func TestRunOnceReturnsFailedPlayerIDs(t *testing.T) {
	t.Parallel()

	syncer := &fakePlayerSyncer{
		fail: map[int]error{
			660271: errors.New("mlb unavailable"),
		},
	}
	service := &Service{
		logger:   slog.Default(),
		schedule: Schedule{Hour: 5, Weekday: time.Monday},
		location: time.UTC,
		now:      time.Now,
		players: fakeActivePlayerSource{
			ids: []int{592450, 660271},
		},
		syncer: syncer,
	}

	err := service.RunOnce(context.Background())

	require.EqualError(t, err, "sync completed with 1 player errors: [660271]")
	require.Equal(t, []int{592450, 660271}, syncer.synced)
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
		now:      time.Now,
		players: fakeActivePlayerSource{
			ids: []int{592450, 660271},
		},
		syncer: syncer,
	}

	err := service.RunOnce(ctx)

	require.ErrorIs(t, err, context.Canceled)
	require.Empty(t, syncer.synced)
}
