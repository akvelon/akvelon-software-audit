package tracing

import (
	"fmt"
	"io"
	"time"

	"github.com/astaxie/beego"
	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegerClientConfig "github.com/uber/jaeger-client-go/config"
)

// InitTracer returns an instance of Jaeger Tracer
func InitTracer(service string) (opentracing.Tracer, io.Closer) {
	cfg := &jaegerClientConfig.Configuration{
		Sampler: &jaegerClientConfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegerClientConfig.ReporterConfig{
			LogSpans: true,
		},
	}

	hostPort := fmt.Sprintf("%s:%s", beego.AppConfig.String("jaeger-agent-host"), beego.AppConfig.String("jaeger-agent-port"))
	sender, err := jaeger.NewUDPTransport(hostPort, 0)

	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}

	tracer, closer, err := cfg.New(service,
		jaegerClientConfig.Reporter(
			jaeger.NewRemoteReporter(
				sender,
				jaeger.ReporterOptions.BufferFlushInterval(1*time.Second),
				jaeger.ReporterOptions.Logger(jaeger.StdLogger),
			),
		),
		jaegerClientConfig.Logger(jaeger.StdLogger))

	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}
