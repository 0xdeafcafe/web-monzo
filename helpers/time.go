package helpers

import "time"

// CompareDates checks if two time's are on the same day
func CompareDates(a time.Time, b time.Time) bool {
	aDay, aMonth, aYear := a.Date()
	bDay, bMonth, bYear := b.Date()

	return aDay == bDay && aMonth == bMonth && aYear == bYear
}
