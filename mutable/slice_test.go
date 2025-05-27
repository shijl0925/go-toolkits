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
