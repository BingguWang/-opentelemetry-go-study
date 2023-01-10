package main

import (
	"context"
	"flag"
	"github.com/BingguWang/opentelemetry-go-study/grpc_trace_otlp/op"
	pb "github.com/BingguWang/opentelemetry-go-study/grpc_trace_otlp/server/proto"
	"github.com/BingguWang/opentelemetry-go-study/grpc_trace_otlp/server/service"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel/propagation"

	lambdadetector "go.opentelemetry.io/contrib/detectors/aws/lambda"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
)

/**
测试grpc-go结合 openTelemetry , otlp export

//TODO 暂时没有成功
*/
var (
	serv = flag.String("service", "score_service", "service name") // 服务名
	host = flag.String("host", "localhost", "listening host")      // 服务的host
	port = flag.String("port", "50051", "The server port")         // 服务的port
)

// grpc有拦截器，所以我们把trace的逻辑写到拦截器里就行
func main() {
	flag.Parse()
	ctx := context.Background()

	// 创建导出器exporter
	exp, err := op.NewExporter(ctx)
	if err != nil {
		panic(err)
	}

	// resource
	detector := lambdadetector.NewResourceDetector()
	resource, _ := detector.Detect(ctx)

	// 创建TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exp),
		//sdktrace.WithResource(op.NewResource()),
		sdktrace.WithResource(resource),
		sdktrace.WithIDGenerator(xray.NewIDGenerator()),
	)
	defer func() {
		_ = tp.Shutdown(ctx)
	}()

	// 设置全局的TracerProvider
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	addr := net.JoinHostPort(*host, *port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor(otelgrpc.WithTracerProvider(tp))), //设置拦截器进行埋点
		//grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor(otelgrpc.WithTracerProvider(tp))),
	)

	// 注册服务
	pb.RegisterScoreServiceServer(grpcServer, &service.ScoreServiceImpl{})
	log.Println("listen on : ", addr)

	// 开启
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}

}
