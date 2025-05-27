package itertools

import "github.com/shijl0925/go-toolkits/internal/constraints"

// Map applies a function to each element in a slice and returns a new slice with the results.
func Map[T any, U any](s []T, fn func(T) U) []U {
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}

// Filter applies a function to each element in a slice and returns a new slice with the elements that satisfy the predicate.
func Filter[T any](s []T, fn func(T) bool) []T {
	result := make([]T, 0, len(s))
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// DropWhile remove elements from the left while the predicate is true.
// 丢弃原有切片中直到函数返回false之前的所有元素，然后返回后面所有元素。
// Example:
//
//	s := []int{1, 2, 3, 4, 5}
//	r := DropWhile(s, func(i int) bool { return > 3 }) // r == []int{4, 5}
func DropWhile[T any](s []T, fn func(T) bool) []T {
	i := 0
	for ; i < len(s); i++ {
		if !fn(s[i]) {
			break
		}
	}
	return s[i:]
}

func Reduce[T any, U any](s []T, fn func(U, T) U, initial U) U {
	acc := initial
	for _, v := range s {
		acc = fn(acc, v)
	}
	return acc
}

type Tuple[T any, U any] struct {
	A T
	B U
}

// Zip returns a new slice of tuples, where each tuple contains the corresponding elements from two slices.
func Zip[T any, U any](a []T, b []U) []Tuple[T, U] {
	length := len(a)
	if len(b) < length {
		length = len(b)
	}

	result := make([]Tuple[T, U], length)

	for i := 0; i < length; i++ {
		result[i] = Tuple[T, U]{a[i], b[i]}
	}

	return result
}

// CombineToMap 将两个切片 keys 和 values 合并为一个 map。
// 键类型 K 必须是 comparable，值类型 V 可以是任意类型。
// 注意：若 keys 中存在重复键，则后面的值将覆盖前面的值,
// 若values数量大于 keys 的长度，则多余的键值对将被忽略。
func CombineToMap[K comparable, V any](keys []K, values []V) map[K]V {
	length := len(keys)
	if len(values) < length {
		length = len(values)
	}
	result := make(map[K]V, length)
	for i := 0; i < length; i++ {
		result[keys[i]] = values[i]
	}
	return result
}

// All returns true if all elements in a slice satisfy the predicate.
func All[T any](s []T, fn func(T) bool) bool {
	for _, v := range s {
		if !fn(v) {
			return false
		}
	}
	return true
}

// Any returns true if any element in a slice satisfies the predicate.
func Any[T any](s []T, fn func(T) bool) bool {
	for _, v := range s {
		if fn(v) {
			return true
		}
	}
	return false
}

// GroupBy groups elements of the slice into a map based on a key function
func GroupBy[T any, U comparable](slice []T, fn func(T) U) map[U][]T {
	result := make(map[U][]T)
	for _, v := range slice {
		key := fn(v)
		if _, ok := result[key]; !ok {
			result[key] = []T{}
		}
		result[key] = append(result[key], v)
	}
	return result
}

//nolint:typecheck
func Range[T constraints.Integer | constraints.Float](start, end, step T) []T {
	var result []T

	if step == 0 {
		return []T{}
	}

	// 判断是否递增序列
	ascending := step > 0

	// 判断是否满足生成条件
	if ascending && start >= end {
		return []T{}
	}
	if !ascending && start <= end {
		return []T{}
	}

	for i := start; (ascending && i < end) || (!ascending && i > end); i += step {
		result = append(result, i)
	}

	return result
}
