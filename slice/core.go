package slice

import (
	"golang.org/x/exp/constraints"
)

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

// removeHead removes the first element from a slice and returns the remaining elements.
// 删除第一个元素后组成的切片, 原始 slice 不会被改变。
func removeHead[T any](s []T) []T {
	if len(s) == 0 {
		return s
	}
	return s[1:]
}

// removeTail removes the last element from a slice and returns the remaining elements.
// 删除最后一个元素后组成的切片, 原始 slice 不会被改变。
func removeTail[T any](s []T) []T {
	if len(s) == 0 {
		return s
	}
	return s[:len(s)-1]
}

// Pop removes and returns the last element from a slice.
// 删除并返回最后一个元素和剩余元素组成的切片, 原始 slice 不会被改变。
// 如果 slice 为空，则返回一个零值和空切片。
func Pop[T any](s []T) (T, []T) {
	if len(s) == 0 {
		var zero T
		return zero, s
	}
	return s[len(s)-1], s[:len(s)-1]
}

// Drop returns the last n elements from a slice.
// 返回最后一个 n 个元素组成的切片, 原始 slice 不会被改变。
func Drop[T any](s []T, n int) []T {
	if n < 0 || n > len(s) {
		panic("wrong drop number")
	}
	return s[len(s)-n:]
}

// Remove removes and returns the first occurrence of an element from a slice.
// 删除第一个出现的元素后组成的切片, 原始 slice 会被改变。
func Remove[T comparable](s []T, element T) []T {
	for i, v := range s {
		if v == element {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

// Delete removes and returns the element at a specific index from a slice.
// 删除位于索引 index 的元素, 原始 slice 会被改变。
func Delete[T any](s []T, index int) []T {
	if index < 0 || index > len(s) {
		panic("index out of range")
	}
	return append(s[:index], s[index+1:]...)
}

// DeleteMany removes and returns the elements between start and end from a slice.
// 删除位于索引 start 和 end 之间的元素, 原始 slice 会被改变。
func DeleteMany[T any](s []T, start, end int) []T {
	if start < 0 || end > len(s) {
		panic("index out of range")
	}
	if end < start {
		panic("end must be greater than start")
	}

	return append(s[:start], s[end:]...)
}

// Append appends one or more elements to a slice.
// 追加一个或多个元素到 slice, 原始 slice 不会被改变。
func Append[T any](src []T, elements ...T) []T {
	return append(src, elements...)
}

// Extend extends a slice by appending one or more elements.
// 扩展 slice，追加一个新切片, 原始 slice 不会被改变。
func Extend[T any](src []T, elements []T) []T {
	return append(src, elements...)
}

// Insert inserts an element at a specific index in a slice.
// 在 slice 的指定位置插入一个元素, 原始 slice 不会被改变。
func Insert[T any](src []T, element T, index int) []T {
	if index < 0 || index > len(src) {
		panic("index out of range")
	}
	return append(src[:index], append([]T{element}, src[index:]...)...)
}

// AddV1  inserts an element at a specific index in a slice.
// 在 slice 的指定位置插入一个元素, 原始 slice 不会被改变。
func AddV1[T any](src []T, element T, index int) []T { // 性能比Insert好
	if index < 0 || index > len(src) {
		panic("index out of range")
	}
	var zeroValue T
	src = append(src, zeroValue)
	for i := len(src) - 1; i > index; i-- {
		if i-1 >= 0 {
			src[i] = src[i-1]
		}
	}

	src[index] = element

	return src
}

// AddV2  inserts an element at a specific index in a slice.
// 在 slice 的指定位置插入一个元素, 原始 slice 不会被改变。
func AddV2[T any](src []T, element T, index int) []T { // 性能比AddV1好
	if index < 0 || index > len(src) {
		panic("index out of range")
	}
	sn := len(src)
	var zeroValue T
	src = append(src, zeroValue)

	if index < sn {
		copy(src[index+1:], src[index:sn])
	}
	src[index] = element

	return src
}

// InsertMany inserts multiple elements at a specific index in a slice.
// 在 slice 的指定位置插入新切片, 原始 slice 不会被改变。
func InsertMany[T any](src []T, elements []T, index int) []T {
	if index < 0 || index > len(src) {
		panic("index out of range")
	}
	return append(src[:index], append(elements, src[index:]...)...)
}

// AddMany adds multiple elements to a slice at a specific index.
// 在 slice 的指定位置添加新切片, 原始 slice 不会被改变。
func AddMany[T any](src []T, elements []T, index int) []T { // 性能比InsertMany好
	if index < 0 || index > len(src) {
		panic("index out of range")
	}
	en := len(elements)
	sn := len(src)
	src = append(src, make([]T, en)...)

	// 移动元素
	if index < sn {
		copy(src[index+en:], src[index:sn])
	}

	// 插入新元素
	copy(src[index:index+en], elements)
	return src
}

// Reverse returns a new slice with the elements in reverse order.
// 返回一个新的 slice，元素顺序相反, 原始 slice 不会被改变。
func Reverse[T any](s []T) []T {
	length := len(s)
	result := make([]T, 0, length)

	for i := length - 1; i >= 0; i-- {
		result = append(result, s[i])
	}

	return result
}

// ReverseSelf reverses the elements in a slice in place.
// 反转 slice 本身, 原始 slice 不会被改变。
func ReverseSelf[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// Repeat returns a new slice with n copies of the given element.
// 返回一个新 slice，其中包含 n 个给定元素的副本, 原始 slice 不会被改变。
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
	for _, v := range src {
		if v == dst {
			return true
		}
	}

	return false
}

// ContainsAny 判断 src 里面是否存在 dst 中的任何一个元素
func ContainsAny[T comparable](src, dst []T) bool {
	srcMap := toMap(src)
	for _, value := range dst {
		if _, ok := srcMap[value]; ok {
			return true
		}
	}

	return false
}

// ContainsAll 判断 src 里面是否存在 dst 中的所有元素
func ContainsAll[T comparable](src, dst []T) bool {
	srcMap := toMap(src)
	for _, value := range dst {
		if _, ok := srcMap[value]; !ok {
			return false
		}
	}
	return true
}

// Unique removes duplicate elements from a slice.
// 删除 slice 中的重复元素, 原始 slice 不会被改变 。
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

func toMap[T comparable](s []T) map[T]struct{} {
	result := make(map[T]struct{}, len(s))
	for _, v := range s {
		result[v] = struct{}{}
	}
	return result
}

// DiffSet returns the difference between two slices.
// 返回 src 和 dst 的差集, 即 src 中存在，dst 中不存在的元素。
func DiffSet[T comparable](src, dst []T) []T {
	srcMap := toMap(src)
	for _, value := range dst {
		delete(srcMap, value)
	}

	result := make([]T, 0, len(srcMap))

	for _, value := range src {
		if _, ok := srcMap[value]; ok {
			result = append(result, value)
		}
	}

	return result
}

// IntersectSet returns the intersection between two slices.
// 返回 src 和 dst 的交集, 即 src 中存在，dst 中也存在元素。
func IntersectSet[T comparable](src, dst []T) []T {
	dstMap := toMap(dst)

	result := make([]T, 0, len(src))

	for _, value := range src {
		if _, ok := dstMap[value]; ok {
			result = append(result, value)
		}
	}
	return result
}

// UnionSet returns the union of two slices.
// 返回 src 和 dst 的并集, 即 src 和 dst 的并集。
func UnionSet[T comparable](src, dst []T) []T {
	result := append(src, dst...)
	return Unique(result)
}

// Index 返回和 dst 相等的第一个元素下标
// -1 表示没找到
func Index[T comparable](src []T, dst T) int {
	for i, v := range src {
		if v == dst {
			return i
		}
	}
	return -1
}

// LastIndex 返回和 dst 相等的最后一个元素下标
// -1 表示没找到
func LastIndex[T comparable](src []T, dst T) int {
	for i := len(src) - 1; i >= 0; i-- {
		if src[i] == dst {
			return i
		}
	}
	return -1
}

// IndexAll 返回和 dst 相等的所有元素的下标
func IndexAll[T comparable](src []T, dst T) []int {
	result := make([]int, 0, len(src))
	for i, v := range src {
		if v == dst {
			result = append(result, i)
		}
	}
	return result
}
