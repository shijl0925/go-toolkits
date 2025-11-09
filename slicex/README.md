# slicex - Go 切片工具包
slicex 是一个用于处理 Go 语言切片的工具包，提供了丰富的切片操作函数，支持泛型类型，能够处理各种数据类型。
功能列表

## 基本数学运算
* `Sum[T constraints.Number](s []T) T` - 计算切片元素总和
* `Avg[T constraints.RealNumber](s []T) float64` - 计算切片元素平均值
* `Max[T constraints.Ordered](s []T) T` - 获取切片中的最大值
* `Min[T constraints.Ordered](s []T) T` - 获取切片中的最小值

## 切片元素操作
* `RemoveHead[T any](s []T) []T` - 移除切片的第一个元素
* `RemoveTail[T any](s []T) []T` - 移除切片的最后一个元素
* `Pop[T any](s []T) (T, bool)` - 弹出切片的最后一个元素
* `Shift[T any](s []T) (T, bool)` - 弹出切片的第一个元素
* `Drop[T any](s []T, n int) ([]T, bool)` - 从切片末尾删除n个元素
* `DropLeft[T any](s []T, n int) ([]T, bool)` - 从切片开头删除n个元素
* `Remove[T comparable](s []T, element T) []T` - 移除切片中所有指定元素
* `Delete[T any](s []T, index int) ([]T, bool)` - 删除指定索引的元素
* `DeleteAt[T any](s []T, index int) ([]T, bool)` - 删除指定索引的元素(性能优化版)
* `DeleteMany[T any](s []T, start, end int) ([]T, bool)` - 删除指定范围内的元素

## 切片添加和插入
* `Append[T any](src []T, elements ...T) []T` - 向切片追加元素
* `Extend[T any](src []T, elements []T) []T` - 扩展切片
* `InsertAt[T any](src []T, element T, index int) ([]T, error)` - 在指定位置插入元素
* `AddMany[T any](src []T, elements []T, index int) ([]T, error)` - 在指定位置添加多个元素

## 切片转换和重组
* `Reverse[T any](s []T) []T` - 反转切片
* `Repeat[T any](s T, n int) []T` - 重复元素n次
* `Product[T any](slices ...[]T) [][]T` - 计算多个切片的笛卡尔积
* `Chunk[T any](s []T, size int) [][]T` - 将切片分块
* `Compact[T comparable](s []T) []T` - 移除假值元素（注意：会修改原始切片）
* `FlattenDeep(slice any) any` - 递归展平嵌套切片
* `Combine[K comparable, V any](keys []K, values []V) map[K]V` - 将两个切片组合成map

## 切片比较和查找
* `Contains[T comparable](src []T, dst T) bool` - 检查切片是否包含指定元素
* `ContainsAny[T comparable](src, dst []T) bool` - 检查切片是否包含任意指定元素
* `ContainsAll[T comparable](src, dst []T) bool` - 检查切片是否包含所有指定元素
* `Index[T comparable](src []T, dst T) int` - 查找元素第一次出现的索引
* `LastIndex[T comparable](src []T, dst T) int` - 查找元素最后一次出现的索引
* `IndexAll[T comparable](src []T, dst T) []int` - 查找元素所有出现的索引
* `Equal[T comparable](s1, s2 []T) bool` - 检查两个切片是否相等
* `EqualUnordered[T comparable](s1, s2 []T) bool` - 检查两个切片是否在无序情况下相等

## 切片去重和集合操作
* `Unique[T comparable](s []T) []T` - 移除重复元素
* `FindUniques[T comparable](s []T) []T` - 查找唯一不重复的元素
* `Difference[T comparable](src, dst []T) []T` - 计算两个切片的差集
* `Intersect[T comparable](src, dst []T) []T` - 计算两个切片的交集
* `Union[T comparable](lists ...[]T) []T` - 计算多个切片的并集
* `Without[T comparable](s []T, excludes ...T) []T` - 排除指定元素

## 排序和计数
* `SortSlice[T any](s []T, less func(a, b T) bool)` - 自定义排序（注意：会原地修改切片）
* `IsAscending[T constraints.Ordered](s []T) bool` - 检查切片是否升序排列
* `IsDescending[T constraints.Ordered](s []T) bool` - 检查切片是否降序排列
* `Counter[T comparable](s []T) map[T]int` - 统计元素出现次数

## 切片连接
* `Concat[T any](s ...[]T) []T` - 连接多个切片

## 元素替换
* `Replace[T comparable](s []T, oldElement, newElement T, n int) []T` - 替换切片中的元素

## 拼接操作
* `JoinSlice[T any](sep string, s []T) string` - 将切片连接成字符串

## 填充操作
* `RightPadding[T any](s []T, v T, n int) []T` - 右侧填充
* `LeftPadding[T any](s []T, v T, n int) []T` - 左侧填充

## 安装
```shell
go get github.com/shijl0925/go-toolkits/slicex
```