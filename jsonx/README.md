# jsonx - Go JSON 工具包

`jsonx` 是一个用于处理 JSON 数据的工具包，提供了便捷的 JSON 编码和解码功能，简化了 Go 语言中 JSON 操作的常用场景。

## 功能列表

### JSON 编码功能
* `Encode(v any) ([]byte, error)` - 将数据编码为 JSON 字节切片
* `Dumps(v any) (string, error)` - 将数据编码为 JSON 字符串
* `DumpsPretty(v any, indent int) (string, error)` - 将数据编码为带格式的 JSON 字符串
* `EncodeToWriter(v any, w io.Writer) error` - 将数据编码为 JSON 并写入到 io.Writer

### JSON 解码功能
* `Decode(data []byte, ptr any) error` - 将 JSON 字节切片解码到数据指针
* `Loads(s string, ptr any) error` - 将 JSON 字符串解码到数据指针
* `DecodeFromReader(r io.Reader, ptr any) error` - 从 io.Reader 中解码 JSON 到数据指针

### 文件操作功能
* `WriteFile(filePath string, data any) error` - 将数据写入 JSON 文件
* `ReadFile(filePath string, v any) error` - 从 JSON 文件读取数据

## 安装

```shell
go get github.com/shijl0925/go-toolkits/jsonx
```
