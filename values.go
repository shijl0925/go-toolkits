package toolkits

import "reflect"

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
