package utils

import "time"

func RoundToNearestInterval(t time.Time, interval time.Duration) time.Time {
	return t.Round(interval)
}

func GetRangeBetween(from, to time.Time, interval time.Duration) []time.Time {
	from = RoundToNearestInterval(from, interval)
	to = RoundToNearestInterval(to, interval)
	var rangeBetween []time.Time
	for t := from; t.Before(to); t = t.Add(interval) {
		rangeBetween = append(rangeBetween, t)
	}
	return rangeBetween
}
