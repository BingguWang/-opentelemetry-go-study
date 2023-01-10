package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/BingguWang/opentelemetry-go-study/grpc_trace_jaeger/client/service"
	pb "github.com/BingguWang/opentelemetry-go-study/grpc_trace_jaeger/server/proto"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	jg "github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
)

var (
	serv = flag.String("service", "my_service", "service name") // 服务名
	host = flag.String("host", "localhost", "listening host")   // 服务的host
	port = flag.String("port", "50051", "The server port")      // 服务的port
)

func main() {
	addr := net.JoinHostPort(*host, *port)

	// 设置jaeger配置信息
	cfg := &jaegerConfig.Configuration{
		ServiceName: *serv, //对其发起请求的的调用链，叫什么服务
		Sampler: &jaegerConfig.SamplerConfig{ //采样策略的配置
			Type:  "const",
			Param: 1,
			/**
			"const" : 0 or 1 for always false/true respectively
			"probabilistic" sampler, a probability between 0 and 1
			"rateLimiting" sampler, the number of spans per second
			"remote" sampler, param is the same as for "probabilistic"
			*/
		},
		Reporter: &jaegerConfig.ReporterConfig{ //配置客户端如何上报trace信息，所有字段都是可选的
			LogSpans:           true,
			LocalAgentHostPort: "localhost:6831", // jaeger本地agent的地址,span会被send到此agent
		},
		//Token配置
		Tags: []opentracing.Tag{ //设置tag，token等信息可存于此
			//opentracing.Tag{Key: "token", Value: token}, //设置token
		},
	}
	// 根据jaeger配置来创建得到tracer
	tracer, closer, err := cfg.NewTracer(jaegerConfig.Logger(jg.StdLogger))
	defer closer.Close()
	if err != nil {
		panic(fmt.Sprintf("ERROR: fail init Jaeger: %v\n", err))
	}

	cc, err := grpc.DialContext(context.Background(), addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)),
	)
	if err != nil {
		panic(err)
	}

	//client
	client := pb.NewScoreServiceClient(cc)

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	id := service.CallScoreService(ctx, client)
	fmt.Println(id)
}
