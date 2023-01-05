package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/BingguWang/grpc-go-study/client/interceptor"
	mr "github.com/BingguWang/grpc-go-study/resolver"
	"github.com/BingguWang/grpc-go-study/server/utils"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
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
	resolver.Register(&mr.ExampleResolverBuilder{}) // 注册自定义的命名解析器
	conn, err := grpc.Dial(
		//*addr,
		fmt.Sprintf("%s:///%s", mr.ExampleScheme, mr.ExampleServiceName),               // Dial to "example:///resolver.example.grpc.io"
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`), // 设置负载均衡策略

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

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second).UTC()) // 设置调用的截止时间
	// 可以在你想要取消RPC调用的时候调用cancel方法,那样就会通知道另一方，思考问题：context的状态是如何在客户端服务端进行的同步的？
	defer cancel()

	// 一元通信
	if _, err := client.AddScoreByUserID(ctx, &pb.AddScoreByUserIDReq{
		UserID: 1,
	}); err != nil { // 如果超时了这里会收到error code = DeadlineExceeded
		code := status.Code(err)
		if code == codes.PermissionDenied {
			fmt.Println(err.Error())
			errorStatus := status.Convert(err)
			for _, detail := range errorStatus.Details() {
				fmt.Println("--", detail)
				switch t := detail.(type) {
				case *errdetails.BadRequest_FieldViolation:
					log.Fatalf("请求错误: %v \n", t)
				default:
					log.Fatal(err)
				}
			}
		}
		log.Fatal(err)
	}

	// 服务端流通信
	//service.CallStreamScoreListByUser(ctx, client, &pb.GetScoreListByUserIDReq{
	//	UserID: 1,
	//})

	// 客户端流通信
	//service.CallStreamAddScoreByUser(ctx, client)

	// 双向流通信
	//service.CallStreamBidirectional(ctx, client)
}
