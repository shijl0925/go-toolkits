# filex - Go 文件操作工具包

`filex` 是一个用于处理文件和目录操作的工具包，提供了丰富的文件系统操作函数，包括文件创建、复制、删除、路径处理等功能。

## 功能列表

### 文件/目录存在性检查
* `IsExist(path string` - 检查文件或目录是否存在
* `IsSameStat(path1, path2 string) (bool, error)` - 检查两个路径是否指向同一个文件或目录

### 路径操作
* `GetAbsPath(path string) (string, error)` - 获取路径的绝对路径
* `GetBaseName(path string) string` - 获取路径的文件名部分
* `GetDirName(path string) string` - 获取路径的目录部分
* `GetExtension(path string) string` - 获取文件扩展名
* `SplitText(path string) (string, string)` - 分割路径为文件名（不含扩展名）和扩展名

### 文件/目录创建
* `CreateFile(path string) error` - 创建文件
* `CreateDir(absPath string) error` - 创建目录

### 文件/目录复制
* `CopyFile(srcPath string, dstPath string) error` - 复制文件
* `CopyDir(srcPath string, dstPath string) error` - 复制目录

### 文件/目录类型检查
* `IsDir(path string) (bool, error)` - 检查路径是否为目录
* `IsFile(path string) (bool, error)` - 检查路径是否为文件
* `IsLink(path string) (bool, error)` - 检查路径是否为符号链接

### 目录内容操作
* `ListDir(path string) ([]string, error)` - 列出目录中的所有文件和子目录
* `WalkV1(root string) []WalkResult` - 遍历目录（深度优先）
* `WalkV2(root string) ([]*WalkResult, error)` - 遍历目录（使用 filepath.Walk）

### 文件/目录删除
* `RemoveFile(path string) error` - 删除文件
* `RemoveDir(path string) error` - 删除目录

### 文件内容读写
* `ReadLines(path string) ([]string, error)` - 逐行读取文件内容
* `ReadFileToString(path string) (string, error)` - 读取文件内容为字符串
* `WriteStringToFile(filepath string, content string) error` - 将字符串写入文件

### 文件信息获取
* `FileMode(path string) (fs.FileMode, error)` - 获取文件模式和权限
* `FileSize(path string) (int64, error)` - 获取文件大小（字节）

## 安装

```shell
go get github.com/shijl0925/go-toolkits/filex
```


## 使用示例

```go
package main

import (
    "fmt"
    "github.com/shijl0925/go-toolkits/filex"
)

func main() {
    // 检查文件是否存在
    exists := filex.IsExist("test.txt")
    fmt.Println("File exists:", exists)
    
    // 获取文件信息
    if exists {
        size, err := filex.FileSize("test.txt")
        if err != nil {
            fmt.Println("Error getting file size:", err)
        } else {
            fmt.Println("File size:", size, "bytes")
        }
    }
    
    // 创建目录
    err := filex.CreateDir("./test_dir")
    if err != nil {
        fmt.Println("Error creating directory:", err)
    }
    
    // 写入文件
    err = filex.WriteStringToFile("./test_dir/test.txt", "Hello, World!")
    if err != nil {
        fmt.Println("Error writing file:", err)
    }
    
    // 读取文件
    content, err := filex.ReadFileToString("./test_dir/test.txt")
    if err != nil {
        fmt.Println("Error reading file:", err)
    } else {
        fmt.Println("File content:", content)
    }
    
    // 列出目录内容
    files, err := filex.ListDir("./test_dir")
    if err != nil {
        fmt.Println("Error listing directory:", err)
    } else {
        fmt.Println("Directory contents:", files)
    }
}
```
