package log

import (
	"log/slog"

	"github.com/hun9k/gapi/conf"
)

// 日志组件
type log struct {
	*slog.Logger
}

var logInstance *log

func logSingle() *log {
	if logInstance == nil {
		logInstance = newLogger()
	}

	return logInstance
}

func newLogger() *log {
	// 声明选项
	options := &slog.HandlerOptions{}

	// 设置最低级别日志
	switch conf.Get[string]("app.mod") {
	case conf.APP_MODE_PROD:
		options.Level = slog.LevelInfo
	// case CONF_APP_MODE_DEV, CONF_APP_MODE_TEST:
	// 	fallthrough
	default:
		options.Level = slog.LevelDebug
	}

	// 基于类型设置handler
	var h slog.Handler
	switch conf.Get[string]("log.format") {
	case "json":
		h = slog.NewJSONHandler(writerSingle(), options)
	case "text":
		fallthrough
	default:
		h = slog.NewTextHandler(writerSingle(), options)
	}

	lg := slog.New(h)
	// 全局logger
	slog.SetDefault(lg)

	return &log{
		lg,
	}
}
