package tracing

import (
	"context"
	"net/http"
	"runtime"
	"strings"

	opentracing "github.com/opentracing/opentracing-go"
	opentracinglog "github.com/opentracing/opentracing-go/log"
	"github.com/lab46/monorepo/gopkg/errors"
)

// StartSpanFromHTTPRequest function
func StartSpanFromHTTPRequest(r *http.Request, operationName string, startSpanOption ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	ctx := r.Context()
	span, spanCtx := opentracing.StartSpanFromContext(ctx, operationName, startSpanOption...)
	return span, spanCtx
}

// StartSpanFromContext function
func StartSpanFromContext(ctx context.Context, operationName string, startSpanOption ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	span, spanCtx := opentracing.StartSpanFromContext(ctx, operationName, startSpanOption...)
	return span, spanCtx
}

// StartSpanFromContextWithRuntime function
func StartSpanFromContextWithRuntime(ctx context.Context, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	pc, _, _, _ := runtime.Caller(1)
	operationName := runtime.FuncForPC(pc).Name()[strings.LastIndex(runtime.FuncForPC(pc).Name(), "/")+1:]
	return opentracing.StartSpanFromContext(ctx, operationName, opts...)
}

// LogError function
func LogError(span opentracing.Span, err error) {
	if err == nil {
		return
	}
	// check what is the error type
	switch err.(type) {
	case *errors.Errs:
		e := err.(*errors.Errs)
		span.LogFields(opentracinglog.String("error", err.Error()))
		span.LogKV(e.GetFields().ToArrayInterface()...)
	default:
		span.LogFields(opentracinglog.Error(err))
	}
}
