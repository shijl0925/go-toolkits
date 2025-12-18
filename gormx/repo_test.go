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

	//t.Run("test select one record by options with joins", func(t *testing.T) {
	//	// 测试带JOIN的查询
	//	query, user := gormx.NewQuery[User]()
	//	query.Eq(&user.Status, 1).Eq(&user.Email, "test1@example.com").
	//		Join("LEFT JOIN posts ON users.id = posts.user_id").
	//		Select(&user.ID, &user.Name)
	//
	//	result, err := repo.SelectOneByOpts(query.ToOptions()...)
	//
	//	if err != nil {
	//		t.Errorf("Expected no error, got: %v", err)
	//	}
	//
	//	if result.Status != 1 {
	//		t.Errorf("Expected status: 1, got: %d", result.Status)
	//	}
	//})

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
			if int(result.Status) == status {
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
