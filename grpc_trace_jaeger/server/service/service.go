package service

import (
	"context"
	"encoding/json"
	pb "github.com/BingguWang/opentelemetry-go-study/grpc_trace_jaeger/server/proto"
	"log"
)

type ScoreServiceImpl struct {
	pb.UnimplementedScoreServiceServer
}
type Stu struct {
	Name string
}

func (a *Stu) Say() string {
	return a.Name
}
func (s *ScoreServiceImpl) AddScoreByUserID(ctx context.Context, req *pb.AddScoreByUserIDReq) (*pb.AddScoreByUserIDResp, error) {
	log.Println("req is :", ToJsonString(req))

	//var stu *Stu
	//stu.Say() // panic null pointer

	resp := &pb.AddScoreByUserIDResp{UserID: req.UserID}
	return resp, nil
}
func ToJsonString(v interface{}) string {
	marshal, _ := json.Marshal(v)
	return string(marshal)
}
