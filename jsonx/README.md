# jsonx

`jsonx` 是对 Go 标准库 `encoding/json` 的轻量封装，提供更简洁的编解码、Reader/Writer 处理和 JSON 文件读写接口。它适合替代项目里重复出现的 `json.Marshal` / `json.Unmarshal` 样板代码。

## 安装

```bash
go get github.com/shijl0925/go-toolkits/jsonx
```

## 核心能力

### 编码

- `Encode`
- `Dumps`
- `DumpsPretty`
- `EncodeToWriter`

### 解码

- `Decode`
- `Loads`
- `DecodeFromReader`

### 文件操作

- `WriteFile`
- `ReadFile`

## 快速示例

```go
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

payload, _ := jsonx.Dumps(User{Name: "Alice", Age: 18})

var user User
_ = jsonx.Loads(payload, &user)
```

## 使用说明

### Dumps 与 DumpsPretty

- `Dumps`：生成紧凑 JSON 字符串
- `DumpsPretty`：生成带缩进的 JSON 字符串
- 当 `indent <= 0` 时，`DumpsPretty` 会退化为 `Dumps`

### Reader / Writer 场景

- `EncodeToWriter` 适合直接写 HTTP 响应、文件或 buffer
- `DecodeFromReader` 适合读取请求体、文件流或其他输入流

### 文件读写

- `WriteFile` 会以创建/覆盖方式写入目标文件
- `ReadFile` 会从文件中读取并解码到目标对象

## 注意事项

- `jsonx` 仍然基于标准库 `encoding/json`，因此标签、导出字段规则与标准库保持一致。
- `WriteFile` 不会自动创建父目录；调用前请确保目录已经存在。
- 处理大文件或连续 JSON 流时，优先使用 `EncodeToWriter` / `DecodeFromReader`，避免一次性把内容全部放入内存。
- 解码函数要求传入目标对象的指针。

## 适用场景

- 简化 JSON 编解码样板
- 写配置文件、缓存文件、导出文件
- HTTP/消息队列场景下的对象读写

