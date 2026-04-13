package stringx

import (
	"fmt"
	"regexp"
	"strings"
	"unicode" //nolint:typecheck
)

var (
	// bearer:disable go_lang_permissive_regex_validation
	splitWordReg = regexp.MustCompile(`([a-z])([A-Z0-9])|([a-zA-Z])([0-9])|([0-9])([a-zA-Z])|([A-Z])([A-Z])([a-z])`)
	// bearer:disable go_lang_permissive_regex_validation
	splitNumberLetterReg = regexp.MustCompile(`([0-9])([a-zA-Z])`)
)

// Capitalize converts the first character of a string to upper case and the remaining to lower case.
// 将字符串的第一个字符转为大写，其余保持小写。
//
// Example
// s := Capitalize("hello")
// fmt.Println(s) // "Hello"
func Capitalize(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	return strings.ToUpper(string(r[:1])) + strings.ToLower(string(r[1:]))
}

// Reverse reverses a string.
// 反转字符串
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func Substring(s string, offset int, length int) string {
	if length <= 0 {
		return ""
	}
	r := []rune(s)
	size := len(r)
	if offset < 0 {
		offset += size
	}
	if offset < 0 || offset >= size {
		return ""
	}
	end := offset + length
	if end > size {
		end = size
	}
	return string(r[offset:end])
}

// IsLower returns true if all characters in the string are lower case.
// 判断字符串是否全部为小写。
func IsLower(s string) bool {
	return hasLetter(s) && strings.ToLower(s) == s
}

// IsUpper returns true if all characters in the string are upper case.
// 判断字符串是否全部为大写。
func IsUpper(s string) bool {
	return hasLetter(s) && strings.ToUpper(s) == s
}

// IsTitle returns true if all characters in the string are title case.
// 判断字符串是否全部为标题格式。
func IsTitle(s string) bool {
	return hasLetter(s) && strings.ToTitle(s) == s
}

// hasLetter reports whether s contains at least one Unicode letter.
func hasLetter(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

// Partition splits a string into three parts using a separator.
// 分割字符串为三个部分，使用分隔符。
func Partition(s, sep string) (string, string, string, error) {
	if sep == "" {
		return "", "", s, fmt.Errorf("partition: sep cannot be empty string")
	}

	parts := strings.SplitN(s, sep, 2)
	if len(parts) == 1 {
		return parts[0], "", "", nil
	}
	return parts[0], sep, parts[1], nil
}

// RightPartition mimics Python's str.rpartition() behavior.
// It splits the string s into the part before the last occurrence of sep,
// the sep itself, and the part after.
// If sep is not found, it returns ("", "", s).
func RightPartition(s, sep string) (string, string, string, error) {
	if sep == "" {
		return "", "", s, fmt.Errorf("RightPartition: sep cannot be empty string")
	}

	idx := strings.LastIndex(s, sep)
	if idx == -1 {
		return "", "", s, nil
	}
	first := s[:idx]
	second := s[idx+len(sep):]
	return first, sep, second, nil
}

// SwapCase 将输入字符串中的 ASCII 英文字母大小写互换。
// 注意：本函数仅处理 a-z 和 A-Z 字符，其他字符保持不变。
func SwapCase(s string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case 'a' <= r && r <= 'z':
			return r - 'a' + 'A'

		case 'A' <= r && r <= 'Z':
			return r - 'A' + 'a'
		default:
			return r
		}
	}, s)
}

// FormatMap mimics Python's str.format_map() behavior.
// It replaces each occurrence of '{key}' in the format string with the value of the corresponding key in the map.
//
// Example:
//
//	FormatMap("Hello, {name}!", map[string]any{"name": "John"}) // returns "Hello, John!"
func FormatMap(format string, m map[string]any) string {
	pairs := make(map[string]string)
	pos := 0
	for {
		idx := strings.Index(format[pos:], "{")
		if idx == -1 {
			break
		}

		idx += pos
		end := idx + 1

		for end < len(format) && format[end] != '}' {
			end++
		}

		if end >= len(format) {
			break
		}
		key := format[idx+1 : end]
		//if key == "" {
		//	pos = end + 1
		//	continue
		//}

		if value, ok := m[key]; ok {
			strKey := "{" + key + "}"
			pairs[strKey] = fmt.Sprint(value)
		}
		pos = end + 1
	}

	for key, value := range pairs {
		format = strings.ReplaceAll(format, key, value)
	}

	return format
}

// SplitString splits string into an array of its words.
func SplitString(str string) []string {
	// Step 1: 处理驼峰命名和其他混合大小写/数字结构
	str = splitWordReg.ReplaceAllString(str, `$1$3$5$7 $2$4$6$8$9`)

	// Step 2: 分离数字和字母组合
	// example: Int8Value => Int 8Value => Int 8 Value
	str = splitNumberLetterReg.ReplaceAllString(str, "$1 $2")

	var result strings.Builder
	for _, r := range str {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
		} else {
			result.WriteRune(' ')
		}
	}

	return strings.Fields(result.String())
}

func ToPascal(s string) string {
	words := SplitString(s)
	for i, item := range words {
		item = Capitalize(item)
		words[i] = item
	}
	return strings.Join(words, "")
}

// ToCamel converts a string to CamelCase
func ToCamel(s string) string {
	words := SplitString(s)
	for i, item := range words {
		if i == 0 {
			item = strings.ToLower(item)
		} else {
			item = Capitalize(item)
		}

		words[i] = item
	}
	return strings.Join(words, "")
}

// ToKebab converts string to kebab-case.
func ToKebab(str string) string {
	items := SplitString(str)
	for i := range items {
		items[i] = strings.ToLower(items[i])
	}
	return strings.Join(items, "-")
}

// ToSnake converts a string to snake_case
func ToSnake(s string) string {
	words := SplitString(s)
	for i, item := range words {
		words[i] = strings.ToLower(item)
	}
	return strings.Join(words, "_")
}
