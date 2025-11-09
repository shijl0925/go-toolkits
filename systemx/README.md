# systemx - Go 系统操作工具包

`systemx` 是一个用于处理系统级操作的工具包，提供了跨平台的系统信息获取、环境变量管理、命令执行和进程控制功能。

## 功能列表

### 操作系统检测
* [IsWindows() bool](https://github.com/shijl0925/go-toolkits/blob/main/systemx/core.go#L15-L17) - 检查当前操作系统是否为 Windows
* [IsLinux() bool](https://github.com/shijl0925/go-toolkits/blob/main/systemx/core.go#L20-L22) - 检查当前操作系统是否为 Linux
* [IsMac() bool](https://github.com/shijl0925/go-toolkits/blob/main/systemx/core.go#L25-L27) - 检查当前操作系统是否为 macOS

### 环境变量操作
* [GetOsEnv(key string) string](https://github.com/shijl0925/go-toolkits/blob/main/systemx/core.go#L30-L32) - 获取指定环境变量的值
* `SetOsEnv(key, value string) error` - 设置环境变量的值
* [RemoveOsEnv(key string) error](https://github.com/shijl0925/go-toolkits/blob/main/systemx/core.go#L40-L42) - 删除指定环境变量

### 命令执行
* `ExecCommand(command string, opts ...Option) (stdout, stderr string, err error)` - 执行系统命令并返回输出结果

### 进程管理
* `StartProcess(command string, args ...string) (int, error)` - 启动一个新进程并返回进程ID
* [StopProcess(pid int) error](https://github.com/shijl0925/go-toolkits/blob/main/systemx/core.go#L92-L99) - 通过进程ID停止进程（发送 kill 信号）
* [KillProcess(pid int) error](https://github.com/shijl0925/go-toolkits/blob/main/systemx/core.go#L102-L109) - 通过进程ID强制终止进程

## 安装

```shell
go get github.com/shijl0925/go-toolkits/systemx
```


## 使用示例

```go
package main

import (
    "fmt"
    "github.com/shijl0925/go-toolkits/systemx"
)

func main() {
    // 检查操作系统类型
    if systemx.IsWindows() {
        fmt.Println("Running on Windows")
    } else if systemx.IsLinux() {
        fmt.Println("Running on Linux")
    } else if systemx.IsMac() {
        fmt.Println("Running on macOS")
    }
    
    // 环境变量操作
    err := systemx.SetOsEnv("MY_VAR", "test_value")
    if err != nil {
        fmt.Printf("Error setting environment variable: %v\n", err)
    }
    
    value := systemx.GetOsEnv("MY_VAR")
    fmt.Printf("MY_VAR = %s\n", value)
    
    // 执行系统命令
    stdout, stderr, err := systemx.ExecCommand("echo Hello World")
    if err != nil {
        fmt.Printf("Error executing command: %v\n", err)
    }
    if stderr != "" {
        fmt.Printf("Command stderr: %s\n", stderr)
    }
    fmt.Printf("Command stdout: %s\n", stdout)
    
    // 进程管理
    pid, err := systemx.StartProcess("sleep", "10")
    if err != nil {
        fmt.Printf("Error starting process: %v\n", err)
    } else {
        fmt.Printf("Started process with PID: %d\n", pid)
        
        // 终止进程
        err = systemx.KillProcess(pid)
        if err != nil {
            fmt.Printf("Error killing process: %v\n", err)
        }
    }
}
```


## 特性说明

### 跨平台支持
- 自动检测操作系统类型并适配相应的命令执行方式
- 在 Linux/macOS 上使用 `/bin/bash -c` 执行命令
- 在 Windows 上使用 `powershell.exe` 执行命令

### 命令执行选项
- 支持通过 `Option` 函数自定义 `exec.Cmd` 的配置
- 返回标准输出、错误输出和错误信息，便于调试

### 进程管理
- 提供启动、停止和强制终止进程的功能
- 通过进程ID进行进程控制，操作精确