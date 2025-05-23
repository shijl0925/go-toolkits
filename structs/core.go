package structs

import (
	"reflect"
)

// StructToMap convert structs to map[string]any by reflect.
func StructToMap(s any) map[string]any {
	s = structVal(s)

	result := make(map[string]any)

	val := reflect.ValueOf(s)
	t := val.Type()

	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)
		value := val.Field(i).Interface()
		fieldName := fieldType.Name

		tagName, tagOpts := parseTag(fieldType.Tag.Get("json"))
		if tagName != "" {
			fieldName = tagName
		}

		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				continue
			}
			value = v.Elem().Interface()
		} else if v.Kind() == reflect.Struct {
			if tagOpts.Has("flatten") {
				// 判断变量是否为匿名结构体
				if fieldType.Anonymous {
					for k, v := range StructToMap(value) {
						result[k] = v
					}
					continue
				}
				value = StructToMap(value)
			}
		}

		if tagOpts.Has("omitempty") {
			zero := reflect.Zero(val.Field(i).Type()).Interface()
			if reflect.DeepEqual(value, zero) {
				continue
			}
		}

		result[fieldName] = value
	}

	return result
}

func structVal(s any) any {
	v := reflect.ValueOf(s)

	// if pointer get the underlying element≤
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		panic("not struct")
	}

	return v.Interface()
}
