package gapi

type Application struct{}

// 全局_app对象
var _app *Application

func App() *Application {
	if _app == nil {
		_app = newApp()
	}

	return _app
}

// 创建应用
func newApp() *Application {
	_app = &Application{}
	return _app
}

// 运行应用
func (app *Application) Run() error {
	// HTTP服务监听
	if Conf().HttpService.Enable {
		Log().Debug("HTTP service is listening", "addr", Conf().HttpService.Addr)
		HttpSvc().Run(Conf().HttpService.Addr)
	}

	return nil
}
