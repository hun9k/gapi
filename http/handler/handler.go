package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hun9k/gapi/conf"
	"github.com/hun9k/gapi/log"
)

// single instance mode
var handlers = map[string]*gin.Engine{}

func Inst(ns ...string) *gin.Engine {
	name := "" // default key
	if len(ns) > 0 {
		name = ns[0]
	}
	if handlers == nil {
		handlers[name] = newHandler()
	}
	return handlers[name]
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
	gin.DefaultWriter = log.WriterInstance()
}
