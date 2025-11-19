package tmpls

const Apis_group = `package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/hun9k/gapi"
)

func routerGroup(path string) *gin.RouterGroup {
	return gapi.RouterGroup(path)

	// ues middleware for group path
	// return gapi.RouterGroup(path, func(ctx *gin.Context) {}, func(ctx *gin.Context) {})
}

`
