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

// TestBaseRepo_SelectOneByMap 测试 SelectOneByMap 方法
func TestBaseRepo_SelectOneByMap(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	repo := &gormx.BaseRepo[User]{}

	t.Run("test select one record by map", func(t *testing.T) {
		// 测试正常查询单条记录
		columnMap := map[string]interface{}{
			"name": "testuser1",
		}

		result, err := repo.SelectOneByMap(columnMap)

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

	t.Run("test select one record by map with no matches", func(t *testing.T) {
		// 测试无匹配记录的情况
		columnMap := map[string]interface{}{
			"name": "nonexistent_user",
		}

		_, err := repo.SelectOneByMap(columnMap)

		if err == nil {
			t.Error("Expected error for no matching records")
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Errorf("Expected gorm.ErrRecordNotFound, got: %v", err)
		}
	})

	t.Run("test select one record by map with multiple matches", func(t *testing.T) {
		// 测试匹配多条记录的情况
		columnMap := map[string]interface{}{
			"status": 1, // 假设有多条status为1的记录
		}

		_, err := repo.SelectOneByMap(columnMap)

		if err == nil {
			t.Error("Expected error for multiple matching records")
		}

		if !strings.Contains(err.Error(), "expected one result") {
			t.Errorf("Expected error message about expecting one result, got: %v", err)
		}
	})

	t.Run("test select one record by map with complex conditions", func(t *testing.T) {
		// 测试复杂查询条件
		columnMap := map[string]interface{}{
			"status": 1,
			"age":    35,
		}

		result, err := repo.SelectOneByMap(columnMap)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if result.Status != 1 {
			t.Errorf("Expected status: 1, got: %d", result.Status)
		}

		if result.Age != 35 {
			t.Errorf("Expected age: 35, got: %d", result.Age)
		}
	})

	t.Run("test select one record by map with empty map", func(t *testing.T) {
		// 测试空map的情况
		columnMap := map[string]interface{}{}

		_, err := repo.SelectOneByMap(columnMap)

		if err == nil {
			t.Error("Expected error for empty map")
		}
	})
}

// TestBaseRepo_SelectListByMap 测试 SelectListByMap 方法
func TestBaseRepo_SelectListByMap(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	repo := &gormx.BaseRepo[User]{}

	t.Run("test select list by map with normal conditions", func(t *testing.T) {
		// 测试正常查询多条记录
		columnMap := map[string]interface{}{
			"status": 1,
		}

		results, err := repo.SelectListByMap(columnMap)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(results) == 0 {
			t.Error("Expected to find users with status 1, but got empty result")
		}

		for _, user := range results {
			if user.Status != 1 {
				t.Errorf("Expected status to be 1, got: %d", user.Status)
			}
		}
	})

	t.Run("test select list by map with no matches", func(t *testing.T) {
		// 测试无匹配记录的情况
		columnMap := map[string]interface{}{
			"name": "nonexistent_user",
		}

		results, err := repo.SelectListByMap(columnMap)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(results) != 0 {
			t.Errorf("Expected empty result, got: %d records", len(results))
		}
	})

	t.Run("test select list by map with complex conditions", func(t *testing.T) {
		// 测试复杂查询条件
		columnMap := map[string]interface{}{
			"status": 1,
			"age":    35,
		}

		results, err := repo.SelectListByMap(columnMap)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(results) == 0 {
			t.Error("Expected to find users with status 1 and age 35, but got empty result")
		}

		for _, user := range results {
			if user.Status != 1 || user.Age != 35 {
				t.Errorf("Expected status=1 and age=35, got status=%d, age=%d", user.Status, user.Age)
			}
		}
	})

	t.Run("test select list by map with empty map", func(t *testing.T) {
		// 测试空map的情况
		columnMap := map[string]interface{}{}

		results, err := repo.SelectListByMap(columnMap)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(results) == 0 {
			t.Error("Expected to find all users, but got empty result")
		}
	})
}

// TestBaseRepo_SelectPage 测试 SelectPage 方法
func TestBaseRepo_SelectPage(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	// 创建 UserRepo 实例
	userRepo := &gormx.BaseRepo[User]{}

	t.Run("test select page with normal pagination", func(t *testing.T) {
		// 测试正常分页查询
		page := 1
		pageSize := 2

		users, total, err := userRepo.SelectPage(page, pageSize)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果数量
		if len(users) > pageSize {
			t.Errorf("Expected at most %d results, got: %d", pageSize, len(users))
		}

		// 验证总数
		if total == 0 {
			t.Error("Expected total count to be greater than 0")
		}
	})

	t.Run("test select page with out of range page", func(t *testing.T) {
		// 测试超出范围的页码
		page := 100
		pageSize := 2

		users, total, err := userRepo.SelectPage(page, pageSize)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果为空
		if len(users) != 0 {
			t.Errorf("Expected empty result, got: %d records", len(users))
		}

		// 验证总数
		if total != 0 {
			t.Error("Expected total count is 0")
		}
	})

	t.Run("test select page with zero page size", func(t *testing.T) {
		// 测试页大小为0的情况
		page := 1
		pageSize := 0

		users, total, err := userRepo.SelectPage(page, pageSize)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果为空
		if len(users) != 0 {
			t.Errorf("Expected empty result, got: %d records", len(users))
		}

		// 验证总数
		if total != 0 {
			t.Error("Expected total count is 0")
		}
	})

	t.Run("test select page with conditions", func(t *testing.T) {
		// 测试带条件的分页查询
		page := 1
		pageSize := 2

		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Status, 1)

		users, total, err := userRepo.SelectPage(page, pageSize, query.ToOptions()...)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果数量
		if len(users) > pageSize {
			t.Errorf("Expected at most %d results, got: %d", pageSize, len(users))
		}

		// 验证结果是否满足条件
		for _, user := range users {
			if user.Status != 1 {
				t.Errorf("Expected status to be 1, got: %d", user.Status)
			}
		}

		// 验证总数
		if total == 0 {
			t.Error("Expected total count to be greater than 0")
		}
	})

	t.Run("test select page with order", func(t *testing.T) {
		// 测试带排序的分页查询
		page := 1
		pageSize := 2

		query, user := gormx.NewQuery[User]()
		query.OrderDesc(&user.CreatedAt)

		users, total, err := userRepo.SelectPage(page, pageSize, query.ToOptions()...)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证结果数量
		if len(users) > pageSize {
			t.Errorf("Expected at most %d results, got: %d", pageSize, len(users))
		}

		// 验证结果是否按创建时间降序排列
		if len(users) > 1 {
			for i := 1; i < len(users); i++ {
				if users[i-1].CreatedAt.Before(users[i].CreatedAt) {
					t.Error("Users are not sorted by created_at in descending order")
				}
			}
		}

		// 验证总数
		if total == 0 {
			t.Error("Expected total count to be greater than 0")
		}
	})
}

// TestBaseRepo_SelectCount 测试 SelectCount 方法
func TestBaseRepo_SelectCount(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	// 创建 UserRepo 实例
	userRepo := &gormx.BaseRepo[User]{}

	t.Run("test select count with no conditions", func(t *testing.T) {
		// 测试无条件统计总数
		count, err := userRepo.SelectCount()

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if count == 0 {
			t.Error("Expected count to be greater than 0")
		}
	})

	t.Run("test select count with conditions", func(t *testing.T) {
		// 测试带条件统计总数
		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Status, 1)

		count, err := userRepo.SelectCount(query.ToOptions()...)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if count == 0 {
			t.Error("Expected count to be greater than 0")
		}
	})

	t.Run("test select count with complex conditions", func(t *testing.T) {
		// 测试复杂条件统计总数
		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Status, 1).Gt(&user.Age, 18)

		count, err := userRepo.SelectCount(query.ToOptions()...)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if count == 0 {
			t.Error("Expected count to be greater than 0")
		}
	})

	t.Run("test select count with no matching records", func(t *testing.T) {
		// 测试无匹配记录的情况
		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Name, "nonexistent_user")

		count, err := userRepo.SelectCount(query.ToOptions()...)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if count != 0 {
			t.Errorf("Expected count to be 0, got: %d", count)
		}
	})

	t.Run("test select count with order", func(t *testing.T) {
		// 测试带排序条件统计总数（排序不影响总数）
		query, user := gormx.NewQuery[User]()
		query.OrderDesc(&user.CreatedAt)

		count, err := userRepo.SelectCount(query.ToOptions()...)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if count == 0 {
			t.Error("Expected count to be greater than 0")
		}
	})

	t.Run("test select count with limit", func(t *testing.T) {
		// 测试带limit条件统计总数（limit不影响总数）
		query, _ := gormx.NewQuery[User]()
		query.Limit(2)

		count, err := userRepo.SelectCount(query.ToOptions()...)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if count == 0 {
			t.Error("Expected count to be greater than 0")
		}
	})
}

// TestBaseRepo_Exists 测试 Exists 方法
func TestBaseRepo_Exists(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	// 创建 UserRepo 实例
	userRepo := &gormx.BaseRepo[User]{}

	t.Run("test exists with no conditions", func(t *testing.T) {
		// 测试无条件判断记录是否存在
		exists := userRepo.Exists()

		if !exists {
			t.Error("Expected records to exist")
		}
	})

	t.Run("test exists with conditions that match", func(t *testing.T) {
		// 测试带条件判断记录是否存在（条件匹配）
		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Name, "testuser1")

		exists := userRepo.Exists(query.ToOptions()...)

		if !exists {
			t.Error("Expected record to exist")
		}
	})

	t.Run("test exists with conditions that do not match", func(t *testing.T) {
		// 测试带条件判断记录是否存在（条件不匹配）
		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Name, "nonexistent_user")

		exists := userRepo.Exists(query.ToOptions()...)

		if exists {
			t.Error("Expected no record to exist")
		}
	})

	t.Run("test exists with complex conditions", func(t *testing.T) {
		// 测试复杂条件判断记录是否存在
		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Status, 1).Gt(&user.Age, 18)

		exists := userRepo.Exists(query.ToOptions()...)

		if !exists {
			t.Error("Expected record to exist")
		}
	})

	t.Run("test exists with order", func(t *testing.T) {
		// 测试带排序条件判断记录是否存在（排序不影响结果）
		query, user := gormx.NewQuery[User]()
		query.OrderDesc(&user.CreatedAt)

		exists := userRepo.Exists(query.ToOptions()...)

		if !exists {
			t.Error("Expected record to exist")
		}
	})

	t.Run("test exists with limit", func(t *testing.T) {
		// 测试带limit条件判断记录是否存在（limit不影响结果）
		query, _ := gormx.NewQuery[User]()
		query.Limit(1)

		exists := userRepo.Exists(query.ToOptions()...)

		if !exists {
			t.Error("Expected record to exist")
		}
	})
}

// TestBaseRepo_Insert 测试 Insert 方法
func TestBaseRepo_Insert(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	// 创建 UserRepo 实例
	userRepo := &gormx.BaseRepo[User]{}

	t.Run("test insert a new record", func(t *testing.T) {
		// 测试插入一条新记录
		newUser := &User{
			Name:   "newuser",
			Email:  "newuser@example.com",
			Age:    25,
			Status: 1,
		}

		err := userRepo.Insert(newUser)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证记录是否插入成功
		if newUser.ID == 0 {
			t.Error("Expected ID to be set after insertion")
		}

		// 查询插入的记录以确认存在
		insertedUser, err := userRepo.SelectOneById(int(newUser.ID))
		if err != nil {
			t.Errorf("Expected no error when querying inserted record, got: %v", err)
		}

		if insertedUser.Name != "newuser" {
			t.Errorf("Expected name to be 'newuser', got: %s", insertedUser.Name)
		}
	})

	t.Run("test insert with nil item", func(t *testing.T) {
		// 测试插入 nil 记录
		err := userRepo.Insert(nil)

		if err == nil {
			t.Error("Expected error for nil item")
		}

		if err.Error() != "item cannot be nil" {
			t.Errorf("Expected error message 'item cannot be nil', got: %v", err)
		}
	})

	t.Run("test insert with duplicate unique field", func(t *testing.T) {
		// 测试插入具有重复唯一字段的记录（假设 Email 是唯一字段）
		duplicateUser := &User{
			Name:   "duplicateuser",
			Email:  "testuser1@example.com", // 假设这个邮箱已存在
			Age:    30,
			Status: 1,
		}

		err := userRepo.Insert(duplicateUser)

		// 根据数据库约束，可能会返回错误
		if err == nil {
			t.Log("Insert succeeded, but expected potential constraint violation")
		}
	})
}

// TestBaseRepo_InsertBatch 测试 InsertBatch 方法
func TestBaseRepo_InsertBatch(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	// 创建 UserRepo 实例
	userRepo := &gormx.BaseRepo[User]{}

	t.Run("test insert batch with multiple records", func(t *testing.T) {
		// 测试批量插入多条记录
		users := []*User{
			{
				Name:   "batchuser1",
				Email:  "batchuser1@example.com",
				Age:    25,
				Status: 1,
			},
			{
				Name:   "batchuser2",
				Email:  "batchuser2@example.com",
				Age:    30,
				Status: 1,
			},
		}

		err := userRepo.InsertBatch(users)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证记录是否插入成功
		for _, user := range users {
			if user.ID == 0 {
				t.Error("Expected ID to be set after insertion")
			}

			// 查询插入的记录以确认存在
			insertedUser, err := userRepo.SelectOneById(int(user.ID))
			if err != nil {
				t.Errorf("Expected no error when querying inserted record, got: %v", err)
			}

			if insertedUser.Name != user.Name {
				t.Errorf("Expected name to be '%s', got: %s", user.Name, insertedUser.Name)
			}
		}
	})

	t.Run("test insert batch with empty slice", func(t *testing.T) {
		// 测试插入空切片
		var users []*User

		err := userRepo.InsertBatch(users)

		if err == nil {
			t.Error("Expected error for empty slice")
		}

		if err.Error() != "items cannot be empty" {
			t.Errorf("Expected error message 'items cannot be empty', got: %v", err)
		}
	})

	t.Run("test insert batch with nil slice", func(t *testing.T) {
		// 测试插入 nil 切片
		var users []*User

		err := userRepo.InsertBatch(users)

		if err == nil {
			t.Error("Expected error for nil slice")
		}

		if err.Error() != "items cannot be empty" {
			t.Errorf("Expected error message 'items cannot be empty', got: %v", err)
		}
	})

	t.Run("test insert batch with duplicate unique fields", func(t *testing.T) {
		// 测试批量插入具有重复唯一字段的记录（假设 Email 是唯一字段）
		users := []*User{
			{
				Name:   "duplicateuser1",
				Email:  "testuser1@example.com", // 假设这个邮箱已存在
				Age:    25,
				Status: 1,
			},
			{
				Name:   "duplicateuser2",
				Email:  "testuser2@example.com", // 假设这个邮箱也已存在
				Age:    30,
				Status: 1,
			},
		}

		err := userRepo.InsertBatch(users)

		// 根据数据库约束，可能会返回错误
		if err == nil {
			t.Log("InsertBatch succeeded, but expected potential constraint violation")
		}
	})
}

// TestBaseRepo_InsertInBatches 测试 InsertInBatches 方法
func TestBaseRepo_InsertInBatches(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	// 创建 UserRepo 实例
	userRepo := &gormx.BaseRepo[User]{}

	t.Run("test insert in batches with multiple records", func(t *testing.T) {
		// 测试按批次插入多条记录
		users := []*User{
			{
				Name:   "batchuser1",
				Email:  "batchuser1@example.com",
				Age:    25,
				Status: 1,
			},
			{
				Name:   "batchuser2",
				Email:  "batchuser2@example.com",
				Age:    30,
				Status: 1,
			},
			{
				Name:   "batchuser3",
				Email:  "batchuser3@example.com",
				Age:    35,
				Status: 1,
			},
		}

		batchSize := 2
		err := userRepo.InsertInBatches(users, batchSize)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证记录是否插入成功
		for _, user := range users {
			if user.ID == 0 {
				t.Error("Expected ID to be set after insertion")
			}

			// 查询插入的记录以确认存在
			insertedUser, err := userRepo.SelectOneById(int(user.ID))
			if err != nil {
				t.Errorf("Expected no error when querying inserted record, got: %v", err)
			}

			if insertedUser.Name != user.Name {
				t.Errorf("Expected name to be '%s', got: %s", user.Name, insertedUser.Name)
			}
		}
	})

	t.Run("test insert in batches with empty slice", func(t *testing.T) {
		// 测试插入空切片
		var users []*User

		batchSize := 2
		err := userRepo.InsertInBatches(users, batchSize)

		if err == nil {
			t.Error("Expected error for empty slice")
		}

		if err.Error() != "items cannot be empty" {
			t.Errorf("Expected error message 'items cannot be empty', got: %v", err)
		}
	})

	t.Run("test insert in batches with invalid batch size", func(t *testing.T) {
		// 测试无效的批次大小
		users := []*User{
			{
				Name:   "batchuser1",
				Email:  "batchuser1@example.com",
				Age:    25,
				Status: 1,
			},
		}

		batchSize := 0
		err := userRepo.InsertInBatches(users, batchSize)

		if err == nil {
			t.Error("Expected error for invalid batch size")
		}

		if err.Error() != "batchSize must be greater than 0" {
			t.Errorf("Expected error message 'batchSize must be greater than 0', got: %v", err)
		}
	})

	t.Run("test insert in batches with duplicate unique fields", func(t *testing.T) {
		// 测试按批次插入具有重复唯一字段的记录（假设 Email 是唯一字段）
		users := []*User{
			{
				Name:   "duplicateuser1",
				Email:  "testuser1@example.com", // 假设这个邮箱已存在
				Age:    25,
				Status: 1,
			},
			{
				Name:   "duplicateuser2",
				Email:  "testuser2@example.com", // 假设这个邮箱也已存在
				Age:    30,
				Status: 1,
			},
		}

		batchSize := 2
		err := userRepo.InsertInBatches(users, batchSize)

		// 根据数据库约束，可能会返回错误
		if err == nil {
			t.Log("InsertInBatches succeeded, but expected potential constraint violation")
		}
	})
}

// TestBaseRepo_InsertOrUpdate 测试 InsertOrUpdate 方法
func TestBaseRepo_InsertOrUpdate(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	// 创建 UserRepo 实例
	userRepo := &gormx.BaseRepo[User]{}

	t.Run("test insert or update with new record", func(t *testing.T) {
		// 测试插入新记录
		newUser := &User{
			Name:   "newuser",
			Email:  "newuser@example.com",
			Age:    25,
			Status: 1,
		}

		err := userRepo.InsertOrUpdate(newUser)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证记录是否插入成功
		if newUser.ID == 0 {
			t.Error("Expected ID to be set after insertion")
		}

		// 查询插入的记录以确认存在
		insertedUser, err := userRepo.SelectOneById(int(newUser.ID))
		if err != nil {
			t.Errorf("Expected no error when querying inserted record, got: %v", err)
		}

		if insertedUser.Name != "newuser" {
			t.Errorf("Expected name to be 'newuser', got: %s", insertedUser.Name)
		}
	})

	t.Run("test insert or update with existing record", func(t *testing.T) {
		// 测试更新已存在的记录
		existingUser := &User{
			Name:   "testuser1",
			Email:  "updateduser@example.com",
			Age:    30,
			Status: 1,
		}

		err := userRepo.InsertOrUpdate(existingUser)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 查询更新后的记录以确认更新成功
		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Name, "testuser1")
		query.Eq(&user.Age, 30)
		query.Eq(&user.Status, 1)

		updatedUser, err := userRepo.SelectOneByOpts(query.ToOptions()...)
		if err != nil {
			t.Errorf("Expected no error when querying updated record, got: %v", err)
		}

		if updatedUser.Email != "updateduser@example.com" {
			t.Errorf("Expected email to be 'updateduser@example.com', got: %s", updatedUser.Email)
		}
	})

	t.Run("test insert or update with nil item", func(t *testing.T) {
		// 测试传入 nil 记录
		err := userRepo.InsertOrUpdate(nil)

		// 根据实现，GORM 的 Create 方法可能会接受 nil，但这里我们期望返回错误
		if err == nil {
			t.Log("InsertOrUpdate succeeded with nil item, but expected potential error")
		}
	})
}

// TestBaseRepo_Update 测试 Update 方法
func TestBaseRepo_Update(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	// 创建 UserRepo 实例
	userRepo := &gormx.BaseRepo[User]{}

	t.Run("test update an existing record", func(t *testing.T) {
		// 测试更新已存在的记录
		// 先查询一条记录
		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Name, "testuser1")

		existingUser, err := userRepo.SelectOneByOpts(query.ToOptions()...)
		if err != nil {
			t.Errorf("Expected no error when querying existing record, got: %v", err)
		}

		// 修改记录字段
		existingUser.Email = "updated@example.com"
		existingUser.Age = 35

		// 执行更新
		err = userRepo.Update(&existingUser)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 查询更新后的记录以确认更新成功
		updatedUser, err := userRepo.SelectOneById(int(existingUser.ID))
		if err != nil {
			t.Errorf("Expected no error when querying updated record, got: %v", err)
		}

		if updatedUser.Email != "updated@example.com" {
			t.Errorf("Expected email to be 'updated@example.com', got: %s", updatedUser.Email)
		}

		if updatedUser.Age != 35 {
			t.Errorf("Expected age to be 35, got: %d", updatedUser.Age)
		}
	})

	t.Run("test update with nil item", func(t *testing.T) {
		// 测试传入 nil 记录
		err := userRepo.Update(nil)

		if err == nil {
			t.Error("Expected error for nil item")
		}

		if err.Error() != "item cannot be nil" {
			t.Errorf("Expected error message 'item cannot be nil', got: %v", err)
		}
	})
}

// TestBaseRepo_UpdateById 测试 UpdateById 方法
func TestBaseRepo_UpdateById(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	// 创建 UserRepo 实例
	userRepo := &gormx.BaseRepo[User]{}

	t.Run("test update by id with existing record", func(t *testing.T) {
		// 测试通过ID更新已存在的记录
		// 先查询一条记录获取其ID
		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Name, "testuser1")

		existingUser, err := userRepo.SelectOneByOpts(query.ToOptions()...)
		if err != nil {
			t.Errorf("Expected no error when querying existing record, got: %v", err)
		}

		// 定义要更新的字段
		updateVars := map[string]interface{}{
			"email": "updated@example.com",
			"age":   35,
		}

		// 执行更新
		err = userRepo.UpdateById(int(existingUser.ID), updateVars)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 查询更新后的记录以确认更新成功
		updatedUser, err := userRepo.SelectOneById(int(existingUser.ID))
		if err != nil {
			t.Errorf("Expected no error when querying updated record, got: %v", err)
		}

		if updatedUser.Email != "updated@example.com" {
			t.Errorf("Expected email to be 'updated@example.com', got: %s", updatedUser.Email)
		}

		if updatedUser.Age != 35 {
			t.Errorf("Expected age to be 35, got: %d", updatedUser.Age)
		}
	})

	t.Run("test update by id with non-existent id", func(t *testing.T) {
		// 测试通过不存在的ID更新记录
		nonExistentId := 999999
		updateVars := map[string]interface{}{
			"email": "nonexistent@example.com",
		}

		err := userRepo.UpdateById(nonExistentId, updateVars)

		// 根据实现，GORM 可能不会返回错误，但记录不会被更新
		if err != nil {
			t.Logf("UpdateById returned error for non-existent ID: %v", err)
		}
	})

	t.Run("test update by id with restricted fields", func(t *testing.T) {
		// 测试更新受限字段（如ID、created_at等）
		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Name, "testuser1")

		existingUser, err := userRepo.SelectOneByOpts(query.ToOptions()...)
		if err != nil {
			t.Errorf("Expected no error when querying existing record, got: %v", err)
		}

		// 尝试更新受限字段
		updateVars := map[string]interface{}{
			"id":         999,
			"created_at": "2023-01-01",
			"email":      "restricted@example.com",
		}

		err = userRepo.UpdateById(int(existingUser.ID), updateVars)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 查询更新后的记录以确认受限字段未被更新
		updatedUser, err := userRepo.SelectOneById(int(existingUser.ID))
		if err != nil {
			t.Errorf("Expected no error when querying updated record, got: %v", err)
		}

		// 验证受限字段未被更新
		if updatedUser.ID != existingUser.ID {
			t.Errorf("Expected ID to remain unchanged, got: %d", updatedUser.ID)
		}

		if updatedUser.Email != "restricted@example.com" {
			t.Errorf("Expected email to be updated, got: %s", updatedUser.Email)
		}
	})
}

// TestBaseRepo_UpdateByOpts 测试 UpdateByOpts 方法
func TestBaseRepo_UpdateByOpts(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	// 创建 UserRepo 实例
	userRepo := &gormx.BaseRepo[User]{}

	t.Run("test update by opts with existing records", func(t *testing.T) {
		// 测试通过条件更新已存在的记录
		// 定义要更新的字段
		updateVars := map[string]interface{}{
			"email": "updated@example.com",
			"age":   35,
		}

		// 定义更新条件
		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Status, 1)

		// 执行更新
		err := userRepo.UpdateByOpts(updateVars, query.ToOptions()...)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 查询更新后的记录以确认更新成功
		users, err := userRepo.SelectListByOpts(query.ToOptions()...)
		if err != nil {
			t.Errorf("Expected no error when querying updated records, got: %v", err)
		}

		for _, user := range users {
			if user.Email != "updated@example.com" {
				t.Errorf("Expected email to be 'updated@example.com', got: %s", user.Email)
			}

			if user.Age != 35 {
				t.Errorf("Expected age to be 35, got: %d", user.Age)
			}
		}
	})

	t.Run("test update by opts with no conditions", func(t *testing.T) {
		// 测试无条件更新记录（应返回错误）
		updateVars := map[string]interface{}{
			"email": "nocondition@example.com",
		}

		err := userRepo.UpdateByOpts(updateVars)

		if err == nil {
			t.Error("Expected error for no conditions")
		}

		if err.Error() != "cannot update records without conditions, please provide where conditions" {
			t.Errorf("Expected error message 'cannot update records without conditions, please provide where conditions', got: %v", err)
		}
	})

	t.Run("test update by opts with non-matching conditions", func(t *testing.T) {
		// 测试通过不匹配的条件更新记录
		updateVars := map[string]interface{}{
			"email": "nomatch@example.com",
		}

		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Name, "nonexistent_user")

		err := userRepo.UpdateByOpts(updateVars, query.ToOptions()...)

		// 根据实现，GORM 可能不会返回错误，但记录不会被更新
		if err != nil {
			t.Logf("UpdateByOpts returned error for non-matching conditions: %v", err)
		}
	})

	t.Run("test update by opts with restricted fields", func(t *testing.T) {
		// 测试更新受限字段（如ID、created_at等）
		updateVars := map[string]interface{}{
			"id":         999,
			"created_at": "2023-01-01",
			"email":      "restricted@example.com",
		}

		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Status, 1)

		err := userRepo.UpdateByOpts(updateVars, query.ToOptions()...)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 查询更新后的记录以确认受限字段未被更新
		users, err := userRepo.SelectListByOpts(query.ToOptions()...)
		if err != nil {
			t.Errorf("Expected no error when querying updated records, got: %v", err)
		}

		for _, user := range users {
			// 验证受限字段未被更新
			if user.ID == 999 {
				t.Errorf("Expected ID to remain unchanged, got: %d", user.ID)
			}

			if user.Email != "restricted@example.com" {
				t.Errorf("Expected email to be updated, got: %s", user.Email)
			}
		}
	})
}

// TestBaseRepo_Upsert 测试 Upsert 方法
func TestBaseRepo_Upsert(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	// 创建 UserRepo 实例
	userRepo := &gormx.BaseRepo[User]{}

	t.Run("test upsert with new record", func(t *testing.T) {
		// 测试插入新记录
		newUser := &User{
			Name:   "newuser",
			Email:  "newuser@example.com",
			Age:    25,
			Status: 1,
		}

		// 定义更新字段（仅在冲突时使用）
		updateVars := map[string]interface{}{
			"email": "updated@example.com",
		}

		err := userRepo.Upsert(newUser, updateVars)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证记录是否插入成功
		if newUser.ID == 0 {
			t.Error("Expected ID to be set after upsert")
		}

		// 查询插入的记录以确认存在
		insertedUser, err := userRepo.SelectOneById(int(newUser.ID))
		if err != nil {
			t.Errorf("Expected no error when querying inserted record, got: %v", err)
		}

		if insertedUser.Name != "newuser" {
			t.Errorf("Expected name to be 'newuser', got: %s", insertedUser.Name)
		}
	})

	t.Run("test upsert with existing record", func(t *testing.T) {
		// 查询已存在的用户记录
		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Name, "testuser1")
		existingUser, err := userRepo.SelectOneByOpts(query.ToOptions()...)
		if err != nil {
			t.Fatalf("Failed to query existing user: %v", err)
		}

		// 定义更新字段（仅在冲突时使用）
		updateVars := map[string]interface{}{
			"email": "conflict-updated@example.com",
		}

		err = userRepo.Upsert(&existingUser, updateVars)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 查询更新后的记录以确认更新成功
		updatedUser, err := userRepo.SelectOneById(int(existingUser.ID))
		if err != nil {
			t.Errorf("Expected no error when querying updated record, got: %v", err)
		}

		if updatedUser.Email != "conflict-updated@example.com" {
			t.Errorf("Expected email to be 'conflict-updated@example.com', got: %s", updatedUser.Email)
		}
	})

	t.Run("test upsert with nil item", func(t *testing.T) {
		// 测试传入 nil 记录
		var user *User
		updateVars := map[string]interface{}{
			"email": "nil@example.com",
		}

		err := userRepo.Upsert(user, updateVars)

		// 根据实现，GORM 的 Create 方法可能会接受 nil，但这里我们期望返回错误
		if err == nil {
			t.Log("Upsert succeeded with nil item, but expected potential error")
		}
	})
}

// TestBaseRepo_GetOrCreate 测试 GetOrCreate 方法
func TestBaseRepo_GetOrCreate(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	// 创建 UserRepo 实例
	userRepo := &gormx.BaseRepo[User]{}

	t.Run("test get or create with existing record", func(t *testing.T) {
		// 测试获取已存在的记录
		whereCond := map[string]interface{}{
			"name": "testuser1",
		}
		assignAttrs := map[string]interface{}{
			"email": "updated@example.com",
		}

		user, err := userRepo.GetOrCreate(whereCond, assignAttrs)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if user.Name != "testuser1" {
			t.Errorf("Expected name to be 'testuser1', got: %s", user.Name)
		}

		// 验证记录未被更新（因为记录已存在）
		if user.Email == "updated@example.com" {
			t.Errorf("Expected email to remain unchanged, got: %s", user.Email)
		}
	})

	t.Run("test get or create with new record", func(t *testing.T) {
		// 测试创建新记录
		whereCond := map[string]interface{}{
			"name": "newuser",
		}
		assignAttrs := map[string]interface{}{
			"email":  "newuser@example.com",
			"age":    25,
			"status": 1,
		}

		user, err := userRepo.GetOrCreate(whereCond, assignAttrs)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if user.Name != "newuser" {
			t.Errorf("Expected name to be 'newuser', got: %s", user.Name)
		}

		if user.Email != "newuser@example.com" {
			t.Errorf("Expected email to be 'newuser@example.com', got: %s", user.Email)
		}

		if user.Age != 25 {
			t.Errorf("Expected age to be 25, got: %d", user.Age)
		}

		if user.Status != 1 {
			t.Errorf("Expected status to be 1, got: %d", user.Status)
		}

		// 验证记录是否插入成功
		if user.ID == 0 {
			t.Error("Expected ID to be set after creation")
		}
	})

	t.Run("test get or create with empty where condition", func(t *testing.T) {
		// 测试空查询条件（应返回错误）
		whereCond := map[string]interface{}{}
		assignAttrs := map[string]interface{}{
			"name":   "emptyconditionuser",
			"email":  "empty@example.com",
			"age":    30,
			"status": 1,
		}

		_, err := userRepo.GetOrCreate(whereCond, assignAttrs)

		if err == nil {
			t.Error("Expected error for no conditions")
		}

		if err.Error() != "cannot get or create record without conditions, please provide where conditions" {
			t.Errorf("Expected error message 'cannot get or create record without conditions, please provide where conditions', got: %v", err)
		}

		//if err != nil {
		//	t.Errorf("Expected no error, got: %v", err)
		//}
		//
		//if user.Name != "emptyconditionuser" {
		//	t.Errorf("Expected name to be 'emptyconditionuser', got: %s", user.Name)
		//}
		//
		//if user.Email != "empty@example.com" {
		//	t.Errorf("Expected email to be 'empty@example.com', got: %s", user.Email)
		//}
		//
		//if user.Age != 30 {
		//	t.Errorf("Expected age to be 30, got: %d", user.Age)
		//}
		//
		//if user.Status != 1 {
		//	t.Errorf("Expected status to be 1, got: %d", user.Status)
		//}
		//
		//// 验证记录是否插入成功
		//if user.ID == 0 {
		//	t.Error("Expected ID to be set after creation")
		//}
	})

	t.Run("test get or create with conflicting unique field", func(t *testing.T) {
		// 测试创建具有冲突唯一字段的记录（假设 Email 是唯一字段）
		whereCond := map[string]interface{}{
			"email": "testuser1@example.com", // 假设这个邮箱已存在
		}
		assignAttrs := map[string]interface{}{
			"name":   "conflictuser",
			"age":    25,
			"status": 1,
		}

		user, err := userRepo.GetOrCreate(whereCond, assignAttrs)

		// 根据数据库约束，可能会返回错误
		if err != nil {
			t.Logf("GetOrCreate returned error for conflicting unique field: %v", err)
		} else {
			// 如果没有错误，验证返回的是已存在的记录
			if user.Email != "testuser1@example.com" {
				t.Errorf("Expected email to be 'testuser1@example.com', got: %s", user.Email)
			}
		}
	})
}

func TestBaseRepo_UpdateOrCreate(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	// 创建 UserRepo 实例
	userRepo := &gormx.BaseRepo[User]{}

	t.Run("test update or create with existing record", func(t *testing.T) {
		// 测试更新已存在的记录
		whereCond := map[string]interface{}{
			"name": "testuser1",
		}
		assignAttrs := map[string]interface{}{
			"email": "updated@example.com",
			"age":   35,
		}

		user, err := userRepo.UpdateOrCreate(whereCond, assignAttrs)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证记录是否更新成功
		if user.Name != "testuser1" {
			t.Errorf("Expected name to be 'testuser1', got: %s", user.Name)
		}

		if user.Email != "updated@example.com" {
			t.Errorf("Expected email to be 'updated@example.com', got: %s", user.Email)
		}

		if user.Age != 35 {
			t.Errorf("Expected age to be 35, got: %d", user.Age)
		}
	})

	t.Run("test update or create with new record", func(t *testing.T) {
		// 测试创建新记录
		whereCond := map[string]interface{}{
			"name": "newuser",
		}
		assignAttrs := map[string]interface{}{
			"email":  "newuser@example.com",
			"age":    25,
			"status": 1,
		}

		user, err := userRepo.UpdateOrCreate(whereCond, assignAttrs)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证记录是否创建成功
		if user.Name != "newuser" {
			t.Errorf("Expected name to be 'newuser', got: %s", user.Name)
		}

		if user.Email != "newuser@example.com" {
			t.Errorf("Expected email to be 'newuser@example.com', got: %s", user.Email)
		}

		if user.Age != 25 {
			t.Errorf("Expected age to be 25, got: %d", user.Age)
		}

		if user.Status != 1 {
			t.Errorf("Expected status to be 1, got: %d", user.Status)
		}

		// 验证记录是否插入成功
		if user.ID == 0 {
			t.Error("Expected ID to be set after creation")
		}
	})

	t.Run("test update or create with empty where condition", func(t *testing.T) {
		// 测试空查询条件（应创建新记录）
		whereCond := map[string]interface{}{}
		assignAttrs := map[string]interface{}{
			"name":   "emptyconditionuser",
			"email":  "empty@example.com",
			"age":    30,
			"status": 1,
		}

		user, err := userRepo.UpdateOrCreate(whereCond, assignAttrs)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证记录是否创建成功
		if user.Name != "emptyconditionuser" {
			t.Errorf("Expected name to be 'emptyconditionuser', got: %s", user.Name)
		}

		if user.Email != "empty@example.com" {
			t.Errorf("Expected email to be 'empty@example.com', got: %s", user.Email)
		}

		if user.Age != 30 {
			t.Errorf("Expected age to be 30, got: %d", user.Age)
		}

		if user.Status != 1 {
			t.Errorf("Expected status to be 1, got: %d", user.Status)
		}

		// 验证记录是否插入成功
		if user.ID == 0 {
			t.Error("Expected ID to be set after creation")
		}
	})

	t.Run("test update or create with conflicting unique field", func(t *testing.T) {
		// 测试创建具有冲突唯一字段的记录（假设 Email 是唯一字段）
		whereCond := map[string]interface{}{
			"email": "testuser1@example.com", // 假设这个邮箱已存在
		}
		assignAttrs := map[string]interface{}{
			"name":   "conflictuser",
			"age":    25,
			"status": 1,
		}

		user, err := userRepo.UpdateOrCreate(whereCond, assignAttrs)

		// 根据数据库约束，可能会返回错误
		if err != nil {
			t.Logf("UpdateOrCreate returned error for conflicting unique field: %v", err)
		} else {
			// 如果没有错误，验证返回的是已存在的记录
			if user.Email != "testuser1@example.com" {
				t.Errorf("Expected email to be 'testuser1@example.com', got: %s", user.Email)
			}
		}
	})
}

func TestBaseRepo_Delete(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	// 创建 UserRepo 实例
	userRepo := &gormx.BaseRepo[User]{}

	t.Run("test delete existing record", func(t *testing.T) {
		// 插入一条记录用于删除
		newUser := &User{
			Name:   "deletetestuser",
			Email:  "delete@example.com",
			Age:    25,
			Status: 1,
		}
		err := userRepo.Insert(newUser)
		if err != nil {
			t.Fatalf("Failed to insert test user: %v", err)
		}

		// 执行删除
		err = userRepo.Delete(newUser)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证记录是否被删除
		_, err = userRepo.SelectOneById(int(newUser.ID))
		if err == nil {
			t.Error("Expected record to be deleted, but it still exists")
		}
	})

	t.Run("test delete non-existent record", func(t *testing.T) {
		// 尝试删除一个不存在的记录
		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Name, "non_existent_user")
		nonExistentUser, _ := userRepo.SelectOneByOpts(query.ToOptions()...)

		err := userRepo.Delete(&nonExistentUser)

		// 根据实现，GORM 可能不会返回错误，但记录不会被删除
		if err != nil {
			t.Logf("Delete returned error for non-existent record: %v", err)
		}
	})

	t.Run("test delete with nil item", func(t *testing.T) {
		// 测试传入 nil 记录
		err := userRepo.Delete(nil)

		if err == nil {
			t.Error("Expected error for nil item")
		}

		if err.Error() != "item cannot be nil" {
			t.Errorf("Expected error message 'item cannot be nil', got: %v", err)
		}
	})
}

func TestBaseRepo_DeleteById(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	// 创建 UserRepo 实例
	userRepo := &gormx.BaseRepo[User]{}

	t.Run("test delete by id with existing record", func(t *testing.T) {
		// 插入一条记录用于删除
		newUser := &User{
			Name:   "deletebyiduser",
			Email:  "deletebyid@example.com",
			Age:    30,
			Status: 1,
		}
		err := userRepo.Insert(newUser)
		if err != nil {
			t.Fatalf("Failed to insert test user: %v", err)
		}

		// 执行删除
		err = userRepo.DeleteById(int(newUser.ID))
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证记录是否被删除
		_, err = userRepo.SelectOneById(int(newUser.ID))
		if err == nil {
			t.Error("Expected record to be deleted, but it still exists")
		}
	})

	t.Run("test delete by id with non-existent id", func(t *testing.T) {
		// 尝试删除一个不存在的 ID
		nonExistentId := 999999
		err := userRepo.DeleteById(nonExistentId)

		// 根据实现，GORM 可能不会返回错误，但记录不会被删除
		if err != nil {
			t.Logf("DeleteById returned error for non-existent ID: %v", err)
		}
	})
}

func TestBaseRepo_DeleteBatchIds(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	// 创建 UserRepo 实例
	userRepo := &gormx.BaseRepo[User]{}

	t.Run("test delete batch ids with existing records", func(t *testing.T) {
		// 插入多条记录用于删除
		users := []*User{
			{Name: "batchuser1", Email: "batch1@example.com", Age: 25, Status: 1},
			{Name: "batchuser2", Email: "batch2@example.com", Age: 30, Status: 1},
		}
		for _, user := range users {
			err := userRepo.Insert(user)
			if err != nil {
				t.Fatalf("Failed to insert test user: %v", err)
			}
		}

		// 获取插入记录的 IDs
		var ids []int
		for _, user := range users {
			ids = append(ids, int(user.ID))
		}

		// 执行批量删除
		err := userRepo.DeleteBatchIds(ids)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证记录是否被删除
		for _, id := range ids {
			_, err := userRepo.SelectOneById(id)
			if err == nil {
				t.Errorf("Expected record with ID %d to be deleted, but it still exists", id)
			}
		}
	})

	t.Run("test delete batch ids with non-existent ids", func(t *testing.T) {
		// 尝试删除不存在的 IDs
		nonExistentIds := []int{999997, 999998, 999999}
		err := userRepo.DeleteBatchIds(nonExistentIds)

		// 根据实现，GORM 可能不会返回错误，但记录不会被删除
		if err != nil {
			t.Logf("DeleteBatchIds returned error for non-existent IDs: %v", err)
		}
	})
}

func TestBaseRepo_DeleteByOpts(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	// 创建 UserRepo 实例
	userRepo := &gormx.BaseRepo[User]{}

	t.Run("test delete by opts with existing records", func(t *testing.T) {
		// 插入多条记录用于删除
		users := []*User{
			{Name: "deletebyoptsuser1", Email: "opts1@example.com", Age: 25, Status: 1},
			{Name: "deletebyoptsuser2", Email: "opts2@example.com", Age: 30, Status: 1},
		}
		for _, user := range users {
			err := userRepo.Insert(user)
			if err != nil {
				t.Fatalf("Failed to insert test user: %v", err)
			}
		}

		// 定义删除条件
		query, user := gormx.NewQuery[User]()
		query.Eq(&user.Status, 1)

		// 执行删除
		err := userRepo.DeleteByOpts(query.ToOptions()...)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证记录是否被删除
		remainingUsers, err := userRepo.SelectListByOpts(query.ToOptions()...)
		if err != nil {
			t.Errorf("Expected no error when querying remaining records, got: %v", err)
		}
		if len(remainingUsers) != 0 {
			t.Errorf("Expected all records with status=1 to be deleted, but found %d remaining", len(remainingUsers))
		}
	})

	t.Run("test delete by opts with no conditions", func(t *testing.T) {
		// 测试无条件删除（应返回错误）
		err := userRepo.DeleteByOpts()

		if err == nil {
			t.Error("Expected error for no conditions")
		}

		if err.Error() != "delete operation requires where conditions" {
			t.Errorf("Expected error message 'delete operation requires where conditions', got: %v", err)
		}
	})
}

func TestBaseRepo_DeleteByMap(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	// 创建 UserRepo 实例
	userRepo := &gormx.BaseRepo[User]{}

	t.Run("test delete by map with existing records", func(t *testing.T) {
		// 插入多条记录用于删除
		users := []*User{
			{Name: "deletebymapuser1", Email: "map1@example.com", Age: 25, Status: 1},
			{Name: "deletebymapuser2", Email: "map2@example.com", Age: 30, Status: 1},
		}
		for _, user := range users {
			err := userRepo.Insert(user)
			if err != nil {
				t.Fatalf("Failed to insert test user: %v", err)
			}
		}

		// 定义删除条件
		columnMap := map[string]interface{}{
			"status": 1,
		}

		// 执行删除
		err := userRepo.DeleteByMap(columnMap)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证记录是否被删除
		remainingUsers, err := userRepo.SelectListByMap(columnMap)
		if err != nil {
			t.Errorf("Expected no error when querying remaining records, got: %v", err)
		}
		if len(remainingUsers) != 0 {
			t.Errorf("Expected all records with status=1 to be deleted, but found %d remaining", len(remainingUsers))
		}
	})

	t.Run("test delete by map with empty map", func(t *testing.T) {
		// 测试空 map 删除（应返回错误）
		columnMap := map[string]interface{}{}
		err := userRepo.DeleteByMap(columnMap)

		if err == nil {
			t.Error("Expected error for empty map")
		}

		if err.Error() != "cannot delete records without conditions" {
			t.Errorf("Expected error message 'cannot delete records without conditions', got: %v", err)
		}
	})
}
