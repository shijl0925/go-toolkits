//nolint:typecheck
package algorithm

import "golang.org/x/exp/constraints"

// BinarySearch return the index of target within a sorted slice, use binary search (recursive call itself).
// If not found return -1.
// 二分查找
func BinarySearch[T constraints.Ordered](sortedSlice []T, target T, lowIndex, highIndex int) int {
	if lowIndex > highIndex || len(sortedSlice) == 0 {
		return -1
	}
	midIndex := lowIndex + (highIndex-lowIndex)/2
	if sortedSlice[midIndex] > target {
		return BinarySearch(sortedSlice, target, lowIndex, midIndex-1)
	} else if sortedSlice[midIndex] < target {
		return BinarySearch(sortedSlice, target, midIndex+1, highIndex)
	}
	return midIndex
}

// BinaryIterativeSearch return the index of target within a sorted slice, use binary search (no recursive).
// If not found return -1.
// 二分查找
func BinaryIterativeSearch[T constraints.Ordered](sortedSlice []T, target T) int {
	i, j := 0, len(sortedSlice)
	for i < j {
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
