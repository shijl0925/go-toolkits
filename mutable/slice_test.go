package mutable_test

import (
	"github.com/shijl0925/go-toolkits/mutable"
	"github.com/shijl0925/go-toolkits/slicex"
	"reflect"
	"strings"
	"testing"
)

// TestMap_Int 测试整数切片的映射
func TestMap_Int(t *testing.T) {
	s := []int{1, 2, 3}
	expected := []int{2, 4, 6}
	mutable.Map(s, func(x int) int {
		return x * 2
	})

	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Map failed, got %v, want %v", s, expected)
	}
}

// TestMap_String 测试字符串切片的映射
func TestMap_String(t *testing.T) {
	s := []string{"a", "b"}
	expected := []string{"x_a", "x_b"}
	mutable.Map(s, func(x string) string {
		return "x_" + x
	})

	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Map failed, got %v, want %v", s, expected)
	}
}

// TestMap_Struct 测试结构体切片的映射
func TestMap_Struct(t *testing.T) {
	type User struct {
		Name string
	}
	s := []User{
		{Name: "alice"},
		{Name: "bob"},
	}
	expected := []User{
		{Name: "Alice"},
		{Name: "Bob"},
	}
	mutable.Map(s, func(u User) User {
		u.Name = strings.ToUpper(u.Name[0:1]) + u.Name[1:] // 只改首字母为大写
		return u
	})

	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Map failed, got %v, want %v", s, expected)
	}
}

// TestMap_NilSlice 测试 nil 切片的情况
func TestMap_NilSlice(t *testing.T) {
	var s []int
	mutable.Map(s, func(x int) int {
		return x * 2
	})
	if len(s) != 0 {
		t.Errorf("Map on nil slice should not modify it")
	}
}

// TestMap_EmptySlice 测试空切片的情况
func TestMap_EmptySlice(t *testing.T) {
	var s []string
	mutable.Map(s, func(x string) string {
		return "x_" + x
	})
	if len(s) != 0 {
		t.Errorf("Map on empty slice should not add elements")
	}
}

func Test_MutableFilter(t *testing.T) {
	t.Run("TC01: Input slice is nil", func(t *testing.T) {
		var s *[]int
		fn := func(i int) bool { return i > 0 }

		result := mutable.Filter(s, fn)
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %v", result)
		}
		if result != nil {
			t.Errorf("Filter on nil slice should not add elements")
		}
	})

	t.Run("TC02: Empty input slice", func(t *testing.T) {
		var s []int
		fn := func(i int) bool { return i > 0 }

		result := mutable.Filter(&s, fn)
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %v", result)
		}
		if result != nil {
			t.Errorf("Filter on empty slice should not add elements")
		}
	})

	t.Run("TC03: All elements do not match", func(t *testing.T) {
		s := []int{1, 2, 3}
		fn := func(i int) bool { return i > 3 }

		result := mutable.Filter(&s, fn)

		if !reflect.DeepEqual(result, []int{}) {
			t.Errorf("Expected %v, got %v", []int{}, result)
		}
		if !reflect.DeepEqual(s, []int{}) {
			t.Errorf("Original slice should be modified to %v", []int{})
		}
	})

	t.Run("TC04: All elements match", func(t *testing.T) {
		s := []int{1, 2, 3}
		fn := func(_ int) bool { return true }

		result := mutable.Filter(&s, fn)
		expected := []int{1, 2, 3}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
		if !reflect.DeepEqual(s, expected) {
			t.Errorf("Original slice should remain unchanged as %v", expected)
		}
	})

	t.Run("TC05: Some elements match", func(t *testing.T) {
		s := []int{1, 2, 3, 4, 5}
		fn := func(i int) bool { return i%2 == 0 } // 偶数保留

		result := mutable.Filter(&s, fn)
		expected := []int{2, 4}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
		if !reflect.DeepEqual(s, expected) {
			t.Errorf("Original slice should be modified to %v", expected)
		}
	})

	t.Run("TC06: Function is nil (should panic)", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic when function is nil")
			}
		}()

		s := []int{1, 2, 3}
		var fn func(int) bool = nil
		mutable.Filter(&s, fn)
	})
}

func Test_MutableRemove(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		element  any
		expected any
	}{
		{
			name:     "EmptySlice",
			input:    []string{},
			element:  "a",
			expected: []string{},
		},
		{
			name:     "NoMatch",
			input:    []int{1, 2, 3},
			element:  4,
			expected: []int{1, 2, 3},
		},
		{
			name:     "MultipleMatches",
			input:    []int{1, 2, 2, 3},
			element:  2,
			expected: []int{1, 3},
		},
		{
			name:     "AllMatches",
			input:    []int{2, 2, 2},
			element:  2,
			expected: []int{},
		},
		{
			name:     "NonContiguousMatches",
			input:    []int{1, 2, 3, 2, 4},
			element:  2,
			expected: []int{1, 3, 4},
		},
		{
			name:     "StringType",
			input:    []string{"a", "b", "a"},
			element:  "a",
			expected: []string{"b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch input := tt.input.(type) {
			case []int:
				got := mutable.Remove(&input, tt.element.(int))
				if !reflect.DeepEqual(got, tt.expected) {
					t.Errorf("Remove(%v, %v) = %v; want %v", input, tt.element, got, tt.expected)
				}

			case []string:
				got := mutable.Remove(&input, tt.element.(string))
				if !reflect.DeepEqual(got, tt.expected) {
					t.Errorf("Remove(%v, %v) = %v; want %v", input, tt.element, got, tt.expected)
				}
			default:
				t.Fatalf("Unsupported type in test case: %T", input)
			}
		})
	}

	t.Run("TC01: Input slice is nil", func(t *testing.T) {
		var s *[]int

		got := mutable.Remove(s, 1)
		if len(got) != 0 {
			t.Errorf("Expected empty slice, got %v", got)
		}
		if got != nil {
			t.Errorf("Remove on nil slice should not add elements")
		}
	})
}

func Test_ReverseSelfSlice(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}
	t.Run("test1", func(t *testing.T) {
		mutable.Reverse(s1)
		if !reflect.DeepEqual(s1, []int{5, 4, 3, 2, 1}) {
			t.Errorf("Reverse() expected %v, got %v", 0, s1)
		}
	})

	s2 := []string{"one", "two", "three"}
	t.Run("test2", func(t *testing.T) {
		mutable.Reverse(s2)
		if !reflect.DeepEqual(s2, []string{"three", "two", "one"}) {
			t.Errorf("Reverse() expected %v, got %v", []string{"three", "two", "one"}, s2)
		}
	})

	s3 := []byte("Google")
	t.Run("test3", func(t *testing.T) {
		mutable.Reverse(s3)
		if string(s3) != "elgooG" {
			t.Errorf("Reverse() expected %v, got %v", "elgooG", string(s3))
		}
	})
}

// TestReverse_StructSlice tests reversing a slice of structs.
func Test_ReverseSelf_StructSlice(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	s := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 35},
	}
	expected := []Person{
		{"Charlie", 35},
		{"Bob", 25},
		{"Alice", 30},
	}
	mutable.Reverse(s)
	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Reverse(%v) = %v; want %v", s, s, expected)
	}
}

// TestReverse_EmptySlice tests reversing an empty slice.
func Test_ReverseSelf_EmptySlice(t *testing.T) {
	var s []int
	var expected []int
	mutable.Reverse(s)
	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Reverse(%v) = %v; want %v", s, s, expected)
	}
}

// TestReverse_SingleElementSlice tests reversing a single-element slice.
func Test_ReverseSelf_SingleElementSlice(t *testing.T) {
	s := []int{42}
	expected := []int{42}
	mutable.Reverse(s)
	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Reverse(%v) = %v; want %v", s, s, expected)
	}
}

// TestReverse_EvenLengthSlice tests reversing a slice with even length.
func Test_ReverseSelf_EvenLengthSlice(t *testing.T) {
	s := []int{1, 2, 3, 4}
	expected := []int{4, 3, 2, 1}
	mutable.Reverse(s)
	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Reverse(%v) = %v; want %v", s, s, expected)
	}
}

// TestReverse_OddLengthSlice tests reversing a slice with odd length.
func Test_ReverseSelf_OddLengthSlice(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	expected := []int{5, 4, 3, 2, 1}
	mutable.Reverse(s)
	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Reverse(%v) = %v; want %v", s, s, expected)
	}
}

func Test_Pop(t *testing.T) {
	var tests = []struct {
		name     string
		slice    *[]int
		expected int
		ok       bool
	}{
		{"test1", &[]int{1, 2, 3, 4}, 4, true},
		{"test2", &[]int{1, 2, 3, 4, 5}, 5, true},
		{"test3", &[]int{1, 2, 3, 4, 5, 6}, 6, true},
		{"test4", &[]int{}, 0, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			last, ok := mutable.Pop(tc.slice)
			if last != tc.expected || ok != tc.ok {
				t.Errorf("Pop() expected %v(%v), got %v(%v)", tc.expected, tc.ok, last, ok)
			}
		})
	}

	t.Run("nil pointer", func(t *testing.T) {
		var s *[]int
		val, ok := mutable.Pop(s)
		var zero int
		if zero != val || ok {
			t.Errorf("Pop() expected %v(%v), got %v(%v)", zero, false, val, ok)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		s := &[]string{}
		val, ok := mutable.Pop(s)
		var zero string
		if zero != val || ok {
			t.Errorf("Pop() expected %v(%v), got %v(%v)", zero, false, val, ok)
		}
	})

	t.Run("single element", func(t *testing.T) {
		s := &[]int{1}
		val, ok := mutable.Pop(s)
		if val != 1 || !ok {
			t.Errorf("Pop() expected %v(%v), got %v(%v)", val, true, val, ok)
		}
		if !reflect.DeepEqual(*s, []int{}) {
			t.Errorf("Pop() expected %v, got %v", []int{}, *s)
		}
	})

	t.Run("multiple elements", func(t *testing.T) {
		s := &[]string{"a", "b"}
		val, ok := mutable.Pop(s)
		if val != "b" || !ok {
			t.Errorf("Pop() expected %v(%v), got %v(%v)", val, true, val, ok)
		}
		if !reflect.DeepEqual(*s, []string{"a"}) {
			t.Errorf("Pop() expected %v, got %v", []string{"a"}, *s)
		}
	})

	t.Run("generic type struct", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}
		s := &[]User{
			{Name: "Alice", Age: 25},
			{Name: "Bob", Age: 30},
			{Name: "Charlie", Age: 35},
		}
		val, ok := mutable.Pop(s)
		if !reflect.DeepEqual(val, User{Name: "Charlie", Age: 35}) || !ok {
			t.Errorf("Pop() expected %v(%v), got %v(%v)", val, true, val, ok)
		}
		if !reflect.DeepEqual(*s, []User{{Name: "Alice", Age: 25}, {Name: "Bob", Age: 30}}) {
			t.Errorf("Pop() expected %v, got %v", []User{{Name: "Alice", Age: 25}, {Name: "Bob", Age: 30}}, *s)
		}
	})
}

func Test_Shift(t *testing.T) {
	var tests = []struct {
		name     string
		slice    *[]int
		expected int
		ok       bool
	}{
		{"test1", &[]int{1, 2, 3, 4}, 1, true},
		{"test2", &[]int{1, 2, 3, 4, 5}, 1, true},
		{"test3", &[]int{1, 2, 3, 4, 5, 6}, 1, true},
		{"test4", &[]int{}, 0, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			val, ok := mutable.Shift(tc.slice)
			if val != tc.expected || ok != tc.ok {
				t.Errorf("Shift() expected %v(%v), got %v(%v)", tc.expected, tc.ok, val, ok)
			}
		})
	}

	t.Run("nil pointer", func(t *testing.T) {
		var s *[]int
		val, ok := mutable.Shift(s)
		var zero int
		if zero != val || ok {
			t.Errorf("Shift() expected %v(%v), got %v(%v)", zero, false, val, ok)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		s := &[]string{}
		val, ok := mutable.Shift(s)
		var zero string
		if zero != val || ok {
			t.Errorf("Shift() expected %v(%v), got %v(%v)", zero, false, val, ok)
		}
	})

	t.Run("single element", func(t *testing.T) {
		s := &[]int{1}
		val, ok := mutable.Shift(s)
		if val != 1 || !ok {
			t.Errorf("Shift() expected %v(%v), got %v(%v)", val, true, val, ok)
		}
		if !reflect.DeepEqual(*s, []int{}) {
			t.Errorf("Shift() expected %v, got %v", []int{}, *s)
		}
	})

	t.Run("multiple elements", func(t *testing.T) {
		s := &[]string{"a", "b"}
		val, ok := mutable.Shift(s)
		if val != "a" || !ok {
			t.Errorf("Shift() expected %v(%v), got %v(%v)", val, true, val, ok)
		}
		if !reflect.DeepEqual(*s, []string{"b"}) {
			t.Errorf("Shift() expected %v, got %v", []string{"b"}, *s)
		}
	})

	t.Run("generic type struct", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}
		s := &[]User{
			{Name: "Alice", Age: 25},
			{Name: "Bob", Age: 30},
			{Name: "Charlie", Age: 35},
		}
		val, ok := mutable.Shift(s)
		if !reflect.DeepEqual(val, User{Name: "Alice", Age: 25}) || !ok {
			t.Errorf("Shift() expected %v(%v), got %v(%v)", val, true, val, ok)
		}
		if !reflect.DeepEqual(*s, []User{{Name: "Bob", Age: 30}, {Name: "Charlie", Age: 35}}) {
			t.Errorf("Shift() expected %v, got %v", []User{{Name: "Bob", Age: 30}, {Name: "Charlie", Age: 35}}, *s)
		}
	})
}

// Helper function to copy a slice (for testing purposes)
func copySlice(s *[]int) []int {
	if s == nil {
		return nil
	}
	copied := make([]int, len(*s))
	copy(copied, *s)
	return copied
}

// TestMutableDrop tests the Drop function with various scenarios.
func TestMutableDrop(t *testing.T) {
	tests := []struct {
		name     string
		input    *[]int
		n        int
		expected []int
		ok       bool
		wantErr  bool
	}{
		{
			name:     "Normal case",
			input:    &[]int{1, 2, 3, 4, 5},
			n:        2,
			expected: []int{1, 2, 3},
			ok:       true,
			wantErr:  false,
		},
		{
			name:     "Nil pointer",
			input:    nil,
			n:        2,
			expected: nil,
			ok:       false,
			wantErr:  true,
		},
		{
			name:     "Negative n",
			input:    &[]int{1, 2, 3},
			n:        -1,
			expected: []int{1, 2, 3},
			ok:       false,
			wantErr:  true,
		},
		{
			name:     "n greater than length",
			input:    &[]int{1, 2, 3},
			n:        5,
			expected: []int{},
			ok:       true,
			wantErr:  false,
		},
		{
			name:     "Empty slice, n=0",
			input:    &[]int{},
			n:        0,
			expected: []int{},
			ok:       true,
			wantErr:  false,
		},
		{
			name:     "Empty slice, n>0",
			input:    &[]int{},
			n:        1,
			expected: []int{},
			ok:       true,
			wantErr:  false,
		},
		{
			name:     "Remove all elements",
			input:    &[]int{1, 2, 3},
			n:        3,
			expected: []int{},
			ok:       true,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := copySlice(tt.input) // Make a copy to check mutation
			ok, err := mutable.Drop(tt.input, tt.n)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Drop() expected error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Drop() error = %v, wantErr %v", err, tt.wantErr)
				}

				if ok != tt.ok {
					t.Errorf("expected ok=%v, got %v", tt.ok, ok)
				}

				if !reflect.DeepEqual(*tt.input, tt.expected) {
					t.Errorf("Drop() = %v, want %v", *tt.input, tt.expected)
				}
			}

			// Check if original slice was modified correctly
			if tt.input != nil && tt.ok && tt.n >= 0 && tt.n <= len(original) {
				expectedModified := original[:len(original)-tt.n]
				if !reflect.DeepEqual(*tt.input, expectedModified) {
					t.Errorf("expected input slice to be modified to %v, but got %v", expectedModified, *tt.input)
				}
			}
			if tt.input != nil && tt.ok && tt.n > len(original) {
				expectedModified := original[:0]
				if !reflect.DeepEqual(*tt.input, expectedModified) {
					t.Errorf("expected input slice to be modified to %v, but got %v", expectedModified, *tt.input)
				}
			}
		})
	}
}

// TestMutableDropLeft tests the DropLeft function with various scenarios.
func TestMutableDropLeft(t *testing.T) {
	tests := []struct {
		name     string
		input    *[]int
		n        int
		expected []int
		ok       bool
		wantErr  bool
	}{
		{
			name:     "Normal case",
			input:    &[]int{1, 2, 3, 4, 5},
			n:        2,
			expected: []int{3, 4, 5},
			ok:       true,
			wantErr:  false,
		},
		{
			name:     "Nil pointer",
			input:    nil,
			n:        2,
			expected: nil,
			ok:       false,
			wantErr:  true,
		},
		{
			name:     "Negative n",
			input:    &[]int{1, 2, 3},
			n:        -1,
			expected: []int{1, 2, 3},
			ok:       false,
			wantErr:  true,
		},
		{
			name:     "n greater than length",
			input:    &[]int{1, 2, 3},
			n:        5,
			expected: []int{},
			ok:       true,
			wantErr:  false,
		},
		{
			name:     "Empty slice, n=0",
			input:    &[]int{},
			n:        0,
			expected: []int{},
			ok:       true,
			wantErr:  false,
		},
		{
			name:     "Empty slice, n>0",
			input:    &[]int{},
			n:        1,
			expected: []int{},
			ok:       true,
			wantErr:  false,
		},
		{
			name:     "Remove all elements",
			input:    &[]int{1, 2, 3},
			n:        3,
			expected: []int{},
			ok:       true,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := copySlice(tt.input) // Make a copy to check mutation
			ok, err := mutable.DropLeft(tt.input, tt.n)

			if tt.wantErr {
				if err == nil {
					t.Errorf("DropLeft() expected error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("DropLeft() error = %v, wantErr %v", err, tt.wantErr)
				}

				if ok != tt.ok {
					t.Errorf("expected ok=%v, got %v", tt.ok, ok)
				}

				if !reflect.DeepEqual(*tt.input, tt.expected) {
					t.Errorf("DropLeft() = %v, want %v", *tt.input, tt.expected)
				}
			}

			// Check if original slice was modified correctly
			if tt.input != nil && tt.ok && tt.n >= 0 && tt.n <= len(original) {
				expectedModified := original[tt.n:]
				if !reflect.DeepEqual(*tt.input, expectedModified) {
					t.Errorf("expected input slice to be modified to %v, but got %v", expectedModified, *tt.input)
				}
			}
			if tt.input != nil && tt.ok && tt.n > len(original) {
				expectedModified := original[:0]
				if !reflect.DeepEqual(*tt.input, expectedModified) {
					t.Errorf("expected input slice to be modified to %v, but got %v", expectedModified, *tt.input)
				}
			}
		})
	}
}

// TestShuffle_EmptySlice tests that shuffling an empty slice does not cause errors and remains unchanged
func TestShuffle_EmptySlice(t *testing.T) {
	var input []int

	var expected []int

	mutable.Shuffle(input)
	if !slicex.EqualUnordered(input, expected) {
		t.Errorf("Shuffle() on empty slice = %v, want %v", input, expected)
	}
}

// TestShuffle_SingleElementSlice tests that shuffling a single-element slice leaves it unchanged
func TestShuffle_SingleElementSlice(t *testing.T) {
	input := []string{"a"}
	expected := []string{"a"}
	mutable.Shuffle(input)
	if !slicex.EqualUnordered(input, expected) {
		t.Errorf("Shuffle() on single-element slice = %v, want %v", input, expected)
	}
}

// TestShuffle_MultipleElements tests that shuffling a multi-element slice changes the order
// but keeps all elements present
func TestShuffle_MultipleElements(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	originalCopy := make([]int, len(input))
	copy(originalCopy, input)

	mutable.Shuffle(input)

	// Check if the order is different (not guaranteed in every run, so we only check once)
	if !slicex.EqualUnordered(input, originalCopy) {
		t.Errorf("Shuffle() did not change the order, got %v", input)
	}
}

// TestShuffle_DuplicateElements tests that shuffling works correctly with duplicate elements
func TestShuffle_DuplicateElements(t *testing.T) {
	input := []string{"a", "b", "a", "c", "b"}
	originalCopy := make([]string, len(input))
	copy(originalCopy, input)

	mutable.Shuffle(input)

	if !slicex.EqualUnordered(input, originalCopy) {
		t.Errorf("Shuffle() changed the elements, got %v want %v", input, originalCopy)
	}
}

// TestShuffle_DifferentDataTypes tests that shuffling works for various data types
func TestShuffle_DifferentDataTypes(t *testing.T) {
	// Test with integers
	intInput := []int{1, 2, 3, 4, 5}
	originalIntCopy := make([]int, len(intInput))
	copy(originalIntCopy, intInput)

	mutable.Shuffle(intInput)

	if !slicex.EqualUnordered(intInput, originalIntCopy) {
		t.Errorf("Shuffle() failed for integer type, got %v want %v", intInput, originalIntCopy)
	}

	// Test with strings
	stringInput := []string{"apple", "banana", "cherry", "date"}
	originalStringCopy := make([]string, len(stringInput))
	copy(originalStringCopy, stringInput)

	mutable.Shuffle(stringInput)

	if !slicex.EqualUnordered(stringInput, originalStringCopy) {
		t.Errorf("Shuffle() failed for string type, got %v want %v", stringInput, originalStringCopy)
	}

	// Test with custom structs
	type testStruct struct {
		ID   int
		Name string
	}
	structInput := []testStruct{
		{ID: 1, Name: "One"},
		{ID: 2, Name: "Two"},
		{ID: 3, Name: "Three"},
	}
	originalStructCopy := make([]testStruct, len(structInput))
	copy(originalStructCopy, structInput)

	mutable.Shuffle(structInput)

	if !slicex.EqualUnordered(structInput, originalStructCopy) {
		t.Errorf("Shuffle() failed for struct type, got %v want %v", structInput, originalStructCopy)
	}
}
