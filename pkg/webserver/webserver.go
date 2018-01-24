package webserver

import (
	"encoding/json"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/lab46/example/pkg/env"
	"github.com/lab46/example/pkg/log"
	"github.com/lab46/example/pkg/router"
	"github.com/prometheus/client_golang/prometheus"
)

type Options struct {
	Port    string
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
	if len(opt.Port) > 0 {
		// check if first string is ":"
		if opt.Port[:1] != ":" {
			opt.Port = ":" + opt.Port
		}
	} else {
		opt.Port = ":9000"
	}
}

func New(opt Options) *WebServer {
	checkOptions(&opt)
	port := opt.Port

	r := router.New(router.Options{Timeout: opt.Timeout})
	// provide metrics endpoint for prometheus metrics
	r.Handle("/metrics", prometheus.Handler())
	// provide service healthcheck
	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		currentEnv := env.GetCurrentServiceEnv()
		configDir := env.GetConfigDir()
		logLevel := log.GetLevel()

		response := map[string]interface{}{
			"environemnt": currentEnv,
			"config":      configDir,
			"log_level":   logLevel,
		}
		jsonResp, _ := json.Marshal(response)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
	})
	web := WebServer{
		router: r,
		port:   port,
		birth:  time.Now(),
	}
	return &web
}

func (w *WebServer) Router() *router.Router {
	return w.router
}

func (w *WebServer) Run() error {
	log.Infof("Webserver running on port %s", w.port)
	return http.ListenAndServe(w.port, w.router)
}
