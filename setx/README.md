# setx

`setx` 基于 `map[T]struct{}` 实现泛型集合，提供常见的增删查、集合运算和切片互转能力。对于需要“去重、成员判断、集合比较”的业务，它通常比手写 map 更直观。

## 安装

```bash
go get github.com/shijl0925/go-toolkits/setx
```

## 核心类型

- `Set[T comparable]`
- `New[T comparable](elements ...interface{}) Set[T]`

## 常用能力

### 基本操作

- `Add`
- `Remove`
- `Exists`
- `Len`
- `IsEmpty`
- `Clear`
- `Pop`
- `Equal`

### 转换与遍历

- `Keys`
- `ToSlice`
- `Iterate`
- `Update`

### 集合运算

- `Difference`
- `Intersect`
- `Union`
- `ContainAny`
- `ContainAll`

## 快速示例

```go
s := setx.New[int](1, 2, 3)
s.Add(4)
s.Remove(2)

exists := s.Exists(3)           // true
items := s.ToSlice()            // 顺序不固定
common := setx.Intersect(s, setx.New[int](3, 4, 5))
```

## 设计特点

### New 的入参

`New` 支持两种常见初始化方式：

```go
setx.New[int](1, 2, 3)
setx.New[int]([]int{1, 2, 3})
```

如果传入的参数类型不匹配，当前实现不会返回错误，而是打印提示并跳过该值。

### 返回顺序

`Keys`、`ToSlice`、`Pop`、`Difference`、`Intersect`、`Union` 的结果顺序都**不保证稳定**，这是基于 Go map 的天然特性。

### Update 的语义

`Update(dst)` 会把 `dst` 中的元素合并到当前集合，并返回更新后的当前集合。

## 注意事项

- `setx` **不是并发安全容器**。如果集合会被多个 goroutine 同时读写，请自行加锁。
- `Pop` 会删除并返回一个任意元素，适合“取一个出来处理”的场景，不适合依赖顺序的逻辑。
- 集合运算函数返回切片而不是 `Set[T]`，便于直接消费或序列化。

## 适用场景

- 标签/权限/特征值去重
- 两组数据的交集、差集分析
- 快速判断“是否包含任一项 / 是否完全包含”

