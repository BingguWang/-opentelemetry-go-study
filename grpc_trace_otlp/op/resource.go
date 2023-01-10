package op

import (
	"github.com/BingguWang/opentelemetry-go-study/grpc_trace_otlp/server/service"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

// NewResource returns a resource describing this application.
func NewResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),

		resource.NewWithAttributes(
			semconv.SchemaURL, // 使用semconv包为资源属性提供常规名称。
			semconv.ServiceNameKey.String(service.Name),
			semconv.ServiceVersionKey.String("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)

	return r
}
