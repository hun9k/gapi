package http

import (
	"net/http"
	"sync"

	"github.com/hun9k/gapi/conf"
	"github.com/hun9k/gapi/log"
	"github.com/quic-go/quic-go/http3"
)

// http listen
func Listen() {
	// http diabled
	if !conf.Http().Enable {
		log.Info("HTTP is disabled")
		return
	}

	wg := &sync.WaitGroup{}
	// http/2 or http/1.1 with no tls
	wg.Go(func() {
		log.Info("HTTP will listen", "addr", conf.Http().Addr)
		if err := http.ListenAndServe(conf.Http().Addr, handlerSingle()); err != nil {
			log.Error("HTTP listen error", "error", err)
		}
	})

	if conf.Http().Tls.Enable {
		// http/2 or http/1.1 with tls
		wg.Go(func() {
			log.Info("HTTPS will listen", "addr", conf.Http().Tls.Addr)
			if err := http.ListenAndServeTLS(conf.Http().Tls.Addr, conf.Http().Tls.CertFile, conf.Http().Tls.KeyFile, handlerSingle()); err != nil {
				log.Error("HTTPS listen error", "error", err)
			}
		})
	}

	// http/3
	if conf.Http().Http3.Enable {
		wg.Go(func() {
			log.Info("HTTP/3 will listen", "addr", conf.Http().Tls.Addr)
			if err := http3.ListenAndServeQUIC(conf.Http().Tls.Addr, conf.Http().Tls.CertFile, conf.Http().Tls.KeyFile, handlerSingle()); err != nil {
				log.Error("HTTP/3 listen error", "error", err)
			}
		})
	}

	wg.Wait()
}
