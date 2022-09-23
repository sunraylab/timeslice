// Copyright 2022 by lolorenzo77. All rights reserved.
// Use of this source code is governed by MIT licence that can be found in the LICENSE file.

package timeline

import (
	"testing"
	"time"
)

func TestSplit(t *testing.T) {

	// a timeslice staring 20220801 0h00:00 and 7 days long
	ts := MakeTimeslice(time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC), 7*Day)

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
		t.Errorf("ProgressDate fails: want 20201221 12:00 got %v", dt1)
	}

	dt1 = ts.WhatTime(0.25)
	if !dt1.Equal(time.Date(2020, 12, 21, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("ProgressDate fails: want 20201221 00:00 got %v", dt1)
	}

	ts = MakeTimeslice(time.Date(2020, 12, 20, 14, 35, 0, 0, time.UTC), 1*time.Hour)
	dt1 = ts.WhatTime(0.5)
	if !dt1.Equal(time.Date(2020, 12, 20, 15, 5, 0, 0, time.UTC)) {
		t.Errorf("ProgressDate fails: want 20201220 15:05 got %v", dt1)
	}

	ts = MakeTimeslice(time.Date(2020, 12, 20, 14, 35, 0, 0, time.UTC), 1*time.Hour)
	dt1 = ts.WhatTime(-1)
	if !dt1.Equal(ts.From) {
		t.Errorf("ProgressDate fails: want FROM got %v", dt1)
	}

	ts = MakeTimeslice(time.Date(2020, 12, 20, 14, 35, 0, 0, time.UTC), 1*time.Hour)
	dt1 = ts.WhatTime(10)
	if !dt1.Equal(ts.To) {
		t.Errorf("ProgressDate fails: want TO got %v", dt1)
	}
}

func TestMiddle(t *testing.T) {
	tim := time.Date(2020, 12, 20, 14, 35, 0, 0, time.UTC)
	ts := MakeTimeslice(tim, 1*time.Hour)
	got := ts.Middle()
	if !got.Equal(time.Date(2020, 12, 20, 15, 5, 0, 0, time.UTC)) {
		t.Errorf("Middle fails: want 2020-12-20 15:05:00 got %v", got)
	}

	got = TimeSlice{}.Middle()
	if !got.IsZero() {
		t.Errorf("Middle fails: want zero time got %v", got)
	}

	got = TimeSlice{From: tim, To: tim}.Middle()
	if !got.Equal(tim) {
		t.Errorf("Middle fails: want %v got %v", tim, got)
	}
}
