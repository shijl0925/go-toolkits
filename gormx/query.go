package gormx

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/shijl0925/go-toolkits/stringx"
	"gorm.io/gorm"
	"log"
	"reflect"
	"strings"
	"sync"
)

type Query[T any] struct {
	opts     []DBOption
	instance *T

	fieldCache map[uintptr]string // 字段地址到字段名的缓存
	cacheMutex sync.RWMutex       // 保护缓存的读写锁

	validFields map[string]bool // 有效字段名缓存
	fieldsOnce  sync.Once       // 确保字段列表只计算一次
}

func NewQuery[T any]() (*Query[T], *T) {
	var instance T
	return &Query[T]{
		instance:    &instance,
		fieldCache:  make(map[uintptr]string),
		validFields: make(map[string]bool),
	}, &instance
}

// ToOptions 转换为 DBOption 列表
func (q *Query[T]) ToOptions() []DBOption {
	return q.opts
}

// ToSQLAndArgs 生成并返回查询的SQL语句和参数值（不执行）
// 返回SQL语句和参数值
// 如何使用:
//
//	query, r := repo.NewQuery[model.Role]()
//	query.Eq(&r.Status, 1)
//	query.Like(&r.Name, "admin")
//	sql, args := query.ToSQLAndArgs()
//	fmt.Printf("Generated SQL: %v, args: %v\n", sql, args)
func (q *Query[T]) ToSQLAndArgs() (string, []interface{}) {
	opts := q.ToOptions()
	db := GetDb(opts...).Model(new(T))
	// 使用DryRun模式生成SQL但不执行
	stmt := db.Session(&gorm.Session{DryRun: true}).Find(q.instance).Statement
	return stmt.SQL.String(), stmt.Vars
}

// getModelFieldNames 获取模型的所有字段名
func (q *Query[T]) getModelFieldNames() []string {
	var fieldNames []string

	// 获取结构体类型
	structType := reflect.TypeOf(*q.instance)
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}

	// 递归处理所有字段
	fieldNames = q.extractFieldNames(structType)
	return fieldNames
}

// extractFieldNames 递归提取字段名
func (q *Query[T]) extractFieldNames(structType reflect.Type) []string {
	var fieldNames []string

	// 遍历所有字段
	for i := 0; i < structType.NumField(); i++ {
		structField := structType.Field(i)

		// 跳过未导出字段
		if !structField.IsExported() {
			continue
		}

		// 如果是嵌套结构体且匿名字段，递归处理
		if structField.Type.Kind() == reflect.Struct && structField.Anonymous {
			// 递归处理匿名嵌套结构体的字段
			nestedFieldNames := q.extractFieldNames(structField.Type)
			fieldNames = append(fieldNames, nestedFieldNames...)
			continue
		}

		// 获取字段名
		var fieldName string

		if tag := structField.Tag.Get("gorm"); tag != "" {
			// 获取gorm标签中的column名称
			if columnName := extractColumnFromGormTag(tag); columnName != "" {
				fieldName = columnName
			} else {
				// 如果没有gorm标签，使用字段名的蛇形命名
				fieldName = stringx.ToSnake(structField.Name)
			}
		} else {
			// 如果没有gorm标签，使用字段名的蛇形命名
			fieldName = stringx.ToSnake(structField.Name)
		}

		fieldNames = append(fieldNames, fieldName)
	}

	return fieldNames
}

// isValidModelField 检查字段是否属于模型的有效字段（带缓存）
// 修改说明：
//  1. 移除了对带 "." 字段名的复杂校验逻辑。
//  2. 新增：如果字段名包含 ".", 直接认为它是有效的。
//     这样做的原因是：在 JOIN 查询中，用户可能会引用其他表的字段（如 "users.name"），
//     这些字段显然不在当前模型 (T) 中，但却是合法的 SQL 字段引用。
//     因此，我们信任用户提供的 "table.column" 格式是正确的。
func (q *Query[T]) isValidModelField(fieldName string) bool {
	// 使用 Once 确保字段列表只计算一次
	q.fieldsOnce.Do(func() {
		validFields := q.getModelFieldNames()
		q.cacheMutex.Lock()
		for _, field := range validFields {
			q.validFields[field] = true
		}
		q.cacheMutex.Unlock()
	})

	// 读取缓存的结果
	q.cacheMutex.RLock()
	defer q.cacheMutex.RUnlock()

	// 直接检查字段名
	if q.validFields[fieldName] {
		return true
	}

	// 如果字段名包含 ".", 则检查 "." 后面的部分是否是有效字段
	if strings.Contains(fieldName, ".") {
		// 简单的格式检查：确保只有一个点且不以点开头或结尾
		parts := strings.Split(fieldName, ".")
		if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
			return true
		}
	}

	return false
}

// 处理字段标识符输入
func (q *Query[T]) resolveFieldName(field interface{}) string {
	var validFieldName string
	switch v := field.(type) {
	case string:
		// 检查字段名是否是有效字段
		if q.isValidModelField(v) {
			validFieldName = v
		}
	case nil:
		validFieldName = ""
	default:
		validFieldName = q.getFieldNameByReflection(field)
	}

	// 只在字段名确实无效时才打印警告（即输入是字符串但无效）
	if validFieldName == "" {
		if str, ok := field.(string); ok && str != "" {
			log.Printf("invalid field name: " + str)
		}
	}
	return validFieldName
}

// 通过反射获取字段名
func (q *Query[T]) getFieldNameByReflection(field interface{}) string {
	// 检查是否为指针类型
	fieldPtr := reflect.ValueOf(field)
	if fieldPtr.Kind() != reflect.Ptr {
		return ""
	}

	// 获取结构体实例
	structType := reflect.TypeOf(*q.instance)
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}

	// 获取字段指针地址
	targetAddr := fieldPtr.Pointer()

	// 使用统一的实例获取字段信息
	instanceValue := reflect.ValueOf(q.instance).Elem()

	// 递归查找字段
	return q.findFieldByAddress(structType, instanceValue, targetAddr)
}

// 递归查找字段地址匹配的字段名（带缓存）
func (q *Query[T]) findFieldByAddress(structType reflect.Type, instanceValue reflect.Value, targetAddr uintptr) string {
	// 先尝试从缓存中获取
	q.cacheMutex.RLock()
	if cachedName, exists := q.fieldCache[targetAddr]; exists {
		q.cacheMutex.RUnlock()
		return cachedName
	}
	q.cacheMutex.RUnlock()

	// 缓存未命中，执行原始查找逻辑
	for i := 0; i < structType.NumField(); i++ {
		structField := structType.Field(i)
		fieldValue := instanceValue.Field(i)

		// 如果是嵌套结构体，递归处理
		if structField.Type.Kind() == reflect.Struct && structField.Anonymous {
			if nestedName := q.findFieldByAddress(structField.Type, fieldValue, targetAddr); nestedName != "" {
				// 找到结果后缓存
				q.cacheMutex.Lock()
				q.fieldCache[targetAddr] = nestedName
				q.cacheMutex.Unlock()

				return nestedName
			}
			continue
		}

		// 获取字段地址进行比较
		fieldAddr := fieldValue.Addr().Pointer()
		if fieldAddr == targetAddr {
			// 计算字段名
			var fieldName string

			if tag := structField.Tag.Get("gorm"); tag != "" {
				// 获取gorm标签中的column名称
				if columnName := extractColumnFromGormTag(tag); columnName != "" {
					fieldName = columnName
				} else {
					// 如果没有gorm标签，使用字段名的蛇形命名
					fieldName = stringx.ToSnake(structField.Name)
				}
			} else {
				// 如果没有gorm标签，使用字段名的蛇形命名
				fieldName = stringx.ToSnake(structField.Name)
			}

			// 缓存结果
			q.cacheMutex.Lock()
			q.fieldCache[targetAddr] = fieldName
			q.cacheMutex.Unlock()

			return fieldName
		}
	}
	return ""
}

// 从gorm标签中提取列名
func extractColumnFromGormTag(tag string) string {
	// 解析gorm标签，提取column部分
	parts := strings.Split(tag, ";")
	for _, part := range parts {
		if strings.HasPrefix(part, "column:") {
			return strings.TrimPrefix(part, "column:")
		}
	}
	return ""
}

// Eq 添加 = 条件
// 该方法用于添加等于条件到查询中，判断字段值是否等于指定值
//
// 参数:
//   - field: 要比较的字段，可以是字段指针(如 &user.Status)或字段名字符串(如 "status")
//   - value: 要比较的值
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 查询状态等于1的用户
//	query, user := repo.NewQuery[model.User]()
//	query.Eq(&user.Status, 1)
//	var users []model.User
//	err := query.Find(&users)
//	// 生成SQL: SELECT * FROM tb_users WHERE status = 1
//
//	// 查询用户名等于"vben"的用户
//	query, u := repo.NewQuery[model.User]()
//	query.Eq(&u.UserName, "vben")
//	var userList []model.User
//	err = query.Find(&userList)
//	// 生成SQL: SELECT * FROM tb_users WHERE user_name = 'vben'
//
//	// 使用字符串字段名
//	query, role := repo.NewQuery[model.Role]()
//	query.Eq("name", "admin")
//	var roles []model.Role
//	err = query.Find(&roles)
//	// 生成SQL: SELECT * FROM tb_cmdb_roles WHERE name = 'admin'
//
// 注意事项:
//   - 字段名会经过验证，无效字段会产生警告但不会导致panic
//   - 可以与其他查询条件组合使用
//   - value参数的类型应该与字段类型兼容
func (q *Query[T]) Eq(field interface{}, value interface{}) *Query[T] {
	fieldName := q.resolveFieldName(field)
	q.opts = append(q.opts, Where(fieldName+" = ?", value))
	return q
}

// Ne 添加 != 条件
// 该方法用于添加不等于条件到查询中，判断字段值是否不等于指定值
//
// 参数:
//   - field: 要比较的字段，可以是字段指针(如 &user.Status)或字段名字符串(如 "status")
//   - value: 要比较的值
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 注意事项:
//   - 字段名会经过验证，无效字段会产生警告但不会导致panic
//   - 可以与其他查询条件组合使用
//   - value参数的类型应该与字段类型兼容
func (q *Query[T]) Ne(field interface{}, value interface{}) *Query[T] {
	fieldName := q.resolveFieldName(field)
	q.opts = append(q.opts, Where(fieldName+" <> ?", value))
	return q
}

// Gt 添加 > 条件
// 该方法用于添加大于条件到查询中，判断字段值是否大于指定值
//
// 参数:
//   - field: 要比较的字段，可以是字段指针(如 &user.Age)或字段名字符串(如 "age")
//   - value: 要比较的值
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 注意事项:
//   - 字段名会经过验证，无效字段会产生警告但不会导致panic
//   - 可以与其他查询条件组合使用
//   - value参数的类型应该与字段类型兼容
func (q *Query[T]) Gt(field interface{}, value interface{}) *Query[T] {
	fieldName := q.resolveFieldName(field)
	q.opts = append(q.opts, Where(fieldName+" > ?", value))
	return q
}

// Ge 添加 >= 条件
// 该方法用于添加大于等于条件到查询中，判断字段值是否大于等于指定值
//
// 参数:
//   - field: 要比较的字段，可以是字段指针(如 &user.Age)或字段名字符串(如 "age")
//   - value: 要比较的值
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 注意事项:
//   - 字段名会经过验证，无效字段会产生警告但不会导致panic
//   - 可以与其他查询条件组合使用
//   - value参数的类型应该与字段类型兼容
func (q *Query[T]) Ge(field interface{}, value interface{}) *Query[T] {
	fieldName := q.resolveFieldName(field)
	q.opts = append(q.opts, Where(fieldName+" >= ?", value))
	return q
}

// Lt 添加 < 条件
// 该方法用于添加小于条件到查询中，判断字段值是否小于指定值
//
// 参数:
//   - field: 要比较的字段，可以是字段指针(如 &user.Age)或字段名字符串(如 "age")
//   - value: 要比较的值
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 注意事项:
//   - 字段名会经过验证，无效字段会产生警告但不会导致panic
//   - 可以与其他查询条件组合使用
//   - value参数的类型应该与字段类型兼容
func (q *Query[T]) Lt(field interface{}, value interface{}) *Query[T] {
	fieldName := q.resolveFieldName(field)
	q.opts = append(q.opts, Where(fieldName+" < ?", value))
	return q
}

// Le 添加 <= 条件
// 该方法用于添加小于等于条件到查询中，判断字段值是否小于等于指定值
//
// 参数:
//   - field: 要比较的字段，可以是字段指针(如 &user.Age)或字段名字符串(如 "age")
//   - value: 要比较的值
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 注意事项:
//   - 字段名会经过验证，无效字段会产生警告但不会导致panic
//   - 可以与其他查询条件组合使用
//   - value参数的类型应该与字段类型兼容
func (q *Query[T]) Le(field interface{}, value interface{}) *Query[T] {
	fieldName := q.resolveFieldName(field)
	q.opts = append(q.opts, Where(fieldName+" <= ?", value))
	return q
}

// Like 添加 like 模糊匹配条件
// 该方法用于添加LIKE模糊匹配条件到查询中，判断字段值是否包含指定的字符串模式
//
// 参数:
//   - field: 要比较的字段，可以是字段指针(如 &user.Name)或字段名字符串(如 "name")
//   - value: 要匹配的字符串值，必须是string类型
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 查询用户名包含"admin"的用户
//	query, user := repo.NewQuery[model.User]()
//	query.Like(&user.UserName, "admin")
//	var users []model.User
//	err := query.Find(&users)
//	// 生成SQL: SELECT * FROM tb_users WHERE user_name LIKE '%admin%'
//
//	// 查询邮箱包含"gmail"的用户
//	query, u := repo.NewQuery[model.User]()
//	query.Like(&u.Email, "gmail")
//	var userList []model.User
//	err = query.Find(&userList)
//	// 生成SQL: SELECT * FROM tb_users WHERE email LIKE '%gmail%'
//
//	// 使用字符串字段名
//	query, role := repo.NewQuery[model.Role]()
//	query.Like("name", "manager")
//	var roles []model.Role
//	err = query.Find(&roles)
//	// 生成SQL: SELECT * FROM tb_cmdb_roles WHERE name LIKE '%manager%'
//
// 注意事项:
//   - value参数必须是string类型，否则会panic
//   - 字段名会经过验证，无效字段会产生警告但不会导致panic
//   - 可以与其他查询条件组合使用
//   - 自动在匹配值前后添加%通配符，实现包含匹配
func (q *Query[T]) Like(field interface{}, value interface{}) *Query[T] {
	strVal, ok := value.(string)
	if !ok {
		panic("like value must be a string")
	}
	fieldName := q.resolveFieldName(field)
	q.opts = append(q.opts, Where(fieldName+" LIKE ?", "%"+strVal+"%"))
	return q
}

// Regexp 添加正则表达式匹配条件
// 该方法用于添加正则表达式匹配条件到查询中，判断字段值是否匹配指定的正则表达式模式
//
// 参数:
//   - field: 要比较的字段，可以是字段指针(如 &user.Name)或字段名字符串(如 "name")
//   - pattern: 正则表达式模式字符串
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 查询用户名符合特定模式的用户(以字母开头，后面跟数字)
//	query, user := repo.NewQuery[model.User]()
//	query.Regexp(&user.UserName, "^[a-zA-Z][0-9]+$")
//	var users []model.User
//	err := query.Find(&users)
//	// 生成SQL: SELECT * FROM tb_users WHERE user_name REGEXP '^[a-zA-Z][0-9]+$'
//
//	// 查询邮箱匹配特定域名的用户
//	query, u := repo.NewQuery[model.User]()
//	query.Regexp(&u.Email, "@gmail\\.com$")
//	var userList []model.User
//	err = query.Find(&userList)
//	// 生成SQL: SELECT * FROM tb_users WHERE email REGEXP '@gmail\.com$'
//
//	// 使用字符串字段名
//	query, role := repo.NewQuery[model.Role]()
//	query.Regexp("remark", ".*管理员.*")
//	var roles []model.Role
//	err = query.Find(&roles)
//	// 生成SQL: SELECT * FROM tb_cmdb_roles WHERE remark REGEXP '.*管理员.*'
//
// 注意事项:
//   - pattern参数必须是有效的正则表达式模式
//   - 字段名会经过验证，无效字段会产生警告但不会导致panic
//   - 可以与其他查询条件组合使用
//   - 不同数据库对正则表达式的支持可能有所不同(如MySQL使用REGEXP，PostgreSQL使用~)
func (q *Query[T]) Regexp(field interface{}, pattern string) *Query[T] {
	fieldName := q.resolveFieldName(field)

	// 根据数据库类型选择合适的正则表达式操作符
	// 默认使用 REGEXP (适用于 MySQL)
	operator := "REGEXP"

	// 如果需要支持其他数据库，可以根据需要添加判断逻辑
	// 例如对于 PostgreSQL 可以使用 ~ 操作符

	q.opts = append(q.opts, Where(fieldName+" "+operator+" ?", pattern))
	return q
}

// IsNull 添加 is null 条件
// 该方法用于添加IS NULL条件到查询中，判断字段值是否为NULL
//
// 参数:
//   - field: 要检查的字段，可以是字段指针(如 &user.Email)或字段名字符串(如 "email")
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 查询所有邮箱为空的用户
//	query, user := repo.NewQuery[model.User]()
//	query.IsNull(&user.Email)
//	var users []model.User
//	err := query.Find(&users)
//	// 生成SQL: SELECT * FROM tb_users WHERE email IS NULL
//
//	// 查询所有没有创建时间的帖子
//	query, post := repo.NewQuery[model.Post]()
//	query.IsNull(&post.CreatedAt)
//	var posts []model.Post
//	err = query.Find(&posts)
//	// 生成SQL: SELECT * FROM tb_posts WHERE created_at IS NULL
//
//	// 使用字符串字段名
//	query, u := repo.NewQuery[model.User]()
//	query.IsNull("phone")
//	var userList []model.User
//	err = query.Find(&userList)
//
// 注意事项:
//   - 该条件用于筛选字段值为NULL的记录
//   - 字段名会经过验证，无效字段会产生警告但不会导致panic
//   - 可以与其他查询条件组合使用
func (q *Query[T]) IsNull(field interface{}) *Query[T] {
	fieldName := q.resolveFieldName(field)
	q.opts = append(q.opts, Where(fieldName+" IS NULL"))
	return q
}

// IsNotNull 添加 is not null 条件
// 该方法用于添加IS NOT NULL条件到查询中，判断字段值是否不为NULL
//
// 参数:
//   - field: 要检查的字段，可以是字段指针(如 &user.Email)或字段名字符串(如 "email")
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 查询所有邮箱不为空的用户
//	query, user := repo.NewQuery[model.User]()
//	query.IsNotNull(&user.Email)
//	var users []model.User
//	err := query.Find(&users)
//	// 生成SQL: SELECT * FROM tb_users WHERE email IS NOT NULL
//
//	// 查询所有有创建时间的帖子
//	query, post := repo.NewQuery[model.Post]()
//	query.IsNotNull(&post.CreatedAt)
//	var posts []model.Post
//	err = query.Find(&posts)
//	// 生成SQL: SELECT * FROM tb_posts WHERE created_at IS NOT NULL
//
//	// 使用字符串字段名
//	query, u := repo.NewQuery[model.User]()
//	query.IsNotNull("phone")
//	var userList []model.User
//	err = query.Find(&userList)
//
// 注意事项:
//   - 该条件用于筛选字段值不为NULL的记录
//   - 字段名会经过验证，无效字段会产生警告但不会导致panic
//   - 可以与其他查询条件组合使用
func (q *Query[T]) IsNotNull(field interface{}) *Query[T] {
	fieldName := q.resolveFieldName(field)
	q.opts = append(q.opts, Where(fieldName+" IS NOT NULL"))
	return q
}

// buildInPlaceholders 构建IN查询的占位符字符串和参数列表
// 该方法用于处理切片类型的参数，生成对应的占位符和参数列表
//
// 参数:
//   - value: 要处理的切片值
//
// 返回:
//   - placeholderStr: 占位符字符串，如 "?,?,?"
//   - args: 参数列表
func (q *Query[T]) buildInPlaceholders(value interface{}) (string, []interface{}) {
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Slice {
		placeholders := make([]string, rv.Len())
		args := make([]interface{}, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			placeholders[i] = "?"
			args[i] = rv.Index(i).Interface()
		}
		return strings.Join(placeholders, ","), args
	}
	// 非切片情况，返回单个占位符和原始值
	return "?", []interface{}{value}
}

// In 添加 in 条件 IN (?, ?, ?)
// 该方法用于添加IN条件到查询中，判断字段值是否在指定的值列表中
//
// 参数:
//   - field: 要比较的字段，可以是字段指针(如 &user.Status)或字段名字符串(如 "status")
//   - value: 要匹配的值列表，可以是数组、切片或任何可迭代的对象
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 查询状态为1或2的用户
//	query, user := repo.NewQuery[model.User]()
//	query.In(&user.Status, []int{1, 2})
//	var users []model.User
//	err := query.Find(&users)
//	// 生成SQL: SELECT * FROM tb_users WHERE status IN (1, 2)
//
//	// 查询特定ID列表的用户
//	query, u := repo.NewQuery[model.User]()
//	query.In(&u.ID, []int{1, 2, 3, 4, 5})
//	var userList []model.User
//	err = query.Find(&userList)
//	// 生成SQL: SELECT * FROM tb_users WHERE id IN (1, 2, 3, 4, 5)
//
//	// 使用字符串字段名
//	query, role := repo.NewQuery[model.Role]()
//	query.In("name", []string{"admin", "user", "guest"})
//	var roles []model.Role
//	err = query.Find(&roles)
//	// 生成SQL: SELECT * FROM tb_cmdb_roles WHERE name IN ('admin', 'user', 'guest')
//
// 注意事项:
//   - value参数可以是任何可迭代的类型，如数组、切片等
//   - 字段名会经过验证，无效字段会产生警告但不会导致panic
//   - 可以与其他查询条件组合使用
//   - 如果value为空，会生成IN ()条件，这在某些数据库中可能导致语法错误
func (q *Query[T]) In(field interface{}, value interface{}) *Query[T] {
	fieldName := q.resolveFieldName(field)
	// 生成占位符字符串和参数列表
	placeholderStr, args := q.buildInPlaceholders(value)
	q.opts = append(q.opts, func(db *gorm.DB) *gorm.DB {
		return db.Where(fieldName+" IN ("+placeholderStr+")", args...)
	})
	return q
}

// NotIn 添加 not in 条件 NOT IN (?, ?, ?)
// 该方法用于添加NOT IN条件到查询中，判断字段值是否不在指定的值列表中
//
// 参数:
//   - field: 要比较的字段，可以是字段指针(如 &user.Status)或字段名字符串(如 "status")
//   - value: 要排除的值列表，可以是数组、切片或任何可迭代的对象
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 查询状态不为1或2的用户
//	query, user := repo.NewQuery[model.User]()
//	query.NotIn(&user.Status, []int{1, 2})
//	var users []model.User
//	err := query.Find(&users)
//	// 生成SQL: SELECT * FROM tb_users WHERE status NOT IN (1, 2)
//
//	// 查询除了特定ID列表外的用户
//	query, u := repo.NewQuery[model.User]()
//	query.NotIn(&u.ID, []int{1, 2, 3, 4, 5})
//	var userList []model.User
//	err = query.Find(&userList)
//	// 生成SQL: SELECT * FROM tb_users WHERE id NOT IN (1, 2, 3, 4, 5)
//
//	// 使用字符串字段名
//	query, role := repo.NewQuery[model.Role]()
//	query.NotIn("name", []string{"admin", "user", "guest"})
//	var roles []model.Role
//	err = query.Find(&roles)
//	// 生成SQL: SELECT * FROM tb_cmdb_roles WHERE name NOT IN ('admin', 'user', 'guest')
//
// 注意事项:
//   - value参数可以是任何可迭代的类型，如数组、切片等
//   - 字段名会经过验证，无效字段会产生警告但不会导致panic
//   - 可以与其他查询条件组合使用
//   - 如果value为空，会生成NOT IN ()条件，这在某些数据库中可能导致语法错误
func (q *Query[T]) NotIn(field interface{}, value interface{}) *Query[T] {
	fieldName := q.resolveFieldName(field)
	// 生成占位符字符串和参数列表
	placeholderStr, args := q.buildInPlaceholders(value)
	q.opts = append(q.opts, func(db *gorm.DB) *gorm.DB {
		return db.Where(fieldName+" NOT IN ("+placeholderStr+")", args...)
	})
	return q
}

// Between 添加 between 条件
// 该方法用于添加BETWEEN条件到查询中，判断字段值是否在指定范围内（包含边界值）
//
// 参数:
//   - field: 要比较的字段，可以是字段指针(如 &user.Age)或字段名字符串(如 "age")
//   - start: 范围起始值
//   - end: 范围结束值
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 查询年龄在18到65岁之间的用户
//	query, user := repo.NewQuery[model.User]()
//	query.Between(&user.Age, 18, 65)
//	var users []model.User
//	err := query.Find(&users)
//	// 生成SQL: SELECT * FROM tb_users WHERE age BETWEEN 18 AND 65
//
//	// 查询创建时间在指定日期范围内的角色
//	query, role := repo.NewQuery[model.Role]()
//	startTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
//	endTime := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)
//	query.Between(&role.CreatedAt, startTime, endTime)
//	var roles []model.Role
//	err = query.Find(&roles)
//	// 生成SQL: SELECT * FROM tb_cmdb_roles WHERE created_at BETWEEN '2023-01-01 00:00:00' AND '2023-12-31 23:59:59'
//
//	// 使用字符串字段名
//	query, emp := repo.NewQuery[model.Employee]()
//	query.Between("salary", 5000, 10000)
//	var employees []model.Employee
//	err = query.Find(&employees)
//
// 注意事项:
//   - BETWEEN条件包含起始值和结束值（即 >= start AND <= end）
//   - 字段名会经过验证，无效字段会产生警告但不会导致panic
//   - 可以与其他查询条件组合使用
//   - start和end参数类型应该与字段类型兼容
func (q *Query[T]) Between(field interface{}, start, end interface{}) *Query[T] {
	fieldName := q.resolveFieldName(field)
	q.opts = append(q.opts, Where(fieldName+" BETWEEN ? AND ?", start, end))
	return q
}

// NotBetween 添加 not between 条件
// 该方法用于添加NOT BETWEEN条件到查询中，判断字段值是否不在指定范围内（不包含边界值）
//
// 参数:
//   - field: 要比较的字段，可以是字段指针(如 &user.Age)或字段名字符串(如 "age")
//   - start: 范围起始值
//   - end: 范围结束值
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 查询年龄不在18到65岁之间的用户
//	query, user := repo.NewQuery[model.User]()
//	query.NotBetween(&user.Age, 18, 65)
//	var users []model.User
//	err := query.Find(&users)
//	// 生成SQL: SELECT * FROM tb_users WHERE age NOT BETWEEN 18 AND 65
//
//	// 查询创建时间不在指定日期范围内的角色
//	query, role := repo.NewQuery[model.Role]()
//	startTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
//	endTime := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)
//	query.NotBetween(&role.CreatedAt, startTime, endTime)
//	var roles []model.Role
//	err = query.Find(&roles)
//	// 生成SQL: SELECT * FROM tb_cmdb_roles WHERE created_at NOT BETWEEN '2023-01-01 00:00:00' AND '2023-12-31 23:59:59'
//
//	// 使用字符串字段名
//	query, emp := repo.NewQuery[model.Employee]()
//	query.NotBetween("salary", 5000, 10000)
//	var employees []model.Employee
//	err = query.Find(&employees)
//
// 注意事项:
//   - NOT BETWEEN条件不包含起始值和结束值（即 < start OR > end）
//   - 字段名会经过验证，无效字段会产生警告但不会导致panic
//   - 可以与其他查询条件组合使用
//   - start和end参数类型应该与字段类型兼容
func (q *Query[T]) NotBetween(field interface{}, start, end interface{}) *Query[T] {
	fieldName := q.resolveFieldName(field)
	q.opts = append(q.opts, Where(fieldName+" NOT BETWEEN ? AND ?", start, end))
	return q
}

// resolveFieldNames 处理字段参数，支持字符串、字段指针和切片类型
func (q *Query[T]) resolveFieldNames(fields []interface{}) []string {
	var fieldNames []string

	for _, field := range fields {
		// 处理切片类型参数
		rv := reflect.ValueOf(field)
		if rv.Kind() == reflect.Slice {
			// 展开切片元素
			for i := 0; i < rv.Len(); i++ {
				element := rv.Index(i).Interface()
				fieldName := q.resolveFieldName(element)
				if fieldName != "" {
					fieldNames = append(fieldNames, fieldName)
				}
			}
		} else {
			// 处理普通参数
			fieldName := q.resolveFieldName(field)
			if fieldName != "" {
				fieldNames = append(fieldNames, fieldName)
			}
		}
	}

	return fieldNames
}

// OrderDesc 排序：ORDER BY 字段1,字段2 Desc (逆序)
// 该方法用于添加降序排序条件到查询中，相当于SQL中的ORDER BY ... DESC子句
//
// 参数:
//   - fields: 要排序的字段列表，支持字段指针(如 &user.CreatedAt)和字段名字符串(如 "created_at")
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	 query, user := repo.NewQuery[model.User]()
//	 query.Eq(&user.Status, 1)
//	 query.OrderDesc(&user.CreatedAt) // 按创建时间倒序
//		var users []model.User
//		err := query.Find(&users)
//		// 生成SQL: SELECT * FROM tb_users WHERE status = 1 ORDER BY created_at desc
//
// ...
//
//	 // 按状态降序，再按创建时间降序排列角色
//	 query, role := repo.NewQuery[model.Role]()
//	 query.OrderDesc(&role.Status, &role.CreatedAt)
//		var roles []model.Role
//		err = query.Find(&roles)
//		// 生成SQL: SELECT * FROM tb_cmdb_roles ORDER BY status desc, created_at desc
//
//		// 使用字符串字段名
//		query, u := repo.NewQuery[model.User]()
//		query.OrderDesc("status", "created_at")
//		var userList []model.User
//		err = query.Find(&userList)
//
// 注意事项:
//   - 可以传入多个字段进行多重排序
//   - 排序字段会按传入顺序依次应用
//   - 字段名会经过验证，无效字段会被忽略
//   - 可以与其他查询条件组合使用
func (q *Query[T]) OrderDesc(fields ...interface{}) *Query[T] {
	fieldNames := q.resolveFieldNames(fields)

	for _, fieldName := range fieldNames {
		// 创建 fieldName 的副本，避免闭包捕获问题
		field := fieldName
		if field == "" {
			continue
		}
		q.opts = append(q.opts, func(db *gorm.DB) *gorm.DB {
			return db.Order(field + " DESC")
		})
	}
	return q
}

// OrderAsc 排序：ORDER BY 字段1,字段2 Asc (顺序)
// 如何使用:
// query, user := repo.NewQuery[model.User]()
// query.Eq(&user.Status, 1)
// query.OrderAsc(&user.CreatedAt) // 按创建时间顺序
// ...
// // 按状态顺序，再按创建时间顺序排列角色
// query, role := repo.NewQuery[model.Role]()
// query.OrderAsc(&role.Status, &role.CreatedAt)
func (q *Query[T]) OrderAsc(fields ...interface{}) *Query[T] {
	fieldNames := q.resolveFieldNames(fields)

	for _, fieldName := range fieldNames {
		// 创建 fieldName 的副本，避免闭包捕获问题
		field := fieldName
		if field == "" {
			continue
		}
		q.opts = append(q.opts, func(db *gorm.DB) *gorm.DB {
			return db.Order(field + " ASC")
		})
	}
	return q
}

// First 查询一条数据
// 该方法用于执行查询并返回匹配条件的第一条记录，将结果填充到目标变量中
//
// 参数:
//   - destination: 指向目标变量的指针，用于存储查询结果，类型应为*T
//
// 返回:
//   - error: 错误信息，如果查询成功则返回nil
//
// 使用示例:
//
//	// 查询用户名为"vben"的用户
//	query, u := repo.NewQuery[model.User]()
//	query.Eq(&u.UserName, "vben")
//	var user model.User
//	err := query.First(&user)
//	if err != nil {
//		// 处理错误，如记录不存在等
//	}
//	// user现在包含用户名为"vben"的第一个用户的信息
//
//	// 查询特定ID的角色
//	query, r := repo.NewQuery[model.Role]()
//	query.Eq(&r.ID, 1)
//	var role model.Role
//	err = query.First(&role)
//	// role现在包含ID为1的角色信息
//
// 注意事项:
//   - destination必须是指向结构体的指针
//   - 可以与各种查询条件组合使用，如Eq、Ne、Like、In等
//   - 如果没有匹配的记录，会返回"gorm.ErrRecordNotFound"错误
//   - 查询结果是按数据库默认排序的第一条记录
func (q *Query[T]) First(destination *T) error {
	opts := q.ToOptions()
	db := GetDb(opts...)
	if db.Statement.Model == nil && db.Statement.Table == "" {
		db = db.Model(new(T))
	}

	return db.First(destination).Error
}

// Find 查询多条数据
// 该方法用于执行查询并返回匹配条件的所有记录，将结果填充到目标切片中
//
// 参数:
//   - destination: 指向目标切片的指针，用于存储查询结果，类型应为*[]T
//
// 返回:
//   - error: 错误信息，如果查询成功则返回nil
//
// 使用示例:
//
//	// 查询所有状态为1的用户
//	query, user := repo.NewQuery[model.User]()
//	query.Eq(&user.Status, 1)
//	var users []model.User
//	err := query.Find(&users)
//	if err != nil {
//		// 处理错误
//	}
//	// users现在包含所有状态为1的用户
//
//	// 根据用户ID查询帖子
//	query, p := repo.NewQuery[model.Post]()
//	userId := 1
//	query.Eq(&p.UserID, userId)
//	var posts []model.Post
//	err = query.Find(&posts)
//	if err != nil {
//		// 处理错误
//	}
//	// posts现在包含所有用户ID为1的帖子
//
//	// 复合条件查询
//	query, role := repo.NewQuery[model.Role]()
//	query.Eq(&role.Status, 1)
//	query.Like(&role.Name, "admin")
//	var roles []model.Role
//	err = query.Find(&roles)
//	// roles现在包含所有状态为1且名称包含"admin"的角色
//
// 注意事项:
//   - destination必须是指向切片的指针
//   - 可以与各种查询条件组合使用，如Eq、Ne、Like、In等
//   - 可以与排序、分页、分组等功能组合使用
//   - 如果没有匹配的记录，不会返回错误，而是返回空切片
func (q *Query[T]) Find(destination *[]T) error {
	opts := q.ToOptions()
	db := GetDb(opts...)
	if db.Statement.Model == nil && db.Statement.Table == "" {
		db = db.Model(new(T))
	}
	return db.Find(destination).Error
}

// Distinct 添加 distinct 函数
// 该方法用于添加DISTINCT关键字到查询中，对查询结果进行去重操作
//
// 参数:
//   - fields: 可选的字段列表，支持字段指针(如 &user.Status)和字段名字符串(如 "status")
//     如果不提供参数，则对整个查询结果去重
//     如果提供参数，则只对指定字段进行去重
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 对整个查询结果去重
//	query, user := repo.NewQuery[model.User]()
//	query.Distinct()
//	var users []model.User
//	err := query.Scan(&users)
//	if err != nil {
//		// 处理错误
//	}
//	// 生成SQL: SELECT DISTINCT * FROM tb_users
//
//	// 对特定字段去重
//	query, user := repo.NewQuery[model.User]()
//	query.Distinct(&user.Status)
//	var users []model.User
//	err = query.Scan(&users)
//	if err != nil {
//		// 处理错误
//	}
//	// 生成SQL: SELECT DISTINCT status FROM tb_users
//
//	// 对多个字段组合去重
//	query, role := repo.NewQuery[model.Role]()
//	query.Distinct(&role.Status, &role.Name)
//	var roles []model.Role
//	err = query.Scan(&roles)
//	// 生成SQL: SELECT DISTINCT status, name FROM tb_cmdb_roles
//
// 注意事项:
//   - 不带参数时，DISTINCT作用于所有字段
//   - 带参数时，DISTINCT只作用于指定字段
//   - 可以与其他查询条件组合使用，如Where、Order等
//   - 字段名会经过验证，无效字段会被忽略
func (q *Query[T]) Distinct(fields ...interface{}) *Query[T] {
	if len(fields) == 0 {
		// 无参数时，对整个查询去重
		q.opts = append(q.opts, func(db *gorm.DB) *gorm.DB {
			return db.Distinct("*")
		})
		return q
	}

	fieldNames := q.resolveFieldNames(fields)

	if len(fieldNames) > 0 {
		q.opts = append(q.opts, func(db *gorm.DB) *gorm.DB {
			// 将 []string 转换为 []interface{}
			args := make([]interface{}, len(fieldNames))
			for i, name := range fieldNames {
				args[i] = name
			}
			return db.Distinct(args...)
		})
	}

	return q
}

// Select 查询字段
// 该方法用于指定查询时需要返回的字段，相当于SQL中的SELECT子句
//
// 参数:
//   - fields: 要查询的字段列表，支持字段指针(如 &user.Name)和字段名字符串(如 "name")
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
// 只查询用户的名字和邮箱字段
//
//	    // 方法一
//		query, user := repo.NewQuery[model.User]()
//		query.Select(&user.Name, &user.Email)  // 等同于 query.Select("name", "email")
//		query.Eq(&user.Status, 1)  // 只查询活跃用户
//
//		var users []model.User
//		err := query.Scan(&users)
//		if err != nil {
//		  // 处理错误
//		}
//		// 生成SQL: SELECT name, email FROM users WHERE status = 1
//
//		// 方法二
//		query, user := repo.NewQuery[model.User]()
//		query.Select(&user.Name, &user.Email)
//		query.Eq(&user.Status, 1)
//		opts := query.ToOptions()
//		user, err := userService.List(opts...)
//		if err != nil {
//		  // 处理错误
//		}
//		// 生成SQL: SELECT name, email FROM users WHERE status = 1
//
// 注意事项:
//   - 如果不指定Select，则默认查询所有字段(*)
//   - 可以多次调用Select，后面的字段会追加到之前的选择中
//   - 字段名会经过验证，无效字段会被忽略
//   - 支持使用字段指针和字符串两种方式指定字段
func (q *Query[T]) Select(fields ...interface{}) *Query[T] {
	if len(fields) == 0 {
		return q
	}

	fieldNames := q.resolveFieldNames(fields)

	if len(fieldNames) > 0 {
		q.opts = append(q.opts, func(db *gorm.DB) *gorm.DB {
			return db.Select(fieldNames)
		})
	}

	return q
}

// Scan 扫描查询结果
func (q *Query[T]) Scan(destination interface{}) error {
	// 检查目标是否为nil
	if destination == nil {
		return errors.New("destination cannot be nil")
	}

	// 检查目标是否为零值（未初始化的指针）
	rv := reflect.ValueOf(destination)
	if !rv.IsValid() || (rv.Kind() == reflect.Ptr && rv.IsNil()) {
		return errors.New("cannot scan into nil or invalid destination")
	}

	opts := q.ToOptions()
	db := GetDb(opts...)
	// 确保查询有正确的模型上下文
	if db.Statement.Model == nil && db.Statement.Table == "" {
		db = db.Model(new(T))
	}
	return db.Scan(destination).Error
}

// RawRows 执行原生SQL查询并返回数据库行结果
// 该方法允许执行自定义SQL查询，返回*gorm.Rows结果集，适用于复杂查询场景。
//
// 参数:
//   - sql: 原生 SQL 查询语句，可包含占位符 ?
//   - args: SQL 查询参数，可变参数。按顺序替换SQL中的占位符
//
// 返回:
//   - *sql.Rows: 数据库行结果集，可以通过Rows.Scan方法将结果映射到结构体或变量中
//   - error: 错误信息，如果查询成功则返回nil
//
// 使用说明:
// query, _ := repo.NewQuery[model.User]()
// rows, err := query.RawRows("SELECT id, name FROM tb_users WHERE status = ?", 1)
//
//	if err != nil {
//		// 处理错误
//	}
//
// defer rows.Close()
//
//	for rows.Next() {
//		var id int
//		var name string
//
//		if err := rows.Scan(&id, &name); err != nil {
//			// 处理错误
//		}
//		// 处理结果
//		fmt.Printf("id: %d, name: %s\n", id, name)
//	}
//
// 注意事项:
// - 调用方需要负责关闭返回的Rows资源(通常使用defer rows.Close())
// - 需要手动处理每一行数据的扫描(使用rows.Scan())
// - SQL注入风险: 请务必对SQL查询参数进行参数化处理(占位符)而非字符串拼接，防止SQL注入攻击。
func (q *Query[T]) RawRows(sql string, args ...interface{}) (*sql.Rows, error) {
	opts := q.ToOptions()
	db := GetDb(opts...)
	// 确保查询有正确的模型上下文
	if db.Statement.Model == nil && db.Statement.Table == "" {
		db = db.Model(new(T))
	}
	return db.Raw(sql, args...).Rows()
}

// Pluck 查询单个字段值列表
// 该方法用于从查询结果中提取指定字段的值，返回一个包含该字段所有值的切片
//
// 参数:
//   - field: 要提取的字段，可以是字段指针(如 &user.Username)或字段名字符串(如 "username")
//   - destination: 目标切片变量，用于存储查询结果，类型应与字段类型匹配
//
// 返回:
//   - error: 错误信息，如果查询成功则返回nil
//
// 如何使用:
//
//	// 提取所有用户的用户名
//	query, user := repo.NewQuery[model.User]()
//	query.Eq(&user.Status, 1)
//	var usernames []string
//	err := query.Pluck(&user.Username, &usernames)
//	if err != nil {
//		// 处理错误
//	}
//	// usernames现在包含所有状态为1的用户的用户名
//
//	// 提取所有用户ID
//	query, u := repo.NewQuery[model.User]()
//	var ids []int
//	err = query.Pluck("id", &ids)
//	if err != nil {
//		// 处理错误
//	}
//	// ids现在包含所有用户的ID
//
//	// 提取用户邮箱列表
//	query, usr := repo.NewQuery[model.User]()
//	query.IsNotNull(&usr.Email)
//	var emails []string
//	err = query.Pluck(&usr.Email, &emails)
//	// emails现在包含所有有邮箱的用户的邮箱地址
//
// 注意事项:
//   - destination必须是指向切片的指针
//   - 字段类型应与目标切片元素类型兼容
//   - 如果字段无效，可能会生成字段名为空的查询条件
//   - 可以与其他查询条件组合使用，如Where、Order等
func (q *Query[T]) Pluck(field interface{}, destination interface{}) error {
	// 检查目标是否为nil
	if destination == nil {
		return errors.New("destination cannot be nil")
	}

	// 检查目标是否为指针
	rv := reflect.ValueOf(destination)
	if rv.Kind() != reflect.Ptr {
		return errors.New("destination must be a pointer to a slice")
	}

	// 检查目标是否为指向切片的指针
	if rv.Elem().Kind() != reflect.Slice {
		return errors.New("destination must be a pointer to a slice")
	}

	fieldName := q.resolveFieldName(field)
	opts := q.ToOptions()
	db := GetDb(opts...)
	// 确保查询有正确的模型上下文
	if db.Statement.Model == nil && db.Statement.Table == "" {
		db = db.Model(new(T))
	}
	return db.Pluck(fieldName, destination).Error
}

// Preload 预加载关联数据
// 该方法用于预加载关联模型的数据，解决N+1查询问题，在查询主模型时同时加载关联模型数据
//
// 参数:
//   - field: 关联字段名，通常是结构体中定义的关联字段名（如"Roles"、"Profile"等）
//   - args: 可选的关联查询参数，可以是排序、条件等
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 预加载用户的角色信息
//
// query, _ := repo.NewQuery[model.User]()
// query.Preload("Roles")
// var users []model.User
// err := query.Find(&users)
//
// 注意事项:
//   - Preload用于解决N+1查询问题，提高查询效率
//   - field参数应为结构体中定义的关联字段名
//   - 可以链式调用Preload预加载多个关联
//   - 支持嵌套预加载，如"Roles.Permissions"
//   - 可以传递查询条件和排序参数来过滤预加载的数据
func (q *Query[T]) Preload(field string, args ...interface{}) *Query[T] {
	q.opts = append(q.opts, func(db *gorm.DB) *gorm.DB {
		return db.Preload(field, args...)
	})
	return q
}

// GroupBy 实现分组功能
// 该方法用于添加GROUP BY子句到查询中，将查询结果按照指定字段进行分组
//
// 参数:
//   - field: 要分组的字段，可以是字段指针(如 &user.Status)或字段名字符串(如 "status")
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 按用户状态分组
//	query, user := repo.NewQuery[model.User]()
//	query.GroupBy(&user.Status)
//	query.Select(&user.Status) // 明确指定要查询的列
//	var users []model.User
//	err := query.Find(&users)
//	// 生成SQL: SELECT status FROM tb_users GROUP BY status
//
//	// 结合聚合函数使用
//	query, r := repo.NewQuery[model.Role]()
//	query.GroupBy(&r.Status)
//	query.Select(&r.Status)
//	query.Count(&r.ID)
//	var results []struct {
//		Status int `json:"status"`
//		Count  int `json:"count"`
//	}
//	err = query.Scan(&results)
//	// 生成SQL: SELECT status, COUNT(id) as count FROM tb_cmdb_roles GROUP BY status
//
// 注意事项:
//   - GroupBy通常与聚合函数（如Count、Sum、Avg等）一起使用
//   - 可以多次调用GroupBy添加多个分组字段
//   - 如果字段无效，会生成字段名为空的GROUP BY子句（如: GROUP BY ）
func (q *Query[T]) GroupBy(field interface{}) *Query[T] {
	fieldName := q.resolveFieldName(field)
	q.opts = append(q.opts, func(db *gorm.DB) *gorm.DB {
		return db.Group(fieldName)
	})
	return q
}

// Having 添加having条件
// 该方法用于添加HAVING子句到查询中，通常与GroupBy配合使用，用于对分组后的结果进行过滤
//
// 参数:
//   - condition: HAVING条件字符串，可以包含占位符?
//   - args: 条件参数，用于替换condition中的占位符
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 按状态分组并筛选出记录数大于1的组
//	query, role := repo.NewQuery[model.Role]()
//	query.GroupBy(&role.Status)
//	query.Select("status", "COUNT(*) as count")
//	query.Having("COUNT(*) > ?", 1)
//	var results []struct {
//		Status int `json:"status"`
//		Count  int `json:"count"`
//	}
//	err := query.Scan(&results)
//	// 生成SQL: SELECT status, COUNT(id) FROM tb_cmdb_roles GROUP BY status HAVING COUNT(*) > 1
//
//	// 按部门分组并筛选出平均薪资大于5000的部门
//	query, emp := repo.NewQuery[model.Employee]()
//	query.GroupBy(&emp.Department)
//	query.Select("department", "AVG(salary) as avg_salary")
//	query.Having("AVG(salary) > ?", 5000)
//	var deptSalaries []struct {
//		Department string  `json:"department"`
//		AvgSalary  float64 `json:"avg_salary"`
//	}
//	err = query.Scan(&deptSalaries)
//	// 生成SQL: SELECT department, AVG(salary) as avg_salary FROM employees GROUP BY department HAVING AVG(salary) > 5000
//
// 注意事项:
//   - HAVING子句通常与GROUP BY一起使用
//   - condition中的参数会自动转义，防止SQL注入攻击
//   - 可以多次调用Having添加多个条件，它们会以AND逻辑连接
func (q *Query[T]) Having(condition string, args ...interface{}) *Query[T] {
	q.opts = append(q.opts, func(db *gorm.DB) *gorm.DB {
		return db.Having(condition, args...)
	})
	return q
}

// Limit 添加 limit 函数
// 该方法用于限制查询结果的数量，相当于SQL中的LIMIT子句
//
// 参数:
//   - limit: 要限制的结果数量
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 限制查询结果为10条记录
//	query, user := repo.NewQuery[model.User]()
//	query.Limit(10)
//	var users []model.User
//	err := query.Find(&users)
//	// 生成SQL: SELECT * FROM tb_users LIMIT 10
//
//	// 结合Offset实现分页功能
//	query, _ = repo.NewQuery[model.User]()
//	query.Offset(20).Limit(10) // 跳过前20条，取10条记录
//	var pageUsers []model.User
//	err = query.Find(&pageUsers)
//	// 生成SQL: SELECT * FROM tb_users OFFSET 20 LIMIT 10
//
// 注意事项:
//   - Limit值应该为正整数
//   - 通常与Offset配合使用实现分页功能
func (q *Query[T]) Limit(limit int) *Query[T] {
	q.opts = append(q.opts, func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	})
	return q
}

// Offset 添加 offset 函数
// 该方法用于跳过查询结果的前N条记录，相当于SQL中的OFFSET子句
//
// 参数:
//   - offset: 要跳过的记录数量
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 跳过前5条记录
//	query, user := repo.NewQuery[model.User]()
//	query.Offset(5)
//	var users []model.User
//	err := query.Find(&users)
//	// 生成SQL: SELECT * FROM tb_users OFFSET 5
//
//	// 结合Limit实现分页功能
//	query, _ = repo.NewQuery[model.User]()
//	query.Offset(20).Limit(10) // 跳过前20条，取10条记录
//	var pageUsers []model.User
//	err = query.Find(&pageUsers)
//	// 生成SQL: SELECT * FROM tb_users OFFSET 20 LIMIT 10
//
// 注意事项:
//   - Offset值应该为非负整数
//   - 通常与Limit配合使用实现分页功能
func (q *Query[T]) Offset(offset int) *Query[T] {
	q.opts = append(q.opts, func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset)
	})
	return q
}

// Count 添加count聚合函数
// 该方法用于计算记录数量，生成COUNT()聚合函数查询
//
// 参数:
//   - field: 要计算数量的字段，可以是字段指针(如 &user.ID)或字段名字符串(如 "id")
//     如果字段为空或无效，则默认统计所有记录数(*)
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 统计所有用户数量
//	query, user := repo.NewQuery[model.User]()
//	query.Count(&user.ID)
//	var count int64
//	err := query.Scan(&count)
//	if err != nil {
//		// 处理错误
//	}
//	fmt.Printf("用户总数: %d\n", count)
//	// 生成SQL: SELECT COUNT(id) as count FROM tb_users
//
//	// 分组统计不同状态的用户数量
//	query, r := repo.NewQuery[model.Role]()
//	query.GroupBy(&r.Status)
//	query.Select(&r.Status)
//	query.Count(&r.ID)
//	var results []struct {
//		Status int `json:"status"`
//		Count  int `json:"count"`
//	}
//	err = query.Scan(&results)
//	// 生成SQL: SELECT status, COUNT(id) as count FROM tb_cmdb_roles GROUP BY status
//
//	// 统计所有记录数
//	countQuery, _ := repo.NewQuery[model.User]()
//	countQuery.Count(nil) // 或者不传入有效字段
//	var totalCount int64
//	countQuery.Scan(&totalCount)
//	// 生成SQL: SELECT COUNT(*) as count FROM tb_users
//
// 注意事项:
//   - 该方法会构建包含COUNT()的SELECT语句
//   - 当与其他字段一起使用Select时，会将COUNT()作为附加字段添加到SELECT子句中
//   - 如果没有指定有效字段，则默认使用"*"统计所有记录
//   - 扫描结果的目标变量类型应该是数值类型(如int, int64等)
func (q *Query[T]) Count(field interface{}) *Query[T] {
	fieldName := q.resolveFieldName(field)
	if fieldName == "" {
		fieldName = "*"
	}
	q.opts = append(q.opts, func(db *gorm.DB) *gorm.DB {
		// 构建select语句
		var selectClause string

		if len(db.Statement.Selects) > 0 {
			// 如果已有select字段，则追加COUNT字段
			selects := make([]string, len(db.Statement.Selects))
			copy(selects, db.Statement.Selects)
			selects = append(selects, "COUNT("+fieldName+") as count")
			selectClause = strings.Join(selects, ", ")
		} else {
			// 只使用COUNT字段
			selectClause = "COUNT(" + fieldName + ") as count"
		}

		return db.Select(selectClause)
	})

	return q
}

// Sum 添加sum聚合函数
// 该方法用于计算指定字段的总和，生成SUM()聚合函数查询
//
// 参数:
//   - field: 要计算总和的字段，可以是字段指针(如 &user.Score)或字段名字符串(如 "score")
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 计算用户积分总和
//	query, user := repo.NewQuery[model.User]()
//	query.Sum(&user.Score)
//	var totalScore float64
//	err := query.Scan(&totalScore)
//	if err != nil {
//		// 处理错误
//	}
//	fmt.Printf("用户积分总和为: %.0f\n", totalScore)
//	// 生成SQL: SELECT SUM(score) FROM tb_users
//
//	// 计算已完成订单的金额总和
//	orderQuery, order := repo.NewQuery[model.Order]()
//	orderQuery.Eq(&order.Status, "completed")
//	orderQuery.Sum(&order.Amount)
//	var totalAmount float64
//	err = orderQuery.Scan(&totalAmount)
//	fmt.Printf("已完成订单金额总和为: %.2f\n", totalAmount)
//	// 生成SQL: SELECT SUM(amount) FROM tb_orders WHERE status = 'completed'
//
//	// 使用字符串字段名
//	productQuery, _ := repo.NewQuery[model.Product]()
//	productQuery.Sum("price")
//	var totalPrice float64
//	productQuery.Scan(&totalPrice)
//
// 注意事项:
//   - 该方法会覆盖之前的Select设置，只查询SUM()聚合结果
//   - 如果字段无效，会生成字段名为空的条件（如: SELECT SUM() FROM ...）
//   - 扫描结果的目标变量类型应与字段类型兼容
func (q *Query[T]) Sum(field interface{}) *Query[T] {
	fieldName := q.resolveFieldName(field)
	q.opts = append(q.opts, func(db *gorm.DB) *gorm.DB {
		// 构建select语句
		var selectClause string

		if len(db.Statement.Selects) > 0 {
			// 如果已有select字段，则追加SUM字段
			selects := make([]string, len(db.Statement.Selects))
			copy(selects, db.Statement.Selects)
			selects = append(selects, "SUM("+fieldName+")")
			selectClause = strings.Join(selects, ", ")
		} else {
			// 只使用SUM字段
			selectClause = "SUM(" + fieldName + ")"
		}

		return db.Select(selectClause)
	})
	return q
}

// Avg 添加avg聚合函数
// 该方法用于计算指定字段的平均值，生成AVG()聚合函数查询
//
// 参数:
//   - field: 要计算平均值的字段，可以是字段指针(如 &user.Age)或字段名字符串(如 "age")
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 计算用户年龄的平均值
//	query, user := repo.NewQuery[model.User]()
//	query.Avg(&user.Age)
//	var avgAge float64
//	err := query.Scan(&avgAge)
//	if err != nil {
//		// 处理错误
//	}
//	fmt.Printf("用户平均年龄为: %.2f\n", avgAge)
//	// 生成SQL: SELECT AVG(age) FROM tb_users
//
//	// 计算订单金额的平均值
//	orderQuery, order := repo.NewQuery[model.Order]()
//	orderQuery.Eq(&order.Status, "completed")
//	orderQuery.Avg(&order.Amount)
//	var avgAmount float64
//	err = orderQuery.Scan(&avgAmount)
//	fmt.Printf("已完成订单的平均金额为: %.2f\n", avgAmount)
//	// 生成SQL: SELECT AVG(amount) FROM tb_orders WHERE status = 'completed'
//
//	// 使用字符串字段名
//	productQuery, _ := repo.NewQuery[model.Product]()
//	productQuery.Avg("price")
//	var avgPrice float64
//	productQuery.Scan(&avgPrice)
//
// 注意事项:
//   - 该方法会覆盖之前的Select设置，只查询AVG()聚合结果
//   - 如果字段无效，会生成字段名为空的条件（如: SELECT AVG() FROM ...）
//   - 扫描结果的目标变量类型应与字段类型兼容
func (q *Query[T]) Avg(field interface{}) *Query[T] {
	fieldName := q.resolveFieldName(field)
	q.opts = append(q.opts, func(db *gorm.DB) *gorm.DB {
		// 构建select语句
		var selectClause string

		if len(db.Statement.Selects) > 0 {
			// 如果已有select字段，则追加AVG字段
			selects := make([]string, len(db.Statement.Selects))
			copy(selects, db.Statement.Selects)
			selects = append(selects, "AVG("+fieldName+")")
			selectClause = strings.Join(selects, ", ")
		} else {
			// 只使用AVG字段
			selectClause = "AVG(" + fieldName + ")"
		}

		return db.Select(selectClause)
	})
	return q
}

// Max 添加max聚合函数
// 该方法用于计算指定字段的最大值，生成MAX()聚合函数查询
//
// 参数:
//   - field: 要计算最大值的字段，可以是字段指针(如 &user.Age)或字段名字符串(如 "age")
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 计算用户年龄的最大值
//	query, user := repo.NewQuery[model.User]()
//	query.Max(&user.Age)
//	var maxAge float64
//	err := query.Scan(&maxAge)
//	if err != nil {
//		// 处理错误
//	}
//	fmt.Printf("用户最大年龄为: %.0f\n", maxAge)
//	// 生成SQL: SELECT MAX(age) FROM tb_users
//
//	// 计算订单金额的最大值
//	orderQuery, order := repo.NewQuery[model.Order]()
//	orderQuery.Eq(&order.Status, "completed")
//	orderQuery.Max(&order.Amount)
//	var maxAmount float64
//	err = orderQuery.Scan(&maxAmount)
//	fmt.Printf("已完成订单的最大金额为: %.2f\n", maxAmount)
//	// 生成SQL: SELECT MAX(amount) FROM tb_orders WHERE status = 'completed'
//
//	// 使用字符串字段名
//	productQuery, _ := repo.NewQuery[model.Product]()
//	productQuery.Max("price")
//	var maxPrice float64
//	productQuery.Scan(&maxPrice)
//
// 注意事项:
//   - 该方法会覆盖之前的Select设置，只查询MAX()聚合结果
//   - 如果字段无效，会生成字段名为空的条件（如: SELECT MAX() FROM ...）
//   - 扫描结果的目标变量类型应与字段类型兼容
func (q *Query[T]) Max(field interface{}) *Query[T] {
	fieldName := q.resolveFieldName(field)
	q.opts = append(q.opts, func(db *gorm.DB) *gorm.DB {
		// 构建select语句
		var selectClause string

		if len(db.Statement.Selects) > 0 {
			// 如果已有select字段，则追加MAX字段
			selects := make([]string, len(db.Statement.Selects))
			copy(selects, db.Statement.Selects)
			selects = append(selects, "MAX("+fieldName+")")
			selectClause = strings.Join(selects, ", ")
		} else {
			// 只使用MAX字段
			selectClause = "MAX(" + fieldName + ")"
		}

		return db.Select(selectClause)
	})
	return q
}

// Min 添加min聚合函数
// 该方法用于计算指定字段的最小值，生成MIN()聚合函数查询
//
// 参数:
//   - field: 要计算最小值的字段，可以是字段指针(如 &user.Age)或字段名字符串(如 "age")
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 计算用户年龄的最小值
//	query, user := repo.NewQuery[model.User]()
//	query.Min(&user.Age)
//	var minAge float64
//	err := query.Scan(&minAge)
//	if err != nil {
//		// 处理错误
//	}
//	fmt.Printf("用户最小年龄为: %.0f\n", minAge)
//	// 生成SQL: SELECT MIN(age) FROM tb_users
//
//	// 计算订单金额的最小值
//	orderQuery, order := repo.NewQuery[model.Order]()
//	orderQuery.Eq(&order.Status, "completed")
//	orderQuery.Min(&order.Amount)
//	var minAmount float64
//	err = orderQuery.Scan(&minAmount)
//	fmt.Printf("已完成订单的最小金额为: %.2f\n", minAmount)
//	// 生成SQL: SELECT MIN(amount) FROM tb_orders WHERE status = 'completed'
//
//	// 使用字符串字段名
//	productQuery, _ := repo.NewQuery[model.Product]()
//	productQuery.Min("price")
//	var minPrice float64
//	productQuery.Scan(&minPrice)
//
// 注意事项:
//   - 该方法会覆盖之前的Select设置，只查询MIN()聚合结果
//   - 如果字段无效，会生成字段名为空的条件（如: SELECT MIN() FROM ...）
//   - 扫描结果的目标变量类型应与字段类型兼容
func (q *Query[T]) Min(field interface{}) *Query[T] {
	fieldName := q.resolveFieldName(field)
	q.opts = append(q.opts, func(db *gorm.DB) *gorm.DB {
		// 构建select语句
		var selectClause string

		if len(db.Statement.Selects) > 0 {
			// 如果已有select字段，则追加MIN字段
			selects := make([]string, len(db.Statement.Selects))
			copy(selects, db.Statement.Selects)
			selects = append(selects, "MIN("+fieldName+")")
			selectClause = strings.Join(selects, ", ")
		} else {
			// 只使用MIN字段
			selectClause = "MIN(" + fieldName + ")"
		}

		return db.Select(selectClause)
	})
	return q
}

// Join 添加join条件
func (q *Query[T]) Join(query string, args ...interface{}) *Query[T] {
	q.opts = append(q.opts, func(db *gorm.DB) *gorm.DB {
		return db.Joins(query, args...)
	})
	return q
}

// LeftJoin 添加left join条件
func (q *Query[T]) LeftJoin(query string, args ...interface{}) *Query[T] {
	q.opts = append(q.opts, func(db *gorm.DB) *gorm.DB {
		return db.Joins("LEFT "+query, args...)
	})
	return q
}

// SubQueryEq 子查询等于条件
// 该方法用于添加等于子查询的条件，将字段与另一个查询结果进行比较，判断字段值是否等于子查询结果
//
// 参数:
//   - field: 要比较的字段，可以是字段指针(如 &user.DepartmentID)或字段名字符串(如 "department_id")
//   - subDB: 子查询的*gorm.DB对象，定义了要执行的子查询
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 查询属于特定部门ID的用户
//	subDB := global.DBMysql.Model(&model.Department{}).Select("id").Where("name = ?", "IT部")
//	query, user := repo.NewQuery[model.User]()
//	query.SubQueryEq(&user.DepartmentID, subDB)
//	var users []model.User
//	err := query.Find(&users)
//	// 生成SQL: SELECT * FROM tb_users WHERE department_id = (SELECT id FROM tb_departments WHERE name = 'IT部')
//
// 注意事项:
//   - 子查询(subDB)应该返回单个值，以便与主查询字段进行等于比较
//   - 如果字段无效，会生成字段名为空的条件（如: "" = (...)）
//   - 子查询返回多个值时，数据库会报错
func (q *Query[T]) SubQueryEq(field interface{}, subDB *gorm.DB) *Query[T] {
	fieldName := q.resolveFieldName(field)
	q.opts = append(q.opts, func(db *gorm.DB) *gorm.DB {
		return db.Where(fieldName+" = (?)", subDB)
	})
	return q
}

// SubQueryIn 子查询in条件
// 该方法用于添加IN子查询条件，将字段与另一个查询结果进行比较，判断字段值是否在子查询结果中
//
// 参数:
//   - field: 要比较的字段，可以是字段指针(如 &user.DepartmentID)或字段名字符串(如 "department_id")
//   - subDB: 子查询的*gorm.DB对象，定义了要执行的子查询
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 查询属于特定用户名的帖子(用户和帖子的关系为一对多)
//	subDB := global.DBMysql.Model(&model.User{}).Select("id").Where("username = ?", "vben")
//	query, p := repo.NewQuery[model.Post]()
//	query.SubQueryIn(&p.UserID, subDB)
//	var posts []model.Post
//	err := query.Find(&posts)
//	// 生成SQL: SELECT * FROM tb_posts WHERE user_id IN (SELECT id FROM tb_users WHERE username = 'vben')
//
//	// 查询属于多个部门的用户(部门和用户的关系为一对多)
//	deptSubQuery := global.DBMysql.Model(&model.Department{}).Select("id").Where("active = ?", true)
//	userQuery, u := repo.NewQuery[model.User]()
//	userQuery.SubQueryIn(&u.DepartmentID, deptSubQuery)
//	var users []model.User
//	err = userQuery.Find(&users)
//	// 生成SQL: SELECT * FROM tb_users WHERE department_id IN (SELECT id FROM tb_departments WHERE active = true)
//
// // 查询属于特定角色的用户(用户和角色关系为多对多)
//
//	subDB := global.DBMysql.Model(&model.RoleUser{}).
//	  Select("user_id").
//	  Where("role_id IN (SELECT id FROM tb_cmdb_roles WHERE name = ?)", "Admin")
//	userQuery, u := repo.NewQuery[model.User]()
//	userQuery.SubQueryIn(&u.ID, subDB)
//	var users []model.User
//	err = userQuery.Find(&users)
//	// 生成SQL: SELECT * FROM tb_users WHERE id IN (SELECT user_id FROM tb_cmdb_role_users WHERE role_id IN (SELECT id FROM tb_cmdb_roles WHERE name = 'Admin'))
//
// 注意事项:
//   - 子查询(subDB)应该返回单列数据，以便与主查询字段进行比较
//   - 如果字段无效，会生成字段名为空的条件（如: "" IN (...)）
func (q *Query[T]) SubQueryIn(field interface{}, subDB *gorm.DB) *Query[T] {
	fieldName := q.resolveFieldName(field)
	q.opts = append(q.opts, func(db *gorm.DB) *gorm.DB {
		return db.Where(fieldName+" IN (?)", subDB)
	})
	return q
}

//// SubQueryIn 子查询in条件
//func (q *Query[T]) SubQueryIn(field interface{}, subQuery interface{ ToOptions() []DBOption }) *Query[T] {
//    fieldName := q.resolveFieldName(field)
//    subOpts := subQuery.ToOptions()
//    q.opts = append(q.opts, func(db *gorm.DB) *gorm.DB {
//        subDB := GetDb(subOpts...)
//		// 确保子查询有正确的模型上下文
//        // 由于无法通过接口获取具体类型，需要在调用时确保子查询已正确设置模型
//        return db.Where(fieldName+" IN (?)", subDB)
//    })
//    return q
//}

// Or 添加OR条件
// 该方法用于添加OR逻辑条件到查询中，支持自定义查询条件和参数
//
// 参数:
//   - query: 查询条件，可以是字符串形式的SQL条件或其它GORM支持的查询格式
//   - args: 查询参数，用于替换查询条件中的占位符
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 基本OR条件查询
//	query, user := repo.NewQuery[model.User]()
//	query.Eq(&user.Status, 1)
//	query.Or("age > ? AND age < ?", 18, 65)
//	var users []model.User
//	err := query.Find(&users)
//	// 生成SQL: SELECT * FROM tb_users WHERE status = 1 OR (age > 18 AND age < 65)
//
//	// 复杂OR条件
//	query.Or("department = ? OR salary > ?", "IT", 50000)
//
// 注意事项:
//   - SQL参数会自动转义，防止SQL注入攻击
//   - OR条件会与之前添加的条件形成逻辑或关系
func (q *Query[T]) Or(query interface{}, args ...interface{}) *Query[T] {
	q.opts = append(q.opts, func(db *gorm.DB) *gorm.DB {
		return db.Or(query, args...)
	})
	return q
}

// Not 添加NOT条件
// 该方法用于添加NOT逻辑条件到查询中，支持自定义查询条件和参数
//
// 参数:
//   - query: 查询条件，可以是字符串形式的SQL条件、结构体或map等GORM支持的查询格式
//   - args: 查询参数，用于替换查询条件中的占位符
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 基本NOT条件查询
//	query, user := repo.NewQuery[model.User]()
//	query.Not("status = ?", 0)
//	var users []model.User
//	err := query.Find(&users)
//	// 生成SQL: SELECT * FROM tb_users WHERE NOT (status = 0)
//
//	// 使用结构体条件
//	query, _ = repo.NewQuery[model.User]()
//	query.Not(model.User{Status: 0})
//	var activeUsers []model.User
//	err = query.Find(&activeUsers)
//	// 生成SQL: SELECT * FROM tb_users WHERE NOT (status = 0)
//
//	// 复杂NOT条件
//	query, _ = repo.NewQuery[model.User]()
//	query.Not("age < ? OR department = ?", 18, "IT")
//	var adultUsers []model.User
//	err = query.Find(&adultUsers)
//	// 生成SQL: SELECT * FROM tb_users WHERE NOT (age < 18 OR department = 'IT')
//
// 注意事项:
//   - SQL参数会自动转义，防止SQL注入攻击
//   - NOT条件会与之前添加的条件形成逻辑非关系
//   - 可以多次调用Not添加多个NOT条件，它们会以AND逻辑连接
func (q *Query[T]) Not(query interface{}, args ...interface{}) *Query[T] {
	q.opts = append(q.opts, func(db *gorm.DB) *gorm.DB {
		return db.Not(query, args...)
	})
	return q
}

// addSubQueryCondition 添加子查询条件
func (q *Query[T]) addSubQueryCondition(field interface{}, operator string, sql string, args ...interface{}) *Query[T] {
	if sql == "" {
		log.Printf("Empty SQL in subquery condition")
		return q
	}
	fieldName := q.resolveFieldName(field)
	if fieldName == "" {
		log.Printf("Invalid field in subquery condition: " + fmt.Sprintf("%v", field))
	}

	safeField := strings.TrimSpace(fieldName)

	// normalize and whitelist operator
	op := strings.ToUpper(strings.TrimSpace(operator))
	allowedOperators := map[string]bool{
		"IN": true, "NOT IN": true,
		">": true, ">=": true,
		"<": true, "<=": true,
		"=": true, "!=": true, "<>": true,
	}

	if !allowedOperators[op] {
		log.Printf("invalid operator in subquery condition: " + fmt.Sprintf("%v", op))
		return q
	}

	normalizedSQL := strings.TrimSpace(sql)
	q.opts = append(q.opts, func(db *gorm.DB) *gorm.DB {
		return db.Where(safeField+" "+op+" (?)", gorm.Expr(normalizedSQL, args...))
	})
	return q
}

// InSql 添加IN SQL条件
// 该方法用于添加IN子查询条件，将字段与一个SQL子查询结果进行比较，判断字段值是否在子查询结果中
//
// 参数:
//   - field: 要比较的字段，可以是字段指针(如 &user.DepartmentID)或字段名字符串(如 "department_id")
//   - sql: 子查询SQL语句，可以包含占位符 ?
//   - args: SQL查询参数，用于替换SQL中的占位符
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 查询属于特定用户ID的文章(用户和文章关系为一对多)
//	subQuery := "SELECT id FROM tb_users WHERE username = ?"
//	query, p := repo.NewQuery[model.Post]()
//	query.InSql(&p.UserID, subQuery, "vben")
//	var posts []model.Post
//	err := query.Find(&posts)
//	// 生成SQL: SELECT * FROM tb_posts WHERE user_id IN (SELECT id FROM tb_users WHERE username = 'vben')
//
//	// 使用字符串字段名
//	query.InSql("user_id", subQuery, "vben")
//
// //查询属于特定角色的用户(用户和角色关系为多对多)
//
//	subQuery := "SELECT user_id FROM tb_cmdb_role_users WHERE role_id IN (SELECT id FROM tb_cmdb_roles WHERE name = ?)"
//	query, user := repo.NewQuery[model.User]()
//	query.InSql(&user.ID, subQuery, "Admin")
//	var users []model.User
//	err := query.Find(&users)
//	// 生成SQL: SELECT * FROM tb_users WHERE id IN (SELECT user_id FROM tb_cmdb_role_users WHERE role_id IN (SELECT id FROM tb_cmdb_roles WHERE name = 'Admin'))
//
// 注意事项:
//   - SQL参数会自动转义，防止SQL注入攻击
//   - 如果SQL语句为空字符串，则该条件会被忽略
//   - 如果字段无效，会生成字段名为空的条件（如: "" IN (SELECT ...)）
func (q *Query[T]) InSql(field interface{}, sql string, args ...interface{}) *Query[T] {
	return q.addSubQueryCondition(field, "IN", sql, args...)
}

// NotInSql 添加NOT IN SQL条件
// 该方法用于添加NOT IN子查询条件，将字段与一个SQL子查询结果进行比较，判断字段值是否不在子查询结果中
//
// 参数:
//   - field: 要比较的字段，可以是字段指针(如 &user.DepartmentID)或字段名字符串(如 "department_id")
//   - sql: 子查询SQL语句，可以包含占位符 ?
//   - args: SQL查询参数，用于替换SQL中的占位符
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 查询不属于特定部门ID列表的用户
//	subQuery := "SELECT id FROM tb_departments WHERE company_id = ?"
//	query, user := repo.NewQuery[model.User]()
//	query.NotInSql(&user.DepartmentID, subQuery, 1001)
//	var users []model.User
//	err := query.Find(&users)
//	// 生成SQL: SELECT * FROM tb_users WHERE department_id NOT IN (SELECT id FROM tb_departments WHERE company_id = 1001)
//
//	// 使用字符串字段名
//	query.NotInSql("department_id", subQuery, 1001)
//
// 注意事项:
//   - SQL参数会自动转义，防止SQL注入攻击
//   - 如果SQL语句为空字符串，则该条件会被忽略
//   - 如果字段无效，会生成字段名为空的条件（如: "" NOT IN (SELECT ...)）
func (q *Query[T]) NotInSql(field interface{}, sql string, args ...interface{}) *Query[T] {
	return q.addSubQueryCondition(field, "NOT IN", sql, args...)
}

// GtSql 添加> SQL条件
// 该方法用于添加大于的SQL子查询条件，将字段与一个SQL子查询进行比较
//
// 参数:
//   - field: 要比较的字段，可以是字段指针(如 &user.Age)或字段名字符串(如 "age")
//   - sql: 子查询SQL语句，可以包含占位符 ?
//   - args: SQL查询参数，用于替换SQL中的占位符
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 查询年龄大于某个子查询结果的用户
//	subQuery := "SELECT AVG(age) FROM tb_users WHERE department = ?"
//	query, user := repo.NewQuery[model.User]()
//	query.GtSql(&user.Age, subQuery, "IT")
//	var users []model.User
//	err := query.Find(&users)
//	// 生成SQL: SELECT * FROM tb_users WHERE age > (SELECT AVG(age) FROM tb_users WHERE department = 'IT')
//
//	// 使用字符串字段名
//	query.GtSql("age", subQuery, "IT")
//
// 注意事项:
//   - SQL参数会自动转义，防止SQL注入攻击
//   - 如果SQL语句为空字符串，则该条件会被忽略
//   - 如果字段无效，会生成字段名为空的条件（如: "" > (SELECT ...)）
func (q *Query[T]) GtSql(field interface{}, sql string, args ...interface{}) *Query[T] {
	return q.addSubQueryCondition(field, ">", sql, args...)
}

// GeSql 添加>= SQL条件
// 该方法用于添加大于等于的SQL子查询条件，将字段与一个SQL子查询进行比较
//
// 参数:
//   - field: 要比较的字段，可以是字段指针(如 &user.Age)或字段名字符串(如 "age")
//   - sql: 子查询SQL语句，可以包含占位符 ?
//   - args: SQL查询参数，用于替换SQL中的占位符
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 查询年龄大于等于某个子查询结果的用户
//	subQuery := "SELECT AVG(age) FROM tb_users WHERE department = ?"
//	query, user := repo.NewQuery[model.User]()
//	query.GeSql(&user.Age, subQuery, "IT")
//	var users []model.User
//	err := query.Find(&users)
//	// 生成SQL: SELECT * FROM tb_users WHERE age >= (SELECT AVG(age) FROM tb_users WHERE department = 'IT')
//
//	// 使用字符串字段名
//	query.GeSql("age", subQuery, "IT")
//
// 注意事项:
//   - SQL参数会自动转义，防止SQL注入攻击
//   - 如果SQL语句为空字符串，则该条件会被忽略
//   - 如果字段无效，会生成字段名为空的条件（如: "" >= (SELECT ...)）
func (q *Query[T]) GeSql(field interface{}, sql string, args ...interface{}) *Query[T] {
	return q.addSubQueryCondition(field, ">=", sql, args...)
}

// LtSql 添加< SQL条件
// 该方法用于添加小于的SQL子查询条件，将字段与一个SQL子查询进行比较
//
// 参数:
//   - field: 要比较的字段，可以是字段指针(如 &user.Age)或字段名字符串(如 "age")
//   - sql: 子查询SQL语句，可以包含占位符 ?
//   - args: SQL查询参数，用于替换SQL中的占位符
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 查询年龄小于某个子查询结果的用户
//	subQuery := "SELECT AVG(age) FROM tb_users WHERE department = ?"
//	query, user := repo.NewQuery[model.User]()
//	query.LtSql(&user.Age, subQuery, "IT")
//	var users []model.User
//	err := query.Find(&users)
//	// 生成SQL: SELECT * FROM tb_users WHERE age < (SELECT AVG(age) FROM tb_users WHERE department = 'IT')
//
//	// 使用字符串字段名
//	query.LtSql("age", subQuery, "IT")
//
// 注意事项:
//   - SQL参数会自动转义，防止SQL注入攻击
//   - 如果SQL语句为空字符串，则该条件会被忽略
//   - 如果字段无效，会生成字段名为空的条件（如: "" < (SELECT ...)）
func (q *Query[T]) LtSql(field interface{}, sql string, args ...interface{}) *Query[T] {
	return q.addSubQueryCondition(field, "<", sql, args...)
}

// LeSql 添加<= SQL条件
// 该方法用于添加小于等于的SQL子查询条件，将字段与一个SQL子查询进行比较
//
// 参数:
//   - field: 要比较的字段，可以是字段指针(如 &user.Age)或字段名字符串(如 "age")
//   - sql: 子查询SQL语句，可以包含占位符 ?
//   - args: SQL查询参数，用于替换SQL中的占位符
//
// 返回:
//   - *Query[T]: 返回查询构建器本身，支持链式调用
//
// 使用示例:
//
//	// 查询年龄小于等于某个子查询结果的用户
//	subQuery := "SELECT AVG(age) FROM tb_users WHERE department = ?"
//	query, user := repo.NewQuery[model.User]()
//	query.LeSql(&user.Age, subQuery, "IT")
//	var users []model.User
//	err := query.Find(&users)
//	// 生成SQL: SELECT * FROM tb_users WHERE age <= (SELECT AVG(age) FROM tb_users WHERE department = 'IT')
//
//	// 使用字符串字段名
//	query.LeSql("age", subQuery, "IT")
//
// 注意事项:
//   - SQL参数会自动转义，防止SQL注入攻击
//   - 如果SQL语句为空字符串，则该条件会被忽略
//   - 如果字段无效，会生成字段名为空的条件（如: "" <= (SELECT ...)）
func (q *Query[T]) LeSql(field interface{}, sql string, args ...interface{}) *Query[T] {
	return q.addSubQueryCondition(field, "<=", sql, args...)
}
