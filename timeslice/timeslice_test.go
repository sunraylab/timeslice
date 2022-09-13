package timeslice

import (
	"fmt"
	"testing"
	"time"

	"github.com/sunraylab/timeline/duration"
)

func TestSplit(t *testing.T) {

	// a timeslice staring 20220801 0h00:00 and 7 days long
	ts := MakeTimeslice(time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC), 7*duration.Day)

	// split in one
	tss0, err := ts.Split(7 * 24 * time.Hour)
	if err != nil || len(tss0) != 1 || tss0[0].Equal(ts) != 1 {
		t.Errorf("split in 1 error: %+v", tss0)
	}

	// daily split
	// want: retourner 7 tranches d'un jour
	tss1, err := ts.Split(24 * time.Hour)
	if err != nil || len(tss1) != 7 {
		t.Errorf("split in 7 error: %+v", tss1)
	}
	if tss1[1].From.Equal(time.Date(2022, 9, 1, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("split in 7 error from")
	}
	if tss1[1].To.Equal(time.Date(2022, 10, 1, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("split in 7 error from")
	}

}

func TestWhatTime(t *testing.T) {

	ts := MakeTimeslice(time.Date(2020, 12, 20, 12, 0, 0, 0, time.UTC), 48*time.Hour)
	dt1 := ts.WhatTime(0.5)
	if !dt1.Equal(time.Date(2020, 12, 21, 12, 0, 0, 0, time.UTC)) {
		t.Errorf("ProgressDate fails: want 20201221 12:00 get %v", dt1)
	}

	dt1 = ts.WhatTime(0.25)
	if !dt1.Equal(time.Date(2020, 12, 21, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("ProgressDate fails: want 20201221 00:00 get %v", dt1)
	}

	ts = MakeTimeslice(time.Date(2020, 12, 20, 14, 35, 0, 0, time.UTC), 1*time.Hour)
	dt1 = ts.WhatTime(0.5)
	if !dt1.Equal(time.Date(2020, 12, 20, 15, 5, 0, 0, time.UTC)) {
		t.Errorf("ProgressDate fails: want 20201220 15:05 get %v", dt1)
	}

	ts = MakeTimeslice(time.Date(2020, 12, 20, 14, 35, 0, 0, time.UTC), 1*time.Hour)
	dt1 = ts.WhatTime(-1)
	if !dt1.Equal(ts.From) {
		t.Errorf("ProgressDate fails: want FROM get %v", dt1)
	}

	ts = MakeTimeslice(time.Date(2020, 12, 20, 14, 35, 0, 0, time.UTC), 1*time.Hour)
	dt1 = ts.WhatTime(10)
	if !dt1.Equal(ts.To) {
		t.Errorf("ProgressDate fails: want TO get %v", dt1)
	}
}

func TestScanSingle(t *testing.T) {
	var get string
	var cursor time.Time

	ts := MakeTimeslice(time.Date(2020, 1, 1, 8, 0, 0, 0, time.UTC), 2*time.Hour)
	mask := 24 * time.Hour

	// get... nothing
	for ts.Scan(&cursor, mask, false); !cursor.IsZero(); ts.Scan(&cursor, mask, false) {
		get += cursor.String()
		fmt.Printf("%v ", cursor)
	}
	if get != "" {
		t.Errorf("Scan fails: %s", get)
	}

	// get... force the single date
	get = ""
	for ts.Scan(&cursor, mask, true); !cursor.IsZero(); ts.Scan(&cursor, mask, true) {
		get += cursor.String()
	}
	if get != "2020-01-01 08:00:00 +0000 UTC2020-01-01 10:00:00 +0000 UTC" {
		t.Errorf("Scan fails: %s", get)
	}

	// get... the single date
	get = ""
	ts = MakeTimeslice(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), 2*time.Hour)
	for ts.Scan(&cursor, mask, false); !cursor.IsZero(); ts.Scan(&cursor, mask, false) {
		get += cursor.String()
	}
	if get != "2020-01-01 00:00:00 +0000 UTC" {
		t.Errorf("Scan fails: %s", get)
	}
}

func TestScanChrono(t *testing.T) {
	var get string
	var cursor time.Time

	ts := MakeTimeslice(time.Date(2020, 1, 1, 8, 0, 0, 0, time.UTC), 2*time.Hour)
	mask := 1 * time.Hour

	// get matching boundaries...
	for ts.Scan(&cursor, mask, false); !cursor.IsZero(); ts.Scan(&cursor, mask, false) {
		get += cursor.String()
	}
	if get != "2020-01-01 08:00:00 +0000 UTC2020-01-01 09:00:00 +0000 UTC2020-01-01 10:00:00 +0000 UTC" {
		t.Errorf("Scan fails: %s", get)
	}

	// ...and unget not matching ones
	ts.From = ts.From.Add(35 * time.Minute)

	get = ""
	cursor = time.Time{}
	for ts.Scan(&cursor, mask, false); !cursor.IsZero(); ts.Scan(&cursor, mask, false) {
		get += cursor.String()
	}
	if get != "2020-01-01 09:00:00 +0000 UTC2020-01-01 10:00:00 +0000 UTC" {
		t.Errorf("Scan fails: %s", get)
	}

	// get boundaries even if not matching the mask
	ts.To = ts.To.Add(-20 * time.Minute)

	get = ""
	cursor = time.Time{}
	for ts.Scan(&cursor, mask, true); !cursor.IsZero(); ts.Scan(&cursor, mask, true) {
		get += cursor.String()
	}
	if get != "2020-01-01 08:35:00 +0000 UTC2020-01-01 09:00:00 +0000 UTC2020-01-01 09:40:00 +0000 UTC" {
		t.Errorf("Scan fails: %s", get)
	}
}

func TestScanAntiChrono(t *testing.T) {
	var get string
	var cursor time.Time

	ts := MakeTimeslice(time.Date(2020, 1, 31, 8, 0, 0, 0, time.UTC), -2*time.Hour)
	mask := 1 * time.Hour

	// get matching boundaries...
	for ts.Scan(&cursor, mask, false); !cursor.IsZero(); ts.Scan(&cursor, mask, false) {
		get += cursor.String()
	}
	if get != "2020-01-31 08:00:00 +0000 UTC2020-01-31 07:00:00 +0000 UTC2020-01-31 06:00:00 +0000 UTC" {
		t.Errorf("Scan fails: %s", get)
	}

	// ...and unget not matching ones
	ts.From = ts.From.Add(-35 * time.Minute)

	get = ""
	cursor = time.Time{}
	for ts.Scan(&cursor, mask, false); !cursor.IsZero(); ts.Scan(&cursor, mask, false) {
		get += cursor.String()
	}
	if get != "2020-01-31 07:00:00 +0000 UTC2020-01-31 06:00:00 +0000 UTC" {
		t.Errorf("Scan fails: %s", get)
	}

	// get boundaries even if not matching the mask
	ts.To = ts.To.Add(20 * time.Minute)

	get = ""
	cursor = time.Time{}
	for ts.Scan(&cursor, mask, true); !cursor.IsZero(); ts.Scan(&cursor, mask, true) {
		get += cursor.String()
	}
	if get != "2020-01-31 07:25:00 +0000 UTC2020-01-31 07:00:00 +0000 UTC2020-01-31 06:20:00 +0000 UTC" {
		t.Errorf("Scan fails: %s", get)
	}
}
