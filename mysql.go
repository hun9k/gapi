package gapi

import (
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var _mysql *gorm.DB

func MySQL() *gorm.DB {
	if _mysql == nil {
		o, err := newMySQL()
		if err != nil {
			Log().Warn("mysql failed", "error", err.Error())
		}
		_mysql = o
	}
	return _mysql
}

func newMySQL() (*gorm.DB, error) {
	// logger 设置
	logLevel := logger.Info
	switch Conf().App.Mode {
	case CONF_APP_MODE_PROD:
		logLevel = logger.Warn
	// case CONF_APP_MODE_TEST, CONF_APP_MODE_DEV:
	// 	fallthrough
	default:
		logLevel = logger.Info
	}
	lger := logger.New(
		log.New(logWtr(), "\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			// Colorful:                  false,
			// IgnoreRecordNotFoundError: false,
			// ParameterizedQueries:      false,
			LogLevel: logLevel,
		},
	)

	// connections
	db, err := gorm.Open(mysql.Open(Conf().MySQL.DSN), &gorm.Config{
		// SkipDefaultTransaction:    false,
		// DefaultTransactionTimeout: 0,
		// DefaultContextTimeout:     0,
		// NamingStrategy:            nil,
		// FullSaveAssociations:      false,
		Logger: lger,
		// NowFunc: func() time.Time {
		// 	panic("TODO")
		// },
		// DryRun:                                   false,
		// PrepareStmt:                              false,
		// PrepareStmtMaxSize:                       0,
		// PrepareStmtTTL:                           0,
		// DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: true,
		// IgnoreRelationshipsWhenMigrating:         false,
		// DisableNestedTransaction:                 false,
		// AllowGlobalUpdate:                        false,
		// QueryFields:                              false,
		// CreateBatchSize:                          0,
		// TranslateError:                           false,
		// PropagateUnscoped:                        false,
		// ClauseBuilders:                           map[string]clause.ClauseBuilder{},
		// ConnPool:                                 nil,
		// Dialector:                                nil,
		// Plugins:                                  map[string]gorm.Plugin{},
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
