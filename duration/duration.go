// Copyright 2022 by lolorenzo77. All rights reserved.
// Use of this source code is governed by MIT licence that can be found in the LICENSE file.

/*
duration package provides an extension of the default time.Duration.

Days, Month and Year methods provides usefull calculations based on the following typical time values:

	1 day = 24 hours
	1 year = 365.25 days in average because of leap years = 8766 hours
	1 month = 1 year / 12 = 30.4375 days = 730.5 hours
	1 quarter = 1 year / 4 = 121.75 days = 2922 hours

duration provides also a special formating function to produce an output to give an order of magnitude of the duration,
at a choosen order, in a human-reading way like:

	'[-][100Y][12M][31d][24h][60m][99s][~]' with:
	- only non-zero components are output
	- starting with the biggest time components available, only the following choosen number of components are output
	- all trailing non zero components produce the `~` output meaning the duration in not exact
*/
package duration

import (
	"fmt"
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
type Duration time.Duration

// Adjust the duration accordint to the factor
func (d Duration) Adjust(factor float64) Duration {
	return Duration(float64(d) * factor)
}

// Days returns the number of days,
// assuming one day is 24h
func (d Duration) Days() float64 {
	return float64(d) / float64(Day)
}

// Days returns the number of weeks,
// assuming one week is 7 days
func (d Duration) Weeks() float64 {
	return float64(d) / float64(Week)
}

// Months returns the average number of months,
// assuming an average Month is a Year by 12
func (d Duration) Months() float64 {
	return float64(d) / float64(Month)
}

// Months returns the average number of months,
// assuming an average Month is a Year by 12
func (d Duration) Quarters() float64 {
	return float64(d) / float64(Quarter)
}

// Years returns the average number of years,
// assuming an average year is 365.25 days, so 730.5 hours, because of leap years
func (d Duration) Years() float64 {
	return float64(d) / float64(Year)
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
func (pd *Duration) FormatOrderOfMagnitude(maxorder uint) (str string) {
	if pd == nil {
		return "infinite"
	}
	leftd := Duration(*pd)
	if leftd == 0 {
		return "0"
	}
	if leftd < 0 {
		str = "-"
		leftd = -leftd
	}
	if time.Duration(leftd).Seconds() < 1 {
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
			leftd = leftd - Duration(cint)*Duration(Year)
		case 1: // months
			cint = int(leftd.Months())
			leftd = leftd - Duration(cint)*Duration(Month)
		case 2: // days
			cint = int(leftd.Days())
			leftd = leftd - Duration(cint)*Duration(Day)
		case 3: // hours
			cint = int(time.Duration(leftd).Hours())
			leftd = leftd - Duration(cint)*Duration(time.Hour)
		case 4: // minutes
			cint = int(time.Duration(leftd).Minutes())
			leftd = leftd - Duration(cint)*Duration(time.Minute)
		case 5: // seconds
			cint = int(time.Duration(leftd).Seconds())
			leftd = leftd - Duration(cint)*Duration(time.Second)
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
func (d Duration) String() string {
	return d.FormatOrderOfMagnitude(3)
}
