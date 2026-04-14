# structs

`structs` 专注于一件事：把 Go 结构体安全地转换成 `map[string]any`。它支持 JSON 标签、`omitempty`、嵌套结构体和指针字段，适合作为日志、动态序列化、模板上下文或通用数据拼装的辅助工具。

## 安装

```bash
go get github.com/shijl0925/go-toolkits/structs
```

## 核心函数

- `StructToMap(s any) map[string]any`

## 快速示例

```go
type Address struct {
    City string `json:"city"`
}

type User struct {
    Name    string   `json:"name"`
    Age     int      `json:"age,omitempty"`
    Address *Address `json:"address"`
}

u := User{
    Name:    "Alice",
    Address: &Address{City: "Shanghai"},
}

result := structs.StructToMap(u)
// map[string]any{"name":"Alice", "address": Address{City: "Shanghai"}}
```

## 标签规则

### 基本 JSON 标签

如果字段存在 `json:"name"`，生成 map 时会优先使用标签名作为 key。

### omitempty

当字段带有 `omitempty` 且值为零值时，该字段不会写入结果。

### flatten

当字段是结构体并带有 `flatten` 选项时，会继续递归转换该字段。

- **匿名结构体字段 + flatten**：会把子字段合并到当前层级
- **普通结构体字段 + flatten**：会把该字段转换成嵌套 `map[string]any`

## 行为说明

- 支持传入结构体或结构体指针。
- 如果传入的是 `nil` 指针、非结构体类型，结果会是空 map。
- 指针字段在非空时会自动解引用；为 `nil` 时会被跳过。
- 嵌套结构体会按标签规则继续处理。

## 适用场景

- 生成日志上下文
- 构造模板渲染数据
- 把领域对象转成通用键值结构
- 与 `jsonx`、`mapx` 组合做二次转换

## 注意事项

- 当前模块只提供 **struct -> map** 能力，不包含反向映射。
- 反射实现更适合配置、日志、辅助转换等场景；高频热点路径建议评估性能。

