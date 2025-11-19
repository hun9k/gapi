package gapi

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

// 日志Writter
var _logWtr io.Writer

func logWtr() io.Writer {
	if _logWtr == nil {
		_logWtr = newLogWtr()
	}

	return _logWtr
}

func newLogWtr() io.Writer {
	var lw io.Writer
	switch Conf().Log.Output {
	case "file":
		lw = &lumberjack.Logger{
			Filename:   Conf().Log.File.Filename,   // 日志文件
			MaxSize:    Conf().Log.File.MaxSize,    // 最大文件尺寸(MB)
			MaxBackups: Conf().Log.File.MaxBackups, // 最多备份文件数
			MaxAge:     Conf().Log.File.MaxAge,     // 最大保存天数
			Compress:   Conf().Log.File.Compress,   // 是否压缩旧日志
		}
	case "std":
		fallthrough
	default:
		lw = os.Stdout
	}

	return lw
}

// 日志组件

var _log *slog.Logger

func Log() *slog.Logger {
	if _log != nil {
		return _log
	}

	_log = newLogger()

	return _log
}

func newLogger() *slog.Logger {
	// 声明选项
	options := &slog.HandlerOptions{}

	// 设置最低级别日志
	switch Conf().App.Mode {
	case APP_MODE_PROD:
		options.Level = slog.LevelInfo
	// case CONF_APP_MODE_DEV, CONF_APP_MODE_TEST:
	// 	fallthrough
	default:
		options.Level = slog.LevelDebug
	}

	// 基于类型设置handler
	var h slog.Handler
	switch Conf().Log.Format {
	case "json":
		h = slog.NewJSONHandler(logWtr(), options)
	case "text":
		fallthrough
	default:
		h = slog.NewTextHandler(logWtr(), options)
	}

	lg := slog.New(h)
	// 全局logger
	slog.SetDefault(lg)

	return lg
}
