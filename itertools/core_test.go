package itertools_test

import (
	"reflect"
	"testing"

	toolkits "github.com/shijl0925/go-toolkits"
	"github.com/shijl0925/go-toolkits/itertools"
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

// TestFind 是 Find 函数的单元测试。
func TestFind(t *testing.T) {
	t.Run("Nil slice returns false", func(t *testing.T) {
		var s []int = nil
		result := itertools.Find(s, func(x int) bool { return x > 0 })
		if result {
			t.Errorf("Expected false, got true")
		}
	})

	t.Run("Empty slice returns false", func(t *testing.T) {
		s := []int{}
		result := itertools.Find(s, func(x int) bool { return x > 0 })
		if result {
			t.Errorf("Expected false, got true")
		}
	})

	t.Run("No element matches condition", func(t *testing.T) {
		s := []int{1, 2, 3}
		result := itertools.Find(s, func(x int) bool { return x > 5 })
		if result {
			t.Errorf("Expected false, got true")
		}
	})

	t.Run("First element matches", func(t *testing.T) {
		s := []int{1, 2, 3}
		result := itertools.Find(s, func(x int) bool { return x == 1 })
		if !result {
			t.Errorf("Expected true, got false")
		}
	})

	t.Run("Middle element matches", func(t *testing.T) {
		s := []int{1, 2, 3}
		result := itertools.Find(s, func(x int) bool { return x == 2 })
		if !result {
			t.Errorf("Expected true, got false")
		}
	})

	t.Run("Last element matches", func(t *testing.T) {
		s := []int{1, 2, 3}
		result := itertools.Find(s, func(x int) bool { return x == 3 })
		if !result {
			t.Errorf("Expected true, got false")
		}
	})

	t.Run("Multiple elements match", func(t *testing.T) {
		s := []string{"a", "ab", "abc"}
		result := itertools.Find(s, func(x string) bool { return len(x) >= 2 })
		if !result {
			t.Errorf("Expected true, got false")
		}
	})
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

// TestDropWhile tests the DropWhile function with various test cases.
func TestDropWhile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     []int
		predicate func(int) bool
		expected  []int
	}{
		{
			name:      "EmptySlice",
			input:     []int{},
			predicate: func(_ int) bool { return true },
			expected:  []int{},
		},
		{
			name:      "AllElementsMatch",
			input:     []int{1, 2, 3, 4},
			predicate: func(x int) bool { return x < 5 },
			expected:  []int{},
		},
		{
			name:      "SomeElementsMatch",
			input:     []int{1, 2, 3, 4, 5, 6},
			predicate: func(x int) bool { return x <= 3 },
			expected:  []int{4, 5, 6},
		},
		{
			name:      "NoElementsMatch",
			input:     []int{4, 5, 6},
			predicate: func(x int) bool { return x < 4 },
			expected:  []int{4, 5, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := itertools.DropWhile(tt.input, tt.predicate)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestDropWhile_GenericString tests DropWhile with a string slice to verify generic behavior.
func TestDropWhile_GenericString(t *testing.T) {
	input := []string{"a", "aa", "aaa", "b"}
	predicate := func(s string) bool { return len(s) < 3 }
	expected := []string{"aaa", "b"}

	result := itertools.DropWhile(input, predicate)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
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

func Test_CombineToMap(t *testing.T) {
	tests := []struct {
		name   string
		slice1 []int
		slice2 []string
		want   map[int]string
	}{
		{
			name:   "equal length slices",
			slice1: []int{1, 2, 3},
			slice2: []string{"one", "two", "three"},
			want: map[int]string{
				1: "one",
				2: "two",
				3: "three",
			},
		},
		{
			name:   "first slice longer",
			slice1: []int{1, 2, 3, 4},
			slice2: []string{"one", "two", "three"},
			want: map[int]string{
				1: "one",
				2: "two",
				3: "three",
			},
		},
		{
			name:   "second slice longer",
			slice1: []int{1, 2},
			slice2: []string{"one", "two", "three", "four"},
			want: map[int]string{
				1: "one",
				2: "two",
			},
		},
		{
			name:   "empty first slice",
			slice1: []int{},
			slice2: []string{"one", "two", "three"},
			want:   map[int]string{},
		},
		{
			name:   "empty second slice",
			slice1: []int{1, 2, 3},
			slice2: []string{},
			want:   map[int]string{},
		},
		{
			name:   "both slices empty",
			slice1: []int{},
			slice2: []string{},
			want:   map[int]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := itertools.CombineToMap(tt.slice1, tt.slice2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CombineToMap() = %v, want %v", got, tt.want)
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
			keyFunc: func(_ int) int { return 1 },
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

	t.Run("test2", func(t *testing.T) {
		type Location struct {
			address string
			date    string
		}
		slice := []Location{
			{
				address: "5412 N CLARK",
				date:    "07/01/2012",
			},
			{
				address: "5148 N CLARK",
				date:    "07/04/2012",
			},
			{
				address: "5800 E 58TH",
				date:    "07/02/2012",
			},
			{
				address: "2122 N CLARK",
				date:    "07/03/2012",
			},
			{
				address: "5645 N RAVENSWOOD",
				date:    "07/02/2012",
			},
			{
				address: "1060 W ADDISON",
				date:    "07/02/2012",
			},
			{
				address: "4801 N BROADWAY",
				date:    "07/01/2012",
			},
			{
				address: "1039 W GRANVILLE",
				date:    "07/04/2012",
			},
		}
		wanted := map[string][]Location{
			"07/01/2012": []Location{
				{
					address: "5412 N CLARK",
					date:    "07/01/2012",
				},
				{
					address: "4801 N BROADWAY",
					date:    "07/01/2012",
				},
			},
			"07/02/2012": []Location{
				{
					address: "5800 E 58TH",
					date:    "07/02/2012",
				},
				{
					address: "5645 N RAVENSWOOD",
					date:    "07/02/2012",
				},
				{
					address: "1060 W ADDISON",
					date:    "07/02/2012",
				},
			},
			"07/03/2012": []Location{
				{
					address: "2122 N CLARK",
					date:    "07/03/2012",
				},
			},
			"07/04/2012": []Location{
				{
					address: "5148 N CLARK",
					date:    "07/04/2012",
				},
				{
					address: "1039 W GRANVILLE",
					date:    "07/04/2012",
				},
			},
		}
		if got := itertools.GroupBy(slice, func(x Location) string { return x.date }); !reflect.DeepEqual(got, wanted) {
			t.Errorf("GroupBy() expected %v, got %v", wanted, got)
		}
	})
}

func Test_Range_Int(t *testing.T) {
	testCases := []struct {
		start, end, step int
		expected         []int
	}{
		{1, 5, 1, []int{1, 2, 3, 4}},
		{5, 1, -1, []int{5, 4, 3, 2}},
		{1, 5, 0, []int{}},
		{5, 1, 1, []int{}},
		{1, 5, -1, []int{}},
		{2, 2, 1, []int{}},
		{0, 10, 2, []int{0, 2, 4, 6, 8}},
		{-5, 5, 2, []int{-5, -3, -1, 1, 3}},
		{5, -5, -2, []int{5, 3, 1, -1, -3}},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			result := itertools.Range(tc.start, tc.end, tc.step)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Range(%d, %d, %d) = %v, want %v", tc.start, tc.end, tc.step, result, tc.expected)
			}
		})
	}
}

func Test_Range_Float(t *testing.T) {
	testCases := []struct {
		start, end, step float64
		expected         []float64
	}{
		{1.0, 5.0, 1.0, []float64{1.0, 2.0, 3.0, 4.0}},
		{5.0, 1.0, -1.0, []float64{5.0, 4.0, 3.0, 2.0}},
		{10.0, 0.0, -2.5, []float64{10.0, 7.5, 5.0, 2.5}},
		{0.0, 1.0, 0.2, []float64{0.0, 0.2, 0.4, 0.6, 0.8}},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			result := itertools.Range(tc.start, tc.end, tc.step)
			for i := range result {
				if !toolkits.EqualFloat64(result[i], tc.expected[i], 9) {
					t.Errorf("Range(%v, %v, %v) = %v; want %v", tc.start, tc.end, tc.step, result, tc.expected)
				}
			}
		})
	}
}
