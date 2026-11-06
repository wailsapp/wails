package main

import (
	"sort"
	"time"
)

// WeeklyPoint is a cumulative star count at the beginning of an ISO week.
type WeeklyPoint struct {
	Date  time.Time
	Count int
}

// BuildWeeklyHistory groups current stargazer timestamps into Monday-based
// weeks and converts them into a cumulative history suitable for plotting.
func BuildWeeklyHistory(stars []Star) []WeeklyPoint {
	counts := make(map[time.Time]int)
	for _, star := range stars {
		if star.StarredAt.IsZero() {
			continue
		}
		date := star.StarredAt.UTC()
		daysSinceMonday := (int(date.Weekday()) + 6) % 7
		week := time.Date(date.Year(), date.Month(), date.Day()-daysSinceMonday, 0, 0, 0, 0, time.UTC)
		counts[week]++
	}

	weeks := make([]time.Time, 0, len(counts))
	for week := range counts {
		weeks = append(weeks, week)
	}
	sort.Slice(weeks, func(i, j int) bool { return weeks[i].Before(weeks[j]) })

	history := make([]WeeklyPoint, 0, len(weeks))
	total := 0
	for _, week := range weeks {
		total += counts[week]
		history = append(history, WeeklyPoint{Date: week, Count: total})
	}
	return history
}

// ReconcileCurrentCount adds today's authoritative count to the reconstructed
// history when the API list and repository total differ or the latest weekly
// point is stale. This keeps the endpoint's current number honest without
// inventing historical points.
func ReconcileCurrentCount(history []WeeklyPoint, currentCount int, now time.Time) []WeeklyPoint {
	if currentCount < 0 {
		return history
	}
	now = now.UTC()
	if len(history) == 0 {
		return []WeeklyPoint{{Date: now, Count: currentCount}}
	}
	last := &history[len(history)-1]
	if last.Date.Year() == now.Year() && last.Date.YearDay() == now.YearDay() {
		last.Count = currentCount
		return history
	}
	history = append(history, WeeklyPoint{Date: now, Count: currentCount})
	return history
}
