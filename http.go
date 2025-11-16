package gapi

import (
	"github.com/gin-gonic/gin"
)

var _httpSvc *gin.Engine

// 获取HTTP服务对象
func HttpSvc() *gin.Engine {
	if _httpSvc != nil {
		return _httpSvc
	}

	_httpSvc = newHttpService()
	return _httpSvc
}

// HttpSvc 's alias
var HttpRouter, Router = HttpSvc, HttpSvc

// 新建HTTP服务对象
func newHttpService() *gin.Engine {
	// 模式
	switch Conf().App.Mode {
	case CONF_APP_MODE_TEST:
		gin.SetMode(gin.TestMode)
	case CONF_APP_MODE_PROD:
		gin.SetMode(gin.ReleaseMode)
	case CONF_APP_MODE_DEV:
		fallthrough
	default:
		gin.SetMode(gin.DebugMode)
	}

	// http 核心对象
	engine := gin.New()

	// 日志
	gin.DefaultWriter = logWtr()
	engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		// Output: logWriter,
		// Formatter: func(params gin.LogFormatterParams) string { return "" },
	}))

	// 恢复
	engine.Use(gin.Recovery())

	return engine
}
