package structs

import (
	"reflect"
)

// StructToMap convert structs to map[string]any by reflect.
func StructToMap(s any) map[string]any {
	result := make(map[string]any)

	val := reflect.ValueOf(s)
	t := val.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := val.Field(i).Interface()
		name := field.Name

		tagName, tagOpts := parseTag(field.Tag.Get("json"))
		if tagName != "" {
			name = tagName
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
				if field.Anonymous {
					for k, v := range StructToMap(value) {
						result[k] = v
					}
					continue
				} else {
					value = StructToMap(value)
				}

			}
		}

		if tagOpts.Has("omitempty") {
			zero := reflect.Zero(val.Field(i).Type()).Interface()
			if reflect.DeepEqual(value, zero) {
				continue
			}
		}

		result[name] = value
	}

	return result
}
