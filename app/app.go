package app

import (
	"sync"

	"github.com/hun9k/gapi/conf"
	"github.com/hun9k/gapi/services/api"
)

// 运行应用
func Run() {
	// run
	wg := &sync.WaitGroup{}

	// dev mode
	wg.Go(func() {
		if conf.Get[string]("app.mode") == conf.APP_MODE_DEV {
		}
	})

	// api service
	wg.Go(func() {
		api.Listen()
	})

	// task service
	wg.Go(func() {
		// TODO
	})

	// wg wait
	wg.Wait()
}
