package slice

import "golang.org/x/exp/constraints"

type RealNumber interface {
	constraints.Integer | constraints.Float
}

type Number interface {
	RealNumber | ~complex64 | ~complex128
}

func Sum[T RealNumber](s []T) T {
	var result T
	for _, v := range s {
		result += v
	}
	return result
}

func Max[T constraints.Ordered](s []T) T {
	result := s[0]
	for _, v := range s[1:] {
		if v > result {
			result = v
		}
	}
	return result
}

func Min[T constraints.Ordered](s []T) T {
	result := s[0]
	for _, v := range s[1:] {
		if v < result {
			result = v
		}
	}
	return result
}
