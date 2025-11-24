package app

import (
	"github.com/hun9k/gapi/conf"
	"github.com/hun9k/gapi/http"
	"github.com/hun9k/gapi/log"
)

// 运行应用
func Run() {

	// HTTP服务监听
	if conf.Http().Enable {
		if errs := http.Listen(); len(errs) > 0 {
			log.Error("http listen error", "errors", errs)
		}
	}
}
