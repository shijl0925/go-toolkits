package timex

import (
	"time"

	"github.com/shijl0925/go-toolkits/internal"
)

type TimeDelta struct {
	Weeks        int
	Days         int
	Hours        int
	Minutes      int
	Seconds      int
	Microseconds int
	Milliseconds int
}

// Duration returns time.Duration that can be added to a time.Time.
// It returns an error if the calculation results in an overflow of time.Duration.
func (td *TimeDelta) Duration() (time.Duration, error) {
	var totalDuration int64 // Accumulate in nanoseconds

	// Helper to add a component and check for overflow
	addComponent := func(currentTotal int64, value int, unitNanos int64) (int64, error) {
		if value == 0 {
			return currentTotal, nil
		}
		termNanos := int64(value) * unitNanos

		// Check for multiplication overflow for the term itself
		if value != 0 && unitNanos != 0 && termNanos/unitNanos != int64(value) {
			return 0, OverflowError // Multiplication overflowed
		}

		// Check for addition overflow
		newTotal := currentTotal + termNanos
		if (termNanos > 0 && newTotal < currentTotal) || (termNanos < 0 && newTotal > currentTotal) {
			return 0, OverflowError // Addition overflowed
		}
		return newTotal, nil
	}

	var err error
	totalDuration, err = addComponent(totalDuration, td.Weeks, int64(7*24*time.Hour))
	if err != nil {
		return 0, err
	}

	totalDuration, err = addComponent(totalDuration, td.Days, int64(24*time.Hour))
	if err != nil {
		return 0, err
	}

	totalDuration, err = addComponent(totalDuration, td.Hours, int64(time.Hour))
	if err != nil {
		return 0, err
	}

	totalDuration, err = addComponent(totalDuration, td.Minutes, int64(time.Minute))
	if err != nil {
		return 0, err
	}

	totalDuration, err = addComponent(totalDuration, td.Seconds, int64(time.Second))
	if err != nil {
		return 0, err
	}

	totalDuration, err = addComponent(totalDuration, td.Milliseconds, int64(time.Millisecond))
	if err != nil {
		return 0, err
	}

	totalDuration, err = addComponent(totalDuration, td.Microseconds, int64(time.Microsecond))
	if err != nil {
		return 0, err
	}

	return time.Duration(totalDuration), nil
}

// Add returns the TimeDelta td+td2.
func (td *TimeDelta) Add(td2 *TimeDelta) TimeDelta {
	return TimeDelta{
		Weeks:        td.Weeks + td2.Weeks,
		Days:         td.Days + td2.Days,
		Hours:        td.Hours + td2.Hours,
		Minutes:      td.Minutes + td2.Minutes,
		Seconds:      td.Seconds + td2.Seconds,
		Milliseconds: td.Milliseconds + td2.Milliseconds,
		Microseconds: td.Microseconds + td2.Microseconds,
	}
}

// Subtract returns the TimeDelta td-td2.
func (td *TimeDelta) Subtract(td2 *TimeDelta) TimeDelta {
	return TimeDelta{
		Weeks:        td.Weeks - td2.Weeks,
		Days:         td.Days - td2.Days,
		Hours:        td.Hours - td2.Hours,
		Minutes:      td.Minutes - td2.Minutes,
		Seconds:      td.Seconds - td2.Seconds,
		Milliseconds: td.Milliseconds - td2.Milliseconds,
		Microseconds: td.Microseconds - td2.Microseconds,
	}
}

// Abs returns the absolute value of td
func (td *TimeDelta) Abs() TimeDelta {
	return TimeDelta{
		Weeks:        internal.Abs(td.Weeks),
		Days:         internal.Abs(td.Days),
		Hours:        internal.Abs(td.Hours),
		Minutes:      internal.Abs(td.Minutes),
		Seconds:      internal.Abs(td.Seconds),
		Milliseconds: internal.Abs(td.Milliseconds),
		Microseconds: internal.Abs(td.Microseconds),
	}
}

// String returns a string representing the TimeDelta's duration in the form "72h3m0.5s".
// If an overflow occurs when calculating the duration, it returns an error message string.
func (td *TimeDelta) String() string {
	d, err := td.Duration()
	if err != nil {
		return "TimeDelta(overflow)" // Or some other indicator of error
	}
	return d.String()
}
