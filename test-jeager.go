package main

import (
	"net/http"
	"github.com/leonguo/wings"
	"github.com/leonguo/wings/tracing"
	"github.com/opentracing/opentracing-go"
	"io"

)

const app_name = "app-test"

type data struct {
	tracer opentracing.Tracer
}

func main() {
	tracer := tracing.Init(app_name)
	d := &data{tracer: tracer}
	app := wings.New(d.tracer)
	app.Handle("/", http.HandlerFunc(d.hello))
	app.Handle("/test", http.HandlerFunc(d.test))
	http.ListenAndServe(":9000", app.Mux)
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