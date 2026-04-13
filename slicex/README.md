# slicex

`slicex` 是一个面向 Go 泛型的切片工具包，覆盖统计、删除、插入、集合运算、排序、重组等高频场景。它更适合**偏函数式、返回新结果**的切片处理；如果你更在意原地修改，可以搭配 `mutable` 模块使用。

## 安装

```bash
go get github.com/shijl0925/go-toolkits/slicex
```

## 适用场景

- 快速完成常见切片变换与清洗
- 用统一 API 替代重复的 for 循环
- 做集合比较、计数、分块、拼接等数据整理工作

## 核心能力

### 统计与比较

- `Sum`、`Avg`、`Max`、`Min`
- `Equal`、`EqualUnordered`
- `IsAscending`、`IsDescending`

### 删除与截取

- `RemoveHead`、`RemoveTail`
- `Pop`、`Shift`
- `Drop`、`DropLeft`
- `Remove`、`Delete`、`DeleteAt`、`DeleteMany`

### 追加与插入

- `Append`、`Extend`
- `InsertAt`、`InsertMany`
- `AddMany`

### 重组与转换

- `Reverse`、`Repeat`、`Concat`
- `Chunk`、`FlattenDeep`
- `Product`、`Combine`
- `JoinSlice`、`Replace`
- `RightPadding`、`LeftPadding`

### 查找与集合运算

- `Contains`、`ContainsAny`、`ContainsAll`
- `Index`、`LastIndex`、`IndexAll`
- `Unique`、`FindUniques`
- `Difference`、`Intersect`、`Union`、`Without`
- `Counter`

### 排序

- `SortSlice`

## 快速示例

```go
nums := []int{1, 2, 2, 3, 4}

unique := slicex.Unique(nums)          // [1 2 3 4]
without := slicex.Without(nums, 2)     // [1 3 4]
chunked := slicex.Chunk(nums, 2)       // [[1 2] [2 3] [4]]
counts := slicex.Counter(nums)         // map[int]int{1:1, 2:2, 3:1, 4:1}
joined := slicex.JoinSlice(",", nums) // "1,2,2,3,4"
```

## 常用说明

### 返回新切片还是修改原切片？

大多数函数会返回一个新的结果视图，便于链式处理；但以下函数需要特别留意：

- `Compact`：会在原底层数组上原地压缩有效元素。
- `SortSlice`：会直接排序传入切片。

如果你明确需要“原地修改”，建议优先查看 `mutable`。

### Delete 与 DeleteAt 的区别

- `Delete`：语义直观。
- `DeleteAt`：实现更偏性能导向，适合高频按索引删除。

### InsertAt / AddMany 何时使用

- 插入单个元素：`InsertAt`
- 插入多个元素：`AddMany`
- 如果你只是在尾部增加元素，直接用 `Append` / `Extend` 更自然。

## 边界行为

- 对空切片进行 `Pop` / `Shift` 时，会返回零值和 `false`。
- `Drop` / `DropLeft` 删除数量超过长度时，会尽可能删除到空切片。
- `Avg` 返回 `float64`，适合统一处理整数与浮点切片。
- `FlattenDeep` 的返回值类型为 `any`，适合处理嵌套深度不固定的数据。

## 示例：集合分析

```go
src := []string{"go", "go", "toolkit", "docs"}
dst := []string{"go", "lint"}

common := slicex.Intersect(src, dst) // [go]
extra := slicex.Difference(src, dst) // [toolkit docs]
all := slicex.Union(src, dst)        // 去重后的并集
```

