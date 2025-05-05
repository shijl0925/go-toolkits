package slice_test

import (
	"github.com/shijl0925/go-toolkits/slice"
	"reflect"
	"testing"
)

func Test_SumSlice(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		if got := slice.Sum([]int{1, 2, 3, 4, 5}); got != 15 {
			t.Errorf("Sum() expected %v, got %v", 15, got)
		}
	})
	t.Run("test2", func(t *testing.T) {
		if got := slice.Sum([]float64{1.1, 2.0, 3.5}); got != 6.6 {
			t.Errorf("Sum() expected %v, got %v", 6.6, got)
		}
	})
}

func Test_MaxSlice(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		if got := slice.Max([]int{-1, -5, -3, -4, -2}); got != -1 {
			t.Errorf("Max() expected %v, got %v", -1, got)
		}
	})
	t.Run("test2", func(t *testing.T) {
		if got := slice.Max([]float64{1.1, 3.5, 2.0, 0.1}); got != 3.5 {
			t.Errorf("Max() expected %v, got %v", 3.5, got)
		}
	})
	t.Run("test3", func(t *testing.T) {
		if got := slice.Max([]int{-1, 5, 3, 4, 2}); got != 5 {
			t.Errorf("Max() expected %v, got %v", 5, got)
		}
	})
}

func Test_MinSlice(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		if got := slice.Min([]int{-1, 5, 3, 4, 2}); got != -1 {
			t.Errorf("Min() expected %v, got %v", -1, got)
		}
	})
	t.Run("test2", func(t *testing.T) {
		if got := slice.Min([]float64{1.1, 3.5, 2.0, 0.1}); got != 0.1 {
			t.Errorf("Min() expected %v, got %v", 0.1, got)
		}
	})
	t.Run("test3", func(t *testing.T) {
		if got := slice.Min([]int{-1, -5, -3, -4, -2}); got != -5 {
			t.Errorf("Min() expected %v, got %v", -5, got)
		}
	})
}

func Test_Insert(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		element  int
		index    int
		expected []int
	}{
		{"test1", []int{}, 1, 0, []int{1}},
		{"test2", []int{1, 2, 3, 4}, 2, 1, []int{1, 2, 2, 3, 4}},
		{"test3", []int{1, 2, 3, 4, 5, 6}, 7, 6, []int{1, 2, 3, 4, 5, 6, 7}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slice.Insert(tc.slice, tc.element, tc.index); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Map() expected %v, got %v", tc.expected, got)
			}
		})
	}
}
