package timeslice

import "time"

// A Time Mask used for the scanning a TimeSlice or to get the time corresponding to the begining of a period.
type TimeMask int

// Available Time Masks used for the scanning a TimeSlice.
const (
	MASK_NONE      TimeMask = 0
	MASK_min       TimeMask = 1
	MASK_SHORTEST  TimeMask = 1
	MASK_MINUTE    TimeMask = 1
	MASK_MINUTEx15 TimeMask = 2
	MASK_HALFHOUR  TimeMask = 3
	MASK_HOUR      TimeMask = 4
	MASK_HOURx4    TimeMask = 5
	MASK_HALFDAY   TimeMask = 6
	MASK_DAY       TimeMask = 7
	MASK_MONTH     TimeMask = 8
	MASK_QUARTER   TimeMask = 9
	MASK_YEAR      TimeMask = 10
	MASK_max       TimeMask = 10
)

func (mask TimeMask) String() string {
	switch mask {
	case MASK_NONE:
		return "none"
	case MASK_MINUTE:
		return "minute"
	case MASK_MINUTEx15:
		return "15 minutes"
	case MASK_HALFHOUR:
		return "half-hour"
	case MASK_HOUR:
		return "hour"
	case MASK_HOURx4:
		return "4 hours"
	case MASK_HALFDAY:
		return "half-day"
	case MASK_DAY:
		return "day"
	case MASK_MONTH:
		return "month"
	case MASK_QUARTER:
		return "quarter"
	case MASK_YEAR:
		return "year"
	}
	return "?"
}

// Apply the mask to a date and returns a masked time and a flag indicating if the given time matches exactly the mask
//
// if mask is MASK_NONE then returned an unchanged time.
func (mask TimeMask) Apply(t time.Time) (masked time.Time, exactMatch bool) {
	// extract time components
	Y := t.Year()
	M := t.Month()
	d := t.Day()
	h := t.Hour()
	m := t.Minute()
	loc := t.Location()

	// apply
	switch mask {
	case MASK_NONE:
		masked = t
	case MASK_MINUTE:
		masked = time.Date(Y, M, d, h, m, 0, 0, loc)
	case MASK_MINUTEx15:
		masked = time.Date(Y, M, d, h, m/15*15, 0, 0, loc)
	case MASK_HALFHOUR:
		masked = time.Date(Y, M, d, h, m/30*30, 0, 0, loc)
	case MASK_HOUR:
		masked = time.Date(Y, M, d, h, 0, 0, 0, loc)
	case MASK_HOURx4:
		masked = time.Date(Y, M, d, h/4*4, 0, 0, 0, loc)
	case MASK_HALFDAY:
		masked = time.Date(Y, M, d, h/12*12, 0, 0, 0, loc)
	case MASK_DAY:
		masked = time.Date(Y, M, d, 0, 0, 0, 0, loc)
	case MASK_MONTH:
		masked = time.Date(Y, M, 1, 0, 0, 0, 0, loc)
	case MASK_QUARTER:
		masked = time.Date(Y, (M/3*3)+1, 1, 0, 0, 0, 0, loc)
	case MASK_YEAR:
		masked = time.Date(Y, 1, 1, 0, 0, 0, 0, loc)
	default:
		panic("unmanaged mask")
	}
	return masked, masked.Equal(t)
}

// Add applies the mask and adds the mask increment to the given time.
func (mask TimeMask) Add(t time.Time) time.Time {
	t, _ = mask.Apply(t)
	switch mask {
	case MASK_MINUTE:
		t = t.Add(time.Minute)
	case MASK_MINUTEx15:
		t = t.Add(time.Minute * 15)
	case MASK_HALFHOUR:
		t = t.Add(time.Minute * 30)
	case MASK_HOUR:
		t = t.Add(time.Hour)
	case MASK_HOURx4:
		t = t.Add(time.Hour * 4)
	case MASK_HALFDAY:
		t = t.Add(time.Hour * 12)
	case MASK_DAY:
		t = t.Add(time.Hour * 24)
	case MASK_MONTH:
		if t.Month() == 12 {
			t = time.Date(t.Year()+1, 1, 1, 0, 0, 0, 0, t.Location())
		} else {
			t = time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
		}
	case MASK_QUARTER:
		if t.Month() > (12 - 3) {
			t = time.Date(t.Year()+1, 1, 1, 0, 0, 0, 0, t.Location())
		} else {
			t = time.Date(t.Year(), t.Month()+(t.Month()/3*3), 1, 0, 0, 0, 0, t.Location())
		}
	case MASK_YEAR:
		t = time.Date(t.Year()+1, 1, 1, 0, 0, 0, 0, t.Location())
	}
	return t
}

// Sub applies the mask and substitute the mask increment to the given time.
func (mask TimeMask) Sub(t time.Time) time.Time {
	t, _ = mask.Apply(t)
	switch mask {
	case MASK_MINUTE:
		return t.Add(-time.Minute)
	case MASK_MINUTEx15:
		return t.Add(-time.Minute * 15)
	case MASK_HALFHOUR:
		return t.Add(-time.Minute * 30)
	case MASK_HOUR:
		return t.Add(-time.Hour)
	case MASK_HOURx4:
		return t.Add(-time.Hour * 4)
	case MASK_HALFDAY:
		return t.Add(-time.Hour * 12)
	case MASK_DAY:
		return t.Add(-time.Hour * 24)
	case MASK_MONTH:
		if t.Month() == 1 {
			t = time.Date(t.Year()-1, 12, 1, 0, 0, 0, 0, t.Location())
		} else {
			t = time.Date(t.Year(), t.Month()-1, 1, 0, 0, 0, 0, t.Location())
		}
		return t
	case MASK_QUARTER:
		if t.Month() <= 3 {
			t = time.Date(t.Year()-1, 12, 1, 0, 0, 0, 0, t.Location())
		} else {
			t = time.Date(t.Year(), t.Month()-(t.Month()/3*3), 1, 0, 0, 0, 0, t.Location())
		}
		return t
	case MASK_YEAR:
		t = time.Date(t.Year()-1, 1, 1, 0, 0, 0, 0, t.Location())
		return t
	}
	return t
}
