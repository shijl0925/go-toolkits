# mutable - Go 可变切片工具包

`mutable` 是一个用于处理可变切片操作的工具包，提供了直接修改切片内容的函数，支持泛型类型。

## 功能列表

### 切片映射和过滤
* `Map[T any](s []T, fn func(T) T)` - 对切片中的每个元素应用函数并就地修改
* `Filter[T any](s *[]T, fn func(T) bool) []T` - 过滤切片中满足条件的元素并就地修改切片
* `Remove[T comparable](s *[]T, element T) []T` - 从切片中移除指定元素并就地修改切片

### 切片元素操作
* `Reverse[T any](s []T)` - 就地反转切片中的元素
* `Pop[T any](s *[]T) (T, bool)` - 弹出切片的最后一个元素
* `Shift[T any](s *[]T) (T, bool)` - 弹出切片的第一个元素
* `Drop[T any](s *[]T, n int) (bool, error)` - 从切片末尾删除n个元素
* `DropLeft[T any](s *[]T, n int) (bool, error)` - 从切片开头删除n个元素

### 切片重组
* `Shuffle[T any](s []T)` - 使用 Fisher-Yates 算法随机打乱切片元素

## 安装

```shell
go get github.com/shijl0925/go-toolkits/mutable
```


## 使用示例

```go
package main

import (
    "fmt"
    "github.com/shijl0925/go-toolkits/mutable"
)

func main() {
    // Map 示例 - 就地修改切片
    numbers := []int{1, 2, 3, 4, 5}
    mutable.Map(numbers, func(x int) int { return x * 2 })
    fmt.Println("Doubled:", numbers) // [2 4 6 8 10]
    
    // Filter 示例 - 过滤偶数
    slice := &[]int{1, 2, 3, 4, 5, 6}
    mutable.Filter(slice, func(x int) bool { return x%2 == 0 })
    fmt.Println("Even numbers:", *slice) // [2 4 6]
    
    // Pop 示例
    s := &[]int{1, 2, 3, 4, 5}
    element, ok := mutable.Pop(s)
    if ok {
        fmt.Println("Popped element:", element) // 5
        fmt.Println("Remaining slice:", *s)     // [1 2 3 4]
    }
    
    // Reverse 示例
    data := []int{1, 2, 3, 4, 5}
    mutable.Reverse(data)
    fmt.Println("Reversed:", data) // [5 4 3 2 1]
}
```
