package gormx

import "gorm.io/gorm"

var globalDb *gorm.DB

func Init(db *gorm.DB) {
	globalDb = db
}
