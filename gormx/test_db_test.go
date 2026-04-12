package gormx_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/shijl0925/go-toolkits/gormx"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	testDialectMySQL    = "mysql"
	testDialectPostgres = "postgres"
)

var liveTestDBAvailable bool

func init() {
	dialect := defaultTestDialect()

	db, err := openTestDB(dialect)
	if err == nil {
		err = pingTestDB(db)
	}
	if err != nil {
		fallbackDB, fallbackErr := openFallbackTestDB(dialect)
		if fallbackErr == nil {
			gormx.Init(fallbackDB)
		}
		return
	}

	if err := migrateTestDB(db); err != nil {
		panic(fmt.Sprintf("failed to migrate %s database: %v", dialect, err))
	}

	liveTestDBAvailable = true
	gormx.Init(db)
}

func withTestDatabase(t *testing.T, dialect string, fn func(t *testing.T)) {
	t.Helper()

	previous := gormx.GetDb()
	db, err := openTestDB(dialect)
	if err == nil {
		err = pingTestDB(db)
	}
	if err != nil {
		t.Skipf("skipping %s database tests: %v", dialect, err)
	}

	if err := migrateTestDB(db); err != nil {
		t.Fatalf("failed to migrate %s database: %v", dialect, err)
	}

	gormx.Init(db)
	t.Cleanup(func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
		if previous != nil {
			gormx.Init(previous)
		}
	})

	fn(t)
}

func requireTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	if !liveTestDBAvailable {
		t.Skip("skipping database-backed gormx tests: test database is unavailable")
	}

	return gormx.GetDb()
}

func migrateTestDB(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Post{}, &Role{}, &UserRole{}, &Profile{})
}

func pingTestDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

func openTestDB(dialect string) (*gorm.DB, error) {
	switch dialect {
	case testDialectMySQL:
		return gorm.Open(mysql.Open(mysqlTestDSN()), &gorm.Config{})
	case testDialectPostgres:
		return gorm.Open(postgres.Open(postgresTestDSN()), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported test dialect: %s", dialect)
	}
}

func openFallbackTestDB(dialect string) (*gorm.DB, error) {
	cfg := &gorm.Config{
		DryRun:               true,
		DisableAutomaticPing: true,
	}

	switch dialect {
	case testDialectMySQL:
		return gorm.Open(mysql.New(mysql.Config{
			DSN:                       mysqlTestDSN(),
			SkipInitializeWithVersion: true,
		}), cfg)
	case testDialectPostgres:
		return gorm.Open(postgres.New(postgres.Config{
			DSN: postgresTestDSN(),
		}), cfg)
	default:
		return nil, fmt.Errorf("unsupported test dialect: %s", dialect)
	}
}

func mysqlTestDSN() string {
	if dsn := os.Getenv("GORMX_TEST_MYSQL_DSN"); dsn != "" {
		return dsn
	}

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		getenvDefault("GORMX_TEST_MYSQL_USER", "root"),
		getenvDefault("GORMX_TEST_MYSQL_PASSWORD", "root@123"),
		getenvDefault("GORMX_TEST_MYSQL_HOST", "127.0.0.1"),
		getenvDefault("GORMX_TEST_MYSQL_PORT", "3306"),
		getenvDefault("GORMX_TEST_MYSQL_DB", "vben"),
	)
}

func postgresTestDSN() string {
	if dsn := os.Getenv("GORMX_TEST_POSTGRES_DSN"); dsn != "" {
		return dsn
	}

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		getenvDefault("GORMX_TEST_POSTGRES_HOST", "127.0.0.1"),
		getenvDefault("GORMX_TEST_POSTGRES_USER", "postgres"),
		getenvDefault("GORMX_TEST_POSTGRES_PASSWORD", "postgres"),
		getenvDefault("GORMX_TEST_POSTGRES_DB", "vben"),
		getenvDefault("GORMX_TEST_POSTGRES_PORT", "5432"),
		getenvDefault("GORMX_TEST_POSTGRES_SSLMODE", "disable"),
		getenvDefault("GORMX_TEST_POSTGRES_TIMEZONE", "UTC"),
	)
}

func getenvDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func defaultTestDialect() string {
	return getenvDefault("GORMX_TEST_DEFAULT_DIALECT", testDialectMySQL)
}
