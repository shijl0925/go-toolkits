package algorithm

import "golang.org/x/exp/constraints"

// BubbleSortV1 冒泡排序V1
func BubbleSortV1[T constraints.Ordered](s []T) {
	n := len(s)
	// 外循环: 未排序区间为 [0, i]
	for i := n - 1; i > 0; i-- {
		// 内循环: 将最大的元素移动到未排序区间的末尾
		for j := 0; j < i; j++ {
			if s[j] > s[j+1] {
				// 交换 s[j]与s[j+1]
				s[j], s[j+1] = s[j+1], s[j]
			}
		}
	}
}

// BubbleSortV2 冒泡排序V1
func BubbleSortV2[T constraints.Ordered](s []T) {
	n := len(s)
	// 外循环: 未排序区间为 [0, n-1]
	for i := 0; i < n; i++ {
		// 内循环: 将最小的元素移动到未排序区间的开头
		for j := i + 1; j < n; j++ {
			if s[i] > s[j] {
				// 交换 s[i]与s[j]
				s[i], s[j] = s[j], s[i]
			}
		}
	}
}

// SelectSort 选择排序
func SelectSort[T constraints.Ordered](s []T) {
	n := len(s)
	// 外循环: 未排序区间为 [i, n-1]
	for i := 0; i < n; i++ {
		// 内循环: 找到未排序区间的最小元素
		k := i
		for j := i + 1; j < n; j++ {
			if s[j] < s[k] {
				// 记录最小元素的索引
				k = j
			}
		}
		// 将该最小元素与当前索引处(i)的元素交换位置
		s[i], s[k] = s[k], s[i]
	}
}

// InsertSort 插入排序
func InsertSort[T constraints.Ordered](s []T) {
	n := len(s)
	// 外循环: 未排序区间未 [0, i]
	for i := 1; i < n; i++ {
		base := s[i]
		j := i - 1
		// 内循环: 将base插入到已排序区间的正确位置
		for ; j >= 0 && s[j] > base; j-- {
			s[j+1] = s[j] // 将s[j]向右移动一位
		}
		s[j+1] = base // 插入base
	}
}
