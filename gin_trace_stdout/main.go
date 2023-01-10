package main

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
	"io"
	"log"
	"net/http"
	"os"
)

// 返回一个 console exporter.
func newExporter(w io.Writer) (tracesdk.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}
func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	exporter, err := newExporter(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithBatcher(exporter),
	)
	if err != nil {
		log.Fatal(err)
	}
	return tp, nil
}

const serviceName = "gin-service"

func main() {
	// 设置一个provider，这里我们用console exporter直接打印结果到控制台
	tp, err := tracerProvider("")
	if err != nil {
		panic(err)
	}
	// 设置为全局
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	r := gin.Default()
	// 关键是使用otelgin中间件
	r.Use(otelgin.Middleware(serviceName))
	r.GET("/ping", func(c *gin.Context) {
		_, span := otel.Tracer("gin-server").Start(c.Request.Context(), "ping", oteltrace.WithAttributes(attribute.String("id", "1")))
		defer span.End()
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
