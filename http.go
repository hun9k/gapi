package gapi

import (
	"github.com/gin-gonic/gin"
)

var _httpSvc *gin.Engine

// 获取HTTP服务对象
func HttpSvc() *gin.Engine {
	if _httpSvc == nil {
		_httpSvc = newHttpService()
	}

	return _httpSvc
}

// HttpSvc 's alias
var HttpRouter, Router = HttpSvc, HttpSvc

var _groups = map[string]*gin.RouterGroup{}

// API router groups, single instance
func RouterGroup(path string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	_, exists := _groups[path]
	if !exists {
		_groups[path] = Router().Group(path, handlers...)
	}

	return _groups[path]
}

// 新建HTTP服务对象
func newHttpService() *gin.Engine {
	// 模式设置
	switch Conf().App.Mode {
	case APP_MODE_TEST:
		gin.SetMode(gin.TestMode)
	case APP_MODE_PROD:
		gin.SetMode(gin.ReleaseMode)
	case APP_MODE_DEV:
		fallthrough
	default:
		gin.SetMode(gin.DebugMode)
	}

	// http 核心对象
	engine := gin.New()

	// 默认中间件
	// 日志
	gin.DefaultWriter = logWtr()
	engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		// Output: logWriter,
		// Formatter: func(params gin.LogFormatterParams) string { return "" },
	}))

	// 恢复
	engine.Use(gin.Recovery())

	// dev 模式

	return engine
}

// func migrateHandler(engine *gin.Engine) {
// 	if Conf().App.Mode == APP_MODE_DEV {
// 		schemas := []any{}
// 		engine.GET("dev/migrate", func(ctx *gin.Context) {
// 			if err := DB().AutoMigrate(schemas...); err != nil {
// 				ctx.JSON(http.StatusOK, Resp{
// 					Status:  1,
// 					Message: err.Error(),
// 				})
// 				return
// 			}
// 			ctx.JSON(http.StatusOK, Resp{
// 				Status:  0,
// 				Message: "migrate success",
// 			})
// 		})
// 	}

// }
