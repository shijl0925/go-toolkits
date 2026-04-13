package mapx

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// Keys 返回 map 里面的所有的 key。
// 需要注意：这些 key 的顺序是随机。
// This function is not concurrency-safe if the map is modified during iteration.
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values 返回 map 里面的所有的 value。
// 需要注意：这些 value 的顺序是随机。
// This function is not concurrency-safe if the map is modified during iteration.
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

// Intersect returns the intersection between two maps.
// 返回两个 map 的交集。
func Intersect[K comparable, V any](src, dst map[K]V) map[K]V {
	if len(src) == 0 || len(dst) == 0 {
		return make(map[K]V)
	}

	// 预估容量为较小 map 的大小
	capHint := len(src)
	if len(dst) < capHint {
		capHint = len(dst)
	}
	result := make(map[K]V, capHint)

	// 遍历较小的 map 以优化性能
	var iterMap, checkMap map[K]V
	if len(src) <= len(dst) {
		iterMap = src
		checkMap = dst
	} else {
		iterMap = dst
		checkMap = src
	}

	for key, iterValue := range iterMap {
		if checkValue, ok := checkMap[key]; ok && reflect.DeepEqual(iterValue, checkValue) {
			result[key] = iterValue // 或 checkValue，它们是相等的
		}
	}

	return result
}

// Merge maps, when the key is same, next value will overwrite previous value.
func Merge[K comparable, V any](maps ...map[K]V) map[K]V {
	size := 0
	for _, m := range maps {
		size += len(m)
	}

	result := make(map[K]V, size)

	for _, m := range maps {
		if m == nil {
			continue
		}

		for k, v := range m {
			result[k] = v
		}
	}

	return result
}

// Chain maps, when the key is same, next value will be ignored.
// It's like a merge, but when the key is same, next value will be ignored.
func Chain[K comparable, V any](maps ...map[K]V) map[K]V {
	size := 0
	for _, m := range maps {
		size += len(m)
	}

	result := make(map[K]V, size)

	for _, m := range maps {
		if m == nil {
			continue
		}

		for k, v := range m {
			if _, ok := result[k]; !ok {
				result[k] = v
			}
		}
	}

	return result
}

// FilterByKey filters a map by key.
// It returns a new map with only the key-value pairs where the key passes the given function.
func FilterByKey[K comparable, V any](m map[K]V, fn func(K) bool) map[K]V {
	result := make(map[K]V, len(m))
	for key, value := range m {
		if fn(key) {
			result[key] = value
		}
	}
	return result
}

// FilterByValue filters a map by value.
// It returns a new map with only the key-value pairs where the value passes the given function.
func FilterByValue[K comparable, V any](m map[K]V, fn func(V) bool) map[K]V {
	result := make(map[K]V, len(m))
	for key, value := range m {
		if fn(value) {
			result[key] = value
		}
	}
	return result
}

// Filter filters a map by key and value
// It return a new map contains all key and value pairs pass the predicate function.
func Filter[K comparable, V any](m map[K]V, fn func(key K, value V) bool) map[K]V {
	result := make(map[K]V, len(m))
	for k, v := range m {
		if fn(k, v) {
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

// SetIfAbsent set the value of the given key with the default value when the key is not present.
// 设置 map 的 key 的 value，如果 key 不存在，则添加 key 和 value。
func SetIfAbsent[K comparable, V any](m map[K]V, key K, defaultValue V) {
	if m == nil {
		return
	}
	if _, ok := m[key]; !ok {
		m[key] = defaultValue
	}
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
	result := make(map[V]K, len(m))
	for key, value := range m {
		if _, ok := result[value]; ok {
			return nil, fmt.Errorf("duplicate value found: %v (from key: %v)", value, key)
		}
		result[value] = key
	}
	return result, nil
}

// ToJson covert map to a json string
func ToJson(m map[string]any) (string, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// MapToStruct 将 map 转换为指定的结构体指针
func MapToStruct(m map[string]any, s any) error {
	val := reflect.ValueOf(s)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return fmt.Errorf("must be a non empty struct pointer")
	}

	elem := val.Elem()
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("must be a struct pointer")
	}

	t := elem.Type()
	for i := 0; i < elem.NumField(); i++ {
		fieldValue := elem.Field(i)
		if !fieldValue.CanSet() {
			continue // 跳过不可导出的字段
		}

		fieldType := t.Field(i)
		key, err := getFieldKey(fieldType, "json")

		if err != nil {
			return err
		}
		if key == "" {
			continue // 跳过被忽略的字段
		}

		value, ok := m[key]

		if !ok {
			continue // 跳过不存在的键
		}

		if fieldType.Anonymous && fieldType.Type.Kind() == reflect.Struct {
			continue // 跳过匿名结构体
		}
		if err := setField(fieldValue, value); err != nil {
			return fmt.Errorf("field %v conversion failed: %v", fieldType.Name, err)
		}
	}
	return nil
}

// getFieldKey 获取结构体字段对应的 map 键名
func getFieldKey(field reflect.StructField, tag string) (string, error) {
	filedTag := field.Tag.Get(tag)
	if filedTag == "-" {
		return "", nil // 忽略该字段
	}

	if filedTag != "" {
		// 处理类似 `json:"name,omitempty"` 的情况
		parts := strings.Split(filedTag, ",")
		return parts[0], nil
	}

	// 默认使用字段名小写
	return strings.ToLower(field.Name), nil
}

// setField 设置结构体字段值，处理类型转换
func setField(field reflect.Value, value any) error {
	targetType := field.Type()
	if value == nil {
		switch targetType.Kind() {
		case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
			if !field.CanSet() {
				// This check should ideally be redundant if MapToStruct filters unexported fields.
				return fmt.Errorf("internal error: cannot set field %s of type %s", field.Type().Name(), targetType)
			}
			field.Set(reflect.Zero(targetType)) // Sets to nil for these types
			return nil
		default:
			// Cannot assign Go's untyped nil to a non-nilable field type (e.g. int, string, struct).
			return fmt.Errorf("cannot assign nil to field of non-nilable type %v", targetType)
		}
	}

	val := reflect.ValueOf(value)

	// 类型完全匹配时直接赋值
	if val.Type().AssignableTo(targetType) {
		field.Set(val)
		return nil
	}

	// 处理指针类型
	if targetType.Kind() == reflect.Ptr {
		if val.IsZero() { // 如果值为零值，则设置为 nil
			field.Set(reflect.Zero(targetType))
			return nil
		}

		elemType := targetType.Elem() // 获取指针指向的类型
		newElem := reflect.New(elemType)
		if err := setField(newElem.Elem(), value); err != nil {
			return err
		}
		field.Set(newElem)
		return nil
	}

	// 处理嵌套结构体或匿名内嵌结构体
	if targetType.Kind() == reflect.Struct {
		if nestedMap, ok := value.(map[string]any); ok {
			newStruct := reflect.New(targetType).Elem()
			if err := MapToStruct(nestedMap, newStruct.Addr().Interface()); err != nil {
				return fmt.Errorf("failed to map to nested struct type %v: %w", targetType, err)
			}
			field.Set(newStruct)
			return nil
		}
	}

	return fmt.Errorf("unsupported type conversion: value of type %v to field of type %v", val.Type(), targetType)
}
