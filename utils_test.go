package toolkits_test

import (
	"github.com/shijl0925/go-toolkits"
	"reflect"
	"strconv"
	"testing"
	"time"
	"unsafe"
)

// 实现 Stringer 接口的测试结构体
type testStringer struct {
	val string
}

func (t testStringer) String() string {
	return t.val
}

// 实现 error 接口的测试结构体
type testError struct {
	msg string
}

func (t testError) Error() string {
	return t.msg
}

func TestSafeString(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: "",
		},
		{
			name:     "string",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "*string",
			input:    func() *string { s := "world"; return &s }(),
			expected: "world",
		},
		{
			name:     "[]byte",
			input:    []byte("bytes"),
			expected: "bytes",
		},
		{
			name:     "time.Duration",
			input:    time.Second * 5,
			expected: strconv.FormatInt(int64(time.Second*5), 10),
		},
		{
			name:     "int",
			input:    123,
			expected: "123",
		},
		{
			name:     "int64",
			input:    int64(987654321),
			expected: "987654321",
		},
		{
			name:     "uint",
			input:    uint(456),
			expected: "456",
		},
		{
			name:     "uint64",
			input:    uint64(1234567890),
			expected: "1234567890",
		},
		{
			name:     "float32",
			input:    float32(3.14),
			expected: "3.14",
		},
		{
			name:     "float64",
			input:    2.71828,
			expected: "2.71828",
		},
		{
			name:     "bool true",
			input:    true,
			expected: "true",
		},
		{
			name:     "bool false",
			input:    false,
			expected: "false",
		},
		{
			name:     "Stringer interface",
			input:    testStringer{val: "custom_string"},
			expected: "custom_string",
		},
		{
			name:     "error interface",
			input:    testError{msg: "an_error"},
			expected: "an_error",
		},
		{
			name:     "nil pointer",
			input:    (*int)(nil),
			expected: "",
		},
		{
			name:     "non-nil pointer to int",
			input:    func() *int { i := 789; return &i }(),
			expected: "789",
		},
		//{
		//	name:     "unsupported type",
		//	input:    struct{}{},
		//	expected: "",
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toolkits.SafeString(tt.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("SafeString(%v) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// isKindOf 函数用于判断任意值的底层类型是否为指定的 reflect.Kind。
func TestIsKindOf(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		kind     reflect.Kind
		expected bool
	}{
		{
			name:     "TC01 - Input is nil",
			input:    nil,
			kind:     reflect.Invalid, // 任意 kind
			expected: false,
		},
		{
			name:     "TC02 - Uninitialized interface",
			input:    func() any { var v any; return v }(),
			kind:     reflect.Invalid, // 任意 kind
			expected: false,
		},
		{
			name:     "TC03 - int type matches Int kind",
			input:    int(5),
			kind:     reflect.Int,
			expected: true,
		},
		{
			name:     "TC04 - string type matches String kind",
			input:    "abc",
			kind:     reflect.String,
			expected: true,
		},
		{
			name:     "TC05 - slice type matches Slice kind",
			input:    []int{},
			kind:     reflect.Slice,
			expected: true,
		},
		{
			name:     "TC06 - map type matches Map kind",
			input:    map[string]int{},
			kind:     reflect.Map,
			expected: true,
		},
		{
			name:     "TC07 - int does not match String kind",
			input:    int(5),
			kind:     reflect.String,
			expected: false,
		},
		{
			name:     "TC08 - struct type matches Struct kind",
			input:    struct{}{},
			kind:     reflect.Struct,
			expected: true,
		},
		{
			name:     "TC09 - Pointer type matches Pointer kind",
			input:    &struct{}{},
			kind:     reflect.Pointer,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toolkits.IsKindOf(tt.input, tt.kind)
			if result != tt.expected {
				t.Errorf("IsKindOf(%v, %v) = %v; want %v", tt.input, tt.kind, result, tt.expected)
			}
		})
	}
}

func TestIsNilValue(t *testing.T) {
	type MyInterface interface{}
	// 空接口，用于测试
	var nilMI MyInterface
	// 非空接口，用于测试
	n := 1
	myInterface := MyInterface(n)

	// nil channel
	var nilCh chan int
	// 非 nil channel
	ch := make(chan int)

	// nil pointer
	var nilPtr *int
	// 非 nil pointer
	ptr := &n

	// nil unsafePointer
	nilUsPtr := unsafe.Pointer(nilPtr)
	// 非 nil unsafePointer
	usPtr := unsafe.Pointer(ptr)

	// nil map
	var nilMp map[string]struct{}
	// 非 nil map
	mp := make(map[string]struct{}, 1)

	// nil slice
	var nilSlice []int
	// 非 nil slice
	slc := make([]int, 1)

	// nil 函数
	type MyFunc func()
	var myFunc MyFunc

	testCases := []struct {
		name string
		val  reflect.Value
		res  bool
	}{
		{
			name: "int 类型",
			val:  reflect.ValueOf(666),
			res:  false,
		},
		{
			name: "string 类型",
			val:  reflect.ValueOf("字符串类型"),
			res:  false,
		},
		{
			name: "bool 类型",
			val:  reflect.ValueOf(true),
			res:  false,
		},
		{
			name: "float 类型",
			val:  reflect.ValueOf(3.14),
			res:  false,
		},
		{
			name: "complex 类型",
			val:  reflect.ValueOf(complex(1, 1)),
			res:  false,
		},
		{
			name: "struct 类型",
			val:  reflect.ValueOf(struct{}{}),
			res:  false,
		},
		{
			name: "array 类型",
			val:  reflect.ValueOf([4]int{}),
			res:  false,
		},
		{
			name: "nil 非法值",
			val:  reflect.ValueOf(nil),
			res:  true,
		},
		{
			name: "interface 类型 - 非空",
			val:  reflect.ValueOf(myInterface),
			res:  false,
		},
		{
			name: "interface 类型 - 空",
			val:  reflect.ValueOf(nilMI),
			res:  true,
		},
		{
			name: "pointer 类型 - 非空",
			val:  reflect.ValueOf(ptr),
			res:  false,
		},
		{
			name: "pointer 类型 - 空",
			val:  reflect.ValueOf(nilPtr),
			res:  true,
		},
		{
			name: "unsafePointer 类型 - 非空",
			val:  reflect.ValueOf(usPtr),
			res:  false,
		},
		{
			name: "unsafePointer 类型 - 空",
			val:  reflect.ValueOf(nilUsPtr),
			res:  true,
		},
		{
			name: "channel 类型 - 非空",
			val:  reflect.ValueOf(ch),
			res:  false,
		},
		{
			name: "channel 类型 - 空",
			val:  reflect.ValueOf(nilCh),
			res:  true,
		},
		{
			name: "map 类型 - 非空",
			val:  reflect.ValueOf(mp),
			res:  false,
		},
		{
			name: "map 类型 - 空",
			val:  reflect.ValueOf(nilMp),
			res:  true,
		},
		{
			name: "slice 类型 - 非空",
			val:  reflect.ValueOf(slc),
			res:  false,
		},
		{
			name: "slice 类型 - 空",
			val:  reflect.ValueOf(nilSlice),
			res:  true,
		},
		{
			name: "func 类型 - 非空",
			val: reflect.ValueOf(func() func() {
				return func() {}
			}),
			res: false,
		},
		{
			name: "func 类型 - 空",
			val:  reflect.ValueOf(myFunc),
			res:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := toolkits.IsNilValue(tc.val)
			if res != tc.res {
				t.Errorf("IsNilValue(%v) = %v; want %v", tc.val, res, tc.res)
			}
		})
	}
}
