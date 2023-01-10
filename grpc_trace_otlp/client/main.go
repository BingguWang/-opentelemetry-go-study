package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/BingguWang/opentelemetry-go-study/grpc_trace_otlp/op"
	pb "github.com/BingguWang/opentelemetry-go-study/grpc_trace_otlp/server/proto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
)

var (
	host = flag.String("host", "localhost", "listening host") // 服务的host
	port = flag.String("port", "50051", "The server port")    // 服务的port
)

func main() {
	flag.Parsed()
	addr := net.JoinHostPort(*host, *port)
	ctx := context.Background()

	// 创建导出器exporter
	exp, err := op.NewExporter(ctx)
	if err != nil {
		panic(err)
	}
	// 创建TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(op.NewResource()),
	)
	defer func() {
		_ = tp.Shutdown(ctx)
	}()
	// 设置全局的TracerProvider
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	client := pb.NewScoreServiceClient(conn)

	// 调用score服务
	ret, err := client.AddScoreByUserID(ctx, &pb.AddScoreByUserIDReq{UserID: 1})
	if err != nil {
		log.Fatalln("call score_service failed: ", err.Error())
	}
	fmt.Println(ret)
}
