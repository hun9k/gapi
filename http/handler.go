package http

import (
	"github.com/gin-gonic/gin"
	"github.com/hun9k/gapi/conf"
	"github.com/hun9k/gapi/log"
)

type RouterGroup struct {
	*gin.RouterGroup
}

type handler struct {
	*gin.Engine
}

func newHandler(opts ...gin.OptionFunc) *handler {
	return &handler{
		gin.New(opts...),
	}
}

// func defaultHandler() *handler {
// 	return &handler{
// 		gin.Default(),
// 	}
// }

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
