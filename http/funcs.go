package http

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/hun9k/gapi/conf"
	"github.com/hun9k/gapi/log"
	"github.com/quic-go/quic-go/http3"
)

// public funcs
var Router = Handler

func Handler() *handler {
	return serviceSingle().handler
}

var _groups = map[string]*RouterGroup{}

// API router groups, single instance
func Group(path string, handlers ...gin.HandlerFunc) *RouterGroup {
	if _, exists := _groups[path]; !exists {
		_groups[path] = &RouterGroup{
			Router().Group(path, handlers...),
		}
	}

	return _groups[path]
}

// service
func Listen() []error {
	errs := []error{}

	wg := &sync.WaitGroup{}
	// http/2 or http/1.1 with no tls
	wg.Go(func() {
		log.Info("HTTP will listen", "addr", conf.Http().Addr)
		if err := http.ListenAndServe(conf.Http().Addr, Handler()); err != nil {
			errs = append(errs, err)
		}
	})

	if conf.Http().Tls.Enable {
		// http/2 or http/1.1 with tls
		wg.Go(func() {
			log.Info("HTTPS will listen", "addr", conf.Http().Tls.Addr)
			if err := http.ListenAndServeTLS(conf.Http().Tls.Addr, conf.Http().Tls.CertFile, conf.Http().Tls.KeyFile, Handler()); err != nil {
				errs = append(errs, err)
			}
		})
	}

	// http/3
	if conf.Http().Http3.Enable {
		wg.Go(func() {
			log.Info("HTTP/3 will listen", "addr", conf.Http().Tls.Addr)
			if err := http3.ListenAndServeQUIC(conf.Http().Tls.Addr, conf.Http().Tls.CertFile, conf.Http().Tls.KeyFile, Handler()); err != nil {
				errs = append(errs, err)
			}
		})
	}

	wg.Wait()

	return errs
}
