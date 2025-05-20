package toolkits

import (
	"math"
	"reflect"
)

// EqualFloat64 比较两个float64是否相等, decimal为精度
//
// Example
//
//	EqualFloat64(1.000000000000001, 1.000000000000002, 15) // true
func EqualFloat64(a, b float64, decimal int) bool {
	// 校验 decimal 合法性
	if decimal < 0 || decimal > 15 {
		panic("decimal must be between 0 and 15 inclusive")
	}

	// 处理 NaN 和 Inf 的情况
	if math.IsNaN(a) || math.IsNaN(b) {
		return false
	}
	if math.IsInf(a, 0) || math.IsInf(b, 0) {
		return a == b
	}
	threshold := math.Pow10(-decimal)
	diff := math.Abs(a - b)
	// 相对误差 + 绝对误差结合判断
	epsilon := threshold * math.Max(math.Max(math.Abs(a), math.Abs(b)), 1.0)
	return diff <= epsilon
}

func isInt(value any) bool {
	intKinds := []reflect.Kind{
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
	}
	for _, kind := range intKinds {
		if IsKindOf(value, kind) {
			return true
		}
	}
	return false
}

func isNumeric(value any) bool {
	intKinds := []reflect.Kind{
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		// 可选：如果需要将复数也视为数值，取消注释以下两行
		// reflect.Complex64, reflect.Complex128
	}
	for _, kind := range intKinds {
		if IsKindOf(value, kind) {
			return true
		}
	}
	return false
}

// isDecimal 判断给定的值是否为十进制数值类型（目前仅支持 int 和 float）
func isDecimal(value any) bool {
	return isInt(value) || isFloat(value)
}

func isFloat(value any) bool {
	floatKinds := []reflect.Kind{
		reflect.Float32, reflect.Float64,
	}
	for _, kind := range floatKinds {
		if IsKindOf(value, kind) {
			return true
		}
	}
	return false
}

func isString(value any) bool {
	return IsKindOf(value, reflect.String)
}

func isMap(value any) bool {
	return IsKindOf(value, reflect.Map)
}

func isSlice(value any) bool {
	return IsKindOf(value, reflect.Slice)
}

func isBool(value any) bool {
	return IsKindOf(value, reflect.Bool)
}

func isStruct(value any) bool {
	return IsKindOf(value, reflect.Struct)
}

func isArray(value any) bool {
	return IsKindOf(value, reflect.Array)
}

func isPointer(value any) bool {
	return IsKindOf(value, reflect.Ptr)
}
