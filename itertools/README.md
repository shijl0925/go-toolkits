# itertools

`itertools` 提供一组偏函数式风格的迭代工具，适合在 Go 中表达“映射、过滤、归约、分组、组合”等数据处理逻辑。它的设计目标不是替代原生 `for`，而是让常见的数据流操作更清晰。

## 安装

```bash
go get github.com/shijl0925/go-toolkits/itertools
```

## 核心能力

### 映射与过滤

- `Map[T, U]`：把 `[]T` 映射成 `[]U`
- `Filter[T]`：保留满足条件的元素
- `DropWhile[T]`：从左向右跳过满足条件的前缀元素

### 聚合

- `Reduce[T, U]`：基于初始值归约
- `All[T]`：所有元素都满足条件时返回 `true`
- `Any[T]`：任一元素满足条件时返回 `true`

### 组合与分组

- `Zip[T, U]`：按位置把两个切片组合为 `[]Tuple[T, U]`
- `CombineToMap[K, V]`：将 keys 和 values 组装为 map
- `GroupBy[T, U]`：按指定规则分组

### 序列生成

- `Range[T]`：生成整数或浮点数序列

## 快速示例

```go
nums := []int{1, 2, 3, 4, 5}

evens := itertools.Filter(nums, func(v int) bool {
    return v%2 == 0
})

squares := itertools.Map(evens, func(v int) int {
    return v * v
})

sum := itertools.Reduce(squares, func(acc, v int) int {
    return acc + v
}, 0)
```

## API 说明

### Zip

`Zip(a, b)` 的结果长度取决于**较短的那个切片**。它不会补零，也不会报错。

```go
pairs := itertools.Zip([]string{"A", "B"}, []int{1, 2, 3})
// 结果只包含 2 组元素
```

### CombineToMap

- 当 `keys` 和 `values` 长度不一致时，超出的部分不会进入结果。
- 如果 `keys` 有重复项，后写入的值会覆盖前面的值。

### Range

- 支持整数和浮点数。
- `step` 需要与范围方向匹配，否则结果可能为空。
- 适合构造简单数值序列，不建议用于高精度金融计算。

### GroupBy

返回结果为 `map[U][]T`，常见于：

- 按状态分组
- 按日期分组
- 按首字母或类别归类

## 适用场景

- 数据清洗与转换
- 统计前的预处理
- 需要让“处理流程”比“循环细节”更突出时
- 与 `slicex` 组合完成复杂切片分析

