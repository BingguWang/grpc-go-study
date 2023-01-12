package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"sync"

	pb "github.com/BingguWang/hystrix-study/retry_test/server/proto"
)

var port = flag.Int("port", 50051, "port number")

type failingServer struct {
	pb.UnimplementedScoreServiceServer
	mu         sync.Mutex
	reqCounter uint // 请求计数器
}

func (s *failingServer) AddScoreByUserID(ctx context.Context, req *pb.AddScoreByUserIDReq) (*pb.AddScoreByUserIDResp, error) {
	s.mu.Lock()
	s.reqCounter++
	s.mu.Unlock()
	if s.reqCounter == 1 {
		log.Printf("第 %v 次 请求失败...failed \n", s.reqCounter)
		// 一个标准的grpc错误至少应该这样写，有code码，和message
		return nil, status.Errorf(codes.Unavailable, "call AddScoreByUserID failed")
	}
	log.Printf("第 %v 次 请求成功...succeed \n", s.reqCounter)
	return &pb.AddScoreByUserIDResp{UserID: req.UserID}, nil
}

func main() {
	flag.Parse()

	address := fmt.Sprintf(":%v", *port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("listen on address", address)

	s := grpc.NewServer()

	failingservice := &failingServer{reqCounter: 0}

	pb.RegisterScoreServiceServer(s, failingservice)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
