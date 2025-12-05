package app

import (
	"sync"

	"github.com/hun9k/gapi/services/api"
)

// 运行应用
func Run() {

	// run services
	wg := &sync.WaitGroup{}

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
