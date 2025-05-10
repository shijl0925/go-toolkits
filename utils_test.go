package toolkits_test

import (
	"github.com/shijl0925/go-toolkits"
	"strconv"
	"testing"
	"time"
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
