# timeline package in go

timeline package provides 3 main types:

- Duration: which is an extension of the default time.Duration struct.
- TimeSlice: representing a range of times bounded by two dates (time.Time) From and To. It accepts infinite boundaries (zero times) and can be chronological or anti-chronological.
- TimeMask: used for scanning a TimeSlice and to get the time corresponding to a rounding o'clock period.

[![Go Reference](https://pkg.go.dev/badge/github.com/sunraylab/timeline/v2.svg)](https://pkg.go.dev/github.com/sunraylab/timeline/v2)

## Usage

### Duration 

`Duration` type extends default `time.Duration` with additional Days, Month and Year methods, and with a special formating function.

`Days`, `Month` and `Year` methods provides usefull calculations based on the following typical time values:

	1 day       = 24 hours
	1 year      = 365.25 days in average because of leap years = 8766 hours
	1 month     = 1 year / 12 = 30.4375 days = 730.5 hours
	1 quarter   = 1 year / 4 = 121.75 days = 2922 hours

Duration provides also a special formating function `FormatOrderOfMagnitude` to produce an output to give an order of magnitude of the duration, at a choosen order, in a human-reading way like ``[-][100Y][12M][31d][24h][60m][99s][~]`` with:

- only non-zero components are output
- starting with the biggest time components available, only the following choosen number of components are output
- all trailing non zero components produce the `~` output meaning the duration in not exact

Example 

```go 
// more than one month
d7 := Duration(1*Month + 4*Day + 2*time.Hour + 35*time.Minute + 25*time.Second)

// Format with Order Of Magnitude of 3
// here the output starts with Months (because there is no years)
// so an order of magnitude of 3 gives 3 time components from Months: Month, Days, and Hours.
// considering lower orders not significant, or dust
fmt.Printf("%s\n", d7.FormatOrderOfMagnitude(3)) // prints: 1M4d2h~
```

## TimeSlice

A TimeSlice can be easily created with literal values:

```go
ts := &TimeSlice{
         From: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
         To:   time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),}
```

A TimeSlice can also be created with a factory function with a defined d duration and a starting time.

```go
ts := MakeTimeSlice(time.Now(), 1 * Day)
```

TimeSlice provides basic functions to proceed with infinite boundaries (zero times) and with chronological or anti-chronological direction.

TimeSlice advanced features: 
- TimeSlice can be splitted in sub timelices of a specified duration
- get the exact time at a selected progress rate within the timeslice boundaries
- considering a certain time, get its position within the timeslice boundaries
- TimeSlice can be scanned with a mask to go through all its starting minutes, all its starting hours...

## TimeMask 

The TimeMask type provides the following scanning possibilities:
```go
	MASK_MINUTE    
	MASK_MINUTEx15 
	MASK_HALFHOUR 
	MASK_HOUR      
	MASK_HOURx4    
	MASK_HALFDAY   
	MASK_DAY       
	MASK_MONTH     
	MASK_QUARTER   
	MASK_YEAR      
```

## Installing 

```bash 
go get -u github.com/sunraylab/timeline/v2@latest
```

## Changelog

- v2.4.x:
  - new features TimeSlice.Query() and ParseFromToQuery()
  - new feature TimeSlice.IsOverlapping()
  - new features TimeSlice.FormatFrom and TimeSlice.FormatTo
  - bug fix on GetTimeFormat
  - new time min/max helpers

- v2.3.1:
  - new feature TimeSlice.BoundIn() 
  - fix bug on WhereIs with antichonological timeslice 

- v2.3.0:
  - new feature TimeSlice.ForceDirection() 
  - new feature TimeSlice.ShiftIn() 
  - new feature TimeSlice.Bound() 
  - new feature Duration.Abd() 
  - timeslice can now be used with chaining method
  - change signature to use time.Duration rather than timeline.Duration every time we use it as a a parameter
  - Extend and Move functions have been renamed with new signatures
  - fix bug on Direction() when From boundaries was infinite
  - fix bug on String where To boundaries was badly formated in some cases

- v2.2.2:
  - add feature timeslice.WhereIs

- v2.2.1:
  - break with the v2.1, major change with Duration which is now embeding time.Duration rather than to be of his type. Add feat. Infinite durations.
  add feat. timeslice.IsZero

- v2.1.3:
  - added IsZero

- v2.1.2 : 
  - String return an empty string for nil duration

- v2.1.0:
  - Equal is replaced by Compare and returns a type, no more an int. + some refactoring

- v2.0.0:
  - the module has been streamlined with its differents packages merged into the timeline package
  - func Shift() has been added to the TimeSlice
  - func MoveTo, MoveFrom, ExtendTo, ExtendFrom have been renamed in ToMove, FromMove, ToExtend, FromExtend to be less confusing

- v1.3.0: 
  - go 1.19.1
  - updating the doc and adding copyright info

- v1.2.1: 
  - fix scanmask

- v1.2.0: 
  - add function GetTimeFormat() to TimeMask

- v1.1.0: 
  - provides the TimeMask type 
  - fix Scan function
  - add function String() to Duration
  - add function Adjust() to Duration

## Licence

[MIT](LICENSE)