package gapi

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

// 日志Writter
var logWriter io.Writer

func logWtr() io.Writer {
	if logWriter == nil {
		logWriter = newLogWriter()
	}

	return logWriter
}

func newLogWriter() io.Writer {
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
type Logger struct {
	*slog.Logger
}

var logger *Logger

func Log() *Logger {
	if logger == nil {
		logger = newLogger()
	}

	return logger
}

func newLogger() *Logger {
	// 声明选项
	options := &slog.HandlerOptions{}

	// 设置最低级别日志
	switch Conf().App.Mode {
	case "prod":
		options.Level = slog.LevelInfo
	case "dev":
		fallthrough
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

	return &Logger{
		lg,
	}
}
