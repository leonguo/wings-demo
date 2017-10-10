package main

import (
	"net/http"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"time"
	"fmt"
	"os"
	"io"
	"github.com/leonguo/wings/util"
	"github.com/opentracing/opentracing-go"
)

const app_name = "app-test"

type TracedServeMux struct {
	tracer opentracing.Tracer
}

func main() {
	cfg := jaegerconfig.Configuration{
		Sampler: &jaegerconfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegerconfig.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  "47.90.63.215:6831",
		},
	}
	tracer, closer, err := cfg.New(
		"tt",
		jaegerconfig.Logger(jaegerlog.StdLogger),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer closer.Close()
	m := &TracedServeMux{tracer: tracer}
	http.HandleFunc("/", m.hello)
	http.ListenAndServe(":9000", nil)
}

func (d *TracedServeMux) hello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//span := tracer.StartSpan("new_span")
	//span.SetTag("test", "dd")
	//defer span.Finish()
	//opentracing.SetGlobalTracer(tracer)
	//parent := opentracing.GlobalTracer().StartSpan("hello")
	//defer parent.Finish()
	dd := util.NewDatabase(d.tracer)
	dd.Get(ctx, "21354")
	io.WriteString(w, "Hello world!")

}
