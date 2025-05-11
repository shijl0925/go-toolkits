package stringx_test

import (
	"github.com/shijl0925/go-toolkits/stringx"
	"testing"
)

// TestCapitalize_NormalCases 测试正常输入情况
func TestCapitalize_NormalCases(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "Hello"},
		{"HELLO", "HELLO"},
		{"hELLo", "HELLo"},
		{"a", "A"},
		{"123abc", "123abc"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := stringx.Capitalize(tt.input)
			if result != tt.expected {
				t.Errorf("Capitalize(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestCapitalize_EmptyString 测试空字符串是否引发 panic
func TestCapitalize_EmptyString(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for empty string input")
		}
	}()
	stringx.Capitalize("")
}
