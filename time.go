package timeline

import "time"

// MaxTime is a time helper returning the greatest time of t1 and t2
func MaxTime(t1 time.Time, t2 time.Time) time.Time {
	if t2.After(t1) {
		return t2
	}
	return t1
}

// MinTime is a time helper returning the lowest time of t1 and t2
func MinTime(t1 time.Time, t2 time.Time) time.Time {
	if t2.Before(t1) {
		return t2
	}
	return t1
}
