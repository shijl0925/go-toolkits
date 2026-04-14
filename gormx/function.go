package gormx

import "gorm.io/gorm"

type DBOption func(*gorm.DB) *gorm.DB

func Where(query interface{}, args ...interface{}) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where(query, args...)
	}
}

func GetDb(opts ...DBOption) *gorm.DB {
	db := globalDb
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

// isValidFieldName validates that a field name contains only safe characters
// (letters, digits, underscores, and at most one dot for table-qualified names)
// to prevent SQL injection through field name parameters.
func isValidFieldName(field string) bool {
	// The empty-string check must come first; the subsequent index expression
	// is safe because short-circuit evaluation skips it when field == "".
	if field == "" || field[len(field)-1] == '.' {
		return false
	}
	dotCount := 0
	for i, r := range field {
		if r == '.' {
			dotCount++
			if dotCount > 1 || i == 0 {
				return false
			}
			continue
		}
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_') {
			return false
		}
	}
	return true
}
