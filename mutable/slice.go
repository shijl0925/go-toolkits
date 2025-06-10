package mutable

import (
	"fmt"
	"math/rand"
)

// Map applies a function to each element of a slice.
// It modifies the slice in place.
func Map[T any](s []T, fn func(T) T) {
	for i, v := range s {
		s[i] = fn(v)
	}
}

func Filter[T any](s *[]T, fn func(T) bool) []T {
	// 输入切片为 nil, 会 panic, 所以需要返回 nil
	if s == nil {
		return nil
	}

	i := 0
	for _, v := range *s {
		if fn(v) {
			(*s)[i] = v
			i++
		}
	}
	*s = (*s)[:i]
	return *s
}

func Remove[T comparable](s *[]T, element T) []T {
	if s == nil {
		return nil
	}

	i := 0
	for _, v := range *s {
		if v != element {
			(*s)[i] = v
			i++
		}
	}
	*s = (*s)[:i]
	return *s
}

// Reverse reverses the elements in a slice in place.
// 反转 slice 本身
//
// Example:
//
// s := []int{1, 2, 3, 4, 5}
// Reverse(s)  // s == []int{5, 4, 3, 2, 1}
func Reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// Pop removes and returns the last element from a slice.
// 从切片 s 中弹出最后一个元素，并返回该元素和弹出状态。
// 如果 slice 为空，则返回一个零值和false。
// 注意：原始 slice 会被改变。
//
// Example:
//
//	s := &[]int{1, 2, 3, 4, 5}
//	r, ok := Pop(s) //  r == 5, ok ==true, s == []int{1, 2, 3, 4}
func Pop[T any](s *[]T) (T, bool) {
	if s == nil || len(*s) == 0 {
		var zero T
		return zero, false
	}
	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v, true
}

// Shift removes the slice's first element and returns the removed element.
// 从切片 s 中弹出第一个元素，并返回该元素和弹出状态。
// 如果 slice 为空，则返回一个零值和false。
// 注意：原始 slice 会被改变。
//
// Example:
//
//	s := &[]int{1, 2, 3, 4, 5}
//	r, ok := Shift(s) //  r == 1, ok ==true, s == []int{2, 3, 4, 5}
func Shift[T any](s *[]T) (T, bool) {
	if s == nil || len(*s) == 0 {
		var zero T
		return zero, false
	}
	v := (*s)[0]
	*s = (*s)[1:]
	return v, true
}

// Drop removes the n elements from the slice and returns the remaining elements.
// 删除并返回切片的前 n 个元素和剩余元素组成的切片, 原始 slice 会被改变。
// Example:
//
//	s := &[]int{1, 2, 3, 4, 5}
//	r, _ := Drop(s, 2) //  r == true  s ==  []int{1, 2, 3}
func Drop[T any](s *[]T, n int) (bool, error) {
	if s == nil {
		return false, fmt.Errorf("cannot drop from a nil slice")
	}

	if n < 0 {
		return false, fmt.Errorf("number of elements to drop cannot be negative")
	}

	if n > len(*s) {
		n = len(*s)
	}

	*s = (*s)[:len(*s)-n]
	return true, nil
}

// DropLeft removes the n elements from the slice and returns the remaining elements.
// 删除并返回切片的左侧 n 个元素和剩余元素组成的切片, 原始 slice 会被改变。
// Example:
//
//	s := &[]int{1, 2, 3, 4, 5}
//	r, _ := DropLeft(s, 2) //  r == true  s ==  []int{3, 4, 5}
func DropLeft[T any](s *[]T, n int) (bool, error) {
	if s == nil {
		return false, fmt.Errorf("cannot drop from a nil slice")
	}

	if n < 0 {
		return false, fmt.Errorf("number of elements to drop cannot be negative")
	}

	if n > len(*s) {
		n = len(*s)
	}

	*s = (*s)[n:]
	return true, nil
}

// Shuffle returns an array of shuffled values. Uses the Fisher-Yates shuffle algorithm.
func Shuffle[T any](s []T) {
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
}
