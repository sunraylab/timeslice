package duration

import (
	"testing"
	"time"
)

func TestDuration1(t *testing.T) {

	durD := Duration(1 * Month)
	m := time.Duration(durD).Minutes()
	h := time.Duration(durD).Hours()
	d := durD.Days()
	M := durD.Months()
	Y := durD.Years()
	if m != 43830 || h != 730.5 || d != 30.4375 || M != 1.0 || int(Y*10000) != 833 {
		t.Errorf("Duration Fails: %f, %f, %f, %f, %f", m, h, d, M, Y)
	}
}

func TestDuration2(t *testing.T) {

	t1 := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)
	dur1 := t2.Sub(t1)
	dur2 := time.Duration(2 * Day)

	if dur1 != dur2 {
		t.Errorf("Duration Fails: %v, %v, %v, %v", t1, t2, dur1, dur2)
	}
}
