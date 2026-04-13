# cryptor

`cryptor` 提供两类基础能力：**Base64 编解码**和 **SHA-256 摘要计算**。它适合文件校验、文本摘要、简单编码转换等轻量场景。

## 安装

```bash
go get github.com/shijl0925/go-toolkits/cryptor
```

## 核心函数

### Base64

- `Base64StdEncode`
- `Base64StdDecode`

### SHA-256

- `Sha256Stream`
- `Sha256String`
- `Sha256File`

## 快速示例

```go
encoded := cryptor.Base64StdEncode("hello")
decoded, _ := cryptor.Base64StdDecode(encoded)

sum1 := cryptor.Sha256String("hello")
sum2, err := cryptor.Sha256File("./demo.txt")
_ = decoded
_ = sum1
_ = sum2
_ = err
```

## 使用说明

### Base64 编解码

这里使用的是**标准 Base64 编码表**，适合普通文本、配置和网络传输中的基础编码场景。

### Sha256Stream 与 Sha256String

这两个函数都用于计算字符串的 SHA-256 十六进制摘要，可按个人偏好选择调用名称。

### Sha256File

- 用于计算文件内容的 SHA-256 摘要
- 会校验路径是否合法，并拒绝目录路径
- 当前实现限制文件大小不超过 **1GB**

## 安全提示

- Base64 **不是加密**，只是编码方式，不能用于保护敏感数据。
- SHA-256 是单向摘要，适合做完整性校验，不适合直接替代密码学协议。
- 如果你需要可逆加密、签名或带密钥摘要，请在业务侧引入更专门的安全方案。

