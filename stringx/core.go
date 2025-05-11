package stringx

import "strings"

// Capitalize capitalizes the first letter of a string.
// 将字符串的第一个字符转为大写，其余保持不变
//
// Example
// s := Capitalize("hello")
// fmt.Println(s) // "Hello"
func Capitalize(s string) string {
	if len(s) == 0 {
		panic("string is empty")
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
