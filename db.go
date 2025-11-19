package gapi

import "gorm.io/gorm"

var _db *gorm.DB

func DB(k ...string) *gorm.DB {
	if _db == nil {
		db, err := newDB()
		if err != nil {
			Log().Warn("DB error", "error", err)
		}
		_db = db
	}

	return _db
}

func newDB() (*gorm.DB, error) {
	return newMySQL()
}
