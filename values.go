package toolkits

import (
	"fmt"
	"strconv"
)

func Int(v any) (int, error) {
	val, ok := v.(int)
	if !ok {
		return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "int", v)
	}
	return val, nil
}

func AsInt(v any) (int, error) {
	switch v := v.(type) {
	case int:
		return v, nil
	case string:
		res, err := strconv.ParseInt(v, 10, 64)
		return int(res), err
	}
	return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "int", v)
}

func Uint(v any) (uint, error) {
	val, ok := v.(uint)
	if !ok {
		return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "uint", v)
	}
	return val, nil
}

func AsUint(v any) (uint, error) {
	switch v := v.(type) {
	case uint:
		return v, nil
	case string:
		res, err := strconv.ParseUint(v, 10, 64)
		return uint(res), err
	}
	return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "uint", v)
}

func Int8(v any) (int8, error) {
	val, ok := v.(int8)
	if !ok {
		return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "int8", v)
	}
	return val, nil
}

func AsInt8(v any) (int8, error) {
	switch v := v.(type) {
	case int8:
		return v, nil
	case string:
		res, err := strconv.ParseInt(v, 10, 64)
		return int8(res), err
	}
	return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "int8", v)
}

func Uint8(v any) (uint8, error) {
	val, ok := v.(uint8)
	if !ok {
		return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "uint8", v)
	}
	return val, nil
}

func AsUint8(v any) (uint8, error) {
	switch v := v.(type) {
	case uint8:
		return v, nil
	case string:
		res, err := strconv.ParseUint(v, 10, 8)
		return uint8(res), err
	}
	return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "uint8", v)
}

func Int16(v any) (int16, error) {
	val, ok := v.(int16)
	if !ok {
		return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "int16", v)
	}
	return val, nil
}

func AsInt16(v any) (int16, error) {
	switch v := v.(type) {
	case int16:
		return v, nil
	case string:
		res, err := strconv.ParseInt(v, 10, 16)
		return int16(res), err
	}
	return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "int16", v)
}

func Uint16(v any) (uint16, error) {
	val, ok := v.(uint16)
	if !ok {
		return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "uint16", v)
	}
	return val, nil
}

func AsUint16(v any) (uint16, error) {
	switch v := v.(type) {
	case uint16:
		return v, nil
	case string:
		res, err := strconv.ParseUint(v, 10, 16)
		return uint16(res), err
	}
	return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "uint16", v)
}

func Int32(v any) (int32, error) {
	val, ok := v.(int32)
	if !ok {
		return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "int32", v)
	}
	return val, nil
}

func AsInt32(v any) (int32, error) {
	switch v := v.(type) {
	case int32:
		return v, nil
	case string:
		res, err := strconv.ParseInt(v, 10, 32)
		return int32(res), err
	}
	return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "int32", v)
}

func Uint32(v any) (uint32, error) {
	val, ok := v.(uint32)
	if !ok {
		return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "uint32", v)
	}
	return val, nil
}

func AsUint32(v any) (uint32, error) {
	switch v := v.(type) {
	case uint32:
		return v, nil
	case string:
		res, err := strconv.ParseUint(v, 10, 32)
		return uint32(res), err
	}
	return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "uint32", v)
}

func Int64(v any) (int64, error) {
	val, ok := v.(int64)
	if !ok {
		return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "int64", v)
	}
	return val, nil
}

func AsInt64(v any) (int64, error) {
	switch v := v.(type) {
	case int64:
		return v, nil
	case string:
		res, err := strconv.ParseInt(v, 10, 64)
		return res, err
	}
	return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "int64", v)
}

func Uint64(v any) (uint64, error) {
	val, ok := v.(uint64)
	if !ok {
		return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "uint64", v)
	}
	return val, nil
}

func AsUint64(v any) (uint64, error) {
	switch v := v.(type) {
	case uint64:
		return v, nil
	case string:
		res, err := strconv.ParseUint(v, 10, 64)
		return res, err
	}
	return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "uint64", v)
}

func Float32(v any) (float32, error) {
	val, ok := v.(float32)
	if !ok {
		return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "float32", v)
	}
	return val, nil
}

func AsFloat32(v any) (float32, error) {
	switch v := v.(type) {
	case float32:
		return v, nil
	case string:
		res, err := strconv.ParseFloat(v, 32)
		return float32(res), err
	}
	return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "float32", v)
}

func Float64(v any) (float64, error) {
	val, ok := v.(float64)
	if !ok {
		return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "float64", v)
	}
	return val, nil
}

func AsFloat64(v any) (float64, error) {
	switch v := v.(type) {
	case float64:
		return v, nil
	case string:
		res, err := strconv.ParseFloat(v, 64)
		return res, err
	}
	return 0, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "float64", v)
}

func String(v any) (string, error) {
	val, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "string", v)
	}
	return val, nil
}

func AsString(v any) (string, error) {
	res, err := SafeToString(v)
	if err != nil {
		return "", fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "string", v)
	}
	return res, nil
}

func Bool(v any) (bool, error) {
	val, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "bool", v)
	}
	return val, nil
}

func Bytes(v any) ([]byte, error) {
	val, ok := v.([]byte)
	if !ok {
		return nil, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "[]byte", v)
	}
	return val, nil
}

func AsBytes(v any) ([]byte, error) {
	switch v := v.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	}
	return nil, fmt.Errorf("类型转换失败，预期类型:%s, 实际值:%#v", "[]byte", v)
}
