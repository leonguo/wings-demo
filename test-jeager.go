package main

import (
	"net/http"
	"github.com/leonguo/wings/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"io"
)

const app_name = "app-test"

type data struct {
	tracer opentracing.Tracer
}

type TracedServeMux1 struct {
	mux    *http.ServeMux
	tracer opentracing.Tracer
}

// NewServeMux creates a new TracedServeMux.
func NewServeMux1(tracer opentracing.Tracer) *TracedServeMux1 {
	return &TracedServeMux1{
		mux:    http.NewServeMux(),
		tracer: tracer,
	}
}

// Handle implements http.ServeMux#Handle
func (tm *TracedServeMux1) Handle(pattern string, handler http.Handler) {
	middleware := nethttp.Middleware(
		tm.tracer,
		handler,
		nethttp.OperationNameFunc(func(r *http.Request) string {
			return "HTTP " + r.Method + " " + pattern
		}))
	tm.mux.Handle(pattern, middleware)
}

// ServeHTTP implements http.ServeMux#ServeHTTP
func (tm *TracedServeMux1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tm.mux.ServeHTTP(w, r)
}
func main() {
	tracer := tracing.InitTrace(app_name)
	d := &data{tracer: tracer}
	m := NewServeMux1(d.tracer)
	m.Handle("/", http.HandlerFunc(d.hello))
	http.ListenAndServe(":9000", m.mux)
}

func (d *data) hello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dd := tracing.NewDatabase(d.tracer)
	dd.FirstFunction(ctx)
	dd.Get(ctx, "ddd231")
	dd.SecondFunction(ctx)
	io.WriteString(w, "Hello world!")

}