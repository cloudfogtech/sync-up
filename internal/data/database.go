package data

import (
	"gorm.io/gorm"
)

type DB struct {
	orm *gorm.DB
}

func NewDB(gorm *gorm.DB) *DB {
	return &DB{
		orm: gorm,
	}
}

func (db *DB) Migrate() {
	err := db.orm.AutoMigrate(Log{})
	if err != nil {
		panic(err)
	}
}
