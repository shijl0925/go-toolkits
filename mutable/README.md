# mutable

`mutable` 是 `slicex` 的互补模块：它更强调**原地修改切片**，适合在意内存分配、希望直接修改底层数据的场景。

## 安装

```bash
go get github.com/shijl0925/go-toolkits/mutable
```

## 核心能力

### 原地映射与过滤

- `Map`
- `Filter`
- `Remove`

### 原地元素操作

- `Reverse`
- `Pop`
- `Shift`
- `Drop`
- `DropLeft`

### 原地重组

- `Shuffle`

## 快速示例

```go
nums := []int{1, 2, 3, 4}
mutable.Map(nums, func(v int) int { return v * 10 })
// nums == []int{10, 20, 30, 40}

mutable.Reverse(nums)
// nums == []int{40, 30, 20, 10}
```

## 为什么有些函数要传 `*[]T`？

像 `Filter`、`Remove`、`Pop`、`Shift`、`Drop`、`DropLeft` 这类操作不仅会修改元素，还会改变切片长度，所以需要接收 `*[]T` 才能把新的切片头部信息回写给调用方。

```go
items := []int{1, 2, 3, 4, 5}
mutable.Filter(&items, func(v int) bool { return v%2 == 1 })
// items == []int{1, 3, 5}
```

## 行为说明

- `Map` 和 `Reverse` 直接修改传入切片。
- `Shuffle` 使用标准库随机打乱顺序。
- `Pop` / `Shift` 在空切片或 `nil` 指针场景下返回零值和 `false`。
- `Drop` / `DropLeft` 对 `nil` 指针或负数删除量会返回错误。
- `Filter` / `Remove` 接收 `nil` 指针时返回 `nil`，不会 panic。

## 何时使用 mutable？

推荐在以下场景使用：

- 数据量较大，希望减少额外分配
- 希望“直接修改原切片”而不是接收新结果
- 处理流程明确，不担心原数据被修改带来的副作用

如果你更偏好不改变输入值的写法，请改用 `slicex`。

