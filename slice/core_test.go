package slice_test

import (
	"github.com/shijl0925/go-toolkits/slice"
	"testing"
)

func Test_MapSliceArray(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		if got := slice.Sum([]int{1, 2, 3, 4, 5}); got != 15 {
			t.Errorf("Reverse() expected %v, got %v", 15, got)
		}
	})
	t.Run("test2", func(t *testing.T) {
		if got := slice.Sum([]float64{1.1, 2.0, 3.5}); got != 6.6 {
			t.Errorf("Reverse() expected %v, got %v", 15, got)
		}
	})
}
