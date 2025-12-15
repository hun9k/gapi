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
		db, err := newDB(key)
		if err != nil {
			log.Error("DB error", "error", err)
		}
		dbs[key] = db
	}

	return dbs[key]
}

func newDB(key string) (*gorm.DB, error) {
	switch strings.ToLower(conf.Get[string](fmt.Sprintf("db.%s.driver", key))) {
	case "mysql":
		return newMySQL(key)
	}
	return &gorm.DB{}, nil
}
