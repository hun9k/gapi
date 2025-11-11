package http

import (
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type HttpService struct {
	router *gin.Engine
}

// 全局存储
var httpService *HttpService

func NewHttpService() *HttpService {
	httpService = &HttpService{}
	httpService.
		modeInit().
		routerInit()

	return httpService
}

func (s *HttpService) modeInit() *HttpService {
	switch strings.ToLower(viper.GetString("mode")) {
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	case "dev":
		fallthrough
	default:
		gin.SetMode(gin.DebugMode)
	}
	return s
}

func (s *HttpService) Listen() {
	slog.Info("HTTP service is listening", "addr", viper.GetString("httpService.addr"))
	s.router.Run(viper.GetString("httpService.addr"))
}
