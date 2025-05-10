package toolkits

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func SafeString(v any) (str string, err error) {
	if v == nil {
		return "", nil
	}
	switch value := v.(type) {
	case string:
		str = value
	case *string:
		str = *value
	case []byte:
		str = string(value)
	case time.Duration:
		str = strconv.FormatInt(int64(value), 10)
	case fmt.Stringer:
		str = value.String()
	case error:
		str = value.Error()
	default:
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				return "", nil
			}
			return SafeString(rv.Elem().Interface())
		}

		switch rv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			str = strconv.FormatInt(rv.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			str = strconv.FormatUint(rv.Uint(), 10)
		case reflect.Float32:
			str = strconv.FormatFloat(rv.Float(), 'f', -1, 32)
		case reflect.Float64:
			str = strconv.FormatFloat(rv.Float(), 'f', -1, 64)
		case reflect.Bool:
			str = strconv.FormatBool(rv.Bool())
		default:
			return "", fmt.Errorf("unsupported type: %T", v)
		}
	}
	return
}
