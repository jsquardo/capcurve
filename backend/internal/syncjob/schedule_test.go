package syncjob

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIsInSeason(t *testing.T) {
	t.Parallel()

	location := time.FixedZone("EST", -5*60*60)

	require.False(t, IsInSeason(time.Date(2026, time.March, 31, 12, 0, 0, 0, location)))
	require.True(t, IsInSeason(time.Date(2026, time.April, 1, 0, 0, 0, 0, location)))
	require.True(t, IsInSeason(time.Date(2026, time.October, 31, 23, 59, 0, 0, location)))
	require.False(t, IsInSeason(time.Date(2026, time.November, 1, 0, 0, 0, 0, location)))
}

func TestScheduleNextRunUsesDailyCadenceDuringSeason(t *testing.T) {
	t.Parallel()

	location := time.FixedZone("EST", -5*60*60)
	schedule := Schedule{
		Hour:    5,
		Minute:  0,
		Weekday: time.Monday,
	}

	now := time.Date(2026, time.July, 10, 3, 30, 0, 0, location)
	nextRun := schedule.NextRun(now)

	require.Equal(t, time.Date(2026, time.July, 10, 5, 0, 0, 0, location), nextRun)
}

func TestScheduleNextRunAdvancesToNextDayAfterDailyWindowPasses(t *testing.T) {
	t.Parallel()

	location := time.FixedZone("EST", -5*60*60)
	schedule := Schedule{
		Hour:    5,
		Minute:  0,
		Weekday: time.Monday,
	}

	now := time.Date(2026, time.August, 10, 7, 0, 0, 0, location)
	nextRun := schedule.NextRun(now)

	require.Equal(t, time.Date(2026, time.August, 11, 5, 0, 0, 0, location), nextRun)
}

func TestScheduleNextRunUsesWeeklyOffseasonCadence(t *testing.T) {
	t.Parallel()

	location := time.FixedZone("EST", -5*60*60)
	schedule := Schedule{
		Hour:    5,
		Minute:  0,
		Weekday: time.Monday,
	}

	now := time.Date(2026, time.January, 7, 9, 0, 0, 0, location)
	nextRun := schedule.NextRun(now)

	require.Equal(t, time.Date(2026, time.January, 12, 5, 0, 0, 0, location), nextRun)
}

func TestScheduleNextRunAdvancesToFollowingWeekAfterWeeklyWindowPasses(t *testing.T) {
	t.Parallel()

	location := time.FixedZone("EST", -5*60*60)
	schedule := Schedule{
		Hour:    5,
		Minute:  0,
		Weekday: time.Monday,
	}

	now := time.Date(2026, time.January, 5, 8, 0, 0, 0, location)
	nextRun := schedule.NextRun(now)

	require.Equal(t, time.Date(2026, time.January, 12, 5, 0, 0, 0, location), nextRun)
}
