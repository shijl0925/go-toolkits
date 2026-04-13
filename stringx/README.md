# stringx

`stringx` 提供一组实用的字符串处理函数，覆盖大小写处理、子串截取、分割、模板替换与命名风格转换。它适合做接口参数整理、展示文本处理以及代码生成前的数据规范化。

## 安装

```bash
go get github.com/shijl0925/go-toolkits/stringx
```

## 核心能力

### 基本处理

- `Capitalize`
- `Reverse`
- `Substring`
- `SwapCase`

### 大小写检查

- `IsLower`
- `IsUpper`
- `IsTitle`

### 分割与格式化

- `Partition`
- `RightPartition`
- `SplitString`
- `FormatMap`

### 命名风格转换

- `ToPascal`
- `ToCamel`
- `ToKebab`
- `ToSnake`

## 快速示例

```go
name := stringx.Capitalize("heLLo")      // "Hello"
slug := stringx.ToKebab("UserProfileV2") // "user-profile-v2"
words := stringx.SplitString("HTTPServer2Port")
text := stringx.FormatMap("hello, {name}", map[string]any{"name": "gopher"})
```

## 使用说明

### Substring

`Substring(s, offset, length)` 基于 **rune** 处理，适合中文等 Unicode 文本。

- `offset >= 0`：从左到右取
- `offset < 0`：从尾部反向偏移
- `length <= 0`：返回空字符串
- 越界时不会 panic，而是返回空字符串或截断后的内容

### Partition / RightPartition

- `Partition`：按第一次出现的分隔符切成三段
- `RightPartition`：按最后一次出现的分隔符切成三段
- 分隔符为空时会返回错误

### SplitString

`SplitString` 适合拆分以下风格：

- `camelCase`
- `PascalCase`
- `snake_case`
- `kebab-case`
- 含数字的混合命名，例如 `UserV2Info`

### SwapCase

`SwapCase` 仅对 **ASCII 英文字母** 做大小写互换，其他字符保持不变。

### FormatMap

`FormatMap` 的行为类似 Python `format_map()`：会把形如 `{key}` 的占位符替换成 map 中对应的值；找不到的 key 会保留原样。

## 适用场景

- 接口字段名转换
- 代码生成前的命名统一
- 简单模板文本拼装
- 用户输入规范化

