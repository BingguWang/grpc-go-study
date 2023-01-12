package service

import (
	"context"
	"fmt"
	pb "github.com/BingguWang/grpc-go-study/server/proto"
	"github.com/BingguWang/grpc-go-study/server/utils"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	//time.Sleep(5*time.Second) // 延时，测试客户调用设置的截止时间是否生效, result:生效
	if in.UserID == 1 {
		errStatus := status.New(codes.PermissionDenied, "权限拒绝") // 利用grpc的状态包自定义错误
		details, err := errStatus.WithDetails(
			&errdetails.BadRequest_FieldViolation{
				Field:       "UserID",
				Description: fmt.Sprintf("UserID为%v的用户不是靓仔", in.UserID),
			},
		)
		if err != nil {
			return nil, err
		}
		//return nil, errors.New(details.Message() + utils.ToJsonString(details.Details()))
		return nil, details.Err()
	}
	return &pb.AddScoreByUserIDResp{UserID: in.UserID}, nil
}

func (*server) AddStreamScoreByUserID(stream pb.ScoreService_AddStreamScoreByUserIDServer) error {
	log.Println("call AddStreamScoreByUserID...")
	var count int
	time.Sleep(5 * time.Second) // 延时，测试客户端调用时设置的截止时间是否生效， result:生效
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
	//time.Sleep(5 * time.Second) // 延时，测试客户端调用时设置的截止时间是否生效， result:生效
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
	//time.Sleep(5 * time.Second) // 延时，测试客户端调用时设置的截止时间是否生效
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
