package timex

import (
	"time"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type TimeDelta struct {
	Weeks        int
	Days         int
	Hours        int
	Minutes      int
	Seconds      int
	Microseconds int
	Milliseconds int
}

// Duration returns time.Duration can be added to time.Date.
func (td *TimeDelta) Duration() time.Duration {
	return time.Duration(td.Weeks)*7*24*time.Hour +
		time.Duration(td.Days)*24*time.Hour +
		time.Duration(td.Hours)*time.Hour +
		time.Duration(td.Minutes)*time.Minute +
		time.Duration(td.Seconds)*time.Second +
		time.Duration(td.Microseconds)*time.Microsecond +
		time.Duration(td.Milliseconds)*time.Millisecond

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
		Weeks:        abs(td.Weeks),
		Days:         abs(td.Days),
		Hours:        abs(td.Hours),
		Minutes:      abs(td.Minutes),
		Seconds:      abs(td.Seconds),
		Milliseconds: abs(td.Milliseconds),
		Microseconds: abs(td.Microseconds),
	}
}

// String returns a string representing the TimeDelta's duration in the form "72h3m0.5s".
func (td *TimeDelta) String() string {
	return td.Duration().String()
}
