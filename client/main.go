package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/BingguWang/grpc-go-study/client/interceptor"
	"github.com/BingguWang/grpc-go-study/client/service"
	"github.com/BingguWang/grpc-go-study/server/utils"
	"log"
	"time"

	pb "github.com/BingguWang/grpc-go-study/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	fmt.Println(utils.ToJsonString(addr))
	conn, err := grpc.Dial(
		*addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptor.MyUnaryClientInterceptor),   // 设置客户端一元拦截器
		grpc.WithStreamInterceptor(interceptor.MyStreamClientInterceptor), // 设置客户端流拦截器
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewScoreServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 一元通信
	//client.AddScoreByUserID(ctx, &pb.AddScoreByUserIDReq{
	//	UserID: 1,
	//})

	// 服务端流通信
	//service.CallStreamScoreListByUser(ctx, client, &pb.GetScoreListByUserIDReq{
	//	UserID: 1,
	//})

	// 客户端流通信
	//service.CallStreamAddScoreByUser(ctx, client)

	// 双向流通信
	service.CallStreamBidirectional(ctx, client)
}
