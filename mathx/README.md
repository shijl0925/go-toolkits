# mathx - Go 数学工具包

`mathx` 是一个用于处理数学运算的工具包，提供了浮点数精度处理、百分比转换等常用数学功能。

## 功能列表

### 浮点数处理
* `RoundToFloat[T float64 | float32](f T, n int) float64` - 将浮点数四舍五入到指定精度
* `FloatToPercent[T float64 | float32](f T, n uint) string` - 将浮点数转换为指定精度的百分比字符串
* `PercentToFloat(s string) (float64, error)` - 将百分比字符串转换为浮点数

## 安装

```shell
go get github.com/shijl0925/go-toolkits/mathx
```


## 使用示例

```go
package main

import (
    "fmt"
    "github.com/shijl0925/go-toolkits/mathx"
)

func main() {
    // 浮点数四舍五入
    result := mathx.RoundToFloat(3.14159, 2)
    fmt.Println("Rounded:", result) // 3.14
    
    // 浮点数转百分比
    percent := mathx.FloatToPercent(0.12345, 2)
    fmt.Println("Percent:", percent) // 12.35%
    
    // 百分比转浮点数
    f, err := mathx.PercentToFloat("12.35%")
    if err != nil {
        panic(err)
    }
    fmt.Println("Float:", f) // 0.1235
}
```
