package http

import (
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	return httpService.router
}

func (s *HttpService) routerInit() *HttpService {
	s.router = gin.Default()
	return s
}
