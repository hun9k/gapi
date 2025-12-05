package log

import (
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/hun9k/gapi/conf"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Writter
var writers = map[string]io.Writer{}

const (
	WRITER_NAME_DFT = "default"
)

func WriterInstance(ns ...string) io.Writer {
	name := WRITER_NAME_DFT
	if len(ns) > 0 {
		name = ns[0]
	}
	if writers[name] == nil {
		writers[name] = newWriter(name)
	}

	return writers[name]
}

func newWriter(name string) io.Writer {
	// std
	if conf.Get[string](fmt.Sprintf("log.%s.writer", name)) == conf.LOG_WRITER_STDOUT {
		return os.Stdout
	}

	// parse writerUrl
	writerUrl := conf.Get[string](fmt.Sprintf("log.%s.writer", name))
	out, err := url.Parse(writerUrl)
	if err != nil {
		slog.Warn(fmt.Sprintf("log.%s.writer parse error", name), "error", err)
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
