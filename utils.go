package toolkits

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"
)

var (
	// ErrNilValue indicates that a value is nil.
	ErrNilValue = fmt.Errorf("nil value")

	// ErrNilPointer indicates that a pointer is nil.
	ErrNilPointer = fmt.Errorf("nil pointer")

	// ErrType indicates that a type is not supported.
	ErrType = fmt.Errorf("unsupported type")

	// ErrValOut indicates that a value is out of range.
	ErrValOut = fmt.Errorf("value out of range")

	// ErrUnsignedInt indicates that a negative value is attempted to be converted to an unsigned integer.
	ErrUnsignedInt = fmt.Errorf("cannot convert negative value to unsigned integer")
)

func SafeToString(v any) (str string, err error) {
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
			return SafeToString(rv.Elem().Interface())
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
			return "", ErrType
		}
	}
	return
}

func SafeToInt64(v any) (int64, error) {
	if v == nil {
		return 0, ErrNilValue
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return 0, ErrNilPointer
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := val.Int()
		return v, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v := val.Uint()
		return int64(v), nil
	case reflect.Float32, reflect.Float64:
		v := val.Float()
		return int64(v), nil
	case reflect.Complex64, reflect.Complex128:
		v := val.Complex()
		return int64(real(v)), nil
	case reflect.String:
		s := val.String()
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, strconv.ErrSyntax
		}
		return v, nil
	case reflect.Bool:
		{
			if val.Bool() {
				return 1, nil
			} else {
				return 0, nil
			}
		}
	default:
		return 0, ErrType
	}
}

func SafeToUint64(v any) (uint64, error) {
	if v == nil {
		return 0, ErrNilValue
	}
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return 0, ErrNilPointer
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := val.Int()
		if v < 0 {
			return 0, ErrUnsignedInt
		}
		return uint64(v), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val.Uint(), nil
	case reflect.Float32, reflect.Float64:
		v := val.Float()
		if v < 0 {
			return 0, ErrUnsignedInt
		}
		return uint64(v), nil
	case reflect.Complex64, reflect.Complex128:
		v := real(val.Complex())
		if v < 0 {
			return 0, ErrUnsignedInt
		}
		return uint64(v), nil
	case reflect.String:
		s := val.String()
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return 0, strconv.ErrSyntax
		}
		return v, nil
	case reflect.Bool:
		{
			if val.Bool() {
				return 1, nil
			} else {
				return 0, nil
			}
		}
	default:
		return 0, ErrType
	}
}

func SafeToBool(v any) (bool, error) {
	if v == nil {
		return false, ErrNilValue
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return false, ErrNilPointer
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
		return false, ErrType
	}
}

func SafeToFloat64(v any) (float64, error) {
	if v == nil {
		return 0, ErrNilValue
	}
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return 0, ErrNilPointer
		}
		val = val.Elem()
	}
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(val.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(val.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return val.Float(), nil
	case reflect.Complex64, reflect.Complex128:
		return real(val.Complex()), nil
	case reflect.String:
		v := val.String()
		if v == "" {
			return 0, nil
		}
		floatValue, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, err
		}
		return floatValue, nil
	case reflect.Bool:
		if val.Bool() {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, ErrType
	}
}


func SafeToBytes(v any) ([]byte, error) {
	if v == nil {
		return nil, ErrNilValue
	}

	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, ErrNilPointer
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.String:
		return []byte(val.String()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		number := val.Int()
		buf := bytes.NewBuffer(nil)
		err := binary.Write(buf, binary.BigEndian, number)
		return buf.Bytes(), err
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		number := val.Uint()
		buf := bytes.NewBuffer(nil)
		err := binary.Write(buf, binary.BigEndian, number)
		return buf.Bytes(), err
	case reflect.Float32:
		number := float32(val.Float())
		bits := math.Float32bits(number)
		bs := make([]byte, 4)
		binary.BigEndian.PutUint32(bs, bits)
		return bs, nil
	case reflect.Float64:
		number := val.Float()
		bits := math.Float64bits(number)
		bs := make([]byte, 8)
		binary.BigEndian.PutUint64(bs, bits)
		return bs, nil
	case reflect.Bool:
		return strconv.AppendBool([]byte{}, val.Bool()), nil
	default:
		return json.Marshal(v)
	}
}

// SafeToInterface converts reflect value to its interface type.
func SafeToInterface(v reflect.Value) (interface{}, bool) {
	if !v.IsValid() {
		return nil, false
	}

	if v.CanInterface() {
		return v.Interface(), true
	}

	switch v.Kind() {
	case reflect.Bool:
		return v.Bool(), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int(), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint(), true
	case reflect.Float32, reflect.Float64:
		return v.Float(), true
	case reflect.Complex64, reflect.Complex128:
		return v.Complex(), true
	case reflect.String:
		return v.String(), true
	case reflect.Ptr, reflect.Interface:
		return SafeToInterface(v.Elem())
	default:
		return nil, false
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
