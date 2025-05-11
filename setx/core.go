package setx

type Set[T comparable] map[T]struct{}

// NewSet creates a new set from a list of elements.
// 构造函数
func NewSet[T comparable](size int, a ...T) Set[T] {
	s := make(Set[T], size)
	for _, v := range a {
		s.Add(v)
	}
	return s
}

// Add adds an element to the set.
func (s Set[T]) Add(element T) {
	s[element] = struct{}{}
}

// Remove removes an element from the set.
func (s Set[T]) Remove(element T) {
	delete(s, element)
}

// Len returns the number of elements in the set.
func (s Set[T]) Len() int {
	return len(s)
}

// Exists returns true if the element exists in the set.
func (s Set[T]) Exists(element T) bool {
	_, ok := s[element]
	return ok
}

// Keys returns the keys of the set.
// 返回一个切片，包含set中的所有元素。
// 顺序是不确定的
func (s Set[T]) Keys() []T {
	result := make([]T, 0, s.Len())
	for element := range s {
		result = append(result, element)
	}
	return result
}

// DiffSet returns the difference between two sets.
// 返回src集合与dst集合的差集, 即 src 中存在，dst 中不存在的元素。
// 并且返回值的顺序是不确定的
func DiffSet[T comparable](src, dst Set[T]) []T {
	result := make([]T, 0, src.Len())

	for element := range src {
		if !dst.Exists(element) {
			result = append(result, element)
		}
	}

	return result
}

// IntersectSet returns the intersection between two sets.
// 返回src集合与dst集合的交集, 即 src 和 dst 都存在的元素。
// 并且返回值的顺序是不确定的
func IntersectSet[T comparable](src, dst Set[T]) []T {
	result := make([]T, 0, src.Len())

	for element := range src {
		if dst.Exists(element) {
			result = append(result, element)
		}
	}

	return result
}

// UnionSet returns the union between two sets.
// 返回src集合与dst集合的并集, 即 src 和 dst中所有元素的集合。
// 并且返回值的顺序是不确定的
func UnionSet[T comparable](src, dst Set[T]) []T {
	for element := range src {
		if !dst.Exists(element) {
			dst.Add(element)
		}
	}

	result := make([]T, 0, dst.Len())
	for element := range dst {
		result = append(result, element)
	}

	return result
}
