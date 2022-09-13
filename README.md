# Timeline module in go

Timeline module provides 2 packages

- Duration : which is an extension of the default time.Duration package.
- TimeSlice : TimeSlice represents a range of times bounded by two dates (time.Time) From and To. It accepts infinite boundaries (zero times) and can be chronological or anti-chronological.

## Usage

### Duration 

`Duration` type extends default `time.Duration` with additional Days, Month and Year methods, and with a special formating function.

`Days`, ``Month`` and ``Year`` methods provides usefull calculations based on the following typical time values:

	1 day       = 24 hours
	1 year      = 365.25 days in average because of leap years = 8766 hours
	1 month     = 1 year / 12 = 30.4375 days = 730.5 hours
	1 quarter   = 1 year / 4 = 121.75 days = 2922 hours

Duration provides also a special formating function `FormatOrderOfMagnitude` to produce an output to give an order of magnitude of the duration, at a choosen order, in a human-reading way like ``[-][100Y][12M][31d][24h][60m][99s][~]`` with:

- only non-zero components are output
- starting with the biggest time components available, only the following choosen number of components are output
- all trailing non zero components produce the `~` output meaning the duration in not exact

For example 
```go 
	// more than one month
	d7 := Duration(1*Month + 4*Day + 2*time.Hour + 35*time.Minute + 25*time.Second)
    // Format with Order Of Magnitude of 3
    // here the output starts with Months (because there is no years)
    // so an order of magnitude of 3 gives 3 time components from Months: Month, Days, and Hours.
    // considering lower orders not significant, or dust
	fmt.Printf("%s\n", d7.FormatOrderOfMagnitude(3)) // prints: 1M4d2h~
```

## Timeslice

A timeSlice can be easily created with literal values:

```go
ts := &TimeSlice{
         From: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
         To:   time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
```

A TimeSlice can also be created with a factory function with a defined d duration and a starting time.
```go
ts := BuildTimeSlice(time.Now(), 1 * Day)
```

TimeSlice provides basic functions to proceed with infinite boundaries (zero times) and with chronological or anti-chronological direction.

In addition a TimeSlice provides advanced features: 
- a TimeSlice can be splitted in sub timelices of a specified duration
- you can get the exact time at a selected progress rate within the timeslice boundaries
- considering a certain time, you can get its position within the timeslice boundaries
- a TimeSlice can be scanned with a mask to go through all its starting minutes, or all its starting hours...

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
go get -u github.com/sunraylab/timeline@latest
```

## Changelog

- v1.1.0 : 
  - provides the TimeMask type 
  - fix Scan function
  - add function String() to Duration
  - add function Adjust() to Duration

## Licence

[MIT](LICENSE)