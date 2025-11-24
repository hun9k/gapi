package http

import (
	"github.com/gin-gonic/gin"
)

type service struct {
	handler *handler
}

var _service *service

// 获取HTTP服务对象
func serviceSingle() *service {
	if _service == nil {
		_service = newService()
	}

	return _service
}

// 新建HTTP服务对象
func newService() *service {
	// 模式设置
	handlerSetMode()
	handlerSetLogWriter()

	// http 核心对象
	handler := newHandler()

	// 默认中间件
	// Log
	handler.Use(gin.Logger())
	// 恢复
	handler.Use(gin.Recovery())

	return &service{
		handler,
	}
}
