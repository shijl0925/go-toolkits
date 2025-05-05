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
	result := make([]T, 0, len(s))
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
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

// GroupBy groups elements of the slice into a map based on a key function
func GroupBy[T any, U comparable](slice []T, fn func(T) U) map[U][]T {
	result := make(map[U][]T)
	for _, v := range slice {
		key := fn(v)
		result[key] = append(result[key], v)
	}
	return result
}
