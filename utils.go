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

func SafeToInt(v any) (int, error) {
	if v == nil {
		return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:nil", "int")
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:nil", "int")
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := val.Int()

		if v < math.MinInt64 || v > math.MaxInt64 {
			return 0, fmt.Errorf("类型转换失败，预期类型:%s, 值超出范围:%d", "int", v)
		}

		return int(v), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v := val.Uint()

		if v > uint64(math.MaxInt64) {
			return 0, fmt.Errorf("类型转换失败，预期类型:%s, 值超出范围:%d", "int", v)
		}

		return int(v), nil
	case reflect.Float32, reflect.Float64:
		v := val.Float()

		if math.Abs(v-math.Trunc(v)) > 1e-9 {
			return 0, fmt.Errorf("类型转换失败，预期类型:%s, 值精度丢失:%f", "int", v)
		}

		if v < float64(math.MinInt64) || v > float64(math.MaxInt64) {
			return 0, fmt.Errorf("类型转换失败，预期类型:%s, 值超出范围:%f", "int", v)
		}

		return int(v), nil
	case reflect.Complex64, reflect.Complex128:
		v := val.Complex()
		return int(real(v)), nil
	case reflect.String:
		s := val.String()

		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "int", v)
		}

		if v < math.MinInt64 || v > math.MaxInt64 {
			return 0, fmt.Errorf("类型转换失败，预期类型:%s, 值超出范围:%d", "int", v)
		}

		return int(v), nil
	case reflect.Bool:
		{
			if val.Bool() {
				return 1, nil
			} else {
				return 0, nil
			}
		}
	default:
		return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "int", v)
	}
}

func SafeToBool(v any) (bool, error) {
	if v == nil {
		return false, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:nil", "bool")
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return false, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:nil", "bool")
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Bool:
		return val.Bool(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() != 0, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val.Uint() != 0, nil
	case reflect.Float32, reflect.Float64:
		return val.Float() != 0, nil
	case reflect.Complex64, reflect.Complex128:
		return real(val.Complex()) != 0, nil
	case reflect.String:
		s := val.String()
		if s == "true" {
			return true, nil
		} else if s == "false" {
			return false, nil
		}
		return s != "", nil
	default:
		return false, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "bool", v)
	}
}

// IsKindOf checks if the given value is an instance of the given kind.
// It returns true if the value is an instance of the given kind, false otherwise.
func IsKindOf(s any, t reflect.Kind) bool {
	if s == nil {
		return false
	}

	v := reflect.TypeOf(s)
	if v == nil {
		return false
	}

	return v.Kind() == t
}

func IsNilValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}

	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return v.IsNil()
	default:
		return false
	}
}
