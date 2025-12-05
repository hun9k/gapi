package log

import (
	"io"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/hun9k/gapi/conf"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 日志Writter
var writer io.Writer

func writerSingle() io.Writer {
	if writer == nil {
		writer = newWriter()
	}

	return writer
}

func newWriter() io.Writer {
	// std
	if conf.Get[string]("log.output") == "std" {
		return os.Stdout
	}

	// parse output
	output := conf.Get[string]("log.output")
	out, err := url.Parse(output)
	if err != nil {
		slog.Warn("log.output parse error", "error", err)
		return os.Stdout
	}

	// declare log writer
	var lw io.Writer

	switch strings.ToLower(out.Scheme) {
	case "file":
		maxSize, _ := strconv.Atoi(out.Query().Get("maxSize"))
		maxBackups, _ := strconv.Atoi(out.Query().Get("maxBackups"))
		maxAge, _ := strconv.Atoi(out.Query().Get("maxAge"))
		compress, _ := strconv.ParseBool(out.Query().Get("compress"))
		lw = &lumberjack.Logger{
			Filename:   out.Query().Get("filename"), // 日志文件
			MaxSize:    maxSize,                     // 最大文件尺寸(MB)
			MaxBackups: maxBackups,                  // 最多备份文件数
			MaxAge:     maxAge,                      // 最大保存天数
			Compress:   compress,                    // 是否压缩旧日志
		}

	default:
		lw = os.Stdout
	}

	return lw
}
