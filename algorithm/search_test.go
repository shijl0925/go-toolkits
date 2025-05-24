package algorithm_test

import (
	"fmt"
	"github.com/shijl0925/go-toolkits/algorithm"
	"testing"
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
