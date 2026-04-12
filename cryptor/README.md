# cryptor - Go 加密工具包

`cryptor` 是一个用于处理加密和编码操作的工具包，提供了常用的安全哈希算法和编码功能。

## 功能列表

### Base64 编码/解码
* `Base64StdEncode(s string) string` - 使用标准 Base64 编码字符串
* `Base64StdDecode(s string) (string, error)` - 解码 Base64 编码的字符串

### SHA256 哈希
* `Sha256Stream(s string) string` - 计算字符串的 SHA256 哈希值
* `Sha256String(s string) string` - 计算字符串的 SHA256 哈希值
* `Sha256File(filePath string) (string, error)` - 计算文件的 SHA256 哈希值

## 安装

```shell
go get github.com/shijl0925/go-toolkits/cryptor
```

## 使用示例

```go
package main

import (
    "fmt"
    "github.com/shijl0925/go-toolkits/cryptor"
)

func main() {
    // Base64 编码/解码
    encoded := cryptor.Base64StdEncode("Hello World")
    fmt.Println("Base64 Encoded:", encoded)
    
    decoded, err := cryptor.Base64StdDecode(encoded)
    if err != nil {
        panic(err)
    }
    fmt.Println("Base64 Decoded:", decoded)
    
    // SHA256 哈希
    sha256Hash := cryptor.Sha256Stream("Hello World")
    fmt.Println("SHA256 Hash:", sha256Hash)
}
```
