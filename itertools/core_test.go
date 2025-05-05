package itertools_test

import (
	"github.com/shijl0925/go-toolkits/itertools"
	"reflect"
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
			if got := itertools.Map(tc.slice, tc.fn); !reflect.DeepEqual(got, tc.expected) {
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
			if got := itertools.Filter(tc.slice, tc.fn); !reflect.DeepEqual(got, tc.expected) {
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
			if got := itertools.All(tc.slice, tc.fn); got != tc.expected {
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
			if got := itertools.Any(tc.slice, tc.fn); got != tc.expected {
				t.Errorf("Any() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func Test_ReduceSlice(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		fn       func(int, int) int
		initial  int
		expected int
	}{
		{"test1", []int{1, 2, 3, 4, 5}, func(x, y int) int { return x + y }, 0, 15},
		{"test2", []int{1, 3, 5, 6, 2}, func(x, y int) int { return x * y }, 1, 180},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := itertools.Reduce(tc.slice, tc.fn, tc.initial); got != tc.expected {
				t.Errorf("Reduce() expected %v, got %v", tc.expected, got)
			}
		})
	}
	t.Run("test3", func(t *testing.T) {
		if got := itertools.Reduce([]string{"geeks", "for", "geeks"}, func(x, y string) string { return x + y }, "!"); got != "!geeksforgeeks" {
			t.Errorf("Reduce() expected %v, got %v", "!geeksforgeeks", got)
		}
	})
}

func Test_Zip(t *testing.T) {
	tests := []struct {
		name   string
		slice1 []int
		slice2 []string
		want   []itertools.Tuple[int, string]
	}{
		{
			name:   "equal length slices",
			slice1: []int{1, 2, 3},
			slice2: []string{"one", "two", "three"},
			want: []itertools.Tuple[int, string]{
				{A: 1, B: "one"},
				{A: 2, B: "two"},
				{A: 3, B: "three"},
			},
		},
		{
			name:   "first slice longer",
			slice1: []int{1, 2, 3, 4},
			slice2: []string{"one", "two", "three"},
			want: []itertools.Tuple[int, string]{
				{A: 1, B: "one"},
				{A: 2, B: "two"},
				{A: 3, B: "three"},
			},
		},
		{
			name:   "second slice longer",
			slice1: []int{1, 2},
			slice2: []string{"one", "two", "three", "four"},
			want: []itertools.Tuple[int, string]{
				{A: 1, B: "one"},
				{A: 2, B: "two"},
			},
		},
		{
			name:   "empty first slice",
			slice1: []int{},
			slice2: []string{"one", "two", "three"},
			want:   []itertools.Tuple[int, string]{},
		},
		{
			name:   "empty second slice",
			slice1: []int{1, 2, 3},
			slice2: []string{},
			want:   []itertools.Tuple[int, string]{},
		},
		{
			name:   "both slices empty",
			slice1: []int{},
			slice2: []string{},
			want:   []itertools.Tuple[int, string]{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := itertools.Zip(tt.slice1, tt.slice2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Zip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GroupBy(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		slice := []string{"A", "B", "DEF", "GH", "IJK"}
		wanted := map[int][]string{
			1: {"A", "B"},
			2: {"GH"},
			3: {"DEF", "IJK"},
		}
		if got := itertools.GroupBy(slice, func(x string) int { return len(x) }); !reflect.DeepEqual(got, wanted) {
			t.Errorf("GroupBy() expected %v, got %v", wanted, got)
		}
	})

	tests := []struct {
		name    string
		slice   []int
		keyFunc func(int) int
		want    map[int][]int
	}{
		{
			name:  "group by parity",
			slice: []int{1, 2, 3, 4, 5},
			keyFunc: func(x int) int {
				return x % 2
			},
			want: map[int][]int{
				1: {1, 3, 5},
				0: {2, 4},
			},
		},
		{
			name:    "empty slice",
			slice:   []int{},
			keyFunc: func(x int) int { return 1 },
			want:    map[int][]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := itertools.GroupBy(tt.slice, tt.keyFunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("groupBy() = %v, want %v", got, tt.want)
			}
		})
	}
}
