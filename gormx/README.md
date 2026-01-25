# gormx

基于 GORM 的增强工具库，提供了便捷的数据库操作封装，包括查询构建器、仓库模式、服务层抽象等功能。

## 特性

- **链式查询构建器**: 提供流畅的 API 进行复杂查询构建
- **仓库模式**: 封装基础的 CRUD 操作
- **服务层抽象**: 提供标准的服务层接口和实现
- **关联关系管理**: 支持一对一、一对多、多对多关系管理
- **反射优化**: 使用缓存机制提高字段名解析性能

## 安装

```bash
go get github.com/shijl0925/go-toolkits/gormx
```


## 快速开始

### 初始化

```go
package main

import (
    "github.com/shijl0925/go-toolkits/gormx"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func main() {
    dsn := "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }
    
    // 初始化全局 DB 实例
    gormx.Init(db)
}
```


### 查询构建器

```go
package main

import (
    "github.com/shijl0925/go-toolkits/gormx"
)

type User struct {
    ID       uint   `gorm:"primaryKey"`
    Name     string `gorm:"column:name"`
    Email    string `gorm:"column:email"`
    Age      int    `gorm:"column:age"`
    Status   int    `gorm:"column:status"`
}

func example() {
    // 创建查询构建器
    query, user := gormx.NewQuery[User]()
    
    // 构建查询条件
    query.Eq(&user.Name, "John").
          Gt(&user.Age, 18).
          In(&user.Status, []int{1, 2, 3}).
          OrderDesc(&user.ID)
    
    var users []User
    err := query.Find(&users)
    if err != nil {
        // 处理错误
    }
    
    // 也可以转换为原生 SQL 查看
    sql, args := query.ToSQLAndArgs()
    println("SQL:", sql)
    println("Args:", args)
}
```


### 仓库模式

```go
// 定义仓库接口
type IUserRepo interface {
    gormx.IBaseRepo[User] // 嵌入基础接口
    
    // 扩展用户相关业务方法
    FindByStatus(status int) ([]User, error)
}

// 实现仓库结构
type UserRepo struct {
    gormx.BaseRepo[User] // 继承基础仓库功能
}

// 实现自定义方法
func (r *UserRepo) FindByStatus(status int) ([]User, error) {
    query, user := gormx.NewQuery[User]()
    query.Eq(&user.Status, status)
    return r.SelectListByOpts(query.ToOptions()...)
}

func NewIUserRepo() IUserRepo {
    return &UserRepo{}
}
```


### 服务层

```go
// 定义服务接口
type IUserService interface {
    gormx.IBaseService[User] // 嵌入基础服务接口
    // 扩展用户相关业务方法
    GetUserWithDetails(id int) (*User, error)
}

// 实现服务结构
type UserService struct {
    *gormx.ServiceImplement[IUserRepo, User]
}

func NewIUserService(repo IUserRepo) IUserService {
    return &UserService{
        ServiceImplement: gormx.NewServiceImplement[IUserRepo, User](repo),
    }
}
```


## 功能详解

### 1. 查询构建器

#### 基本查询条件

```go
query, user := gormx.NewQuery[User]()

// 等于
query.Eq(&user.Name, "John")

// 不等于
query.Ne(&user.Status, 0)

// 大于
query.Gt(&user.Age, 18)

// 小于
query.Lt(&user.Salary, 10000)

// 大于等于
query.Ge(&user.Age, 18)

// 小于等于
query.Le(&user.Salary, 10000)

// 包含
query.Like(&user.Name, "%admin%")

// 范围查询
query.Between(&user.Age, 18, 65)

// IN 查询
query.In(&user.Status, []int{1, 2, 3})
query.NotIn(&user.Status, []int{4, 5, 6})

// 空值查询
query.IsNull(&user.Email)
query.IsNotNull(&user.Phone)
```


#### 子查询

```go
// 使用 SQL 子查询
query.InSql(&user.ID, "SELECT user_id FROM orders WHERE total > ?", 1000)

// 使用 GORM 子查询
subQuery, order := gormx.NewQuery[Order]()
subQuery.Select(&order.UserID).Gt(&order.Total, 1000)
query.SubQueryIn(&user.ID, GetDb(subQuery.ToOptions()...).Model(&Order{}).Select(&order.UserID))
```


#### 分组聚合

```go
query, user := gormx.NewQuery[User]()

// 分组
query.GroupBy(&user.Status)

// 排序
query.OrderAsc(&user.Name)
query.OrderDesc(&user.Age)

// 聚合函数
query.Having("COUNT(*) > ?", 10)

// 限制结果数量
query.Limit(10).Offset(20)
```


### 2. 仓库模式

#### 基础 CRUD 操作

```go
repo := NewIUserRepo()

// 查询
user, err := repo.SelectOneById(1)
users, err := repo.SelectListByOpts(opts...)
count, err := repo.SelectCount(opts...)

// 插入
err := repo.Insert(&user)

// 更新
err := repo.Update(&user)
err := repo.UpdateById(1, map[string]interface{}{"name": "new_name"})

// 删除
err := repo.DeleteById(1)
err := repo.DeleteBatchIds([]int{1, 2, 3})
```


#### 高级功能

```go
// 上下文操作
err := repo.InsertInBatches(users, 100) // 批量插入

// 增量操作
err := repo.Increment(1, "login_count", 1) // 自增
err := repo.Decrement(1, "balance", 100)   // 自减

// 条件更新
err := repo.Upsert(&user, map[string]interface{}{"name": "new_name"})
```


### 3. 服务层

服务层提供业务逻辑封装：

```go
type UserService struct {
    *gormx.ServiceImplement[IUserRepo, User]
}

// 基础服务方法
func (s *UserService) Create(user *User) error {
    return s.Repository.Insert(user)
}

func (s *UserService) UpdateUser(id int, updates map[string]interface{}) error {
    return s.Repository.UpdateById(id, updates)
}

func (s *UserService) GetUser(id int) (*User, error) {
    return s.Repository.SelectOneById(id)
}
```


### 4. 关联关系管理

```go
// 创建关联管理器
manager := gormx.NewAssociationManager[user, role](currentUser, "Roles")

// 添加关联
err := manager.Add(role1, role2)

// 移除关联
err := manager.Remove(role1)

// 清空关联
err := manager.Clear()

// 设置关联（替换所有）
err := manager.Set([]role{role1, role2})

// 获取关联数量
count, err := manager.Count()

// 获取所有关联数据
roles, err := manager.All()
```


## API 参考

### Query[T]

查询构建器的主要方法：

| 方法 | 说明 |
|------|------|
| `Eq(field interface{}, value interface{})` | 等于条件 |
| `Ne(field interface{}, value interface{})` | 不等于条件 |
| `Gt(field interface{}, value interface{})` | 大于条件 |
| `Lt(field interface{}, value interface{})` | 小于条件 |
| `Ge(field interface{}, value interface{})` | 大于等于条件 |
| `Le(field interface{}, value interface{})` | 小于等于条件 |
| `Like(field interface{}, value interface{})` | LIKE 条件 |
| `In(field interface{}, values interface{})` | IN 条件 |
| `NotIn(field interface{}, values interface{})` | NOT IN 条件 |
| `Between(field interface{}, start, end interface{})` | BETWEEN 条件 |
| `IsNull(field interface{})` | IS NULL 条件 |
| `IsNotNull(field interface{})` | IS NOT NULL 条件 |
| `OrderAsc(field interface{})` | 升序排序 |
| `OrderDesc(field interface{})` | 降序排序 |
| `GroupBy(field interface{})` | 分组 |
| `Limit(limit int)` | 限制结果数量 |
| `Offset(offset int)` | 偏移量 |
| `Having(condition string, args ...interface{})` | HAVING 条件 |

### BaseRepo[T]

基础仓库方法：

| 方法 | 说明 |
|------|------|
| `SelectOneById(id int)` | 根据 ID 查询单条记录 |
| `SelectOneByOpts(opts ...DBOption)` | 根据选项查询单条记录 |
| `SelectListByOpts(opts ...DBOption)` | 根据选项查询多条记录 |
| `SelectCount(opts ...DBOption)` | 统计记录数量 |
| `Insert(item *T)` | 插入单条记录 |
| `Update(item *T)` | 更新记录 |
| `DeleteById(id int)` | 根据 ID 删除记录 |
| `Sum(field string, opts ...DBOption)` | 求和 |
| `Max(field string, opts ...DBOption)` | 最大值 |
| `Min(field string, opts ...DBOption)` | 最小值 |
| `Avg(field string, opts ...DBOption)` | 平均值 |

### ServiceImplement[T, R]

基础服务实现：

| 方法 | 说明 |
|------|------|
| `GetRepository()` | 获取仓库实例 |
| `GetById(id int) (R, error)` | 根据 ID 获取记录 |
| `Save(entity *R) error` | 保存实体 |
| `RemoveById(id int) error` | 根据 ID 删除实体 |

## 最佳实践

1. **合理使用字段标识**:
    - 推荐使用结构体字段指针，类型安全
    - 字符串形式适合动态字段查询

2. **缓存利用**:
    - 查询构建器内部使用缓存优化字段名解析
    - 避免重复的反射操作

3. **错误处理**:
    - 每个操作都应处理可能的错误
    - 使用事务保证数据一致性

4. **SQL 注入防护**:
    - 所有参数都会被自动转义
    - 避免直接拼接 SQL 字符串

## 贡献

欢迎提交 Issue 和 Pull Request 来改进此项目。

## 许可证

MIT License