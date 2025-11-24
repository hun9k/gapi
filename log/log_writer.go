package log

import (
	"io"
	"os"

	"github.com/hun9k/gapi/conf"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 日志Writter
var wirterInstance io.Writer

func writerSingle() io.Writer {
	if wirterInstance == nil {
		wirterInstance = newLogWriter()
	}

	return wirterInstance
}

func newLogWriter() io.Writer {
	var lw io.Writer
	switch conf.Log().Format {
	case "file":
		lw = &lumberjack.Logger{
			Filename:   conf.Log().File.Filename,   // 日志文件
			MaxSize:    conf.Log().File.MaxSize,    // 最大文件尺寸(MB)
			MaxBackups: conf.Log().File.MaxBackups, // 最多备份文件数
			MaxAge:     conf.Log().File.MaxAge,     // 最大保存天数
			Compress:   conf.Log().File.Compress,   // 是否压缩旧日志
		}
	case "std":
		fallthrough
	default:
		lw = os.Stdout
	}

	return lw
}
