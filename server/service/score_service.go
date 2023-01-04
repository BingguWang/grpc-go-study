package service

import (
	"context"
	"fmt"
	pb "github.com/BingguWang/grpc-go-study/server/proto"
	"github.com/BingguWang/grpc-go-study/server/utils"
	"io"
	"log"
	"sync"
	"time"
)

var (
	serverInstance *server
	once           sync.Once
)

func GetServer() *server {
	once.Do(func() {
		serverInstance = &server{}
	})
	return serverInstance
}

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedScoreServiceServer // 所有的实现类必须内嵌此结构，为了实现向前兼容
}

// 实现存根方法
func (*server) AddScoreByUserID(ctx context.Context, in *pb.AddScoreByUserIDReq) (*pb.AddScoreByUserIDResp, error) {
	log.Println("call AddScoreByUserID...")
	return &pb.AddScoreByUserIDResp{UserID: in.UserID}, nil
}

func (*server) AddStreamScoreByUserID(stream pb.ScoreService_AddStreamScoreByUserIDServer) error {
	log.Println("call AddStreamScoreByUserID...")
	var count int
	for {
		// 从客户端发送的流内接收请求，这里grpc可以保证接收的顺序好客户端发送请求的顺序是一致的
		req, err := stream.Recv()
		if err == io.EOF {
			// 发送响应， 这里选择的是在请求接收完时才发送响应
			fmt.Println("count: ", count)
			return stream.SendAndClose(&pb.AddScoreByUserIDResp{UserID: 1})
		}
		if err != nil {
			return err
		}
		fmt.Println(req.Scores[0].Type)
		fmt.Println(req.Scores[0].Value)
		fmt.Println(req.Scores)
		count++
	}
}

// GetStreamScoreListByUser 服务端流式
func (*server) GetStreamScoreListByUser(in *pb.GetScoreListByUserIDReq, stream pb.ScoreService_GetStreamScoreListByUserServer) error {
	log.Println("call GetStreamScoreListByUser...")
	arr := []*pb.GetScoreListByUserIDResp{
		{
			UserID: 1,
			Scores: []*pb.Score{
				{Type: 1, Value: 100},
				{Type: 2, Value: 120},
				{Type: 3, Value: 130},
			},
		},
		{
			UserID: 2,
			Scores: []*pb.Score{
				{Type: 11, Value: 200},
				{Type: 22, Value: 220},
				{Type: 33, Value: 230},
			},
		},
		{
			UserID: 3,
			Scores: []*pb.Score{
				{Type: 11, Value: 300},
				{Type: 22, Value: 320},
			},
		},
	}
	for _, v := range arr {
		if err := stream.SendMsg(v); err != nil {
			return err
		}
	}
	return nil
}

// AddAndGetScore 双向流
func (*server) AddAndGetScore(stream pb.ScoreService_AddAndGetScoreServer) error {
	log.Println("call AddAndGetScore...")
	lastest := &pb.GetScoreListByUserIDResp{}
	for {
		// 从客户端发送的流接收请求
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		// 获取并新增分数
		fmt.Println("in: ", utils.ToJsonString(in))
		fmt.Println("server recv time: ", time.Now())

		// 返回最新的分数到响应流
		lastest.UserID = in.UserID
		lastest.Scores = append(lastest.Scores, in.Scores...)
		if err := stream.Send(lastest); err != nil {
			return err
		}
	}
}
