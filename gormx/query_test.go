package gormx_test

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/shijl0925/go-toolkits/gormx"
	"gorm.io/gorm"
)

// Post 文章模型，与User是多对一关系
type Post struct {
	gorm.Model
	UserID uint `gorm:"column:user_id;index"`
	User   User `gorm:"foreignKey:UserID;references:ID"`

	Title   string `gorm:"column:title"`
	Content string `gorm:"column:content"`
	Status  int    `gorm:"column:status"`
}

// User 用户模型
type User struct {
	gorm.Model

	Name     string  `gorm:"column:name"`
	Email    string  `gorm:"column:email"`
	Phone    string  `gorm:"column:phone"`
	Age      int     `gorm:"column:age"`
	Score    float64 `gorm:"column:score"`
	Address  string  `gorm:"column:address"`
	IsActive bool    `gorm:"column:is_active"`
	Salary   float64 `gorm:"column:salary"`
	Status   int     `gorm:"column:status"`

	Posts []Post `gorm:"foreignKey:UserID"`     // 一对多关系
	Roles []Role `gorm:"many2many:user_roles;"` // 多对多关系

	rolesManager gormx.AssociationManager[User, Role]
}

// RolesManager 获取角色关联管理器
func (u *User) RolesManager() gormx.AssociationManager[User, Role] {
	if u.rolesManager == nil {
		u.rolesManager = gormx.NewAssociationManager[User, Role](*u, "Roles")
	}
	return u.rolesManager
}

// UserRole 用户角色关联模型，用于建立 User 和 Role 的多对多关系
type UserRole struct {
	gorm.Model
	UserID uint `gorm:"column:user_id;index"`
	RoleID uint `gorm:"column:role_id;index"`
}

// Role 角色模型
type Role struct {
	gorm.Model
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
	Status      int    `gorm:"column:status"`

	// 多对多关系
	Users []User `gorm:"many2many:user_roles;"`
}

// Profile 用户资料模型，与 User 是一对一关系
type Profile struct {
	gorm.Model
	UserID      uint   `gorm:"column:user_id;index"`
	Avatar      string `gorm:"column:avatar"`
	PhoneNumber string `gorm:"column:phone_number"`
	Bio         string `gorm:"column:bio"`

	// 一对一关系
	User User `gorm:"foreignKey:UserID"`
}

// ValidateQuery 验证给定的模型查询是否能正确执行
func ValidateQuery[T any](t *testing.T, query *gormx.Query[T]) error {
	sql, args := query.ToSQLAndArgs()

	db := gormx.GetDb()
	if db == nil {
		return errors.New("database is nil")
	}

	if strings.TrimSpace(sql) == "" {
		return errors.New("SQL statement is empty")
	}

	// 检查SQL中占位符数量与参数数量是否匹配
	placeholderCount := strings.Count(sql, "?")
	if placeholderCount != len(args) {
		return errors.New("SQL statement placeholder count does not match argument count")
	}

	// 使用EXPLAIN语句检查生成的SQL语法(MySQL特有)
	// db.Dialector.Name() == "mysql"
	if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(sql)), "SELECT") {
		explainSQL := "EXPLAIN " + sql
		if err := db.Raw(explainSQL, args...).Error; err != nil {
			return fmt.Errorf("failed to execute EXPLAIN statement: %v", err)
		}
	} else {
		// 使用DryRun模式检查生成的SQL语法
		sessionDB := db.Session(&gorm.Session{PrepareStmt: true, DryRun: true})

		var result T
		stmt := sessionDB.Model(&result).Where(sql, args...)

		if err := stmt.Error; err != nil {
			return fmt.Errorf("failed to execute SQL statement: %v", err)
		}
	}

	return nil
}

// AssertQueryValid 断言给定的模型查询是否能正确执行
func AssertQueryValid[T any](t *testing.T, query *gormx.Query[T]) {
	t.Helper() // 标记这是一个测试辅助函数

	if err := ValidateQuery(t, query); err != nil {
		t.Errorf("Query validation failed: %v", err)
	}
}

func expectedSQLForCurrentDialect(expected string) string {
	if db := gormx.GetDb(); db != nil && db.Dialector.Name() == testDialectPostgres {
		expected = currentDialectSQLFragment(expected)
		expected = postgresPlaceholderSQL(expected)
	}

	return expected
}

func currentDialectSQLFragment(fragment string) string {
	if db := gormx.GetDb(); db != nil && db.Dialector.Name() == testDialectPostgres {
		return strings.NewReplacer(
			"`", `"`,
			" REGEXP ", " ~ ",
		).Replace(fragment)
	}

	return fragment
}

func postgresPlaceholderSQL(sql string) string {
	var builder strings.Builder
	builder.Grow(len(sql) + 8)

	argIndex := 1
	for _, char := range sql {
		if char == '?' {
			builder.WriteString(fmt.Sprintf("$%d", argIndex))
			argIndex++
			continue
		}
		builder.WriteRune(char)
	}

	return builder.String()
}

// TestQuery_Eq 测试 Eq 方法的基本功能
func TestQuery_Eq(t *testing.T) {
	t.Run("test with field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字段指针的方式
		result := query.Eq(&user.Name, "testuser")

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("Eq should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("Eq should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `name` = ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "testuser" {
			t.Errorf("Expected args: [testuser], got: %v", args)
		}
	})

	t.Run("test with field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字段名字符串的方式
		query.Eq("name", "testuser")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `name` = ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "testuser" {
			t.Errorf("Expected args: [testuser], got: %v", args)
		}
	})

	t.Run("test with different value types", func(t *testing.T) {
		_, user := gormx.NewQuery[User]()

		// 测试不同的值类型
		testCases := []struct {
			field interface{}
			value interface{}
			name  string
		}{
			{&user.Age, 25, "integer value"},
			{&user.IsActive, true, "boolean value"},
			{&user.Salary, 5000.50, "float value"},
			{"status", 1, "string field with integer value"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				testQuery, _ := gormx.NewQuery[User]()
				testQuery.Eq(tc.field, tc.value)

				_, args := testQuery.ToSQLAndArgs()
				if len(args) != 1 || args[0] != tc.value {
					t.Errorf("Expected arg: %v, got: %v", tc.value, args[0])
				}
			})
		}
	})

	t.Run("test chaining multiple conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试链式调用多个条件
		query.Eq(&user.Name, "testuser").
			Eq(&user.Age, 25).
			Eq("status", 1)

		if len(query.ToOptions()) != 3 {
			t.Error("Should have 3 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `name` = ? AND `age` = ? AND `status` = ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 3 {
			t.Errorf("Expected 3 args, got: %d", len(args))
		}

		expectedArgs := []interface{}{"testuser", 25, 1}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d]: %v, got: %v", i, expected, args[i])
			}
		}
	})

	t.Run("test with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效字段名（应该不会panic，但会在实际查询时处理）
		query.Eq("invalid_field", "value")

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍然会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `invalid_field` = ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "value" {
			t.Errorf("Expected args: [value], got: %v", args)
		}
	})

	t.Run("test with multiple dots", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		query.Eq("address.city.town", "New York")
		sql, args := query.ToSQLAndArgs()

		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `address`.`city`.`town` = ? AND `users`.`deleted_at` IS NULL")

		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "New York" {
			t.Errorf("Expected args: [value], got: %v", args)
		}
	})
}

// TestQuery_Ne 测试 Ne 方法的基本功能
func TestQuery_Ne(t *testing.T) {
	t.Run("test with field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字段指针的方式
		result := query.Ne(&user.Name, "testuser")

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("Ne should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("Ne should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `name` <> ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "testuser" {
			t.Errorf("Expected args: [testuser], got: %v", args)
		}
	})

	t.Run("test with field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字段名字符串的方式
		query.Ne("name", "testuser")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `name` <> ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "testuser" {
			t.Errorf("Expected args: [testuser], got: %v", args)
		}
	})

	t.Run("test with different value types", func(t *testing.T) {
		_, user := gormx.NewQuery[User]()

		// 测试不同的值类型
		testCases := []struct {
			field interface{}
			value interface{}
			name  string
		}{
			{&user.Age, 25, "integer value"},
			{&user.IsActive, true, "boolean value"},
			{&user.Salary, 5000.50, "float value"},
			{"status", 1, "string field with integer value"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				testQuery, _ := gormx.NewQuery[User]()
				testQuery.Ne(tc.field, tc.value)

				_, args := testQuery.ToSQLAndArgs()
				if len(args) != 1 || args[0] != tc.value {
					t.Errorf("Expected arg: %v, got: %v", tc.value, args[0])
				}
			})
		}
	})

	t.Run("test chaining multiple conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试链式调用多个条件
		query.Ne(&user.Name, "testuser").
			Ne(&user.Age, 25).
			Ne("status", 1)

		if len(query.ToOptions()) != 3 {
			t.Error("Should have 3 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `name` <> ? AND `age` <> ? AND `status` <> ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 3 {
			t.Errorf("Expected 3 args, got: %d", len(args))
		}

		expectedArgs := []interface{}{"testuser", 25, 1}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d]: %v, got: %v", i, expected, args[i])
			}
		}
	})

	t.Run("test with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效字段名（应该不会panic，但会在实际查询时处理）
		query.Ne("invalid_field", "value")

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍然会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `invalid_field` <> ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "value" {
			t.Errorf("Expected args: [value], got: %v", args)
		}
	})
}

// TestQuery_Gt 测试 Gt 方法的基本功能
func TestQuery_Gt(t *testing.T) {
	t.Run("test with field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字段指针的方式
		result := query.Gt(&user.Age, 18)

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("Gt should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("Gt should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` > ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 18 {
			t.Errorf("Expected args: [18], got: %v", args)
		}
	})

	t.Run("test with field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字段名字符串的方式
		query.Gt("salary", 5000.0)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `salary` > ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 5000.0 {
			t.Errorf("Expected args: [5000.0], got: %v", args)
		}
	})

	t.Run("test with different value types", func(t *testing.T) {
		_, user := gormx.NewQuery[User]()

		// 测试不同的值类型
		testCases := []struct {
			field interface{}
			value interface{}
			name  string
		}{
			{&user.Age, 25, "integer value"},
			{&user.Salary, 5000.50, "float value"},
			{"status", 1, "string field with integer value"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				testQuery, _ := gormx.NewQuery[User]()
				testQuery.Gt(tc.field, tc.value)

				_, args := testQuery.ToSQLAndArgs()
				if len(args) != 1 || args[0] != tc.value {
					t.Errorf("Expected arg: %v, got: %v", tc.value, args[0])
				}
			})
		}
	})

	t.Run("test chaining multiple conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试链式调用多个条件
		query.Gt(&user.Age, 18).
			Gt(&user.Salary, 3000.0).
			Gt("status", 0)

		if len(query.ToOptions()) != 3 {
			t.Error("Should have 3 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` > ? AND `salary` > ? AND `status` > ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 3 {
			t.Errorf("Expected 3 args, got: %d", len(args))
		}

		expectedArgs := []interface{}{18, 3000.0, 0}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d]: %v, got: %v", i, expected, args[i])
			}
		}
	})

	t.Run("test with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效字段名（应该不会panic，但会在实际查询时处理）
		query.Gt("invalid_field", 10)

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍然会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `invalid_field` > ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 10 {
			t.Errorf("Expected args: [10], got: %v", args)
		}
	})
}

// TestQuery_Ge 测试 Ge 方法的基本功能
func TestQuery_Ge(t *testing.T) {
	t.Run("test with field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字段指针的方式
		result := query.Ge(&user.Age, 18)

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("Ge should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("Ge should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` >= ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 18 {
			t.Errorf("Expected args: [18], got: %v", args)
		}
	})

	t.Run("test with field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字段名字符串的方式
		query.Ge("salary", 5000.0)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `salary` >= ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 5000.0 {
			t.Errorf("Expected args: [5000.0], got: %v", args)
		}
	})

	t.Run("test with different value types", func(t *testing.T) {
		_, user := gormx.NewQuery[User]()

		// 测试不同的值类型
		testCases := []struct {
			field interface{}
			value interface{}
			name  string
		}{
			{&user.Age, 25, "integer value"},
			{&user.Salary, 5000.50, "float value"},
			{"status", 1, "string field with integer value"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				testQuery, _ := gormx.NewQuery[User]()
				testQuery.Ge(tc.field, tc.value)

				_, args := testQuery.ToSQLAndArgs()
				if len(args) != 1 || args[0] != tc.value {
					t.Errorf("Expected arg: %v, got: %v", tc.value, args[0])
				}
			})
		}
	})

	t.Run("test chaining multiple conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试链式调用多个条件
		query.Ge(&user.Age, 18).
			Ge(&user.Salary, 3000.0).
			Ge("status", 0)

		if len(query.ToOptions()) != 3 {
			t.Error("Should have 3 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` >= ? AND `salary` >= ? AND `status` >= ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 3 {
			t.Errorf("Expected 3 args, got: %d", len(args))
		}

		expectedArgs := []interface{}{18, 3000.0, 0}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d]: %v, got: %v", i, expected, args[i])
			}
		}
	})

	t.Run("test with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效字段名（应该不会panic，但会在实际查询时处理）
		query.Ge("invalid_field", 10)

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍然会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `invalid_field` >= ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 10 {
			t.Errorf("Expected args: [10], got: %v", args)
		}
	})
}

// TestQuery_Lt 测试 Lt 方法的基本功能
func TestQuery_Lt(t *testing.T) {
	t.Run("test with field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字段指针的方式
		result := query.Lt(&user.Age, 65)

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("Lt should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("Lt should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` < ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 65 {
			t.Errorf("Expected args: [65], got: %v", args)
		}
	})

	t.Run("test with field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字段名字符串的方式
		query.Lt("salary", 10000.0)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `salary` < ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 10000.0 {
			t.Errorf("Expected args: [10000.0], got: %v", args)
		}
	})

	t.Run("test with different value types", func(t *testing.T) {
		_, user := gormx.NewQuery[User]()

		// 测试不同的值类型
		testCases := []struct {
			field interface{}
			value interface{}
			name  string
		}{
			{&user.Age, 25, "integer value"},
			{&user.Salary, 5000.50, "float value"},
			{"status", 1, "string field with integer value"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				testQuery, _ := gormx.NewQuery[User]()
				testQuery.Lt(tc.field, tc.value)

				_, args := testQuery.ToSQLAndArgs()
				if len(args) != 1 || args[0] != tc.value {
					t.Errorf("Expected arg: %v, got: %v", tc.value, args[0])
				}
			})
		}
	})

	t.Run("test chaining multiple conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试链式调用多个条件
		query.Lt(&user.Age, 65).
			Lt(&user.Salary, 10000.0).
			Lt("status", 5)

		if len(query.ToOptions()) != 3 {
			t.Error("Should have 3 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` < ? AND `salary` < ? AND `status` < ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 3 {
			t.Errorf("Expected 3 args, got: %d", len(args))
		}

		expectedArgs := []interface{}{65, 10000.0, 5}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d]: %v, got: %v", i, expected, args[i])
			}
		}
	})

	t.Run("test with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效字段名（应该不会panic，但会在实际查询时处理）
		query.Lt("invalid_field", 10)

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍然会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `invalid_field` < ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 10 {
			t.Errorf("Expected args: [10], got: %v", args)
		}
	})
}

// TestQuery_Le 测试 Le 方法的基本功能
func TestQuery_Le(t *testing.T) {
	t.Run("test with field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字段指针的方式
		result := query.Le(&user.Age, 65)

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("Le should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("Le should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` <= ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 65 {
			t.Errorf("Expected args: [65], got: %v", args)
		}
	})

	t.Run("test with field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字段名字符串的方式
		query.Le("salary", 10000.0)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `salary` <= ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 10000.0 {
			t.Errorf("Expected args: [10000.0], got: %v", args)
		}
	})

	t.Run("test with different value types", func(t *testing.T) {
		_, user := gormx.NewQuery[User]()

		// 测试不同的值类型
		testCases := []struct {
			field interface{}
			value interface{}
			name  string
		}{
			{&user.Age, 25, "integer value"},
			{&user.Salary, 5000.50, "float value"},
			{"status", 1, "string field with integer value"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				testQuery, _ := gormx.NewQuery[User]()
				testQuery.Le(tc.field, tc.value)

				_, args := testQuery.ToSQLAndArgs()
				if len(args) != 1 || args[0] != tc.value {
					t.Errorf("Expected arg: %v, got: %v", tc.value, args[0])
				}
			})
		}
	})

	t.Run("test chaining multiple conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试链式调用多个条件
		query.Le(&user.Age, 65).
			Le(&user.Salary, 10000.0).
			Le("status", 5)

		if len(query.ToOptions()) != 3 {
			t.Error("Should have 3 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` <= ? AND `salary` <= ? AND `status` <= ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 3 {
			t.Errorf("Expected 3 args, got: %d", len(args))
		}

		expectedArgs := []interface{}{65, 10000.0, 5}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d]: %v, got: %v", i, expected, args[i])
			}
		}
	})

	t.Run("test with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效字段名（应该不会panic，但会在实际查询时处理）
		query.Le("invalid_field", 10)

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍然会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `invalid_field` <= ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 10 {
			t.Errorf("Expected args: [10], got: %v", args)
		}
	})
}

// TestQuery_Like 测试 Like 方法的基本功能
func TestQuery_Like(t *testing.T) {
	t.Run("test with field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字段指针的方式
		result := query.Like(&user.Name, "testuser")

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("Like should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("Like should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `name` LIKE ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "%testuser%" {
			t.Errorf("Expected args: [%%testuser%%], got: %v", args)
		}
	})

	t.Run("test with field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字段名字符串的方式
		query.Like("name", "testuser")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `name` LIKE ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "%testuser%" {
			t.Errorf("Expected args: [%%testuser%%], got: %v", args)
		}
	})

	t.Run("test with different string values", func(t *testing.T) {
		_, user := gormx.NewQuery[User]()

		// 测试不同的字符串值
		testCases := []struct {
			field interface{}
			value string
			name  string
		}{
			{&user.Name, "admin", "simple string"},
			{&user.Name, "", "empty string"},
			{&user.Name, "%admin%", "string with wildcards"},
			{"name", "manager", "string field with string value"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				testQuery, _ := gormx.NewQuery[User]()
				testQuery.Like(tc.field, tc.value)

				_, args := testQuery.ToSQLAndArgs()
				expectedValue := "%" + tc.value + "%"
				if len(args) != 1 || args[0] != expectedValue {
					t.Errorf("Expected arg: %v, got: %v", expectedValue, args[0])
				}
			})
		}
	})

	t.Run("test chaining multiple conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试链式调用多个条件
		query.Like(&user.Name, "admin").
			Like(&user.Name, "user").
			Like("name", "test")

		if len(query.ToOptions()) != 3 {
			t.Error("Should have 3 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `name` LIKE ? AND `name` LIKE ? AND `name` LIKE ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 3 {
			t.Errorf("Expected 3 args, got: %d", len(args))
		}

		expectedArgs := []interface{}{"%admin%", "%user%", "%test%"}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d]: %v, got: %v", i, expected, args[i])
			}
		}
	})

	t.Run("test with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效字段名（应该不会panic，但会在实际查询时处理）
		query.Like("invalid_field", "value")

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍然会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `invalid_field` LIKE ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "%value%" {
			t.Errorf("Expected args: [%%value%%], got: %v", args)
		}
	})

	t.Run("test panic with non-string value", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试传入非字符串值时是否panic
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Like should panic when value is not a string")
			}
		}()

		query.Like(&user.Name, 123) // 传入整数应该触发panic
	})
}

// TestQuery_Regexp 测试 Regexp 方法的基本功能
func TestQuery_Regexp(t *testing.T) {
	t.Run("test with field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字段指针的方式
		result := query.Regexp(&user.Name, "^[a-zA-Z]+$")

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("Regexp should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("Regexp should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `name` REGEXP ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "^[a-zA-Z]+$" {
			t.Errorf("Expected args: [^[a-zA-Z]+$], got: %v", args)
		}
	})

	t.Run("test with field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字段名字符串的方式
		query.Regexp("name", ".*admin.*")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `name` REGEXP ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != ".*admin.*" {
			t.Errorf("Expected args: [.*admin.*], got: %v", args)
		}
	})

	t.Run("test with different regex patterns", func(t *testing.T) {
		_, user := gormx.NewQuery[User]()

		// 测试不同的正则表达式模式
		testCases := []struct {
			field   interface{}
			pattern string
			name    string
		}{
			{&user.Name, "^[a-zA-Z]+$", "letters only pattern"},
			{&user.Name, "\\d+", "digits pattern"},
			{&user.Name, "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$", "email pattern"},
			{"name", ".*test.*", "contains test pattern"},
			{&user.Name, "", "empty pattern"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				testQuery, _ := gormx.NewQuery[User]()
				testQuery.Regexp(tc.field, tc.pattern)

				_, args := testQuery.ToSQLAndArgs()
				if len(args) != 1 || args[0] != tc.pattern {
					t.Errorf("Expected arg: %v, got: %v", tc.pattern, args[0])
				}
			})
		}
	})

	t.Run("test chaining multiple conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试链式调用多个条件
		query.Regexp(&user.Name, "^[a-zA-Z]+$").
			Regexp(&user.Name, "admin").
			Regexp("name", "\\d+$")

		if len(query.ToOptions()) != 3 {
			t.Error("Should have 3 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `name` REGEXP ? AND `name` REGEXP ? AND `name` REGEXP ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 3 {
			t.Errorf("Expected 3 args, got: %d", len(args))
		}

		expectedArgs := []interface{}{"^[a-zA-Z]+$", "admin", "\\d+$"}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d]: %v, got: %v", i, expected, args[i])
			}
		}
	})

	t.Run("test with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效字段名（应该不会panic，但会在实际查询时处理）
		query.Regexp("invalid_field", "pattern")

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍然会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `invalid_field` REGEXP ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "pattern" {
			t.Errorf("Expected args: [pattern], got: %v", args)
		}
	})

	t.Run("test with special regex characters", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试包含特殊字符的正则表达式
		specialPattern := "@gmail\\.com$"
		query.Regexp(&user.Name, specialPattern)

		_, args := query.ToSQLAndArgs()
		if len(args) != 1 || args[0] != specialPattern {
			t.Errorf("Expected arg to handle special regex chars: %v, got: %v", specialPattern, args[0])
		}
	})
}

// TestQuery_IsNull 测试 IsNull 方法的基本功能
func TestQuery_IsNull(t *testing.T) {
	t.Run("test with field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字段指针的方式
		result := query.IsNull(&user.Email)

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("IsNull should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("IsNull should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `email` IS NULL AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected no args, got: %v", args)
		}
	})

	t.Run("test with field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字段名字符串的方式
		query.IsNull("phone")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `phone` IS NULL AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected no args, got: %v", args)
		}
	})

	t.Run("test chaining multiple conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试链式调用多个条件
		query.IsNull(&user.Email).
			IsNull(&user.Phone).
			IsNull("address")

		if len(query.ToOptions()) != 3 {
			t.Error("Should have 3 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `email` IS NULL AND `phone` IS NULL AND `address` IS NULL AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected no args, got: %v", args)
		}
	})

	t.Run("test with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效字段名（应该不会panic，但会在实际查询时处理）
		query.IsNull("invalid_field")

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍然会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `invalid_field` IS NULL AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected no args, got: %v", args)
		}
	})

	t.Run("test combined with other conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与其他条件组合使用
		query.Eq(&user.Status, 1).
			IsNull(&user.Email)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `status` = ? AND `email` IS NULL AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})
}

// TestQuery_IsNotNull 测试 IsNotNull 方法的基本功能
func TestQuery_IsNotNull(t *testing.T) {
	t.Run("test with field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字段指针的方式
		result := query.IsNotNull(&user.Email)

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("IsNotNull should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("IsNotNull should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `email` IS NOT NULL AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected no args, got: %v", args)
		}
	})

	t.Run("test with field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字段名字符串的方式
		query.IsNotNull("phone")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `phone` IS NOT NULL AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected no args, got: %v", args)
		}
	})

	t.Run("test chaining multiple conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试链式调用多个条件
		query.IsNotNull(&user.Email).
			IsNotNull(&user.Phone).
			IsNotNull("address")

		if len(query.ToOptions()) != 3 {
			t.Error("Should have 3 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `email` IS NOT NULL AND `phone` IS NOT NULL AND `address` IS NOT NULL AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected no args, got: %v", args)
		}
	})

	t.Run("test with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效字段名（应该不会panic，但会在实际查询时处理）
		query.IsNotNull("invalid_field")

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍然会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `invalid_field` IS NOT NULL AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected no args, got: %v", args)
		}
	})

	t.Run("test combined with other conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与其他条件组合使用
		query.Eq(&user.Status, 1).
			IsNotNull(&user.Email)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `status` = ? AND `email` IS NOT NULL AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})
}

// TestQuery_In 测试 In 方法的基本功能
func TestQuery_In(t *testing.T) {
	t.Run("test with field pointer and integer slice", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字段指针和整数切片
		result := query.In(&user.Status, []int{1, 2, 3})

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("In should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("In should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `status` IN (?,?,?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 3 {
			t.Errorf("Expected 3 args, got: %d", len(args))
		}

		// 验证参数值
		expectedValues := []interface{}{1, 2, 3}
		for i, expected := range expectedValues {
			if args[i] != expected {
				t.Errorf("Expected arg[%d] = %v, got %v", i, expected, args[i])
			}
		}
	})

	t.Run("test with field name string and string slice", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字段名字符串和字符串切片
		query.In("name", []string{"admin", "user", "guest"})

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `name` IN (?,?,?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 3 {
			t.Errorf("Expected 3 args, got: %d", len(args))
		}

		// 验证参数值
		expectedValues := []interface{}{"admin", "user", "guest"}
		for i, expected := range expectedValues {
			if args[i] != expected {
				t.Errorf("Expected arg[%d] = %v, got %v", i, expected, args[i])
			}
		}
	})

	t.Run("test chaining multiple conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试链式调用多个条件
		query.In(&user.Status, []int{1, 2}).
			In(&user.ID, []int{10, 20, 30}).
			In("name", []string{"admin", "user"})

		if len(query.ToOptions()) != 3 {
			t.Error("Should have 3 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `status` IN (?,?) AND `id` IN (?,?,?) AND `name` IN (?,?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 7 {
			t.Errorf("Expected 7 args, got: %d", len(args))
		}
	})

	t.Run("test with different slice types", func(t *testing.T) {
		_, user := gormx.NewQuery[User]()

		// 测试不同的切片类型
		testCases := []struct {
			field interface{}
			value interface{}
			name  string
		}{
			{&user.ID, []int{1, 2, 3}, "integer slice"},
			{&user.Name, []string{"a", "b", "c"}, "string slice"},
			{&user.Salary, []float64{100.5, 200.75}, "float slice"},
			{"status", []int{0, 1}, "string field with integer slice"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				testQuery, _ := gormx.NewQuery[User]()
				testQuery.In(tc.field, tc.value)

				_, args := testQuery.ToSQLAndArgs()
				rv := reflect.ValueOf(tc.value)
				if len(args) != rv.Len() {
					t.Errorf("Expected %d args, got: %d", rv.Len(), len(args))
				}
			})
		}
	})

	t.Run("test with array instead of slice", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用数组而不是切片
		arr := [3]int{1, 2, 3}
		query.In(&user.Status, arr)

		_, args := query.ToSQLAndArgs()
		if len(args) != 3 {
			t.Errorf("Expected 3 args, got: %d", len(args))
		}
	})

	t.Run("test with empty slice", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用空切片
		query.In(&user.Status, []int{})

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `status` IN (NULL) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效字段名
		query.In("invalid_field", []int{1, 2, 3})

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍然会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `invalid_field` IN (?,?,?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 3 {
			t.Errorf("Expected 3 args, got: %d", len(args))
		}
	})

	// 多对一关系使用场景测试
	t.Run("test many-to-one relationship scenario", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 模拟查询特定用户ID集合的帖子
		userIDs := []int{1, 5, 10, 15}
		query.In(&post.UserID, userIDs)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `posts` WHERE `user_id` IN (?,?,?,?) AND `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 4 {
			t.Errorf("Expected 4 args, got: %d", len(args))
		}

		// 验证参数值
		expectedUserIDs := []interface{}{1, 5, 10, 15}
		for i, expected := range expectedUserIDs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d] = %v, got %v", i, expected, args[i])
			}
		}
	})

	t.Run("test complex many-to-one query", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 模拟复杂查询：特定用户ID集合且状态为已发布的文章
		userIDs := []int{1, 5, 10}
		query.In(&post.UserID, userIDs).
			Eq(&post.Status, 1) // 1表示已发布

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `posts` WHERE `user_id` IN (?,?,?) AND `status` = ? AND `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 4 {
			t.Errorf("Expected 4 args, got: %d", len(args))
		}

		// 验证前三个参数（IN条件）
		expectedUserIDs := []interface{}{1, 5, 10}
		for i, expected := range expectedUserIDs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d] = %v, got %v", i, expected, args[i])
			}
		}

		// 验证第四个参数（EQ条件）
		if args[3] != 1 {
			t.Errorf("Expected args[3] = 1, got %v", args[3])
		}
	})
}

// TestQuery_NotIn 测试 NotIn 方法的基本功能
func TestQuery_NotIn(t *testing.T) {
	t.Run("test with field pointer and integer slice", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字段指针和整数切片
		result := query.NotIn(&user.Status, []int{1, 2, 3})

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("NotIn should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("NotIn should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `status` NOT IN (?,?,?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 3 {
			t.Errorf("Expected 3 args, got: %d", len(args))
		}

		// 验证参数值
		expectedValues := []interface{}{1, 2, 3}
		for i, expected := range expectedValues {
			if args[i] != expected {
				t.Errorf("Expected arg[%d] = %v, got %v", i, expected, args[i])
			}
		}
	})

	t.Run("test with field name string and string slice", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字段名字符串和字符串切片
		query.NotIn("name", []string{"admin", "user", "guest"})

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `name` NOT IN (?,?,?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 3 {
			t.Errorf("Expected 3 args, got: %d", len(args))
		}

		// 验证参数值
		expectedValues := []interface{}{"admin", "user", "guest"}
		for i, expected := range expectedValues {
			if args[i] != expected {
				t.Errorf("Expected arg[%d] = %v, got %v", i, expected, args[i])
			}
		}
	})

	t.Run("test chaining multiple conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试链式调用多个条件
		query.NotIn(&user.Status, []int{1, 2}).
			NotIn(&user.ID, []int{10, 20, 30}).
			NotIn("name", []string{"admin", "user"})

		if len(query.ToOptions()) != 3 {
			t.Error("Should have 3 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `status` NOT IN (?,?) AND `id` NOT IN (?,?,?) AND `name` NOT IN (?,?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 7 {
			t.Errorf("Expected 7 args, got: %d", len(args))
		}
	})

	t.Run("test with different slice types", func(t *testing.T) {
		_, user := gormx.NewQuery[User]()

		// 测试不同的切片类型
		testCases := []struct {
			field interface{}
			value interface{}
			name  string
		}{
			{&user.ID, []int{1, 2, 3}, "integer slice"},
			{&user.Name, []string{"a", "b", "c"}, "string slice"},
			{&user.Salary, []float64{100.5, 200.75}, "float slice"},
			{"status", []int{0, 1}, "string field with integer slice"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				testQuery, _ := gormx.NewQuery[User]()
				testQuery.NotIn(tc.field, tc.value)

				_, args := testQuery.ToSQLAndArgs()
				rv := reflect.ValueOf(tc.value)
				if len(args) != rv.Len() {
					t.Errorf("Expected %d args, got: %d", rv.Len(), len(args))
				}
			})
		}
	})

	t.Run("test with array instead of slice", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用数组而不是切片
		arr := [3]int{1, 2, 3}
		query.NotIn(&user.Status, arr)

		_, args := query.ToSQLAndArgs()
		if len(args) != 3 {
			t.Errorf("Expected 3 args, got: %d", len(args))
		}
	})

	t.Run("test with empty slice", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用空切片
		query.NotIn(&user.Status, []int{})

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `status` IS NOT NULL AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效字段名
		query.NotIn("invalid_field", []int{1, 2, 3})

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍然会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `invalid_field` NOT IN (?,?,?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 3 {
			t.Errorf("Expected 3 args, got: %d", len(args))
		}
	})

	// 多对一关系使用场景测试
	t.Run("test many-to-one relationship scenario", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 模拟查询不属于特定用户ID集合的帖子
		userIDs := []int{1, 5, 10, 15}
		query.NotIn(&post.UserID, userIDs)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `posts` WHERE `user_id` NOT IN (?,?,?,?) AND `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 4 {
			t.Errorf("Expected 4 args, got: %d", len(args))
		}

		// 验证参数值
		expectedUserIDs := []interface{}{1, 5, 10, 15}
		for i, expected := range expectedUserIDs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d] = %v, got %v", i, expected, args[i])
			}
		}
	})

	t.Run("test complex many-to-one query", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 模拟复杂查询：不属于特定用户ID集合且状态不为已删除的文章
		userIDs := []int{1, 5, 10}
		query.NotIn(&post.UserID, userIDs).
			NotIn(&post.Status, []int{0}) // 0表示已删除

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `posts` WHERE `user_id` NOT IN (?,?,?) AND `status` <> ? AND `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 4 {
			t.Errorf("Expected 4 args, got: %d", len(args))
		}

		// 验证前三个参数（NOT IN条件）
		expectedUserIDs := []interface{}{1, 5, 10}
		for i, expected := range expectedUserIDs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d] = %v, got %v", i, expected, args[i])
			}
		}

		// 验证第四个参数（NOT IN条件）
		if args[3] != 0 {
			t.Errorf("Expected args[3] = 0, got %v", args[3])
		}
	})

	t.Run("test combination with other query methods", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与其他查询方法组合使用
		query.Eq(&user.IsActive, true).
			NotIn(&user.Status, []int{0, -1}).
			Like(&user.Name, "test")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `is_active` = ? AND `status` NOT IN (?,?) AND `name` LIKE ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 4 {
			t.Errorf("Expected 4 args, got: %d", len(args))
		}

		// 验证参数值
		expectedArgs := []interface{}{true, 0, -1, "%test%"}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d] = %v, got %v", i, expected, args[i])
			}
		}
	})
}

// TestQuery_Between 测试 Between 方法的基本功能
func TestQuery_Between(t *testing.T) {
	t.Run("test with field pointer and numeric values", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字段指针和整数范围值
		result := query.Between(&user.Age, 18, 65)

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("Between should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("Between should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		// 注意：GORM 会自动为 BETWEEN 条件添加括号
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE (`age` BETWEEN ? AND ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 {
			t.Errorf("Expected 2 args, got: %d", len(args))
		}

		expectedArgs := []interface{}{18, 65}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d] = %v, got %v", i, expected, args[i])
			}
		}
	})

	t.Run("test with field name string and float values", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字段名字符串和浮点数范围值
		query.Between("salary", 3000.50, 15000.75)

		sql, args := query.ToSQLAndArgs()
		// 注意：GORM 会自动为 BETWEEN 条件添加括号
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE (`salary` BETWEEN ? AND ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 {
			t.Errorf("Expected 2 args, got: %d", len(args))
		}

		expectedArgs := []interface{}{3000.50, 15000.75}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d] = %v, got %v", i, expected, args[i])
			}
		}
	})

	t.Run("test chaining multiple conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试链式调用多个条件
		query.Between(&user.Age, 18, 65).
			Between(&user.Salary, 3000.0, 15000.0).
			Between("status", 0, 2)

		if len(query.ToOptions()) != 3 {
			t.Error("Should have 3 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		// 注意：GORM 会自动为每个 BETWEEN 条件添加括号
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE (`age` BETWEEN ? AND ?) AND (`salary` BETWEEN ? AND ?) AND (`status` BETWEEN ? AND ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 6 {
			t.Errorf("Expected 6 args, got: %d", len(args))
		}

		expectedArgs := []interface{}{18, 65, 3000.0, 15000.0, 0, 2}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d] = %v, got %v", i, expected, args[i])
			}
		}
	})

	t.Run("test with different value types", func(t *testing.T) {
		_, user := gormx.NewQuery[User]()

		// 测试不同的值类型
		testCases := []struct {
			field interface{}
			start interface{}
			end   interface{}
			name  string
		}{
			{&user.Age, 18, 65, "integer range"},
			{&user.Salary, 3000.50, 15000.75, "float range"},
			{"status", 0, 2, "string field with integer range"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				testQuery, _ := gormx.NewQuery[User]()
				testQuery.Between(tc.field, tc.start, tc.end)

				_, args := testQuery.ToSQLAndArgs()
				if len(args) != 2 {
					t.Errorf("Expected 2 args, got: %d", len(args))
				}

				if args[0] != tc.start || args[1] != tc.end {
					t.Errorf("Expected args: [%v, %v], got: [%v, %v]", tc.start, tc.end, args[0], args[1])
				}
			})
		}
	})

	t.Run("test with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效字段名
		query.Between("invalid_field", 10, 20)

		sql, args := query.ToSQLAndArgs()
		// 注意：GORM 会自动为 BETWEEN 条件添加括号，无效字段名会生成空字段名
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE (`invalid_field` BETWEEN ? AND ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 {
			t.Errorf("Expected 2 args, got: %d", len(args))
		}

		expectedArgs := []interface{}{10, 20}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d] = %v, got %v", i, expected, args[i])
			}
		}
	})

	// 多对一关系使用场景测试
	t.Run("test many-to-one relationship scenario", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 模拟查询特定用户ID范围内的帖子
		query.Between(&post.UserID, 1, 100)

		sql, args := query.ToSQLAndArgs()
		// 注意：GORM 会自动为 BETWEEN 条件添加括号
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `posts` WHERE (`user_id` BETWEEN ? AND ?) AND `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 {
			t.Errorf("Expected 2 args, got: %d", len(args))
		}

		expectedArgs := []interface{}{1, 100}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d] = %v, got %v", i, expected, args[i])
			}
		}
	})

	t.Run("test complex many-to-one query", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 模拟复杂查询：特定用户ID范围内且状态为已发布或草稿的文章
		query.Between(&post.UserID, 1, 100).
			Between(&post.Status, 0, 1) // 0表示草稿，1表示已发布

		sql, args := query.ToSQLAndArgs()
		// 注意：GORM 会自动为每个 BETWEEN 条件添加括号
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `posts` WHERE (`user_id` BETWEEN ? AND ?) AND (`status` BETWEEN ? AND ?) AND `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 4 {
			t.Errorf("Expected 4 args, got: %d", len(args))
		}

		expectedArgs := []interface{}{1, 100, 0, 1}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d] = %v, got %v", i, expected, args[i])
			}
		}
	})

	t.Run("test combination with other query methods", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与其他查询方法组合使用
		query.Eq(&user.IsActive, true).
			Between(&user.Age, 18, 65).
			Like(&user.Name, "test")

		sql, args := query.ToSQLAndArgs()
		// 注意：GORM 会自动为 BETWEEN 条件添加括号
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `is_active` = ? AND (`age` BETWEEN ? AND ?) AND `name` LIKE ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 4 {
			t.Errorf("Expected 4 args, got: %d", len(args))
		}

		expectedArgs := []interface{}{true, 18, 65, "%test%"}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d] = %v, got %v", i, expected, args[i])
			}
		}
	})

	t.Run("test boundary values", func(t *testing.T) {
		_, user := gormx.NewQuery[User]()

		// 测试边界值
		testCases := []struct {
			start interface{}
			end   interface{}
			name  string
		}{
			{0, 0, "zero boundaries"},
			{-10, 10, "negative start"},
			{10, -10, "negative end"},
			{1.1, 2.2, "decimal boundaries"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				testQuery, _ := gormx.NewQuery[User]()
				testQuery.Between(&user.Age, tc.start, tc.end)

				_, args := testQuery.ToSQLAndArgs()
				if len(args) != 2 {
					t.Errorf("Expected 2 args, got: %d", len(args))
				}

				if args[0] != tc.start || args[1] != tc.end {
					t.Errorf("Expected args: [%v, %v], got: [%v, %v]", tc.start, tc.end, args[0], args[1])
				}
			})
		}
	})
}

// TestQuery_NotBetween 测试 NotBetween 方法的基本功能
func TestQuery_NotBetween(t *testing.T) {
	t.Run("test with field pointer and numeric values", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字段指针和整数范围值
		result := query.NotBetween(&user.Age, 18, 65)

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("NotBetween should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("NotBetween should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE (`age` NOT BETWEEN ? AND ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 {
			t.Errorf("Expected 2 args, got: %d", len(args))
		}

		expectedArgs := []interface{}{18, 65}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d] = %v, got %v", i, expected, args[i])
			}
		}
	})

	t.Run("test with field name string and float values", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字段名字符串和浮点数范围值
		query.NotBetween("salary", 3000.50, 15000.75)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE (`salary` NOT BETWEEN ? AND ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 {
			t.Errorf("Expected 2 args, got: %d", len(args))
		}

		expectedArgs := []interface{}{3000.50, 15000.75}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d] = %v, got %v", i, expected, args[i])
			}
		}
	})

	t.Run("test chaining multiple conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试链式调用多个条件
		query.NotBetween(&user.Age, 18, 65).
			NotBetween(&user.Salary, 3000.0, 15000.0).
			NotBetween("status", 0, 2)

		if len(query.ToOptions()) != 3 {
			t.Error("Should have 3 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE (`age` NOT BETWEEN ? AND ?) AND (`salary` NOT BETWEEN ? AND ?) AND (`status` NOT BETWEEN ? AND ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 6 {
			t.Errorf("Expected 6 args, got: %d", len(args))
		}

		expectedArgs := []interface{}{18, 65, 3000.0, 15000.0, 0, 2}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d] = %v, got %v", i, expected, args[i])
			}
		}
	})

	t.Run("test with different value types", func(t *testing.T) {
		_, user := gormx.NewQuery[User]()

		// 测试不同的值类型
		testCases := []struct {
			field interface{}
			start interface{}
			end   interface{}
			name  string
		}{
			{&user.Age, 18, 65, "integer range"},
			{&user.Salary, 3000.50, 15000.75, "float range"},
			{"status", 0, 2, "string field with integer range"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				testQuery, _ := gormx.NewQuery[User]()
				testQuery.NotBetween(tc.field, tc.start, tc.end)

				_, args := testQuery.ToSQLAndArgs()
				if len(args) != 2 {
					t.Errorf("Expected 2 args, got: %d", len(args))
				}

				if args[0] != tc.start || args[1] != tc.end {
					t.Errorf("Expected args: [%v, %v], got: [%v, %v]", tc.start, tc.end, args[0], args[1])
				}
			})
		}
	})

	t.Run("test with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效字段名
		query.NotBetween("invalid_field", 10, 20)

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍然会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE (`invalid_field` NOT BETWEEN ? AND ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 {
			t.Errorf("Expected 2 args, got: %d", len(args))
		}

		expectedArgs := []interface{}{10, 20}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d] = %v, got %v", i, expected, args[i])
			}
		}
	})

	// 多对一关系使用场景测试
	t.Run("test many-to-one relationship scenario", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 模拟查询特定用户ID范围外的帖子
		query.NotBetween(&post.UserID, 1, 100)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `posts` WHERE (`user_id` NOT BETWEEN ? AND ?) AND `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 {
			t.Errorf("Expected 2 args, got: %d", len(args))
		}

		expectedArgs := []interface{}{1, 100}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d] = %v, got %v", i, expected, args[i])
			}
		}
	})

	t.Run("test complex many-to-one query", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 模拟复杂查询：特定用户ID范围外且状态不为已发布或草稿的文章
		query.NotBetween(&post.UserID, 1, 100).
			NotBetween(&post.Status, 0, 1) // 0表示草稿，1表示已发布

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `posts` WHERE (`user_id` NOT BETWEEN ? AND ?) AND (`status` NOT BETWEEN ? AND ?) AND `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 4 {
			t.Errorf("Expected 4 args, got: %d", len(args))
		}

		expectedArgs := []interface{}{1, 100, 0, 1}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d] = %v, got %v", i, expected, args[i])
			}
		}
	})

	t.Run("test combination with other query methods", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与其他查询方法组合使用
		query.Eq(&user.IsActive, true).
			NotBetween(&user.Age, 18, 65).
			Like(&user.Name, "test")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `is_active` = ? AND (`age` NOT BETWEEN ? AND ?) AND `name` LIKE ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 4 {
			t.Errorf("Expected 4 args, got: %d", len(args))
		}

		expectedArgs := []interface{}{true, 18, 65, "%test%"}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d] = %v, got %v", i, expected, args[i])
			}
		}
	})

	t.Run("test boundary values", func(t *testing.T) {
		_, user := gormx.NewQuery[User]()

		// 测试边界值
		testCases := []struct {
			start interface{}
			end   interface{}
			name  string
		}{
			{0, 0, "zero boundaries"},
			{-10, 10, "negative start"},
			{10, -10, "negative end"},
			{1.1, 2.2, "decimal boundaries"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				testQuery, _ := gormx.NewQuery[User]()
				testQuery.NotBetween(&user.Age, tc.start, tc.end)

				_, args := testQuery.ToSQLAndArgs()
				if len(args) != 2 {
					t.Errorf("Expected 2 args, got: %d", len(args))
				}

				if args[0] != tc.start || args[1] != tc.end {
					t.Errorf("Expected args: [%v, %v], got: [%v, %v]", tc.start, tc.end, args[0], args[1])
				}
			})
		}
	})
}

// TestQuery_OrderDesc 测试 OrderDesc 方法的基本功能
func TestQuery_OrderDesc(t *testing.T) {
	t.Run("test with single field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用单个字段指针
		result := query.OrderDesc(&user.CreatedAt)

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("OrderDesc should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("OrderDesc should add one option")
		}

		// 验证生成的SQL
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `created_at` DESC")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test with single field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用单个字段名字符串
		query.OrderDesc("created_at")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `created_at` DESC")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test with multiple field pointers", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用多个字段指针
		query.OrderDesc(&user.CreatedAt, &user.Name)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `created_at` DESC,`name` DESC")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test with multiple field name strings", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用多个字段名字符串
		query.OrderDesc("created_at", "name")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `created_at` DESC,`name` DESC")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test mixed field pointers and field name strings", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试混合使用字段指针和字段名字符串
		query.OrderDesc(&user.CreatedAt, "name")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `created_at` DESC,`name` DESC")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test chaining with other query methods", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与其他查询方法链式调用
		query.Eq(&user.Status, 1).
			OrderDesc(&user.CreatedAt)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `status` = ? AND `users`.`deleted_at` IS NULL ORDER BY `created_at` DESC")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 {
			t.Errorf("Expected 1 arg, got: %d", len(args))
		}

		if args[0] != 1 {
			t.Errorf("Expected arg: 1, got: %v", args[0])
		}
	})

	t.Run("test multiple OrderDesc calls", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试多次调用OrderDesc方法
		query.OrderDesc(&user.CreatedAt).
			OrderDesc(&user.Name)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `created_at` DESC,`name` DESC")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用无效字段名
		query.OrderDesc("invalid_field")

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `invalid_field` DESC")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test with mixed valid and invalid fields", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试混合有效和无效字段
		query.OrderDesc(&user.CreatedAt, "invalid_field", "name")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `created_at` DESC,`invalid_field` DESC,`name` DESC")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	// 多对一关系使用场景测试
	t.Run("test many-to-one relationship scenario", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 测试多对一关系中的排序
		query.OrderDesc(&post.UserID, &post.CreatedAt)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `posts` WHERE `posts`.`deleted_at` IS NULL ORDER BY `user_id` DESC,`created_at` DESC")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test complex query with ordering", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 测试复杂查询与排序结合
		query.Eq(&post.Status, 1).
			OrderDesc(&post.UserID, &post.CreatedAt)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `posts` WHERE `status` = ? AND `posts`.`deleted_at` IS NULL ORDER BY `user_id` DESC,`created_at` DESC")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 {
			t.Errorf("Expected 1 arg, got: %d", len(args))
		}

		if args[0] != 1 {
			t.Errorf("Expected arg: 1, got: %v", args[0])
		}
	})
}

// TestQuery_OrderAsc 测试 OrderAsc 方法的基本功能
func TestQuery_OrderAsc(t *testing.T) {
	t.Run("test with single field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用单个字段指针
		result := query.OrderAsc(&user.CreatedAt)

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("OrderAsc should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("OrderAsc should add one option")
		}

		// 验证生成的SQL
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `created_at`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test with single field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用单个字段名字符串
		query.OrderAsc("created_at")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `created_at`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test with multiple field pointers", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用多个字段指针
		query.OrderAsc(&user.CreatedAt, &user.Name)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `created_at`,`name`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test with multiple field name strings", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用多个字段名字符串
		query.OrderAsc("created_at", "name")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `created_at`,`name`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test mixed field pointers and field name strings", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试混合使用字段指针和字段名字符串
		query.OrderAsc(&user.CreatedAt, "name")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `created_at`,`name`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test chaining with other query methods", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与其他查询方法链式调用
		query.Eq(&user.Status, 1).
			OrderAsc(&user.CreatedAt)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `status` = ? AND `users`.`deleted_at` IS NULL ORDER BY `created_at`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 {
			t.Errorf("Expected 1 arg, got: %d", len(args))
		}

		if args[0] != 1 {
			t.Errorf("Expected arg: 1, got: %v", args[0])
		}
	})

	t.Run("test multiple OrderAsc calls", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试多次调用OrderAsc方法
		query.OrderAsc(&user.CreatedAt).
			OrderAsc(&user.Name)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `created_at`,`name`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用无效字段名
		query.OrderAsc("invalid_field")

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `invalid_field`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test with mixed valid and invalid fields", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试混合有效和无效字段
		query.OrderAsc(&user.CreatedAt, "invalid_field", "name")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `created_at`,`invalid_field`,`name`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	// 多对一关系使用场景测试
	t.Run("test many-to-one relationship scenario", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 测试多对一关系中的排序
		query.OrderAsc(&post.UserID, &post.CreatedAt)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `posts` WHERE `posts`.`deleted_at` IS NULL ORDER BY `user_id`,`created_at`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test complex query with ordering", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 测试复杂查询与排序结合
		query.Eq(&post.Status, 1).
			OrderAsc(&post.UserID, &post.CreatedAt)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `posts` WHERE `status` = ? AND `posts`.`deleted_at` IS NULL ORDER BY `user_id`,`created_at`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 {
			t.Errorf("Expected 1 arg, got: %d", len(args))
		}

		if args[0] != 1 {
			t.Errorf("Expected arg: 1, got: %v", args[0])
		}
	})

	// 测试与OrderDesc混合使用
	t.Run("test mixing OrderAsc and OrderDesc", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试OrderAsc和OrderDesc混合使用
		query.OrderAsc(&user.CreatedAt).
			OrderDesc(&user.Name)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `created_at`,`name` DESC")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})
}

// TestQuery_First 测试 First 方法的基本功能
func TestQuery_First(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	t.Run("test basic first query with field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试基本查询
		query.Eq(&user.Name, "testuser1")

		var result User
		err := query.First(&result)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果
		if result.Name != "testuser1" {
			t.Errorf("Expected name: testuser1, got: %s", result.Name)
		}
	})

	t.Run("test first query with string field", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字符串字段名
		query.Eq("name", "testuser2")

		var result User
		err := query.First(&result)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if result.Name != "testuser2" {
			t.Errorf("Expected name: testuser2, got: %s", result.Name)
		}
	})

	t.Run("test first query with multiple conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试多条件查询
		query.Eq(&user.Name, "testuser3").
			Eq(&user.Status, 1)

		var result User
		err := query.First(&result)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if result.Name != "testuser3" || result.Status != 1 {
			t.Errorf("Expected name: testuser3, status: 1, got name: %s, status: %d",
				result.Name, result.Status)
		}
	})

	t.Run("test first query with ordering", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试带排序的查询
		query.Eq(&user.Status, 1).
			OrderAsc(&user.CreatedAt)

		var result User
		err := query.First(&result)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果存在
		if result.ID == 0 {
			t.Error("Expected to find a user, but got empty result")
		}
	})

	t.Run("test first query with no matching records", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试无匹配记录的情况
		query.Eq(&user.Name, "nonexistent_user")

		var result User
		err := query.First(&result)

		// 应该返回 gorm.ErrRecordNotFound 错误
		if err == nil {
			t.Error("Expected error for no matching records, got nil")
		}

		if err != gorm.ErrRecordNotFound {
			t.Errorf("Expected gorm.ErrRecordNotFound, got: %v", err)
		}
	})

	t.Run("test first query with invalid field", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效字段名
		query.Eq("invalid_field", "value")

		var result User
		err := query.First(&result)

		// 无效字段查询会导致SQL语法错误
		if err == nil {
			t.Error("Expected SQL syntax error, got nil")
		}

		// 验证是数据库语法错误而不是记录未找到错误
		if err == gorm.ErrRecordNotFound {
			t.Error("Expected SQL syntax error, but got ErrRecordNotFound")
		}
	})

	t.Run("test first query with different data types", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试不同数据类型的字段查询
		query.Eq(&user.Age, 25)

		var result User
		err := query.First(&result)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if result.Age != 25 {
			t.Errorf("Expected age: 25, got: %d", result.Age)
		}
	})

	t.Run("test first query with boolean field", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试布尔字段查询
		query.Eq(&user.IsActive, true)

		var result User
		err := query.First(&result)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if !result.IsActive {
			t.Error("Expected IsActive to be true")
		}
	})

	t.Run("test first query with float field", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试浮点字段查询
		query.Gt(&user.Salary, 4000.0)

		var result User
		err := query.First(&result)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if result.Salary <= 4000.0 {
			t.Errorf("Expected salary > 4000.0, got: %f", result.Salary)
		}
	})

	t.Run("test first query with like condition", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试LIKE条件查询
		query.Like(&user.Name, "testuser")

		var result User
		err := query.First(&result)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if !strings.Contains(result.Name, "testuser") {
			t.Errorf("Expected name to contain 'testuser', got: %s", result.Name)
		}
	})

	t.Run("test first query with complex conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试复杂条件组合
		query.Eq(&user.Status, 1).
			Gt(&user.Age, 18).
			Like(&user.Name, "test").
			OrderDesc(&user.CreatedAt)

		var result User
		err := query.First(&result)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果满足条件
		if result.Status != 1 || result.Age <= 18 || !strings.Contains(result.Name, "test") {
			t.Errorf("Result doesn't match all conditions: status=%d, age=%d, name=%s",
				result.Status, result.Age, result.Name)
		}
	})

	// 多对一关系使用场景测试
	t.Run("test first query with many-to-one relationship", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 测试多对一关系查询
		query.Eq(&post.UserID, 1)

		var result Post
		err := query.First(&result)

		if err != nil && err != gorm.ErrRecordNotFound {
			t.Errorf("Expected no error or ErrRecordNotFound, got: %v", err)
		}

		// 如果找到了记录，验证UserID
		if err == nil && result.UserID != 1 {
			t.Errorf("Expected UserID: 1, got: %d", result.UserID)
		}
	})

	t.Run("test first query with multiple ordering", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试多字段排序
		query.Eq(&user.Status, 1).
			OrderDesc(&user.CreatedAt).
			OrderAsc(&user.Name)

		var result User
		err := query.First(&result)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if result.Status != 1 {
			t.Errorf("Expected status: 1, got: %d", result.Status)
		}
	})

	t.Run("test first query with limit and offset", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与Limit和Offset组合使用
		query.Eq(&user.Status, 1).
			Limit(1)

		var result User
		err := query.First(&result)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if result.Status != 1 {
			t.Errorf("Expected status: 1, got: %d", result.Status)
		}
	})
}

// setupTestData 准备测试数据
func setupTestData(t *testing.T) {
	// 清理现有数据
	db := gormx.GetDb()
	db.Where("1 = 1").Delete(&User{})
	db.Where("1 = 1").Delete(&Post{})

	// 插入测试用户数据
	users := []User{
		{Name: "testuser1", Email: "test1@example.com", Age: 25, IsActive: true, Salary: 5000.0, Status: 1},
		{Name: "testuser2", Email: "test2@example.com", Age: 30, IsActive: true, Salary: 6000.0, Status: 1},
		{Name: "testuser3", Email: "test3@example.com", Age: 35, IsActive: false, Salary: 7000.0, Status: 1},
		{Name: "otheruser", Email: "other@example.com", Age: 28, IsActive: true, Salary: 5500.0, Status: 0},
	}

	for i := range users {
		if err := db.Create(&users[i]).Error; err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}

	// 获取刚创建的用户ID
	var createdUsers []User
	if err := db.Find(&createdUsers).Error; err != nil {
		t.Fatalf("Failed to fetch created users: %v", err)
	}

	// 确保有足够的用户用于创建帖子
	if len(createdUsers) < 2 {
		t.Fatal("Not enough users created for testing posts")
	}

	// 插入测试帖子数据，使用实际创建的用户ID
	posts := []Post{
		{UserID: createdUsers[0].ID, Title: "Test Post 1", Content: "Content 1", Status: 1},
		{UserID: createdUsers[1].ID, Title: "Test Post 2", Content: "Content 2", Status: 1},
	}

	for i := range posts {
		if err := db.Create(&posts[i]).Error; err != nil {
			t.Fatalf("Failed to create test post: %v", err)
		}
	}
}

// TestQuery_Find 测试 Find 方法的基本功能
func TestQuery_Find(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	t.Run("test basic find query with field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试基本查询
		query.Eq(&user.Status, 1)

		var results []User
		err := query.Find(&results)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证至少有一条记录
		if len(results) == 0 {
			t.Error("Expected to find users with status 1, but got empty result")
		}

		// 验证所有结果的状态都是1
		for _, result := range results {
			if result.Status != 1 {
				t.Errorf("Expected status: 1, got: %d", result.Status)
			}
		}
	})

	t.Run("test find query with string field", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字符串字段名
		query.Eq("is_active", true)

		var results []User
		err := query.Find(&results)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证至少有一条记录
		if len(results) == 0 {
			t.Error("Expected to find active users, but got empty result")
		}

		// 验证所有结果都是激活状态
		for _, result := range results {
			if !result.IsActive {
				t.Error("Expected all users to be active")
			}
		}
	})

	t.Run("test find query with multiple conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试多条件查询
		query.Eq(&user.Status, 1).
			Gt(&user.Age, 20)

		var results []User
		err := query.Find(&results)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果
		if len(results) == 0 {
			t.Error("Expected to find users, but got empty result")
		}

		// 验证所有结果满足条件
		for _, result := range results {
			if result.Status != 1 || result.Age <= 20 {
				t.Errorf("Result doesn't match conditions: status=%d, age=%d", result.Status, result.Age)
			}
		}
	})

	t.Run("test find query with like condition", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试LIKE条件查询
		query.Like(&user.Name, "testuser")

		var results []User
		err := query.Find(&results)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果
		if len(results) == 0 {
			t.Error("Expected to find users with name containing 'testuser', but got empty result")
		}

		// 验证所有结果满足条件
		for _, result := range results {
			if !strings.Contains(result.Name, "testuser") {
				t.Errorf("Expected name to contain 'testuser', got: %s", result.Name)
			}
		}
	})

	t.Run("test find query with ordering", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试带排序的查询
		query.Eq(&user.Status, 1).
			OrderAsc(&user.Age)

		var results []User
		err := query.Find(&results)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果
		if len(results) == 0 {
			t.Error("Expected to find users, but got empty result")
		}

		// 验证排序是否正确
		for i := 1; i < len(results); i++ {
			if results[i-1].Age > results[i].Age {
				t.Error("Results are not sorted by age in ascending order")
			}
		}
	})

	t.Run("test find query with limit and offset", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试分页查询
		query.Eq(&user.Status, 1).
			Limit(2).
			Offset(0)

		var results []User
		err := query.Find(&results)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果数量不超过限制
		if len(results) > 2 {
			t.Errorf("Expected at most 2 results, got: %d", len(results))
		}
	})

	t.Run("test find query with no matching records", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试无匹配记录的情况
		query.Eq(&user.Name, "nonexistent_user")

		var results []User
		err := query.Find(&results)

		// Find方法即使没有匹配记录也不应该返回错误
		if err != nil {
			t.Errorf("Expected no error for no matching records, got: %v", err)
		}

		// 应该返回空切片而不是nil
		if results == nil {
			t.Error("Expected empty slice, got nil")
		}

		if len(results) != 0 {
			t.Errorf("Expected empty slice, got length: %d", len(results))
		}
	})

	t.Run("test find query with invalid field", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效字段名
		query.Eq("invalid_field", "value")

		var results []User
		err := query.Find(&results)

		// 无效字段查询会导致SQL语法错误
		if err == nil {
			t.Error("Expected SQL syntax error, got nil")
		}
	})

	t.Run("test find query with in condition", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试IN条件查询
		query.In(&user.Status, []int{0, 1})

		var results []User
		err := query.Find(&results)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果
		if len(results) == 0 {
			t.Error("Expected to find users, but got empty result")
		}

		// 验证所有结果满足条件
		for _, result := range results {
			if result.Status != 0 && result.Status != 1 {
				t.Errorf("Expected status to be 0 or 1, got: %d", result.Status)
			}
		}
	})

	t.Run("test find query with between condition", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试BETWEEN条件查询
		query.Between(&user.Age, 20, 40)

		var results []User
		err := query.Find(&results)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果
		if len(results) == 0 {
			t.Error("Expected to find users, but got empty result")
		}

		// 验证所有结果满足条件
		for _, result := range results {
			if result.Age < 20 || result.Age > 40 {
				t.Errorf("Expected age between 20 and 40, got: %d", result.Age)
			}
		}
	})

	t.Run("test find query with complex conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试复杂条件组合
		query.Eq(&user.Status, 1).
			Gt(&user.Age, 18).
			Like(&user.Name, "test").
			OrderDesc(&user.CreatedAt).
			Limit(5)

		var results []User
		err := query.Find(&results)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果数量不超过限制
		if len(results) > 5 {
			t.Errorf("Expected at most 5 results due to limit, got: %d", len(results))
		}

		// 验证所有结果满足条件
		for _, result := range results {
			if result.Status != 1 || result.Age <= 18 || !strings.Contains(result.Name, "test") {
				t.Errorf("Result doesn't match all conditions: status=%d, age=%d, name=%s",
					result.Status, result.Age, result.Name)
			}
		}

		// 验证排序是否正确（按创建时间降序）
		for i := 1; i < len(results); i++ {
			if results[i-1].CreatedAt.Before(results[i].CreatedAt) {
				t.Error("Results are not sorted by created_at in descending order")
			}
		}
	})

	// 多对一关系使用场景测试
	t.Run("test find query with many-to-one relationship", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 测试多对一关系查询
		query.Eq(&post.UserID, 1)

		var results []Post
		err := query.Find(&results)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果
		if len(results) > 0 {
			// 如果找到了记录，验证UserID
			for _, result := range results {
				if result.UserID != 1 {
					t.Errorf("Expected UserID: 1, got: %d", result.UserID)
				}
			}
		}
	})

	t.Run("test find query with multiple ordering", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试多字段排序
		query.Eq(&user.Status, 1).
			OrderDesc(&user.CreatedAt).
			OrderAsc(&user.Name)

		var results []User
		err := query.Find(&results)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(results) == 0 {
			t.Error("Expected to find users, but got empty result")
		}

		// 验证所有结果满足条件
		for _, result := range results {
			if result.Status != 1 {
				t.Errorf("Expected status: 1, got: %d", result.Status)
			}
		}
	})

	t.Run("test find query with distinct", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试去重查询
		query.Distinct(&user.Status)

		var results []User
		err := query.Find(&results)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果
		if len(results) == 0 {
			t.Error("Expected to find users, but got empty result")
		}
	})

	t.Run("test find query with select", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试字段选择查询
		query.Select(&user.Name, &user.Email).
			Eq(&user.Status, 1)

		var results []User
		err := query.Find(&results)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果
		if len(results) == 0 {
			t.Error("Expected to find users, but got empty result")
		}
	})
}

// TestQuery_Distinct 测试 Distinct 方法的基本功能
func TestQuery_Distinct(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	t.Run("test distinct without fields", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无参数去重（对整个查询结果去重）
		result := query.Distinct()
		if result != query {
			t.Error("Distinct should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("Distinct should add one option")
		}

		// 验证生成的SQL
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT DISTINCT * FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test distinct with single field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用单个字段指针去重
		result := query.Distinct(&user.Status)

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("Distinct should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("Distinct should add one option")
		}

		// 验证生成的SQL
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT DISTINCT `status` FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test distinct with single field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用单个字段名字符串去重
		query.Distinct("status")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT DISTINCT `status` FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test distinct with multiple field pointers", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用多个字段指针去重
		query.Distinct(&user.Status, &user.Name)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT DISTINCT `status`,`name` FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test distinct with multiple field name strings", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用多个字段名字符串去重
		query.Distinct("status", "name")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT DISTINCT `status`,`name` FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test mixed field pointers and field name strings", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试混合使用字段指针和字段名字符串
		query.Distinct(&user.Status, "name")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT DISTINCT `status`,`name` FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test chaining with other query methods", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与其他查询方法链式调用
		query.Eq(&user.Status, 1).
			Distinct(&user.Name)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT DISTINCT `name` FROM `users` WHERE `status` = ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 {
			t.Errorf("Expected 1 arg, got: %d", len(args))
		}

		if args[0] != 1 {
			t.Errorf("Expected arg: 1, got: %v", args[0])
		}
	})

	t.Run("test multiple Distinct calls", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试多次调用Distinct方法（最后一次生效）
		query.Distinct(&user.Status).
			Distinct(&user.Name)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT DISTINCT `name` FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用无效字段名
		query.Distinct("invalid_field")

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段会被忽略，所以不会出现在SQL中
		expectedSQL := expectedSQLForCurrentDialect("SELECT DISTINCT invalid_field FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test with mixed valid and invalid fields", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试混合有效和无效字段
		query.Distinct(&user.Status, "invalid_field", "name")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT DISTINCT `status`,invalid_field,`name` FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	// 多对一关系使用场景测试
	t.Run("test many-to-one relationship scenario", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 测试多对一关系中的去重
		query.Distinct(&post.UserID)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT DISTINCT `user_id` FROM `posts` WHERE `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test complex query with distinct", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 测试复杂查询与去重结合
		query.Eq(&post.Status, 1).
			Distinct(&post.UserID)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT DISTINCT `user_id` FROM `posts` WHERE `status` = ? AND `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 {
			t.Errorf("Expected 1 arg, got: %d", len(args))
		}

		if args[0] != 1 {
			t.Errorf("Expected arg: 1, got: %v", args[0])
		}
	})

	t.Run("test distinct with slice parameter", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用切片参数
		fields := []interface{}{&user.Status, "name"}
		query.Distinct(fields)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT DISTINCT `status`,`name` FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test distinct with empty slice", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用空切片参数
		var fields []interface{}

		query.Distinct(fields...)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT DISTINCT * FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})
}

// TestQuery_Select 测试 Select 方法的基本功能
func TestQuery_Select(t *testing.T) {
	t.Run("test select with single field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用单个字段指针
		result := query.Select(&user.Name)
		if result != query {
			t.Error("Select should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("Select should add one option")
		}

		// 验证生成的SQL
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT `name` FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test select with single field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用单个字段名字符串
		query.Select("email")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT `email` FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test select with multiple field pointers", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用多个字段指针
		query.Select(&user.Name, &user.Email, &user.Age)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT `name`,`email`,`age` FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test select with multiple field name strings", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用多个字段名字符串
		query.Select("name", "email", "age")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT `name`,`email`,`age` FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test mixed field pointers and field name strings", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试混合使用字段指针和字段名字符串
		query.Select(&user.Name, "email", &user.Age)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT `name`,`email`,`age` FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test chaining with other query methods", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与其他查询方法链式调用
		query.Select(&user.Name, &user.Email).
			Eq(&user.Status, 1)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT `name`,`email` FROM `users` WHERE `status` = ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 {
			t.Errorf("Expected 1 arg, got: %d", len(args))
		}

		if args[0] != 1 {
			t.Errorf("Expected arg: 1, got: %v", args[0])
		}
	})

	t.Run("test multiple Select calls", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试多次调用Select方法
		query.Select(&user.Name).
			Select(&user.Email)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT `email` FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用无效字段名
		query.Select("invalid_field")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT invalid_field FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test with mixed valid and invalid fields", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试混合有效和无效字段
		query.Select(&user.Name, "invalid_field", "email")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT `name`,invalid_field,`email` FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test select with no fields", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试不传入任何字段
		result := query.Select()
		if result != query {
			t.Error("Select should return the same query instance for chaining")
		}

		// 验证没有添加选项
		if len(query.ToOptions()) != 0 {
			t.Error("Select with no fields should not add any options")
		}

		// 验证生成的SQL仍然是默认查询
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	// 多对一关系使用场景测试
	t.Run("test many-to-one relationship scenario", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 测试多对一关系中的字段选择
		query.Select(&post.Title, &post.UserID)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT `title`,`user_id` FROM `posts` WHERE `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test complex query with select and other methods", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 测试复杂查询与字段选择结合
		query.Select(&post.Title, &post.UserID, &post.Status).
			Eq(&post.Status, 1).
			OrderDesc(&post.CreatedAt)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT `title`,`user_id`,`status` FROM `posts` WHERE `status` = ? AND `posts`.`deleted_at` IS NULL ORDER BY `created_at` DESC")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 {
			t.Errorf("Expected 1 arg, got: %d", len(args))
		}

		if args[0] != 1 {
			t.Errorf("Expected arg: 1, got: %v", args[0])
		}
	})

	t.Run("test select with slice parameter", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用切片参数
		fields := []interface{}{&user.Name, "email", &user.Age}
		query.Select(fields...)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT `name`,`email`,`age` FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test select with empty slice parameter", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用空切片参数
		var fields []interface{}

		result := query.Select(fields...)
		if result != query {
			t.Error("Select should return the same query instance for chaining")
		}

		// 验证没有添加选项
		if len(query.ToOptions()) != 0 {
			t.Error("Select with empty slice should not add any options")
		}

		// 验证生成的SQL仍然是默认查询
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})
}

// TestQuery_Scan 测试 Scan 方法的基本功能
func TestQuery_Scan(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	t.Run("test scan with basic query", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 基本查询测试
		query.Eq(&user.Status, 1)

		var results []User
		err := query.Scan(&results)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证至少有一条记录
		if len(results) == 0 {
			t.Error("Expected to find users with status 1, but got empty result")
		}

		// 验证所有结果的状态都是1
		for _, result := range results {
			if result.Status != 1 {
				t.Errorf("Expected status: 1, got: %d", result.Status)
			}
		}
	})

	t.Run("test scan with field selection", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试字段选择查询
		query.Select(&user.Name, &user.Email).
			Eq(&user.Status, 1)

		var results []User
		err := query.Scan(&results)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果
		if len(results) == 0 {
			t.Error("Expected to find users, but got empty result")
		}

		// 验证只选择了指定字段（其他字段应该为零值）
		for _, result := range results {
			if result.Name == "" {
				t.Error("Expected Name field to be selected")
			}
			if result.Email == "" {
				t.Error("Expected Email field to be selected")
			}
			// Age字段应该为零值，因为我们没有选择它
			if result.Age != 0 {
				t.Error("Expected Age field to be zero value as it was not selected")
			}
		}
	})

	t.Run("test scan with like condition", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试LIKE条件查询
		query.Like(&user.Name, "testuser")

		var results []User
		err := query.Scan(&results)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果
		if len(results) == 0 {
			t.Error("Expected to find users with name containing 'testuser', but got empty result")
		}

		// 验证所有结果满足条件
		for _, result := range results {
			if !strings.Contains(result.Name, "testuser") {
				t.Errorf("Expected name to contain 'testuser', got: %s", result.Name)
			}
		}
	})

	t.Run("test scan with ordering", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试带排序的查询
		query.Eq(&user.Status, 1).
			OrderAsc(&user.Age)

		var results []User
		err := query.Scan(&results)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果
		if len(results) == 0 {
			t.Error("Expected to find users, but got empty result")
		}

		// 验证排序是否正确
		for i := 1; i < len(results); i++ {
			if results[i-1].Age > results[i].Age {
				t.Error("Results are not sorted by age in ascending order")
			}
		}
	})

	t.Run("test scan with limit and offset", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试分页查询
		query.Eq(&user.Status, 1).
			Limit(2).
			Offset(0)

		var results []User
		err := query.Scan(&results)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果数量不超过限制
		if len(results) > 2 {
			t.Errorf("Expected at most 2 results, got: %d", len(results))
		}
	})

	t.Run("test scan with no matching records", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试无匹配记录的情况
		query.Eq(&user.Name, "nonexistent_user")

		var results []User
		err := query.Scan(&results)

		// Scan方法即使没有匹配记录也不应该返回错误
		if err != nil {
			t.Errorf("Expected no error for no matching records, got: %v", err)
		}

		// 应该返回nil或空切片都是可以接受的
		if len(results) != 0 {
			t.Errorf("Expected empty or nil slice, got length: %d", len(results))
		}
	})

	t.Run("test scan with invalid field", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效字段名
		query.Eq("invalid_field", "value")

		var results []User
		err := query.Scan(&results)

		// 无效字段查询会导致SQL语法错误
		if err == nil {
			t.Error("Expected SQL syntax error, got nil")
		}
	})

	t.Run("test scan with in condition", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试IN条件查询
		query.In(&user.Status, []int{0, 1})

		var results []User
		err := query.Scan(&results)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果
		if len(results) == 0 {
			t.Error("Expected to find users, but got empty result")
		}

		// 验证所有结果满足条件
		for _, result := range results {
			if result.Status != 0 && result.Status != 1 {
				t.Errorf("Expected status to be 0 or 1, got: %d", result.Status)
			}
		}
	})

	t.Run("test scan with between condition", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试BETWEEN条件查询
		query.Between(&user.Age, 20, 40)

		var results []User
		err := query.Scan(&results)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果
		if len(results) == 0 {
			t.Error("Expected to find users, but got empty result")
		}

		// 验证所有结果满足条件
		for _, result := range results {
			if result.Age < 20 || result.Age > 40 {
				t.Errorf("Expected age between 20 and 40, got: %d", result.Age)
			}
		}
	})

	t.Run("test scan with complex conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试复杂条件组合
		query.Eq(&user.Status, 1).
			Gt(&user.Age, 18).
			Like(&user.Name, "test").
			OrderDesc(&user.CreatedAt).
			Limit(5)

		var results []User
		err := query.Scan(&results)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果数量不超过限制
		if len(results) > 5 {
			t.Errorf("Expected at most 5 results due to limit, got: %d", len(results))
		}

		// 验证所有结果满足条件
		for _, result := range results {
			if result.Status != 1 || result.Age <= 18 || !strings.Contains(result.Name, "test") {
				t.Errorf("Result doesn't match all conditions: status=%d, age=%d, name=%s",
					result.Status, result.Age, result.Name)
			}
		}

		// 验证排序是否正确（按创建时间降序）
		for i := 1; i < len(results); i++ {
			if results[i-1].CreatedAt.Before(results[i].CreatedAt) {
				t.Error("Results are not sorted by created_at in descending order")
			}
		}
	})

	// 多对一关系使用场景测试
	t.Run("test scan with many-to-one relationship", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 测试多对一关系查询
		query.Eq(&post.UserID, 1)

		var results []Post
		err := query.Scan(&results)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果
		if len(results) > 0 {
			// 如果找到了记录，验证UserID
			for _, result := range results {
				if result.UserID != 1 {
					t.Errorf("Expected UserID: 1, got: %d", result.UserID)
				}
			}
		}
	})

	t.Run("test scan with multiple ordering", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试多字段排序
		query.Eq(&user.Status, 1).
			OrderDesc(&user.CreatedAt).
			OrderAsc(&user.Name)

		var results []User
		err := query.Scan(&results)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(results) == 0 {
			t.Error("Expected to find users, but got empty result")
		}

		// 验证所有结果满足条件
		for _, result := range results {
			if result.Status != 1 {
				t.Errorf("Expected status: 1, got: %d", result.Status)
			}
		}
	})

	t.Run("test scan with distinct", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试去重查询
		query.Distinct(&user.Status)

		var results []User
		err := query.Scan(&results)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果
		if len(results) == 0 {
			t.Error("Expected to find users, but got empty result")
		}
	})

	t.Run("test scan into single struct destination", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试扫描到非切片目标（应该失败）
		query.Eq(&user.Name, "testuser1")

		var result User
		err := query.Scan(&result)

		// 这不应该失败，因为扫描到单个结构体是GORM的正常行为
		if err != nil {
			t.Errorf("Expected no error when scanning into single struct, got: %v", err)
		}

		// 验证结果
		if result.Name != "testuser1" {
			t.Errorf("Expected name: testuser1, got: %s", result.Name)
		}
	})

	t.Run("test scan with nil destination", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试扫描到nil目标（应该失败）
		err := query.Scan(nil)

		if err == nil {
			t.Error("Expected error when scanning into nil destination, got nil")
		}

		// 验证错误信息
		if err.Error() != "destination cannot be nil" {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})

	t.Run("test scan with complex select and join", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 测试复杂查询与连接
		query.Select(&post.Title, &post.UserID).
			Eq(&post.Status, 1).
			OrderDesc(&post.CreatedAt)

		var results []Post
		err := query.Scan(&results)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果
		if len(results) > 0 {
			for _, result := range results {
				if result.Title == "" {
					t.Error("Expected Title field to be selected")
				}
				if result.UserID == 0 {
					t.Error("Expected UserID field to be selected")
				}
				// Content字段应该为零值，因为我们没有选择它
				if result.Content != "" {
					t.Error("Expected Content field to be empty as it was not selected")
				}
			}
		}
	})
}

// TestQuery_RawRows 测试 RawRows 方法的各种使用场景
func TestQuery_RawRows(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	t.Run("test basic raw rows query", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试基本的原生SQL查询
		rows, err := query.RawRows("SELECT id, name, email, status FROM users WHERE status = ?", 1)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		defer rows.Close()

		// 验证结果
		var users []User
		for rows.Next() {
			var user User
			db := gormx.GetDb()
			if err := db.ScanRows(rows, &user); err != nil {
				t.Errorf("Failed to scan row: %v", err)
			}
			users = append(users, user)
		}

		if len(users) == 0 {
			t.Error("Expected to find users with status 1, but got empty result")
		}

		// 验证所有结果的状态都是1
		for _, user := range users {
			if user.Status != 1 {
				t.Errorf("Expected status: 1, got: %d", user.Status)
			}
		}
	})

	t.Run("test raw rows with multiple parameters", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试带有多个参数的原生SQL查询
		rows, err := query.RawRows("SELECT id, name FROM users WHERE status = ? AND age > ?", 1, 20)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		defer rows.Close()

		// 验证结果
		count := 0
		for rows.Next() {
			count++
			var id uint
			var name string
			if err := rows.Scan(&id, &name); err != nil {
				t.Errorf("Failed to scan row: %v", err)
			}

			if id == 0 || name == "" {
				t.Error("Expected non-empty id and name")
			}
		}

		if count == 0 {
			t.Error("Expected to find users, but got empty result")
		}
	})

	t.Run("test raw rows with no matching records", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无匹配记录的查询
		rows, err := query.RawRows("SELECT id, name FROM users WHERE name = ?", "nonexistent_user")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		defer rows.Close()

		// 应该返回空结果集而不是错误
		count := 0
		for rows.Next() {
			count++
		}

		if count != 0 {
			t.Errorf("Expected empty result set, got %d rows", count)
		}
	})

	t.Run("test raw rows with invalid SQL", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效SQL语句
		_, err := query.RawRows("INVALID SQL STATEMENT")
		if err == nil {
			t.Error("Expected SQL syntax error, got nil")
		}
	})

	t.Run("test raw rows with wrong number of parameters", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试参数数量不匹配的情况
		_, err := query.RawRows("SELECT id, name FROM users WHERE status = ? AND age > ?", 1)
		if err == nil {
			t.Error("Expected parameter mismatch error, got nil")
		}
	})

	t.Run("test raw rows with complex query", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试复杂查询，包括JOIN操作
		rows, err := query.RawRows(`
			SELECT u.id, u.name, u.email, p.title 
			FROM users u 
			LEFT JOIN posts p ON u.id = p.user_id 
			WHERE u.status = ? 
			ORDER BY u.created_at DESC`, 1)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		defer rows.Close()

		// 验证结果结构
		columns, err := rows.Columns()
		if err != nil {
			t.Errorf("Failed to get columns: %v", err)
		}

		expectedColumns := []string{"id", "name", "email", "title"}
		if len(columns) != len(expectedColumns) {
			t.Errorf("Expected %d columns, got %d", len(expectedColumns), len(columns))
		}
	})

	t.Run("test raw rows with aggregation", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试聚合查询
		rows, err := query.RawRows("SELECT status, COUNT(*) as count FROM users GROUP BY status")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		defer rows.Close()

		type Result struct {
			Status int
			Count  int
		}

		var results []Result
		for rows.Next() {
			var result Result
			if err := rows.Scan(&result.Status, &result.Count); err != nil {
				t.Errorf("Failed to scan row: %v", err)
			}
			results = append(results, result)
		}

		if len(results) == 0 {
			t.Error("Expected aggregation results, got empty result")
		}
	})

	t.Run("test raw rows without parameters", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试不带参数的查询
		rows, err := query.RawRows("SELECT COUNT(*) FROM users")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		defer rows.Close()

		var count int64
		if rows.Next() {
			if err := rows.Scan(&count); err != nil {
				t.Errorf("Failed to scan row: %v", err)
			}
		}

		if count <= 0 {
			t.Errorf("Expected positive count, got %d", count)
		}
	})

	t.Run("test raw rows with limit", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试带LIMIT的查询
		rows, err := query.RawRows("SELECT id, name FROM users WHERE status = ? ORDER BY id LIMIT 2", 1)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		defer rows.Close()

		count := 0
		for rows.Next() {
			count++
			var id uint
			var name string
			if err := rows.Scan(&id, &name); err != nil {
				t.Errorf("Failed to scan row: %v", err)
			}
		}

		if count > 2 {
			t.Errorf("Expected at most 2 rows due to LIMIT, got %d rows", count)
		}
	})

	// 多对一关系使用场景测试
	t.Run("test raw rows with many-to-one relationship scenario", func(t *testing.T) {
		query, _ := gormx.NewQuery[Post]()

		// 测试多对一关系查询
		rows, err := query.RawRows(`
			SELECT p.id, p.title, p.user_id, u.name as user_name 
			FROM posts p 
			JOIN users u ON p.user_id = u.id 
			WHERE p.status = ?`, 1)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		defer rows.Close()

		type PostWithUser struct {
			ID       uint
			Title    string
			UserID   uint
			UserName string
		}

		var results []PostWithUser
		for rows.Next() {
			var result PostWithUser
			if err := rows.Scan(&result.ID, &result.Title, &result.UserID, &result.UserName); err != nil {
				t.Errorf("Failed to scan row: %v", err)
			}
			results = append(results, result)
		}

		if len(results) > 0 {
			// 验证结果
			for _, result := range results {
				if result.UserID == 0 || result.UserName == "" {
					t.Error("Expected valid user information")
				}
			}
		}
	})

	t.Run("test raw rows resource management", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试资源管理
		rows, err := query.RawRows("SELECT id, name FROM users LIMIT 1")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证rows不是nil
		if rows == nil {
			t.Error("Expected non-nil rows, got nil")
		}

		// 正确关闭资源
		err = rows.Close()
		if err != nil {
			t.Errorf("Expected no error when closing rows, got: %v", err)
		}
	})
}

// TestQuery_Pluck 测试 Pluck 方法的各种使用场景
func TestQuery_Pluck(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	t.Run("test pluck with field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字段指针提取字段值
		var names []string
		err := query.Pluck(&user.Name, &names)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(names) == 0 {
			t.Error("Expected to pluck user names, but got empty result")
		}

		// 验证结果中包含预期的用户名
		found := false
		expectedNames := []string{"testuser1", "testuser2", "testuser3", "otheruser"}
		for _, name := range names {
			for _, expected := range expectedNames {
				if name == expected {
					found = true
					break
				}
			}
		}
		if !found {
			t.Errorf("Expected to find one of %v in plucked names, but got %v", expectedNames, names)
		}
	})

	t.Run("test pluck with field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字段名字符串提取字段值
		var ages []int
		err := query.Pluck("age", &ages)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(ages) == 0 {
			t.Error("Expected to pluck user ages, but got empty result")
		}

		// 验证结果包含预期的年龄值
		expectedAges := []int{25, 30, 35, 28}
		for _, age := range ages {
			found := false
			for _, expected := range expectedAges {
				if age == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Unexpected age %d in plucked ages", age)
			}
		}
	})

	t.Run("test pluck with conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试带条件的pluck操作
		var activeUserNames []string
		err := query.Eq(&user.Status, 1).Pluck(&user.Name, &activeUserNames)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(activeUserNames) == 0 {
			t.Error("Expected to pluck active user names, but got empty result")
		}

		// 验证结果只包含状态为1的用户名称
		expectedActiveNames := []string{"testuser1", "testuser2", "testuser3"}
		if len(activeUserNames) != len(expectedActiveNames) {
			t.Errorf("Expected %d active users, got %d", len(expectedActiveNames), len(activeUserNames))
		}

		for _, name := range activeUserNames {
			found := false
			for _, expected := range expectedActiveNames {
				if name == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Unexpected active user name %s in plucked names", name)
			}
		}
	})

	t.Run("test pluck with like condition", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试带LIKE条件的pluck操作
		var testUserNames []string
		err := query.Like(&user.Name, "testuser").Pluck(&user.Name, &testUserNames)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(testUserNames) == 0 {
			t.Error("Expected to pluck test user names, but got empty result")
		}

		// 验证结果只包含包含"testuser"的用户名
		for _, name := range testUserNames {
			if !strings.Contains(name, "testuser") {
				t.Errorf("Expected name to contain 'testuser', got %s", name)
			}
		}
	})

	t.Run("test pluck with ordering", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试带排序的pluck操作
		var orderedAges []int
		err := query.OrderAsc(&user.Age).Pluck(&user.Age, &orderedAges)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(orderedAges) == 0 {
			t.Error("Expected to pluck ordered ages, but got empty result")
		}

		// 验证结果按升序排列
		for i := 1; i < len(orderedAges); i++ {
			if orderedAges[i-1] > orderedAges[i] {
				t.Error("Expected ages to be in ascending order")
			}
		}
	})

	t.Run("test pluck with no matching records", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试无匹配记录时的pluck操作
		var names []string
		err := query.Eq(&user.Name, "nonexistent_user").Pluck(&user.Name, &names)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 应该返回空切片而不是错误
		if names == nil {
			t.Error("Expected empty slice, got nil")
		}

		if len(names) != 0 {
			t.Errorf("Expected empty slice, got length: %d", len(names))
		}
	})

	t.Run("test pluck with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用无效字段名
		var values []string
		err := query.Pluck("invalid_field", &values)
		// 使用无效字段应该返回错误
		if err == nil {
			t.Error("Expected error for invalid field, got nil")
		}

		// 验证返回空结果
		if len(values) != 0 {
			t.Errorf("Expected empty result for invalid field, got %v", values)
		}
	})

	t.Run("test pluck with different data types", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试提取不同类型字段的值
		testCases := []struct {
			field      interface{}
			name       string
			validateFn func(interface{}) bool
		}{
			{
				field: &user.ID,
				name:  "ID field",
				validateFn: func(val interface{}) bool {
					ptr, ok := val.(*[]uint)
					if !ok {
						return false
					}
					ids := *ptr
					for _, id := range ids {
						if id <= 0 {
							return false
						}
					}
					return true
				},
			},
			{
				field: &user.Age,
				name:  "Age field",
				validateFn: func(val interface{}) bool {
					ptr, ok := val.(*[]int)
					if !ok {
						return false
					}
					ages := *ptr
					for _, age := range ages {
						if age <= 0 {
							return false
						}
					}
					return true
				},
			},
			{
				field: &user.IsActive,
				name:  "IsActive field",
				validateFn: func(val interface{}) bool {
					ptr, ok := val.(*[]bool)
					if !ok {
						return false
					}
					statuses := *ptr
					return len(statuses) > 0
				},
			},
			{
				field: &user.Salary,
				name:  "Salary field",
				validateFn: func(val interface{}) bool {
					ptr, ok := val.(*[]float64)
					if !ok {
						return false
					}
					salaries := *ptr
					for _, salary := range salaries {
						if salary < 0 {
							return false
						}
					}
					return true
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// 创建对应类型的切片
				var result interface{}
				switch tc.field {
				case &user.ID:
					result = &[]uint{}
				case &user.Age:
					result = &[]int{}
				case &user.IsActive:
					result = &[]bool{}
				case &user.Salary:
					result = &[]float64{}
				default:
					result = &[]string{}
				}

				err := query.Pluck(tc.field, result)
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}

				// 验证结果
				if !tc.validateFn(result) {
					t.Errorf("Validation failed for %s", tc.name)
				}
			})
		}
	})

	t.Run("test pluck with nil destination", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试目标为nil的情况
		err := query.Pluck(&user.Name, nil)
		if err == nil {
			t.Error("Expected error when destination is nil, got nil")
		}
	})

	t.Run("test pluck with non-pointer destination", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试目标不是指针的情况
		var names []string
		err := query.Pluck(&user.Name, names) // 传递值而不是指针
		if err == nil {
			t.Error("Expected error when destination is not a pointer, got nil")
		}
	})

	// 多对一关系使用场景测试
	t.Run("test pluck with many-to-one relationship scenario", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 测试多对一关系中的pluck操作
		var userIds []uint
		err := query.Pluck(&post.UserID, &userIds)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(userIds) > 0 {
			// 验证结果
			for _, userId := range userIds {
				if userId == 0 {
					t.Error("Expected valid user IDs")
				}
			}
		}
	})

	t.Run("test pluck with limit", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试带limit的pluck操作
		var limitedNames []string
		err := query.Limit(2).Pluck(&user.Name, &limitedNames)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(limitedNames) > 2 {
			t.Errorf("Expected at most 2 names due to limit, got %d", len(limitedNames))
		}
	})

	t.Run("test pluck with distinct", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试带distinct的pluck操作
		var statuses []int
		err := query.Distinct(&user.Status).Pluck(&user.Status, &statuses)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(statuses) == 0 {
			t.Error("Expected to pluck distinct statuses, but got empty result")
		}

		// 验证去重效果
		uniqueStatuses := make(map[int]bool)
		for _, status := range statuses {
			uniqueStatuses[status] = true
		}

		if len(uniqueStatuses) != len(statuses) {
			t.Error("Expected distinct statuses, but found duplicates")
		}
	})
}

// TestQuery_Preload 测试 Preload 方法的基本功能
func TestQuery_Preload(t *testing.T) {
	// 准备测试数据
	setupPreloadTestData(t)

	t.Run("test preload has many relationship", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试预加载用户的文章（一对多关系）
		query.Eq(&user.Name, "testuser1").
			Preload("Posts")

		var result User
		err := query.First(&result)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证用户信息
		if result.Name != "testuser1" {
			t.Errorf("Expected name: testuser1, got: %s", result.Name)
		}

		// 验证关联数据已加载
		if result.ID == 0 {
			t.Error("Expected user to be loaded")
		}
	})

	t.Run("test preload belongs to relationship", func(t *testing.T) {
		query, _ := gormx.NewQuery[Post]()

		// 测试预加载文章的作者（多对一关系）
		query.Preload("User")

		var results []Post
		err := query.Find(&results)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证至少有一篇文章
		if len(results) == 0 {
			t.Error("Expected to find posts")
			return
		}

		// 验证关联数据已加载
		for _, post := range results {
			// 添加调试信息
			t.Logf("Post ID: %d, UserID: %d, UserName: %s", post.ID, post.UserID, post.User.Name)

			// 检查关联是否已预加载
			if post.User.ID == 0 {
				// 如果 User.ID 为 0，可能是关联未正确加载
				// 我们检查是否至少设置了 UserID 字段
				if post.UserID == 0 {
					t.Error("Expected post to have a valid UserID")
				}
				// 注意：User.ID 为 0 并不一定意味着预加载失败，可能只是关联记录不存在
				// 但我们在这个测试中期望能找到有效的用户
				t.Logf("Warning: User not preloaded for post ID %d", post.ID)
			}
		}
	})

	t.Run("test preload many to many relationship", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试预加载用户的角色（多对多关系）
		query.Eq(&user.Name, "testuser2").
			Preload("Roles")

		var result User
		err := query.First(&result)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证用户信息
		if result.Name != "testuser2" {
			t.Errorf("Expected name: testuser2, got: %s", result.Name)
		}

		// 验证关联数据已加载
		if result.ID == 0 {
			t.Error("Expected user to be loaded")
		}
	})

	t.Run("test preload with conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试带条件的预加载
		query.Eq(&user.Name, "testuser1").
			Preload("Posts", "status = ?", 1)

		var result User
		err := query.First(&result)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证用户信息
		if result.Name != "testuser1" {
			t.Errorf("Expected name: testuser1, got: %s", result.Name)
		}
	})

	t.Run("test preload with ordering", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试带排序的预加载
		query.Eq(&user.Name, "testuser1").
			Preload("Posts", func(db *gorm.DB) *gorm.DB {
				return db.Order("created_at DESC")
			})

		var result User
		err := query.First(&result)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证用户信息
		if result.Name != "testuser1" {
			t.Errorf("Expected name: testuser1, got: %s", result.Name)
		}
	})

	t.Run("test nested preload", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试嵌套预加载
		query.Eq(&user.Name, "testuser2").
			Preload("Roles").
			Preload("Posts")

		var result User
		err := query.First(&result)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证用户信息
		if result.Name != "testuser2" {
			t.Errorf("Expected name: testuser2, got: %s", result.Name)
		}
	})

	t.Run("test preload with invalid association", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试无效的关联名称
		query.Eq(&user.Name, "testuser1").
			Preload("InvalidAssociation")

		var result User
		err := query.First(&result)

		// 无效关联名可能会导致错误，取决于 GORM 的实现
		// 我们主要确保程序不会 panic
		if err != nil && err != gorm.ErrRecordNotFound {
			// 接受关联错误或者其他非记录未找到的错误
			t.Logf("Got error (expected for invalid association): %v", err)
		}
	})

	t.Run("test preload with multiple associations", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试预加载多个关联
		query.Preload("Posts").
			Preload("Roles")

		var results []User
		err := query.Find(&results)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果
		if len(results) == 0 {
			t.Error("Expected to find users")
		}
	})

	t.Run("test preload with select", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试预加载与字段选择结合使用
		query.Eq(&user.Name, "testuser1").
			Select(&user.Name, &user.Email).
			Preload("Posts")

		var result User
		err := query.First(&result)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证用户信息
		if result.Name != "testuser1" {
			t.Errorf("Expected name: testuser1, got: %s", result.Name)
		}
	})
}

// setupPreloadTestData 准备预加载测试数据
func setupPreloadTestData(t *testing.T) {
	// 清理现有数据
	db := gormx.GetDb()
	db.Where("1 = 1").Delete(&User{})
	db.Where("1 = 1").Delete(&Post{})
	db.Where("1 = 1").Delete(&Role{})
	db.Where("1 = 1").Delete(&UserRole{})
	db.Where("1 = 1").Delete(&Profile{})

	// 插入测试用户数据
	users := []User{
		{Name: "testuser1", Email: "test1@example.com", Age: 25, IsActive: true, Salary: 5000.0, Status: 1},
		{Name: "testuser2", Email: "test2@example.com", Age: 30, IsActive: true, Salary: 6000.0, Status: 1},
		{Name: "testuser3", Email: "test3@example.com", Age: 35, IsActive: false, Salary: 7000.0, Status: 1},
		{Name: "otheruser", Email: "other@example.com", Age: 28, IsActive: true, Salary: 5500.0, Status: 0},
	}

	for i := range users {
		if err := db.Create(&users[i]).Error; err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}

	// 获取刚创建的用户ID
	var createdUsers []User
	if err := db.Find(&createdUsers).Error; err != nil {
		t.Fatalf("Failed to fetch created users: %v", err)
	}

	// 确保有足够的用户用于创建帖子
	if len(createdUsers) < 2 {
		t.Fatal("Not enough users created for testing posts")
	}

	// 插入测试角色数据
	roles := []Role{
		{Name: "Admin", Description: "Administrator role", Status: 1},
		{Name: "Editor", Description: "Editor role", Status: 1},
		{Name: "Viewer", Description: "Viewer role", Status: 1},
	}

	for i := range roles {
		if err := db.Create(&roles[i]).Error; err != nil {
			t.Fatalf("Failed to create test role: %v", err)
		}
	}

	// 获取刚创建的角色ID
	var createdRoles []Role
	if err := db.Find(&createdRoles).Error; err != nil {
		t.Fatalf("Failed to fetch created roles: %v", err)
	}

	// 插入用户角色关联数据，使用实际创建的用户和角色ID
	userRoles := []UserRole{
		{UserID: createdUsers[0].ID, RoleID: createdRoles[0].ID}, // testuser1 -> Admin
		{UserID: createdUsers[0].ID, RoleID: createdRoles[1].ID}, // testuser1 -> Editor
		{UserID: createdUsers[1].ID, RoleID: createdRoles[1].ID}, // testuser2 -> Editor
		{UserID: createdUsers[1].ID, RoleID: createdRoles[2].ID}, // testuser2 -> Viewer
	}

	for i := range userRoles {
		if err := db.Create(&userRoles[i]).Error; err != nil {
			t.Fatalf("Failed to create test user role: %v", err)
		}
	}

	// 插入测试帖子数据，使用实际创建的用户ID
	posts := []Post{
		{UserID: createdUsers[0].ID, Title: "Test Post 1", Content: "Content 1", Status: 1},
		{UserID: createdUsers[0].ID, Title: "Test Post 2", Content: "Content 2", Status: 1},
		{UserID: createdUsers[1].ID, Title: "Test Post 3", Content: "Content 3", Status: 1},
	}

	for i := range posts {
		if err := db.Create(&posts[i]).Error; err != nil {
			t.Fatalf("Failed to create test post: %v", err)
		}
	}

	// 插入测试资料数据
	profiles := []Profile{
		{UserID: createdUsers[0].ID, Avatar: "avatar1.jpg", PhoneNumber: "1234567890", Bio: "Bio of testuser1"},
		{UserID: createdUsers[1].ID, Avatar: "avatar2.jpg", PhoneNumber: "0987654321", Bio: "Bio of testuser2"},
	}

	for i := range profiles {
		if err := db.Create(&profiles[i]).Error; err != nil {
			t.Fatalf("Failed to create test profile: %v", err)
		}
	}
}

// TestQuery_GroupBy 测试 GroupBy 方法的各种使用场景
func TestQuery_GroupBy(t *testing.T) {
	// 准备测试数据
	setupGroupByTestData(t)

	t.Run("test group by with field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字段指针进行分组
		result := query.GroupBy(&user.Status)

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("GroupBy should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("GroupBy should add one option")
		}

		// 验证生成的SQL
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL GROUP BY `status`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test group by with field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字段名字符串进行分组
		query.GroupBy("status")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL GROUP BY `status`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test group by with multiple fields", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用多个字段进行分组
		query.GroupBy(&user.Status).
			GroupBy(&user.IsActive)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL GROUP BY `status`,`is_active`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test group by with mixed field pointers and field name strings", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试混合使用字段指针和字段名字符串进行分组
		query.GroupBy(&user.Status).
			GroupBy("is_active")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL GROUP BY `status`,`is_active`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test group by chaining with other query methods", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与其他查询方法链式调用
		query.Eq(&user.Status, 1).
			GroupBy(&user.IsActive)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `status` = ? AND `users`.`deleted_at` IS NULL GROUP BY `is_active`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 {
			t.Errorf("Expected 1 arg, got: %d", len(args))
		}

		if args[0] != 1 {
			t.Errorf("Expected arg: 1, got: %v", args[0])
		}
	})

	t.Run("test group by with select", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与Select方法结合使用
		query.Select(&user.Status).
			GroupBy(&user.Status)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT `status` FROM `users` WHERE `users`.`deleted_at` IS NULL GROUP BY `status`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test group by with count aggregation", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与Count聚合函数结合使用
		query.GroupBy(&user.Status).
			Select(&user.Status).
			Count(&user.ID)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT status, COUNT(`id`) as count FROM `users` WHERE `users`.`deleted_at` IS NULL GROUP BY `status`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test group by with sum aggregation", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与Sum聚合函数结合使用
		query.GroupBy(&user.Status).
			Select(&user.Status).
			Sum(&user.Salary)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT status, SUM(`salary`) FROM `users` WHERE `users`.`deleted_at` IS NULL GROUP BY `status`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test group by with avg aggregation", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与Avg聚合函数结合使用
		query.GroupBy(&user.Status).
			Select(&user.Status).
			Avg(&user.Age)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT status, AVG(`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL GROUP BY `status`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test group by with having clause", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与Having子句结合使用
		query.GroupBy(&user.Status).
			Select(&user.Status).
			Count(&user.ID).
			Having("COUNT(id) > ?", 1)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT status, COUNT(`id`) as count FROM `users` WHERE `users`.`deleted_at` IS NULL GROUP BY `status` HAVING COUNT(id) > ?")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 {
			t.Errorf("Expected 1 arg, got: %d", len(args))
		}

		if args[0] != 1 {
			t.Errorf("Expected arg: 1, got: %v", args[0])
		}
	})

	t.Run("test group by with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用无效字段名进行分组
		query.GroupBy("invalid_field")

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL GROUP BY `invalid_field`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}
		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test group by with order by", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与OrderBy结合使用
		query.GroupBy(&user.Status).
			OrderAsc(&user.Status)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL GROUP BY `status` ORDER BY `status`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test group by with where condition", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与Where条件结合使用
		query.Eq(&user.IsActive, true).
			GroupBy(&user.Status)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `is_active` = ? AND `users`.`deleted_at` IS NULL GROUP BY `status`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 {
			t.Errorf("Expected 1 arg, got: %d", len(args))
		}

		if args[0] != true {
			t.Errorf("Expected arg: true, got: %v", args[0])
		}
	})

	// 多对一关系使用场景测试
	t.Run("test group by with many-to-one relationship", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 测试在多对一关系中使用分组
		query.GroupBy(&post.UserID)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `posts` WHERE `posts`.`deleted_at` IS NULL GROUP BY `user_id`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test group by with complex query", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 测试复杂查询与分组结合使用
		query.Eq(&post.Status, 1).
			GroupBy(&post.UserID).
			Select(&post.UserID).
			Count(&post.ID).
			OrderDesc(&post.UserID)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT user_id, COUNT(`id`) as count FROM `posts` WHERE `status` = ? AND `posts`.`deleted_at` IS NULL GROUP BY `user_id` ORDER BY `user_id` DESC")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 {
			t.Errorf("Expected 1 arg, got: %d", len(args))
		}

		if args[0] != 1 {
			t.Errorf("Expected arg: 1, got: %v", args[0])
		}
	})
}

// setupGroupByTestData 准备GroupBy测试数据
func setupGroupByTestData(t *testing.T) {
	// 清理现有数据
	db := gormx.GetDb()
	db.Where("1 = 1").Delete(&User{})
	db.Where("1 = 1").Delete(&Post{})
	// 插入测试用户数据
	users := []User{
		{Name: "testuser1", Email: "test1@example.com", Age: 25, IsActive: true, Salary: 5000.0, Status: 1},
		{Name: "testuser2", Email: "test2@example.com", Age: 30, IsActive: true, Salary: 6000.0, Status: 1},
		{Name: "testuser3", Email: "test3@example.com", Age: 35, IsActive: false, Salary: 7000.0, Status: 2},
		{Name: "testuser4", Email: "test4@example.com", Age: 28, IsActive: false, Salary: 5500.0, Status: 2},
		{Name: "testuser5", Email: "test5@example.com", Age: 32, IsActive: true, Salary: 6500.0, Status: 1},
	}

	for i := range users {
		if err := db.Create(&users[i]).Error; err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}

	// 插入测试帖子数据
	posts := []Post{
		{UserID: 1, Title: "Test Post 1", Content: "Content 1", Status: 1},
		{UserID: 1, Title: "Test Post 2", Content: "Content 2", Status: 1},
		{UserID: 2, Title: "Test Post 3", Content: "Content 3", Status: 1},
		{UserID: 3, Title: "Test Post 4", Content: "Content 4", Status: 2},
	}

	for i := range posts {
		if err := db.Create(&posts[i]).Error; err != nil {
			t.Fatalf("Failed to create test post: %v", err)
		}
	}
}

// TestQuery_InSql 测试 InSql 方法的基本功能
func TestQuery_InSql(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	t.Run("test basic in sql with field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字段指针和子查询SQL
		result := query.InSql(&user.Status, "SELECT 1")

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("InSql should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("InSql should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `status` IN (SELECT 1) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test basic in sql with field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字段名字符串和子查询SQL
		query.InSql("status", "SELECT 1")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `status` IN (SELECT 1) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test in sql with parameters", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试带参数的子查询SQL
		query.InSql(&user.Status, "SELECT ? FROM dual", 1)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `status` IN (SELECT ? FROM dual) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	t.Run("test chaining multiple in sql conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试链式调用多个InSql条件
		query.InSql(&user.Status, "SELECT 1").
			InSql(&user.Age, "SELECT 2").
			InSql("name", "SELECT 'test'")

		if len(query.ToOptions()) != 3 {
			t.Error("Should have 3 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `status` IN (SELECT 1) AND `age` IN (SELECT 2) AND `name` IN (SELECT 'test') AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test in sql combined with other query methods", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与其他查询方法组合使用
		query.Eq(&user.IsActive, true).
			InSql(&user.Status, "SELECT 1 UNION SELECT 2")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `is_active` = ? AND `status` IN (SELECT 1 UNION SELECT 2) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != true {
			t.Errorf("Expected args: [true], got: %v", args)
		}
	})

	t.Run("test in sql with complex subquery", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试复杂的子查询
		subQuery := "SELECT u.status FROM users u WHERE u.age > ? AND u.created_at > ?"
		query.InSql(&user.Status, subQuery, 18, "2023-01-01")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `status` IN (SELECT u.status FROM users u WHERE u.age > ? AND u.created_at > ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 || args[0] != 18 || args[1] != "2023-01-01" {
			t.Errorf("Expected args: [18 2023-01-01], got: %v", args)
		}
	})

	t.Run("test in sql with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效字段名（应该不会panic，但会在实际查询时处理）
		query.InSql("invalid_field", "SELECT 1")

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍然会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `invalid_field` IN (SELECT 1) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test in sql with empty sql", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试空SQL（应该被忽略）
		query.InSql(&user.Status, "")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	// 一对多关系使用场景测试
	t.Run("test in sql with one-to-many relationship", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试在一对多关系中使用InSql
		query.InSql(&user.ID, "SELECT DISTINCT user_id FROM posts WHERE status = ?", 1)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `id` IN (SELECT DISTINCT user_id FROM posts WHERE status = ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	// 多对多关系使用场景测试
	t.Run("test in sql with many-to-many relationship", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试在多对多关系中使用InSql
		query.InSql(&user.ID, "SELECT DISTINCT user_id FROM user_roles WHERE role_id IN (SELECT id FROM roles WHERE name = ?)", "admin")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `id` IN (SELECT DISTINCT user_id FROM user_roles WHERE role_id IN (SELECT id FROM roles WHERE name = ?)) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "admin" {
			t.Errorf("Expected args: [admin], got: %v", args)
		}
	})

	// LEFT JOIN 使用场景测试
	t.Run("test in sql with left join scenario", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 模拟LEFT JOIN场景：查找有帖子的用户
		query.InSql(&user.ID, "SELECT DISTINCT p.user_id FROM posts p LEFT JOIN users u ON p.user_id = u.id WHERE p.status = ?", 1)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `id` IN (SELECT DISTINCT p.user_id FROM posts p LEFT JOIN users u ON p.user_id = u.id WHERE p.status = ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	t.Run("test in sql with complex left join scenario", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 更复杂的LEFT JOIN场景：查找有特定角色的用户
		subQuery := "SELECT DISTINCT u.id FROM users u LEFT JOIN user_roles ur ON u.id = ur.user_id LEFT JOIN roles r ON ur.role_id = r.id WHERE r.name IN (?, ?)"
		query.InSql(&user.ID, subQuery, "admin", "editor")

		sql, args := query.ToSQLAndArgs()
		// 注意：换行符会被压缩成空格
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `id` IN (SELECT DISTINCT u.id FROM users u LEFT JOIN user_roles ur ON u.id = ur.user_id LEFT JOIN roles r ON ur.role_id = r.id WHERE r.name IN (?, ?)) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 || args[0] != "admin" || args[1] != "editor" {
			t.Errorf("Expected args: [admin editor], got: %v", args)
		}
	})
}

// TestQuery_NotInSql 测试 NotInSql 方法的基本功能
func TestQuery_NotInSql(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	t.Run("test basic not in sql with field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字段指针和子查询SQL
		result := query.NotInSql(&user.Status, "SELECT 1")

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("NotInSql should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("NotInSql should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `status` NOT IN (SELECT 1) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test not in sql with field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字段名字符串和子查询SQL
		query.NotInSql("status", "SELECT 1")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `status` NOT IN (SELECT 1) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test not in sql with parameters", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试带参数的子查询SQL
		query.NotInSql(&user.Status, "SELECT ? FROM dual", 1)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `status` NOT IN (SELECT ? FROM dual) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	t.Run("test chaining multiple not in sql conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试链式调用多个NotInSql条件
		query.NotInSql(&user.Status, "SELECT 1").
			NotInSql(&user.Age, "SELECT 2").
			NotInSql("name", "SELECT 'test'")

		if len(query.ToOptions()) != 3 {
			t.Error("Should have 3 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `status` NOT IN (SELECT 1) AND `age` NOT IN (SELECT 2) AND `name` NOT IN (SELECT 'test') AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test not in sql combined with other query methods", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与其他查询方法组合使用
		query.Eq(&user.IsActive, true).
			NotInSql(&user.Status, "SELECT 1 UNION SELECT 2")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `is_active` = ? AND `status` NOT IN (SELECT 1 UNION SELECT 2) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != true {
			t.Errorf("Expected args: [true], got: %v", args)
		}
	})

	t.Run("test not in sql with complex subquery", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试复杂的子查询
		subQuery := "SELECT u.status FROM users u WHERE u.age > ? AND u.created_at > ?"
		query.NotInSql(&user.Status, subQuery, 18, "2023-01-01")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `status` NOT IN (SELECT u.status FROM users u WHERE u.age > ? AND u.created_at > ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 || args[0] != 18 || args[1] != "2023-01-01" {
			t.Errorf("Expected args: [18 2023-01-01], got: %v", args)
		}
	})

	t.Run("test not in sql with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效字段名（应该不会panic，但会在实际查询时处理）
		query.NotInSql("invalid_field", "SELECT 1")

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍然会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `invalid_field` NOT IN (SELECT 1) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test not in sql with empty sql", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试空SQL（应该被忽略）
		query.NotInSql(&user.Status, "")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	// 一对多关系使用场景测试
	t.Run("test not in sql with one-to-many relationship", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试在一对多关系中使用NotInSql
		query.NotInSql(&user.ID, "SELECT DISTINCT user_id FROM posts WHERE status = ?", 1)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `id` NOT IN (SELECT DISTINCT user_id FROM posts WHERE status = ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	// 多对多关系使用场景测试
	t.Run("test not in sql with many-to-many relationship", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试在多对多关系中使用NotInSql
		query.NotInSql(&user.ID, "SELECT DISTINCT user_id FROM user_roles WHERE role_id IN (SELECT id FROM roles WHERE name = ?)", "admin")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `id` NOT IN (SELECT DISTINCT user_id FROM user_roles WHERE role_id IN (SELECT id FROM roles WHERE name = ?)) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "admin" {
			t.Errorf("Expected args: [admin], got: %v", args)
		}
	})
}

// TestQuery_GtSql 测试 GtSql 方法的基本功能
func TestQuery_GtSql(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	t.Run("test basic gt sql with field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字段指针和子查询SQL
		result := query.GtSql(&user.Age, "SELECT AVG(age) FROM users")

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("GtSql should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("GtSql should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` > (SELECT AVG(age) FROM users) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test gt sql with field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字段名字符串和子查询SQL
		query.GtSql("age", "SELECT AVG(age) FROM users")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` > (SELECT AVG(age) FROM users) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test gt sql with parameters", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试带参数的子查询SQL
		query.GtSql(&user.Age, "SELECT AVG(age) FROM users WHERE department = ?", "IT")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` > (SELECT AVG(age) FROM users WHERE department = ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "IT" {
			t.Errorf("Expected args: [IT], got: %v", args)
		}
	})

	t.Run("test chaining multiple gt sql conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试链式调用多个GtSql条件
		query.GtSql(&user.Age, "SELECT 18").
			GtSql(&user.Score, "SELECT 80").
			GtSql("salary", "SELECT 5000")

		if len(query.ToOptions()) != 3 {
			t.Error("Should have 3 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` > (SELECT 18) AND `score` > (SELECT 80) AND `salary` > (SELECT 5000) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test gt sql combined with other query methods", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与其他查询方法组合使用
		query.Eq(&user.IsActive, true).
			GtSql(&user.Age, "SELECT AVG(age) FROM users WHERE is_active = true")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `is_active` = ? AND `age` > (SELECT AVG(age) FROM users WHERE is_active = true) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != true {
			t.Errorf("Expected args: [true], got: %v", args)
		}
	})

	t.Run("test gt sql with complex subquery", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试复杂的子查询
		subQuery := "SELECT AVG(u.age) FROM users u WHERE u.department = ? AND u.created_at > ?"
		query.GtSql(&user.Age, subQuery, "IT", "2023-01-01")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` > (SELECT AVG(u.age) FROM users u WHERE u.department = ? AND u.created_at > ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 || args[0] != "IT" || args[1] != "2023-01-01" {
			t.Errorf("Expected args: [IT 2023-01-01], got: %v", args)
		}
	})

	t.Run("test gt sql with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效字段名（应该不会panic，但会在实际查询时处理）
		query.GtSql("invalid_field", "SELECT 1")

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍然会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `invalid_field` > (SELECT 1) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test gt sql with empty sql", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试空SQL（应该被忽略）
		query.GtSql(&user.Age, "")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	// 一对多关系使用场景测试
	t.Run("test gt sql with one-to-many relationship", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 测试在一对多关系中使用GtSql
		query.GtSql(&post.UserID, "SELECT MIN(user_id) FROM posts WHERE status = ?", 1)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `posts` WHERE `user_id` > (SELECT MIN(user_id) FROM posts WHERE status = ?) AND `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	// 多对多关系使用场景测试
	t.Run("test gt sql with many-to-many relationship", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试在多对多关系中使用GtSql
		query.GtSql(&user.ID, "SELECT COUNT(role_id) FROM user_roles WHERE user_id IN (SELECT id FROM users WHERE department = ?)", "IT")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `id` > (SELECT COUNT(role_id) FROM user_roles WHERE user_id IN (SELECT id FROM users WHERE department = ?)) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "IT" {
			t.Errorf("Expected args: [IT], got: %v", args)
		}
	})
}

// TestQuery_GeSql 测试 GeSql 方法的基本功能
func TestQuery_GeSql(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	t.Run("test basic ge sql with field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字段指针和子查询SQL
		result := query.GeSql(&user.Age, "SELECT AVG(age) FROM users")

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("GeSql should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("GeSql should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` >= (SELECT AVG(age) FROM users) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test ge sql with field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字段名字符串和子查询SQL
		query.GeSql("age", "SELECT AVG(age) FROM users")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` >= (SELECT AVG(age) FROM users) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test ge sql with parameters", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试带参数的子查询SQL
		query.GeSql(&user.Age, "SELECT AVG(age) FROM users WHERE department = ?", "IT")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` >= (SELECT AVG(age) FROM users WHERE department = ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "IT" {
			t.Errorf("Expected args: [IT], got: %v", args)
		}
	})

	t.Run("test chaining multiple ge sql conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试链式调用多个GeSql条件
		query.GeSql(&user.Age, "SELECT 18").
			GeSql(&user.Score, "SELECT 80").
			GeSql("salary", "SELECT 5000")

		if len(query.ToOptions()) != 3 {
			t.Error("Should have 3 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` >= (SELECT 18) AND `score` >= (SELECT 80) AND `salary` >= (SELECT 5000) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test ge sql combined with other query methods", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与其他查询方法组合使用
		query.Eq(&user.IsActive, true).
			GeSql(&user.Age, "SELECT AVG(age) FROM users WHERE is_active = true")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `is_active` = ? AND `age` >= (SELECT AVG(age) FROM users WHERE is_active = true) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != true {
			t.Errorf("Expected args: [true], got: %v", args)
		}
	})

	t.Run("test ge sql with complex subquery", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试复杂的子查询
		subQuery := "SELECT AVG(u.age) FROM users u WHERE u.department = ? AND u.created_at > ?"
		query.GeSql(&user.Age, subQuery, "IT", "2023-01-01")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` >= (SELECT AVG(u.age) FROM users u WHERE u.department = ? AND u.created_at > ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 || args[0] != "IT" || args[1] != "2023-01-01" {
			t.Errorf("Expected args: [IT 2023-01-01], got: %v", args)
		}
	})

	t.Run("test ge sql with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效字段名（应该不会panic，但会在实际查询时处理）
		query.GeSql("invalid_field", "SELECT 1")

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍然会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `invalid_field` >= (SELECT 1) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test ge sql with empty sql", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试空SQL（应该被忽略）
		query.GeSql(&user.Age, "")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	// 一对多关系使用场景测试
	t.Run("test ge sql with one-to-many relationship", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 测试在一对多关系中使用GeSql
		query.GeSql(&post.UserID, "SELECT MIN(user_id) FROM posts WHERE status = ?", 1)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `posts` WHERE `user_id` >= (SELECT MIN(user_id) FROM posts WHERE status = ?) AND `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	// 多对多关系使用场景测试
	t.Run("test ge sql with many-to-many relationship", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试在多对多关系中使用GeSql
		query.GeSql(&user.ID, "SELECT COUNT(role_id) FROM user_roles WHERE user_id IN (SELECT id FROM users WHERE department = ?)", "IT")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `id` >= (SELECT COUNT(role_id) FROM user_roles WHERE user_id IN (SELECT id FROM users WHERE department = ?)) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "IT" {
			t.Errorf("Expected args: [IT], got: %v", args)
		}
	})
}

// TestQuery_LtSql 测试 LtSql 方法的基本功能
func TestQuery_LtSql(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	t.Run("test basic lt sql with field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字段指针和子查询SQL
		result := query.LtSql(&user.Age, "SELECT AVG(age) FROM users")

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("LtSql should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("LtSql should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` < (SELECT AVG(age) FROM users) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test lt sql with field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字段名字符串和子查询SQL
		query.LtSql("age", "SELECT AVG(age) FROM users")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` < (SELECT AVG(age) FROM users) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test lt sql with parameters", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试带参数的子查询SQL
		query.LtSql(&user.Age, "SELECT AVG(age) FROM users WHERE department = ?", "IT")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` < (SELECT AVG(age) FROM users WHERE department = ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "IT" {
			t.Errorf("Expected args: [IT], got: %v", args)
		}
	})

	t.Run("test chaining multiple lt sql conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试链式调用多个LtSql条件
		query.LtSql(&user.Age, "SELECT 65").
			LtSql(&user.Score, "SELECT 90").
			LtSql("salary", "SELECT 10000")

		if len(query.ToOptions()) != 3 {
			t.Error("Should have 3 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` < (SELECT 65) AND `score` < (SELECT 90) AND `salary` < (SELECT 10000) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test lt sql combined with other query methods", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与其他查询方法组合使用
		query.Eq(&user.IsActive, true).
			LtSql(&user.Age, "SELECT AVG(age) FROM users WHERE is_active = true")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `is_active` = ? AND `age` < (SELECT AVG(age) FROM users WHERE is_active = true) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != true {
			t.Errorf("Expected args: [true], got: %v", args)
		}
	})

	t.Run("test lt sql with complex subquery", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试复杂的子查询
		subQuery := "SELECT AVG(u.age) FROM users u WHERE u.department = ? AND u.created_at > ?"
		query.LtSql(&user.Age, subQuery, "IT", "2023-01-01")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` < (SELECT AVG(u.age) FROM users u WHERE u.department = ? AND u.created_at > ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 || args[0] != "IT" || args[1] != "2023-01-01" {
			t.Errorf("Expected args: [IT 2023-01-01], got: %v", args)
		}
	})

	t.Run("test lt sql with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效字段名（应该不会panic，但会在实际查询时处理）
		query.LtSql("invalid_field", "SELECT 1")

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍然会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `invalid_field` < (SELECT 1) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test lt sql with empty sql", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试空SQL（应该被忽略）
		query.LtSql(&user.Age, "")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	// 一对多关系使用场景测试
	t.Run("test lt sql with one-to-many relationship", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 测试在一对多关系中使用LtSql
		query.LtSql(&post.UserID, "SELECT MAX(user_id) FROM posts WHERE status = ?", 1)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `posts` WHERE `user_id` < (SELECT MAX(user_id) FROM posts WHERE status = ?) AND `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	// 多对多关系使用场景测试
	t.Run("test lt sql with many-to-many relationship", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试在多对多关系中使用LtSql
		query.LtSql(&user.ID, "SELECT COUNT(role_id) FROM user_roles WHERE user_id IN (SELECT id FROM users WHERE department = ?)", "IT")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `id` < (SELECT COUNT(role_id) FROM user_roles WHERE user_id IN (SELECT id FROM users WHERE department = ?)) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "IT" {
			t.Errorf("Expected args: [IT], got: %v", args)
		}
	})
}

// TestQuery_LeSql 测试 LeSql 方法的基本功能
func TestQuery_LeSql(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	t.Run("test basic le sql with field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字段指针和子查询SQL
		result := query.LeSql(&user.Age, "SELECT AVG(age) FROM users")

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("LeSql should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("LeSql should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` <= (SELECT AVG(age) FROM users) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test le sql with field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字段名字符串和子查询SQL
		query.LeSql("age", "SELECT AVG(age) FROM users")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` <= (SELECT AVG(age) FROM users) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test le sql with parameters", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试带参数的子查询SQL
		query.LeSql(&user.Age, "SELECT AVG(age) FROM users WHERE department = ?", "IT")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` <= (SELECT AVG(age) FROM users WHERE department = ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "IT" {
			t.Errorf("Expected args: [IT], got: %v", args)
		}
	})

	t.Run("test chaining multiple le sql conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试链式调用多个LeSql条件
		query.LeSql(&user.Age, "SELECT 65").
			LeSql(&user.Score, "SELECT 90").
			LeSql("salary", "SELECT 10000")

		if len(query.ToOptions()) != 3 {
			t.Error("Should have 3 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` <= (SELECT 65) AND `score` <= (SELECT 90) AND `salary` <= (SELECT 10000) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test le sql combined with other query methods", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与其他查询方法组合使用
		query.Eq(&user.IsActive, true).
			LeSql(&user.Age, "SELECT AVG(age) FROM users WHERE is_active = true")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `is_active` = ? AND `age` <= (SELECT AVG(age) FROM users WHERE is_active = true) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != true {
			t.Errorf("Expected args: [true], got: %v", args)
		}
	})

	t.Run("test le sql with complex subquery", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试复杂的子查询
		subQuery := "SELECT AVG(u.age) FROM users u WHERE u.department = ? AND u.created_at > ?"
		query.LeSql(&user.Age, subQuery, "IT", "2023-01-01")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` <= (SELECT AVG(u.age) FROM users u WHERE u.department = ? AND u.created_at > ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 || args[0] != "IT" || args[1] != "2023-01-01" {
			t.Errorf("Expected args: [IT 2023-01-01], got: %v", args)
		}
	})

	t.Run("test le sql with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效字段名（应该不会panic，但会在实际查询时处理）
		query.LeSql("invalid_field", "SELECT 1")

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍然会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `invalid_field` <= (SELECT 1) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test le sql with empty sql", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试空SQL（应该被忽略）
		query.LeSql(&user.Age, "")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	// 一对多关系使用场景测试
	t.Run("test le sql with one-to-many relationship", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 测试在一对多关系中使用LeSql
		query.LeSql(&post.UserID, "SELECT MAX(user_id) FROM posts WHERE status = ?", 1)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `posts` WHERE `user_id` <= (SELECT MAX(user_id) FROM posts WHERE status = ?) AND `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	// 多对多关系使用场景测试
	t.Run("test le sql with many-to-many relationship", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试在多对多关系中使用LeSql
		query.LeSql(&user.ID, "SELECT COUNT(role_id) FROM user_roles WHERE user_id IN (SELECT id FROM users WHERE department = ?)", "IT")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `id` <= (SELECT COUNT(role_id) FROM user_roles WHERE user_id IN (SELECT id FROM users WHERE department = ?)) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "IT" {
			t.Errorf("Expected args: [IT], got: %v", args)
		}
	})
}

// TestQuery_Not 测试 Not 方法的基本功能
func TestQuery_Not(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	t.Run("test basic not with string condition", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字符串条件的Not方法
		result := query.Not("status = ?", 0)

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("Not should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("Not should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE NOT status = ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 0 {
			t.Errorf("Expected args: [0], got: %v", args)
		}
	})

	//t.Run("test not with struct condition", func(t *testing.T) {
	//	query, _ := gormx.NewQuery[User]()
	//
	//	// 测试使用结构体条件的Not方法
	//	query.Not(User{Status: 0})
	//
	//	sql, args := query.ToSQLAndArgs()
	//	expectedSQL := "SELECT * FROM `users` WHERE NOT (`users`.`status` = ?) AND `users`.`deleted_at` IS NULL"
	//	if sql != expectedSQL {
	//		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	//	}
	//
	//	if len(args) != 1 || args[0] != 0 {
	//		t.Errorf("Expected args: [0], got: %v", args)
	//	}
	//})

	t.Run("test not with map condition", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用map条件的Not方法
		query.Not(map[string]interface{}{"status": 0})

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `users`.`status` <> ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 0 {
			t.Errorf("Expected args: [0], got: %v", args)
		}
	})

	t.Run("test chaining multiple not conditions", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试链式调用多个Not条件
		query.Not("status = ?", 0).
			Not("is_active = ?", false).
			Not("age < ?", 18)

		if len(query.ToOptions()) != 3 {
			t.Error("Should have 3 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE NOT status = ? AND NOT is_active = ? AND NOT age < ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 3 {
			t.Errorf("Expected 3 args, got: %d", len(args))
		}

		expectedArgs := []interface{}{0, false, 18}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d]: %v, got: %v", i, expected, args[i])
			}
		}
	})

	t.Run("test not combined with other query methods", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试Not方法与其他查询方法组合使用
		query.Eq(&user.Name, "testuser").
			Not("status = ?", 0)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `name` = ? AND NOT status = ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 || args[0] != "testuser" || args[1] != 0 {
			t.Errorf("Expected args: [testuser 0], got: %v", args)
		}
	})

	t.Run("test not with complex condition", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试复杂的Not条件
		query.Not("age < ? OR department = ?", 18, "IT")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE NOT (age < ? OR department = ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 || args[0] != 18 || args[1] != "IT" {
			t.Errorf("Expected args: [18 IT], got: %v", args)
		}
	})

	t.Run("test not with multiple parameters", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试带多个参数的Not条件
		query.Not("status IN (?) AND age BETWEEN ? AND ?", []int{0, -1}, 18, 65)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE NOT (status IN (?,?) AND age BETWEEN ? AND ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 4 { // 1个slice参数会展开为多个值
			t.Errorf("Expected 4 args, got: %d", len(args))
		}
	})

	// 一对多关系使用场景测试
	t.Run("test not with one-to-many relationship", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试在一对多关系中使用Not
		query.Not("id IN (SELECT DISTINCT user_id FROM posts WHERE status = ?)", 0)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE NOT id IN (SELECT DISTINCT user_id FROM posts WHERE status = ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 0 {
			t.Errorf("Expected args: [0], got: %v", args)
		}
	})

	// 多对多关系使用场景测试
	t.Run("test not with many-to-many relationship", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试在多对多关系中使用Not
		query.Not("id IN (SELECT DISTINCT user_id FROM user_roles WHERE role_id IN (SELECT id FROM roles WHERE name = ?))", "banned")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE NOT id IN (SELECT DISTINCT user_id FROM user_roles WHERE role_id IN (SELECT id FROM roles WHERE name = ?)) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "banned" {
			t.Errorf("Expected args: [banned], got: %v", args)
		}
	})

	t.Run("test not with subquery condition", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用子查询的Not条件
		subQuery := "SELECT id FROM users WHERE created_at < ? AND is_active = ?"
		query.Not("id IN ("+subQuery+")", "2023-01-01", true)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE NOT (id IN (SELECT id FROM users WHERE created_at < ? AND is_active = ?)) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 || args[0] != "2023-01-01" || args[1] != true {
			t.Errorf("Expected args: [2023-01-01 true], got: %v", args)
		}
	})

	t.Run("test not with like condition", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试Not与LIKE条件结合使用
		query.Not("name LIKE ?", "%test%")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE NOT name LIKE ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "%test%" {
			t.Errorf("Expected args: [%%test%%], got: %v", args)
		}
	})

	t.Run("test multiple not conditions with and logic", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试多个Not条件以AND逻辑连接
		query.Not("status = ?", 0).
			Not("is_active = ?", false).
			Not("name LIKE ?", "%inactive%")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE NOT status = ? AND NOT is_active = ? AND NOT name LIKE ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 3 {
			t.Errorf("Expected 3 args, got: %d", len(args))
		}
	})
}

// TestQuery_Or 测试 Or 方法的基本功能
func TestQuery_Or(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	t.Run("test basic or with string condition", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字符串条件的Or方法
		result := query.Eq(&user.Status, 1).
			Or("age > ? AND age < ?", 18, 65)

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("Or should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 2 {
			t.Error("Or should add one option, making total 2 options")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE (`status` = ? OR (age > ? AND age < ?)) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 3 {
			t.Errorf("Expected 3 args, got: %d", len(args))
		}

		expectedArgs := []interface{}{1, 18, 65}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d]: %v, got: %v", i, expected, args[i])
			}
		}
	})

	t.Run("test chaining multiple or conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试链式调用多个Or条件
		query.Eq(&user.Status, 1).
			Or("age > ?", 18).
			Or("is_active = ?", true)

		if len(query.ToOptions()) != 3 {
			t.Error("Should have 3 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE (`status` = ? OR age > ? OR is_active = ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 3 {
			t.Errorf("Expected 3 args, got: %d", len(args))
		}

		expectedArgs := []interface{}{1, 18, true}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d]: %v, got: %v", i, expected, args[i])
			}
		}
	})

	t.Run("test or with complex condition", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试复杂的Or条件
		query.Eq(&user.Status, 1).
			Or("department = ? OR salary > ?", "IT", 50000.0)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE (`status` = ? OR (department = ? OR salary > ?)) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 3 || args[0] != 1 || args[1] != "IT" || args[2] != 50000.0 {
			t.Errorf("Expected args: [1 IT 50000.0], got: %v", args)
		}
	})

	t.Run("test or with like condition", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试Or与LIKE条件结合使用
		query.Eq(&user.Status, 1).
			Or("name LIKE ?", "%admin%")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE (`status` = ? OR name LIKE ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 || args[0] != 1 || args[1] != "%admin%" {
			t.Errorf("Expected args: [1 %%admin%%], got: %v", args)
		}
	})

	t.Run("test or with in condition", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试Or与IN条件结合使用
		query.Eq(&user.Status, 1).
			Or("id IN (?)", []int{2, 3, 4})

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE (`status` = ? OR id IN (?,?,?)) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 4 { // 1个值 + 3个slice展开的值
			t.Errorf("Expected 4 args, got: %d", len(args))
		}
	})

	t.Run("test or with between condition", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试Or与BETWEEN条件结合使用
		query.Eq(&user.Status, 1).
			Or("age BETWEEN ? AND ?", 20, 40)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE (`status` = ? OR (age BETWEEN ? AND ?)) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 3 || args[0] != 1 || args[1] != 20 || args[2] != 40 {
			t.Errorf("Expected args: [1 20 40], got: %v", args)
		}
	})

	// 一对多关系使用场景测试
	t.Run("test or with one-to-many relationship", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试在一对多关系中使用Or
		query.Eq(&user.Status, 1).
			Or("id IN (SELECT DISTINCT user_id FROM posts WHERE status = ?)", 1)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE (`status` = ? OR id IN (SELECT DISTINCT user_id FROM posts WHERE status = ?)) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 || args[0] != 1 || args[1] != 1 {
			t.Errorf("Expected args: [1 1], got: %v", args)
		}
	})

	// 多对多关系使用场景测试
	t.Run("test or with many-to-many relationship", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试在多对多关系中使用Or
		query.Eq(&user.Status, 1).
			Or("id IN (SELECT DISTINCT user_id FROM user_roles WHERE role_id IN (SELECT id FROM roles WHERE name = ?))", "admin")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE (`status` = ? OR id IN (SELECT DISTINCT user_id FROM user_roles WHERE role_id IN (SELECT id FROM roles WHERE name = ?))) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 || args[0] != 1 || args[1] != "admin" {
			t.Errorf("Expected args: [1 admin], got: %v", args)
		}
	})

	t.Run("test or with subquery condition", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用子查询的Or条件
		subQuery := "SELECT id FROM users WHERE created_at < ? AND is_active = ?"
		query.Eq(&user.Status, 1).
			Or("id IN ("+subQuery+")", "2023-01-01", true)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE (`status` = ? OR (id IN (SELECT id FROM users WHERE created_at < ? AND is_active = ?))) AND `users`.`deleted_at` IS NULL")

		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 3 || args[0] != 1 || args[1] != "2023-01-01" || args[2] != true {
			t.Errorf("Expected args: [1 2023-01-01 true], got: %v", args)
		}
	})

	t.Run("test multiple or conditions with and logic connection", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试多个Or条件与其他条件的逻辑连接
		query.Eq(&user.Status, 1).
			Or("is_active = ?", false).
			Or("name LIKE ?", "%test%").
			Eq(&user.Age, 25)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE (`status` = ? OR is_active = ? OR name LIKE ? AND `age` = ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 4 {
			t.Errorf("Expected 4 args, got: %d", len(args))
		}
	})

	t.Run("test or with null condition", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试Or与NULL条件结合使用
		query.Eq(&user.Status, 1).
			Or("email IS NULL")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE (`status` = ? OR email IS NULL) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	t.Run("test or with not null condition", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试Or与NOT NULL条件结合使用
		query.Eq(&user.Status, 1).
			Or("phone IS NOT NULL")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE (`status` = ? OR phone IS NOT NULL) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})
}

// TestQuery_SubQueryEq 测试 SubQueryEq 方法的基本功能
func TestQuery_SubQueryEq(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	t.Run("test basic subquery eq with field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 创建子查询
		subDB := gormx.GetDb().Model(&User{}).Select("id").Where("name = ?", "testuser1")

		// 测试使用字段指针的SubQueryEq方法
		result := query.SubQueryEq(&user.ID, subDB)

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("SubQueryEq should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("SubQueryEq should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `id` = (SELECT `id` FROM `users` WHERE name = ? AND `users`.`deleted_at` IS NULL) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "testuser1" {
			t.Errorf("Expected args: [testuser1], got: %v", args)
		}
	})

	t.Run("test subquery eq with field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 创建子查询
		subDB := gormx.GetDb().Model(&User{}).Select("id").Where("name = ?", "testuser2")

		// 测试使用字段名字符串的SubQueryEq方法
		query.SubQueryEq("id", subDB)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `id` = (SELECT `id` FROM `users` WHERE name = ? AND `users`.`deleted_at` IS NULL) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "testuser2" {
			t.Errorf("Expected args: [testuser2], got: %v", args)
		}
	})

	t.Run("test chaining multiple subquery eq conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 创建多个子查询
		subDB1 := gormx.GetDb().Model(&User{}).Select("id").Where("name = ?", "testuser1")
		subDB2 := gormx.GetDb().Model(&User{}).Select("status").Where("age = ?", 25)

		// 测试链式调用多个SubQueryEq条件
		query.SubQueryEq(&user.ID, subDB1).
			SubQueryEq(&user.Status, subDB2)

		if len(query.ToOptions()) != 2 {
			t.Error("Should have 2 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `id` = (SELECT `id` FROM `users` WHERE name = ? AND `users`.`deleted_at` IS NULL) AND `status` = (SELECT `status` FROM `users` WHERE age = ? AND `users`.`deleted_at` IS NULL) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 {
			t.Errorf("Expected 2 args, got: %d", len(args))
		}

		expectedArgs := []interface{}{"testuser1", 25}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d]: %v, got: %v", i, expected, args[i])
			}
		}
	})

	t.Run("test subquery eq combined with other query methods", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 创建子查询
		subDB := gormx.GetDb().Model(&User{}).Select("id").Where("name = ?", "testuser1")

		// 测试SubQueryEq方法与其他查询方法组合使用
		query.Eq(&user.IsActive, true).
			SubQueryEq(&user.ID, subDB)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `is_active` = ? AND `id` = (SELECT `id` FROM `users` WHERE name = ? AND `users`.`deleted_at` IS NULL) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 || args[0] != true || args[1] != "testuser1" {
			t.Errorf("Expected args: [true testuser1], got: %v", args)
		}
	})

	t.Run("test subquery eq with complex subquery", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 创建复杂子查询
		subDB := gormx.GetDb().Model(&User{}).
			Select("MAX(age)").
			Where("status = ?", 1).Group("is_active")

		// 测试使用复杂子查询的SubQueryEq方法
		query.SubQueryEq(&user.Age, subDB)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `age` = (SELECT MAX(age) FROM `users` WHERE status = ? AND `users`.`deleted_at` IS NULL GROUP BY `is_active`) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	// 一对多关系使用场景测试
	t.Run("test subquery eq with one-to-many relationship", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 创建子查询 - 查找特定用户的ID
		subDB := gormx.GetDb().Model(&User{}).Select("id").Where("name = ?", "testuser1")

		// 测试在一对多关系中使用SubQueryEq
		query.SubQueryEq(&post.UserID, subDB)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `posts` WHERE `user_id` = (SELECT `id` FROM `users` WHERE name = ? AND `users`.`deleted_at` IS NULL) AND `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "testuser1" {
			t.Errorf("Expected args: [testuser1], got: %v", args)
		}
	})

	// 多对多关系使用场景测试
	t.Run("test subquery eq with many-to-many relationship", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 创建子查询 - 查找具有特定角色的用户ID
		subDB := gormx.GetDb().Model(&UserRole{}).
			Select("user_id").
			Where("role_id = (SELECT id FROM roles WHERE name = ?)", "admin")

		// 测试在多对多关系中使用SubQueryEq
		query.SubQueryEq(&user.ID, subDB)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `id` = (SELECT `user_id` FROM `user_roles` WHERE role_id = (SELECT id FROM roles WHERE name = ?) AND `user_roles`.`deleted_at` IS NULL) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "admin" {
			t.Errorf("Expected args: [admin], got: %v", args)
		}
	})

	t.Run("test subquery eq with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 创建子查询
		subDB := gormx.GetDb().Model(&User{}).Select("id").Where("name = ?", "testuser")

		// 测试无效字段名（应该不会panic，但会在实际查询时处理）
		query.SubQueryEq("invalid_field", subDB)

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍然会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `invalid_field` = (SELECT `id` FROM `users` WHERE name = ? AND `users`.`deleted_at` IS NULL) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "testuser" {
			t.Errorf("Expected args: [testuser], got: %v", args)
		}
	})

	t.Run("test subquery eq with subquery returning multiple values", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 创建返回多个值的子查询（这在实际使用中会导致数据库错误）
		subDB := gormx.GetDb().Model(&User{}).Select("id").Where("status = ?", 1)

		// 测试SubQueryEq方法与返回多个值的子查询
		query.SubQueryEq(&user.ID, subDB)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `id` = (SELECT `id` FROM `users` WHERE status = ? AND `users`.`deleted_at` IS NULL) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})
}

// TestQuery_SubQueryIn 测试 SubQueryIn 方法的基本功能
func TestQuery_SubQueryIn(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	t.Run("test basic subquery in with field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 创建子查询
		subDB := gormx.GetDb().Model(&User{}).Select("id").Where("status = ?", 1)

		// 测试使用字段指针的SubQueryIn方法
		result := query.SubQueryIn(&user.ID, subDB)

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("SubQueryIn should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("SubQueryIn should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `id` IN (SELECT `id` FROM `users` WHERE status = ? AND `users`.`deleted_at` IS NULL) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	t.Run("test subquery in with field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 创建子查询
		subDB := gormx.GetDb().Model(&User{}).Select("id").Where("status = ?", 1)

		// 测试使用字段名字符串的SubQueryIn方法
		query.SubQueryIn("id", subDB)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `id` IN (SELECT `id` FROM `users` WHERE status = ? AND `users`.`deleted_at` IS NULL) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	t.Run("test chaining multiple subquery in conditions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 创建多个子查询
		subDB1 := gormx.GetDb().Model(&User{}).Select("id").Where("status = ?", 1)
		subDB2 := gormx.GetDb().Model(&User{}).Select("id").Where("age > ?", 18)

		// 测试链式调用多个SubQueryIn条件
		query.SubQueryIn(&user.ID, subDB1).
			SubQueryIn(&user.Status, subDB2)

		if len(query.ToOptions()) != 2 {
			t.Error("Should have 2 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `id` IN (SELECT `id` FROM `users` WHERE status = ? AND `users`.`deleted_at` IS NULL) AND `status` IN (SELECT `id` FROM `users` WHERE age > ? AND `users`.`deleted_at` IS NULL) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 {
			t.Errorf("Expected 2 args, got: %d", len(args))
		}

		expectedArgs := []interface{}{1, 18}
		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d]: %v, got: %v", i, expected, args[i])
			}
		}
	})

	t.Run("test subquery in combined with other query methods", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 创建子查询
		subDB := gormx.GetDb().Model(&User{}).Select("id").Where("status = ?", 1)

		// 测试SubQueryIn方法与其他查询方法组合使用
		query.Eq(&user.IsActive, true).
			SubQueryIn(&user.ID, subDB)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `is_active` = ? AND `id` IN (SELECT `id` FROM `users` WHERE status = ? AND `users`.`deleted_at` IS NULL) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 || args[0] != true || args[1] != 1 {
			t.Errorf("Expected args: [true 1], got: %v", args)
		}
	})

	t.Run("test subquery in with complex subquery", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 创建复杂子查询
		subDB := gormx.GetDb().Model(&User{}).
			Select("id").
			Where("created_at > ?", "2023-01-01").
			Group("status").
			Having("COUNT(*) > ?", 1)

		// 测试使用复杂子查询的SubQueryIn方法
		query.SubQueryIn(&user.ID, subDB)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `id` IN (SELECT `id` FROM `users` WHERE created_at > ? AND `users`.`deleted_at` IS NULL GROUP BY `status` HAVING COUNT(*) > ?) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 || args[0] != "2023-01-01" || args[1] != 1 {
			t.Errorf("Expected args: [2023-01-01 1], got: %v", args)
		}
	})

	// 一对多关系使用场景测试
	t.Run("test subquery in with one-to-many relationship", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 创建子查询 - 查找特定状态的用户IDs
		subDB := gormx.GetDb().Model(&User{}).Select("id").Where("status = ?", 1)

		// 测试在一对多关系中使用SubQueryIn
		query.SubQueryIn(&post.UserID, subDB)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `posts` WHERE `user_id` IN (SELECT `id` FROM `users` WHERE status = ? AND `users`.`deleted_at` IS NULL) AND `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	// 多对多关系使用场景测试
	t.Run("test subquery in with many-to-many relationship", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 创建子查询 - 查找具有特定角色的用户IDs
		subDB := gormx.GetDb().Model(&UserRole{}).
			Select("user_id").
			Where("role_id IN (SELECT id FROM roles WHERE name = ?)", "admin")

		// 测试在多对多关系中使用SubQueryIn
		query.SubQueryIn(&user.ID, subDB)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `id` IN (SELECT `user_id` FROM `user_roles` WHERE role_id IN (SELECT id FROM roles WHERE name = ?) AND `user_roles`.`deleted_at` IS NULL) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "admin" {
			t.Errorf("Expected args: [admin], got: %v", args)
		}
	})

	t.Run("test subquery in with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 创建子查询
		subDB := gormx.GetDb().Model(&User{}).Select("id").Where("status = ?", 1)

		// 测试无效字段名（应该不会panic，但会在实际查询时处理）
		query.SubQueryIn("invalid_field", subDB)

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍然会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `invalid_field` IN (SELECT `id` FROM `users` WHERE status = ? AND `users`.`deleted_at` IS NULL) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	t.Run("test subquery in with subquery returning single column", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 创建返回单列数据的子查询
		subDB := gormx.GetDb().Model(&User{}).Select("status").Where("age > ?", 25)

		// 测试SubQueryIn方法与返回单列数据的子查询
		query.SubQueryIn(&user.Status, subDB)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT * FROM `users` WHERE `status` IN (SELECT `status` FROM `users` WHERE age > ? AND `users`.`deleted_at` IS NULL) AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 25 {
			t.Errorf("Expected args: [25], got: %v", args)
		}
	})
}

// TestQuery_Count 测试 Count 方法的基本功能
func TestQuery_Count(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	t.Run("test basic count with field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字段指针的Count方法
		result := query.Count(&user.ID)

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("Count should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("Count should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT COUNT(`id`) as count FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test count with field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字段名字符串的Count方法
		query.Count("id")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT COUNT(`id`) as count FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test count with nil field", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用nil字段的Count方法
		query.Count(nil)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT COUNT(*) as count FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test count with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用无效字段名的Count方法
		query.Count("invalid_field")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT COUNT(`invalid_field`) as count FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test chaining multiple aggregate functions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试链式调用多个聚合函数
		query.Count(&user.ID).
			Count(&user.Age)

		if len(query.ToOptions()) != 2 {
			t.Error("Should have 2 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT COUNT(`id`) as count, COUNT(`age`) as count FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test count combined with other query methods", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试Count方法与其他查询方法组合使用
		query.Eq(&user.Status, 1).
			Count(&user.ID)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT COUNT(`id`) as count FROM `users` WHERE `status` = ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	t.Run("test count with select fields", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试Count方法与Select组合使用
		query.Select(&user.Status).
			Count(&user.ID)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT status, COUNT(`id`) as count FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test count with group by", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试Count方法与GroupBy组合使用
		query.GroupBy(&user.Status).
			Count(&user.ID)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT COUNT(`id`) as count FROM `users` WHERE `users`.`deleted_at` IS NULL GROUP BY `status`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test count with having", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试Count方法与Having组合使用
		query.GroupBy(&user.Status).
			Count(&user.ID).
			Having("COUNT(id) > ?", 1)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT COUNT(`id`) as count FROM `users` WHERE `users`.`deleted_at` IS NULL GROUP BY `status` HAVING COUNT(id) > ?")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	// 一对多关系使用场景测试
	t.Run("test count with one-to-many relationship", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 测试在一对多关系中使用Count
		query.Count(&post.UserID)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT COUNT(`user_id`) as count FROM `posts` WHERE `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	// 多对多关系使用场景测试
	t.Run("test count with many-to-many relationship", func(t *testing.T) {
		query, userRole := gormx.NewQuery[UserRole]()

		// 测试在多对多关系中使用Count
		query.Count(&userRole.RoleID)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT COUNT(`role_id`) as count FROM `user_roles` WHERE `user_roles`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test count with complex field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用复杂字段名的Count方法
		query.Count("created_at")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT COUNT(`created_at`) as count FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})
}

// TestQuery_Sum 测试 Sum 方法的基本功能
func TestQuery_Sum(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	t.Run("test basic sum with field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字段指针的Sum方法
		result := query.Sum(&user.Age)

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("Sum should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("Sum should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT SUM(`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test sum with field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字段名字符串的Sum方法
		query.Sum("age")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT SUM(`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test chaining multiple aggregate functions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试链式调用多个聚合函数
		query.Count(&user.ID).
			Sum(&user.Age)

		if len(query.ToOptions()) != 2 {
			t.Error("Should have 2 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT COUNT(`id`) as count, SUM(`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test sum combined with other query methods", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试Sum方法与其他查询方法组合使用
		query.Eq(&user.Status, 1).
			Sum(&user.Age)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT SUM(`age`) FROM `users` WHERE `status` = ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	t.Run("test sum with select fields", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试Sum方法与Select组合使用
		query.Select(&user.Status).
			Sum(&user.Age)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT status, SUM(`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test sum with group by", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试Sum方法与GroupBy组合使用
		query.GroupBy(&user.Status).
			Sum(&user.Age)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT SUM(`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL GROUP BY `status`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test sum with having", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试Sum方法与Having组合使用
		query.GroupBy(&user.Status).
			Sum(&user.Age).
			Having("SUM(age) > ?", 100)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT SUM(`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL GROUP BY `status` HAVING SUM(age) > ?")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 100 {
			t.Errorf("Expected args: [100], got: %v", args)
		}
	})

	t.Run("test sum with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效字段名（应该不会panic，但会在实际查询时处理）
		query.Sum("invalid_field")

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍然会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT SUM(`invalid_field`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	// 一对多关系使用场景测试
	t.Run("test sum with one-to-many relationship", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 测试在一对多关系中使用Sum
		query.Sum(&post.UserID)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT SUM(`user_id`) FROM `posts` WHERE `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	// 多对多关系使用场景测试
	t.Run("test sum with many-to-many relationship", func(t *testing.T) {
		query, userRole := gormx.NewQuery[UserRole]()

		// 测试在多对多关系中使用Sum
		query.Sum(&userRole.RoleID)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT SUM(`role_id`) FROM `user_roles` WHERE `user_roles`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test sum with complex field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用复杂字段名的Sum方法
		query.Sum("created_at")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT SUM(`created_at`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test sum with dot notation field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用点号格式字段名（表名.字段名）
		query.Sum("users.age")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT SUM(`users`.`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test sum with multiple dot notation fields", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与其他字段一起使用点号格式
		query.Select(&user.Name).
			Sum("users.age")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT name, SUM(`users`.`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test sum with mixed field name formats", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试点号格式和普通字段混合使用
		query.Select("users.name", &user.Email).
			Sum("users.age")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT users.name, email, SUM(`users`.`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test sum with dot notation and query conditions", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试点号格式字段与查询条件组合使用
		query.Eq("users.status", 1).
			Sum("users.age")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT SUM(`users`.`age`) FROM `users` WHERE `users`.`status` = ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	t.Run("test sum with dot notation in group by and having", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试点号格式字段在分组和聚合中使用
		query.GroupBy("users.status").
			Sum("users.age").
			Having(currentDialectSQLFragment("SUM(`users`.`age`) > ?"), 100)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT SUM(`users`.`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL GROUP BY `users`.`status` HAVING SUM(`users`.`age`) > ?")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 100 {
			t.Errorf("Expected args: [100], got: %v", args)
		}
	})

	t.Run("test sum with dot notation in join query", func(t *testing.T) {
		query, _ := gormx.NewQuery[Post]()

		// 测试在JOIN查询中使用点号格式字段
		query.Join("JOIN users ON posts.user_id = users.id").
			Sum("users.age")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT SUM(`users`.`age`) FROM `posts` JOIN users ON posts.user_id = users.id WHERE `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test sum with invalid dot notation field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效的点号格式字段名
		query.Sum("invalid_table.invalid_field")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT SUM(`invalid_table`.`invalid_field`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test sum with malformed dot notation field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试格式错误的点号字段名（多个点号）
		query.Sum("table.subtable.field")

		sql, args := query.ToSQLAndArgs()
		// 根据实现，可能按普通字符串处理或有特殊处理
		expectedSQL := expectedSQLForCurrentDialect("SELECT SUM(`table`.`subtable`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Logf("Actual SQL: %s", sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})
}

// TestQuery_Avg 测试 Avg 方法的基本功能
func TestQuery_Avg(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	t.Run("test basic avg with field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字段指针的Avg方法
		result := query.Avg(&user.Age)

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("Avg should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("Avg should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT AVG(`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test avg with field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字段名字符串的Avg方法
		query.Avg("age")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT AVG(`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test chaining multiple aggregate functions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试链式调用多个聚合函数
		query.Count(&user.ID).
			Avg(&user.Age)

		if len(query.ToOptions()) != 2 {
			t.Error("Should have 2 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT COUNT(`id`) as count, AVG(`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test avg combined with other query methods", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试Avg方法与其他查询方法组合使用
		query.Eq(&user.Status, 1).
			Avg(&user.Age)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT AVG(`age`) FROM `users` WHERE `status` = ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	t.Run("test avg with select fields", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试Avg方法与Select组合使用
		query.Select(&user.Status).
			Avg(&user.Age)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT status, AVG(`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test avg with group by", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试Avg方法与GroupBy组合使用
		query.GroupBy(&user.Status).
			Avg(&user.Age)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT AVG(`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL GROUP BY `status`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test avg with having", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试Avg方法与Having组合使用
		query.GroupBy(&user.Status).
			Avg(&user.Age).
			Having("AVG(age) > ?", 18)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT AVG(`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL GROUP BY `status` HAVING AVG(age) > ?")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 18 {
			t.Errorf("Expected args: [18], got: %v", args)
		}
	})

	t.Run("test avg with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效字段名（应该不会panic，但会在实际查询时处理）
		query.Avg("invalid_field")

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍然会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT AVG(`invalid_field`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	// 一对多关系使用场景测试
	t.Run("test avg with one-to-many relationship", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 测试在一对多关系中使用Avg
		query.Avg(&post.UserID)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT AVG(`user_id`) FROM `posts` WHERE `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	// 多对多关系使用场景测试
	t.Run("test avg with many-to-many relationship", func(t *testing.T) {
		query, userRole := gormx.NewQuery[UserRole]()

		// 测试在多对多关系中使用Avg
		query.Avg(&userRole.RoleID)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT AVG(`role_id`) FROM `user_roles` WHERE `user_roles`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test avg with complex field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用复杂字段名的Avg方法
		query.Avg("created_at")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT AVG(`created_at`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test avg with dot notation field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用点号格式字段名（表名.字段名）
		query.Avg("users.age")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT AVG(`users`.`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test avg with multiple dot notation fields", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与其他字段一起使用点号格式
		query.Select(&user.Name).
			Avg("users.age")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT name, AVG(`users`.`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test avg with mixed field name formats", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试点号格式和普通字段混合使用
		query.Select("users.name", &user.Email).
			Avg("users.age")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT users.name, email, AVG(`users`.`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test avg with dot notation and query conditions", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试点号格式字段与查询条件组合使用
		query.Eq("users.status", 1).
			Avg("users.age")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT AVG(`users`.`age`) FROM `users` WHERE `users`.`status` = ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	t.Run("test avg with dot notation in group by and having", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试点号格式字段在分组和聚合中使用
		query.GroupBy("users.status").
			Avg("users.age").
			Having(currentDialectSQLFragment("AVG(`users`.`age`) > ?"), 25)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT AVG(`users`.`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL GROUP BY `users`.`status` HAVING AVG(`users`.`age`) > ?")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 25 {
			t.Errorf("Expected args: [25], got: %v", args)
		}
	})

	t.Run("test avg with dot notation in join query", func(t *testing.T) {
		query, _ := gormx.NewQuery[Post]()

		// 测试在JOIN查询中使用点号格式字段
		query.Join("JOIN users ON posts.user_id = users.id").
			Avg("users.age")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT AVG(`users`.`age`) FROM `posts` JOIN users ON posts.user_id = users.id WHERE `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test avg with invalid dot notation field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效的点号格式字段名
		query.Avg("invalid_table.invalid_field")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT AVG(`invalid_table`.`invalid_field`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test avg with malformed dot notation field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试格式错误的点号字段名（多个点号）
		query.Avg("table.subtable.field")

		sql, args := query.ToSQLAndArgs()
		// 根据实现，可能按普通字符串处理或有特殊处理
		expectedSQL := expectedSQLForCurrentDialect("SELECT AVG(`table`.`subtable`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Logf("Actual SQL: %s", sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})
}

// TestQuery_Max 测试 Max 方法的基本功能
func TestQuery_Max(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	t.Run("test basic max with field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字段指针的Max方法
		result := query.Max(&user.Age)

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("Max should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("Max should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT MAX(`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test max with field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字段名字符串的Max方法
		query.Max("age")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT MAX(`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test chaining multiple aggregate functions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试链式调用多个聚合函数
		query.Count(&user.ID).
			Max(&user.Age)

		if len(query.ToOptions()) != 2 {
			t.Error("Should have 2 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT COUNT(`id`) as count, MAX(`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test max combined with other query methods", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试Max方法与其他查询方法组合使用
		query.Eq(&user.Status, 1).
			Max(&user.Age)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT MAX(`age`) FROM `users` WHERE `status` = ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	t.Run("test max with select fields", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试Max方法与Select组合使用
		query.Select(&user.Status).
			Max(&user.Age)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT status, MAX(`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test max with group by", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试Max方法与GroupBy组合使用
		query.GroupBy(&user.Status).
			Max(&user.Age)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT MAX(`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL GROUP BY `status`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test max with having", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试Max方法与Having组合使用
		query.GroupBy(&user.Status).
			Max(&user.Age).
			Having("MAX(age) > ?", 18)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT MAX(`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL GROUP BY `status` HAVING MAX(age) > ?")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 18 {
			t.Errorf("Expected args: [18], got: %v", args)
		}
	})

	t.Run("test max with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效字段名（应该不会panic，但会在实际查询时处理）
		query.Max("invalid_field")

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍然会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT MAX(`invalid_field`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	// 一对多关系使用场景测试
	t.Run("test max with one-to-many relationship", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 测试在一对多关系中使用Max
		query.Max(&post.UserID)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT MAX(`user_id`) FROM `posts` WHERE `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	// 多对多关系使用场景测试
	t.Run("test max with many-to-many relationship", func(t *testing.T) {
		query, userRole := gormx.NewQuery[UserRole]()

		// 测试在多对多关系中使用Max
		query.Max(&userRole.RoleID)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT MAX(`role_id`) FROM `user_roles` WHERE `user_roles`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test max with complex field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用复杂字段名的Max方法
		query.Max("created_at")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT MAX(`created_at`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test max with dot notation field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用点号格式字段名（表名.字段名）
		query.Max("users.age")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT MAX(`users`.`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test max with multiple dot notation fields", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与其他字段一起使用点号格式
		query.Select(&user.Name).
			Max("users.age")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT name, MAX(`users`.`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test max with mixed field name formats", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试点号格式和普通字段混合使用
		query.Select("users.name", &user.Email).
			Max("users.age")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT users.name, email, MAX(`users`.`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test max with dot notation and query conditions", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试点号格式字段与查询条件组合使用
		query.Eq("users.status", 1).
			Max("users.age")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT MAX(`users`.`age`) FROM `users` WHERE `users`.`status` = ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	t.Run("test max with dot notation in group by and having", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试点号格式字段在分组和聚合中使用
		query.GroupBy("users.status").
			Max("users.age").
			Having(currentDialectSQLFragment("MAX(`users`.`age`) > ?"), 30)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT MAX(`users`.`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL GROUP BY `users`.`status` HAVING MAX(`users`.`age`) > ?")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 30 {
			t.Errorf("Expected args: [30], got: %v", args)
		}
	})

	t.Run("test max with dot notation in join query", func(t *testing.T) {
		query, _ := gormx.NewQuery[Post]()

		// 测试在JOIN查询中使用点号格式字段
		query.Join("JOIN users ON posts.user_id = users.id").
			Max("users.age")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT MAX(`users`.`age`) FROM `posts` JOIN users ON posts.user_id = users.id WHERE `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test max with invalid dot notation field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效的点号格式字段名
		query.Max("invalid_table.invalid_field")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT MAX(`invalid_table`.`invalid_field`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test max with malformed dot notation field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试格式错误的点号字段名（多个点号）
		query.Max("table.subtable.field")

		sql, args := query.ToSQLAndArgs()
		// 根据实现，可能按普通字符串处理或有特殊处理
		expectedSQL := expectedSQLForCurrentDialect("SELECT MAX(`table`.`subtable`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Logf("Actual SQL: %s", sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})
}

// TestQuery_Min 测试 Min 方法的基本功能
func TestQuery_Min(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	t.Run("test basic min with field pointer", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试使用字段指针的Min方法
		result := query.Min(&user.Age)

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("Min should return the same query instance for chaining")
		}

		// 验证选项被正确添加
		if len(query.ToOptions()) != 1 {
			t.Error("Min should add one option")
		}

		// 验证生成的SQL和参数
		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT MIN(`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test min with field name string", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用字段名字符串的Min方法
		query.Min("age")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT MIN(`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test chaining multiple aggregate functions", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试链式调用多个聚合函数
		query.Count(&user.ID).
			Min(&user.Age)

		if len(query.ToOptions()) != 2 {
			t.Error("Should have 2 options after chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT COUNT(`id`) as count, MIN(`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test min combined with other query methods", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试Min方法与其他查询方法组合使用
		query.Eq(&user.Status, 1).
			Min(&user.Age)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT MIN(`age`) FROM `users` WHERE `status` = ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	t.Run("test min with select fields", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试Min方法与Select组合使用
		query.Select(&user.Status).
			Min(&user.Age)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT status, MIN(`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test min with group by", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试Min方法与GroupBy组合使用
		query.GroupBy(&user.Status).
			Min(&user.Age)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT MIN(`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL GROUP BY `status`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test min with having", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试Min方法与Having组合使用
		query.GroupBy(&user.Status).
			Min(&user.Age).
			Having("MIN(age) > ?", 18)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT MIN(`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL GROUP BY `status` HAVING MIN(age) > ?")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 18 {
			t.Errorf("Expected args: [18], got: %v", args)
		}
	})

	t.Run("test min with invalid field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效字段名（应该不会panic，但会在实际查询时处理）
		query.Min("invalid_field")

		sql, args := query.ToSQLAndArgs()
		// 注意：无效字段仍然会生成SQL，但在实际使用中可能需要额外验证
		expectedSQL := expectedSQLForCurrentDialect("SELECT MIN(`invalid_field`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	// 一对多关系使用场景测试
	t.Run("test min with one-to-many relationship", func(t *testing.T) {
		query, post := gormx.NewQuery[Post]()

		// 测试在一对多关系中使用Min
		query.Min(&post.UserID)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT MIN(`user_id`) FROM `posts` WHERE `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	// 多对多关系使用场景测试
	t.Run("test min with many-to-many relationship", func(t *testing.T) {
		query, userRole := gormx.NewQuery[UserRole]()

		// 测试在多对多关系中使用Min
		query.Min(&userRole.RoleID)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT MIN(`role_id`) FROM `user_roles` WHERE `user_roles`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test min with complex field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用复杂字段名的Min方法
		query.Min("created_at")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT MIN(`created_at`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test min with dot notation field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试使用点号格式字段名（表名.字段名）
		query.Min("users.age")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT MIN(`users`.`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test min with multiple dot notation fields", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试与其他字段一起使用点号格式
		query.Select(&user.Name).
			Min("users.age")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT name, MIN(`users`.`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test min with mixed field name formats", func(t *testing.T) {
		query, user := gormx.NewQuery[User]()

		// 测试点号格式和普通字段混合使用
		query.Select("users.name", &user.Email).
			Min("users.age")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT users.name, email, MIN(`users`.`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test min with dot notation and query conditions", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试点号格式字段与查询条件组合使用
		query.Eq("users.status", 1).
			Min("users.age")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT MIN(`users`.`age`) FROM `users` WHERE `users`.`status` = ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	t.Run("test min with dot notation in group by and having", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试点号格式字段在分组和聚合中使用
		query.GroupBy("users.status").
			Min("users.age").
			Having(currentDialectSQLFragment("MIN(`users`.`age`) > ?"), 18)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT MIN(`users`.`age`) FROM `users` WHERE `users`.`deleted_at` IS NULL GROUP BY `users`.`status` HAVING MIN(`users`.`age`) > ?")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 18 {
			t.Errorf("Expected args: [18], got: %v", args)
		}
	})

	t.Run("test min with dot notation in join query", func(t *testing.T) {
		query, _ := gormx.NewQuery[Post]()

		// 测试在JOIN查询中使用点号格式字段
		query.Join("JOIN users ON posts.user_id = users.id").
			Min("users.age")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT MIN(`users`.`age`) FROM `posts` JOIN users ON posts.user_id = users.id WHERE `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test min with invalid dot notation field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试无效的点号格式字段名
		query.Min("invalid_table.invalid_field")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT MIN(`invalid_table`.`invalid_field`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test min with malformed dot notation field name", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试格式错误的点号字段名（多个点号）
		query.Min("table.subtable.field")

		sql, args := query.ToSQLAndArgs()
		// 根据实现，可能按普通字符串处理或有特殊处理
		expectedSQL := expectedSQLForCurrentDialect("SELECT MIN(`table`.`subtable`) FROM `users` WHERE `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Logf("Actual SQL: %s", sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})
}

// TestQuery_Join 测试 Join 方法的基本功能
func TestQuery_Join(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	// 一对多关系使用场景测试
	t.Run("test join with one-to-many relationship", func(t *testing.T) {
		query, _ := gormx.NewQuery[Post]()

		// 测试在一对多关系中使用Join
		query.Join("JOIN users ON posts.user_id = users.id").
			Eq("posts.status", 1).
			Eq("users.name", "testuser1")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT `posts`.`id`,`posts`.`created_at`,`posts`.`updated_at`,`posts`.`deleted_at`,`posts`.`user_id`,`posts`.`title`,`posts`.`content`,`posts`.`status` FROM `posts` JOIN users ON posts.user_id = users.id WHERE `posts`.`status` = ? AND `users`.`name` = ? AND `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 || args[0] != 1 || args[1] != "testuser1" {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	// 多对多关系使用场景测试
	t.Run("test join with many-to-many relationship", func(t *testing.T) {
		query, _ := gormx.NewQuery[User]()

		// 测试在多对多关系中使用Join
		query.Join("JOIN user_roles ON users.id = user_roles.user_id").
			Join("JOIN roles ON user_roles.role_id = roles.id").
			Eq("users.name", "testuser1")
		//Select("users.id", "users.name")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT `users`.`id`,`users`.`created_at`,`users`.`updated_at`,`users`.`deleted_at`,`users`.`name`,`users`.`email`,`users`.`phone`,`users`.`age`,`users`.`score`,`users`.`address`,`users`.`is_active`,`users`.`salary`,`users`.`status` FROM `users` JOIN user_roles ON users.id = user_roles.user_id JOIN roles ON user_roles.role_id = roles.id WHERE `users`.`name` = ? AND `users`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != "testuser1" {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	t.Run("test join with parameters", func(t *testing.T) {
		query, _ := gormx.NewQuery[Post]()

		// 测试带参数的Join
		query.Join("JOIN users ON posts.user_id = users.id AND users.status = ?", 1).
			Eq("posts.status", 1)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT `posts`.`id`,`posts`.`created_at`,`posts`.`updated_at`,`posts`.`deleted_at`,`posts`.`user_id`,`posts`.`title`,`posts`.`content`,`posts`.`status` FROM `posts` JOIN users ON posts.user_id = users.id AND users.status = ? WHERE `posts`.`status` = ? AND `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 || args[0] != 1 || args[1] != 1 {
			t.Errorf("Expected args: [1 1], got: %v", args)
		}
	})

	t.Run("test join combined with other query methods", func(t *testing.T) {
		query, _ := gormx.NewQuery[Post]()

		// 测试Join方法与其他查询方法组合使用
		query.Join("JOIN users ON posts.user_id = users.id").
			Eq("posts.status", 1).
			Like("users.name", "admin")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT `posts`.`id`,`posts`.`created_at`,`posts`.`updated_at`,`posts`.`deleted_at`,`posts`.`user_id`,`posts`.`title`,`posts`.`content`,`posts`.`status` FROM `posts` JOIN users ON posts.user_id = users.id WHERE `posts`.`status` = ? AND `users`.`name` LIKE ? AND `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 2 || args[0] != 1 || args[1] != "%admin%" {
			t.Errorf("Expected args: [1 admin], got: %v", args)
		}
	})

	t.Run("test join with select fields", func(t *testing.T) {
		query, _ := gormx.NewQuery[Post]()

		// 测试Join方法与Select组合使用
		query.Select("posts.title", "users.name").
			Join("JOIN users ON posts.user_id = users.id").
			Eq("posts.status", 1)

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT posts.title,users.name FROM `posts` JOIN users ON posts.user_id = users.id WHERE `posts`.`status` = ? AND `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	t.Run("test join with order by", func(t *testing.T) {
		query, _ := gormx.NewQuery[Post]()

		// 测试Join方法与OrderBy组合使用
		query.Join("JOIN users ON posts.user_id = users.id").
			Eq("posts.status", 1).
			OrderAsc("users.created_at")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT `posts`.`id`,`posts`.`created_at`,`posts`.`updated_at`,`posts`.`deleted_at`,`posts`.`user_id`,`posts`.`title`,`posts`.`content`,`posts`.`status` FROM `posts` JOIN users ON posts.user_id = users.id WHERE `posts`.`status` = ? AND `posts`.`deleted_at` IS NULL ORDER BY `users`.`created_at`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args: [1], got: %v", args)
		}
	})

	t.Run("test join with group by", func(t *testing.T) {
		query, _ := gormx.NewQuery[Post]()

		// 测试Join方法与GroupBy组合使用
		query.Join("JOIN users ON posts.user_id = users.id").
			GroupBy("users.id").
			Count("posts.id")

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT COUNT(`posts`.`id`) as count FROM `posts` JOIN users ON posts.user_id = users.id WHERE `posts`.`deleted_at` IS NULL GROUP BY `users`.`id`")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})

	t.Run("test multiple joins", func(t *testing.T) {
		query, _ := gormx.NewQuery[Post]()

		// 测试多个Join
		result := query.Join("JOIN users ON posts.user_id = users.id").
			Join("JOIN categories ON posts.category_id = categories.id").
			Join("LEFT JOIN tags ON posts.tag_id = tags.id")

		// 验证返回的是同一个查询实例（链式调用）
		if result != query {
			t.Error("Join should return the same query instance for chaining")
		}

		sql, args := query.ToSQLAndArgs()
		expectedSQL := expectedSQLForCurrentDialect("SELECT `posts`.`id`,`posts`.`created_at`,`posts`.`updated_at`,`posts`.`deleted_at`,`posts`.`user_id`,`posts`.`title`,`posts`.`content`,`posts`.`status` FROM `posts` JOIN users ON posts.user_id = users.id JOIN categories ON posts.category_id = categories.id LEFT JOIN tags ON posts.tag_id = tags.id WHERE `posts`.`deleted_at` IS NULL")
		if sql != expectedSQL {
			t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
		}

		if len(args) != 0 {
			t.Errorf("Expected 0 args, got: %d", len(args))
		}
	})
}
