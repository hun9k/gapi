package gapi

import (
	"log/slog"

	"github.com/spf13/viper"
)

type App struct {
	Name        string
	Mode        string
	HttpService *HttpService
}

// 创建应用
func NewApp() *App {
	return &App{
		Name: viper.GetString("name"),
		Mode: viper.GetString("mode"),
	}
}

// 运行应用
func (app *App) Run() {
	slog.Info("application is running", "name", app.Name)

	// HTTP服务监听
	if viper.GetBool("httpService.enabled") {
		httpService := &HttpService{
			Addr: viper.GetString("httpService.addr"),
		}
		httpService.Listen()
	}
}

// 创建并运行应用
func Run() {
	NewApp().Run()
}
