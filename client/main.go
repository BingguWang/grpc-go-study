package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/BingguWang/grpc-go-study/client/interceptor"
	"github.com/BingguWang/grpc-go-study/client/service"
	"github.com/BingguWang/grpc-go-study/server/utils"
	"google.golang.org/grpc/credentials"
	"log"
	"time"

	pb "github.com/BingguWang/grpc-go-study/server/proto"
	"google.golang.org/grpc"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	fmt.Println(utils.ToJsonString(addr))
	creds, err := credentials.NewClientTLSFromFile(
		"/home/wangbing/grpc-test/key/server.pem",
		"x.binggu.example.com", // 填""也可以
	) // 读取并解析服务器端给的公钥证书，创建启用 TLS 的证书
	if err != nil {
		log.Fatalf("加载证书失败 %v \n", err)
	}
	conn, err := grpc.Dial(
		*addr,
		//grpc.WithTransportCredentials(insecure.NewCredentials()), // 没有认证
		grpc.WithTransportCredentials(creds),
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
