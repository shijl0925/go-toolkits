package gormx_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shijl0925/go-toolkits/gormx"
	"gorm.io/gorm"
)

// Compile-time proof that BaseRepo[T] satisfies the full IBaseRepo[T] interface,
// including the opts ...DBOption parameters that were added to fix the issue.
// This line will not compile if any method signature is missing or mismatched.
var _ gormx.IBaseRepo[User] = (*gormx.BaseRepo[User])(nil)

// TestBaseRepo_Insert_WithTransaction_Rollback verifies that Insert (and SelectOneByOpts)
// correctly use the DB passed via opts rather than globalDb:
//  1. A record inserted inside a transaction is visible within that transaction.
//  2. After the transaction is rolled back the record is no longer visible via the
//     main (global) DB — proving the operation actually ran inside the supplied tx.
func TestBaseRepo_Insert_WithTransaction_Rollback(t *testing.T) {
	setupTestData(t)

	repo := &gormx.BaseRepo[User]{}
	db := gormx.GetDb()

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}

	txOpt := func(db *gorm.DB) *gorm.DB { return tx }

	newUser := &User{Name: "tx_rollback_user", Email: "tx_rollback@example.com", Age: 25, Status: 1}
	if err := repo.Insert(newUser, txOpt); err != nil {
		tx.Rollback()
		t.Fatalf("Insert inside tx failed: %v", err)
	}

	// Record must be visible within the same transaction (proves txOpt was used for Insert).
	inTxUser, err := repo.SelectOneByOpts(
		txOpt,
		func(db *gorm.DB) *gorm.DB { return db.Where("name = ?", "tx_rollback_user") },
	)
	if err != nil {
		tx.Rollback()
		t.Fatalf("record should be visible inside tx: %v", err)
	}
	if inTxUser.Name != "tx_rollback_user" {
		tx.Rollback()
		t.Fatalf("expected tx_rollback_user, got %q", inTxUser.Name)
	}

	// Roll back the transaction.
	if err := tx.Rollback().Error; err != nil {
		t.Fatalf("rollback failed: %v", err)
	}

	// After rollback the record must not be visible through the global DB.
	q, u := gormx.NewQuery[User]()
	q.Eq(&u.Name, "tx_rollback_user")
	_, err = repo.SelectOneByOpts(q.ToOptions()...)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("expected ErrRecordNotFound after rollback, got: %v", err)
	}
}

// TestBaseRepo_Insert_WithTransaction_Commit verifies the symmetric case:
// a record inserted inside a committed transaction is durably visible through
// the global DB afterwards — proving the tx opt was propagated to the write path.
func TestBaseRepo_Insert_WithTransaction_Commit(t *testing.T) {
	setupTestData(t)

	repo := &gormx.BaseRepo[User]{}
	db := gormx.GetDb()

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}

	txOpt := func(db *gorm.DB) *gorm.DB { return tx }

	newUser := &User{Name: "tx_commit_user", Email: "tx_commit@example.com", Age: 30, Status: 1}
	if err := repo.Insert(newUser, txOpt); err != nil {
		tx.Rollback()
		t.Fatalf("Insert inside tx failed: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		t.Fatalf("commit failed: %v", err)
	}

	// After commit the record must be visible through the global DB.
	q, u := gormx.NewQuery[User]()
	q.Eq(&u.Name, "tx_commit_user")
	result, err := repo.SelectOneByOpts(q.ToOptions()...)
	if err != nil {
		t.Errorf("expected committed record to be visible, got: %v", err)
	}
	if result.Name != "tx_commit_user" {
		t.Errorf("expected name tx_commit_user, got %q", result.Name)
	}
}

// TestBaseRepo_MultipleOps_SameTransaction verifies that two separate BaseRepo calls
// (Insert + SelectCount) participate in the same transaction when the same txOpt is
// passed to both.  Rolling back the transaction makes all writes disappear atomically.
func TestBaseRepo_MultipleOps_SameTransaction(t *testing.T) {
	setupTestData(t)

	repo := &gormx.BaseRepo[User]{}
	db := gormx.GetDb()

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}

	txOpt := func(db *gorm.DB) *gorm.DB { return tx }

	u1 := &User{Name: "tx_atomic_user1", Email: "atomic1@example.com", Age: 20, Status: 1}
	u2 := &User{Name: "tx_atomic_user2", Email: "atomic2@example.com", Age: 21, Status: 1}

	if err := repo.Insert(u1, txOpt); err != nil {
		tx.Rollback()
		t.Fatalf("insert u1: %v", err)
	}
	if err := repo.Insert(u2, txOpt); err != nil {
		tx.Rollback()
		t.Fatalf("insert u2: %v", err)
	}

	// Both records must be visible inside the transaction (opts propagation for reads).
	inTxCount, err := repo.SelectCount(
		txOpt,
		func(db *gorm.DB) *gorm.DB {
			return db.Where("name IN ?", []string{"tx_atomic_user1", "tx_atomic_user2"})
		},
	)
	if err != nil {
		tx.Rollback()
		t.Fatalf("SelectCount inside tx: %v", err)
	}
	if inTxCount != 2 {
		tx.Rollback()
		t.Fatalf("expected 2 records inside tx, got %d", inTxCount)
	}

	// Roll back — both inserts must be undone atomically.
	if err := tx.Rollback().Error; err != nil {
		t.Fatalf("rollback failed: %v", err)
	}

	// Neither record should be visible through the global DB.
	outsideCount, err := repo.SelectCount(
		func(db *gorm.DB) *gorm.DB {
			return db.Where("name IN ?", []string{"tx_atomic_user1", "tx_atomic_user2"})
		},
	)
	if err != nil {
		t.Errorf("SelectCount outside tx: %v", err)
	}
	if outsideCount != 0 {
		t.Errorf("expected 0 records after rollback, got %d", outsideCount)
	}
}

// TestBaseRepo_UsesPassedContext verifies that a context-scoped DB passed via opts
// is actually forwarded to the underlying GORM call.  If globalDb were used instead,
// the already-cancelled context would have no effect and the call would succeed.
func TestBaseRepo_UsesPassedContext(t *testing.T) {
	setupTestData(t)

	repo := &gormx.BaseRepo[User]{}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before any DB call

	ctxOpt := func(db *gorm.DB) *gorm.DB {
		return db.WithContext(ctx)
	}

	_, err := repo.SelectListByOpts(ctxOpt)
	if err == nil {
		t.Error("expected an error propagated from the cancelled context, got nil")
	}
}
