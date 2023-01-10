package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/BingguWang/opentelemetry-go-study/grpc_trace_jaeger/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
)

var (
	serv = flag.String("service", "score_service", "service name") // 服务名
	host = flag.String("host", "localhost", "listening host")      // 服务的host
	port = flag.String("port", "50051", "The server port")         // 服务的port
)

func main() {
	addr := net.JoinHostPort(*host, *port)

	cc, err := grpc.DialContext(context.Background(), addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	//client
	client := pb.NewScoreServiceClient(cc)

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	id, err := client.AddScoreByUserID(ctx, &pb.AddScoreByUserIDReq{UserID: 1})
	if err != nil {
		panic(err)
	}
	fmt.Println(id)

}
