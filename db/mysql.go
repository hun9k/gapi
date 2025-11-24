package db

import (
	"log"
	"time"

	"github.com/hun9k/gapi/conf"
	gapiLog "github.com/hun9k/gapi/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func mySQLNew() (*gorm.DB, error) {
	// logger 设置
	logLevel := logger.Info
	switch conf.App().Mode {
	case conf.AM_PROD:
		logLevel = logger.Warn
	// case CONF_APP_MODE_TEST, CONF_APP_MODE_DEV:
	// 	fallthrough
	default:
		logLevel = logger.Info
	}
	lger := logger.New(
		log.New(gapiLog.LogWriter(), "\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			// Colorful:                  false,
			// IgnoreRecordNotFoundError: false,
			// ParameterizedQueries:      false,
			LogLevel: logLevel,
		},
	)

	// connections
	db, err := gorm.Open(mysql.Open(conf.MySQL().DSN), &gorm.Config{
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
