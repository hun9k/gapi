package http

import (
	"github.com/gin-gonic/gin"
	"github.com/hun9k/gapi/conf"
	"github.com/hun9k/gapi/log"
)

// 获取Handler
func Handler() *gin.Engine {
	return handlerSingle()
}

// single instance mode
var _handler *gin.Engine

func handlerSingle() *gin.Engine {
	if _handler == nil {
		_handler = handlerNew()
	}
	return _handler
}

func handlerNew() *gin.Engine {
	handlerSetLogWriter()
	handlerSetMode()
	return gin.Default()
}

func handlerSetMode() {
	switch conf.App().Mode {
	case conf.APP_MODE_TEST:
		gin.SetMode(gin.TestMode)
	case conf.APP_MODE_PROD:
		gin.SetMode(gin.ReleaseMode)
	case conf.APP_MODE_DEV:
		fallthrough
	default:
		gin.SetMode(gin.DebugMode)
	}
}

func handlerSetLogWriter() {
	// 日志Writer
	gin.DefaultWriter = log.LogWriter()
}
