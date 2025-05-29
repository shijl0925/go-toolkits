package mutable_test

import (
	"github.com/shijl0925/go-toolkits/mutable"
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
		fn := func(i int) bool { return true }

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
