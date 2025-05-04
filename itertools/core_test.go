package itertools

import (
	"testing"
)

func Test_MapSliceArray(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		fn       func(int) int
		expected []int
	}{
		{"test1", []int{}, func(x int) int { return x * 3 }, []int{}},
		{"test2", []int{1, 2, 3, 4}, func(x int) int { return x * 3 }, []int{3, 6, 9, 12}},
		{"test3", []int{1, 2, 3, 4, 5, 6}, func(x int) int { return x % 2 }, []int{1, 0, 1, 0, 1, 0}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := Map(tc.slice, tc.fn); !equal(got, tc.expected) {
				t.Errorf("Map() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func Test_FilterSlice(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		fn       func(int) bool
		expected []int
	}{
		{"test1", []int{}, func(x int) bool { return x%3 == 0 }, []int{}},
		{"test2", []int{1, 2, 3, 4}, func(x int) bool { return x > 3 }, []int{4}},
		{"test3", []int{1, 2, 3, 4, 5, 6}, func(x int) bool { return x%3 == 0 }, []int{3, 6}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := Filter(tc.slice, tc.fn); !equal(got, tc.expected) {
				t.Errorf("Filter() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func Test_AllSlice(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		fn       func(int) bool
		expected bool
	}{
		{"test1", []int{}, func(x int) bool { return x%2 == 0 }, true},
		{"test2", []int{2, 4, 6}, func(x int) bool { return x%2 == 0 }, true},
		{"test3", []int{2, 3, 6}, func(x int) bool { return x%2 == 0 }, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := All(tc.slice, tc.fn); got != tc.expected {
				t.Errorf("All() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func Test_AnySlice(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		fn       func(int) bool
		expected bool
	}{
		{"test1", []int{}, func(x int) bool { return x%2 == 0 }, false},
		{"test2", []int{1, -2, -3, -4, -5}, func(x int) bool { return x > 0 }, true},
		{"test3", []int{1, 2, 4, 5}, func(x int) bool { return x%3 == 0 }, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := Any(tc.slice, tc.fn); got != tc.expected {
				t.Errorf("Any() expected %v, got %v", tc.expected, got)
			}
		})
	}
}
