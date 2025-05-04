package itertools

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
	var result []T
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

func Reduce[T any](s []T, fn func(T, T) T) T {
	result := s[0]
	for _, v := range s[1:] {
		result = fn(result, v)
	}
	return result
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

// Reverse returns a new slice with the elements in reverse order.
func Reverse[T any](s []T) []T {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
