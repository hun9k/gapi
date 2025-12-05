package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hun9k/gapi/conf"
	"github.com/hun9k/gapi/log"
)

// single instance mode
var handlers map[string]*gin.Engine

func Instance(ks ...string) *gin.Engine {
	key := "" // default key
	if len(ks) > 0 {
		key = ks[0]
	}
	if handlers == nil {
		handlers[key] = newHandler()
	}
	return handlers[key]
}

func newHandler() *gin.Engine {
	initGinHanler()

	return gin.Default()
}

func initGinHanler() {
	handlerSetLogWriter()
	handlerSetMode()
}

func handlerSetMode() {
	switch conf.Get[string]("app.mode") {
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
