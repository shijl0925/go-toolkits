# itertools - Go 迭代工具包
itertools 是一个用于处理 Go 语言切片的函数式编程工具包，提供了丰富的迭代操作函数，支持泛型类型，能够处理各种数据类型。
功能列表

## 映射和过滤
* `Map[T any, U any](s []T, fn func(T) U) []U` - 对切片中每个元素应用函数并返回结果切片
* `Filter[T any](s []T, fn func(T) bool) []T` - 过滤切片中满足条件的元素
* `DropWhile[T any](s []T, fn func(T) bool) []T` - 从左侧丢弃满足条件的元素直到遇到不满足条件的元素

## 聚合操作
* `Reduce[T any, U any](s []T, fn func(U, T) U, initial U) U` - 对切片进行归约操作，将二元操作函数依次应用于初始值和切片中的每个元素
* `All[T any](s []T, fn func(T) bool) bool` - 检查切片中所有元素是否满足条件
* `Any[T any](s []T, fn func(T) bool) bool` - 检查切片中是否有元素满足条件

## 组合操作
* `Zip[T any, U any](a []T, b []U) []Tuple[T, U]` - 将两个切片组合成元组切片，长度为两个切片中较短的那个
* `CombineToMap[K comparable, V any](keys []K, values []V) map[K]V` - 将两个切片组合成map，若keys中存在重复键则后面的值覆盖前面的值

## 分组操作
* `GroupBy[T any, U comparable](slice []T, fn func(T) U) map[U][]T` - 根据键函数对切片元素进行分组，返回以键函数结果为键、元素切片为值的map

## 范围生成
* `Range[T constraints.Integer | constraints.Float](start, end, step T) []T` - 生成指定范围内的数字序列，支持整数和浮点数类型

## 安装
```shell
go get github.com/shijl0925/go-toolkits/itertools
```