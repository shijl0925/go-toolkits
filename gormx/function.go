package gormx

import "gorm.io/gorm"

type DBOption func(*gorm.DB) *gorm.DB

func Where(query interface{}, args ...interface{}) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where(query, args...)
	}
}

// UseDB returns a DBOption that replaces the base *gorm.DB with the provided
// instance, regardless of what globalDb is set to.  Use it to inject a
// request-scoped DB or an active transaction into any BaseRepo / Query call:
//
//	// request-scoped DB carrying the gin context
//	repo.SelectListByOpts(gormx.UseDB(orm.GetDB(c)))
//
//	// active transaction
//	tx := gormx.GetDb().Begin()
//	repo.Insert(&item, gormx.UseDB(tx))
//	repo.UpdateById(id, vars, gormx.UseDB(tx))
//	tx.Commit()
//
// Any DBOption arguments that follow UseDB in the same call are applied to
// the substituted DB in order, just as they would be with globalDb.
func UseDB(db *gorm.DB) DBOption {
	return func(_ *gorm.DB) *gorm.DB {
		return db
	}
}

func GetDb(opts ...DBOption) *gorm.DB {
	db := globalDb
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}
