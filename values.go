package toolkits

import (
	"fmt"
	"math"
	"reflect"
)

func AnyToInt(v any) (int, error) {
	i, err := SafeToInt64(v)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

func AnyToUint(v any) (uint, error) {
	i, err := SafeToUint64(v)
	if err != nil {
		return 0, err
	}
	return uint(i), nil
}

func AnyToInt8(v any) (int8, error) {
	i, err := SafeToInt64(v)
	if err != nil {
		return 0, err
	}
	if i < math.MinInt8 || i > math.MaxInt8 {
		return 0, fmt.Errorf("value overflow int8. input: %v", i)
	}
	return int8(i), nil
}

func AnyToUint8(v any) (uint8, error) {
	i, err := SafeToUint64(v)
	if err != nil {
		return 0, err
	}
	if i > math.MaxUint8 {
		return 0, fmt.Errorf("value overflow uint8. input: %v", i)
	}
	return uint8(i), nil
}

func AnyToInt16(v any) (int16, error) {
	i, err := SafeToInt64(v)
	if err != nil {
		return 0, err
	}
	if i < math.MinInt16 || i > math.MaxInt16 {
		return 0, fmt.Errorf("value overflow int16. input: %v", i)
	}
	return int16(i), nil
}

func AnyToUint16(v any) (uint16, error) {
	i, err := SafeToUint64(v)
	if err != nil {
		return 0, err
	}
	if i > math.MaxUint16 {
		return 0, fmt.Errorf("value overflow uint16. input: %v", i)
	}
	return uint16(i), nil
}

func AnyToInt32(v any) (int32, error) {
	i, err := SafeToInt64(v)
	if err != nil {
		return 0, err
	}
	if i < math.MinInt32 || i > math.MaxInt32 {
		return 0, fmt.Errorf("value overflow int32. input: %v", i)
	}
	return int32(i), nil
}

func AnyToUint32(v any) (uint32, error) {
	i, err := SafeToUint64(v)
	if err != nil {
		return 0, err
	}
	if i > math.MaxUint32 {
		return 0, fmt.Errorf("value overflow uint32. input: %v", i)
	}
	return uint32(i), nil
}

func AnyToInt64(v any) (int64, error) {
	return SafeToInt64(v)
}

func AnyToUint64(v any) (uint64, error) {
	return SafeToUint64(v)
}

func AnyToFloat32(v any) (float32, error) {
	f, err := SafeToFloat64(v)
	if err != nil {
		return 0, err
	}
	if f < -math.MaxFloat32 || f > math.MaxFloat32 {
		return 0, fmt.Errorf("value overflow float32. input: %v", f)
	}
	return float32(f), nil
}

func AnyToFloat64(v any) (float64, error) {
	return SafeToFloat64(v)
}

func AnyToString(v any) (string, error) {
	return SafeToString(v)
}

func AnyToBool(v any) (bool, error) {
	return SafeToBool(v)
}

func AnyToBytes(v any) ([]byte, error) {
	s, err := SafeToString(v)
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}

func AnyToInterface(value any) (interface{}, bool) {
	return SafeToInterface(reflect.ValueOf(value))
}

// Cast Converts any(interface{}) type to a given type. If conversion fails, it returns the zero value
// of the given type. This is a type-safe version of type assertion.
func Cast[T any](value any) (T, bool) {
	if v, ok := value.(T); ok {
		return v, true
	}
	var zero T
	return zero, false
}
