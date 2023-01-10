package op

import (
	"context"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/trace"
)

/**
导出器，
导出器是 openTelemetry的三个组件之一，exporter这个组件一般是由第三方库实现
*/

// 这里我们使用输出到控制台的exporter：OTLP, 可以根据业务需要使用别的导出器，比如jaeger,zipkin，Prometheus

// NewExporter 返回一个 exporter.
func NewExporter(ctx context.Context) (trace.SpanExporter, error) {
	//New otlp exporter
	opts := []otlptracegrpc.Option{
		//otlptracegrpc.WithEndpoint(""), // 未设置默认是localhost:4317
		otlptracegrpc.WithInsecure(),
	}
	return otlptracegrpc.New(ctx, opts...)
}
