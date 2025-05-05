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
	for _, v := range s[1:] {
		if v > result {
			result = v
		}
	}
	return result
}

// Min returns the minimum value in the slice.
func Min[T constraints.Ordered](s []T) T {
	result := s[0]
	for _, v := range s[1:] {
		if v < result {
			result = v
		}
	}
	return result
}

// Insert inserts an element at a specific index in a slice.
func Insert[T any](src []T, element T, index int) []T {
	return append(src[:index], append([]T{element}, src[index:]...)...)
}
