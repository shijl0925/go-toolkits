# mapx

`mapx` 提供一组围绕 Go `map` 的实用工具，覆盖键值提取、合并、过滤、排序以及与结构体 / JSON 之间的转换。适合在业务代码里减少重复样板。

## 安装

```bash
go get github.com/shijl0925/go-toolkits/mapx
```

## 核心能力

### 基本操作

- `Keys`
- `Values`
- `HasKey`
- `GetOrDefault`
- `SetIfAbsent`

### 合并与集合

- `Merge`：后面的 map 覆盖前面的同名 key
- `Chain`：保留第一个出现的值，忽略后续重复 key
- `Intersect`：仅保留 **key 和 value 都相同** 的项

### 过滤与排序

- `FilterByKey`
- `FilterByValue`
- `Filter`
- `SortByKey`
- `SortByValue`

### 转换能力

- `InvertWithErr`
- `ToJson`
- `MapToStruct`

## 快速示例

```go
m1 := map[string]int{"a": 1, "b": 2}
m2 := map[string]int{"b": 20, "c": 3}

merged := mapx.Merge(m1, m2) // map[a:1 b:20 c:3]
chained := mapx.Chain(m1, m2) // map[a:1 b:2 c:3]

filtered := mapx.FilterByValue(merged, func(v int) bool {
    return v >= 3
})
```

## 使用建议

### Merge 与 Chain 如何选择？

- 想让**后值覆盖前值**：用 `Merge`
- 想让**首个值优先**：用 `Chain`

### Intersect 的判断规则

`Intersect` 不是只比较 key，而是要求 **两个 map 中同一个 key 的 value 也相等**，才会进入结果。

### SortByKey / SortByValue 的返回值

排序函数返回的是 `[]KV[K, V]`，不是 map。这样设计是为了保留排序后的顺序，适合直接用于展示、遍历或进一步转换。

### InvertWithErr 的适用前提

只有当原 map 的 value 全部唯一时，`InvertWithErr` 才能成功。如果存在重复 value，会返回错误。

## 注意事项

- Go 原生 `map` **不是并发安全容器**。`Keys`、`Values` 等函数也不解决并发读写问题。
- `Keys` 和 `Values` 的顺序天然不稳定；如需稳定输出，请先排序。
- `SetIfAbsent` 会直接修改原 map。
- `MapToStruct` 适合结构清晰、字段可映射的场景；复杂类型转换建议先做预处理。

## 示例：排序后输出

```go
items := mapx.SortByValue(map[string]int{
    "go":   3,
    "docs": 1,
    "test": 2,
}, func(a, b int) bool {
    return a < b
})
```

