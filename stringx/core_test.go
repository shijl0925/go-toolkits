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
			first, second, third, err := stringx.Partition(tt.s, tt.sep)
			if err != nil {
				t.Errorf("Partition(%q, %q) failed: %v", tt.s, tt.sep, err)
			}
			if !reflect.DeepEqual([3]string{first, second, third}, tt.expected) {
				t.Errorf("Partition(%q, %q) = %v, want %v", tt.s, tt.sep, [3]string{first, second, third}, tt.expected)
			}
		})
	}
}

/*func TestPartition_InvalidSep(t *testing.T) {
	t.Run("Invalid sep", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic due to invalid sep, but did not get one")
			}
		}()

		stringx.Partition("abc", "")
	})
}*/

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
			first, second, third, err := stringx.RightPartition(tt.s, tt.sep)
			if err != nil {
				t.Errorf("RightPartition(%q, %q) failed: %v", tt.s, tt.sep, err)
			}
			if !reflect.DeepEqual([3]string{first, second, third}, tt.expected) {
				t.Errorf("RightPartition(%q, %q) = %v, want %v", tt.s, tt.sep, [3]string{first, second, third}, tt.expected)
			}
		})
	}
}

/*func TestRightPartition_InvalidSep(t *testing.T) {
	t.Run("Invalid sep", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic due to invalid sep, but did not get one")
			}
		}()

		stringx.RightPartition("abc", "")
	})
}*/

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
		{
			name:     "Non-string values are formatted",
			format:   "count={count}, active={active}",
			inputMap: map[string]any{"count": 5, "active": true},
			expected: "count=5, active=true",
		},
		{
			name:     "Nil value is formatted",
			format:   "value={value}",
			inputMap: map[string]any{"value": nil},
			expected: "value=<nil>",
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

func TestSplitString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"TC01: Camel case", "camelCase", []string{"camel", "Case"}},
		{"TC02: Pascal case", "PascalCase", []string{"Pascal", "Case"}},
		{"TC03: Number and letters", "Int8Value", []string{"Int", "8", "Value"}},
		{"TC04: Special characters", "hello-world!", []string{"hello", "world"}},
		{"TC05: Acronym in word", "HTTPRequest", []string{"HTTP", "Request"}},
		{"TC06: Empty input", "", []string{}},
		{"TC07: Mixed lower and upper", "aBcDefGhi", []string{"a", "Bc", "Def", "Ghi"}},
		{"TC08: Letters with numbers", "ABC123XYZ", []string{"ABC", "123", "XYZ"}},
		{"TC09: Underscore and number", "user_id_123", []string{"user", "id", "123"}},
		{"TC10: Acronyms mixed", "MyURLIsHTTPS", []string{"My", "URL", "IsHTTPS"}},
		{"TC11", "AnyKind of_string", []string{"Any", "Kind", "of", "string"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stringx.SplitString(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("SplitString(%q) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestAllCase(t *testing.T) {
	type output struct {
		PascalCase string
		CamelCase  string
		KebabCase  string
		SnakeCase  string
	}
	name := ""
	tests := []struct {
		name   string
		input  string
		output output
	}{
		{name: name, output: output{}},
		{name: name, input: ".", output: output{}},
		{name: name, input: "Hello world!", output: output{
			PascalCase: "HelloWorld",
			CamelCase:  "helloWorld",
			KebabCase:  "hello-world",
			SnakeCase:  "hello_world",
		}},
		{name: name, input: "A", output: output{
			PascalCase: "A",
			CamelCase:  "a",
			KebabCase:  "a",
			SnakeCase:  "a",
		}},
		{name: name, input: "a", output: output{
			PascalCase: "A",
			CamelCase:  "a",
			KebabCase:  "a",
			SnakeCase:  "a",
		}},
		{name: name, input: "foo", output: output{
			PascalCase: "Foo",
			CamelCase:  "foo",
			KebabCase:  "foo",
			SnakeCase:  "foo",
		}},
		{name: name, input: "snake_case", output: output{
			PascalCase: "SnakeCase",
			CamelCase:  "snakeCase",
			KebabCase:  "snake-case",
			SnakeCase:  "snake_case",
		}},
		{name: name, input: "SNAKE_CASE", output: output{
			PascalCase: "SnakeCase",
			CamelCase:  "snakeCase",
			KebabCase:  "snake-case",
			SnakeCase:  "snake_case",
		}},
		{name: name, input: "kebab-case", output: output{
			PascalCase: "KebabCase",
			CamelCase:  "kebabCase",
			KebabCase:  "kebab-case",
			SnakeCase:  "kebab_case",
		}},
		{name: name, input: "PascalCase", output: output{
			PascalCase: "PascalCase",
			CamelCase:  "pascalCase",
			KebabCase:  "pascal-case",
			SnakeCase:  "pascal_case",
		}},
		{name: name, input: "camelCase", output: output{
			PascalCase: "CamelCase",
			CamelCase:  "camelCase",
			KebabCase:  `camel-case`,
			SnakeCase:  "camel_case",
		}},
		{name: name, input: "Title Case", output: output{
			PascalCase: "TitleCase",
			CamelCase:  "titleCase",
			KebabCase:  "title-case",
			SnakeCase:  "title_case",
		}},
		{name: name, input: "point.case", output: output{
			PascalCase: "PointCase",
			CamelCase:  "pointCase",
			KebabCase:  "point-case",
			SnakeCase:  "point_case",
		}},
		{name: name, input: "snake_case_with_more_words", output: output{
			PascalCase: "SnakeCaseWithMoreWords",
			CamelCase:  "snakeCaseWithMoreWords",
			KebabCase:  "snake-case-with-more-words",
			SnakeCase:  "snake_case_with_more_words",
		}},
		{name: name, input: "SNAKE_CASE_WITH_MORE_WORDS", output: output{
			PascalCase: "SnakeCaseWithMoreWords",
			CamelCase:  "snakeCaseWithMoreWords",
			KebabCase:  "snake-case-with-more-words",
			SnakeCase:  "snake_case_with_more_words",
		}},
		{name: name, input: "kebab-case-with-more-words", output: output{
			PascalCase: "KebabCaseWithMoreWords",
			CamelCase:  "kebabCaseWithMoreWords",
			KebabCase:  "kebab-case-with-more-words",
			SnakeCase:  "kebab_case_with_more_words",
		}},
		{name: name, input: "PascalCaseWithMoreWords", output: output{
			PascalCase: "PascalCaseWithMoreWords",
			CamelCase:  "pascalCaseWithMoreWords",
			KebabCase:  "pascal-case-with-more-words",
			SnakeCase:  "pascal_case_with_more_words",
		}},
		{name: name, input: "camelCaseWithMoreWords", output: output{
			PascalCase: "CamelCaseWithMoreWords",
			CamelCase:  "camelCaseWithMoreWords",
			KebabCase:  "camel-case-with-more-words",
			SnakeCase:  "camel_case_with_more_words",
		}},
		{name: name, input: "Title Case With More Words", output: output{
			PascalCase: "TitleCaseWithMoreWords",
			CamelCase:  "titleCaseWithMoreWords",
			KebabCase:  "title-case-with-more-words",
			SnakeCase:  "title_case_with_more_words",
		}},
		{name: name, input: "point.case.with.more.words", output: output{
			PascalCase: "PointCaseWithMoreWords",
			CamelCase:  "pointCaseWithMoreWords",
			KebabCase:  "point-case-with-more-words",
			SnakeCase:  "point_case_with_more_words",
		}},
		{name: name, input: "snake_case__with___multiple____delimiters", output: output{
			PascalCase: "SnakeCaseWithMultipleDelimiters",
			CamelCase:  "snakeCaseWithMultipleDelimiters",
			KebabCase:  "snake-case-with-multiple-delimiters",
			SnakeCase:  "snake_case_with_multiple_delimiters",
		}},
		{name: name, input: "SNAKE_CASE__WITH___multiple____DELIMITERS", output: output{
			PascalCase: "SnakeCaseWithMultipleDelimiters",
			CamelCase:  "snakeCaseWithMultipleDelimiters",
			KebabCase:  "snake-case-with-multiple-delimiters",
			SnakeCase:  "snake_case_with_multiple_delimiters",
		}},
		{name: name, input: "kebab-case--with---multiple----delimiters", output: output{
			PascalCase: "KebabCaseWithMultipleDelimiters",
			CamelCase:  "kebabCaseWithMultipleDelimiters",
			KebabCase:  "kebab-case-with-multiple-delimiters",
			SnakeCase:  "kebab_case_with_multiple_delimiters",
		}},
		{name: name, input: "Title Case  With   Multiple    Delimiters", output: output{
			PascalCase: "TitleCaseWithMultipleDelimiters",
			CamelCase:  "titleCaseWithMultipleDelimiters",
			KebabCase:  "title-case-with-multiple-delimiters",
			SnakeCase:  "title_case_with_multiple_delimiters",
		}},
		{name: name, input: "point.case..with...multiple....delimiters", output: output{
			PascalCase: "PointCaseWithMultipleDelimiters",
			CamelCase:  "pointCaseWithMultipleDelimiters",
			KebabCase:  "point-case-with-multiple-delimiters",
			SnakeCase:  "point_case_with_multiple_delimiters",
		}},
		{name: name, input: " leading space", output: output{
			PascalCase: "LeadingSpace",
			CamelCase:  "leadingSpace",
			KebabCase:  "leading-space",
			SnakeCase:  "leading_space",
		}},
		{name: name, input: "   leading spaces", output: output{
			PascalCase: "LeadingSpaces",
			CamelCase:  "leadingSpaces",
			KebabCase:  "leading-spaces",
			SnakeCase:  "leading_spaces",
		}},
		{name: name, input: "\t\t\r\n leading whitespaces", output: output{
			PascalCase: "LeadingWhitespaces",
			CamelCase:  "leadingWhitespaces",
			KebabCase:  "leading-whitespaces",
			SnakeCase:  "leading_whitespaces",
		}},
		{name: name, input: "trailing space ", output: output{
			PascalCase: "TrailingSpace",
			CamelCase:  "trailingSpace",
			KebabCase:  "trailing-space",
			SnakeCase:  "trailing_space",
		}},
		{name: name, input: "trailing spaces   ", output: output{
			PascalCase: "TrailingSpaces",
			CamelCase:  "trailingSpaces",
			KebabCase:  "trailing-spaces",
			SnakeCase:  "trailing_spaces",
		}},
		{name: name, input: "trailing whitespaces\t\t\r\n", output: output{
			PascalCase: "TrailingWhitespaces",
			CamelCase:  "trailingWhitespaces",
			KebabCase:  "trailing-whitespaces",
			SnakeCase:  "trailing_whitespaces",
		}},
		{name: name, input: " on both sides ", output: output{
			PascalCase: "OnBothSides",
			CamelCase:  "onBothSides",
			KebabCase:  "on-both-sides",
			SnakeCase:  "on_both_sides",
		}},
		{name: name, input: "    many on both sides  ", output: output{
			PascalCase: "ManyOnBothSides",
			CamelCase:  "manyOnBothSides",
			KebabCase:  "many-on-both-sides",
			SnakeCase:  "many_on_both_sides",
		}},
		{name: name, input: "\r whitespaces on both sides\t\t\r\n", output: output{
			PascalCase: "WhitespacesOnBothSides",
			CamelCase:  "whitespacesOnBothSides",
			KebabCase:  "whitespaces-on-both-sides",
			SnakeCase:  "whitespaces_on_both_sides",
		}},
		{name: name, input: "  extraSpaces in_This TestCase Of MIXED_CASES\t", output: output{
			PascalCase: "ExtraSpacesInThisTestCaseOfMixedCases",
			CamelCase:  "extraSpacesInThisTestCaseOfMixedCases",
			KebabCase:  "extra-spaces-in-this-test-case-of-mixed-cases",
			SnakeCase:  "extra_spaces_in_this_test_case_of_mixed_cases",
		}},
		{name: name, input: "CASEBreak", output: output{
			PascalCase: "CaseBreak",
			CamelCase:  "caseBreak",
			KebabCase:  "case-break",
			SnakeCase:  "case_break",
		}},
		{name: name, input: "ID", output: output{
			PascalCase: "Id",
			CamelCase:  "id",
			KebabCase:  "id",
			SnakeCase:  "id",
		}},
		{name: name, input: "userID", output: output{
			PascalCase: "UserId",
			CamelCase:  "userId",
			KebabCase:  "user-id",
			SnakeCase:  "user_id",
		}},
		{name: name, input: "JSON_blob", output: output{
			PascalCase: "JsonBlob",
			CamelCase:  "jsonBlob",
			KebabCase:  "json-blob",
			SnakeCase:  "json_blob",
		}},
		{name: name, input: "HTTPStatusCode", output: output{
			PascalCase: "HttpStatusCode",
			CamelCase:  "httpStatusCode",
			KebabCase:  "http-status-code",
			SnakeCase:  "http_status_code",
		}},
		{name: name, input: "FreeBSD and SSLError are not golang initialisms", output: output{
			PascalCase: "FreeBsdAndSslErrorAreNotGolangInitialisms",
			CamelCase:  "freeBsdAndSslErrorAreNotGolangInitialisms",
			KebabCase:  "free-bsd-and-ssl-error-are-not-golang-initialisms",
			SnakeCase:  "free_bsd_and_ssl_error_are_not_golang_initialisms",
		}},
		{name: name, input: "David's Computer", output: output{
			PascalCase: "DavidSComputer",
			CamelCase:  "davidSComputer",
			KebabCase:  "david-s-computer",
			SnakeCase:  "david_s_computer",
		}},
		{name: name, input: "http200", output: output{
			PascalCase: "Http200",
			CamelCase:  "http200",
			KebabCase:  "http-200",
			SnakeCase:  "http_200",
		}},
		{name: name, input: "NumberSplittingVersion1.0r3", output: output{
			PascalCase: "NumberSplittingVersion10R3",
			CamelCase:  "numberSplittingVersion10R3",
			KebabCase:  "number-splitting-version-1-0-r3",
			SnakeCase:  "number_splitting_version_1_0_r3",
		}},
		{name: name, input: "When you have a comma, odd results", output: output{
			PascalCase: "WhenYouHaveACommaOddResults",
			CamelCase:  "whenYouHaveACommaOddResults",
			KebabCase:  "when-you-have-a-comma-odd-results",
			SnakeCase:  "when_you_have_a_comma_odd_results",
		}},
		{name: name, input: "Ordinal numbers work: 1st 2nd and 3rd place", output: output{
			PascalCase: "OrdinalNumbersWork1St2NdAnd3RdPlace",
			CamelCase:  "ordinalNumbersWork1St2NdAnd3RdPlace",
			KebabCase:  "ordinal-numbers-work-1-st-2-nd-and-3-rd-place",
			SnakeCase:  "ordinal_numbers_work_1_st_2_nd_and_3_rd_place",
		}},
		{name: name, input: "BadUTF8\xe2\xe2\xa1", output: output{
			PascalCase: "BadUtf8",
			CamelCase:  "badUtf8",
			KebabCase:  "bad-utf-8",
			SnakeCase:  "bad_utf_8",
		}},
		{name: name, input: "IDENT3", output: output{
			PascalCase: "Ident3",
			CamelCase:  "ident3",
			KebabCase:  "ident-3",
			SnakeCase:  "ident_3",
		}},
		{name: name, input: "LogRouterS3BucketName", output: output{
			PascalCase: "LogRouterS3BucketName",
			CamelCase:  "logRouterS3BucketName",
			KebabCase:  "log-router-s3-bucket-name",
			SnakeCase:  "log_router_s3_bucket_name",
		}},
		{name: name, input: "PINEAPPLE", output: output{
			PascalCase: "Pineapple",
			CamelCase:  "pineapple",
			KebabCase:  "pineapple",
			SnakeCase:  "pineapple",
		}},
		{name: name, input: "Int8Value", output: output{
			PascalCase: "Int8Value",
			CamelCase:  "int8Value",
			KebabCase:  "int-8-value",
			SnakeCase:  "int_8_value",
		}},
		{name: name, input: "first.last", output: output{
			PascalCase: "FirstLast",
			CamelCase:  "firstLast",
			KebabCase:  "first-last",
			SnakeCase:  "first_last",
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pascal := stringx.ToPascal(test.input)
			if pascal != test.output.PascalCase {
				t.Errorf("ToPascal(%q) = %q; expected %q", test.input, pascal, test.output.PascalCase)
			}
			camel := stringx.ToCamel(test.input)
			if camel != test.output.CamelCase {
				t.Errorf("ToCamel(%q) = %q; expected %q", test.input, camel, test.output.CamelCase)
			}
			kebab := stringx.ToKebab(test.input)
			if kebab != test.output.KebabCase {
				t.Errorf("ToKebab(%q) = %q; expected %q", test.input, kebab, test.output.KebabCase)
			}
			snake := stringx.ToSnake(test.input)
			if snake != test.output.SnakeCase {
				t.Errorf("ToSnake(%q) = %q; expected %q", test.input, snake, test.output.SnakeCase)
			}
		})
	}
}
