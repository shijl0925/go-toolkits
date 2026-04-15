package gormx_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shijl0925/go-toolkits/gormx"
	"gorm.io/gorm"
)

// TestUseDB_ReplacesGlobalDB verifies the fundamental contract of UseDB:
// GetDb(UseDB(other)) must return `other`, not globalDb.
// We open a second connection (a dry-run one) so the test does not need a
// second real database server, while still proving that the pointer identity
// changes from globalDb to the injected DB.
func TestUseDB_ReplacesGlobalDB(t *testing.T) {
	requireTestDB(t) // skip when no live DB

	globalDB := gormx.GetDb()

	// Use a Session-derived DB as "another" DB — it is a distinct *gorm.DB.
	otherDB := globalDB.Session(&gorm.Session{NewDB: true})
	if otherDB == globalDB {
		t.Fatal("test prerequisite failed: Session() must return a new *gorm.DB instance")
	}

	got := gormx.GetDb(gormx.UseDB(otherDB))
	if got != otherDB {
		t.Errorf("GetDb(UseDB(otherDB)) returned wrong instance: expected the injected DB, got something else")
	}
	// Also confirm it is not globalDb.
	if got == globalDB {
		t.Errorf("GetDb(UseDB(otherDB)) still returned globalDb — UseDB is not substituting the DB")
	}
}

// TestUseDB_SubsequentOptsApplied verifies that DBOption arguments placed
// *after* UseDB are applied to the substituted DB, not discarded.
// We use a dry-run Session to capture the resulting *gorm.Statement without
// executing any SQL; the presence of the WHERE clause proves the extra opt ran.
func TestUseDB_SubsequentOptsApplied(t *testing.T) {
	requireTestDB(t)

	repo := &gormx.BaseRepo[User]{}

	// SelectCount with UseDB as first opt and a WHERE condition as second opt.
	// The WHERE condition is applied to the UseDB-injected DB, not globalDb.
	injected := gormx.GetDb().Session(&gorm.Session{NewDB: true})

	// Two opts: UseDB to substitute the base, then a plain Where to narrow it.
	count, err := repo.SelectCount(
		gormx.UseDB(injected),
		func(db *gorm.DB) *gorm.DB { return db.Where("1 = 0") }, // always false
	)
	if err != nil {
		t.Fatalf("SelectCount with UseDB + extra opt: %v", err)
	}
	// The WHERE 1=0 opt must have been applied — result must be 0.
	if count != 0 {
		t.Errorf("expected count 0 (WHERE 1=0), got %d; extra opt was not applied after UseDB", count)
	}
}

// TestUseDB_WithTransaction_Rollback is the same scenario as
// TestBaseRepo_Insert_WithTransaction_Rollback but expressed using the named
// gormx.UseDB helper instead of an inline function literal.
// This validates that the helper is a drop-in ergonomic replacement.
func TestUseDB_WithTransaction_Rollback(t *testing.T) {
	setupTestData(t)

	repo := &gormx.BaseRepo[User]{}

	tx := gormx.GetDb().Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}

	// UseDB(tx) is the named equivalent of: func(db *gorm.DB) *gorm.DB { return tx }
	newUser := &User{Name: "usedb_rollback_user", Email: "usedb_rb@example.com", Age: 22, Status: 1}
	if err := repo.Insert(newUser, gormx.UseDB(tx)); err != nil {
		tx.Rollback()
		t.Fatalf("Insert inside tx via UseDB failed: %v", err)
	}

	// Record must be visible within the same transaction via UseDB.
	inTxUser, err := repo.SelectOneByOpts(
		gormx.UseDB(tx),
		func(db *gorm.DB) *gorm.DB { return db.Where("name = ?", "usedb_rollback_user") },
	)
	if err != nil {
		tx.Rollback()
		t.Fatalf("record should be visible inside tx: %v", err)
	}
	if inTxUser.Name != "usedb_rollback_user" {
		tx.Rollback()
		t.Fatalf("expected usedb_rollback_user, got %q", inTxUser.Name)
	}

	// Roll back.
	if err := tx.Rollback().Error; err != nil {
		t.Fatalf("rollback failed: %v", err)
	}

	// After rollback the record must be gone from globalDb.
	q, u := gormx.NewQuery[User]()
	q.Eq(&u.Name, "usedb_rollback_user")
	_, err = repo.SelectOneByOpts(q.ToOptions()...)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("expected ErrRecordNotFound after rollback, got: %v", err)
	}
}

// TestUseDB_WithTransaction_Commit mirrors the commit scenario, using UseDB.
func TestUseDB_WithTransaction_Commit(t *testing.T) {
	setupTestData(t)

	repo := &gormx.BaseRepo[User]{}

	tx := gormx.GetDb().Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}

	newUser := &User{Name: "usedb_commit_user", Email: "usedb_cm@example.com", Age: 27, Status: 1}
	if err := repo.Insert(newUser, gormx.UseDB(tx)); err != nil {
		tx.Rollback()
		t.Fatalf("Insert inside tx via UseDB failed: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		t.Fatalf("commit failed: %v", err)
	}

	// After commit the record must be visible through globalDb.
	q, u := gormx.NewQuery[User]()
	q.Eq(&u.Name, "usedb_commit_user")
	result, err := repo.SelectOneByOpts(q.ToOptions()...)
	if err != nil {
		t.Errorf("expected committed record to be visible via globalDb, got: %v", err)
	}
	if result.Name != "usedb_commit_user" {
		t.Errorf("expected name usedb_commit_user, got %q", result.Name)
	}
}

// TestUseDB_WithContext_CancelledContext verifies that UseDB correctly injects a
// context-scoped DB built with gorm.WithContext into the repo call chain.
// Because the context is already cancelled when the query runs, the call must
// return an error — proving UseDB propagates the context-scoped DB all the way
// to the underlying GORM operation.
func TestUseDB_WithContext_CancelledContext(t *testing.T) {
	setupTestData(t)

	repo := &gormx.BaseRepo[User]{}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancelled before any DB activity

	// Build a context-scoped DB and inject it via UseDB.
	ctxDB := gormx.GetDb().WithContext(ctx)

	_, err := repo.SelectListByOpts(gormx.UseDB(ctxDB))
	if err == nil {
		t.Error("expected an error from the cancelled context when using UseDB, got nil")
	}
}
