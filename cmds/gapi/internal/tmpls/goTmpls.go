package tmpls

var RoutersPing = `package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hun9k/gapi/services/api"
)

func init() {
	api.Router().GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
`

const Main = `package main

import (
	// 初始化路由
	_ "{{.Path}}/routers"

	"github.com/hun9k/gapi/app"
)

func main() {
	// 应用运行
	app.Run()
}

`
