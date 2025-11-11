package gapi

import (
	"log/slog"

	"github.com/hun9k/gapi/http"

	"github.com/spf13/viper"
)

type Application struct {
	httpService *http.HttpService
}

// 应用模式
const (
	MODE_DEV  = "dev"
	MODE_TEST = "test"
	MODE_PROD = "prod"
)

// 全局app对象
var app *Application

func App() *Application {
	return app
}

// 创建应用
func NewApp() *Application {
	app = &Application{}

	return app
}

func (app *Application) AddHttpService() *Application {
	// HTTP服务
	if viper.GetBool("httpService.enabled") {
		app.httpService = http.NewHttpService()
	}

	return app
}

// 运行应用
func (app *Application) Run() error {
	slog.Info("application is running", "name", viper.GetString("name"))

	// HTTP服务监听
	if app.httpService != nil {
		app.httpService.Listen()
	}

	return nil
}
