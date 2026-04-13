# mathx

`mathx` 聚焦于几个常见但容易重复实现的浮点工具：四舍五入、百分比格式化以及百分比字符串解析。模块很小，但在报表、展示层和配置处理里非常实用。

## 安装

```bash
go get github.com/shijl0925/go-toolkits/mathx
```

## 核心函数

- `RoundToFloat[T float64 | float32](f T, n int) float64`
- `FloatToPercent[T float64 | float32](f T, n uint) string`
- `PercentToFloat(s string) (float64, error)`

## 快速示例

```go
rounded := mathx.RoundToFloat(3.14159, 2)   // 3.14
percent := mathx.FloatToPercent(0.1234, 2)  // "12.34%"
value, err := mathx.PercentToFloat("12.34%")
_ = value
_ = err
```

## 使用说明

### RoundToFloat

- 支持 `float32` 和 `float64`
- 返回值统一为 `float64`
- `n` 表示保留的小数位数

### FloatToPercent

把 `0.1234` 这样的比例值转成带 `%` 的字符串，适合直接展示。

### PercentToFloat

- 输入必须包含 `%` 后缀
- 会自动处理前后空格，以及类似 `"50 %"` 这样的写法
- 解析成功后返回比例值，例如 `"12.5%" -> 0.125`

## 注意事项

- `PercentToFloat` 输入不合法时会返回错误，应显式处理。
- 该模块适合展示与通用业务计算，不建议替代高精度数值库。

