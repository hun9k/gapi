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
	DEFAULT_NAME = "default"
)

func Inst(ns ...string) *gorm.DB {
	name := DEFAULT_NAME
	if len(ns) > 0 {
		name = ns[0]
	}
	if dbs[name] == nil {
		db, err := newDB(name)
		if err != nil {
			log.Error("DB error", "error", err)
		}
		dbs[name] = db
	}

	return dbs[name]
}

func newDB(name string) (*gorm.DB, error) {
	switch strings.ToLower(conf.Get[string](fmt.Sprintf("db.%s.driver", name))) {
	case "mysql":
		return newMySQL(name)
	}
	return nil, nil
}
