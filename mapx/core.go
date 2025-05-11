package mapx

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"sort"
)

// Keys 返回 map 里面的所有的 key。
// 需要注意：这些 key 的顺序是随机。
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values 返回 map 里面的所有的 value。
// 需要注意：这些 value 的顺序是随机。
func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

type KV[K any, V any] struct {
	Key   K
	Value V
}

// SortByKey 按 key 排序。
func SortByKey[K constraints.Ordered, V any](a, b KV[K, V]) bool {
	return a.Key < b.Key
}

// SortByValue 根据 value 排序。
func SortByValue[K any, V constraints.Ordered](a, b KV[K, V]) bool {
	return a.Value < b.Value
}

// defaultKeySort 是 SortMap 的默认排序逻辑
func defaultKeySort[K constraints.Ordered, V any](a, b KV[K, V]) bool {
	return a.Key < b.Key
}

// SortMap 对 map[K]V 进行排序。
//
// Example:
//
//	m := map[string]int{
//		"b": 2,
//		"a": 3,
//		"c": 1,
//	}
//
// sorted1 := SortMap(m)
// fmt.Println(sorted1) // []mapx.KV[string,int]{Key: "a", Value: 3} mapx.KV[string,int]{Key: "b", Value: 2} mapx.KV[string,int]{Key: "c", Value: 1}]
//
// sorted2 := SortMap(m, func(a, b mapx.KV[string, int]) bool {return a.Value < b.Value})
// fmt.Println(sorted2) //  []mapx.KV[string,int]{Key: "c", Value: 1} mapx.KV[string,int]{Key: "b", Value: 2} mapx.KV[string,int]{Key: "a", Value: 3}]
func SortMap[K constraints.Ordered, V any](m map[K]V, opts ...func(a, b KV[K, V]) bool) []KV[K, V] {
	result := make([]KV[K, V], 0, len(m))
	for key, value := range m {
		result = append(result, KV[K, V]{Key: key, Value: value})
	}

	var lessFunc func(i, j int) bool

	if len(opts) > 0 {
		lessFunc = func(i, j int) bool {
			return opts[0](result[i], result[j])
		}
	} else {
		lessFunc = func(i, j int) bool {
			return defaultKeySort(result[i], result[j])
		}
	}
	sort.Slice(result, lessFunc)

	return result
}

// InvertWithErr 接收一个 map[K]V 类型的输入，并返回一个反转后的 map[V]K, 并在遇到重复值时返回错误。
func InvertWithErr[K, V comparable](m map[K]V) (map[V]K, error) {
	invMap := make(map[V]K, len(m))
	for key, value := range m {
		if _, ok := invMap[value]; ok {
			return nil, fmt.Errorf("duplicate value found: %v", value)
		}
		invMap[value] = key
	}
	return invMap, nil
}
