package main

import (
	"net/http"
	"github.com/leonguo/wings/tracing"
	"github.com/opentracing/opentracing-go"
	"io"
)

const app_name = "app-test"

type data struct {
	tracer opentracing.Tracer
}

func main() {
	tracer := tracing.InitTrace(app_name)
	d := &data{tracer: tracer}
	m := tracing.NewServeMux(d.tracer)
	m.Handle("/", http.HandlerFunc(d.hello))
	m.Handle("/test", http.HandlerFunc(d.test))
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

func (d *data) test(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "nothing Hello world!")
}