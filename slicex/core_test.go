package slicex_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/shijl0925/go-toolkits/slicex"
)

func Test_SumSlice(t *testing.T) {
	t.Run("Sum nil", func(t *testing.T) {
		got := slicex.Sum[int](nil)
		if got != 0 {
			t.Errorf("Sum() expected %v, got %v", 0, got)
		}
	})
	tests := []struct {
		name     string
		slice    any
		expected any
	}{
		{"TC01: Positive integers", []int{1, 2, 3, 4, 5}, 15},
		{"TC03: Float numbers", []float64{1.1, 2.0, 3.5}, 6.6},
		{"TC03: Negative and zero", []int{1, 0, -1}, 0},
		{"TC04: Empty slice", []int{}, 0},
		{"TC05: Unsigned integers", []uint{1, 2, 3, 4, 5}, uint(15)},
		{"TC06: Complex numbers", []complex128{1 + 2i, 3 + 4i}, complex(4, 6)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch input := tt.slice.(type) {
			case []int:
				got := slicex.Sum(input)
				if got != tt.expected {
					t.Errorf("Sum() expected %v, got %v", tt.expected, got)
				}
			case []float64:
				got := slicex.Sum(input)
				if got != tt.expected {
					t.Errorf("Sum() expected %v, got %v", tt.expected, got)
				}
			case []uint:
				got := slicex.Sum(input)
				if got != tt.expected {
					t.Errorf("Sum() expected %v, got %v", tt.expected, got)
				}
			case []complex128:
				got := slicex.Sum(input)
				if got != tt.expected {
					t.Errorf("Sum() expected %v, got %v", tt.expected, got)
				}
			default:
				t.Fatalf("Unsupported input type: %T", input)
			}
		})
	}
}

func TestAvg_Nil(t *testing.T) {
	got := slicex.Avg[float64](nil)
	if got != 0 {
		t.Errorf("Avg() expected %v, got %v", 0, got)
	}
}

// TestAvg_Int 测试 int 类型的平均值计算
func TestAvg_Int(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected float64
	}{
		{"empty slice", []int{}, 0},
		{"single element", []int{100}, 100},
		{"positive numbers", []int{1, 2, 3}, 2.0},
		{"negative and zero", []int{-1, 0, 1}, 0.0},
		{"all zeros", []int{0, 0, 0}, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slicex.Avg(tt.input)
			if got != tt.expected {
				t.Errorf("Avg(%v) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}

// TestAvg_Float64 测试 float64 类型的平均值计算
func TestAvg_Float64(t *testing.T) {
	tests := []struct {
		name     string
		input    []float64
		expected float64
	}{
		{"empty slice", []float64{}, 0},
		{"single element", []float64{2.5}, 2.5},
		{"multiple values", []float64{2.5, 3.5}, 3.0},
		{"negative and positive", []float64{-1.5, 1.5, 1.5}, 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slicex.Avg(tt.input)
			if got != tt.expected {
				t.Errorf("Avg(%v) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}

// TestAvg_Uint 测试 uint 类型的平均值计算
func TestAvg_Uint(t *testing.T) {
	tests := []struct {
		name     string
		input    []uint
		expected float64
	}{
		{"empty slice", []uint{}, 0},
		{"single element", []uint{10}, 10},
		{"multiple elements", []uint{1, 2, 3}, 2.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slicex.Avg(tt.input)
			if got != tt.expected {
				t.Errorf("Avg(%v) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func Test_MaxSlice(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []float64
		expected float64
	}{
		{"test1", []float64{-1, -5, -3, -4, -2}, -1},
		{"test2", []float64{1.1, 3.5, 2.0, 0.1}, 3.5},
		{"test3", []float64{-1, 5, 3, 4, 2}, 5},
		{"test4", []float64{}, 0.0},
		{"test5", nil, 0.0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slicex.Max(tc.slice); got != tc.expected {
				t.Errorf("Max() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

// TestMax_EmptySlice_ReturnsZeroValue tests the case when input slice is empty.
func TestMax_EmptySlice_ReturnsZeroValue(t *testing.T) {
	t.Run("TC01", func(t *testing.T) {
		got := slicex.Max([]int{})
		if got != 0 {
			t.Errorf("Max() expected %v, got %v", 0, got)
		}
	})
	t.Run("TC02", func(t *testing.T) {
		got := slicex.Max([]float64{})
		if got != 0.0 {
			t.Errorf("Max() expected %v, got %v", 0.0, got)
		}
	})
	t.Run("TC03", func(t *testing.T) {
		got := slicex.Max([]string{})
		if got != "" {
			t.Errorf("Max() expected %v, got %v", "", got)
		}
	})
	t.Run("TC04", func(t *testing.T) {
		got := slicex.Max[int](nil)
		if got != 0 {
			t.Errorf("Max() expected %v, got %v", 0, got)
		}
	})
}

// TestMax_SingleElement_ReturnsThatElement tests when there's only one element.
func TestMax_SingleElement_ReturnsThatElement(t *testing.T) {
	t.Run("TC01", func(t *testing.T) {
		got := slicex.Max([]int{99})
		if got != 99 {
			t.Errorf("Max() expected %v, got %v", 99, got)
		}
	})
	t.Run("TC02", func(t *testing.T) {
		got := slicex.Max([]string{"z"})
		if got != "z" {
			t.Errorf("Max() expected %v, got %v", "z", got)
		}
	})
	t.Run("TC03", func(t *testing.T) {
		got := slicex.Max([]float64{3.14})
		if got != 3.14 {
			t.Errorf("Max() expected %v, got %v", 3.14, got)
		}
	})
}

// TestMax_AllElementsSame_ReturnsSameValue tests when all elements are same.
func TestMax_AllElementsSame_ReturnsSameValue(t *testing.T) {
	t.Run("TC01", func(t *testing.T) {
		got := slicex.Max([]int{10, 10, 10})
		if got != 10 {
			t.Errorf("Max() expected %v, got %v", 10, got)
		}
	})
	t.Run("TC02", func(t *testing.T) {
		got := slicex.Max([]string{"a", "a", "a"})
		if got != "a" {
			t.Errorf("Max() expected %v, got %v", "a", got)
		}
	})
}

// TestMax_MaxAtBeginning_Middle_End tests different positions of max value.
func TestMax_MaxAtBeginning_Middle_End(t *testing.T) {
	t.Run("TC01", func(t *testing.T) {
		got := slicex.Max([]int{10, 5, 3})
		if got != 10 {
			t.Errorf("Max() expected %v, got %v", 10, got)
		}
	})
	t.Run("TC02", func(t *testing.T) {
		got := slicex.Max([]int{5, 10, 3})
		if got != 10 {
			t.Errorf("Max() expected %v, got %v", 10, got)
		}
	})
	t.Run("TC03", func(t *testing.T) {
		got := slicex.Max([]int{5, 3, 10})
		if got != 10 {
			t.Errorf("Max() expected %v, got %v", 10, got)
		}
	})
}

// TestMax_NegativeNumbers tests with negative numbers.
func TestMax_NegativeNumbers(t *testing.T) {
	t.Run("TC01", func(t *testing.T) {
		got := slicex.Max([]int{-5, -10, -1})
		if got != -1 {
			t.Errorf("Max() expected %v, got %v", -1, got)
		}
	})
}

// TestMax_StringLexicographicalOrder tests lexicographical order for strings.
func TestMax_StringLexicographicalOrder(t *testing.T) {
	t.Run("TC01", func(t *testing.T) {
		got := slicex.Max([]string{"a", "c", "b"})
		if got != "c" {
			t.Errorf("Max() expected %v, got %v", "c", got)
		}
	})
	t.Run("TC02", func(t *testing.T) {
		got := slicex.Max([]string{"apple", "zoo", "banana"})
		if got != "zoo" {
			t.Errorf("Max() expected %v, got %v", "zoo", got)
		}
	})
}

func Test_MinSlice(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []float64
		expected float64
	}{
		{"test1", []float64{-1, 5, 3, 4, 2}, -1},
		{"test2", []float64{1.1, 3.5, 2.0, 0.1}, 0.1},
		{"test3", []float64{-1, -5, -3, -4, -2}, -5},
		{"test4", []float64{}, 0.0},
		{"test5", nil, 0.0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slicex.Min(tc.slice); got != tc.expected {
				t.Errorf("Min() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

// TestMin_EmptySlice_ReturnsZeroValue tests the case when input slice is empty.
func TestMin_EmptySlice_ReturnsZeroValue(t *testing.T) {
	t.Run("TC01", func(t *testing.T) {
		got := slicex.Min([]int{})
		if got != 0 {
			t.Errorf("Min() expected %v, got %v", 0, got)
		}
	})
	t.Run("TC02", func(t *testing.T) {
		got := slicex.Min([]float64{})
		if got != 0.0 {
			t.Errorf("Min() expected %v, got %v", 0.0, got)
		}
	})
	t.Run("TC03", func(t *testing.T) {
		got := slicex.Min([]string{})
		if got != "" {
			t.Errorf("Min() expected %v, got %v", "", got)
		}
	})
	t.Run("TC04", func(t *testing.T) {
		got := slicex.Min[int](nil)
		if got != 0 {
			t.Errorf("Min() expected %v, got %v", 0, got)
		}
	})
}

// TestMin_SingleElement_ReturnsThatElement tests when there's only one element.
func TestMin_SingleElement_ReturnsThatElement(t *testing.T) {
	t.Run("TC01", func(t *testing.T) {
		got := slicex.Min([]int{99})
		if got != 99 {
			t.Errorf("Min() expected %v, got %v", 99, got)
		}
	})
	t.Run("TC02", func(t *testing.T) {
		got := slicex.Min([]string{"z"})
		if got != "z" {
			t.Errorf("Min() expected %v, got %v", "z", got)
		}
	})
	t.Run("TC03", func(t *testing.T) {
		got := slicex.Min([]float64{3.14})
		if got != 3.14 {
			t.Errorf("Min() expected %v, got %v", 3.14, got)
		}
	})
}

// TestMin_AllElementsSame_ReturnsSameValue tests when all elements are same.
func TestMin_AllElementsSame_ReturnsSameValue(t *testing.T) {
	t.Run("TC01", func(t *testing.T) {
		got := slicex.Min([]int{10, 10, 10})
		if got != 10 {
			t.Errorf("Min() expected %v, got %v", 10, got)
		}
	})
	t.Run("TC02", func(t *testing.T) {
		got := slicex.Min([]string{"a", "a", "a"})
		if got != "a" {
			t.Errorf("Min() expected %v, got %v", "a", got)
		}
	})
}

// TestMin_MinAtBeginning_Middle_End tests different positions of min value.
func TestMin_MinAtBeginning_Middle_End(t *testing.T) {
	t.Run("TC01", func(t *testing.T) {
		got := slicex.Min([]int{10, 5, 3})
		if got != 3 {
			t.Errorf("Min() expected %v, got %v", 3, got)
		}
	})
	t.Run("TC02", func(t *testing.T) {
		got := slicex.Min([]int{5, 10, 3})
		if got != 3 {
			t.Errorf("Min() expected %v, got %v", 3, got)
		}
	})
	t.Run("TC03", func(t *testing.T) {
		got := slicex.Min([]int{5, 3, 10})
		if got != 3 {
			t.Errorf("Min() expected %v, got %v", 3, got)
		}
	})
}

// TestMin_NegativeNumbers tests with negative numbers.
func TestMin_NegativeNumbers(t *testing.T) {
	t.Run("TC01", func(t *testing.T) {
		got := slicex.Min([]int{-5, -10, -1})
		if got != -10 {
			t.Errorf("Min() expected %v, got %v", -10, got)
		}
	})
}

// TestMin_StringLexicographicalOrder tests lexicographical order for strings.
func TestMin_StringLexicographicalOrder(t *testing.T) {
	t.Run("TC01", func(t *testing.T) {
		got := slicex.Min([]string{"a", "c", "b"})
		if got != "a" {
			t.Errorf("Min() expected %v, got %v", "a", got)
		}
	})
	t.Run("TC02", func(t *testing.T) {
		got := slicex.Min([]string{"apple", "zoo", "banana"})
		if got != "apple" {
			t.Errorf("Min() expected %v, got %v", "apple", got)
		}
	})
}

func Test_InsertAtV0(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		element  int
		index    int
		expected []int
		wantErr  bool
	}{
		{"test1", []int{}, 1, 0, []int{1}, false},
		{"test2", []int{1, 2, 3, 4}, 2, 1, []int{1, 2, 2, 3, 4}, false},
		{"test3", []int{1, 2, 3, 4, 5, 6}, 7, 6, []int{1, 2, 3, 4, 5, 6, 7}, false},
		{"test4", []int{1, 2, 3, 4}, 2, -1, []int{1, 2, 3, 4}, true},
		{"test5", []int{1, 2, 3, 4}, 2, 5, []int{1, 2, 3, 4}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := slicex.InsertAtV0(tc.slice, tc.element, tc.index)

			if tc.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if !reflect.DeepEqual(got, tc.expected) {
					t.Errorf("Insert() expected %v, got %v", tc.expected, got)
				}
			}
		})
	}
}

func Test_InsertAtV1(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		element  int
		index    int
		expected []int
		wantErr  bool
	}{
		{"test1", []int{}, 1, 0, []int{1}, false},
		{"test2", []int{1, 2, 3, 4}, 2, 1, []int{1, 2, 2, 3, 4}, false},
		{"test3", []int{1, 2, 3, 4, 5, 6}, 7, 6, []int{1, 2, 3, 4, 5, 6, 7}, false},
		{"test4", []int{1, 2, 3, 4}, 2, -1, []int{1, 2, 3, 4}, true},
		{"test5", []int{1, 2, 3, 4}, 2, 5, []int{1, 2, 3, 4}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := slicex.InsertAtV1(tc.slice, tc.element, tc.index)

			if tc.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if !reflect.DeepEqual(got, tc.expected) {
					t.Errorf("InsertAtV1() expected %v, got %v", tc.expected, got)
				}
			}
		})
	}
}

func Test_InsertAt(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		element  int
		index    int
		expected []int
		wantErr  bool
	}{
		{"test0", nil, 1, 0, []int{1}, false},
		{"test1", []int{}, 1, 0, []int{1}, false},
		{"test2", []int{1, 2, 3, 4}, 2, 1, []int{1, 2, 2, 3, 4}, false},
		{"test3", []int{1, 2, 3, 4, 5, 6}, 7, 6, []int{1, 2, 3, 4, 5, 6, 7}, false},
		{"test4", []int{1, 2, 3, 4}, 2, -1, []int{1, 2, 3, 4}, true},
		{"test5", []int{1, 2, 3, 4}, 2, 5, []int{1, 2, 3, 4}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := slicex.InsertAt(tc.slice, tc.element, tc.index)

			if tc.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if !reflect.DeepEqual(got, tc.expected) {
					t.Errorf("InsertAtV2() expected %v, got %v", tc.expected, got)
				}
			}
		})
	}
}

func Test_InsertMany(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		elements []int
		index    int
		expected []int
		wantErr  bool
	}{
		{"test1", []int{}, []int{1}, 0, []int{1}, false},
		{"test2", []int{1, 2, 3, 4}, []int{2, 3}, 1, []int{1, 2, 3, 2, 3, 4}, false},
		{"test3", []int{1, 2, 3, 4, 5, 6}, []int{7}, 6, []int{1, 2, 3, 4, 5, 6, 7}, false},
		{"test4", []int{1, 2, 3, 4}, []int{2, 3}, -1, []int{1, 2, 3, 4}, true},
		{"test5", []int{1, 2, 3, 4}, []int{2, 3}, 5, []int{1, 2, 3, 4}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := slicex.InsertMany(tc.slice, tc.elements, tc.index)

			if tc.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if !reflect.DeepEqual(got, tc.expected) {
					t.Errorf("InsertMany() expected %v, got %v", tc.expected, got)
				}
			}
		})
	}
}

func Test_AddMany(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		elements []int
		index    int
		expected []int
		wantErr  bool
	}{
		{"test1", []int{}, []int{1}, 0, []int{1}, false},
		{"test2", []int{1, 2, 3, 4}, []int{2, 3}, 1, []int{1, 2, 3, 2, 3, 4}, false},
		{"test3", []int{1, 2, 3, 4, 5, 6}, []int{7}, 6, []int{1, 2, 3, 4, 5, 6, 7}, false},
		{"test4", []int{1, 2, 3, 4}, []int{2, 3}, -1, []int{1, 2, 3, 4}, true},
		{"test5", []int{1, 2, 3, 4}, []int{2, 3}, 5, []int{1, 2, 3, 4}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := slicex.AddMany(tc.slice, tc.elements, tc.index)
			if tc.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if !reflect.DeepEqual(got, tc.expected) {
					t.Errorf("AddMany() expected %v, got %v", tc.expected, got)
				}
			}
		})
	}
}

func Test_Pop(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		expected int
		ok       bool
	}{
		{"test1", []int{1, 2, 3, 4}, 4, true},
		{"test2", []int{1, 2, 3, 4, 5}, 5, true},
		{"test3", []int{1, 2, 3, 4, 5, 6}, 6, true},
		{"test4", []int{}, 0, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			last, ok := slicex.Pop(tc.slice)
			if last != tc.expected || ok != tc.ok {
				t.Errorf("Pop() expected %v(%v), got %v(%v)", tc.expected, tc.ok, last, ok)
			}
		})
	}

	t.Run("nil pointer", func(t *testing.T) {
		var s []int
		val, ok := slicex.Pop(s)
		var zero int
		if zero != val || ok {
			t.Errorf("Pop() expected %v(%v), got %v(%v)", zero, false, val, ok)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		var s []string

		val, ok := slicex.Pop(s)
		var zero string
		if zero != val || ok {
			t.Errorf("Pop() expected %v(%v), got %v(%v)", zero, false, val, ok)
		}
	})

	t.Run("single element", func(t *testing.T) {
		s := []int{1}
		val, ok := slicex.Pop(s)
		if val != 1 || !ok {
			t.Errorf("Pop() expected %v(%v), got %v(%v)", val, true, val, ok)
		}
		if !reflect.DeepEqual(s, []int{1}) {
			t.Errorf("Pop() expected %v, got %v", []int{}, s)
		}
	})

	t.Run("multiple elements", func(t *testing.T) {
		s := []string{"a", "b"}
		val, ok := slicex.Pop(s)
		if val != "b" || !ok {
			t.Errorf("Pop() expected %v(%v), got %v(%v)", val, true, val, ok)
		}
		if !reflect.DeepEqual(s, []string{"a", "b"}) {
			t.Errorf("Pop() expected %v, got %v", []string{"a"}, s)
		}
	})

	t.Run("generic type struct", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}
		s := []User{
			{Name: "Alice", Age: 25},
			{Name: "Bob", Age: 30},
			{Name: "Charlie", Age: 35},
		}
		val, ok := slicex.Pop(s)
		if !reflect.DeepEqual(val, User{Name: "Charlie", Age: 35}) || !ok {
			t.Errorf("Pop() expected %v(%v), got %v(%v)", val, true, val, ok)
		}
		if !reflect.DeepEqual(s, []User{{Name: "Alice", Age: 25}, {Name: "Bob", Age: 30}, {Name: "Charlie", Age: 35}}) {
			t.Errorf("Pop() expected %v, got %v", []User{{Name: "Alice", Age: 25}, {Name: "Bob", Age: 30}, {Name: "Charlie", Age: 35}}, s)
		}
	})
}

func Test_Shift(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		expected int
		ok       bool
	}{
		{"test1", []int{1, 2, 3, 4}, 1, true},
		{"test2", []int{1, 2, 3, 4, 5}, 1, true},
		{"test3", []int{1, 2, 3, 4, 5, 6}, 1, true},
		{"test4", []int{}, 0, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			val, ok := slicex.Shift(tc.slice)
			if val != tc.expected || ok != tc.ok {
				t.Errorf("Shift() expected %v(%v), got %v(%v)", tc.expected, tc.ok, val, ok)
			}
		})
	}

	t.Run("nil pointer", func(t *testing.T) {
		var s []int
		val, ok := slicex.Shift(s)
		var zero int
		if zero != val || ok {
			t.Errorf("Shift() expected %v(%v), got %v(%v)", zero, false, val, ok)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		var s []string

		val, ok := slicex.Shift(s)
		var zero string
		if zero != val || ok {
			t.Errorf("Shift() expected %v(%v), got %v(%v)", zero, false, val, ok)
		}
	})

	t.Run("single element", func(t *testing.T) {
		s := []int{1}
		val, ok := slicex.Shift(s)
		if val != 1 || !ok {
			t.Errorf("Shift() expected %v(%v), got %v(%v)", val, true, val, ok)
		}
		if !reflect.DeepEqual(s, []int{1}) {
			t.Errorf("Shift() expected %v, got %v", []int{1}, s)
		}
	})

	t.Run("multiple elements", func(t *testing.T) {
		s := []string{"a", "b"}
		val, ok := slicex.Shift(s)
		if val != "a" || !ok {
			t.Errorf("Shift() expected %v(%v), got %v(%v)", val, true, val, ok)
		}
		if !reflect.DeepEqual(s, []string{"a", "b"}) {
			t.Errorf("Shift() expected %v, got %v", []string{"a", "b"}, s)
		}
	})

	t.Run("generic type struct", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}
		s := []User{
			{Name: "Alice", Age: 25},
			{Name: "Bob", Age: 30},
			{Name: "Charlie", Age: 35},
		}
		val, ok := slicex.Shift(s)
		if !reflect.DeepEqual(val, User{Name: "Alice", Age: 25}) || !ok {
			t.Errorf("Shift() expected %v(%v), got %v(%v)", val, true, val, ok)
		}
		if !reflect.DeepEqual(s, []User{{Name: "Alice", Age: 25}, {Name: "Bob", Age: 30}, {Name: "Charlie", Age: 35}}) {
			t.Errorf("Shift() expected %v, got %v", []User{{Name: "Alice", Age: 25}, {Name: "Bob", Age: 30}, {Name: "Charlie", Age: 35}}, s)
		}
	})
}

func Test_Drop(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		n        int
		expected []int
		ok       bool
	}{
		{"test1", []int{1, 2, 3, 4}, 0, []int{1, 2, 3, 4}, true},
		{"test2", []int{1, 2, 3, 4, 5}, 2, []int{1, 2, 3}, true},
		{"test3", []int{1, 2, 3, 4, 5, 6}, 6, []int{}, true},
		{"test4", []int{1, 2, 3, 4}, 5, []int{}, true},
		{"test5", []int{1, 2, 3, 4}, -1, []int{1, 2, 3, 4}, false},
		{"test6", []int{}, 1, []int{}, true},
		{"test7", []int{}, 0, []int{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := slicex.Drop(tt.slice, tt.n)
			if !reflect.DeepEqual(result, tt.expected) || ok != tt.ok {
				t.Errorf("Drop() expected %v(%v), but got %v(%v)", tt.expected, tt.ok, result, ok)
			}
		})
	}
}

func TestDropLeft(t *testing.T) {
	tests := []struct {
		name     string
		s        []int
		n        int
		expected []int
		ok       bool
	}{
		{
			name:     "Normal case, n=2",
			s:        []int{1, 2, 3, 4},
			n:        2,
			expected: []int{3, 4},
			ok:       true,
		},
		{
			name:     "n=0, no removal",
			s:        []int{1, 2, 3, 4},
			n:        0,
			expected: []int{1, 2, 3, 4},
			ok:       true,
		},
		{
			name:     "n equals length of slice",
			s:        []int{1, 2, 3, 4},
			n:        4,
			expected: []int{},
			ok:       true,
		},
		{
			name:     "n negative",
			s:        []int{1, 2, 3, 4},
			n:        -1,
			expected: []int{1, 2, 3, 4},
			ok:       false,
		},
		{
			name:     "n greater than length",
			s:        []int{1, 2, 3, 4},
			n:        5,
			expected: []int{},
			ok:       true,
		},
		{
			name:     "Empty slice, n=0",
			s:        []int{},
			n:        0,
			expected: []int{},
			ok:       true,
		},
		{
			name:     "Empty slice, n=1",
			s:        []int{},
			n:        1,
			expected: []int{},
			ok:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := slicex.DropLeft(tt.s, tt.n)
			if !reflect.DeepEqual(result, tt.expected) || ok != tt.ok {
				t.Errorf("DropLeft() expected %v(%v), but got %v(%v)", tt.expected, tt.ok, result, ok)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		element  int
		expected []int
	}{
		{
			name:     "nil",
			input:    nil,
			element:  5,
			expected: []int{},
		},
		{
			name:     "Empty slice",
			input:    []int{},
			element:  5,
			expected: []int{},
		},
		{
			name:     "Int slice with one match",
			input:    []int{1, 2, 3, 4, 5},
			element:  3,
			expected: []int{1, 2, 4, 5},
		},
		{
			name:     "All elements match",
			input:    []int{1, 1, 1},
			element:  1,
			expected: []int{},
		},
		{
			name:     "Multiple matches in middle",
			input:    []int{2, 2, 3, 2},
			element:  2,
			expected: []int{3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := slicex.Remove(tt.input, tt.element)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Remove() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRemove_String(t *testing.T) {
	input := []string{"a", "b", "c", "b"}
	result := slicex.Remove(input, "b")
	expected := []string{"a", "c"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestRemove_Struct(t *testing.T) {
	type S struct{ X int }
	input := []S{{1}, {2}, {1}}
	result := slicex.Remove(input, S{1})
	expected := []S{{2}}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func Test_Delete(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		index    int
		expected []int
		ok       bool
	}{
		{"test1", []int{1, 2, 3, 4}, 1, []int{1, 3, 4}, true},
		{"test2", []int{1, 2, 3, 4}, 3, []int{1, 2, 3}, true},
		{"test4", []int{1, 2, 3, 4, 5, 6}, 2, []int{1, 2, 4, 5, 6}, true},
		{"test5", []int{}, 0, []int{}, false},
		{"test6", []int{1, 2, 3, 4}, 4, []int{1, 2, 3, 4}, false},
		{"test2", []int{1, 2, 3, 4}, 0, []int{2, 3, 4}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, ok := slicex.Delete(tc.slice, tc.index); !reflect.DeepEqual(got, tc.expected) || ok != tc.ok {
				t.Errorf("Delete() expected %v, %v, got %v, %v", tc.expected, tc.ok, got, ok)
			}
		})
	}
}

func Test_DeleteAt(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		index    int
		expected []int
		ok       bool
	}{
		{"test1", []int{1, 2, 3, 4}, 1, []int{1, 3, 4}, true},
		{"test2", []int{1, 2, 3, 4}, 3, []int{1, 2, 3}, true},
		{"test4", []int{1, 2, 3, 4, 5, 6}, 2, []int{1, 2, 4, 5, 6}, true},
		{"test5", []int{}, 0, []int{}, false},
		{"test6", nil, 0, nil, false},
		{"test7", []int{1, 2, 3, 4}, 4, []int{1, 2, 3, 4}, false},
		{"test8", []int{1, 2, 3, 4}, 0, []int{2, 3, 4}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, ok := slicex.DeleteAt(tc.slice, tc.index); !reflect.DeepEqual(got, tc.expected) || ok != tc.ok {
				t.Errorf("DeleteAt() expected %v, %v, got %v, %v", tc.expected, tc.ok, got, ok)
			}
		})
	}
}

func Test_DeleteMany(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		start    int
		end      int
		expected []int
		ok       bool
	}{
		{"test1", []int{1, 2, 3, 4}, 1, 2, []int{1, 3, 4}, true},
		{"test2", []int{1, 2, 3, 4, 5}, 2, 4, []int{1, 2, 5}, true},
		{"test3", []int{1, 2, 3, 4, 5, 6}, 4, 5, []int{1, 2, 3, 4, 6}, true},
		{"test4", []int{}, 1, 2, []int{}, true},
		{"test5", []int{1, 2, 3, 4}, -1, 2, []int{3, 4}, true},
		{"test6", []int{1, 2, 3, 4}, 1, 4, []int{1}, true},
		{"test6", []int{1, 2, 3, 4}, 1, 5, []int{1}, true},
		{"test7", []int{1, 2, 3, 4}, -1, 5, []int{}, true},
		{"test8", []int{1, 2, 3, 4}, 3, 2, []int{1, 2, 3, 4}, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, ok := slicex.DeleteMany(tc.slice, tc.start, tc.end); !reflect.DeepEqual(got, tc.expected) || ok != tc.ok {
				t.Errorf("DeleteMany() expected %v, %v, got %v, %v", tc.expected, tc.ok, got, ok)
			}
		})
	}
}

func Test_Append(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		element  int
		expected []int
	}{
		{"test1", []int{}, 1, []int{1}},
		{"test2", []int{1, 2, 3, 4}, 5, []int{1, 2, 3, 4, 5}},
		{"test3", []int{1, 2, 3, 4, 5, 6}, 7, []int{1, 2, 3, 4, 5, 6, 7}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slicex.Append(tc.slice, tc.element); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Append() expected %v, got %v", tc.expected, got)
			}
		})
	}

	slice1 := []int{1, 2, 3, 4}
	elements := []int{5, 6, 7}
	want := []int{1, 2, 3, 4, 5, 6, 7}
	t.Run("test0", func(t *testing.T) {
		if got := slicex.Append(slice1, elements...); !reflect.DeepEqual(got, want) {
			t.Errorf("Append() expected %v, got %v", want, got)
		}
	})
}

func Test_Extend(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		elements []int
		expected []int
	}{
		{"test1", []int{}, []int{5, 6, 7}, []int{5, 6, 7}},
		{"test2", []int{1, 2, 3, 4}, []int{5, 6, 7}, []int{1, 2, 3, 4, 5, 6, 7}},
		{"test3", []int{1, 2, 3, 4, 5, 6}, []int{7}, []int{1, 2, 3, 4, 5, 6, 7}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slicex.Extend(tc.slice, tc.elements); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Extend() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func Test_Unique(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		expected []int
	}{
		{"test0", nil, []int{}},
		{"test1", []int{}, []int{}},
		{"test2", []int{1, 2, 3, 4}, []int{1, 2, 3, 4}},
		{"test3", []int{1, 2, 2, 3, 4, 4, 5}, []int{1, 2, 3, 4, 5}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slicex.Unique(tc.slice); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Unique() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func Test_ReverseSlice(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		if got := slicex.Reverse([]int{1, 2, 3, 4, 5}); !reflect.DeepEqual(got, []int{5, 4, 3, 2, 1}) {
			t.Errorf("Reverse() expected %v, got %v", []int{5, 4, 3, 2, 1}, got)
		}
	})
	t.Run("test2", func(t *testing.T) {
		if got := slicex.Reverse([]string{"one", "two", "three"}); !reflect.DeepEqual(got, []string{"three", "two", "one"}) {
			t.Errorf("Reverse() expected %v, got %v", []string{"three", "two", "one"}, got)
		}
	})
	t.Run("test3", func(t *testing.T) {
		if got := slicex.Reverse([]byte("Google")); string(got) != "elgooG" {
			t.Errorf("Reverse() expected %v, got %v", "elgooG", string(got))
		}
	})
}

func Test_RepeatSlice(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		if got := slicex.Repeat(1, 5); !reflect.DeepEqual(got, []int{1, 1, 1, 1, 1}) {
			t.Errorf("Repeat() expected %v, got %v", []int{1, 1, 1, 1, 1}, got)
		}
	})
	t.Run("test2", func(t *testing.T) {
		wanted := [][]int{[]int{1, 2, 3}, []int{1, 2, 3}, []int{1, 2, 3}, []int{1, 2, 3}, []int{1, 2, 3}}
		if got := slicex.Repeat([]int{1, 2, 3}, 5); !reflect.DeepEqual(got, wanted) {
			t.Errorf("Repeat() expected %v, got %v", wanted, got)
		}
	})
	t.Run("test3", func(t *testing.T) {
		wanted := []string{"one", "one", "one"}
		if got := slicex.Repeat("one", 3); !reflect.DeepEqual(got, wanted) {
			t.Errorf("Repeat() expected %v, got %v", wanted, got)
		}
	})
}

func Test_Product(t *testing.T) {
	tests := []struct {
		name   string
		slices [][]string
		want   [][]string
	}{
		{
			name:   "test1",
			slices: [][]string{[]string{"A", "B", "C"}, []string{"one", "two", "three"}},
			want: [][]string{
				{"A", "one"},
				{"A", "two"},
				{"A", "three"},
				{"B", "one"},
				{"B", "two"},
				{"B", "three"},
				{"C", "one"},
				{"C", "two"},
				{"C", "three"},
			},
		},
		{
			name:   "test2",
			slices: [][]string{[]string{"A", "B"}, []string{"C", "D"}, []string{"E", "F"}},
			want: [][]string{
				{"A", "C", "E"},
				{"A", "C", "F"},
				{"A", "D", "E"},
				{"A", "D", "F"},
				{"B", "C", "E"},
				{"B", "C", "F"},
				{"B", "D", "E"},
				{"B", "D", "F"},
			},
		},
		{
			name:   "test3",
			slices: [][]string{[]string{"A", "B"}, []string{"C", "D", "E", "F"}},
			want: [][]string{
				{"A", "C"},
				{"A", "D"},
				{"A", "E"},
				{"A", "F"},
				{"B", "C"},
				{"B", "D"},
				{"B", "E"},
				{"B", "F"},
			},
		},
		{
			name:   "test4",
			slices: [][]string{[]string{"A", "B", "C", "D", "E"}, []string{"F", "G"}},
			want: [][]string{
				{"A", "F"},
				{"A", "G"},
				{"B", "F"},
				{"B", "G"},
				{"C", "F"},
				{"C", "G"},
				{"D", "F"},
				{"D", "G"},
				{"E", "F"},
				{"E", "G"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slicex.Product(tt.slices...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Product() = %v, want %v", got, tt.want)
			}
		})
	}
	t.Run("test5", func(t *testing.T) {
		slices := [][]int{[]int{1, 2, 3}}
		wanted := [][]int{{1}, {2}, {3}}
		if got := slicex.Product(slices...); !reflect.DeepEqual(got, wanted) {
			t.Errorf("Product() expected %v, got %v", wanted, got)
		}
	})
	t.Run("test6", func(t *testing.T) {
		var slices [][]int
		var wanted [][]int
		if got := slicex.Product(slices...); !reflect.DeepEqual(got, wanted) {
			t.Errorf("Product() expected %v, got %v", wanted, got)
		}
	})
}

func Test_Difference(t *testing.T) {
	var tests = []struct {
		name     string
		slice1   []int
		slice2   []int
		expected []int
	}{
		{"test1", []int{1, 2, 3, 4}, []int{3, 5, 4, 6, 7}, []int{1, 2}},
		{"test2", []int{1, 3, 4, 2}, []int{}, []int{1, 3, 4, 2}},
		{"test3", []int{1, 3, 2, 4}, []int{5, 6, 7, 8}, []int{1, 3, 2, 4}},
		{"test4", []int{}, []int{3, 4, 5, 6, 7}, []int{}},
		{"test5", []int{}, []int{}, []int{}},
		{"test6", []int{1, 2, 3, 4}, []int{1, 2, 3, 4}, []int{}},
		{"test7", nil, []int{3, 4, 5, 6, 7}, []int{}},
		{"test8", []int{1, 2, 3, 4}, nil, []int{1, 2, 3, 4}},
		{"test9", nil, nil, []int{}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slicex.Difference(tc.slice1, tc.slice2); !slicex.EqualUnordered(got, tc.expected) {
				t.Errorf("Difference() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func Test_Intersect(t *testing.T) {
	var tests = []struct {
		name     string
		slice1   []int
		slice2   []int
		expected []int
	}{
		{"test1", []int{1, 2, 3, 4}, []int{3, 5, 4, 6, 7}, []int{3, 4}},
		{"test2", []int{1, 3, 4, 2}, []int{}, []int{}},
		{"test3", []int{1, 3, 2, 4}, []int{5, 6, 7, 8}, []int{}},
		{"test4", []int{}, []int{3, 4, 5, 6, 7}, []int{}},
		{"test5", []int{}, []int{}, []int{}},
		{"test6", []int{1, 2, 3, 4}, []int{1, 2, 3, 4}, []int{1, 2, 3, 4}},
		{"test7", nil, []int{3, 4, 5, 6, 7}, []int{}},
		{"test8", []int{1, 2, 3, 4}, nil, []int{}},
		{"test9", nil, nil, []int{}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slicex.Intersect(tc.slice1, tc.slice2); !slicex.EqualUnordered(got, tc.expected) {
				t.Errorf("Intersect() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func Test_Union(t *testing.T) {
	var tests = []struct {
		name     string
		slice1   []int
		slice2   []int
		expected []int
	}{
		{"test1", []int{1, 2, 3, 4}, []int{3, 5, 4, 6, 7}, []int{1, 2, 3, 4, 5, 6, 7}},
		{"test2", []int{1, 3, 4, 2}, []int{}, []int{1, 3, 4, 2}},
		{"test3", []int{1, 3, 2, 2, 4}, []int{5, 6, 7, 8}, []int{1, 3, 2, 4, 5, 6, 7, 8}},
		{"test4", []int{1, 3, 2, 2, 4}, []int{5, 6, 7, 8, 5, 7}, []int{1, 3, 2, 4, 5, 6, 7, 8}},
		{"test5", []int{}, []int{3, 4, 5, 6, 7}, []int{3, 4, 5, 6, 7}},
		{"test6", []int{}, []int{}, []int{}},
		{"test7", nil, []int{3, 4, 5, 6, 7}, []int{3, 4, 5, 6, 7}},
		{"test8", []int{1, 2, 3, 4}, nil, []int{1, 2, 3, 4}},
		{"test9", nil, nil, []int{}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slicex.Union(tc.slice1, tc.slice2); !slicex.EqualUnordered(got, tc.expected) {
				t.Errorf("Union() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func Test_Index(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		e        int
		expected int
	}{
		{"test1", []int{1, 2, 3, 4}, 1, 0},
		{"test2", []int{1, 2, 3, 4}, -3, -1},
		{"test3", []int{1, 2, 3, 4, 4}, 4, 3},
		{"test4", []int{1, 2, 3, 4, 5, 4, 6}, 4, 3},
		{"test5", []int{}, 4, -1},
		{"test6", nil, 7, -1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slicex.Index(tc.slice, tc.e); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Index() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func Test_LastIndex(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		e        int
		expected int
	}{
		{"test1", []int{1, 2, 3, 4}, 2, 1},
		{"test2", []int{1, 2, 3, 4}, -3, -1},
		{"test3", []int{1, 2, 3, 4, 4}, 4, 4},
		{"test4", []int{1, 2, 3, 4, 5, 4, 6}, 4, 5},
		{"test5", []int{}, 4, -1},
		{"test6", nil, 7, -1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slicex.LastIndex(tc.slice, tc.e); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("LastIndex() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func Test_IndexAll(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		e        int
		expected []int
	}{
		{"test1", []int{1, 2, 3, 4}, 2, []int{1}},
		{"test2", []int{1, 2, 3, 4}, -3, []int{}},
		{"test3", []int{1, 2, 3, 4, 4}, 4, []int{3, 4}},
		{"test4", []int{1, 2, 3, 4, 5, 4, 6}, 4, []int{3, 5}},
		{"test5", []int{}, 4, []int{}},
		{"test6", nil, 7, []int{}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slicex.IndexAll(tc.slice, tc.e); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("IndexAll() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

func (t Person) String() string {
	return fmt.Sprintf("%s %s, %d", t.FirstName, t.LastName, t.Age)
}

func Test_JoinSlice(t *testing.T) {
	var s0 []int
	w0 := ""
	t.Run("test0", func(t *testing.T) {
		if got := slicex.JoinSlice(",", s0); !reflect.DeepEqual(got, w0) {
			t.Errorf("JoinSlice() expected %v, got %v", w0, got)
		}
	})

	s1 := []int{1, 2, 3, 4, 5}
	w1 := "1,2,3,4,5"
	t.Run("test1", func(t *testing.T) {
		if got := slicex.JoinSlice(",", s1); !reflect.DeepEqual(got, w1) {
			t.Errorf("JoinSlice() expected %v, got %v", w1, got)
		}
	})

	s2 := []string{"a", "b", "c", "d"}
	w2 := "a,b,c,d"
	t.Run("test2", func(t *testing.T) {
		if got := slicex.JoinSlice(",", s2); !reflect.DeepEqual(got, w2) {
			t.Errorf("JoinSlice() expected %v, got %v", w2, got)
		}
	})

	s3 := []Person{
		{"Alice", "Smith", 18},
		{"Bob", "Johnson", 20},
		{"Charlie", "Brown", 17},
	}
	w3 := "Alice Smith, 18;Bob Johnson, 20;Charlie Brown, 17"
	t.Run("test2", func(t *testing.T) {
		if got := slicex.JoinSlice(";", s3); !reflect.DeepEqual(got, w3) {
			t.Errorf("JoinSlice() expected %v, got %v", w3, got)
		}
	})
}

func TestChunk(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		size   int
		expect [][]int
	}{
		{
			name:   "nil slice and zero size",
			input:  nil,
			size:   0,
			expect: [][]int{},
		},
		{
			name:   "empty slice",
			input:  []int{},
			size:   5,
			expect: [][]int{},
		},
		{
			name:   "positive case with remainder",
			input:  []int{1, 2, 3, 4, 5},
			size:   2,
			expect: [][]int{{1, 2}, {3, 4}, {5}},
		},
		{
			name:   "exact division",
			input:  []int{1, 2, 3, 4},
			size:   4,
			expect: [][]int{{1, 2, 3, 4}},
		},
		{
			name:   "size larger than slice length",
			input:  []int{1, 2},
			size:   5,
			expect: [][]int{{1, 2}},
		},
		{
			name:   "zero size with non-empty slice",
			input:  []int{1, 2, 3},
			size:   0,
			expect: [][]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slicex.Chunk(tt.input, tt.size)
			if !reflect.DeepEqual(got, tt.expect) {
				t.Errorf("Chunk(%v, %d) = %v; want %v", tt.input, tt.size, got, tt.expect)
			}
		})
	}
}

// 测试用例结构体
type testCompactCase[T comparable] struct {
	name     string
	input    []T
	expected []T
}

func runCompactTests[T comparable](t *testing.T, tests []testCompactCase[T]) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := slicex.Compact(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Compact(%v) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCompact(t *testing.T) {
	// int 类型测试
	runCompactTests[int](t, []testCompactCase[int]{
		{"All zeros", []int{0, 0, 0}, []int{}},
		{"Mixed with zeros", []int{0, 1, 0, 2, 0}, []int{1, 2}},
		{"No zeros", []int{1, 2, 3}, []int{1, 2, 3}},
		{"Empty slice", []int{}, []int{}},
		{"Nil slice", nil, nil},
	})

	// string 类型测试
	runCompactTests[string](t, []testCompactCase[string]{
		{"All empty strings", []string{"", "", ""}, []string{}},
		{"Mixed empty strings", []string{"", "a", "", "b"}, []string{"a", "b"}},
		{"No empty strings", []string{"x", "y"}, []string{"x", "y"}},
	})

	// bool 类型测试
	runCompactTests[bool](t, []testCompactCase[bool]{
		{"All false", []bool{false, false}, []bool{}},
		{"Mixed", []bool{false, true, false}, []bool{true}},
	})

	// 指针类型测试
	runCompactTests[*int](t, []testCompactCase[*int]{
		{"Nil pointers", []*int{nil, nil}, []*int{}},
		{"Mixed pointers", []*int{nil, new(int), nil}, []*int{new(int)}},
	})

	// 结构体类型测试
	type S struct{}
	runCompactTests[S](t, []testCompactCase[S]{
		{"All zero structs", []S{{}, {}, {}}, []S{}},
	})

	// interface{} 类型测试
	runCompactTests[any](t, []testCompactCase[any]{
		{"Mixed interface", []any{nil, "hello", 42}, []any{"hello", 42}},
	})
}

// TestConcat_Int 测试整数切片的拼接
func TestConcat_Int(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]int
		expected []int
	}{
		{
			name:     "TC01 - 多个非空切片",
			input:    [][]int{{1, 2}, {3}},
			expected: []int{1, 2, 3},
		},
		{
			name:     "TC02 - 零个切片",
			input:    [][]int{},
			expected: []int{},
		},
		{
			name:     "TC03 - 包含空切片",
			input:    [][]int{{}, {1}},
			expected: []int{1},
		},
		{
			name:     "TC04 - 空切片",
			input:    [][]int{},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := slicex.Concat(tt.input...)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestConcat_String 测试字符串切片的拼接
func TestConcat_String(t *testing.T) {
	input := [][]string{{"a"}, {"b", "c"}}
	expected := []string{"a", "b", "c"}

	result := slicex.Concat(input...)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// TestConcat_Struct 测试结构体切片的拼接
func TestConcat_Struct(t *testing.T) {
	type S struct{}
	input := [][]S{{{}, {}}, {{}}}
	expected := []S{{}, {}, {}}

	result := slicex.Concat(input...)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// TestEqualUnordered 是 EqualUnordered 的单元测试
func TestEqualUnordered(t *testing.T) {
	tests := []struct {
		name string
		s1   []int
		s2   []int
		want bool
	}{
		{
			name: "TC01 - 整型切片顺序不同",
			s1:   []int{1, 2, 3},
			s2:   []int{3, 2, 1},
			want: true,
		},
		{
			name: "TC02 - 存在一个不同元素",
			s1:   []int{1, 2, 3},
			s2:   []int{1, 2, 4},
			want: false,
		},
		{
			name: "TC03 - 包含重复元素",
			s1:   []int{1, 2, 2},
			s2:   []int{2, 1, 2},
			want: true,
		},
		{
			name: "TC04 - 长度不同",
			s1:   []int{1, 2},
			s2:   []int{1, 2, 3},
			want: false,
		},
		{
			name: "TC05 - 空切片比较",
			s1:   []int{},
			s2:   []int{},
			want: true,
		},
		{
			name: "TC06 - 数量不一致",
			s1:   []int{1, 2, 3},
			s2:   []int{3, 2, 2},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slicex.EqualUnordered(tt.s1, tt.s2); got != tt.want {
				t.Errorf("EqualUnordered() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCounter(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected any
	}{
		{
			name:     "Empty slice",
			input:    []int{},
			expected: map[int]int{},
		},
		{
			name:     "Int slice with duplicates",
			input:    []int{1, 2, 2, 3},
			expected: map[int]int{1: 1, 2: 2, 3: 1},
		},
		{
			name:     "String slice with duplicates",
			input:    []string{"a", "b", "a"},
			expected: map[string]int{"a": 2, "b": 1},
		},
		{
			name:     "Float64 slice with duplicates",
			input:    []float64{1.1, 2.2, 1.1},
			expected: map[float64]int{1.1: 2, 2.2: 1},
		},
		{
			name:     "Struct slice",
			input:    []struct{}{{}, {}, {}},
			expected: map[struct{}]int{{}: 3},
		},
		{
			name:     "Interface slice",
			input:    []any{1, "a", 1, "a"},
			expected: map[any]int{1: 2, "a": 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch in := tt.input.(type) {
			case []int:
				got := slicex.Counter(in)
				if !reflect.DeepEqual(got, tt.expected) {
					t.Errorf("Counter() = %v, want %v", got, tt.expected)
				}
			case []string:
				got := slicex.Counter(in)
				if !reflect.DeepEqual(got, tt.expected) {
					t.Errorf("Counter() = %v, want %v", got, tt.expected)
				}
			case []float64:
				got := slicex.Counter(in)
				if !reflect.DeepEqual(got, tt.expected) {
					t.Errorf("Counter() = %v, want %v", got, tt.expected)
				}
			case []struct{}:
				got := slicex.Counter(in)
				if !reflect.DeepEqual(got, tt.expected) {
					t.Errorf("Counter() = %v, want %v", got, tt.expected)
				}
			case []any:
				got := slicex.Counter(in)
				if !reflect.DeepEqual(got, tt.expected) {
					t.Errorf("Counter() = %v, want %v", got, tt.expected)
				}
			default:
				t.Fatalf("Unsupported input type: %T", in)
			}
		})
	}
}

func TestReplace(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		old      string
		new      string
		n        int
		expected []string
	}{
		{
			name:     "test00",
			input:    []string{"a", "b", "a", "c", "d", "a"},
			old:      "a",
			new:      "x",
			n:        0,
			expected: []string{"a", "b", "a", "c", "d", "a"},
		},
		{
			name:     "test01",
			input:    []string{"a", "b", "a", "c", "d", "a"},
			old:      "a",
			new:      "x",
			n:        1,
			expected: []string{"x", "b", "a", "c", "d", "a"},
		},
		{
			name:     "test02",
			input:    []string{"a", "b", "a", "c", "d", "a"},
			old:      "a",
			new:      "x",
			n:        2,
			expected: []string{"x", "b", "x", "c", "d", "a"},
		},
		{
			name:     "test03",
			input:    []string{"a", "b", "a", "c", "d", "a"},
			old:      "a",
			new:      "x",
			n:        3,
			expected: []string{"x", "b", "x", "c", "d", "x"},
		},
		{
			name:     "test04",
			input:    []string{"a", "b", "a", "c", "d", "a"},
			old:      "a",
			new:      "x",
			n:        4,
			expected: []string{"x", "b", "x", "c", "d", "x"},
		},
		{
			name:     "test05",
			input:    []string{"a", "b", "a", "c", "d", "a"},
			old:      "a",
			new:      "x",
			n:        -1,
			expected: []string{"x", "b", "x", "c", "d", "x"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slicex.Replace(tt.input, tt.old, tt.new, tt.n); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Replace() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_RightPadding(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := []int{1, 2, 3, 4, 5, 0, 0, 0}

	result := slicex.RightPadding(input, 0, 3)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("RightPadding() expected %v, got %v", expected, result)
	}
}

func Test_RightPaddingZero(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := []int{1, 2, 3, 4, 5}

	result := slicex.RightPadding(input, 0, 0)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("RightPadding() expected %v, got %v", expected, result)
	}
}

func Test_LeftPadding(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := []int{0, 0, 0, 1, 2, 3, 4, 5}

	result := slicex.LeftPadding(input, 0, 3)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("LeftPadding() expected %v, got %v", expected, result)
	}
}

func Test_LeftPaddingZero(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := []int{1, 2, 3, 4, 5}

	result := slicex.LeftPadding(input, 0, 0)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("LeftPadding() expected %v, got %v", expected, result)
	}
}

func Test_Contains(t *testing.T) {
	tests := []struct {
		name     string
		src      []int
		dst      int
		expected bool
	}{
		{
			name:     "TC01",
			src:      []int{1, 2, 3, 4, 5},
			dst:      2,
			expected: true,
		},
		{
			name:     "TC02",
			src:      []int{1, 2, 3, 4, 5},
			dst:      0,
			expected: false,
		},
		{
			name:     "TC03",
			src:      []int{},
			dst:      5,
			expected: false,
		},
		{
			name:     "TC04",
			src:      nil,
			dst:      5,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slicex.Contains(tt.src, tt.dst)
			if got != tt.expected {
				t.Errorf("ContainsAll() expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func Test_ContainsAny(t *testing.T) {
	tests := []struct {
		name     string
		src      []int
		dst      []int
		expected bool
	}{
		{
			name:     "TC01",
			src:      []int{1, 2, 3, 4, 5},
			dst:      []int{2, 4},
			expected: true,
		},
		{
			name:     "TC02",
			src:      []int{1, 2, 3, 4, 5},
			dst:      []int{2, 6},
			expected: true,
		},
		{
			name:     "TC03",
			src:      []int{1, 2, 3, 4, 5},
			dst:      []int{},
			expected: false,
		},
		{
			name:     "TC04",
			src:      []int{},
			dst:      []int{1, 2, 3, 4, 5},
			expected: false,
		},
		{
			name:     "TC05",
			src:      nil,
			dst:      []int{2, 6},
			expected: false,
		},
		{
			name:     "TC06",
			src:      []int{2, 6},
			dst:      nil,
			expected: false,
		},
		{
			name:     "TC07",
			src:      nil,
			dst:      nil,
			expected: false,
		},
		{
			name:     "TC08",
			src:      []int{},
			dst:      []int{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slicex.ContainsAny(tt.src, tt.dst)
			if got != tt.expected {
				t.Errorf("ContainsAny() expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func Test_ContainsAll(t *testing.T) {
	tests := []struct {
		name     string
		src      []int
		dst      []int
		expected bool
	}{
		{
			name:     "TC01",
			src:      []int{1, 2, 3, 4, 5},
			dst:      []int{2, 4},
			expected: true,
		},
		{
			name:     "TC02",
			src:      []int{1, 2, 3, 4, 5},
			dst:      []int{2, 6},
			expected: false,
		},
		{
			name:     "TC03",
			src:      []int{1, 2, 3, 4, 5},
			dst:      []int{},
			expected: true,
		},
		{
			name:     "TC04",
			src:      []int{},
			dst:      []int{1, 2, 3, 4, 5},
			expected: false,
		},
		{
			name:     "TC05",
			src:      []int{},
			dst:      []int{},
			expected: true,
		},
		{
			name:     "TC06",
			src:      nil,
			dst:      []int{2, 6},
			expected: false,
		},
		{
			name:     "TC07",
			src:      []int{2, 6},
			dst:      nil,
			expected: true,
		},
		{
			name:     "TC08",
			src:      nil,
			dst:      nil,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slicex.ContainsAll(tt.src, tt.dst)
			if got != tt.expected {
				t.Errorf("ContainsAll() expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestFlattenDeep(t *testing.T) {
	input := [][][]string{{{"a", "b"}}, {{"c", "d"}}}
	expected := []string{"a", "b", "c", "d"}

	result := slicex.FlattenDeep(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FlattenDeep() expected %v, got %v", expected, result)
	}
}

// TestRemoveHead_Int 测试整型切片
func TestRemoveHead_Int(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "Normal case",
			input:    []int{1, 2, 3},
			expected: []int{2, 3},
		},
		{
			name:     "Empty slice",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "Nil slice",
			input:    nil,
			expected: nil,
		},
		{
			name:     "Single element",
			input:    []int{42},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := slicex.RemoveHead(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("RemoveHead(%v) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestRemoveHead_String 测试字符串切片
func TestRemoveHead_String(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "Multiple elements",
			input:    []string{"a", "b", "c"},
			expected: []string{"b", "c"},
		},
		{
			name:     "Empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "Nil slice",
			input:    nil,
			expected: nil,
		},
		{
			name:     "Single element",
			input:    []string{"hello"},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := slicex.RemoveHead(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("RemoveHead(%v) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestRemoveHead_Struct 测试结构体切片
func TestRemoveHead_Struct(t *testing.T) {
	type dummy struct {
		X int
	}
	tests := []struct {
		name     string
		input    []dummy
		expected []dummy
	}{
		{
			name:     "Multiple structs",
			input:    []dummy{{X: 1}, {X: 2}},
			expected: []dummy{{X: 2}},
		},
		{
			name:     "Empty slice",
			input:    []dummy{},
			expected: []dummy{},
		},
		{
			name:     "Nil slice",
			input:    nil,
			expected: nil,
		},
		{
			name:     "Single struct",
			input:    []dummy{{X: 3}},
			expected: []dummy{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := slicex.RemoveHead(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("RemoveHead(%+v) = %+v; want %+v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestRemoveTail_Int 测试整数切片的 RemoveTail 函数。
func TestRemoveTail_Int(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "Empty slice",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "Single element",
			input:    []int{1},
			expected: []int{},
		},
		{
			name:     "Multiple elements",
			input:    []int{1, 2, 3},
			expected: []int{1, 2},
		},
		{
			name:     "Nil slice",
			input:    nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := slicex.RemoveTail(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("RemoveTail(%v) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestRemoveTail_String 测试字符串切片的 RemoveTail 函数。
func TestRemoveTail_String(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "Empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "Single element",
			input:    []string{"hello"},
			expected: []string{},
		},
		{
			name:     "Multiple elements",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b"},
		},
		{
			name:     "Nil slice",
			input:    nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := slicex.RemoveTail(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("RemoveTail(%v) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestRemoveTail_Bool 测试布尔值切片的 RemoveTail 函数。
func TestRemoveTail_Bool(t *testing.T) {
	tests := []struct {
		name     string
		input    []bool
		expected []bool
	}{
		{
			name:     "Empty slice",
			input:    []bool{},
			expected: []bool{},
		},
		{
			name:     "Single element",
			input:    []bool{true},
			expected: []bool{},
		},
		{
			name:     "Multiple elements",
			input:    []bool{false, false, true},
			expected: []bool{false, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := slicex.RemoveTail(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("RemoveTail(%v) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestRemoveTail_Struct 测试结构体切片的 RemoveTail 函数。
func TestRemoveTail_Struct(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}

	tests := []struct {
		name     string
		input    []User
		expected []User
	}{
		{
			name:     "Empty slice",
			input:    []User{},
			expected: []User{},
		},
		{
			name:     "Single element",
			input:    []User{{ID: 1, Name: "Alice"}},
			expected: []User{},
		},
		{
			name:     "Multiple elements",
			input:    []User{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}},
			expected: []User{{ID: 1, Name: "Alice"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := slicex.RemoveTail(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("RemoveTail(%v) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	tests := []struct {
		name string
		s1   any
		s2   any
		want bool
	}{
		{
			name: "TC01 - Same int slices",
			s1:   []int{1, 2, 3},
			s2:   []int{1, 2, 3},
			want: true,
		},
		{
			name: "TC02 - Different last element in int slices",
			s1:   []int{1, 2, 3},
			s2:   []int{1, 2, 4},
			want: false,
		},
		{
			name: "TC03 - Different second string element",
			s1:   []string{"a", "b"},
			s2:   []string{"a", "c"},
			want: false,
		},
		{
			name: "TC04 - Both empty slices",
			s1:   []int{},
			s2:   []int{},
			want: true,
		},
		{
			name: "TC05 - Length mismatch",
			s1:   []int{1, 2},
			s2:   []int{},
			want: false,
		},
		{
			name: "TC06 - Same bool slices",
			s1:   []bool{true, false},
			s2:   []bool{true, false},
			want: true,
		},
		{
			name: "TC07 - Different bool values",
			s1:   []bool{true},
			s2:   []bool{false},
			want: false,
		},
		{
			name: "TC08 - Same float64 slices",
			s1:   []float64{1.1, 2.2},
			s2:   []float64{1.1, 2.2},
			want: true,
		},
		{
			name: "TC09 - Different float64 values",
			s1:   []float64{1.1, 2.2},
			s2:   []float64{1.1, 2.3},
			want: false,
		},
		{
			name: "TC10 - Same interface{} slices",
			s1:   []interface{}{1, "a"},
			s2:   []interface{}{1, "a"},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch s1 := tt.s1.(type) {
			case []int:
				s2 := tt.s2.([]int)
				got := slicex.Equal(s1, s2)
				if got != tt.want {
					t.Errorf("Equal() = %v, want %v", got, tt.want)
				}
			case []string:
				s2 := tt.s2.([]string)
				got := slicex.Equal(s1, s2)
				if got != tt.want {
					t.Errorf("Equal() = %v, want %v", got, tt.want)
				}
			case []bool:
				s2 := tt.s2.([]bool)
				got := slicex.Equal(s1, s2)
				if got != tt.want {
					t.Errorf("Equal() = %v, want %v", got, tt.want)
				}
			case []float64:
				s2 := tt.s2.([]float64)
				got := slicex.Equal(s1, s2)
				if got != tt.want {
					t.Errorf("Equal() = %v, want %v", got, tt.want)
				}
			case []interface{}:
				s2 := tt.s2.([]interface{})
				got := slicex.Equal(s1, s2)
				if got != tt.want {
					t.Errorf("Equal() = %v, want %v", got, tt.want)
				}
			default:
				t.Fatalf("Unsupported slice type in test: %T", tt.s1)
			}
		})
	}
}

// TestIsAscending_Types_Int 测试 int 类型
func TestIsAscending_Types_Int(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected bool
	}{
		{"Empty slice", []int{}, true},
		{"Single element", []int{1}, true},
		{"Strictly increasing", []int{1, 2, 3}, true},
		{"Decreasing in middle", []int{1, 3, 2}, false},
		{"All equal", []int{2, 2, 2}, true},
		{"Fully decreasing", []int{3, 2, 1}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slicex.IsAscending(tt.input); got != tt.expected {
				t.Errorf("IsAscending(%v) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}

// TestIsAscending_Types_String 测试 string 类型
func TestIsAscending_Types_String(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected bool
	}{
		{"Empty slice", []string{}, true},
		{"Single element", []string{"a"}, true},
		{"Lexicographical order", []string{"a", "b", "c"}, true},
		{"Out of order", []string{"a", "c", "b"}, false},
		{"All same", []string{"x", "x", "x"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slicex.IsAscending(tt.input); got != tt.expected {
				t.Errorf("IsAscending(%v) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}

// TestIsAscending_Types_Float 测试 float64 类型
func TestIsAscending_Types_Float(t *testing.T) {
	tests := []struct {
		name     string
		input    []float64
		expected bool
	}{
		{"Empty slice", []float64{}, true},
		{"Single element", []float64{1.1}, true},
		{"Strictly increasing", []float64{1.1, 1.2, 1.3}, true},
		{"Decreasing in middle", []float64{1.1, 1.3, 1.2}, false},
		{"All equal", []float64{2.2, 2.2, 2.2}, true},
		{"Fully decreasing", []float64{3.3, 2.2, 1.1}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slicex.IsAscending(tt.input); got != tt.expected {
				t.Errorf("IsAscending(%v) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}

// TestIsDescending 是 IsDescending 函数的单元测试
func TestIsDescending(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected bool
	}{
		{
			name:     "Empty slice",
			input:    []int{},
			expected: true,
		},
		{
			name:     "Single element (int)",
			input:    []int{5},
			expected: true,
		},
		{
			name:     "Fully descending (int)",
			input:    []int{5, 4, 3, 2, 1},
			expected: true,
		},
		{
			name:     "Descending with equal elements (int)",
			input:    []int{5, 5, 4, 4, 3},
			expected: true,
		},
		{
			name:     "Not descending (int)",
			input:    []int{5, 6, 4},
			expected: false,
		},
		{
			name:     "Descending (string)",
			input:    []string{"z", "y", "x"},
			expected: true,
		},
		{
			name:     "Ascending (string)",
			input:    []string{"a", "b", "c"},
			expected: false,
		},
		{
			name:     "Descending (float64)",
			input:    []float64{5.5, 4.4, 3.3},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result bool
			switch v := tt.input.(type) {
			case []int:
				result = slicex.IsDescending(v)
			case []string:
				result = slicex.IsDescending(v)
			case []float64:
				result = slicex.IsDescending(v)
			default:
				t.Fatalf("Unsupported type in test: %T", v)
			}

			if result != tt.expected {
				t.Errorf("Expected %v, got %v for input %v", tt.expected, result, tt.input)
			}
		})
	}
}

// TestCombine 是 Combine 函数的单元测试
func TestCombine(t *testing.T) {
	tests := []struct {
		name   string
		keys   []string
		values []int
		want   map[string]int
	}{
		{
			name:   "正常情况",
			keys:   []string{"a", "b"},
			values: []int{1, 2},
			want:   map[string]int{"a": 1, "b": 2},
		},
		{
			name:   "空切片",
			keys:   []string{},
			values: []int{},
			want:   map[string]int{},
		},
		{
			name:   "长度不一致",
			keys:   []string{"x"},
			values: []int{},
			want:   map[string]int{},
		},
		{
			name:   "nil 切片",
			keys:   nil,
			values: nil,
			want:   map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slicex.Combine(tt.keys, tt.values)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Combine() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCombine_GenericTypes 测试不同泛型类型的组合
func TestCombine_GenericTypes(t *testing.T) {
	// K=int, V=string
	{
		keys := []int{1, 2}
		values := []string{"one", "two"}
		want := map[int]string{1: "one", 2: "two"}

		got := slicex.Combine(keys, values)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Combine(int, string) = %v, want %v", got, want)
		}
	}

	// K=string, V=bool
	{
		keys := []string{"k1", "k2"}
		values := []bool{true, false}
		want := map[string]bool{"k1": true, "k2": false}

		got := slicex.Combine(keys, values)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Combine(string, bool) = %v, want %v", got, want)
		}
	}
}

// TestSortSlice_EmptySlice tests when the input slice is empty.
func TestSortSlice_EmptySlice(t *testing.T) {
	var s []int
	slicex.SortSlice(s, func(a, b int) bool { return a < b })
	if len(s) != 0 {
		t.Errorf("Expected empty slice, got length %d", len(s))
	}
}

// TestSortSlice_NilLessFunc tests when the 'less' function is nil.
func TestSortSlice_NilLessFunc(t *testing.T) {
	s := []int{3, 1, 2}
	slicex.SortSlice(s, nil)
	if !reflect.DeepEqual(s, []int{3, 1, 2}) {
		t.Errorf("Expected no change due to nil less func, got %v", s)
	}
}

// TestSortSlice_IntAscending tests sorting an int slice in ascending order.
func TestSortSlice_IntAscending(t *testing.T) {
	s := []int{3, 1, 2}
	expected := []int{1, 2, 3}
	slicex.SortSlice(s, func(a, b int) bool { return a < b })
	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Expected %v, got %v", expected, s)
	}
}

// TestSortSlice_StringDescending tests sorting a string slice in descending order.
func TestSortSlice_StringDescending(t *testing.T) {
	s := []string{"apple", "banana", "cherry"}
	expected := []string{"cherry", "banana", "apple"}
	slicex.SortSlice(s, func(a, b string) bool { return a > b })
	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Expected %v, got %v", expected, s)
	}
}

type Employee struct {
	FirstName string
	LastName  string
	Age       int
}

func TestSortSlice_StructField(t *testing.T) {
	input := []Employee{
		{"Bob", "Smith", 25},
		{"John", "Doe", 30},
		{"Jane", "Doe", 28},
	}
	expected := []Employee{
		{"Bob", "Smith", 25},
		{"Jane", "Doe", 28},
		{"John", "Doe", 30},
	}

	slicex.SortSlice(input, func(a, b Employee) bool {
		return a.Age < b.Age
	})
	if !reflect.DeepEqual(input, expected) {
		t.Errorf("SortSlice() expected %v, got %v", expected, input)
	}
}

func TestSortSlice_MapField(t *testing.T) {
	input := []map[string]interface{}{
		{"first_name": "Bob", "last_name": "Smith", "age": 25},
		{"first_name": "John", "last_name": "Doe", "age": 30},
		{"first_name": "Jane", "last_name": "Doe", "age": 28},
	}
	expected := []map[string]interface{}{
		{"first_name": "Bob", "last_name": "Smith", "age": 25},
		{"first_name": "Jane", "last_name": "Doe", "age": 28},
		{"first_name": "John", "last_name": "Doe", "age": 30},
	}
	slicex.SortSlice(input, func(a, b map[string]interface{}) bool {
		return a["age"].(int) < b["age"].(int)
	})
	if !reflect.DeepEqual(input, expected) {
		t.Errorf("Expected %v, got %v", expected, input)
	}
}

// TestFindUniques_Int_NoDuplicates tests a slice with no duplicates.
func TestFindUniques_Int_NoDuplicates(t *testing.T) {
	input := []int{1, 2, 3}
	expected := []int{1, 2, 3}

	if got := slicex.FindUniques(input); !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}

// TestFindUniques_Int_WithDuplicates tests a slice with duplicates.
func TestFindUniques_Int_WithDuplicates(t *testing.T) {
	input := []int{1, 2, 1, 3}
	expected := []int{2, 3}

	if got := slicex.FindUniques(input); !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}

// TestFindUniques_String_WithDuplicates tests string slice with duplicates.
func TestFindUniques_String_WithDuplicates(t *testing.T) {
	input := []string{"a", "b", "a", "c"}
	expected := []string{"b", "c"}

	if got := slicex.FindUniques(input); !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}

// TestFindUniques_EmptySlice tests empty input.
func TestFindUniques_EmptySlice(t *testing.T) {
	if got := slicex.FindUniques([]int{}); !reflect.DeepEqual([]int{}, got) {
		t.Errorf("Expected %v, got %v", []int{}, got)
	}
}

// TestFindUniques_AllSameElements tests all elements are same.
func TestFindUniques_AllSameElements(t *testing.T) {
	input := []int{1, 1, 1, 1}

	if got := slicex.FindUniques(input); !reflect.DeepEqual([]int{}, got) {
		t.Errorf("Expected %v, got %v", []int{}, got)
	}
}

// TestFindUniques_InterfaceMixedType tests mixed type in interface slice.
func TestFindUniques_InterfaceMixedType(t *testing.T) {
	input := []interface{}{1, "a", 1, "a"}

	if got := slicex.FindUniques(input); !reflect.DeepEqual([]interface{}{}, got) {
		t.Errorf("Expected %v, got %v", []interface{}{}, got)
	}
}

// TestFindUniques_Int_Interleaved tests interleaved duplicates.
func TestFindUniques_Int_Interleaved(t *testing.T) {
	input := []int{2, 3, 2, 3, 4}
	expected := []int{4}

	if got := slicex.FindUniques(input); !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}

func TestWithout(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		excludes []int
		expected []int
	}{
		{
			name:     "Normal case",
			input:    []int{1, 2, 3, 4},
			excludes: []int{2, 4},
			expected: []int{1, 3},
		},
		{
			name:     "Empty input slice",
			input:    []int{},
			excludes: []int{1, 2},
			expected: []int{},
		},
		{
			name:     "No elements to exclude",
			input:    []int{1, 2, 3},
			excludes: []int{},
			expected: []int{1, 2, 3},
		},
		{
			name:     "All elements excluded",
			input:    []int{1, 2, 3},
			excludes: []int{1, 2, 3},
			expected: []int{},
		},
		{
			name:     "Duplicate elements",
			input:    []int{1, 2, 2, 3},
			excludes: []int{2},
			expected: []int{1, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slicex.Without(tt.input, tt.excludes...)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestWithout_GenericString(t *testing.T) {
	input := []string{"a", "b", "c"}
	excludes := []string{"b"}
	expected := []string{"a", "c"}

	got := slicex.Without(input, excludes...)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}
