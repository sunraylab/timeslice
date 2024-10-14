// Copyright 2022-2024 by larry868. All rights reserved.
// Use of this source code is governed by MIT licence that can be found in the LICENSE file.

package timeline

import (
	"fmt"
	"time"
)

func ExampleDuration_Days() {
	d0 := NewDuration(1 * time.Hour)
	d1 := NewDuration(24 * time.Hour)
	d2 := NewDuration(25 * time.Hour)
	fmt.Printf("%f\n%f\n%f", d0.Days(), d1.Days(), d2.Days())
	// Output:
	// 0.041667
	// 1.000000
	// 1.041667
}

func ExampleDuration_Months() {
	// a single day is a fraction of a month
	d0 := NewDuration(24 * time.Hour)
	// 1 average month is 30.4375 days = 730.5 hours
	n1 := Nanoseconds(730.5 * float64(time.Hour))
	// 1 average year is 365.25 days = 8766 hours
	n2 := NewDuration(8766 * time.Hour)
	fmt.Printf("%f\n%f\n%f\n", d0.Months(), n1.Months(), n2.Months())
	// Output:
	// 0.032854
	// 1.000000
	// 12.000000
}

func ExampleDuration_Years() {
	// half year (1 average month is 30.4375 days)
	n1 := Nanoseconds(6 * 30.4375 * 24 * float64(time.Hour))
	// one full year
	n2 := Nanoseconds(365.25 * 24 * float64(time.Hour))
	fmt.Printf("%f\n%f\n", n1.Years(), n2.Years())
	// Output:
	// 0.500000
	// 1.000000
}

func ExampleDuration_FormatOrderOfMagnitude() {

	// infifnite
	var d Duration
	fmt.Printf("%s\n", d.FormatOrderOfMagnitude(3))
	// not significant
	d1 := NewDuration(1 * time.Millisecond)
	fmt.Printf("%s\n", d1.FormatOrderOfMagnitude(3))
	// some seconds
	d2 := NewDuration(10*time.Second + 75*time.Millisecond)
	fmt.Printf("%s\n", d2.FormatOrderOfMagnitude(3))
	// some minutes without seconds
	d3 := NewDuration(15 * time.Minute)
	fmt.Printf("%s\n", d3.FormatOrderOfMagnitude(3))
	// on day with some minutes but no hours
	d4 := NewDuration(Day + 25*time.Minute)
	fmt.Printf("%s\n", d4.FormatOrderOfMagnitude(3))
	// the same but with an order of magnitude of 2
	fmt.Printf("%s\n", d4.FormatOrderOfMagnitude(2))
	// one month
	d6 := NewDuration(1 * Month)
	fmt.Printf("%s\n", d6.FormatOrderOfMagnitude(3))
	// more than one month, leaving dust of minutes & seconds
	d7 := NewDuration(1*Month + 4*Day + 2*time.Hour + 35*time.Minute + 25*time.Second)
	fmt.Printf("%s\n", d7.FormatOrderOfMagnitude(3))

	// Output:
	// infinite
	// 0s~
	// 10s
	// 15m
	// 1d25m
	// 1d~
	// 1M
	// 1M4d2h~

}
