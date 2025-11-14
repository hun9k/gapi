package gapi

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type GrpcService struct {
	*gin.Engine
}

var grpcService *HttpService

// 获取HTTP服务对象
func GrpcSvc() *HttpService {
	if httpService != nil {
		return httpService
	}

	httpService = newGrpcService()
	return httpService
}

// 新建HTTP服务对象
func newGrpcService() *HttpService {
	engine := gin.New()

	// 日志
	gin.DefaultWriter = logWriter
	engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Output: logWriter,
		// Formatter: func(params gin.LogFormatterParams) string { return "" },
	}))

	// 恢复
	engine.Use(gin.Recovery())

	return &HttpService{
		engine,
	}
}

func (service *GrpcService) Listen() {
	slog.Info("HTTP service is listening", "addr", viper.GetString("httpService.addr"))
	service.Run(viper.GetString("httpService.addr"))
}
