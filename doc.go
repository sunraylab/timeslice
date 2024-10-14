// Copyright 2022-2024 by larry868. All rights reserved.
// Use of this source code is governed by MIT licence that can be found in the LICENSE file.

/*
timeline package provides 3 main types:
  - Duration: an extension of the default time.Duration type.
  - TimeSlice: represents a range of times bounded by two dates (time.Time) From and To. It accepts infinite boundaries (zero times) and can be chronological or anti-chronological.
  - TimeMask: used for scanning a TimeSlice and to get the time corresponding to a rounding o'clock period.

# Duration

Days, Month and Year methods provides usefull calculations based on the following typical time values:

	1 day = 24 hours
	1 year = 365.25 days in average because of leap years = 8766 hours
	1 month = 1 year / 12 = 30.4375 days = 730.5 hours
	1 quarter = 1 year / 4 = 121.75 days = 2922 hours

Duration provides also a special formating function to produce an output to give an order of magnitude of the duration,
at a choosen order, in a human-reading way like:

	'[-][100Y][12M][31d][24h][60m][99s][~]' with:
	- only non-zero components are output
	- starting with the biggest time components available, only the following choosen number of components are output
	- all trailing non zero components produce the `~` output meaning the duration in not exact
*/
package timeline
