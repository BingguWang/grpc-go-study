package main

import (
    pb "github.com/BingguWang/grpc-go-study/scoreService/proto"
    service "github.com/BingguWang/grpc-go-study/scoreService/service"
    "google.golang.org/grpc"
    "net"
)

const address = ":50002"

// 启动积分项目
func main() {
    listener, err := net.Listen("tcp", address)
    if err != nil {
        panic(err)
    }
    server := grpc.NewServer()
    // 注册服务
    pb.RegisterScoreServiceServer(server, &service.ScoreServiceImpl{})
    if err := server.Serve(listener); err != nil {
        panic(err)
    }
}
