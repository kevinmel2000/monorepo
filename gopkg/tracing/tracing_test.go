package tracing

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestStartSpanFromHTTPRequest(t *testing.T) {
	r, err := http.NewRequest("GET", "localhost", nil)
	if err != nil {
		t.Error(err)
	}
	rCtx := context.WithValue(r.Context(), "key1", "value1")
	r = r.WithContext(rCtx)

	span, ctx := StartSpanFromHTTPRequest(r, "operation_name")
	defer span.Finish()
	// expecting same value from context
	val := ctx.Value("key1")
	if val == nil {
		t.Error("value is nil")
	}
}

func TestSpanFromContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), "key1", "value1")
	// replace ctx with ctx from StartSpanFromContext
	span, ctx := StartSpanFromContext(ctx, "operation_name")
	defer span.Finish()
	// expecting same value from context
	val := ctx.Value("key1")
	if val == nil {
		t.Error("value is nil")
	}
}

func TEestLogError(t *testing.T) {
	ctx := context.WithValue(context.Background(), "key1", "value1")
	// replace ctx with ctx from StartSpanFromContext
	span, ctx := StartSpanFromContext(ctx, "operation_name")
	defer span.Finish()
	// expecting same value from context
	val := ctx.Value("key1")
	if val == nil {
		t.Error("value is nil")
	}
	LogError(span, errors.New("some error"))
}
