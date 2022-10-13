// Copyright 2022 by lolorenzo77. All rights reserved.
// Use of this source code is governed by MIT licence that can be found in the LICENSE file.

package timeline

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// Defines the chronological direction of a timeslice:
//   - AntiChronological
//   - Undefined
//   - Chronological
type Direction int

const (
	AntiChronological Direction = -1
	Undefined         Direction = 0
	Chronological     Direction = 1
)

// 2 timeslices comparisons
//   - EQUAL if equal and in the same direction.
//   - DIFFERENT if not equal.
//   - OPPOSITE if equal but in the opposite direction.
type Compare uint8

const (
	DIFFERENT Compare = 0b00000000
	EQUAL     Compare = 0b00000001
	OPPOSITE  Compare = 0b00000010
)

// Binary flag defining the position of a time compared with a timeslice:
//   - TS_UNDEF
//   - TS_OUTSIDE
//   - TS_BEFORE
//   - TS_START
//   - TS_INSIDE
//   - TS_END
//   - TS_AFTER
type TimePosition int

const (
	TS_UNDEF  TimePosition = 0b00000000
	TS_OUT    TimePosition = 0b00010001
	TS_BEFORE TimePosition = 0b00010000
	TS_START  TimePosition = 0b00001000
	TS_WITHIN TimePosition = 0b00000100
	TS_END    TimePosition = 0b00000010
	TS_IN     TimePosition = 0b00001110
	TS_AFTER  TimePosition = 0b00000001
)

// TimeSlice represents a range of times bounded by two dates (time.Time) From and To. Each boundary can be an infinite time.
type TimeSlice struct {
	From time.Time
	To   time.Time
}

// MakeTimeSlice creates and returns a new timeslice with a defined d duration and a starting time.
//   - If d == zero then the timeslice represents a single time.
//   - If d > 0 then the given times represents the begining
//   - If d < 0 then the given times represents the end
//
// panic if the given date is not defined (zero time)
func MakeTimeSlice(dte time.Time, d time.Duration) TimeSlice {
	if dte.IsZero() {
		panic(dte)
	}
	ts := &TimeSlice{From: dte, To: dte.Add(d)}
	return *ts
}

// String returns default formating: "{ from - to : duration }".
//
// An infinite begining prints "past" and an infinite end prints "future".
//   - if a boundary does not have any hours nor minutes nor seconds, then prints only the date.
//   - if a boundary does not have any year nor month nor day, then prints only the time.
func (ts TimeSlice) String() string {
	var strfrom, strto, strdur string
	if ts.From.IsZero() {
		strfrom = "past"
	} else {
		if ts.From.Hour() == 0 && ts.From.Minute() == 0 && ts.From.Second() == 0 {
			strfrom = ts.From.Format("20060102 MST")
		} else {
			strfrom = ts.From.Format("20060102 15:04:05 MST")
		}
	}
	if ts.From.IsZero() {
		strto = "future"
	} else {
		if ts.To.Hour() == 0 && ts.To.Minute() == 0 && ts.To.Second() == 0 {
			strto = ts.To.Format("20060102 MST")
		} else if ts.From.Year() == ts.To.Year() && ts.From.Month() == ts.To.Month() && ts.From.Day() == ts.To.Day() {
			strto = ts.To.Format("15:04:05")
		} else {
			strto = ts.To.Format("20060102 15:04:05 MST")
		}
	}
	strdur = ts.Duration().FormatOrderOfMagnitude(3)
	return fmt.Sprintf("{ %s - %s : %s }", strfrom, strto, strdur)
}

// Moves the begining of the timeslice to the requested time.
//
// Postpone the end of the timeslice if the request time exceeds it, or
// cap the request time to the end of the timeslice, according to the direction.
//
// In case of capped or postponed move the timeslice become a single date.
func (pts *TimeSlice) FromMove(request time.Time, cap bool) {
	if !pts.To.IsZero() {
		if cap {
			if (pts.Direction() == Chronological && request.After(pts.To)) || (pts.Direction() == AntiChronological && request.Before(pts.To)) {
				request = pts.To
			}
		} else {
			if (pts.Direction() == Chronological && request.After(pts.To)) || (pts.Direction() == AntiChronological && request.Before(pts.To)) {
				pts.To = request
			}
		}
	}
	pts.From = request
}

// FromExtend add the duration at the begining of the timeslice.
//   - if the duration is negative then the begining time moves backward.
//   - if *pts.From is infinite, then nothing occurs.
//
// The timeslice direction can change
func (pts *TimeSlice) FromExtend(dur Duration) {
	if !pts.From.IsZero() {
		pts.To = pts.To.Add(dur.Duration)
	}
}

// Moves the end of the timeslice to the requested time.
//
// Move back the begining of the timeslice if the request time exceeds it, or
// cap the request time to the begining of the timeslice, according to the direction.
//
// In case of capped or back move the timeslice become a single date.
func (pts *TimeSlice) ToMove(request time.Time, cap bool) {
	if !pts.From.IsZero() {
		if cap {
			if (pts.Direction() == Chronological && request.Before(pts.From)) || (pts.Direction() == AntiChronological && request.After(pts.From)) {
				request = pts.From
			}
		} else {
			if (pts.Direction() == Chronological && request.Before(pts.From)) || (pts.Direction() == AntiChronological && request.After(pts.From)) {
				pts.From = request
			}
		}
	}
	pts.To = request
}

// ToExtend add the duration at the end of the timeslice.
//   - if the duration is negative then the end time moves backward.
//   - if *pts.To is infinite, then nothing occurs.
//
// The timeslice direction can change.
func (pts *TimeSlice) ToExtend(dur Duration) {
	if !pts.To.IsZero() {
		pts.To = pts.To.Add(dur.Duration)
	}
}

// Shift moves simultaneously both boundaries of the timeslice.
// Move occurs only for finite boundaries.
// Move to the past if dur is negative.
func (pts *TimeSlice) Shift(dur Duration) {
	pts.FromExtend(dur)
	pts.ToExtend(dur)
}

// Middle returns the time at the middle of the timeslice.
//
// Return a zero time if the timeslice has infinite boundaries
func (ts TimeSlice) Middle() time.Time {
	if ts.IsInfinite() {
		return time.Time{}
	}
	d := ts.Duration()
	if d.Duration == 0 {
		return ts.From
	}
	return ts.From.Add(time.Duration(d.Duration / 2.0))
}

// Duration returns the timeslice duration.
//   - returns zero if timeslice boundaries have the exact same times.
func (ts TimeSlice) Duration() Duration {
	var d Duration
	if ts.From.IsZero() || ts.To.IsZero() {
		return d
	}
	d.Duration = ts.To.Sub(ts.From)
	d.IsFinite = true
	return d
}

// WhereIs returns the position of t within the timeslice
//   - returns zero if timeslice boundaries have the exact same times.
func (ts TimeSlice) WhereIs(t time.Time) TimePosition {
	if ts.From.IsZero() && ts.To.IsZero() {
		return TS_UNDEF
	}
	var w TimePosition
	if t.Equal(ts.From) {
		w = w | TS_START
	}
	if t.Before(ts.From) {
		w = w | TS_BEFORE
	}
	if t.Equal(ts.To) {
		w = w | TS_END
	}
	if t.After(ts.To) {
		w = w | TS_AFTER
	}
	if t.After(ts.From) && t.Before(ts.To) {
		w = w | TS_WITHIN
	}

	return w
}

// IsInfinite returns true if at least one boundary is a zero time
func (ts TimeSlice) IsInfinite() bool {
	if ts.From.IsZero() || ts.To.IsZero() {
		return true
	}
	return false
}

// IsZero returns true if both bouraies are a zero time
func (ts TimeSlice) IsZero() bool {
	return ts.From.IsZero() && ts.To.IsZero()
}

// Truncate returns the result of rounding t down to a multiple of dur (since the zero time).
// If dur <= 0, Truncate returns t stripped of any monotonic clock reading but otherwise unchanged.
func (ts TimeSlice) Truncate(dur time.Duration) TimeSlice {
	ts.From = ts.From.Truncate(dur)
	ts.To = ts.To.Truncate(dur)
	return ts
}

// Compare checks if 2 timeslices start and end at the same times, event if they're in a different timezone.
//   - returns EQUAL if equal and in the same direction.
//   - returns DIFFERENT if not equal.
//   - returns OPPOSITE if equal but in the opposite direction.
func (one TimeSlice) Compare(another TimeSlice) Compare {
	if one.From.Equal(another.From) && one.To.Equal(another.To) {
		return EQUAL
	} else if one.From.Equal(another.To) && one.To.Equal(another.From) {
		return OPPOSITE
	}
	return DIFFERENT
}

// Returns the direction of the timeslice.
//
// returns 'Undefined' if both boundaries are infinite or if the timeslice is a single date.
func (ts TimeSlice) Direction() Direction {
	if ts.From.IsZero() && ts.To.IsZero() {
		return Undefined
	}
	if ts.From.IsZero() {
		return AntiChronological
	}
	if ts.To.IsZero() {
		return Chronological
	}

	d := int(ts.To.Sub(ts.From))
	switch {
	case d < 0:
		return AntiChronological
	case d > 0:
		return Chronological
	default:
		return Undefined
	}
}

// Progress returns the progress rate of a given time within the timeslice, with the level of precision of the second.
//
// The progress is calculated from the begining of the timeslice, whatever its direction. The returned rate is always positive.
//
// returns 0.5 if the timeslice has no duration.
//
// for a chronological timeslice:
//   - returns 0 if datetime is before the begining
//   - returns 1 if datetime is after the end
//
// for an anti-chronological timeslice:
//   - returns 0 if datetime is after the begining
//   - returns 1 if datetime is before the end
func (ts TimeSlice) Progress(datetime time.Time) (rate float64) {
	if ts.IsInfinite() {
		return 0.5
	}

	dur := ts.Duration()
	rate = datetime.Sub(ts.From).Seconds() / dur.Seconds()
	if dur.Duration < 0 {
		rate = -rate
	}

	// bound it between 0 an 1
	if rate < 0 {
		rate = 0
	} else if rate > 1 {
		rate = 1
	}
	return rate
}

// WhatTime returns the datetime at a certain rate within the timeslice.
//
// The progress is calculated from the begining of the timeslice, whatever its direction. The returned date is always within the time slice
//
// returns a zero time if the timeslice has an infinite duration.
// If the timeslice is a single date then returns it.
func (ts TimeSlice) WhatTime(rate float64) time.Time {
	if ts.IsInfinite() {
		return time.Time{}
	}

	var t time.Time
	dur := ts.Duration()
	dprog := float64(dur.Duration) * rate
	t = ts.From.Add(time.Duration(dprog))

	// bount it within the timeslice boundaries
	if dur.Duration > 0 && t.After(ts.To) || dur.Duration < 0 && t.Before(ts.To) {
		t = ts.To
	}
	if dur.Duration > 0 && t.Before(ts.From) || dur.Duration < 0 && t.After(ts.From) {
		t = ts.From
	}
	return t
}

// Split a timeslice in multiple timeslices of a d duration.
//
// The end of a slice is the exact time of the begining of the next one.
// The last slice duration can be lower than d duration if thists duration is not a multiple of d.
//
// returns an error if a boundary is infinite.
//
// panic if d is <= 0
func (ts TimeSlice) Split(d time.Duration) ([]TimeSlice, error) {
	if d <= 0 {
		log.Fatalf("TimeSlice.Split with invalid duration: %v", d)
	}

	// check duration of ts
	if ts.IsInfinite() {
		return []TimeSlice{}, errors.New("unable to split an infinite timeslice")
	}
	dur := ts.Duration()
	if dur.Duration < 0 {
		d = -d
	}

	slices := make([]TimeSlice, 0)
	for {
		split := MakeTimeSlice(ts.From, d)
		if dur.Duration > 0 && split.To.After(ts.To) || dur.Duration < 0 && split.To.Before(ts.To) {
			split.To = ts.To
			if split.Duration().Duration != 0 {
				slices = append(slices, split)
			}
			break
		}
		slices = append(slices, split)
		ts.From = split.To
	}
	return slices, nil
}

// GetScanMask returns the best appropriate TimeMask for scanning a timeline and to ensure max Scans in a timeslice.
// The returned mask can be used directly by the scan function.
//   - returns MASK_NONE if the timeslice has infinite duration or maxScans = 0
//   - returns MASK_SHORTEST if the timselice is a single date
func (ts TimeSlice) GetScanMask(maxScans uint) (mask TimeMask) {
	if ts.IsInfinite() || maxScans == 0 {
		return MASK_NONE
	}

	d := ts.Duration()
	// returns MASK_SHORTEST if the timselice is a single date
	if d.Duration == 0 {
		return MASK_SHORTEST
	}
	// calculation on the absolute duration
	if d.Duration < 0 {
		d.Duration = -d.Duration
	}

	//log.Printf("m=%f h=%f d=%f M=%f Y=%f ", time.Duration(d).Minutes(), time.Duration(d).Hours(), d.Days(), d.Months(), d.Years())

	switch {
	case d.Minutes() <= float64(maxScans):
		mask = MASK_MINUTE
	case (d.Hours() * 4) <= float64(maxScans):
		mask = MASK_MINUTEx15
	case (d.Hours() * 2) <= float64(maxScans):
		mask = MASK_HALFHOUR
	case d.Hours() <= float64(maxScans):
		mask = MASK_HOUR
	case (d.Days() * 6) <= float64(maxScans):
		mask = MASK_HOURx4
	case (d.Days() * 2) <= float64(maxScans):
		mask = MASK_HALFDAY
	case d.Days() <= float64(maxScans):
		mask = MASK_DAY
	case d.Months() <= float64(maxScans):
		mask = MASK_MONTH
	case d.Quarters() <= float64(maxScans):
		mask = MASK_QUARTER
	default:
		mask = MASK_YEAR
	}
	return mask
}

// Scan returns next time, within the timeslice boundaries, matching mask.
//
// Scan always starts by the begining of the timeslice. If the begining is infinite then Scan returns a zero date and the cursor is reset to nil.
//
// Scan looks for the next time after the cursor matching the mask and returns it. The cursor moves to this returned time.
// If the matching time is over the timeslice boundary then Scan returns a zero time and reset the cursor.
//
// Use fBoundaries if you want the scanner to returns the boundary even if they do not match the mask.
//
// If the timeslice has an infinite end boundary, then the scan will never returns a nil cursor.
//
// panic if mask not an allowed value
func (ts TimeSlice) Scan(cursor *time.Time, mask TimeMask, fBoundaries bool) time.Time {
	if mask < MASK_min || mask > MASK_max {
		log.Fatalf("invalid scan mask: %d", mask)
	}
	if ts.From.IsZero() {
		return time.Time{}
	}

	newcursor := *cursor
	start := false

	// init the cursor if we start scanning
	if newcursor.IsZero() {
		start = true
		newcursor = ts.From
	}

	// calculate the next cursor according to the mask, and the direction
	var fmatch bool
	if ts.Direction() == AntiChronological {
		// apply the mask to the cursor
		newcursor, fmatch = mask.Apply(newcursor)

		// for the first scan, returns the begining of the timeslice
		// if it's matching with the mask, or if boundary requested
		if start && (fBoundaries || fmatch || mask.Add(newcursor).Equal(ts.From)) {
			newcursor = ts.From
			*cursor = newcursor
			return newcursor
		}

		// move the cursor one step backward, only if already matching the mask
		if fmatch {
			newcursor = mask.Sub(newcursor)
		}

		// check overflow
		if newcursor.Before(ts.To) {
			if fBoundaries && !cursor.Equal(ts.To) {
				newcursor = ts.To
			} else {
				newcursor = time.Time{}
			}
		}
	} else {
		// apply the mask to the cursor
		newcursor, fmatch = mask.Apply(newcursor)

		// for the first scan, returns the begining of the timeslice
		// if it's matching with the mask, or if boundary requested
		if start && (fBoundaries || fmatch) {
			newcursor = ts.From
			*cursor = newcursor
			return newcursor
		}

		// move the cursor one step
		newcursor = mask.Add(newcursor)

		// check end boundary
		if newcursor.After(ts.To) {
			if fBoundaries && !cursor.Equal(ts.To) {
				// returns the end of the timeslice
				newcursor = ts.To
			} else {
				// the scan if finished
				newcursor = time.Time{}
			}
		}
	}

	*cursor = newcursor
	return newcursor
}
