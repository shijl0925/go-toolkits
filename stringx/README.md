# stringx - Go 字符串工具包

`stringx` 是一个用于处理 Go 语言字符串的工具包，提供了丰富的字符串操作函数。

## 功能列表

### 基本字符串操作
* `Capitalize(s string) string` - 将字符串的第一个字符转为大写，其余保持小写
* `Reverse(s string) string` - 反转字符串
* `Substring(s string, offset int, length int) string` - 提取字符串的子串
* `SwapCase(s string) string` - 将字符串中的英文字母大小写互换

### 字符串大小写检查
* `IsLower(s string) bool` - 判断字符串是否全部为小写
* `IsUpper(s string) bool` - 判断字符串是否全部为大写
* `IsTitle(s string) bool` - 判断字符串是否全部为标题格式

### 字符串分割操作
* `Partition(s, sep string) (string, string, string, error)` - 使用分隔符将字符串分割为三个部分
* `RightPartition(s, sep string) (string, string, string, error)` - 从右侧使用分隔符将字符串分割为三个部分
* `SplitString(str string) []string` - 将字符串分割为单词数组

### 字符串格式化
* `FormatMap(format string, m map[string]any) string` - 类似 Python 的 str.format_map() 行为，使用 map 替换格式字符串中的占位符

### 命名风格转换
* `ToPascal(s string) string` - 将字符串转换为 PascalCase 格式
* `ToCamel(s string) string` - 将字符串转换为 camelCase 格式
* `ToKebab(str string) string` - 将字符串转换为 kebab-case 格式
* `ToSnake(s string) string` - 将字符串转换为 snake_case 格式

## 安装

```shell
go get github.com/shijl0925/go-toolkits/stringx
```
