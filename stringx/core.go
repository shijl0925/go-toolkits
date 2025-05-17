package stringx

import (
	"strings"
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
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
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

func Substring(s string, offset int, length uint) string {
	r := []rune(s)
	size := len(r)
	if offset < 0 {
		offset += size
	}
	if offset < 0 || offset >= size {
		return ""
	}
	end := offset + int(length)
	if end > size {
		end = size
	}
	return string(r[offset:end])
}

// IsLower returns true if all characters in the string are lower case.
// 判断字符串是否全部为小写。
func IsLower(s string) bool {
	return strings.ToLower(s) == s
}

// IsUpper returns true if all characters in the string are upper case.
// 判断字符串是否全部为大写。
func IsUpper(s string) bool {
	return strings.ToUpper(s) == s
}

// IsTitle returns true if all characters in the string are title case.
// 判断字符串是否全部为标题格式。
func IsTitle(s string) bool {
	return strings.ToTitle(s) == s
}

// Partition splits a string into three parts using a separator.
// 分割字符串为三个部分，使用分隔符。
func Partition(s, sep string) (string, string, string) {
	if sep == "" {
		panic("Partition: sep cannot be empty string")
	}

	parts := strings.SplitN(s, sep, 2)
	if len(parts) == 1 {
		return parts[0], "", ""
	}
	return parts[0], sep, parts[1]
}

// RightPartition mimics Python's str.rpartition() behavior.
// It splits the string s into the part before the last occurrence of sep,
// the sep itself, and the part after.
// If sep is not found, it returns ("", "", s).
func RightPartition(s, sep string) (string, string, string) {
	if sep == "" {
		panic("Partition: sep cannot be empty string")
	}

	idx := strings.LastIndex(s, sep)
	if idx == -1 {
		return "", "", s
	}
	first := s[:idx]
	second := s[idx+len(sep):]
	return first, sep, second
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
			if strVal, ok := value.(string); ok {
				strKey := "{" + key + "}"
				pairs[strKey] = strVal
			}
		}
		pos = end + 1
	}

	for key, value := range pairs {
		format = strings.ReplaceAll(format, key, value)
	}

	return format
}
