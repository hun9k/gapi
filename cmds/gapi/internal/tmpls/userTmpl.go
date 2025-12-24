package tmpls

const UserMessage = `package {{.resource}}

type loginBody struct {
	Username string ` + "`json:\"username\" binding:\"required\"`" + `
	Password string ` + "`json:\"password\" binding:\"required\"`" + `
}
`
const UserHandler = `package {{.resource}}

import (
	"{{.modPath}}/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hun9k/gapi/db"
	"github.com/hun9k/gapi/log"
	"github.com/hun9k/gapi/utils"
	"gorm.io/gorm"
)

func Login(ctx *gin.Context) {
	// bind body
	body := loginBody{}
	if err := ctx.ShouldBind(&body); err != nil {
		log.Info("get bind body error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   1,
			"message": err.Error(),
		})
		return
	}

	// check user
	user, err := checkUser(body.Username, body.Password)
	if err != nil {
		log.Info("check user error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusUnauthorized, nil)
		return
	}

	// TODO: generate token
	token, err := utils.MkJWT(user.ID)
	if err != nil {
		log.Info("gen token error", "path", ctx.Request.URL.Path, "error", err)
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	// response
	ctx.JSON(http.StatusOK, gin.H{
		"error": 0,
		"token": token, // TODO
	})
}

func checkUser(username, password string) (*{{.modelName}}, error) {
	user, err := gorm.G[{{.modelName}}](db.Inst()).Where("username = ?", username).First(context.Background())
	if err != nil {
		return nil, err
	}

	if bingo, err := utils.VerifyPassword(password, user.Password); err != nil || !bingo {
		return nil, err
	}

	return &user, nil
}
`

const UserRouter = `package {{.resource}}

import "github.com/gin-gonic/gin"

// 自定义路由
func routers(group *gin.RouterGroup) {
	group.Match([]string{"POST", "OPTIONS"}, "login", Login)
}
`

const UserModel = `package {{.package}}

import (
	base "github.com/hun9k/gapi/base"
	utils "github.com/hun9k/gapi/utils"
	"gorm.io/gorm"
)

type User struct {
	// your fields here
	Username string ` + "`gorm:\"size:1024;unique\" json:\"username\"`" + `
	Password string ` + "`gorm:\"size:4096;\" json:\"-\"`" + `

	Fullname string ` + "`gorm:\"size:1024;\" json:\"fullname\"`" + `
	Avater   string ` + "`gorm:\"\" json:\"avater\"`" + `

	base.Model
}

// TableName 指定表名
// func (m *User) TableName() string {
//     return "User"
// }

// Save 相关钩子函数：
// BeforeSave 保存前的钩子
func (m *User) BeforeSave(tx *gorm.DB) error {
	p, err := utils.EncryptPassword(m.Password)
	if err != nil {
		return err
	}
	m.Password = p

	return nil
}
// 
// AfterSave 保存后的钩子
// func (m *User) AfterSave(tx *gorm.DB) error {
//     return nil
// }

// Create 相关钩子函数：
// BeforeCreate 创建前的钩子
// func (m *User) BeforeCreate(tx *gorm.DB) error {
// 	return nil
// }
//
// AfterCreate 创建后的钩子
// func (m *User) AfterCreate(tx *gorm.DB) error {
//     return nil
// }

// Update 相关钩子函数：
// BeforeUpdate 更新前的钩子
// func (m *User) BeforeUpdate(tx *gorm.DB) error {
//     return nil
// }
//
// AfterUpdate 更新后的钩子
// func (m *User) AfterUpdate(tx *gorm.DB) error {
//     return nil
// }

// Delete 相关钩子函数：
// BeforeDelete 删除前的钩子
// func (m *User) BeforeDelete(tx *gorm.DB) error {
//     return nil
// }
//
// AfterDelete 删除后的钩子
// func (m *User) AfterDelete(tx *gorm.DB) error {
//     return nil
// }

// Find 相关钩子函数：
// AfterFind 查询后的钩子
// func (m *User) AfterFind(tx *gorm.DB) error {
//     return nil
// }
`
