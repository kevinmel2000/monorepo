package tracing

import (
	"errors"
	"fmt"
	"io"

	"github.com/lab46/monorepo/gopkg/env"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

// Jaeger struct
type Jaeger struct {
	tracer io.Closer
	config config.Configuration
}

// Close tracer
func (j *Jaeger) Close() error {
	return j.tracer.Close()
}

// JaegerOptions for jaeger
type JaegerOptions struct {
	ServiceName string
	// Pretend only pretend open connection to jaeger
	Pretend bool
	// address can be empty, the default is localhost
	// specifying the address is useful if instance don't have its own jaeger agent
	Address string
	// Env can be empty, the default value is env.GetCurrentServiceEnv()
	Env string
}

// Validate jaeger options
func (o *JaegerOptions) Validate() error {
	if o.ServiceName == "" {
		return errors.New("service name cannot be empty")
	}
	if o.Env == "" {
		o.Env = env.GetCurrentServiceEnv()
	}
	return nil
}

// NewJaeger returns jaeger for tracing
func NewJaeger(opt JaegerOptions) (*Jaeger, error) {
	err := opt.Validate()
	if err != nil {
		return nil, err
	}
	// exit if just pretending
	if opt.Pretend {
		return nil, nil
	}

	sampler := config.SamplerConfig{
		Type:  jaeger.SamplerTypeConst,
		Param: 1,
	}

	conf := config.Configuration{
		Sampler: &sampler,
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}
	logger := jaeger.NullLogger
	metricsFactory := metrics.NullFactory

	tracer, err := conf.InitGlobalTracer(
		fmt.Sprintf("%s-%s", opt.ServiceName, opt.Env),
		config.Logger(logger),
		config.Metrics(metricsFactory),
	)
	if err != nil {
		return nil, err
	}

	j := Jaeger{
		tracer: tracer,
		config: conf,
	}
	return &j, nil
}
