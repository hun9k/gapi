package tmpls

var Resource_routers = `package routers

import (
	"{{.Mod.Path}}/internal/{{.Resource}}"
)

func init() {
	// Rest
	group := routerGroup("{{.Version}}").Group("{{.Resource}}")
	group.POST("", {{.Resource}}.Post)                       // 创建一个资源
	group.DELETE(":id", {{.Resource}}.DeleteId)              // 删除一个资源
	group.DELETE("", {{.Resource}}.Delete)                   // 删除多个资源
	group.PATCH("restore/:id", {{.Resource}}.RestoreId)      // 恢复一个资源
	group.PATCH("restore", {{.Resource}}.Restore)            // 恢复多个资源
	group.PUT(":id", {{.Resource}}.PutId)                    // 更新一个资源
	group.PUT("", {{.Resource}}.Put)                         // 更新多个资源
	group.PATCH(":id", {{.Resource}}.PatchId)                // 更新一个资源的部分字段
	group.PATCH("", {{.Resource}}.Patch)                     // 更新多个资源的部分字段
	group.GET(":id", {{.Resource}}.GetId)                    // 获取一个资源
	group.GET("", {{.Resource}}.Get)                         // 获取多个资源
}

`

var Routers_resources_bare = `package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	// Rest
	group := routerGroup("{{.Version}}").Group("{{.Resource}}")
	group.POST("", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "{{.Resource}} post")
	}) // 创建一个资源
	group.DELETE(":id", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "{{.Resource}} delete by id")
	}) // 删除一个资源
	group.DELETE("", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "{{.Resource}} delete")
	}) // 删除多个资源
	group.PATCH(":id/restore", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "{{.Resource}} restore by id")
	}) // 恢复一个资源
	group.PATCH("restore", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "{{.Resource}} restore")
	}) // 恢复多个资源
	group.PUT(":id", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "{{.Resource}} put by id")
	}) // 更新一个资源
	group.PUT("", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "articles put")
	}) // 更新多个资源
	group.PATCH(":id", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "{{.Resource}} patch by id")
	}) // 更新一个资源的部分字段
	group.PATCH("", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "{{.Resource}} patch")
	}) // 更新多个资源的部分字段
	group.GET(":id", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "{{.Resource}} get by id")
	}) // 获取一个资源
	group.GET("", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "{{.Resource}} get")
	}) // 获取多个资源
}

`
