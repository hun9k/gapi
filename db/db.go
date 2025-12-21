package db

import (
	"fmt"
	"strings"

	"github.com/hun9k/gapi/conf"
	"github.com/hun9k/gapi/log"
	"gorm.io/gorm"
)

var dbs = map[string]*gorm.DB{}

const (
	DEFAULT_KEY = "default"
)

func Inst(keys ...string) *gorm.DB {
	key := DEFAULT_KEY
	if len(keys) > 0 {
		key = keys[0]
	}

	if dbs[key] == nil {
		db, err := NewDB(key)
		if err != nil {
			log.Error("DB error", "error", err)
			return nil
		}
		dbs[key] = db
	}

	return dbs[key]
}

func NewDB(key string) (db *gorm.DB, err error) {
	switch strings.ToLower(conf.Get[string](fmt.Sprintf("db.%s.driver", key))) {
	case "mysql":
		db, err = NewMySQL(key)
	}

	if err != nil {
		panic(err)
	}

	return db, err
}
