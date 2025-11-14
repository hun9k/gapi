package gapi

type Application struct{}

// 全局app对象
var app *Application

func App() *Application {
	if app == nil {
		app = NewApp()
	}

	return app
}

// 创建应用
func NewApp() *Application {
	app = &Application{}
	return app
}

// 运行应用
func (app *Application) Run() error {
	Log().Info("application is running", "name", Conf().App.Name)

	// HTTP服务监听
	if Conf().HttpService.Enabled {
		httpService.Listen()
	}

	return nil
}
