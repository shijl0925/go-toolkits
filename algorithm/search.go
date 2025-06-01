//nolint:typecheck
package algorithm

import "github.com/shijl0925/go-toolkits/internal/constraints"

// BinarySearch return the index of target within a sorted slice, use binary search (recursive call itself).
// If not found return -1.
// 二分查找
func BinarySearch[T constraints.Ordered](s []T, target T, l, h int) int {
	if l > h || len(s) == 0 {
		return -1
	}
	m := l + (h-l)/2
	if s[m] > target {
		return BinarySearch(s, target, l, m-1)
	} else if s[m] < target {
		return BinarySearch(s, target, m+1, h)
	}
	return m
}

// BinaryIterativeSearch return the index of target within a sorted slice, use binary search (no recursive).
// If not found return -1.
// 二分查找
func BinaryIterativeSearch[T constraints.Ordered](sortedSlice []T, target T) int {
	i, j := 0, len(sortedSlice)-1
	for i <= j {
		m := i + (j-i)/2
		if sortedSlice[m] < target {
			i = m + 1
		} else if sortedSlice[m] > target {
			j = m - 1
		} else {
			return m
		}
	}
	return -1
}
