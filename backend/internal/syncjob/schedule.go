package syncjob

import "time"

const (
	inSeasonStartMonth = time.April
	inSeasonEndMonth   = time.October
)

type Schedule struct {
	Hour    int
	Minute  int
	Weekday time.Weekday
}

func (s Schedule) NextRun(now time.Time) time.Time {
	localNow := now
	if IsInSeason(localNow) {
		return nextDailyRun(localNow, s.Hour, s.Minute)
	}

	return nextWeeklyRun(localNow, s.Weekday, s.Hour, s.Minute)
}

func IsInSeason(now time.Time) bool {
	month := now.Month()
	return month >= inSeasonStartMonth && month <= inSeasonEndMonth
}

func nextDailyRun(now time.Time, hour int, minute int) time.Time {
	candidate := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
	if !candidate.After(now) {
		candidate = candidate.AddDate(0, 0, 1)
	}

	return candidate
}

func nextWeeklyRun(now time.Time, weekday time.Weekday, hour int, minute int) time.Time {
	candidate := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())

	daysUntil := (int(weekday) - int(now.Weekday()) + 7) % 7
	candidate = candidate.AddDate(0, 0, daysUntil)
	if !candidate.After(now) {
		candidate = candidate.AddDate(0, 0, 7)
	}

	return candidate
}
