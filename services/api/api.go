package api

import (
	"net/http"
	"sync"

	"github.com/hun9k/gapi/conf"
	"github.com/hun9k/gapi/http/handler"
	"github.com/hun9k/gapi/log"
	"github.com/quic-go/quic-go/http3"
)

// http listen
func Listen() {
	// http diabled
	if !conf.Get[bool]("api.enable") {
		log.Info("API service is not enabled")
		return
	}

	wg := &sync.WaitGroup{}
	// http/2 or http/1.1 with no tls
	wg.Go(func() {
		log.Info("HTTP will listen", "addr", conf.Get[string]("api.addr"))
		if err := http.ListenAndServe(conf.Get[string]("api.addr"), handler.Inst(API_HANDLER_NAME)); err != nil {
			log.Error("HTTP listen error", "error", err)
		}
	})

	if conf.Get[bool]("api.tls.enable") {
		// http/2 or http/1.1 with tls
		wg.Go(func() {
			log.Info("HTTPS will listen", "addr", conf.Get[string]("api.tls.addr"))
			if err := http.ListenAndServeTLS(conf.Get[string]("api.tls.addr"), conf.Get[string]("api.tls.certfile"), conf.Get[string]("api.tls.keyfile"), handler.Inst(API_HANDLER_NAME)); err != nil {
				log.Error("HTTPS listen error", "error", err)
			}
		})
	}

	// http/3
	if conf.Get[bool]("api.http3.enable") {
		wg.Go(func() {
			log.Info("HTTP/3 will listen", "addr", conf.Get[string]("api.tls.addr"))
			if err := http3.ListenAndServeQUIC(conf.Get[string]("api.tls.addr"), conf.Get[string]("api.tls.certfile"), conf.Get[string]("api.tls.keyfile"), handler.Inst(API_HANDLER_NAME)); err != nil {
				log.Error("HTTP/3 listen error", "error", err)
			}
		})
	}

	wg.Wait()
}

const API_HANDLER_NAME = "api"
