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
