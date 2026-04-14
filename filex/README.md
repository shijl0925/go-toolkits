# filex

`filex` 是一个覆盖面较广的文件系统工具包，提供路径处理、文件/目录创建、复制、删除、目录遍历与内容读写等能力，适合脚本工具、配置管理、导入导出和本地任务程序。

## 安装

```bash
go get github.com/shijl0925/go-toolkits/filex
```

## 核心能力

### 路径与存在性

- `IsExist`
- `IsSameStat`
- `GetAbsPath`
- `GetBaseName`
- `GetDirName`
- `GetExtension`
- `SplitText`

### 创建与删除

- `CreateFile`
- `CreateDir`
- `RemoveFile`
- `RemoveDir`

### 复制与遍历

- `CopyFile`
- `CopyDir`
- `ListDir`
- `WalkV1`
- `WalkV2`

### 内容与元信息

- `ReadLines`
- `ReadFileToString`
- `WriteStringToFile`
- `IsDir`
- `IsFile`
- `IsLink`
- `FileMode`
- `FileSize`

## 快速示例

```go
_ = filex.CreateDir("./data")
_ = filex.WriteStringToFile("./data/app.txt", "hello")

content, _ := filex.ReadFileToString("./data/app.txt")
_ = content
```

## 使用说明

### CreateFile / CreateDir

- `CreateFile` 会先尝试创建父目录。
- 如果目标文件已经存在，`CreateFile` 会返回错误。
- `CreateDir` 使用递归创建方式，适合一次性建立多层目录。

### CopyFile / CopyDir

- `CopyFile` 会保留源文件权限。
- `CopyDir` 会递归复制整个目录树。
- 对符号链接采取更保守的处理策略：
  - `CopyFile` 会拒绝复制符号链接文件
  - `CopyDir` 遍历时会跳过符号链接项

### WalkV1 与 WalkV2

- `WalkV1`：返回值更直接，适合简单扫描目录树
- `WalkV2`：基于 `filepath.Walk`，适合与标准库遍历行为保持一致的场景

### ReadLines

`ReadLines` 会按行返回文本内容，适合日志、配置、脚本文件等场景。

## 注意事项

- 涉及文件系统写操作时，请显式处理错误，尤其是权限不足、路径不存在、目标已存在等情况。
- `CreateFile`、`CopyFile`、`CopyDir` 对符号链接较为严格，这是为了避免意外覆盖和路径逃逸风险。
- 大文件拷贝与读取会消耗更多 IO 时间；如有流式处理需求，建议在业务侧结合 `io.Reader` / `io.Writer` 扩展。

