package webserver

import (
	"net/http"
	"strconv"
	"time"

	"github.com/lab46/example/pkg/log"
	"github.com/lab46/example/pkg/router"
	"github.com/prometheus/client_golang/prometheus"
)

type Options struct {
	Port    int
	Timeout time.Duration
}

type WebServer struct {
	router *router.Router
	port   string
	birth  time.Time
}

func checkOptions(opt *Options) {
	if opt.Timeout == time.Duration(0) {
		opt.Timeout = time.Second * 3
	}
	if opt.Port == 0 {
		opt.Port = 9000
	}
}

func New(opt Options) WebServer {
	checkOptions(&opt)
	port := ":" + strconv.Itoa(opt.Port)

	r := router.New(router.Options{Timeout: opt.Timeout})
	// provide metrics endpoint for prometheus metrics
	r.Handle("/metrics", prometheus.Handler())
	// provide service healthcheck
	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	web := WebServer{
		router: r,
		port:   port,
		birth:  time.Now(),
	}
	return web
}

func (w *WebServer) Router() *router.Router {
	return w.router
}

func (w *WebServer) Run() error {
	log.Infof("Webserver running on port %s", w.port)
	return http.ListenAndServe(w.port, w.router)
}
