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

在调用 `GetDb`、`BaseRepo`、`Query[T]` 等能力前，需要先初始化全局 `*gorm.DB`：

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

