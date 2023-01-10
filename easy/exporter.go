package main

import (
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
	"io"
)

/**
导出器，
导出器是 openTelemetry的三个组件之一，exporter这个组件一般是由第三方库实现
*/

// 这里我们使用输出到控制台的exporter：stdouttrace, 导出到console，有更高级的需求可以使用别的导出器，比如jaeger,zipkin，OTLP

// 返回一个 console exporter.
func newExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}
