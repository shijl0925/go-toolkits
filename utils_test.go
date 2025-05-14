package toolkits_test

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/shijl0925/go-toolkits"
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
			name:     "unsupported type",
			input:    struct{}{},
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

func TestSafeToInterface(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
		ok       bool
	}{
		{
			name:     "Invalid Value",
			input:    nil,
			expected: nil,
			ok:       false,
		},
		{
			name:     "Struct CanInterface",
			input:    struct{ X int }{X: 42},
			expected: struct{ X int }{X: 42},
			ok:       true,
		},
		{
			name:     "Bool",
			input:    true,
			expected: true,
			ok:       true,
		},
		{
			name:     "Int",
			input:    123,
			expected: 123,
			ok:       true,
		},
		{
			name:     "Uint",
			input:    uint(456),
			expected: uint(456),
			ok:       true,
		},
		{
			name:     "Float",
			input:    3.14,
			expected: 3.14,
			ok:       true,
		},
		{
			name:     "Complex",
			input:    complex(1, 2),
			expected: complex(1, 2),
			ok:       true,
		},
		{
			name:     "String",
			input:    "hello",
			expected: "hello",
			ok:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var val reflect.Value
			if tt.input == nil {
				val = reflect.Value{}
			} else {
				val = reflect.ValueOf(tt.input)
			}

			got, ok := toolkits.SafeToInterface(val)
			if ok != tt.ok {
				t.Errorf("expected ok=%v, got %v", tt.ok, ok)
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("expected value %v (%T), got %v (%T)", tt.expected, tt.expected, got, got)
			}
		})
	}
}
