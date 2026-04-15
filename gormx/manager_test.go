package gormx_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/shijl0925/go-toolkits/gormx"
)

// TestGenericAssociationManager_Add 测试 Add 方法
func TestGenericAssociationManager_Add(t *testing.T) {
	// 准备测试数据
	setupTestData(t)
	setupPreloadTestData(t)

	userRepo := &gormx.BaseRepo[User]{}
	uq, u := gormx.NewQuery[User]()
	uq.Eq(&u.Name, "testuser1")
	user, _ := userRepo.SelectOneByOpts(uq.ToOptions()...)
	//roleManager := gormx.NewAssociationManager[User, Role](user, "Roles")

	roleRepo := &gormx.BaseRepo[Role]{}
	rq, r := gormx.NewQuery[Role]()
	rq.Eq(&r.Name, "Admin")
	role, _ := roleRepo.SelectOneByOpts(rq.ToOptions()...)

	err := user.RolesManager().Add(role)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// 验证关联是否添加成功
	associatedRoles, err := user.RolesManager().All()
	fmt.Println(associatedRoles)
	if err != nil {
		t.Errorf("Expected no error when fetching associated roles, got: %v", err)
	}

	found := false
	for _, r := range associatedRoles {
		if r.Name == "Admin" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected role 'admin' to be associated with user")
	}
}

// TestGenericAssociationManager_Remove 测试 Remove 方法
func TestGenericAssociationManager_Remove(t *testing.T) {
	// 准备测试数据
	setupTestData(t)
	setupPreloadTestData(t)

	userRepo := &gormx.BaseRepo[User]{}
	uq, u := gormx.NewQuery[User]()
	uq.Eq(&u.Name, "testuser1")
	user, _ := userRepo.SelectOneByOpts(uq.ToOptions()...)
	//roleManager := gormx.NewAssociationManager[User, Role](user, "Roles")

	// 先添加一个角色
	roleRepo := &gormx.BaseRepo[Role]{}
	rq, r := gormx.NewQuery[Role]()
	rq.Eq(&r.Name, "Admin")
	role, _ := roleRepo.SelectOneByOpts(rq.ToOptions()...)

	err := user.RolesManager().Add(role)
	if err != nil {
		t.Fatalf("Failed to add role: %v", err)
	}

	// 然后删除该角色
	err = user.RolesManager().Remove(role)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// 验证关联是否删除成功
	associatedRoles, err := user.RolesManager().All()
	if err != nil {
		t.Errorf("Expected no error when fetching associated roles, got: %v", err)
	}

	for _, r := range associatedRoles {
		if r.Name == "Admin" {
			t.Error("Expected role 'admin' to be removed from user")
		}
	}
}

// TestGenericAssociationManager_Clear 测试 Clear 方法
func TestGenericAssociationManager_Clear(t *testing.T) {
	// 准备测试数据
	setupTestData(t)
	setupPreloadTestData(t)

	//user := User{Name: "testuser"}
	userRepo := &gormx.BaseRepo[User]{}
	uq, u := gormx.NewQuery[User]()
	uq.Eq(&u.Name, "testuser1")
	user, _ := userRepo.SelectOneByOpts(uq.ToOptions()...)
	//roleManager := gormx.NewAssociationManager[User, Role](user, "Roles")

	// 添加多个角色
	roleRepo := &gormx.BaseRepo[Role]{}
	rq, r := gormx.NewQuery[Role]()
	rq.In(&r.Name, []string{"Admin", "Editor"})
	roles, _ := roleRepo.SelectListByOpts(rq.ToOptions()...)

	err := user.RolesManager().Add(roles...)
	if err != nil {
		t.Fatalf("Failed to add roles: %v", err)
	}

	// 清空所有关联
	err = user.RolesManager().Clear()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// 验证关联是否清空成功
	count, err := user.RolesManager().Count()
	if err != nil {
		t.Errorf("Expected no error when counting associated roles, got: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 associated roles, got: %d", count)
	}
}

// TestGenericAssociationManager_Set 测试 Set 方法
func TestGenericAssociationManager_Set(t *testing.T) {
	// 准备测试数据
	setupTestData(t)
	setupPreloadTestData(t)

	userRepo := &gormx.BaseRepo[User]{}
	uq, u := gormx.NewQuery[User]()
	uq.Eq(&u.Name, "testuser1")
	user, _ := userRepo.SelectOneByOpts(uq.ToOptions()...)
	//roleManager := gormx.NewAssociationManager[User, Role](user, "Roles")

	// 设置新的关联角色
	roleRepo := &gormx.BaseRepo[Role]{}
	rq, r := gormx.NewQuery[Role]()
	rq.In(&r.Name, []string{"Admin", "Editor"})
	newRoles, _ := roleRepo.SelectListByOpts(rq.ToOptions()...)

	err := user.RolesManager().Set(newRoles)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// 验证关联是否设置成功
	associatedRoles, err := user.RolesManager().All()
	if err != nil {
		t.Errorf("Expected no error when fetching associated roles, got: %v", err)
	}

	expectedRoles := map[string]bool{
		"Admin":  true,
		"Editor": true,
	}

	for _, r := range associatedRoles {
		if !expectedRoles[r.Name] {
			t.Errorf("Unexpected role found: %s", r.Name)
		}
		delete(expectedRoles, r.Name)
	}

	if len(expectedRoles) > 0 {
		t.Errorf("Expected roles not found: %v", expectedRoles)
	}
}

// TestGenericAssociationManager_Count 测试 Count 方法
func TestGenericAssociationManager_Count(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	userRepo := &gormx.BaseRepo[User]{}
	uq, u := gormx.NewQuery[User]()
	uq.Eq(&u.Name, "testuser1")
	user, _ := userRepo.SelectOneByOpts(uq.ToOptions()...)
	//roleManager := gormx.NewAssociationManager[User, Role](user, "Roles")

	// 添加多个角色
	roleRepo := &gormx.BaseRepo[Role]{}
	rq, r := gormx.NewQuery[Role]()
	rq.In(&r.Name, []string{"Admin", "Editor"})
	roles, _ := roleRepo.SelectListByOpts(rq.ToOptions()...)

	err := user.RolesManager().Add(roles...)
	if err != nil {
		t.Fatalf("Failed to add roles: %v", err)
	}

	// 统计关联数量
	count, err := user.RolesManager().Count()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected 2 associated roles, got: %d", count)
	}
}

// TestGenericAssociationManager_All 测试 All 方法
func TestGenericAssociationManager_All(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

	userRepo := &gormx.BaseRepo[User]{}
	uq, u := gormx.NewQuery[User]()
	uq.Eq(&u.Name, "testuser1")
	user, _ := userRepo.SelectOneByOpts(uq.ToOptions()...)
	//roleManager := gormx.NewAssociationManager[User, Role](user, "Roles")

	// 添加多个角色
	roleRepo := &gormx.BaseRepo[Role]{}
	rq, r := gormx.NewQuery[Role]()
	rq.In(&r.Name, []string{"Admin", "Editor"})
	roles, _ := roleRepo.SelectListByOpts(rq.ToOptions()...)
	err := user.RolesManager().Add(roles...)
	if err != nil {
		t.Fatalf("Failed to add roles: %v", err)
	}

	// 获取所有关联角色
	associatedRoles, err := user.RolesManager().All()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(associatedRoles) != 2 {
		t.Errorf("Expected 2 associated roles, got: %d", len(associatedRoles))
	}

	// 验证返回的角色
	expectedRoles := map[string]bool{
		"Admin":  true,
		"Editor": true,
	}

	for _, r := range associatedRoles {
		if !expectedRoles[r.Name] {
			t.Errorf("Unexpected role found: %s", r.Name)
		}
		delete(expectedRoles, r.Name)
	}

	if len(expectedRoles) > 0 {
		t.Errorf("Expected roles not found: %v", expectedRoles)
	}
}

// --- DB 注入 / 事务测试 ---

// getTestUser 返回名为 "testuser1" 的用户，若不存在或存在多条同名记录则令测试失败。
// 唯一性由 SelectOneByOpts 保证：超过一条匹配时会返回错误。
func getTestUser(t *testing.T) User {
	t.Helper()
	userRepo := &gormx.BaseRepo[User]{}
	uq, u := gormx.NewQuery[User]()
	uq.Eq(&u.Name, "testuser1")
	user, err := userRepo.SelectOneByOpts(uq.ToOptions()...)
	if err != nil {
		t.Fatalf("getTestUser: %v", err)
	}
	return user
}

// getTestRole 返回指定名称的角色，若不存在或存在多条同名记录则令测试失败。
// 唯一性由 SelectOneByOpts 保证：超过一条匹配时会返回错误。
func getTestRole(t *testing.T, name string) Role {
	t.Helper()
	roleRepo := &gormx.BaseRepo[Role]{}
	rq, r := gormx.NewQuery[Role]()
	rq.Eq(&r.Name, name)
	role, err := roleRepo.SelectOneByOpts(rq.ToOptions()...)
	if err != nil {
		t.Fatalf("getTestRole(%q): %v", name, err)
	}
	return role
}

// TestAssociationManager_Add_WithTransaction_Rollback 验证在事务中通过 UseDB 注入的
// AssociationManager.Add 操作会在事务回滚后还原：
//  1. 在事务内添加关联，关联在事务内可见。
//  2. 回滚事务后，关联从全局 DB 中消失。
func TestAssociationManager_Add_WithTransaction_Rollback(t *testing.T) {
	setupTestData(t)
	setupPreloadTestData(t)

	user := getTestUser(t)
	role := getTestRole(t, "Viewer")

	tx := gormx.GetDb().Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}

	// 在事务内添加关联
	txManager := gormx.NewAssociationManager[User, Role](user, "Roles", gormx.UseDB(tx))
	if err := txManager.Add(role); err != nil {
		tx.Rollback()
		t.Fatalf("Add inside tx failed: %v", err)
	}

	// 事务内应能看到该关联
	inTxRoles, err := txManager.All()
	if err != nil {
		tx.Rollback()
		t.Fatalf("All inside tx failed: %v", err)
	}
	found := false
	for _, r := range inTxRoles {
		if r.Name == "Viewer" {
			found = true
			break
		}
	}
	if !found {
		tx.Rollback()
		t.Fatal("expected role 'Viewer' to be visible inside tx")
	}

	// 回滚
	if err := tx.Rollback().Error; err != nil {
		t.Fatalf("rollback failed: %v", err)
	}

	// 回滚后全局 DB 中不应有该关联
	globalManager := gormx.NewAssociationManager[User, Role](user, "Roles")
	afterRoles, err := globalManager.All()
	if err != nil {
		t.Fatalf("All after rollback failed: %v", err)
	}
	for _, r := range afterRoles {
		if r.Name == "Viewer" {
			t.Error("expected role 'Viewer' to be absent after rollback, but it was found")
		}
	}
}

// TestAssociationManager_Add_WithTransaction_Commit 验证在事务提交后，
// 通过 UseDB 注入的 AssociationManager.Add 操作结果在全局 DB 中持久可见。
func TestAssociationManager_Add_WithTransaction_Commit(t *testing.T) {
	setupTestData(t)
	setupPreloadTestData(t)

	user := getTestUser(t)
	role := getTestRole(t, "Viewer")

	tx := gormx.GetDb().Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}

	txManager := gormx.NewAssociationManager[User, Role](user, "Roles", gormx.UseDB(tx))
	if err := txManager.Add(role); err != nil {
		tx.Rollback()
		t.Fatalf("Add inside tx failed: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		t.Fatalf("commit failed: %v", err)
	}

	// 提交后通过全局 DB 应能看到该关联
	globalManager := gormx.NewAssociationManager[User, Role](user, "Roles")
	afterRoles, err := globalManager.All()
	if err != nil {
		t.Fatalf("All after commit failed: %v", err)
	}

	found := false
	for _, r := range afterRoles {
		if r.Name == "Viewer" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected role 'Viewer' to be visible after commit, but not found")
	}

	// 清理：只移除本测试新增的角色，不影响预置关联
	_ = globalManager.Remove(role)
}

// TestAssociationManager_Set_WithTransaction_Rollback 验证在事务中通过 UseDB 注入的
// AssociationManager.Set 操作在回滚后不会改变原始关联状态。
func TestAssociationManager_Set_WithTransaction_Rollback(t *testing.T) {
	setupTestData(t)
	setupPreloadTestData(t)

	user := getTestUser(t)
	roleAdmin := getTestRole(t, "Admin")
	roleViewer := getTestRole(t, "Viewer")

	// 先通过全局 DB 建立初始关联（Admin）
	initialManager := gormx.NewAssociationManager[User, Role](user, "Roles")
	if err := initialManager.Clear(); err != nil {
		t.Fatalf("initial Clear failed: %v", err)
	}
	if err := initialManager.Add(roleAdmin); err != nil {
		t.Fatalf("initial Add failed: %v", err)
	}

	tx := gormx.GetDb().Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}

	// 在事务内将角色替换为 Viewer
	txManager := gormx.NewAssociationManager[User, Role](user, "Roles", gormx.UseDB(tx))
	if err := txManager.Set([]Role{roleViewer}); err != nil {
		tx.Rollback()
		t.Fatalf("Set inside tx failed: %v", err)
	}

	// 回滚
	if err := tx.Rollback().Error; err != nil {
		t.Fatalf("rollback failed: %v", err)
	}

	// 全局 DB 中应仍保留原来的 Admin 关联
	globalManager := gormx.NewAssociationManager[User, Role](user, "Roles")
	afterRoles, err := globalManager.All()
	if err != nil {
		t.Fatalf("All after rollback failed: %v", err)
	}

	adminFound := false
	for _, r := range afterRoles {
		if r.Name == "Viewer" {
			t.Error("expected role 'Viewer' to be absent after rollback, but it was found")
		}
		if r.Name == "Admin" {
			adminFound = true
		}
	}
	if !adminFound {
		t.Error("expected original role 'Admin' to still be present after rollback")
	}

	// 清理
	_ = globalManager.Clear()
}

// TestAssociationManager_WithContext_CancelledContext 验证通过 UseDB 注入携带已取消
// context 的 DB 时，AssociationManager 的操作会正确返回错误，而非使用全局 DB 继续执行。
func TestAssociationManager_WithContext_CancelledContext(t *testing.T) {
	setupTestData(t)
	setupPreloadTestData(t)

	user := getTestUser(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 提前取消

	ctxDB := gormx.GetDb().WithContext(ctx)
	manager := gormx.NewAssociationManager[User, Role](user, "Roles", gormx.UseDB(ctxDB))

	_, err := manager.All()
	if err == nil {
		t.Error("expected an error from the cancelled context, got nil")
	}
}
