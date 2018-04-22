package tracing

import (
	"testing"
)

func NewJaegerTest(t *testing.T) {
	opt := JaegerOptions{
		ServiceName: "logistic",
		Pretend:     true,
	}
	_, err := NewJaeger(opt)
	if err != nil {
		t.Error(err)
	}
}
