package algorithm_test

import (
	"fmt"
	"testing"

	"github.com/shijl0925/go-toolkits/algorithm"
)

func TestBinarySearch(t *testing.T) {
	sortedNumbers := []int{1, 3, 4, 5, 7, 8, 9, 10, 12}
	n1 := algorithm.BinarySearch(sortedNumbers, 6, 0, len(sortedNumbers)-1)
	fmt.Printf("result: %v\n", n1)
	if n1 != -1 {
		t.Errorf("Expected %q, got %q", -1, n1)
	}

	n2 := algorithm.BinarySearch(sortedNumbers, 3, 0, len(sortedNumbers)-1)
	fmt.Printf("result: %v\n", n2)
	if n2 != 1 {
		t.Errorf("Expected %q, got %q", 1, n2)
	}

	n3 := algorithm.BinarySearch(sortedNumbers, 8, 0, len(sortedNumbers)-1)
	fmt.Printf("result: %v\n", n3)
	if n3 != 5 {
		t.Errorf("Expected %q, got %q", -1, n3)
	}
}

func TestBinarySearch_OutOfRangeBounds(t *testing.T) {
	sortedNumbers := []int{1, 3, 4, 5, 7, 8, 9, 10, 12}

	tests := []struct {
		name   string
		target int
		low    int
		high   int
		want   int
	}{
		{
			name:   "high bound exceeds slice",
			target: 8,
			low:    0,
			high:   100,
			want:   5,
		},
		{
			name:   "low bound below zero",
			target: 3,
			low:    -5,
			high:   len(sortedNumbers) - 1,
			want:   1,
		},
		{
			name:   "both bounds outside slice but still searchable",
			target: 12,
			low:    -10,
			high:   999,
			want:   8,
		},
		{
			name:   "low bound stays past slice end",
			target: 12,
			low:    len(sortedNumbers),
			high:   len(sortedNumbers) + 10,
			want:   -1,
		},
		{
			name:   "high bound below zero",
			target: 1,
			low:    -10,
			high:   -1,
			want:   -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := algorithm.BinarySearch(sortedNumbers, tt.target, tt.low, tt.high)
			if got != tt.want {
				t.Errorf("BinarySearch() = %d; want %d", got, tt.want)
			}
		})
	}
}

func TestBinarySearch_EmptySliceWithOutOfRangeBounds(t *testing.T) {
	if got := algorithm.BinarySearch([]int{}, 1, -10, 100); got != -1 {
		t.Errorf("BinarySearch() = %d; want -1", got)
	}
}

func TestBinaryIterativeSearch(t *testing.T) {
	sortedNumbers := []int{1, 3, 4, 5, 7, 8, 9, 10, 12}
	n1 := algorithm.BinaryIterativeSearch(sortedNumbers, 6)
	fmt.Printf("result: %v\n", n1)
	if n1 != -1 {
		t.Errorf("Expected %q, got %q", -1, n1)
	}

	n2 := algorithm.BinaryIterativeSearch(sortedNumbers, 3)
	fmt.Printf("result: %v\n", n2)
	if n2 != 1 {
		t.Errorf("Expected %q, got %q", 1, n2)
	}

	n3 := algorithm.BinaryIterativeSearch(sortedNumbers, 8)
	fmt.Printf("result: %v\n", n3)
	if n3 != 5 {
		t.Errorf("Expected %q, got %q", -1, n3)
	}
}
