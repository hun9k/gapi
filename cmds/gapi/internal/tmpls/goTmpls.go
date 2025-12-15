package tmpls

var ResourceMessages = `package {{.resource}}
`

var ResourceHandlers = `package {{.resource}}

import "github.com/gin-gonic/gin"

func get(ctx *gin.Context) {

}

func getList(ctx *gin.Context) {

}

func create(ctx *gin.Context) {

}

func update(ctx *gin.Context) {

}
func updateList(ctx *gin.Context) {

}

func delete(ctx *gin.Context) {

}
func deleteList(ctx *gin.Context) {

}

func restore(ctx *gin.Context) {

}

func restoreList(ctx *gin.Context) {

}

`

var ResourceRouters = `package {{.resource}}

import gin "github.com/gin-gonic/gin"

func SetupRouter(g *gin.RouterGroup) {
	group := g.Group("{{.resource}}")
	group.GET(":id", get)
	group.GET("", getList)
	group.POST("", create)
	group.PUT(":id", update)
	group.PUT("", updateList)
	group.DELETE(":id", delete)
	group.DELETE("", deleteList)
	group.PUT(":id/restore", restore)
	group.POST("restore", restoreList)
}

`

var ResourceModel = `package {{.package}}

import base "github.com/hun9k/gapi/base"

type {{.model}} struct {
	// your fields here

	base.Model
}

// TableName 指定表名
// func (p *{{.model}}) TableName() string {
//     return "{{.resource}}"
// }

// Create 相关钩子函数：
// BeforeCreate 创建前的钩子
// func (p *{{.model}}) BeforeCreate(tx *gorm.DB) error {
//     return nil
// }
//
// AfterCreate 创建后的钩子
// func (p *{{.model}}) AfterCreate(tx *gorm.DB) error {
//     return nil
// }

// Update 相关钩子函数：
// BeforeUpdate 更新前的钩子
// func (p *{{.model}}) BeforeUpdate(tx *gorm.DB) error {
//     return nil
// }
//
// AfterUpdate 更新后的钩子
// func (p *{{.model}}) AfterUpdate(tx *gorm.DB) error {
//     return nil
// }

// Save 相关钩子函数：
// BeforeSave 保存前的钩子
// func (p *{{.model}}) BeforeSave(tx *gorm.DB) error {
//     return nil
// }
//
// AfterSave 保存后的钩子
// func (p *{{.model}}) AfterSave(tx *gorm.DB) error {
//     return nil
// }

// Delete 相关钩子函数：
// BeforeDelete 删除前的钩子
// func (p *{{.model}}) BeforeDelete(tx *gorm.DB) error {
//     return nil
// }
//
// AfterDelete 删除后的钩子
// func (p *{{.model}}) AfterDelete(tx *gorm.DB) error {
//     return nil
// }

// Find 相关钩子函数：
// AfterFind 查询后的钩子
// func (p *{{.model}}) AfterFind(tx *gorm.DB) error {
//     return nil
// }

`

var ModelsInit = `package models

import (
	"github.com/hun9k/gapi/conf"
	"github.com/hun9k/gapi/db"
)

func init() {
	// migrate db
	migrate()
}

func migrate() {
	if conf.Get[string]("app.mode") == conf.APP_MODE_DEV {
		db.Inst().AutoMigrate({{.modelList}})
	}
}
`

var Routers = `package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hun9k/gapi/services/api"
)

func init() {
	// routers
	routers()
}

func routers() {
	// // version group
	// v1 := api.Router().Group("v1") // .Use(middleware.AuthMiddleware())
	// {
	// 	// platform group
	// 	admin := v1.Group("admin") // .Use(middleware.AuthMiddleware())
	// 	{
	// 		// setup resource routers
	// 		contents.SetupRouter(admin)
	// 	}
	// }

	// ping
	api.Router().GET("ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}

`

const Main = `package main

import (
	// 初始化路由
	_ "{{.modPath}}/handlers"

	"github.com/hun9k/gapi/app"
)

func main() {
	// 应用运行
	app.Run()
}

`
