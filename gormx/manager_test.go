package gormx_test

import (
	"fmt"
	"github.com/shijl0925/go-toolkits/gormx"
	"testing"
)

// TestGenericAssociationManager_Add 测试 Add 方法
func TestGenericAssociationManager_Add(t *testing.T) {
	// 准备测试数据
	setupTestData(t)

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
