package db

import (
	"github.com/hun9k/gapi/log"
	"gorm.io/gorm"
)

var _db *gorm.DB

func DB(k ...string) *gorm.DB {
	if _db == nil {
		db, err := dbNew()
		if err != nil {
			log.Error("DB error", "error", err)
		}
		_db = db
	}

	return _db
}

func dbNew() (*gorm.DB, error) {
	return mySQLNew()
}
