package service

import (
	"context"
	pb "github.com/BingguWang/opentelemetry-go-study/grpc_trace_jaeger/server/proto"
	"go.opentelemetry.io/otel"
)

func CallScoreService(ctx context.Context, client pb.ScoreServiceClient) *pb.AddScoreByUserIDResp {
	newCtx, span := otel.Tracer("").Start(ctx, "CallScoreService")
	defer span.End()
	id, err := client.AddScoreByUserID(newCtx, &pb.AddScoreByUserIDReq{UserID: 1})
	if err != nil {
		panic(err)
	}
	return id
}
