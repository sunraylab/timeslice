// Copyright 2022 by lolorenzo77. All rights reserved.
// Use of this source code is governed by MIT licence that can be found in the LICENSE file.

package timeline

import (
	"fmt"
	"math"
	"time"
)

const (
	Day     time.Duration = 24 * time.Hour          // 1 day = 24 hours
	Week                  = 7 * Day                 // 1 week = 168 hours
	Year                  = 365.25 * 24 * time.Hour // 1 year = 365.25 days in average because of leap years = 8766 hours
	Month                 = Year / 12               // 1 month = 30.4375 days = 730.5 hours
	Quarter               = Year / 4                // 1 quatrer = 121.75 days = 2922 hours
)

// Duration type extends default time.Duration with additional Days, Month and Year methods, and with a special formating function.
// Handles Infinite duration.
type Duration struct {
	time.Duration
	IsFinite bool // by default the Duration is Infinite
}

// NewDuration factory to build a new Duration with an initial timeduration
func NewDuration(timeduration int64) Duration {
	d := &Duration{Duration: time.Duration(timeduration)}
	d.IsFinite = true
	return *d
}

// Adjust the duration accordint to the factor
func (d Duration) Adjust(factor float64) Duration {
	newdur := &Duration{}
	newdur.Duration = time.Duration(float64(d.Duration) * factor)
	return *newdur
}

// Days returns the number of days,
// assuming one day is 24h
func (d Duration) Days() float64 {
	if !d.IsFinite {
		return math.Inf(0)
	}
	return float64(d.Duration) / float64(Day)
}

// Days returns the number of weeks,
// assuming one week is 7 days
func (d Duration) Weeks() float64 {
	if !d.IsFinite {
		return math.Inf(0)
	}
	return float64(d.Duration) / float64(Week)
}

// Months returns the average number of months,
// assuming an average Month is a Year by 12
func (d Duration) Months() float64 {
	if !d.IsFinite {
		return math.Inf(0)
	}
	return float64(d.Duration) / float64(Month)
}

// Months returns the average number of months,
// assuming an average Month is a Year by 12
func (d Duration) Quarters() float64 {
	if !d.IsFinite {
		return math.Inf(0)
	}
	return float64(d.Duration) / float64(Quarter)
}

// Years returns the average number of years,
// assuming an average year is 365.25 days, so 730.5 hours, because of leap years
func (d Duration) Years() float64 {
	if !d.IsFinite {
		return math.Inf(0)
	}
	return float64(d.Duration) / float64(Year)
}

// FormatOrderOfMagnitude formats duration to give an order of magnitude, or order maxorder, in a human-reading way based on the following rules:
//  1. only non zero components are output,
//  2. starting with the biggest time components available, only the first maxorder components are included,
//  3. all trailing non zero components produce the `~` output meaning the duration in not exact
//  4. the result is based on average leadtime for a day, a month, a quarter and year.
//  5. values are truncated
//
// The output will be in the following format:
//
//	 '[-][100Y][12M][31d][24h][60m][99s][~]' with:
//		- only non-zero components are output
//		- starting with the biggest time components available, only the following choosen number of components are output
//		- all trailing non zero components produce the `~` output meaning the duration in not exact
//
// maxorder is bounded between 1 and 6.
//
// Special cases:
//
//	*pd == nil // returns "infinite"
//	*pd == 0 // returns "0"
//	*pd <= 0 // returns a string started with a minus symbol
func (leftd Duration) FormatOrderOfMagnitude(maxorder uint) (str string) {
	if !leftd.IsFinite {
		return "infinite"
	}
	if leftd.Duration == 0 {
		return "0"
	}
	if leftd.Duration < 0 {
		str = "-"
		leftd.Duration = -leftd.Duration
	}
	if leftd.Seconds() < 1 {
		return str + "0s~"
	}

	// bound maxorder to be able to produce at least one component
	if maxorder < 1 {
		maxorder = 1
	}

	dust := ""
	const components = "YMdhms"
	var cint int
	order := uint(0)
	for i := 0; i < 6; i++ {
		switch i {
		case 0: // years
			cint = int(leftd.Years())
			leftd.Duration = leftd.Duration - time.Duration(cint)*Year
		case 1: // months
			cint = int(leftd.Months())
			leftd.Duration = leftd.Duration - time.Duration(cint)*Month
		case 2: // days
			cint = int(leftd.Days())
			leftd.Duration = leftd.Duration - time.Duration(cint)*Day
		case 3: // hours
			cint = int(leftd.Duration.Hours())
			leftd.Duration = leftd.Duration - time.Duration(cint)*time.Hour
		case 4: // minutes
			cint = int(leftd.Duration.Minutes())
			leftd.Duration = leftd.Duration - time.Duration(cint)*time.Minute
		case 5: // seconds
			cint = int(leftd.Duration.Seconds())
			leftd.Duration = leftd.Duration - time.Duration(cint)*time.Second
		}
		// increment order
		if order > 0 {
			order++
		}
		// add the component to the output string only if not zero
		if cint != 0 {
			// output the component only if in the choosen order
			if order <= maxorder {
				str += fmt.Sprintf("%d%s", cint, string(components[i]))
				// start counting the order when the first non-zero component is found
				if order == 0 {
					order = 1
				}
			} else {
				dust = "~"
			}

		}
	}
	return str + dust
}

// Default formating
func (pd *Duration) String() string {
	if pd == nil {
		return "nil"
	}
	return pd.FormatOrderOfMagnitude(3)
}
