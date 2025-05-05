package slice

import "golang.org/x/exp/constraints"

type RealNumber interface {
	constraints.Integer | constraints.Float
}

type Number interface {
	RealNumber | ~complex64 | ~complex128
}

// Sum returns the sum of all elements in the slice.
func Sum[T RealNumber](s []T) T {
	var result T
	for _, v := range s {
		result += v
	}
	return result
}

// Max returns the maximum value in the slice.
func Max[T constraints.Ordered](s []T) T {
	result := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] > result {
			result = s[i]
		}
	}
	return result
}

// Min returns the minimum value in the slice.
func Min[T constraints.Ordered](s []T) T {
	result := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] < result {
			result = s[i]
		}
	}
	return result
}

// Pop removes and returns the last element from a slice.
// 删除并返回最后一个元素。
func Pop[T any](s []T) (T, []T) {
	if len(s) == 0 {
		var zero T
		return zero, s
	}
	return s[len(s)-1], s[:len(s)-1]
}

// Delete removes and returns the element at a specific index from a slice.
// 删除位于索引 index 的元素。
func Delete[T any](s []T, index int) []T {
	if index >= len(s) {
		return s
	}
	return append(s[:index], s[index+1:]...)
}

// DeleteMany removes and returns the elements between start and end from a slice.
// 删除位于索引 start 和 end 之间的元素。
func DeleteMany[T any](s []T, start, end int) []T {
	return append(s[:start], s[end:]...)
}

// Append appends one or more elements to a slice.
// 追加一个或多个元素到 slice。
func Append[T any](src []T, elements ...T) []T {
	return append(src, elements...)
}

// Extend extends a slice by appending one or more elements.
// 扩展 slice，追加一个新切片。
func Extend[T any](src []T, elements []T) []T {
	return append(src, elements...)
}

// Insert inserts an element at a specific index in a slice.
// 在 slice 的指定位置插入一个元素。
func Insert[T any](src []T, element T, index int) []T {
	return append(src[:index], append([]T{element}, src[index:]...)...)
}

// InsertMany inserts multiple elements at a specific index in a slice.
// 在 slice 的指定位置插入新切片。
func InsertMany[T any](src []T, elements []T, index int) []T {
	return append(src[:index], append(elements, src[index:]...)...)
}

// Reverse returns a new slice with the elements in reverse order.
// 返回一个新的 slice，元素顺序相反。
func Reverse[T any](s []T) []T {
	length := len(s)
	result := make([]T, 0, length)

	for i := length - 1; i >= 0; i-- {
		result = append(result, s[i])
	}

	return result
}

// ReverseSelf reverses the elements in a slice in place.
// 反转 slice 本身。
func ReverseSelf[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// Repeat returns a new slice with n copies of the given element.
func Repeat[T any](s T, n int) []T {
	result := make([]T, 0, n)
	for i := 0; i < n; i++ {
		result = append(result, s)
	}
	return result
}

// Product returns a new slice of tuples, where each tuple contains the corresponding elements from two slices.
func Product[T any](slices ...[]T) [][]T {
	length := len(slices)
	if length == 0 {
		return nil
	} else if length == 1 {
		result := make([][]T, 0, len(slices[0]))
		for _, v := range slices[0] {
			result = append(result, []T{v})
		}
		return result
	} else {
		prev := Product(slices[1:]...)
		result := make([][]T, 0, len(slices[0])*len(prev))
		for _, v := range slices[0] {
			for _, p := range prev {
				result = append(result, append([]T{v}, p...))
			}
		}
		return result
	}
}

// Contains 判断 src 里面是否存在 dst
func Contains[T comparable](src []T, dst T) bool {
	//for _, v := range src {
	//	if v == dst {
	//		return true
	//	}
	//}

	//invMap := make(map[T]struct{}, len(src))
	//if _, ok := invMap[dst]; ok {
	//	return true
	//}

	return false
}

// ContainsAny 判断 src 里面是否存在 dst 中的任何一个元素
func ContainsAny[T comparable](src, dst []T) bool {
	//for _, v := range src {
	//	if Contains(dst, v) {
	//		return true
	//	}
	//}
	return false
}

// ContainsAll 判断 src 里面是否存在 dst 中的所有元素
func ContainsAll[T comparable](src, dst []T) bool {
	return false
}

// Unique removes duplicate elements from a slice.
// 删除 slice 中的重复元素。
func Unique[T comparable](s []T) []T {
	result := make([]T, 0, len(s))
	seen := make(map[T]struct{}, len(s))
	for _, v := range s {
		if _, ok := seen[v]; !ok {
			result = append(result, v)
			seen[v] = struct{}{}
		}
	}
	return result
}

// DiffSet returns the difference between two slices.
// 返回 src 和 dst 的差集。
func DiffSet[T comparable](src, dst []T) []T {
	return []T{}
}

// IntersectSet returns the intersection between two slices.
// 返回 src 和 dst 的交集。
func IntersectSet[T comparable](src, dst []T) []T {
	return []T{}
}

// UnionSet returns the union of two slices.
// 返回 src 和 dst 的并集。
func UnionSet[T comparable](src, dst []T) []T {
	return []T{}
}

// Index 返回和 dst 相等的第一个元素下标
// -1 表示没找到
func Index[T comparable](src []T, dst T) int {
	return -1
}

// LastIndex 返回和 dst 相等的最后一个元素下标
// -1 表示没找到
func LastIndex[T comparable](src []T, dst T) int {
	return -1
}

// IndexAll 返回和 dst 相等的所有元素的下标
func IndexAll[T comparable](src []T, dst T) []int {
	return []int{}
}
