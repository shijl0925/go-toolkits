# setx - Go 集合工具包
setx 是一个用于处理 Go 语言集合(set)的工具包，基于 map 实现，提供了丰富的集合操作函数，支持泛型类型。

## 集合类型
* `Set[T comparable]` - 基于 map[T]struct{} 实现的泛型集合类型

## 构造函数
* `New[T comparable](elements ...interface{}) Set[T]` - 创建一个新的集合并添加元素，支持单个元素或元素切片

## 集合基本操作
* `Add(elements ...T)` - 向集合中添加元素
* `Remove(element T)` - 从集合中移除元素
* `Len() int` - 返回集合中元素的数量
* `Exists(element T) bool` - 检查元素是否存在于集合中
* `Clear()` - 清空集合中的所有元素
* `Equal(dst Set[T]) bool` - 检查两个集合是否相等
* `Pop() (v T, ok bool)` - 弹出集合中的一个元素，如果集合为空则返回零值和false
* `IsEmpty() bool` - 检查集合是否为空

## 集合转换操作
* `Keys() []T` - 返回包含集合中所有元素的切片(顺序不确定)
* `ToSlice() []T` - 返回包含集合中所有元素的切片
* `Iterate(fn func(item T))` - 对集合中的每个元素执行函数

## 集合更新操作
* `Update(dst Set[T]) Set[T]` - 将目标集合的元素添加到当前集合中

## 集合运算操作
* `Difference[T comparable](src, dst Set[T]) []T` - 计算两个集合的差集(src中存在但dst中不存在的元素)
* `Intersect[T comparable](src, dst Set[T]) []T` - 计算两个集合的交集(两个集合都存在的元素)
* `Union[T comparable](src, dst Set[T]) []T` - 计算两个集合的并集(两个集合中所有元素的集合)
* `ContainAny[T comparable](src, dst Set[T]) bool` - 检查两个集合是否存在交集
* `ContainAll[T comparable](src, dst Set[T]) bool` - 检查源集合是否包含目标集合的所有元素

## 安装
```shell
go get github.com/shijl0925/go-toolkits/setx
```