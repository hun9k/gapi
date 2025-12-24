package tmpls

const ResMessages = `package {{.resource}}

import (
	"{{.modPath}}/models"
	{{if .iTime}}"time"{{end}}
)

type putBody struct {
{{range .fields}}
	{{.Name}} {{if .IsNonRef}}*{{end}}{{.Type}} {{.Tag}}{{end}}
}

func bodyToModel(body putBody) (model {{.modelName}}, cols []string) {
{{range .fields}}if body.{{.Name}} != nil {
		cols = append(cols, "{{.Col}}")
		model.{{.Name}} = {{if .IsNonRef}}*{{end}}body.{{.Name}}
	}
{{end}}
	return model, cols
}

`

const ResHandlers = `package {{.resource}}

import (
	"{{.modPath}}/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hun9k/gapi/base"
	"github.com/hun9k/gapi/log"
)

func Create(ctx *gin.Context) {
	base.Create[{{.modelName}}](ctx)
}

func Delete(ctx *gin.Context) {
	base.Delete[{{.modelName}}](ctx)
}

func DeleteMany(ctx *gin.Context) {
	base.DeleteMany[{{.modelName}}](ctx)
}

func Update(ctx *gin.Context) {
	// bind body
	body := putBody{}
	if err := ctx.ShouldBind(&body); err != nil {
		log.Info("bind body error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	model, cols := bodyToModel(body)
	base.Update(ctx, model, cols)
}

func UpdateMany(ctx *gin.Context) {
	// bind body
	body := putBody{}
	if err := ctx.ShouldBind(&body); err != nil {
		log.Info("bind body error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	model, cols := bodyToModel(body)
	base.UpdateMany(ctx, model, cols)
}

func GetOne(ctx *gin.Context) {
	base.GetOne[{{.modelName}}](ctx)
}

func Get(ctx *gin.Context) {
	base.Get[{{.modelName}}](ctx)
}

`

const ResRouters = `package {{.resource}}

import "github.com/gin-gonic/gin"

func crudRouters(group *gin.RouterGroup) {
	group.OPTIONS("", nil)       // OPTIONS
	group.OPTIONS(":id", nil)    // OPTIONS
	group.POST("", Create)       // 增
	group.DELETE(":id", Delete)  // 删单id
	group.DELETE("", DeleteMany) // 删多id
	group.PUT(":id", Update)     // 改单id
	group.PUT("", UpdateMany)    // 改多id
	group.GET(":id", GetOne)     // 查单id
	group.GET("", Get)           // 查多id,或过滤条件
}
`

const ResSetup = `package {{.resource}}

import "github.com/gin-gonic/gin"

// 资源中间件
var middlewares = []gin.HandlerFunc{}

// 设置资源路由
func SetupRouter(group *gin.RouterGroup) {
	g := group.Group("{{.routerPrefix}}", middlewares...)

	// CRUD路由
	crudRouters(g)
	
	// 自定义路由
	{{if ne .resource "user"}}// {{end}}routers(g)
}
`

const ResModel = `package {{.package}}

import base "github.com/hun9k/gapi/base"

type {{.model}} struct {
	// your fields here

	base.Model
}

// TableName 指定表名
// func (p *{{.model}}) TableName() string {
//     return "{{.resource}}"
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

const ModelsInit = `package models

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

const PlatSetup = `package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/hun9k/gapi/base"
)

const {{.platform}}RouterVerion = ""
const {{.platform}}RouterPrefix = "{{.platform}}"

{{if .user}}
// 非认证路由
var noCheckAuthPaths = map[string]struct{}{
	"{{if .platform}}/{{.platform}}{{end}}/{{.user}}/login": {},
}
{{end}}

// 中间件列表
var {{.platform}}Middlewares = []gin.HandlerFunc{
	base.CorsDefault(),
{{if .user}}	base.AuthDefault(noCheckAuthPaths),{{end}}
}

`

const PlatRouters = `package handlers
{{$platform := .platform}}{{$modPath := .modPath}}
import (
	{{range .resources}}"{{$modPath}}/handlers{{if $platform}}/{{$platform}}{{end}}/{{.}}"
{{end}}
	"github.com/hun9k/gapi/services/api"
)

func init() {
	// platform group
	platform := api.Router().Group({{.platform}}RouterVerion).Group({{.platform}}RouterPrefix)
	platform.Use({{.platform}}Middlewares...)

	// setup resource routers
	{{range .resources}}{{.}}.SetupRouter(platform)
{{end}}
}
`

const TasksInit = `package tasks

`

const HandlersInit = `package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hun9k/gapi/services/api"
)

func init() {
	api.Router().GET("ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
}

`

const Main = `package main

import (
	// API, Task init
	_ "{{.modPath}}/handlers"

	"github.com/hun9k/gapi/app"
)

func main() {
	// 应用运行
	app.Run()
}

`
