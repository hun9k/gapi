package app

import (
	"sync"

	"github.com/hun9k/gapi/http"
)

// 运行应用
func Run() {
	wg := &sync.WaitGroup{}

	// HTTP监听
	wg.Go(func() {
		http.Listen()
	})

	// task

	wg.Wait()
}
