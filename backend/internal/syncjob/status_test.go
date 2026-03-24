package syncjob

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStatusStoreTracksRunLifecycle(t *testing.T) {
	t.Parallel()

	store := NewStatusStore(true)
	nextRun := time.Date(2026, time.July, 11, 5, 0, 0, 0, time.UTC)
	startedAt := time.Date(2026, time.July, 10, 5, 0, 0, 0, time.UTC)
	completedAt := startedAt.Add(2 * time.Minute)

	store.SetNextRun(nextRun, true)
	store.MarkRunStarted(startedAt, 2026, true)
	store.MarkRunCompleted(completedAt, nil)

	snapshot := store.Snapshot()
	require.True(t, snapshot.Enabled)
	require.False(t, snapshot.Running)
	require.True(t, snapshot.InSeason)
	require.Equal(t, 2026, snapshot.TargetSeasonYear)
	require.Equal(t, nextRun, *snapshot.NextScheduledSyncAt)
	require.Equal(t, startedAt, *snapshot.LastSyncStartedAt)
	require.Equal(t, completedAt, *snapshot.LastSyncCompletedAt)
	require.Equal(t, completedAt, *snapshot.LastSuccessfulSyncAt)
	require.Empty(t, snapshot.LastError)
}

func TestStatusStoreKeepsLastSuccessfulSyncWhenRunFails(t *testing.T) {
	t.Parallel()

	store := NewStatusStore(true)
	successAt := time.Date(2026, time.July, 10, 5, 2, 0, 0, time.UTC)
	failedAt := successAt.Add(24 * time.Hour)

	store.MarkRunCompleted(successAt, nil)
	store.MarkRunCompleted(failedAt, errors.New("sync completed with 1 player errors: [660271]"))

	snapshot := store.Snapshot()
	require.Equal(t, failedAt, *snapshot.LastSyncCompletedAt)
	require.Equal(t, successAt, *snapshot.LastSuccessfulSyncAt)
	require.Equal(t, "sync completed with 1 player errors: [660271]", snapshot.LastError)
}
