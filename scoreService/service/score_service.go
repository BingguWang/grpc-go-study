package service

import (
    "context"
    "encoding/json"
    pb "github.com/BingguWang/grpc-go-study/scoreService/proto"
    "log"
)

type ScoreServiceImpl struct {
    pb.UnimplementedScoreServiceServer // 必须内嵌次结构
}

// @alias =/score/add/byUser
func (s *ScoreServiceImpl) AddScoreByUserID(ctx context.Context, req *pb.AddScoreByUserIDReq) (*pb.AddScoreByUserIDResp, error) {
    log.Println("AddScoreByUserID, req: ", ToJsonString(req))
    // do something...
    return &pb.AddScoreByUserIDResp{}, nil
}

// @alias =/score/list
func (s *ScoreServiceImpl) GetScoreListByUser(ctx context.Context, req *pb.GetScoreListByUserIDReq) (*pb.GetScoreListByUserIDResp, error) {
    log.Println("AddScoreByUserID, req: ", ToJsonString(req))
    // do something...
    return &pb.GetScoreListByUserIDResp{}, nil
}

func ToJsonString(v interface{}) string {
    marshal, _ := json.Marshal(v)
    return string(marshal)
}
