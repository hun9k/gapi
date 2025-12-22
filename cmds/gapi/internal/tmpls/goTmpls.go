package tmpls

var ResourceMessages = `package {{.resource}}

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

var ResourceHandlers = `package {{.resource}}

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

var ResourceRouters = `package {{.resource}}

func crudRouters() {
	router.OPTIONS("", nil)       // OPTIONS
	router.OPTIONS(":id", nil)    // OPTIONS
	router.POST("", Create)       // 增
	router.DELETE(":id", Delete)  // 删单id
	router.DELETE("", DeleteMany) // 删多id
	router.PUT(":id", Update)     // 改单id
	router.PUT("", UpdateMany)    // 改多id
	router.GET(":id", GetOne)     // 查单id
	router.GET("", Get)           // 查多id,或过滤条件
}

`

var ResourceSetup = `package {{.resource}}

import "github.com/gin-gonic/gin"


// 自定义路由
func routers() {
	// router.GET("path", func(c *gin.Context) {})
}

// 资源中间件
var middlewares = []gin.HandlerFunc{}

// 资源路由
var router *gin.RouterGroup

// 设置资源路由
func SetupRouter(g *gin.RouterGroup) {
	// 资源路由组
	router = g.Group("{{.resource}}", middlewares...)
	// CRUD路由
	crudRouters()

	// 自定义路由
	routers()
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

var HandlersInit = `package handlers

`

var PlatformSetup = `package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/hun9k/gapi/base"
)

const adminRouterVerion = ""
const adminRouterPrefix = "{{.platform}}"

// 中间件列表
var adminMiddlewares = []gin.HandlerFunc{
	base.CorsDefault(),
}

`

var PlatformRouters = `package handlers
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
