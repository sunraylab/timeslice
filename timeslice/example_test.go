package timeslice

import (
	"fmt"
	"time"

	"github.com/sunraylab/timeline/duration"
)

func ExampleTimeSlice_MoveFrom_one() {
	// take a date and build a timeslice staring at this date and ending 7 days after
	from := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	ts := MakeTimeslice(from, duration.Week)
	fmt.Println(ts)

	// Move forward the begining by 4 days
	ts.MoveFrom(ts.From.Add(duration.Day*4), false)
	fmt.Println(ts)

	// Move fotrward again the begining by 4 days
	ts.MoveFrom(ts.From.Add(duration.Day*4), false)
	fmt.Println(ts)

	// Output:
	// { 20220101 UTC - 20220108 UTC : 7d }
	// { 20220105 UTC - 20220108 UTC : 3d }
	// { 20220109 UTC - 20220109 UTC : 0 }
}

func ExampleTimeSlice_MoveTo_one() {
	// take a date and build a timeslice staring at this date and ending 7 days after
	from := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	ts := MakeTimeslice(from, duration.Week)
	fmt.Println(ts)

	// Move backward the ending by 4 days
	ts.MoveTo(ts.To.Add(-duration.Day*4), false)
	fmt.Println(ts)

	// Move backward again the ending by 4 days
	ts.MoveTo(ts.To.Add(-duration.Day*4), false)
	fmt.Println(ts)

	// Output:
	// { 20220101 UTC - 20220108 UTC : 7d }
	// { 20220101 UTC - 20220104 UTC : 3d }
	// { 20211231 UTC - 20211231 UTC : 0 }
}

func ExampleTimeSlice_MoveFrom_two() {
	// take a date and build a timeslice staring at this date and ending 7 days after
	from := time.Date(2022, 1, 6, 8, 0, 0, 0, time.UTC)
	ts := MakeTimeslice(from, duration.Day)
	fmt.Println(ts)
	ts.MoveFrom(time.Date(2022, 1, 6, 9, 0, 0, 0, time.UTC), true)
	fmt.Println(ts)
	ts.MoveFrom(time.Date(2022, 1, 7, 9, 0, 0, 0, time.UTC), true)
	fmt.Println(ts)
	// Output:
	// { 20220106 08:00:00 UTC - 20220107 08:00:00 UTC : 1d }
	// { 20220106 09:00:00 UTC - 20220107 08:00:00 UTC : 23h }
	// { 20220107 08:00:00 UTC - 08:00:00 : 0 }
}

func ExampleTimeSlice_MoveTo_two() {
	// take a date and build a timeslice staring at this date and ending 7 days after
	from := time.Date(2022, 1, 6, 8, 0, 0, 0, time.UTC)
	ts := MakeTimeslice(from, duration.Day)
	fmt.Println(ts)
	ts.MoveTo(time.Date(2022, 1, 6, 9, 0, 0, 0, time.UTC), true)
	fmt.Println(ts)
	ts.MoveTo(time.Date(2022, 1, 6, 7, 0, 0, 0, time.UTC), true)
	fmt.Println(ts)
	// Output:
	// { 20220106 08:00:00 UTC - 20220107 08:00:00 UTC : 1d }
	// { 20220106 08:00:00 UTC - 09:00:00 : 1h }
	// { 20220106 08:00:00 UTC - 08:00:00 : 0 }
}

func ExampleTimeSlice_String() {
	ts := MakeTimeslice(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC), duration.Week+time.Hour*31)
	fmt.Println(ts)
	// Output: { 20220101 UTC - 20220109 07:00:00 UTC : 8d7h }
}

func ExampleTimeSlice_Progress() {
	// take a 24 hours timeslice, starting the 2022,1,1 at midnight
	ts := MakeTimeslice(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC), duration.Day)

	// get the corresponding progress date the same day at 6 AM
	rate := ts.Progress(time.Date(2022, 1, 1, 6, 0, 0, 0, time.UTC))
	fmt.Println(rate)
	// Output: 0.25
}

func ExampleTimeSlice_WhatTime() {
	// take a 10 days timeslice, starting the 2022,1,1 at 8AM
	ts := MakeTimeslice(time.Date(2022, 1, 1, 8, 0, 0, 0, time.UTC), duration.Day*10)
	fmt.Println(ts)
	for rate := 0.0; rate <= 1.0; rate += 0.2 {
		t := ts.WhatTime(rate)
		fmt.Println(t.Format("20060102 15:04:05 MST"))
	}

	// Output:
	// { 20220101 08:00:00 UTC - 20220111 08:00:00 UTC : 10d }
	// 20220101 08:00:00 UTC
	// 20220103 08:00:00 UTC
	// 20220105 08:00:00 UTC
	// 20220107 08:00:00 UTC
	// 20220109 08:00:00 UTC
	// 20220111 08:00:00 UTC
}
