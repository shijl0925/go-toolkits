package slicex

import (
	"fmt"
	"github.com/shijl0925/go-toolkits"
	"golang.org/x/exp/constraints"
	"reflect"
	"strings"
)

type RealNumber interface {
	constraints.Integer | constraints.Float
}

type Number interface {
	RealNumber | ~complex64 | ~complex128
}

// Sum returns the sum of all elements in the slice.
//
// Example:
//
//	s := []int{1, 2, 3, 4, 5}
//	r := Sum(s) //  r == 15
func Sum[T RealNumber](s []T) T {
	var result T
	for _, v := range s {
		result += v
	}
	return result
}

// Max returns the maximum value in the slice.
//
// Example:
//
//	s := []int{1, 2, 3, 4, 5}
//	r := Max(s) //  r == 5
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
//
// Example:
//
//	s := []int{1, 2, 3, 4, 5}
//	r := Min(s) //  r == 1
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
// 返回一个新的切片，其中去除了原切片的第一个元素, 原始 slice 不会被改变。
// 如果原切片为空，则返回原切片本身。
//
// Example:
//
//	s := []int{1, 2, 3, 4, 5}
//	r := removeHead(s) //  r == []int{2, 3, 4, 5}
func removeHead[T any](s []T) []T {
	if len(s) == 0 {
		return s
	}
	return s[1:]
}

// removeTail removes the last element from a slice and returns the remaining elements.
// 返回一个新的切片, 其中去除了原切片的最后一个元素, 原始 slice 不会被改变。
// 如果原切片为空，则返回原切片本身。
//
// Example:
//
//	s := []int{1, 2, 3, 4, 5}
//	r := removeTail(s) //  r == []int{1, 2, 3, 4}
func removeTail[T any](s []T) []T {
	if len(s) == 0 {
		return s
	}
	return s[:len(s)-1]
}

// Pop removes and returns the last element from a slice.
// 从切片 s 中弹出最后一个元素，并返回该元素和弹出状态。
// 如果 slice 为空，则返回一个零值和false。
// 注意：原始 slice 会被改变。
//
// Example:
//
//	s := []int{1, 2, 3, 4, 5}
//	r, ok := Pop(s) //  r == 5, ok ==true, s == []int{1, 2, 3, 4}
func Pop[T any](s []T) (T, bool) {
	if len(s) == 0 {
		var zero T
		return zero, false
	}
	v := s[len(s)-1]
	s = removeTail(s)
	return v, true
}

// Drop removes the n elements from the slice and returns the remaining elements.
// 删除并返回切片的前 n 个元素和剩余元素组成的切片, 原始 slice 不会被改变。
// 注意: 返回的切片仍与原切片共享底层数组。
// Example:
//
//	s := []int{1, 2, 3, 4, 5}
//	r, _ := Drop(s, 2) //  r == []int{1, 2, 3}
func Drop[T any](s []T, n int) ([]T, error) {
	if n < 0 || n > len(s) {
		return nil, fmt.Errorf("invalid drop number: %d, length of slice is %d", n, len(s))
	}
	return s[:len(s)-n], nil
}

// DropLeft removes the n elements from the slice and returns the remaining elements.
// 删除并返回切片的左侧 n 个元素和剩余元素组成的切片, 原始 slice 不会被改变。
// 注意: 返回的切片仍与原切片共享底层数组。
// Example:
//
//	s := []int{1, 2, 3, 4, 5}
//	r, _ := DropLeft(s, 2) //  r == []int{3, 4, 5}
func DropLeft[T any](s []T, n int) ([]T, error) {
	if n < 0 || n > len(s) {
		return nil, fmt.Errorf("invalid drop number: %d, length of slice is %d", n, len(s))
	}
	return s[n:], nil
}

// Remove returns a slice with the all occurrence of element removed.
// 返回一个切片，该切片移除了所有等于 element 的元素。
// 注意: 该函数不会修改原始切片内容，而是返回一个新的切片。
// Example:
//
//	s := []int{1, 2, 3, 2}
//	r := Remove(s, 2) // r == []int{1, 3}
func Remove[T comparable](s []T, element T) []T {
	result := make([]T, 0, len(s))
	for _, v := range s {
		if v != element {
			result = append(result, v)
		}
	}
	return result
}

// Delete removes and returns the element at a specific index from a slice.
// 删除位于索引 index 的元素, 原始 slice 会被改变。
// 注意: 返回的切片仍与原切片共享底层数组。
//
// Example:
//
//	s := []int{1, 2, 3, 2}
//	r, _ := Delete(s, 2) // r == []int{1, 2, 2}
func Delete[T any](s []T, index int) ([]T, bool) {
	if index < 0 || index >= len(s) {
		return s, false
	}
	return append(s[:index], s[index+1:]...), true
}

// DeleteAt removes and returns the element at a specific index from a slice.
// 删除位于索引 index 的元素, 原始 slice 会被改变。
// 注意: 返回的切片仍与原切片共享底层数组。
//
// Example:
//
//	s := []int{1, 2, 3, 2}
//	r, _ := DeleteAt(s, 2) // r == []int{1, 2, 2}
func DeleteAt[T any](s []T, index int) ([]T, bool) { // 性能比Delete好
	length := len(s)
	if index < 0 || index >= length {
		return s, false
	}

	//从index位置开始，后面的元素依次往前挪1个位置
	for i := index; i < length-1; i++ {
		s[i] = s[i+1]
	}

	return s[:length-1], true
}

// DeleteMany removes and returns the elements between start and end from a slice.
// 删除位于索引 start 和 end 之间的元素, 原始 slice 会被改变。
// 注意: 返回的切片仍与原切片共享底层数组。
//
// Example:
//
//	s := []int{1, 2, 3, 4, 5, 6}
//	r, _ := DeleteMany(s, 2, 4) // r == []int{1, 2, 5, 6}
func DeleteMany[T any](s []T, start, end int) ([]T, bool) {
	if len(s) == 0 {
		return s, false
	}
	if start < 0 || end > len(s) {
		return s, false
	}
	if end < start {
		return s, false
	}

	return append(s[:start], s[end:]...), true
}

// Append appends one or more elements to a slice.
// 追加一个或多个元素到 slice。
// 注意: 该函数不会修改原始切片内容，而是返回一个新的切片。
//
// Example:
//
//	s := []int{1, 2, 3}
//	r := Append(s, 4, 5) // r == []int{1, 2, 3, 4, 5}
func Append[T any](src []T, elements ...T) []T {
	return append(src, elements...)
}

// Extend extends a slice by appending one or more elements.
// 将切片 elements 中的所有元素追加到 src 切片后并返回新切片。
// 注意: 该函数不会修改原始切片内容，而是返回一个新的切片。
//
// Example:
//
//	s := []int{1, 2, 3}
//	r := Extend(s, []int{4, 5}) // r == []int{1, 2, 3, 4, 5}
func Extend[T any](src []T, elements []T) []T {
	return append(src, elements...)
}

// Insert inserts an element at a specific index in a slice.
// 在 slice 的指定位置插入一个元素, 原始 slice 不会被改变。
// 注意: 该函数不会修改原始切片内容，而是返回一个新的切片。
//
// Example:
//
//	s := []int{1, 2, 3}
//	r, _ := Insert(s, 4, 1) // r == []int{1, 4, 2, 3}
func Insert[T any](src []T, element T, index int) ([]T, error) {
	if index < 0 || index > len(src) {
		return src, fmt.Errorf("index out of range: %d", index)
	}
	return append(src[:index], append([]T{element}, src[index:]...)...), nil
}

// InsertAtV1  inserts an element at a specific index in a slice.
// 在 slice 的指定位置插入一个元素, 原始 slice 不会被改变。
// 注意: 该函数不会修改原始切片内容，而是返回一个新的切片。
//
// Example:
//
//	s := []int{1, 2, 3}
//	r, _ := InsertAtV1(s, 4, 1) // r == []int{1, 4, 2, 3}
func InsertAtV1[T any](src []T, element T, index int) ([]T, error) { // 性能比Insert好
	if index < 0 || index > len(src) {
		return src, fmt.Errorf("index out of range: %d", index)
	}
	var zeroValue T
	src = append(src, zeroValue)
	for i := len(src) - 1; i > index; i-- {
		if i-1 >= 0 {
			src[i] = src[i-1]
		}
	}

	src[index] = element

	return src, nil
}

// InsertAtV2  inserts an element at a specific index in a slice.
// 在 slice 的指定位置插入一个元素, 原始 slice 不会被改变。
// 注意: 该函数不会修改原始切片内容，而是返回一个新的切片。
//
// Example:
//
//	s := []int{1, 2, 3, 4, 5}
//	r, _ := InsertAtV2(slice, 6, 2)  // r == [1, 2, 6, 3, 4, 5]
func InsertAtV2[T any](src []T, element T, index int) ([]T, error) { // 性能比InsertAtV1好
	if index < 0 || index > len(src) {
		return src, fmt.Errorf("index out of range: %d", index)
	}

	var zeroValue T
	src = append(src, zeroValue)

	if index < len(src)-1 {
		copy(src[index+1:], src[index:len(src)-1])
	}
	src[index] = element

	return src, nil
}

// InsertMany inserts multiple elements at a specific index in a slice.
// 在 slice 的指定位置插入新切片, 原始 slice 不会被改变。
// 注意: 该函数不会修改原始切片内容，而是返回一个新的切片。
//
// Example:
//
//	s := []int{1, 2, 3}
//	r, _ := InsertMany([]int{1, 2, 3}, []int{4, 5}, 1)  // r == [1, 4, 5, 2, 3]
func InsertMany[T any](src []T, elements []T, index int) ([]T, error) {
	if index < 0 || index > len(src) {
		return src, fmt.Errorf("index out of range: %d", index)
	}

	return append(src[:index], append(elements, src[index:]...)...), nil
}

// AddMany adds multiple elements to a slice at a specific index.
// 在 slice 的指定位置添加新切片, 原始 slice 不会被改变。
// 注意: 该函数不会修改原始切片内容，而是返回一个新的切片。
//
// Example:
//
// s := []int{1, 2, 3}
// r, _ := AddMany([]int{1, 2, 3}, []int{4, 5}, 1)  // r == [1, 4, 5, 2, 3]
func AddMany[T any](src []T, elements []T, index int) ([]T, error) { // 性能比InsertMany好
	if index < 0 || index > len(src) {
		return src, fmt.Errorf("index out of range: %d", index)
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
	return src, nil
}

// Reverse returns a new slice with the elements in reverse order.
// 返回一个新的 slice，元素顺序相反, 原始 slice 不会被改变。
// 注意: 该函数不会修改原始切片内容，而是返回一个新的切片。
//
// Example:
//
// s := []int{1, 2, 3, 4, 5}
// r := Reverse(s)  // r == []int{5, 4, 3, 2, 1}
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
//
// Example:
//
// s := []int{1, 2, 3, 4, 5}
// ReverseSelf(s)  // s == []int{5, 4, 3, 2, 1}
func ReverseSelf[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// Repeat returns a new slice with n copies of the given element.
// 返回一个新 slice，其中包含 n 个给定元素的副本, 原始 slice 不会被改变。
// 注意: 返回的切片仍与原切片共享底层数组。
//
// Example:
//
// r0 := Repeat(5, 0)  // r0 == []int{}
// r1 := Repeat(5, 3)  // r1 == []int{5, 5, 5}
// r2 := Repeat([]int{1, 2, 3}, 3)  // r2 == [][]int{{1, 2, 3}, {1, 2, 3}, {1, 2, 3}}
// r3 := Repeat("one", 3)  // r3 == []string{"one", "one", "one"}
func Repeat[T any](s T, n int) []T {
	result := make([]T, 0, n)
	for i := 0; i < n; i++ {
		result = append(result, s)
	}
	return result
}

// Product returns a new slice of tuples, where each tuple contains the corresponding elements from two slices.
//
// Example:
//
// s1 := []int{1, 2}
// s2 := []int{3, 4}
// r := Product(s1, s2) // r == [][]int{{1, 3}, {1, 4}, {2, 3}, {2, 4}}
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
//
// Example:
//
//	r := Contains([]int{1, 2, 3}, 2) //  r == true
func Contains[T comparable](src []T, dst T) bool {
	for _, v := range src {
		if v == dst {
			return true
		}
	}

	return false
}

// ContainsAny 判断 src 里面是否存在 dst 中的任何一个元素
//
// Example:
//
//	r := ContainsAny([]int{1, 2, 3}, []int{2, 4}) //  r == true
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
//
// Example:
//
//	r := ContainsAll([]int{1, 2, 3}, []int{2, 3}) //  r == true
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
//
// Example:
//
//	s := []int{1, 2, 3, 1, 2, 4}
//	r := Unique(s) // r == []int{1, 2, 3, 4}
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

// 将 slice 转换为 map
//
// Example:
//
//	s := []int{1, 2, 3, 1, 2, 4}
//	r := ToMap(s) // r == map[int]struct{}{1: {}, 2: {}, 3: {}, 4: {}}
func toMap[T comparable](s []T) map[T]struct{} {
	result := make(map[T]struct{}, len(s))
	for _, v := range s {
		result[v] = struct{}{}
	}
	return result
}

// DiffSet returns the difference between two slices.
// 返回 src 和 dst 的差集, 即 src 中存在，dst 中不存在的元素。
//
// Example:
//
//	src := []int{1, 2, 3, 4, 5}
//	dst := []int{1, 2, 3}
//	r := DiffSet(src, dst) // r == []int{4, 5}
func DiffSet[T comparable](src, dst []T) []T {
	dstMap := toMap(dst)

	result := make([]T, 0, len(src))

	for _, value := range src {
		if _, ok := dstMap[value]; !ok {
			result = append(result, value)
		}
	}

	return result
}

// IntersectSet returns the intersection between two slices.
// 返回 src 和 dst 的交集, 即 src 中存在，dst 中也存在元素。
//
// Example:
//
//	src := []int{1, 2, 3, 4, 5}
//	dst := []int{4, 5, 6, 7, 8}
//	r := IntersectSet(src, dst) // r == []int{4, 5}
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
// 返回 src 和 dst 的并集, 即 src 和 dst中所有元素的集合。
//
// Example:
//
//	src := []int{1, 2, 3, 4, 5}
//	dst := []int{4, 5, 6, 7, 8}
//	r := UnionSet(src, dst) // r == []int{1, 2, 3, 4, 5, 6, 7, 8}
func UnionSet[T comparable](src, dst []T) []T {
	result := append(src, dst...)
	return Unique(result)
}

// Index 返回和 dst 相等的第一个元素下标
// -1 表示没找到
//
// Example:
//
//	src := []int{1, 2, 3, 2, 5}
//	r := Index(src, 2) // r == 1
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
//
// Example:
//
//	src := []int{1, 2, 3, 2, 5}
//	r := LastIndex(src, 2)  // r == 3
func LastIndex[T comparable](src []T, dst T) int {
	for i := len(src) - 1; i >= 0; i-- {
		if src[i] == dst {
			return i
		}
	}
	return -1
}

// IndexAll 返回和 dst 相等的所有元素的下标
//
// Example:
//
//	src := []int{1, 2, 3, 2, 5}
//	r := IndexAll(src, 2)  // r == []int{1, 3}
func IndexAll[T comparable](src []T, dst T) []int {
	result := make([]int, 0, len(src))
	for i, v := range src {
		if v == dst {
			result = append(result, i)
		}
	}
	return result
}

// ShallowClone returns a new slice with the same elements as the original slice.
// 返回一个新的 slice，其中包含原始 slice 的副本, 原始 slice 不会被改变。
// 注意：本函数不实现真正的深拷贝，适用于元素为不可变或非指针类型的场景。
// 若 T 包含指针或引用类型，原切片与新切片将共享这些引用。
// func ShallowClone[T any](s []T) []T {
// 	t := make([]T, 0, len(s))
// 	t = append(t, s...)
// 	return t
// }

// JoinSlice join []any slice to string.
//
// Example:
//
//	s := []string{"a", "b", "c"}
//	s1 := JoinSlice(",", s) // s1 == "a,b,c"
//	t := []int{1, 2, 3}
//	t1 := JoinSlice(",", t) // t1 == "1,2,3"
func JoinSlice[T any](sep string, s []T) string {
	var bf strings.Builder
	for i, v := range s {
		if i > 0 {
			bf.WriteString(sep)
		}
		s, _ := toolkits.SafeToString(v)
		bf.WriteString(s)
	}
	return bf.String()
}
// Chunk creates a slice of elements split into groups the length of size.
// 将切片按照给定的大小进行分组，并返回一个二维切片。
func Chunk[T any](s []T, size int) [][]T {
	result := [][]T{}

	if len(s) == 0 || size <= 0 {
		return result
	}

	for i := 0; i < len(s); i += size {
		end := i + size
		if end > len(s) {
			end = len(s)
		}
		result = append(result, s[i:end])
	}

	return result
}

// Compact returns a new slice with all falsey values removed. The values false, nil, 0, and "" are falsey.
// 返回一个新的 slice，其中包含原始 slice 中所有非零值的元素。注意：本函数会改变原始 slice。
func Compact[T comparable](s []T) []T {
	var zero T
	i := 0
	for _, item := range s {
		if item != zero {
			s[i] = item
			i++
		}
	}
	return s[:i]
}

// Concat returns a new slice with all slices concatenated.
func Concat[T any](s ...[]T) []T {
	var totalLen int
	for _, slice := range s {
		totalLen += len(slice)
	}

	result := make([]T, 0, totalLen)

	for _, item := range s {
		result = append(result, item...)
	}
	return result
}

// Equal returns true if the two slices are equal.
// 判断两个 slice 是否相等。
func Equal[T comparable](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

// EqualUnordered checks if two slices are equal: the same length and all elements' value are equal (unordered).
// 判断两个 slice 是否相等：长度相同且所有元素值相等（不考虑顺序）。
//
// Example:
//
//	s1 := []int{1, 2, 3}
//	s2 := []int{3, 2, 1}
//	EqualUnordered(s1, s2) // true
func EqualUnordered[T comparable](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}

	seen := make(map[T]int, len(s1))
	for _, v := range s1 {
		seen[v]++
	}

	for _, v := range s2 {
		if seen[v] == 0 {
			return false
		}
		seen[v]--
	}

	return true
}

// Counter returns a map of the counts of each value in the slice.
// 返回一个 map，其中键为 slice 中的每个值，值为该值在 slice 中出现的次数。
//
// Example:
//
// s := []string{"a", "b", "a"}
// r := Counter(s) // r == map[string]int{"a": 2, "b": 1}
func Counter[T comparable](s []T) map[T]int {
	if len(s) == 0 {
		return map[T]int{}
	}

	result := make(map[T]int)
	for _, v := range s {
		result[v]++
	}
	return result
}

// Replace returns a new slice with all occurrences of old replaced by new.
// 返回一个新的 slice，其中包含原始 slice 中所有 old 值被 new 替换后的元素。n 表示替换的次数，-1 表示替换所有。
func Replace[T comparable](s []T, old, new T, n int) []T {
	result := make([]T, len(s))
	copy(result, s)

	for i := range result {
		if result[i] == old && n != 0 {
			result[i] = new
			n--
		}
	}
	return result
}

// IsAscending checks if a slice is ascending order.
// 判断一个 slice 是否是升序排列的。
func IsAscending[T constraints.Ordered](s []T) bool {
	for i := 1; i < len(s); i++ {
		if s[i] < s[i-1] {
			return false
		}
	}
	return true
}

// IsDescending checks if a slice is descending order.
// 判断一个 slice 是否是降序排列的。
func IsDescending[T constraints.Ordered](s []T) bool {
	for i := 1; i < len(s); i++ {
		if s[i] > s[i-1] {
			return false
		}
	}
	return true
}

// RightPadding adds padding to the right end of a slice.
func RightPadding[T any](s []T, v T, n int) []T {
	if n == 0 {
		return s
	}

	result := make([]T, len(s)+n)
	copy(result, s)

	for i := len(s); i < len(result); i++ {
		result[i] = v
	}
	return result
}

// LeftPadding adds padding to the left begin of a slice.
func LeftPadding[T any](s []T, v T, n int) []T {
	if n == 0 {
		return s
	}
	result := make([]T, len(s)+n)
	copy(result[n:], s)

	for i := 0; i < n; i++ {
		result[i] = v
	}

	return result
}

func flattenRecursive(value reflect.Value, result reflect.Value) reflect.Value {
	for i := 0; i < value.Len(); i++ {
		item := value.Index(i)
		kind := item.Kind()

		if kind == reflect.Slice {
			result = flattenRecursive(item, result)
		} else {
			result = reflect.Append(result, item)
		}
	}

	return result
}

// FlattenDeep flattens slice recursive.
// 递归地展平切片
//
// Example:
//
// s := [][]int{{1, 2}, {3, 4}, {5, 6}}
// r := FlattenDeep(s)   // r == []int{1, 2, 3, 4, 5, 6}
func FlattenDeep(slice any) any {
	rv := reflect.ValueOf(slice)

	// 获取切片的最里层元素类型
	rt := rv.Type()
	for {
		if rt.Kind() != reflect.Slice {
			break
		}
		rt = rt.Elem()
	}

	tmp := reflect.MakeSlice(reflect.SliceOf(rt), 0, 0)
	result := flattenRecursive(rv, tmp)

	return result.Interface()
}
