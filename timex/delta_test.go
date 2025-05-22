package timex_test

import (
	"github.com/shijl0925/go-toolkits/timex"
	"testing"
	"time"
)

// -----------------------------
// Duration 方法测试
// -----------------------------

func TestTimeDelta_Duration(t *testing.T) {
	tests := []struct {
		name string
		td   timex.TimeDelta
		want time.Duration
	}{
		{
			name: "only seconds",
			td:   timex.TimeDelta{Seconds: 30},
			want: 30 * time.Second,
		},
		{
			name: "mixed units",
			td: timex.TimeDelta{
				Weeks:        1,
				Days:         2,
				Hours:        3,
				Minutes:      4,
				Seconds:      5,
				Milliseconds: 6,
				Microseconds: 7,
			},
			want: time.Duration(1*7*24*time.Hour +
				2*24*time.Hour +
				3*time.Hour +
				4*time.Minute +
				5*time.Second +
				6*time.Millisecond +
				7*time.Microsecond),
		},
		{
			name: "zero value",
			td:   timex.TimeDelta{},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.td.Duration(); got != tt.want {
				t.Errorf("Duration() = %v; want %v", got, tt.want)
			}
		})
	}
}

// -----------------------------
// Add 方法测试
// -----------------------------

func TestTimeDelta_Add(t *testing.T) {
	td1 := timex.TimeDelta{Days: 2, Hours: 3}
	td2 := timex.TimeDelta{Days: 1, Hours: 4}
	expected := timex.TimeDelta{Days: 3, Hours: 7}

	result := td1.Add(&td2)
	if result != expected {
		t.Errorf("Add() = %v; want %v", result, expected)
	}
}

// -----------------------------
// Subtract 方法测试
// -----------------------------

func TestTimeDelta_Subtract(t *testing.T) {
	td1 := timex.TimeDelta{Days: 5, Hours: 3}
	td2 := timex.TimeDelta{Days: 2, Hours: 10}
	expected := timex.TimeDelta{Days: 3, Hours: -7}

	result := td1.Subtract(&td2)
	if result != expected {
		t.Errorf("Subtract() = %v; want %v", result, expected)
	}
}

// -----------------------------
// Abs 方法测试
// -----------------------------

func TestTimeDelta_Abs(t *testing.T) {
	td := timex.TimeDelta{Days: -2, Hours: 3}
	expected := timex.TimeDelta{Days: 2, Hours: 3}

	result := td.Abs()
	if result != expected {
		t.Errorf("Abs() = %v; want %v", result, expected)
	}
}

// -----------------------------
// String 方法测试
// -----------------------------

func TestTimeDelta_String(t *testing.T) {
	td := timex.TimeDelta{Hours: 2, Minutes: 3, Seconds: 1}
	expected := td.Duration().String()

	if got := td.String(); got != expected {
		t.Errorf("String() = %q; want %q", got, expected)
	}
}
