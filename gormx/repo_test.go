package gormx_test

import (
	"errors"
	"github.com/shijl0925/go-toolkits/gormx"
	"gorm.io/gorm"
	"strings"
	"testing"
)

// TestBaseRepo_SelectOneByOpts 测试 SelectOneByOpts 方法
func TestBaseRepo_SelectOneByOpts(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	repo := &gormx.BaseRepo[User]{}

	t.Run("test select one record by options", func(t *testing.T) {
		// 测试正常查询单条记录
		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Name, "testuser1")

		result, err := repo.SelectOneByOpts(query.ToOptions()...)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if result.Name != "testuser1" {
			t.Errorf("Expected name: testuser1, got: %s", result.Name)
		}

		if result.ID == 0 {
			t.Error("Expected ID to be set")
		}
	})

	t.Run("test select one record by options with no matches", func(t *testing.T) {
		// 测试无匹配记录的情况
		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Name, "nonexistent_user")

		_, err := repo.SelectOneByOpts(query.ToOptions()...)

		if err == nil {
			t.Error("Expected error for no matching records")
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Errorf("Expected gorm.ErrRecordNotFound, got: %v", err)
		}
	})

	t.Run("test select one record by options with multiple matches", func(t *testing.T) {
		// 测试匹配多条记录的情况
		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Status, 1) // 假设有多条status为1的记录

		_, err := repo.SelectOneByOpts(query.ToOptions()...)

		if err == nil {
			t.Error("Expected error for multiple matching records")
		}

		if !strings.Contains(err.Error(), "expected one result") {
			t.Errorf("Expected error message about expecting one result, got: %v", err)
		}
	})

	t.Run("test select one record by options with complex conditions", func(t *testing.T) {
		// 测试复杂查询条件
		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Status, 1).
			Gt(&user.Age, 30).
			Like(&user.Name, "test")

		result, err := repo.SelectOneByOpts(query.ToOptions()...)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if result.Status != 1 {
			t.Errorf("Expected status: 1, got: %d", result.Status)
		}

		if result.Age <= 30 {
			t.Errorf("Expected age > 20, got: %d", result.Age)
		}

		if !strings.Contains(result.Name, "test") {
			t.Errorf("Expected name to contain 'test', got: %s", result.Name)
		}
	})

	t.Run("test select one record by options with ordering", func(t *testing.T) {
		// 测试带排序的查询
		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Status, 1).
			Gt(&user.Age, 30).
			OrderAsc(&user.CreatedAt) // 按创建时间升序，确保获取第一条

		result, err := repo.SelectOneByOpts(query.ToOptions()...)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if result.Status != 1 {
			t.Errorf("Expected status: 1, got: %d", result.Status)
		}
	})

	t.Run("test select one record by options with limit and offset", func(t *testing.T) {
		// 测试带limit和offset的查询（虽然这些在方法内部会被覆盖）
		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Status, 1).
			Limit(10).
			Offset(2)

		result, err := repo.SelectOneByOpts(query.ToOptions()...)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if result.Status != 1 {
			t.Errorf("Expected status: 1, got: %d", result.Status)
		}
	})

	t.Run("test select one record by options with joins", func(t *testing.T) {
		// 测试带JOIN的查询
		query, _ := gormx.NewQuery[User]()
		query.Eq("users.status", 1).Eq("users.email", "test1@example.com").
			Join("LEFT JOIN posts ON users.id = posts.user_id").
			Select("users.id", "users.name", "users.status")

		sql, args := query.ToSQLAndArgs()
		t.Logf("SQL: %s, args: %v", sql, args)

		result, err := repo.SelectOneByOpts(query.ToOptions()...)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if result.Status != 1 {
			t.Errorf("Expected status: 1, got: %d", result.Status)
		}
	})

	t.Run("test select one record by options with group by", func(t *testing.T) {
		// 测试带GROUP BY的查询（这类查询可能不适合此方法，但需要测试其行为）
		query, user := gormx.NewQuery[User]()
		query.IsNotNull(&user.Name).
			GroupBy(&user.Status).
			Select(&user.Status)

		_, err := repo.SelectOneByOpts(query.ToOptions()...)

		// 这类查询可能返回多条记录或者错误，取决于数据库实现
		// 我们主要确保不会panic
		t.Logf("Group by query result: %v", err)
	})

	t.Run("test select one record by options with in condition", func(t *testing.T) {
		// 测试IN条件查询
		query, user := gormx.NewQuery[User]()
		query.In(&user.Status, []int{0, 1}).Eq(&user.Name, "testuser1")

		result, err := repo.SelectOneByOpts(query.ToOptions()...)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		found := false
		for _, status := range []int{0, 1} {
			if result.Status == status {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Expected Status to be one of [0, 1], got: %d", result.Status)
		}
	})

	t.Run("test select one record by options with not in condition", func(t *testing.T) {
		// 测试NOT IN条件查询
		query, user := gormx.NewQuery[User]()
		query.NotIn(&user.Status, []int{1}).Eq(&user.Name, "otheruser") // 排除不存在的Status

		result, err := repo.SelectOneByOpts(query.ToOptions()...)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if result.Status < 0 {
			t.Errorf("Expected valid Status, got: %d", result.Status)
		}
	})

	t.Run("test select one record by options with between condition", func(t *testing.T) {
		// 测试BETWEEN条件查询
		query, user := gormx.NewQuery[User]()
		query.Between(&user.Age, 31, 40)

		result, err := repo.SelectOneByOpts(query.ToOptions()...)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if result.Age < 31 || result.Age > 40 {
			t.Errorf("Expected age between 31 and 40, got: %d", result.Age)
		}
	})

	t.Run("test select one record by options with null condition", func(t *testing.T) {
		// 测试NULL条件查询
		query, user := gormx.NewQuery[User]()
		query.IsNotNull(&user.Email).Eq(&user.Name, "testuser1")

		result, err := repo.SelectOneByOpts(query.ToOptions()...)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if result.Email == "" {
			t.Error("Expected non-empty email")
		}
	})

	t.Run("test select one record by options with regexp condition", func(t *testing.T) {
		// 测试正则表达式条件查询
		query, user := gormx.NewQuery[User]()
		query.Regexp(&user.Name, "^other.*")

		result, err := repo.SelectOneByOpts(query.ToOptions()...)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if !strings.HasPrefix(result.Name, "other") {
			t.Errorf("Expected name to start with 'test', got: %s", result.Name)
		}
	})

	t.Run("test select one record by options with subquery", func(t *testing.T) {
		// 测试子查询条件
		query, user := gormx.NewQuery[User]()

		// 创建子查询
		subDB := gormx.GetDb().Model(&User{}).Select("MAX(id)").Where("status = ?", 1)
		query.SubQueryEq(&user.ID, subDB)

		result, err := repo.SelectOneByOpts(query.ToOptions()...)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if result.Status != 1 {
			t.Errorf("Expected status: 1, got: %d", result.Status)
		}
	})

	t.Run("test select one record by options with or condition", func(t *testing.T) {
		// 测试OR条件查询
		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Status, 0).
			Or("age > ? AND age <= ?", 25, 28)

		result, err := repo.SelectOneByOpts(query.ToOptions()...)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果满足任一条件
		if result.Status != 0 || (result.Age <= 25 || result.Age > 28) {
			t.Errorf("Expected status=0 OR (age > 25 AND age<=28), got status=%d, age=%d", result.Status, result.Age)
		}
	})

	t.Run("test select one record by options with not condition", func(t *testing.T) {
		// 测试NOT条件查询
		query, _ := gormx.NewQuery[User]()
		query.Not("status = ?", 1)

		result, err := repo.SelectOneByOpts(query.ToOptions()...)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if result.Status == 1 {
			t.Errorf("Expected status != 0, got: %d", result.Status)
		}
	})
}

// TestBaseRepo_SelectListByOpts 测试 SelectListByOpts 方法
func TestBaseRepo_SelectListByOpts(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	// 创建 UserRepo 实例
	userRepo := &gormx.BaseRepo[User]{}

	t.Run("test select all records", func(t *testing.T) {
		// 测试查询所有记录
		users, err := userRepo.SelectListByOpts()

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(users) == 0 {
			t.Error("Expected to find users, but got empty result")
		}
	})

	t.Run("test select with condition", func(t *testing.T) {
		// 测试带条件查询
		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Status, 1)

		users, err := userRepo.SelectListByOpts(query.ToOptions()...)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果
		if len(users) > 0 {
			for _, user := range users {
				if user.Status != 1 {
					t.Errorf("Expected status to be 1, got: %d", user.Status)
				}
			}
		}
	})

	t.Run("test select with multiple conditions", func(t *testing.T) {
		// 测试多条件查询
		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Status, 1).Gt(&user.Age, 18)

		users, err := userRepo.SelectListByOpts(
			query.ToOptions()...,
		)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果
		if len(users) > 0 {
			for _, user := range users {
				if user.Status != 1 || user.Age <= 18 {
					t.Errorf("Result doesn't match conditions: status=%d, age=%d", user.Status, user.Age)
				}
			}
		}
	})

	t.Run("test select with order", func(t *testing.T) {
		// 测试排序查询
		users, err := userRepo.SelectListByOpts(
			func(db *gorm.DB) *gorm.DB {
				return db.Order("created_at DESC")
			},
		)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果是否按创建时间降序排列
		if len(users) > 1 {
			for i := 1; i < len(users); i++ {
				if users[i-1].CreatedAt.Before(users[i].CreatedAt) {
					t.Error("Users are not sorted by created_at in descending order")
				}
			}
		}
	})

	t.Run("test select with limit", func(t *testing.T) {
		// 测试限制查询结果数量
		users, err := userRepo.SelectListByOpts(
			func(db *gorm.DB) *gorm.DB {
				return db.Limit(2)
			},
		)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果数量不超过限制
		if len(users) > 2 {
			t.Errorf("Expected at most 2 results, got: %d", len(users))
		}
	})

	t.Run("test select with complex query", func(t *testing.T) {
		// 测试复杂查询：带条件、排序和限制
		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Status, 1).Limit(3).OrderAsc(&user.Age)

		users, err := userRepo.SelectListByOpts(
			query.ToOptions()...,
		)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果数量
		if len(users) > 3 {
			t.Errorf("Expected at most 3 results, got: %d", len(users))
		}

		// 验证结果是否满足条件
		for _, user := range users {
			if user.Status != 1 {
				t.Errorf("Expected status to be 1, got: %d", user.Status)
			}
		}

		// 验证结果是否按年龄升序排列
		if len(users) > 1 {
			for i := 1; i < len(users); i++ {
				if users[i-1].Age > users[i].Age {
					t.Error("Users are not sorted by age in ascending order")
				}
			}
		}
	})

	t.Run("test select with no matching records", func(t *testing.T) {
		// 测试无匹配记录的情况
		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Name, "nonexistent")

		users, err := userRepo.SelectListByOpts(
			query.ToOptions()...,
		)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 应该返回空切片而不是错误
		if users == nil {
			t.Error("Expected empty slice, got nil")
		}

		if len(users) != 0 {
			t.Errorf("Expected empty slice, got length: %d", len(users))
		}
	})

	t.Run("test select with join", func(t *testing.T) {
		// 创建 postRepo 实例
		postRepo := &gormx.BaseRepo[Post]{}
		query, _ := gormx.NewQuery[Post]()
		query.Join("JOIN users ON posts.user_id = users.id")
		query.Eq("users.name", "testuser1")

		posts, err := postRepo.SelectListByOpts(
			query.ToOptions()...,
		)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		for _, post := range posts {
			userId := post.UserID
			user, err := userRepo.SelectOneById(int(userId))
			if err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
			if user.Name != "testuser1" {
				t.Errorf("Expected user name to be 'testuser1', got: %s", user.Name)
			}
		}
	})
}
