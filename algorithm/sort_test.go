package algorithm_test

import (
	"github.com/shijl0925/go-toolkits/algorithm"
	"reflect"
	"testing"
)

func Test_BubbleSortV1(t *testing.T) {
	input := []int{4, 1, 3, 1, 5, 2}
	expected := []int{1, 1, 2, 3, 4, 5}
	algorithm.BubbleSortV1(input)

	if !reflect.DeepEqual(expected, input) {
		t.Errorf("Expected %q, got %q", expected, input)
	}
}

func Test_BubbleSortV2(t *testing.T) {
	input := []int{4, 1, 3, 1, 5, 2}
	expected := []int{1, 1, 2, 3, 4, 5}
	algorithm.BubbleSortV2(input)

	if !reflect.DeepEqual(expected, input) {
		t.Errorf("Expected %q, got %q", expected, input)
	}
}

func Test_SelectSort(t *testing.T) {
	input := []int{4, 1, 3, 1, 5, 2}
	expected := []int{1, 1, 2, 3, 4, 5}
	algorithm.SelectSort(input)

	if !reflect.DeepEqual(expected, input) {
		t.Errorf("Expected %q, got %q", expected, input)
	}
}

func Test_InsertSort(t *testing.T) {
	input := []int{4, 1, 3, 1, 5, 2}
	expected := []int{1, 1, 2, 3, 4, 5}
	algorithm.InsertSort(input)

	if !reflect.DeepEqual(expected, input) {
		t.Errorf("Expected %q, got %q", expected, input)
	}
}
