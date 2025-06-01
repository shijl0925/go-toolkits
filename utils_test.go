package toolkits_test

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	toolkits "github.com/shijl0925/go-toolkits"
	"math"
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

// 定义一个不可 marshal 的结构体（包含函数字段）
type UnmarshalableStruct struct {
	F func()
}

func TestSafeToString(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
		wantErr  bool
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: "",
			wantErr:  false,
		},
		{
			name:     "string",
			input:    "hello",
			expected: "hello",
			wantErr:  false,
		},
		{
			name:     "*string",
			input:    func() *string { s := "world"; return &s }(),
			expected: "world",
			wantErr:  false,
		},
		{
			name:     "[]byte",
			input:    []byte("bytes"),
			expected: "bytes",
			wantErr:  false,
		},
		{
			name:     "time.Duration",
			input:    time.Second * 5,
			expected: strconv.FormatInt(int64(time.Second*5), 10),
			wantErr:  false,
		},
		{
			name:     "int",
			input:    123,
			expected: "123",
			wantErr:  false,
		},
		{
			name:     "int64",
			input:    int64(987654321),
			expected: "987654321",
			wantErr:  false,
		},
		{
			name:     "uint",
			input:    uint(456),
			expected: "456",
			wantErr:  false,
		},
		{
			name:     "uint64",
			input:    uint64(1234567890),
			expected: "1234567890",
			wantErr:  false,
		},
		{
			name:     "float32",
			input:    float32(3.14),
			expected: "3.14",
			wantErr:  false,
		},
		{
			name:     "float64",
			input:    2.71828,
			expected: "2.71828",
			wantErr:  false,
		},
		{
			name:     "bool true",
			input:    true,
			expected: "true",
			wantErr:  false,
		},
		{
			name:     "bool false",
			input:    false,
			expected: "false",
			wantErr:  false,
		},
		{
			name:     "Stringer interface",
			input:    testStringer{val: "custom_string"},
			expected: "custom_string",
			wantErr:  false,
		},
		{
			name:     "error interface",
			input:    testError{msg: "an_error"},
			expected: "an_error",
			wantErr:  false,
		},
		{
			name:     "nil pointer",
			input:    (*int)(nil),
			expected: "",
			wantErr:  false,
		},
		{
			name:     "non-nil pointer to int",
			input:    func() *int { i := 789; return &i }(),
			expected: "789",
			wantErr:  false,
		},
		{
			name:     "struct type",
			input:    struct{}{},
			expected: "{}",
			wantErr:  false,
		},
		{
			name:     "unmarshalable struct",
			input:    UnmarshalableStruct{},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toolkits.SafeToString(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if got != tt.expected {
					t.Errorf("SafeToString(%v) = %q, want %q", tt.input, got, tt.expected)
				}
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

func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

func TestSafeToInt64(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected int64
		wantErr  bool
	}{
		{
			name:     "Nil input",
			input:    nil,
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "Nil pointer",
			input:    (*int)(nil),
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "non-nil pointer",
			input:    intPtr(42),
			expected: 42,
			wantErr:  false,
		},
		{
			name:     "Normal int",
			input:    42,
			expected: 42,
			wantErr:  false,
		},
		{
			name:     "uint",
			input:    uint(42),
			expected: 42,
			wantErr:  false,
		},
		{
			name:     "MaxInt64 value",
			input:    int64(math.MaxInt64),
			expected: math.MaxInt64,
			wantErr:  false,
		},
		{
			name:     "MinInt64 value",
			input:    int64(math.MinInt64),
			expected: math.MinInt64,
			wantErr:  false,
		},
		{
			name:     "Uint64 overflow",
			input:    uint64(1 << 63),
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "Integer float",
			input:    3.0,
			expected: 3,
			wantErr:  false,
		},
		{
			name:     "Complex number",
			input:    complex(3, 4),
			expected: 3,
			wantErr:  false,
		},
		{
			name:     "Parseable string",
			input:    "123",
			expected: 123,
			wantErr:  false,
		},
		{
			name:     "Unparseable string",
			input:    "abc",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "Boolean true",
			input:    true,
			expected: 1,
			wantErr:  false,
		},
		{
			name:     "Boolean false",
			input:    false,
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "Unsupported type",
			input:    struct{}{},
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := toolkits.SafeToInt64(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if result != tt.expected {
					t.Errorf("Expected %d but got %d", tt.expected, result)
				}
			}
		})
	}
}
func TestSafeToBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected []byte
		wantErr  error
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
			wantErr:  toolkits.ErrNilValue,
		},
		{
			name:     "nil pointer",
			input:    (*int)(nil),
			expected: nil,
			wantErr:  toolkits.ErrNilPointer,
		},
		{
			name:     "non-nil pointer",
			input:    stringPtr("hello"),
			expected: []byte("hello"),
			wantErr:  nil,
		},
		{
			name:     "string",
			input:    "hello",
			expected: []byte("hello"),
			wantErr:  nil,
		},
		{
			name:  "int",
			input: int(123),
			expected: func() []byte {
				buf := new(bytes.Buffer)
				_ = binary.Write(buf, binary.BigEndian, int64(123))
				return buf.Bytes()
			}(),
			wantErr: nil,
		},
		{
			name:  "uint",
			input: uint(123),
			expected: func() []byte {
				buf := new(bytes.Buffer)
				_ = binary.Write(buf, binary.BigEndian, uint64(123))
				return buf.Bytes()
			}(),
			wantErr: nil,
		},
		{
			name:  "float32",
			input: float32(123.45),
			expected: func() []byte {
				bits := math.Float32bits(float32(123.45))
				bs := make([]byte, 4)
				binary.BigEndian.PutUint32(bs, bits)
				return bs
			}(),
			wantErr: nil,
		},
		{
			name:  "float64",
			input: 123.45,
			expected: func() []byte {
				bits := math.Float64bits(123.45)
				bs := make([]byte, 8)
				binary.BigEndian.PutUint64(bs, bits)
				return bs
			}(),
			wantErr: nil,
		},
		{
			name:     "bool true",
			input:    true,
			expected: []byte("true"),
			wantErr:  nil,
		},
		{
			name:     "bool false",
			input:    false,
			expected: []byte("false"),
			wantErr:  nil,
		},
		{
			name:     "struct",
			input:    struct{ Name string }{Name: "test"},
			expected: []byte(`{"Name":"test"}`),
			wantErr:  nil,
		},
		{
			name:     "unmarshalable type (e.g., chan)",
			input:    make(chan int),
			expected: nil,
			wantErr:  &json.UnsupportedTypeError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toolkits.SafeToBytes(tt.input)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("expected error %v but got nil", tt.wantErr)
				} else if tt.wantErr != err && !reflect.TypeOf(tt.wantErr).AssignableTo(reflect.TypeOf(err)) {
					t.Errorf("expected error %v but got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !bytes.Equal(got, tt.expected) {
				t.Errorf("got %v, expected %v", got, tt.expected)
			}
		})
	}
}

// 测试无效值
func Test_SafeToInterface_Invalid(t *testing.T) {
	var v reflect.Value
	result, ok := toolkits.SafeToInterface(v)
	if result != nil || ok {
		t.Errorf("Expected (nil, false), got (%v, %v)", result, ok)
	}
}

// 测试 CanInterface 为 true 的情况
func Test_SafeToInterface_CanInterface(t *testing.T) {
	v := reflect.ValueOf(42)
	result, ok := toolkits.SafeToInterface(v)
	if !ok || result.(int) != 42 {
		t.Errorf("Expected (42, true), got (%v, %v)", result, ok)
	}
}

// 测试布尔值
func Test_SafeToInterface_Bool(t *testing.T) {
	v := reflect.ValueOf(true)
	result, ok := toolkits.SafeToInterface(v)
	if !ok || result.(bool) != true {
		t.Errorf("Expected (true, true), got (%v, %v)", result, ok)
	}
}

// 测试 Int 类型
func Test_SafeToInterface_Int(t *testing.T) {
	var i int64 = 123
	v := reflect.ValueOf(&i).Elem()
	result, ok := toolkits.SafeToInterface(v)
	if !ok || result.(int64) != 123 {
		t.Errorf("Expected (123, true), got (%v, %v)", result, ok)
	}
}

// 测试 Uint 类型
func Test_SafeToInterface_Uint(t *testing.T) {
	var u uint = 456
	v := reflect.ValueOf(&u).Elem()
	result, ok := toolkits.SafeToInterface(v)
	if !ok || result.(uint) != 456 {
		t.Errorf("Expected (456, true), got (%v, %v)", result, ok)
	}
}

// 测试 Float 类型
func Test_SafeToInterface_Float(t *testing.T) {
	v := reflect.ValueOf(3.14)
	result, ok := toolkits.SafeToInterface(v)
	if !ok || result.(float64) != 3.14 {
		t.Errorf("Expected (3.14, true), got (%v, %v)", result, ok)
	}
}

// 测试 Complex 类型
func Test_SafeToInterface_Complex(t *testing.T) {
	c := complex(1, 2)
	v := reflect.ValueOf(c)
	result, ok := toolkits.SafeToInterface(v)
	if !ok || result.(complex128) != c {
		t.Errorf("Expected (%v, true), got (%v, %v)", c, result, ok)
	}
}

// 测试字符串
func Test_SafeToInterface_String(t *testing.T) {
	v := reflect.ValueOf("hello")
	result, ok := toolkits.SafeToInterface(v)
	if !ok || result.(string) != "hello" {
		t.Errorf("Expected ('hello', true), got (%q, %v)", result, ok)
	}
}

// 测试指针类型
func Test_SafeToInterface_Ptr(t *testing.T) {
	s := "abc"
	v := reflect.ValueOf(&s).Elem()
	result, ok := toolkits.SafeToInterface(v)
	if !ok || result.(string) != "abc" {
		t.Errorf("Expected ('abc', true), got (%q, %v)", result, ok)
	}
}

// 测试接口类型
func Test_SafeToInterface_Interface(t *testing.T) {
	var i interface{} = "xyz"
	v := reflect.ValueOf(i)
	result, ok := toolkits.SafeToInterface(v)
	if !ok || result.(string) != "xyz" {
		t.Errorf("Expected ('xyz', true), got (%q, %v)", result, ok)
	}
}

// 测试结构体类型
func Test_SafeToInterface_Struct(t *testing.T) {
	type S struct{ X int }
	s := S{X: 10}
	v := reflect.ValueOf(s)
	result, ok := toolkits.SafeToInterface(v)
	if !ok || result != s {
		t.Errorf("Expected (nil, false), got (%v, %v)", result, ok)
	}
}

func TestSafeToUint64(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected uint64
		wantErr  bool
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "nil pointer",
			input:    (*int)(nil),
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "negative int",
			input:    -1,
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "positive int",
			input:    123,
			expected: 123,
			wantErr:  false,
		},
		{
			name:     "uint",
			input:    uint(456),
			expected: 456,
			wantErr:  false,
		},
		{
			name:     "float64 positive",
			input:    78.9,
			expected: 78,
			wantErr:  false,
		},
		{
			name:     "float64 negative",
			input:    -1.0,
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "complex128 positive real",
			input:    complex(100, 0),
			expected: 100,
			wantErr:  false,
		},
		{
			name:     "complex128 negative real",
			input:    complex(-1, 0),
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "string valid number",
			input:    "123",
			expected: 123,
			wantErr:  false,
		},
		{
			name:     "string invalid number",
			input:    "-123",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "string invalid",
			input:    "abc",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "bool true",
			input:    true,
			expected: 1,
			wantErr:  false,
		},
		{
			name:     "bool false",
			input:    false,
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "unsupported type slice",
			input:    []int{},
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "pointer to int",
			input:    new(int),
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "pointer to nil int",
			input:    func() *int { var v *int; return v }(),
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toolkits.SafeToUint64(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if got != tt.expected {
					t.Errorf("SafeToUint64(%v) = %q, want %q", tt.input, got, tt.expected)
				}
			}
		})
	}
}

func TestSafeToBool(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected bool
		wantErr  bool
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: false,
			wantErr:  true,
		},
		{
			name:     "nil pointer",
			input:    (*int)(nil),
			expected: false,
			wantErr:  true,
		},
		{
			name:     "pointer to int",
			input:    new(int),
			expected: false,
			wantErr:  false,
		},
		{
			name:     "bool true",
			input:    true,
			expected: true,
			wantErr:  false,
		},
		{
			name:     "bool false",
			input:    false,
			expected: false,
			wantErr:  false,
		},
		{
			name:     "int zero",
			input:    0,
			expected: false,
			wantErr:  false,
		},
		{
			name:     "int non-zero",
			input:    1,
			expected: true,
			wantErr:  false,
		},
		{
			name:     "uint zero",
			input:    uint(0),
			expected: false,
			wantErr:  false,
		},
		{
			name:     "uint non-zero",
			input:    uint(1),
			expected: true,
			wantErr:  false,
		},
		{
			name:     "float zero",
			input:    0.0,
			expected: false,
			wantErr:  false,
		},
		{
			name:     "float non-zero",
			input:    0.1,
			expected: true,
			wantErr:  false,
		},
		{
			name:     "complex real zero",
			input:    complex(0, 0),
			expected: false,
			wantErr:  false,
		},
		{
			name:     "complex real non-zero",
			input:    complex(0.1, 0),
			expected: true,
			wantErr:  false,
		},
		{
			name:     "string true",
			input:    "true",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "string false",
			input:    "false",
			expected: false,
			wantErr:  false,
		},
		{
			name:     "string empty",
			input:    "",
			expected: false,
			wantErr:  true,
		},
		{
			name:     "string non-empty",
			input:    "hello",
			expected: false,
			wantErr:  true,
		},
		{
			name:     "unsupported slice",
			input:    []int{},
			expected: false,
			wantErr:  true,
		},
		{
			name:     "unsupported struct",
			input:    struct{}{},
			expected: false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toolkits.SafeToBool(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if got != tt.expected {
					t.Errorf("SafeToBool(%v) = %v, want %v", tt.input, got, tt.expected)
				}
			}
		})
	}
}

func TestSafeToFloat64(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected float64
		wantErr  bool
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "nil pointer",
			input:    (*int)(nil),
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "int value",
			input:    int(123),
			expected: 123,
			wantErr:  false,
		},
		{
			name:     "uint value",
			input:    uint(456),
			expected: 456,
			wantErr:  false,
		},
		{
			name:     "float32 value",
			input:    float32(3.14),
			expected: float64(float32(3.14)),
			wantErr:  false,
		},
		{
			name:     "float64 value",
			input:    3.14,
			expected: 3.14,
			wantErr:  false,
		},
		{
			name:     "complex64 value",
			input:    complex(100, 0),
			expected: 100,
			wantErr:  false,
		},
		{
			name:     "string numeric",
			input:    "123.45",
			expected: 123.45,
			wantErr:  false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "invalid string",
			input:    "abc",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "bool true",
			input:    true,
			expected: 1,
			wantErr:  false,
		},
		{
			name:     "bool false",
			input:    false,
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "unsupported slice type",
			input:    []int{},
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "non-nil pointer to int",
			input:    func() *int { i := 10; return &i }(),
			expected: 10,
			wantErr:  false,
		},
		{
			name:     "interface with string",
			input:    interface{}("123"),
			expected: 123,
			wantErr:  false,
		},
		{
			name:     "interface with int",
			input:    interface{}(123),
			expected: 123,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toolkits.SafeToFloat64(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if got != tt.expected {
					t.Errorf("SafeToFloat64(%v) = %v, want %v", tt.input, got, tt.expected)
				}
			}
		})
	}
}

// 测试 IsNil
func TestIsNil(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{
			name:     "Interface is nil",
			input:    nil,
			expected: true,
		},
		{
			name:     "Nil Pointer",
			input:    (*int)(nil),
			expected: true,
		},
		{
			name:     "Initialized Pointer",
			input:    new(int),
			expected: false,
		},
		{
			name:     "Nil Slice",
			input:    []int(nil),
			expected: true,
		},
		{
			name:     "Initialized Slice",
			input:    []int{1, 2},
			expected: false,
		},
		{
			name:     "Nil Map",
			input:    map[string]int(nil),
			expected: true,
		},
		{
			name:     "Initialized Map",
			input:    map[string]int{"a": 1},
			expected: false,
		},
		{
			name:     "Nil Chan",
			input:    (chan int)(nil),
			expected: true,
		},
		{
			name:     "Initialized Chan",
			input:    make(chan int),
			expected: false,
		},
		{
			name:     "Nil Func",
			input:    (func())(nil),
			expected: true,
		},
		{
			name:     "Initialized Func",
			input:    func() {},
			expected: false,
		},
		{
			name:     "Int Value",
			input:    42,
			expected: false,
		},
		{
			name:     "Struct Value",
			input:    struct{}{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toolkits.IsNil(tt.input)
			if got != tt.expected {
				t.Errorf("IsNil(%v) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}
