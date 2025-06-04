package setx

import "fmt"

type Set[T comparable] map[T]struct{}

// New creates a new set from a list of elements.
// 构造函数
//
// Example:
//
// s := setx.New(1, 2, 3)
// // or s:= setx.New([]int{1, 2, 3})
// fmt.Println(s) // map[1:{} 2:{} 3:{}]
func New[T comparable](elements ...interface{}) Set[T] {
	s := make(Set[T])
	for _, v := range elements {
		switch val := v.(type) {
		case T:
			s.Add(val)
		case []T:
			s.Add(val...)
		default:
			fmt.Printf("setx.New: invalid type")
			continue
		}
	}
	return s
}

// Add adds an element to the set.
//
// Example:
//
// s := setx.New(1, 2, 3)
// s.Add(4)
// fmt.Println(s) // map[1:{} 2:{} 3:{} 4:{}]
func (s Set[T]) Add(elements ...T) {
	for _, element := range elements {
		s[element] = struct{}{}
	}
}

// Remove removes an element from the set.
//
// Example:
//
// s := setx.New(1, 2, 3)
// s.Remove(3)
// fmt.Println(s) // map[1:{} 2:{}]
func (s Set[T]) Remove(element T) {
	delete(s, element)
}

// Len returns the number of elements in the set.
//
// Example:
//
// s := setx.New(1, 2, 3)
// fmt.Println(s.Len()) // 3
func (s Set[T]) Len() int {
	return len(s)
}

// Exists returns true if the element exists in the set.
//
// Example:
//
// s := setx.New(1, 2, 3)
// fmt.Println(s.Exists(3)) // true
// fmt.Println(s.Exists(0)) // false
func (s Set[T]) Exists(element T) bool {
	_, ok := s[element]
	return ok
}

// Keys returns the keys of the set.
// 返回一个切片，包含set中的所有元素。
// 顺序是不确定的
//
// Example:
//
// s := setx.New(1, 2, 3)
// fmt.Println(s.Keys()) // [1 2 3] or [3 1 2] or [2 3 1]
func (s Set[T]) Keys() []T {
	result := make([]T, 0, s.Len())
	for element := range s {
		result = append(result, element)
	}
	return result
}

// Clear removes all elements from the set.
// 清空集合
//
// Example:
// s := setx.New([]int{1, 2, 3}...)
// s.Clear()
// fmt.Println(s.Len()) // 0
func (s Set[T]) Clear() {
	for element := range s {
		delete(s, element)
	}
}

// Equal returns true if the set is equal to the other set.
// 判断两个集合是否相等，如果两个集合的元素数量和元素值都相同，则返回true，否则返回false。
//
// Example:
//
// s := setx.New([]int{1, 2, 3}...)
// t := setx.New([]int{1, 2, 3}...)
// fmt.Println(s.Equal(t)) // true
// fmt.Println(s.Equal(setx.New([]int{1,2}...))) // false
func (s Set[T]) Equal(dst Set[T]) bool {
	if s.Len() != dst.Len() {
		return false
	}

	for element := range s {
		if !dst.Exists(element) {
			return false
		}
	}

	return true
}

// Iterate call function by every element of set
func (s Set[T]) Iterate(fn func(item T)) {
	for element := range s {
		fn(element)
	}
}

// ToSlice returns a slice containing all values of the set.
func (s Set[T]) ToSlice() []T {
	if s.Len() == 0 {
		return []T{}
	}

	result := make([]T, 0, s.Len())
	s.Iterate(func(item T) {
		result = append(result, item)
	})

	return result
}

// IsEmpty checks the set is empty or not
func (s Set[T]) IsEmpty() bool {
	return s.Len() == 0
}

// Pop delete one element of set then return it, if set is empty, return nil-value of T and false.
func (s Set[T]) Pop() (v T, ok bool) {
	if s.Len() > 0 {
		for element := range s {
			v = element
			delete(s, element)
			return v, true
		}
	}
	return v, false
}

// Update updates the set with the elements of the other set.
// 将dst集合中的元素添加到src集合中，如果src集合中已经存在该元素，则不添加。
// 并且返回值是更新后的src集合
//
// Example:
//
// s := setx.New([]int{1, 2, 3}...)
// t := setx.New([]int{1, 2, 4}...)
// fmt.Println(s.Update(t)) // map[1:{} 2:{} 3:{} 4:{}]
// fmt.Println(s) // map[1:{} 2:{} 3:{} 4:{}]
// fmt.Println(t) // map[1:{} 2:{} 4:{}]
func (s Set[T]) Update(dst Set[T]) Set[T] {
	for element := range dst {
		s[element] = struct{}{}
	}
	return s
}

// Difference returns the difference between two sets.
// 返回src集合与dst集合的差集, 即 src 中存在，dst 中不存在的元素。
// 并且返回值的顺序是不确定的
//
// Example:
//
// s := setx.New(1, 2, 3)
// t := setx.New(1, 2, 4)
// fmt.Println(setx.Difference(s, t)) // [3]
func Difference[T comparable](src, dst Set[T]) []T {
	result := make([]T, 0, src.Len())

	for element := range src {
		if !dst.Exists(element) {
			result = append(result, element)
		}
	}

	return result
}

// Intersect returns the intersection between two sets.
// 返回src集合与dst集合的交集, 即 src 和 dst 都存在的元素。
// 并且返回值的顺序是不确定的
//
// Example:
//
// s := setx.New(1, 2, 3)
// t := setx.New(1, 2, 4)
// fmt.Println(setx.Intersect(s, t)) // [1 2]
func Intersect[T comparable](src, dst Set[T]) []T {
	if src.Len() == 0 || dst.Len() == 0 {
		return []T{}
	}

	var iterSet, checkSet Set[T]
	// Iterate over the smaller set for efficiency
	if src.Len() <= dst.Len() {
		iterSet = src
		checkSet = dst
	} else {
		iterSet = dst
		checkSet = src
	}

	// The capacity of the result slice can be at most the size of the smaller set.
	result := make([]T, 0, iterSet.Len())

	for element := range iterSet {
		if checkSet.Exists(element) {
			result = append(result, element)
		}
	}

	return result
}

// Union returns the union between two sets.
// 返回src集合与dst集合的并集, 即 src 和 dst中所有元素的集合。
// 并且返回值的顺序是不确定的
//
// Example:
//
// s := setx.New(1, 2, 3)
// t := setx.New(1, 2, 4)
// fmt.Println(setx.Union(s, t)) // [1 2 3 4]
func Union[T comparable](src, dst Set[T]) []T {
	// Create a new set to store the union to avoid modifying input sets.
	// Estimate initial capacity for the union set.
	// A reasonable estimate is the size of the larger set, or sum if many are unique.
	// New() already handles Add, which is optimized.
	initialCap := src.Len()
	if dst.Len() > initialCap {
		initialCap = dst.Len()
	}
	// If one set is empty, the union is the other set.
	// NewFromSlice can be used if we convert the other set to slice first.
	// Or, more simply, create a new set and add elements.

	unionSet := make(Set[T], initialCap) // Preallocate based on larger set

	for element := range src {
		unionSet[element] = struct{}{}
	}

	for element := range dst {
		unionSet[element] = struct{}{}
	}

	// Convert the unionSet to a slice. Keys() preallocates correctly.
	return unionSet.Keys()
}
