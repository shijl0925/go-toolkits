# mapx - Go Map 工具包
mapx 是一个用于处理 Go 语言 map 的工具包，提供了丰富的 map 操作函数，支持泛型类型，能够处理各种数据类型。
功能列表

## 基本操作
* `Keys[K comparable, V any](m map[K]V) []K` - 返回 map 中所有的 key（顺序随机，非并发安全）
* `Values[K comparable, V any](m map[K]V) []V` - 返回 map 中所有的 value（顺序随机，非并发安全）
* `HasKey[K comparable, V any](m map[K]V, key K) bool` - 检查 map 是否包含指定 key
* `GetOrDefault[K comparable, V any](m map[K]V, key K, defaultValue V) V` - 获取指定 key 的值，如果不存在则返回默认值
* `SetIfAbsent[K comparable, V any](m map[K]V, key K, defaultValue V)` - 如果 key 不存在，则设置默认值

## Map 合并操作
* `Merge[K comparable, V any](maps ...map[K]V) map[K]V` - 合并多个 map，相同 key 时后面的值会覆盖前面的值
* `Chain[K comparable, V any](maps ...map[K]V) map[K]V` - 链接多个 map，相同 key 时忽略后面的值

## Map 集合操作
* `Intersect[K comparable, V any](src, dst map[K]V) map[K]V` - 返回两个 map 的交集（key 和 value 都相等）

## Map 过滤操作
* `FilterByKey[K comparable, V any](m map[K]V, fn func(K) bool) map[K]V` - 根据 key 条件过滤 map
* `FilterByValue[K comparable, V any](m map[K]V, fn func(V) bool) map[K]V` - 根据 value 条件过滤 map
* `Filter[K comparable, V any](m map[K]V, fn func(key K, value V) bool) map[K]V` - 根据 key 和 value 条件过滤 map

## Map 排序操作
* `SortByKey[K comparable, V any](m map[K]V, less func(a, b K) bool) []KV[K, V]` - 按 key 排序返回 KV 结构体切片
* `SortByValue[K comparable, V any](m map[K]V, less func(a, b V) bool) []KV[K, V]` - 按 value 排序返回 KV 结构体切片

## Map 转换操作
* `InvertWithErr[K, V comparable](m map[K]V) (map[V]K, error)` - 反转 map 的 key 和 value，遇到重复值时返回错误
* `ToJson(m map[string]any) (string, error)` - 将 map 转换为 JSON 字符串
* `MapToStruct(m map[string]any, s any) error` - 将 map 转换为指定的结构体指针

## 安装
```shell
go get github.com/shijl0925/go-toolkits/mapx
```