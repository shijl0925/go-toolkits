# algorithm - Go 算法工具包

`algorithm` 是一个提供常用算法实现的工具包，包括排序算法和搜索算法，支持泛型类型。

## 功能列表

### 排序算法
* `BubbleSortV1[T constraints.Ordered](s []T)` - 冒泡排序实现版本1（将最大元素移动到末尾）
* `BubbleSortV2[T constraints.Ordered](s []T)` - 冒泡排序实现版本2（将最小元素移动到开头）
* `SelectSort[T constraints.Ordered](s []T)` - 选择排序（每次选择最小元素放到已排序区间）
* `InsertSort[T constraints.Ordered](s []T)` - 插入排序（将元素插入到已排序区间的正确位置）

### 搜索算法
* `BinarySearch[T constraints.Ordered](s []T, target T, l, h int) int` - 二分查找（递归实现）
* `BinaryIterativeSearch[T constraints.Ordered](sortedSlice []T, target T) int` - 二分查找（迭代实现）

## 安装

```shell
go get github.com/shijl0925/go-toolkits/algorithm
```


## 使用示例

```go
package main

import (
    "fmt"
    "github.com/shijl0925/go-toolkits/algorithm"
)

func main() {
    // 排序示例
    numbers := []int{64, 34, 25, 12, 22, 11, 90}
    fmt.Println("Original array:", numbers)
    
    // 使用冒泡排序
    algorithm.BubbleSortV1(numbers)
    fmt.Println("Sorted array:", numbers)
    
    // 搜索示例
    sortedNumbers := []int{11, 12, 22, 25, 34, 64, 90}
    target := 25
    
    // 使用二分查找（递归）
    index1 := algorithm.BinarySearch(sortedNumbers, target, 0, len(sortedNumbers)-1)
    fmt.Printf("BinarySearch: %d found at index %d\n", target, index1)
    
    // 使用二分查找（迭代）
    index2 := algorithm.BinaryIterativeSearch(sortedNumbers, target)
    fmt.Printf("BinaryIterativeSearch: %d found at index %d\n", target, index2)
}
```


## 算法说明

### 排序算法

1. **冒泡排序 V1** ([BubbleSortV1](https://github.com/shijl0925/go-toolkits/blob/main/algorithm/sort.go#L6-L18))
    - 通过重复遍历列表，比较相邻元素并交换顺序错误的元素
    - 每次遍历将最大元素"冒泡"到未排序部分的末尾
    - 时间复杂度: O(n²)，空间复杂度: O(1)

2. **冒泡排序 V2** ([BubbleSortV2](https://github.com/shijl0925/go-toolkits/blob/main/algorithm/sort.go#L21-L33))
    - 另一种冒泡排序实现
    - 每次遍历将最小元素"冒泡"到已排序部分的开头
    - 时间复杂度: O(n²)，空间复杂度: O(1)

3. **选择排序** ([SelectSort](https://github.com/shijl0925/go-toolkits/blob/main/algorithm/sort.go#L36-L51))
    - 每次从未排序部分选择最小元素，放到已排序部分的末尾
    - 时间复杂度: O(n²)，空间复杂度: O(1)

4. **插入排序** ([InsertSort](https://github.com/shijl0925/go-toolkits/blob/main/algorithm/sort.go#L54-L66))
    - 将未排序元素逐个插入到已排序部分的正确位置
    - 时间复杂度: O(n²)，空间复杂度: O(1)

### 搜索算法

1. **二分查找** ([BinarySearch](https://github.com/shijl0925/go-toolkits/blob/main/algorithm/search.go#L8-L19))
    - 在已排序数组中查找目标值
    - 通过递归方式实现
    - 时间复杂度: O(log n)，空间复杂度: O(log n)

2. **二分查找（迭代）** ([BinaryIterativeSearch](https://github.com/shijl0925/go-toolkits/blob/main/algorithm/search.go#L24-L37))
    - 在已排序数组中查找目标值
    - 通过迭代方式实现
    - 时间复杂度: O(log n)，空间复杂度: O(1)