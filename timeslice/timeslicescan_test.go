package timeslice

import (
	"fmt"
	"testing"
	"time"

	"github.com/sunraylab/timeline/duration"
)

func TestScanSingle(t *testing.T) {
	var get string
	var cursor time.Time

	ts := MakeTimeslice(time.Date(2020, 1, 1, 8, 0, 0, 0, time.UTC), 2*time.Hour)
	mask := MASK_DAY

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

func TestScanChrono1(t *testing.T) {
	var get string
	var cursor time.Time

	ts := MakeTimeslice(time.Date(2020, 1, 1, 8, 0, 0, 0, time.UTC), 2*time.Hour)
	mask := MASK_HOUR

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

func TestScanChrono2(t *testing.T) {
	var get string
	var cursor time.Time

	ts := MakeTimeslice(time.Date(2022, 1, 6, 7, 30, 0, 0, time.UTC), duration.Month*3)
	mask := MASK_MONTH

	// get matching boundaries...
	for ts.Scan(&cursor, mask, false); !cursor.IsZero(); ts.Scan(&cursor, mask, false) {
		get += cursor.String()
	}
	if get != "2022-02-01 00:00:00 +0000 UTC2022-03-01 00:00:00 +0000 UTC2022-04-01 00:00:00 +0000 UTC" {
		t.Errorf("Scan fails: %s", get)
	}
}

func TestScanAntiChrono(t *testing.T) {
	var get string
	var cursor time.Time

	ts := MakeTimeslice(time.Date(2020, 1, 31, 8, 0, 0, 0, time.UTC), -2*time.Hour)
	mask := MASK_HOUR

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
