// Copyright 2022 by lolorenzo77. All rights reserved.
// Use of this source code is governed by MIT licence that can be found in the LICENSE file.

package timeline

import (
	"testing"
	"time"
)

func TestSplit(t *testing.T) {

	// a timeslice staring 20220801 0h00:00 and 7 days long
	ts := MakeTimeSlice(time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC), 7*Day)

	// split in one
	tss0, err := ts.Split(7 * 24 * time.Hour)
	if err != nil || len(tss0) != 1 || tss0[0].Compare(ts) != 1 {
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

	ts := MakeTimeSlice(time.Date(2020, 12, 20, 12, 0, 0, 0, time.UTC), 48*time.Hour)
	dt1 := ts.WhatTime(0.5)
	if !dt1.Equal(time.Date(2020, 12, 21, 12, 0, 0, 0, time.UTC)) {
		t.Errorf("ProgressDate fails: want 20201221 12:00 got %v", dt1)
	}

	dt1 = ts.WhatTime(0.25)
	if !dt1.Equal(time.Date(2020, 12, 21, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("ProgressDate fails: want 20201221 00:00 got %v", dt1)
	}

	ts = MakeTimeSlice(time.Date(2020, 12, 20, 14, 35, 0, 0, time.UTC), 1*time.Hour)
	dt1 = ts.WhatTime(0.5)
	if !dt1.Equal(time.Date(2020, 12, 20, 15, 5, 0, 0, time.UTC)) {
		t.Errorf("ProgressDate fails: want 20201220 15:05 got %v", dt1)
	}

	ts = MakeTimeSlice(time.Date(2020, 12, 20, 14, 35, 0, 0, time.UTC), 1*time.Hour)
	dt1 = ts.WhatTime(-1)
	if !dt1.Equal(ts.From) {
		t.Errorf("ProgressDate fails: want FROM got %v", dt1)
	}

	ts = MakeTimeSlice(time.Date(2020, 12, 20, 14, 35, 0, 0, time.UTC), 1*time.Hour)
	dt1 = ts.WhatTime(10)
	if !dt1.Equal(ts.To) {
		t.Errorf("ProgressDate fails: want TO got %v", dt1)
	}
}

func TestMiddle(t *testing.T) {
	tim := time.Date(2020, 12, 20, 14, 35, 0, 0, time.UTC)
	ts := MakeTimeSlice(tim, 1*time.Hour)
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

func TestBound(t *testing.T) {

	tim := time.Date(2020, 12, 20, 14, 35, 0, 0, time.UTC)
	timp2 := tim.Add(Day * 2)
	timm2 := tim.Add(-Day * 2)

	// infinite ts
	tsToInfinite := TimeSlice{From: tim}
	tbound := tsToInfinite.Bound(timp2)
	if !tbound.Equal(timp2) {
		t.Errorf("Bound fails: got %v", tbound)
	}

	tsFromInfinite := TimeSlice{To: tim}
	tbound = tsFromInfinite.Bound(timp2)
	if !tbound.Equal(tim) {
		t.Errorf("Bound fails: got %v", tbound)
	}

	// Finite chrono timeslice
	tsFiniteChrono := MakeTimeSlice(tim, Day)
	tbound = tsFiniteChrono.Bound(timp2)
	if !tbound.Equal(tsFiniteChrono.To) {
		t.Errorf("Bound fails: got %v", tbound)
	}
	tbound = tsFiniteChrono.Bound(timm2)
	if !tbound.Equal(tsFiniteChrono.From) {
		t.Errorf("Bound fails: got %v", tbound)
	}

	// Finite antichrono timeslice
	tsFiniteAntoChrono := MakeTimeSlice(tim, -Day)
	tbound = tsFiniteAntoChrono.Bound(timm2)
	if !tbound.Equal(tsFiniteAntoChrono.To) {
		t.Errorf("Bound fails: got %v", tbound)
	}
	tbound = tsFiniteAntoChrono.Bound(timp2)
	if !tbound.Equal(tsFiniteAntoChrono.From) {
		t.Errorf("Bound fails: got %v", tbound)
	}

}

func TestMove(t *testing.T) {

	tim := time.Date(2020, 12, 20, 14, 35, 0, 0, time.UTC)
	timp2 := tim.Add(Day * 2)
	timm2 := tim.Add(-Day * 2)

	// infinite ts
	tsToInfinite := TimeSlice{From: tim}
	if !tsToInfinite.MoveToAt(timp2).To.Equal(timp2) {
		t.Errorf("ToMove fails: got %v", tsToInfinite)
	}

	tsFromInfinite := TimeSlice{To: tim}
	if !tsFromInfinite.MoveToAt(timp2).To.Equal(timp2) {
		t.Errorf("ToMove fails: got %v", tsFromInfinite)
	}

	// Finite chrono timeslice
	tsFiniteChrono := MakeTimeSlice(tim, Day)
	if !tsFiniteChrono.MoveFromAt(timp2).To.Equal(timp2) {
		t.Errorf("FromMove fails: got %v", tsFiniteChrono)
	}
	if !tsFiniteChrono.MoveFromAt(tim).To.Equal(timp2) {
		t.Errorf("FromMove fails: got %v", tsFiniteChrono)
	}
	if !tsFiniteChrono.MoveToAt(timm2).From.Equal(timm2) {
		t.Errorf("ToMove fails: got %v", tsFiniteChrono)
	}
	if !tsFiniteChrono.MoveToAt(tim).From.Equal(timm2) {
		t.Errorf("ToMove fails: got %v", tsFiniteChrono)
	}

	// Finite antichrono timeslice
	tsFiniteAntoChrono := MakeTimeSlice(tim, -Day)
	if !tsFiniteAntoChrono.MoveFromAt(timm2).To.Equal(timm2) {
		t.Errorf("FromMove fails: got %v", tsFiniteAntoChrono)
	}
	if !tsFiniteAntoChrono.MoveFromAt(tim).To.Equal(timm2) {
		t.Errorf("FromMove fails: got %v", tsFiniteAntoChrono)
	}
	if !tsFiniteAntoChrono.MoveToAt(timp2).From.Equal(timp2) {
		t.Errorf("ToMove fails: got %v", tsFiniteAntoChrono)
	}
	if !tsFiniteAntoChrono.MoveToAt(tim).From.Equal(timp2) {
		t.Errorf("ToMove fails: got %v", tsFiniteAntoChrono)
	}

}

// Test all combinations of ShiftIn
func FuzzShiftIn(f *testing.F) {

	// corpus
	for shiftf := -1; shiftf <= 1; shiftf++ {
		for shiftt := -1; shiftt <= 1; shiftt++ {
			for boundf := -1; boundf <= 1; boundf++ {
				for boundt := -1; boundt <= 1; boundt++ {
					f.Add(shiftf, shiftt, boundf, boundt)
				}
			}
		}
	}

	n := 0
	// target
	f.Fuzz(func(t *testing.T, shiftf int, shiftt int, boundf int, boundt int) {
		n++
		tt := time.Date(2022, 06, 10, 0, 0, 0, 0, time.UTC)

		tsbound := TimeSlice{}
		switch boundf {
		case 1:
			tsbound.From = tt.Add(Day * 0)
		case 0: // stay infinite
		case -1:
			tsbound.From = tt.Add(Day * 14)
		}
		switch boundt {
		case 1:
			tsbound.To = tt.Add(Day * 14)
		case 0: // stay infinite
		case -1:
			tsbound.To = tt.Add(Day * 0)
		}

		tsshift := TimeSlice{}
		switch shiftf {
		case 1:
			tsshift.From = tt.Add(Day * 2)
		case 0: // stay infinite
		case -1:
			tsshift.From = tt.Add(Day * 12)
		}
		switch shiftt {
		case 1:
			tsshift.To = tt.Add(Day * 12)
		case 0: // stay infinite
		case -1:
			tsshift.To = tt.Add(Day * 2)
		}

		tsshift0 := tsshift
		tsboundres := tsbound
		tsboundres.ForceDirection(Chronological)

		if tsshift.ShiftIn(Day*10, tsbound) != nil {
			tsshiftres := tsshift
			tsshiftres.ForceDirection(Chronological)

			if !tsshift0.IsInfinite() && tsshift0.Direction() != tsshift.Direction() {
				t.Errorf("[%d] ShiftIn Direction Fails on finite ts; 0:%v, bound:%v, shifted:%v", n, tsshift0, tsbound, tsshift)
			}
			if !tsshiftres.To.IsZero() && !tsboundres.To.IsZero() && tsshiftres.To.After(tsboundres.To) {
				t.Errorf("[%d] ShiftIn To is out of boundaries; 0:%v, bound:%v, shifted:%v", n, tsshift0, tsbound, tsshift)
			}
			if !tsshiftres.From.IsZero() && !tsboundres.From.IsZero() && tsshiftres.From.Before(tsboundres.From) {
				t.Errorf("[%d] ShiftIn From is out of boundaries; 0:%v, bound:%v, shifted:%v", n, tsshift0, tsbound, tsshift)
			}
		} else {
			if !(tsshift0.Duration().Abs().Duration > tsboundres.Duration().Abs().Duration) {
				t.Errorf("[%d] ShiftIn returs nil; 0:%v, bound:%v, shifted:%v", n, tsshift0, tsbound, tsshift)
			}
		}
	})
}
