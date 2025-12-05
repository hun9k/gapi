package log

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/hun9k/gapi/conf"
)

// log.Instance("default").Debug() 等的快捷调用
func Debug(msg string, args ...any) {
	Inst().Debug(msg, args...)
}

func Info(msg string, args ...any) {
	Inst().Info(msg, args...)
}

func Warn(msg string, args ...any) {
	Inst().Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Inst().Error(msg, args...)
}

type logger = slog.Logger

var loggers = map[string]*logger{}

const LOGGER_NAME_DFT = "default"

func Inst(ns ...string) *logger {
	name := LOGGER_NAME_DFT
	if len(ns) > 0 {
		name = ns[0]
	}

	if loggers[name] == nil {
		loggers[name] = newLogger(name)
	}

	return loggers[name]
}

func newLogger(name string) *logger {
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
	switch strings.ToLower(conf.Get[string](fmt.Sprintf("log.%s.format", name))) {
	case conf.LOG_FORMAT_JSON:
		h = slog.NewJSONHandler(WriterInstance(name), options)
	case conf.LOG_FORMAT_TEXT:
		fallthrough
	default:
		h = slog.NewTextHandler(WriterInstance(name), options)
	}

	return (*logger)(slog.New(h))
}
