# structs - Go 结构体工具包

`structs` 是一个用于处理 Go 语言结构体的工具包，提供了结构体与 map 之间转换的功能。

## 功能列表

### 结构体操作
* `StructToMap(s any) map[string]any` - 将结构体转换为 map，支持 JSON 标签和嵌套结构体

## 安装

```shell
go get github.com/shijl0925/go-toolkits/structs
```


## 使用示例

```go
package main

import (
    "fmt"
    "github.com/shijl0925/go-toolkits/structs"
)

type Address struct {
    City  string `json:"city"`
    State string `json:"state"`
}

type Person struct {
    Name    string  `json:"name"`
    Age     int     `json:"age,omitempty"`
    Address Address `json:"address"`
}

func main() {
    person := Person{
        Name: "Alice",
        Age:  0, // 由于 omitempty 标签，此字段会被忽略
        Address: Address{
            City:  "New York",
            State: "NY",
        },
    }
    
    // 将结构体转换为 map
    result := structs.StructToMap(person)
    fmt.Printf("Struct to Map: %+v\n", result)
    // 输出: map[name:Alice address:map[city:New York state:NY]]
}
```


## 特性说明

### JSON 标签支持
- 支持 `json` 标签来指定 map 中的键名
- 支持 `omitempty` 选项，当字段为零值时忽略该字段
- 支持 `flatten` 选项，将嵌套结构体展开到当前层级

### 嵌套结构体处理
- 自动处理嵌套的结构体类型
- 支持匿名结构体的展开（配合 `flatten` 标签）
- 正确处理指针类型的字段

### 类型安全
- 自动解引用指针类型
- 处理 nil 指针情况
- 验证输入是否为有效的结构体类型