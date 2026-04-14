# gormx

`gormx` 是对 GORM 的轻量增强，重点解决三个问题：**类型安全的条件构造**、**可复用的仓储层**、**统一的服务层封装**。如果你的项目已经使用 GORM，希望减少字符串字段名、重复 CRUD 和样板 service 代码，这个模块会更顺手。

## 安装

```bash
go get github.com/shijl0925/go-toolkits/gormx
```

## 核心能力

- **Query[T] 查询构建器**：支持用字段指针构造条件，减少手写列名。
- **BaseRepo[T] 仓储基类**：内置单表常见 CRUD、分页、批量写入、聚合与增减操作。
- **ServiceImplement[Entity, T] 服务层实现**：把 repo 能力提升为统一的 service 接口。
- **AssociationManager[T, R] 关联管理器**：统一处理关联的新增、替换、移除和查询。
- **SQL 预览能力**：可在不执行 SQL 的情况下输出生成结果，便于调试。

## 使用前提

在调用 `GetDb`、`BaseRepo`、`Query[T]` 等能力前，需要先初始化全局 `*gorm.DB`。

### MySQL

```go
package main

import (
    "github.com/shijl0925/go-toolkits/gormx"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func main() {
    db, err := gorm.Open(mysql.Open("user:pass@tcp(127.0.0.1:3306)/demo?parseTime=True"), &gorm.Config{})
    if err != nil {
        panic(err)
    }
    gormx.Init(db)
}
```

### PostgreSQL

```bash
go get gorm.io/driver/postgres
```

```go
package main

import (
    "github.com/shijl0925/go-toolkits/gormx"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func main() {
    dsn := "host=127.0.0.1 user=postgres password=postgres dbname=demo port=5432 sslmode=disable TimeZone=Asia/Shanghai"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic(err)
    }
    gormx.Init(db)
}
```

DSN 字段说明：

| 字段 | 说明 | 常用值示例 |
| --- | --- | --- |
| `host` | 数据库地址 | `127.0.0.1` |
| `port` | 端口 | `5432` |
| `user` | 用户名 | `postgres` |
| `password` | 密码 | `postgres` |
| `dbname` | 数据库名 | `demo` |
| `sslmode` | SSL 模式 | `disable`、`require` |
| `TimeZone` | 时区 | `Asia/Shanghai`、`UTC` |

## 快速开始

### 1. 使用 Query[T] 构造查询

```go
type User struct {
    ID     uint   `gorm:"primaryKey"`
    Name   string `gorm:"column:name"`
    Age    int    `gorm:"column:age"`
    Status int    `gorm:"column:status"`
}

func listAdults() ([]User, error) {
    query, user := gormx.NewQuery[User]()
    query.
        Eq(&user.Status, 1).
        Ge(&user.Age, 18).
        OrderDesc(&user.ID).
        Limit(20)

    var users []User
    if err := query.Find(&users); err != nil {
        return nil, err
    }
    return users, nil
}
```

### 2. 使用 BaseRepo[T] 复用通用数据访问

```go
type UserRepo struct {
    gormx.BaseRepo[User]
}

func (r *UserRepo) FindEnabled() ([]User, error) {
    return r.SelectListByOpts(gormx.Where("status = ?", 1))
}
```

### 3. 使用 ServiceImplement 统一服务层

```go
type IUserRepo interface {
    gormx.IBaseRepo[User]
}

type UserService struct {
    *gormx.ServiceImplement[IUserRepo, User]
}

func NewUserService(repo IUserRepo) *UserService {
    return &UserService{
        ServiceImplement: gormx.NewServiceImplement[IUserRepo, User](repo),
    }
}
```

### 4. 预览生成的 SQL

```go
query, user := gormx.NewQuery[User]()
query.Eq(&user.Name, "alice").In(&user.Status, []int{1, 2})
sql, args := query.ToSQLAndArgs()
```

## API 导览

### Query[T] 常用条件能力

- 比较：`Eq`、`Ne`、`Gt`、`Ge`、`Lt`、`Le`
- 集合：`In`、`NotIn`、`Between`、`NotBetween`
- 模糊：`Like`、`Regexp`
- 空值：`IsNull`、`IsNotNull`
- 排序与分页：`OrderAsc`、`OrderDesc`、`Limit`、`Offset`
- 聚合与分组：`GroupBy`、`Having`、`Count`、`Sum`、`Avg`、`Max`、`Min`
- 查询执行：`First`、`Find`、`Scan`、`Pluck`、`RawRows`
- 关联与扩展：`Preload`、`Join`、`Or`、`Not`
- 子查询：`SubQueryEq`、`SubQueryIn`、`InSql`、`NotInSql`、`GtSql`、`GeSql`、`LtSql`、`LeSql`

### BaseRepo[T] 常用仓储能力

- 单条查询：`SelectOneById`、`SelectOneByOpts`、`SelectOneByMap`
- 列表与分页：`SelectListByIds`、`SelectListByOpts`、`SelectListByMap`、`SelectPage`、`SelectCount`
- 写入：`Insert`、`InsertBatch`、`InsertInBatches`、`InsertOrUpdate`
- 更新：`Update`、`UpdateById`、`UpdateByOpts`、`Upsert`
- 删除：`Delete`、`DeleteById`、`DeleteBatchIds`、`DeleteByOpts`、`DeleteByMap`
- 辅助：`Exists`、`GetOrCreate`、`UpdateOrCreate`、`Increment`、`Decrement`

### AssociationManager[T, R]

- `Add`：追加关联
- `Remove`：移除关联
- `Set`：整体替换关联
- `Clear`：清空关联
- `Count`：统计关联数量
- `All`：读取全部关联数据

## 使用建议

- **优先传字段指针**：例如 `Eq(&user.Name, "alice")`，可避免手写列名带来的拼写问题。
- **跨表场景可传字符串字段**：如 `users.name`，适用于 join 查询。
- **先处理初始化**：未调用 `Init(db)` 时，全局数据库实例不可用。
- **SelectOne 系列会校验唯一性**：查询结果多于一条时会返回错误，而不是静默取第一条。
- **ServiceImplement 依赖 repo 接口能力**：若仓储未实现对应方法，会返回 `repo does not implement ...` 错误。

## 适用场景

- GORM 项目中统一单表 CRUD 入口
- 希望减少字符串列名和样板查询代码
- 按 repo / service 分层组织业务代码
- 需要统一处理模型关联关系

## PostgreSQL 使用场景

`gormx` 对 PostgreSQL 提供原生支持。大多数 API 与 MySQL 完全一致，切换时只需替换驱动和 DSN，少量行为存在差异，具体如下。

### 占位符格式

PostgreSQL 使用 `$1`、`$2` 作为参数占位符，而 MySQL 使用 `?`。gormx 的所有查询构建方法（`Eq`、`In`、`Between` 等）已经由底层 GORM 驱动自动转换，**调用方无需做任何改动**。

以下两段代码在 MySQL 和 PostgreSQL 下的写法完全一样：

```go
query, user := gormx.NewQuery[User]()
query.Eq(&user.Status, 1).Ge(&user.Age, 18)

// MySQL 实际执行: WHERE `status` = ? AND `age` >= ?
// PostgreSQL 实际执行: WHERE "status" = $1 AND "age" >= $2
```

如果你需要通过 `ToSQLAndArgs()` 预览生成的 SQL，在 PostgreSQL 下返回的占位符也会是 `$1`、`$2` 形式。

### 正则匹配（Regexp）

MySQL 使用 `REGEXP` 关键字，PostgreSQL 使用 `~` 操作符。`gormx` 会在运行时自动检测数据库类型，选择正确的操作符，**调用方无需区分**：

```go
query, user := gormx.NewQuery[User]()
query.Regexp(&user.Name, `^[A-Z]`)

// MySQL:      WHERE name REGEXP '^[A-Z]'
// PostgreSQL: WHERE name ~ '^[A-Z]'
```

### 完整示例

以下示例展示了在 PostgreSQL 下使用 gormx 完成常见业务场景的完整流程：

```go
package main

import (
    "fmt"

    "github.com/shijl0925/go-toolkits/gormx"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type Article struct {
    ID      uint   `gorm:"primaryKey"`
    Title   string `gorm:"column:title"`
    Author  string `gorm:"column:author"`
    Status  int    `gorm:"column:status"`
    ViewCnt int    `gorm:"column:view_cnt"`
}

type ArticleRepo struct {
    gormx.BaseRepo[Article]
}

func (r *ArticleRepo) FindPublished() ([]Article, error) {
    return r.SelectListByOpts(gormx.Where("status = ?", 1))
}

func main() {
    dsn := "host=127.0.0.1 user=postgres password=postgres dbname=demo port=5432 sslmode=disable TimeZone=Asia/Shanghai"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic(err)
    }
    gormx.Init(db)

    _ = db.AutoMigrate(&Article{})

    repo := &ArticleRepo{}

    // 1. 条件查询：标题以大写字母开头、状态为已发布、浏览量超过 100
    query, a := gormx.NewQuery[Article]()
    query.
        Regexp(&a.Title, `^[A-Z]`).
        Eq(&a.Status, 1).
        Gt(&a.ViewCnt, 100).
        OrderDesc(&a.ViewCnt).
        Limit(10)

    var articles []Article
    if err := query.Find(&articles); err != nil {
        panic(err)
    }
    fmt.Println("Found articles:", len(articles))

    // 2. 分页查询（第 1 页，每页 20 条）
    list, total, err := repo.SelectPage(1, 20, gormx.Where("status = ?", 1))
    if err != nil {
        panic(err)
    }
    fmt.Printf("Total: %d, Page: %v\n", total, list)

    // 3. 预览生成的 PostgreSQL SQL（占位符为 $1、$2 ...）
    q2, a2 := gormx.NewQuery[Article]()
    q2.Eq(&a2.Status, 1).In(&a2.Author, []string{"Alice", "Bob"})
    sql, args := q2.ToSQLAndArgs()
    fmt.Println("SQL:", sql)
    fmt.Println("Args:", args)

    // 4. 聚合查询：统计各状态文章数量
    q3, a3 := gormx.NewQuery[Article]()
    q3.GroupBy(&a3.Status).Having("COUNT(*) > ?", 5)
    var result []struct {
        Status int
        Count  int
    }
    if err := q3.Scan(&result); err != nil {
        panic(err)
    }
}
```

### 注意事项

- **表名引号**：PostgreSQL 默认区分大小写并使用双引号，GORM 会自动处理；建议 `gorm:"column:xxx"` 使用全小写或 snake_case 命名。
- **JSONB / 数组等扩展类型**：如需使用 PostgreSQL 特有的字段类型（`jsonb`、`text[]` 等），在模型定义上配合 `gorm.io/datatypes` 等扩展包即可，gormx 的查询构建器仍然适用。
- **`Regexp` 区分大小写**：PostgreSQL 的 `~` 默认区分大小写；如需忽略大小写，可改用 `~*` 操作符，此时通过 `RawRows` 或 `Scan` 手写 WHERE 子句即可。

