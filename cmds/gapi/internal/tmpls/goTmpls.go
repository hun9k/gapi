package tmpls

var ResourceMessages = `package {{.resource}}
`

var ResourceHandlers = `package {{.resource}}

import (
	"{{.modPath}}/models"

	"github.com/gin-gonic/gin"
	"github.com/hun9k/gapi/base"
)

func get(ctx *gin.Context) {
	// base get list
	base.Get[{{.modelName}}](ctx)
}

func post(ctx *gin.Context) {
	base.Post[{{.modelName}}](ctx)
}

func put(ctx *gin.Context) {
	base.Put[{{.modelName}}](ctx)
}

func delete(ctx *gin.Context) {
	base.Delete[{{.modelName}}](ctx)
}

func restore(ctx *gin.Context) {
	base.Restore[{{.modelName}}](ctx)
}

`

var ResourceRouters = `package {{.resource}}

import (
	gin "github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.RouterGroup) {
	group := g.Group("{{.resource}}")
	group.GET("", get)            // 查
	group.POST("", post)          // 增
	group.PUT("", put)            // 改
	group.DELETE("", delete)      // 删
	group.PUT("restore", restore) // 恢
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
	"github.com/hun9k/gapi/dao"
	"github.com/hun9k/gapi/db"
)

func init() {
	// models
	models := []any{
		{{.modelList}},
	}

	// migrate db
	if conf.Get[string]("app.mode") == conf.APP_MODE_DEV {
		dao.ModelMigrate(db.Inst(), models...)
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
	// v1 := api.Router().Group("") // .Use(middleware.AuthMiddleware())
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
