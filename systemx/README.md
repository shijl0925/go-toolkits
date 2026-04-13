# systemx

`systemx` 面向操作系统层面的常用任务，提供环境变量管理、命令执行、进程启动/终止以及操作系统识别能力。它适合脚本化工具、运维辅助程序和本地自动化任务。

## 安装

```bash
go get github.com/shijl0925/go-toolkits/systemx
```

## 核心能力

### 平台识别

- `IsWindows`
- `IsLinux`
- `IsMac`

### 环境变量

- `GetOsEnv`
- `SetOsEnv`
- `RemoveOsEnv`

### 命令与进程

- `ExecCommand`
- `StartProcess`
- `StopProcess`
- `KillProcess`
- `Option`

## 快速示例

```go
stdout, stderr, err := systemx.ExecCommand("echo hello")
_ = stdout
_ = stderr
_ = err

pid, err := systemx.StartProcess("sleep", "5")
if err == nil {
    _ = systemx.KillProcess(pid)
}
```

## 使用说明

### ExecCommand

`ExecCommand` 接收的是**完整命令字符串**，内部会先解析参数，再直接调用目标可执行文件，而不是把命令交给 shell。

这带来两个好处：

- 降低命令注入风险
- 更容易拿到明确的 stdout / stderr 结果

### 相对路径限制

为了安全起见，`ExecCommand` / `StartProcess` 不允许使用带路径分隔符的相对命令路径，例如：

- `./tool`
- `bin/tool`

如果你要执行本地程序，请传入：

- 系统 PATH 中可解析的命令名，或
- 绝对路径

### Option 的用途

`Option` 本质上是 `func(*exec.Cmd)`，适合在执行前自定义 `Dir`、`Env`、标准输入输出等设置。

### StopProcess 与 KillProcess

- `StopProcess`：通过 `Signal(os.Kill)` 终止进程
- `KillProcess`：直接调用 `Process.Kill()`

在实际使用上，两者都属于强终止能力，建议谨慎调用。

## 注意事项

- 命令字符串如果包含未闭合引号或非法转义，会返回解析错误。
- `ExecCommand` 失败时会返回 `stderr` 和 `error`，便于排查问题。
- 不同平台上可执行文件解析、信号处理行为可能存在细微差异，跨平台脚本请补充验证。
- 如果进程需要优雅退出，建议业务侧实现更细粒度的控制协议。

