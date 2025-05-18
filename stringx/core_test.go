package stringx_test

import (
	"github.com/shijl0925/go-toolkits/stringx"
	"reflect"
	"testing"
)

// TestCapitalize_NormalCases 测试正常输入情况
func TestCapitalize_NormalCases(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "Hello"},
		{"HELLO", "Hello"},
		{"hELLo", "Hello"},
		{"a", "A"},
		{"123abc", "123abc"},
		{"", ""},
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

func TestPartition(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		sep      string
		expected [3]string
		panicMsg string
	}{
		{
			name:     "Normal case with hyphen",
			s:        "hello-world",
			sep:      "-",
			expected: [3]string{"hello", "-", "world"},
		},
		{
			name:     "Normal case with space",
			s:        "hello world",
			sep:      " ",
			expected: [3]string{"hello", " ", "world"},
		},
		{
			name:     "Separator not found",
			s:        "hello",
			sep:      "-",
			expected: [3]string{"hello", "", ""},
		},
		{
			name:     "Empty input string",
			s:        "",
			sep:      "-",
			expected: [3]string{"", "", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			first, second, third := stringx.Partition(tt.s, tt.sep)
			if !reflect.DeepEqual([3]string{first, second, third}, tt.expected) {
				t.Errorf("Partition(%q, %q) = %v, want %v", tt.s, tt.sep, [3]string{first, second, third}, tt.expected)
			}
		})
	}
}

func TestPartition_InvalidSep(t *testing.T) {
	t.Run("Invalid sep", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic due to invalid sep, but did not get one")
			}
		}()

		stringx.Partition("abc", "")
	})
}

func TestRightPartition(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		sep      string
		expected [3]string
		panicMsg string
	}{
		{
			name:     "Normal case with multiple separators",
			s:        "a/b/c",
			sep:      "/",
			expected: [3]string{"a/b", "/", "c"},
		},
		{
			name:     "Separator not found",
			s:        "abc",
			sep:      "x",
			expected: [3]string{"", "", "abc"},
		},
		{
			name:     "Single separator at middle",
			s:        "abc/def",
			sep:      "/",
			expected: [3]string{"abc", "/", "def"},
		},
		{
			name:     "Empty input string",
			s:        "",
			sep:      "/",
			expected: [3]string{"", "", ""},
		},
		{
			name:     "Multi-character separator",
			s:        "a/b/c/d",
			sep:      "b/c",
			expected: [3]string{"a/", "b/c", "/d"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			first, second, third := stringx.RightPartition(tt.s, tt.sep)
			if !reflect.DeepEqual([3]string{first, second, third}, tt.expected) {
				t.Errorf("RightPartition(%q, %q) = %v, want %v", tt.s, tt.sep, [3]string{first, second, third}, tt.expected)
			}
		})
	}
}

func TestRightPartition_InvalidSep(t *testing.T) {
	t.Run("Invalid sep", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic due to invalid sep, but did not get one")
			}
		}()

		stringx.RightPartition("abc", "")
	})
}

// TestSwapCase runs unit tests for the SwapCase function.
func TestSwapCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "HELLO"},
		{"WORLD", "world"},
		{"GoLang", "gOlANG"},
		{"123!@#", "123!@#"},
		{"", ""},
		{"AbcDef123", "aBCdEF123"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := stringx.SwapCase(test.input)
			if result != test.expected {
				t.Errorf("SwapCase(%q) = %q; expected %q", test.input, result, test.expected)
			}
		})
	}
}

// Test cases for FormatMap function
func TestFormatMap(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		inputMap map[string]any
		expected string
	}{
		{
			name:     "Single placeholder replaced",
			format:   "Hello {name}",
			inputMap: map[string]any{"name": "Alice"},
			expected: "Hello Alice",
		},
		{
			name:     "No placeholders",
			format:   "No placeholder",
			inputMap: nil,
			expected: "No placeholder",
		},
		{
			name:     "Multiple placeholders replaced",
			format:   "{a} and {b}",
			inputMap: map[string]any{"a": "apple", "b": "banana"},
			expected: "apple and banana",
		},
		{
			name:     "Some keys missing",
			format:   "{a} and {b}",
			inputMap: map[string]any{"a": "apple"},
			expected: "apple and {b}",
		},
		{
			name:     "Unclosed placeholder",
			format:   "Unclosed {key",
			inputMap: map[string]any{"key": "value"},
			expected: "Unclosed {key",
		},
		{
			name:     "Empty replacement value",
			format:   "{empty}",
			inputMap: map[string]any{"empty": ""},
			expected: "",
		},
		{
			name:     "Multiple with some missing",
			format:   "{k1} {k2} {k3}",
			inputMap: map[string]any{"k1": "v1", "k3": "v3"},
			expected: "v1 {k2} v3",
		},
		{
			name:     "Escaped closing bracket",
			format:   "{key}}",
			inputMap: map[string]any{"key": "value"},
			expected: "value}",
		},
		{
			name:     "Empty input string",
			format:   "",
			inputMap: nil,
			expected: "",
		},
		{
			name:     "Nested braces",
			format:   "{{escaped}}",
			inputMap: map[string]any{"escaped": "value"},
			expected: "{{escaped}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stringx.FormatMap(tt.format, tt.inputMap)
			if result != tt.expected {
				t.Errorf("FormatMap(%q, %v) = %q; want %q", tt.format, tt.inputMap, result, tt.expected)
			}
		})
	}
}

// TestReverse_EmptyString tests reversing an empty string.
func TestReverse_EmptyString(t *testing.T) {
	input := ""
	expected := ""
	result := stringx.Reverse(input)
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestReverse_SingleCharacter tests reversing a single character string.
func TestReverse_SingleCharacter(t *testing.T) {
	input := "a"
	expected := "a"
	result := stringx.Reverse(input)
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestReverse_EvenLengthString tests reversing a string with even length.
func TestReverse_EvenLengthString(t *testing.T) {
	input := "ab"
	expected := "ba"
	result := stringx.Reverse(input)
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestReverse_OddLengthString tests reversing a string with odd length.
func TestReverse_OddLengthString(t *testing.T) {
	input := "abc"
	expected := "cba"
	result := stringx.Reverse(input)
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestReverse_NormalString tests reversing a normal English string.
func TestReverse_NormalString(t *testing.T) {
	input := "hello"
	expected := "olleh"
	result := stringx.Reverse(input)
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestReverse_UnicodeString tests reversing a string containing Unicode characters.
func TestReverse_UnicodeString(t *testing.T) {
	input := "你好"
	expected := "好你"
	result := stringx.Reverse(input)
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestReverse_NumericString tests reversing a numeric string.
func TestReverse_NumericString(t *testing.T) {
	input := "123456"
	expected := "654321"
	result := stringx.Reverse(input)
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestSubstring 使用多种边界情况和典型情况来验证 Substring 函数的行为
func TestSubstring(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		offset   int
		length   int
		expected string
	}{
		{"TC00: Normal case", "abcdef", 2, -1, ""},
		{"TC00: Normal case", "abcdef", 2, 0, ""},
		{"TC01: Normal case", "abcdef", 2, 3, "cde"},
		{"TC02: Negative offset", "abcdef", -2, 2, "ef"},
		{"TC03: Offset too small", "abcdef", -10, 2, ""},
		{"TC04: Length exceeds", "abcdef", 4, 10, "ef"},
		{"TC05: Offset equals length", "abcdef", 6, 1, ""},
		{"TC06: Offset greater than length", "abcdef", 7, 1, ""},
		{"TC07: Empty string", "", 0, 1, ""},
		{"TC08: Offset at -len(s)", "abcdef", -6, 3, "abc"},
		{"TC09: Negative offset with large length", "abcdef", -5, 10, "bcdef"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stringx.Substring(tt.s, tt.offset, tt.length)
			if result != tt.expected {
				t.Errorf("Substring(%q, %d, %d) = %q; want %q", tt.s, tt.offset, tt.length, result, tt.expected)
			}
		})
	}
}
