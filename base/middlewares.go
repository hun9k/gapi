package base

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsDefault() gin.HandlerFunc {
	return cors.Default() // 允许所有来源
}
