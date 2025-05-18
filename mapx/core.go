package mapx

import (
	"fmt"
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

// HasKey checks if map has key or not.
func HasKey[K comparable, V any](m map[K]V, key K) bool {
	_, ok := m[key]
	return ok
}

// Merge maps, when the key is same, next value will overwrite previous value.
func Merge[K comparable, V any](maps ...map[K]V) map[K]V {
	size := 0
	for _, m := range maps {
		size += len(m)
	}

	result := make(map[K]V, size)

	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}

	return result
}

// GetOrDefault returns the value of the given key or a default value if the key is not present.
func GetOrDefault[K comparable, V any](m map[K]V, key K, defaultValue V) V {
	if value, ok := m[key]; ok {
		return value
	}
	return defaultValue
}

type KV[K comparable, V any] struct {
	Key   K
	Value V
}

// SortByKey 按 key 排序。
//
// Example:
//
//	m := map[string]int{
//		"b": 2,
//		"a": 3,
//		"c": 1,
//	}
//
// result := SortByKey(m, func(a, b string) bool { return a < b })
// fmt.Println(result) // []mapx.KV[string,int]{Key: "a", Value: 3} mapx.KV[string,int]{Key: "b", Value: 2} mapx.KV[string,int]{Key: "c", Value: 1}]
func SortByKey[K comparable, V any](m map[K]V, less func(a, b K) bool) []KV[K, V] {
	result := make([]KV[K, V], len(m))

	keys := Keys[K, V](m)
	sort.Slice(keys, func(i, j int) bool { return less(keys[i], keys[j]) })

	for i, k := range keys {
		result[i] = KV[K, V]{Key: k, Value: m[k]}
	}

	return result
}

// SortByValue 按 value 排序。
//
// Example:
//
//	m := map[string]int{
//		"b": 2,
//		"a": 3,
//		"c": 1,
//	}
//
// result := SortByValue(m, func(a, b int) bool {return a < b})
// fmt.Println(result) //  []mapx.KV[string,int]{Key: "c", Value: 1} mapx.KV[string,int]{Key: "b", Value: 2} mapx.KV[string,int]{Key: "a", Value: 3}]
func SortByValue[K comparable, V any](m map[K]V, less func(a, b V) bool) []KV[K, V] {
	result := make([]KV[K, V], 0, len(m))
	for k, v := range m {
		result = append(result, KV[K, V]{Key: k, Value: v})
	}

	sort.Slice(result, func(i, j int) bool { return less(result[i].Value, result[j].Value) })

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
