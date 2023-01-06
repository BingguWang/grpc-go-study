package main

import (
	"flag"
	"fmt"
	pb "github.com/BingguWang/grpc-go-study/server/proto"
	"github.com/BingguWang/grpc-go-study/server/service"
	"github.com/BingguWang/grpc-go-study/server/utils"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

/**
这里是TLS双向认证
*/
func main() {
	flag.Parse()
	listen, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 双向TLS认证
	opts := utils.GetBothSideTlsServerOpts()
	server := grpc.NewServer(opts...) // 传入服务器
	pb.RegisterScoreServiceServer(server, service.GetServer())
	log.Printf("server listening at %v", listen.Addr())

	// 输出注册完的serviceInfo看下
	fmt.Println(utils.ToJsonString(server.GetServiceInfo()))

	if err := server.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
