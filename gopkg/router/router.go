package router

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	traceLog "github.com/opentracing/opentracing-go/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/lab46/monorepo/gopkg/http/httputil"
	"github.com/lab46/monorepo/gopkg/log"
	"github.com/lab46/monorepo/gopkg/tracing"
)

// Router of vehicle-insurance
type Router struct {
	opt Options
	r   *mux.Router
}

// TimeoutOptions struct
type TimeoutOptions struct {
	Duration time.Duration
	Response interface{}
}

// RecoverOptions struct
type RecoverOptions struct {
}

var defaultTimeoutResponse []byte

// Options for router
type Options struct {
	Timeout TimeoutOptions
}

var once sync.Once

func init() {
	once.Do(func() {
		defaultTimeoutResponse, _ = json.Marshal(map[string]interface{}{
			"errors": []string{"request timed out"},
		})
	})
}

// New router
func New(opt Options) *Router {
	muxRouter := mux.NewRouter()
	rtr := &Router{
		r:   muxRouter,
		opt: opt,
	}
	return rtr
}

// URLParam get param from rest request
func URLParam(r *http.Request, key string) string {
	params := mux.Vars(r)
	return params[key]
}

// SubRouter return a new Router with path prefix
func (rtr Router) SubRouter(pathPrefix string) *Router {
	muxSubrouter := rtr.r.PathPrefix(pathPrefix).Subrouter()
	return &Router{
		r:   muxSubrouter,
		opt: rtr.opt,
	}
}

// handle middleware
// the timeout middleware should cover timeout budget
func (rtr *Router) handle(h http.HandlerFunc, traceOperation ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opt := rtr.opt
		// cancel context
		if opt.Timeout.Duration > 0 {
			ctx, cancel := context.WithTimeout(r.Context(), opt.Timeout.Duration*time.Second)
			defer cancel()
			r = r.WithContext(ctx)
		}

		// trace the request
		if len(traceOperation) > 0 {
			span, ctx := tracing.StartSpanFromHTTPRequest(r, traceOperation[0])
			defer span.Finish()
			span.LogFields(traceLog.String("IP", httputil.GetClientIPAddress(r)))
			r = r.WithContext(ctx)
		}

		doneChan := make(chan bool)
		go func() {
			h(w, r)
			doneChan <- true
		}()
		select {
		case <-r.Context().Done():
			jsonResp, err := json.Marshal(rtr.opt.Timeout.Response)
			if err != nil {
				log.Errorf("[router][timeout] failed to marshal response for timeout. Err: %s", err.Error())
				jsonResp = defaultTimeoutResponse
			}
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write(jsonResp)
			return
		case <-doneChan:
			return
		}
	}
}

func sanitizeStatusCode(status int) string {
	code := strconv.Itoa(status)
	return code
}

// Get function
func (rtr Router) Get(pattern string, h http.HandlerFunc, traceOperation ...string) {
	log.Debugf("[router][get] %s", pattern)
	rtr.r.HandleFunc(pattern, prometheus.InstrumentHandlerFunc(pattern, rtr.handle(h, traceOperation...))).Methods("GET")
}

// Post function
func (rtr Router) Post(pattern string, h http.HandlerFunc, traceOperation ...string) {
	log.Debugf("[router][post] %s", pattern)
	rtr.r.HandleFunc(pattern, prometheus.InstrumentHandlerFunc(pattern, rtr.handle(h, traceOperation...))).Methods("POST")
}

// Put function
func (rtr Router) Put(pattern string, h http.HandlerFunc, traceOperation ...string) {
	log.Debugf("[router][put] %s", pattern)
	rtr.r.HandleFunc(pattern, prometheus.InstrumentHandlerFunc(pattern, rtr.handle(h, traceOperation...))).Methods("PUT")
}

// Delete function
func (rtr Router) Delete(pattern string, h http.HandlerFunc, traceOperation ...string) {
	log.Debugf("[router][delete] %s", pattern)
	rtr.r.HandleFunc(pattern, prometheus.InstrumentHandlerFunc(pattern, rtr.handle(h, traceOperation...))).Methods("DELETE")
}

// Patch function
func (rtr Router) Patch(pattern string, h http.HandlerFunc, traceOperation ...string) {
	log.Debugf("[router][patch] %s", pattern)
	rtr.r.HandleFunc(pattern, prometheus.InstrumentHandlerFunc(pattern, rtr.handle(h, traceOperation...))).Methods("PATCH")
}

// Handle function
func (rtr Router) Handle(pattern string, h http.Handler) {
	log.Debugf("[router][handle] %s", pattern)
	rtr.r.Handle(pattern, h)
}

// ServeHTTP function
func (rtr Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rtr.r.ServeHTTP(w, r)
}
