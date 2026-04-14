# algorithm

`algorithm` 提供几种基础排序与二分查找实现，适合作为教学、面试练习、基础算法验证或对比实验使用。它们都是简洁直观的实现版本，而不是以极致性能为目标的生产级排序库。

## 安装

```bash
go get github.com/shijl0925/go-toolkits/algorithm
```

## 包含的算法

### 排序（原地修改）

- `BubbleSortV1`
- `BubbleSortV2`
- `SelectSort`
- `InsertSort`

### 搜索

- `BinarySearch`
- `BinaryIterativeSearch`

## 快速示例

```go
nums := []int{64, 34, 25, 12, 22, 11, 90}
algorithm.InsertSort(nums)

idx := algorithm.BinaryIterativeSearch(nums, 25)
```

## 使用说明

### 排序函数

所有排序函数都会直接修改传入切片。

| 函数 | 特点 | 时间复杂度 | 空间复杂度 |
| --- | --- | --- | --- |
| `BubbleSortV1` | 经典冒泡，最大值逐轮后移 | O(n²) | O(1) |
| `BubbleSortV2` | 变体实现，最小值逐轮前移 | O(n²) | O(1) |
| `SelectSort` | 每轮选择最小值放到前面 | O(n²) | O(1) |
| `InsertSort` | 对近乎有序数据更友好 | O(n²) | O(1) |

### 二分查找

- 二分查找要求输入切片已经有序。
- `BinarySearch` 是递归实现，调用时需要传入左右边界。
- `BinaryIterativeSearch` 是迭代实现，调用更直接。
- 未找到目标值时，两者都返回 `-1`。

## 何时使用

- 学习和理解基础算法过程
- 编写小型练习或单元测试示例
- 对比不同简单排序策略的行为

## 注意事项

- 对于生产场景的通用排序，优先考虑标准库 `sort` / `slices`。
- `BinarySearch` 会在边界越界时自动修正范围，但前提仍然是输入已排序。

