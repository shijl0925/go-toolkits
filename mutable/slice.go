package mutable

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
// ReverseSelf(s)  // s == []int{5, 4, 3, 2, 1}
func Reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
