package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"log"
	"os"
	"os/signal"
)

/**
简单的使用openTelemetry
*/
var (
	service *EasyService
)

func init() {
	l := log.New(os.Stdout, "", 0)
	service = &EasyService{
		r: os.Stdin,
		l: l,
	}
}
func main() {
	ctx := context.Background()

	// 创建导出器exporter,这里是console explorer，直接导出到控制台
	exp, err := newExporter(os.Stdout)
	if err != nil {
		panic(err)
	}

	// 创建TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(newResource()),
	)
	defer func() {
		_ = tp.Shutdown(ctx)
	}()

	// 设置全局的TracerProvider
	otel.SetTracerProvider(tp)
	// 之后就可以在关键操作(函数里)创建tracer，然后由tracer去创建span

	// 开启业务
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	errCh := make(chan error)
	go func(c context.Context) {
		errCh <- service.EasyServiceHandler(c)
	}(ctx)

	select {
	case <-sigCh:
		fmt.Println("\ngoodbye")
		return
	case err := <-errCh:
		if err != nil {
			log.Fatal(err)
		}
	}
}
