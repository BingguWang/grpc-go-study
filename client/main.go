package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/BingguWang/grpc-go-study/client/interceptor"
	"github.com/BingguWang/grpc-go-study/etcdv3"
	pb "github.com/BingguWang/grpc-go-study/server/proto"
	"github.com/BingguWang/grpc-go-study/server/utils"
	"github.com/patrickmn/go-cache"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	svc  = flag.String("service", "score_service", "service name")
	reg  = flag.String("reg", "http://localhost:2379", "register etcd address") // 注册中心etcd的地址
	//reg  = flag.String("reg", "http://127.0.0.1:2379,http://127.0.0.1:12379,http://127.0.0.1:22379", "register etcd address")
)

/**
 有了注册中心后，客户端只要知道服务名，不需要知道服务地址，解析服务名的工作也交给注册中心，客户端不再需要知道服务名-地址的映射关系
客户把服务名给注册中心，由注册中心去解析出服务地址
*/
var clientCounterCache *cache.Cache

const (
	succeedKey = "succeed"
	failedKey  = "failed"
	limitedKey = "limited"
)

func init() {
	// 初始化一个请求计数器
	if clientCounterCache == nil {
		clientCounterCache = cache.New(cache.NoExpiration, cache.DefaultExpiration)
	}
	clientCounterCache.Set(succeedKey, 0, cache.NoExpiration)
	clientCounterCache.Set(failedKey, 0, cache.NoExpiration)
	clientCounterCache.Set(limitedKey, 0, cache.NoExpiration)
}
func main() {
	flag.Parse()
	// Set up a connection to the server.
	fmt.Println(utils.ToJsonString(addr))

	//resolver.Register(&mr.ExampleResolverBuilder{}) // 注册自定义的grpc命名解析器
	// 使用自定义的etcd命名解析器
	r := etcdv3.NewResolver(*reg, *svc)
	resolver.Register(r)

	opts := utils.GetOneSideTlsClientOpts()
	conn, err := grpc.Dial(
		//*addr,
		r.Scheme()+"://authority/"+*svc, // etcd的命名解析，格式要写对 scheme名称://authority/servicename
		//fmt.Sprintf("%s:///%s", mr.ExampleScheme, mr.ExampleServiceName),               // grpc的命名解析
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`), // 设置负载均衡策略

		//grpc.WithTransportCredentials(insecure.NewCredentials()), // 没有认证
		grpc.WithUnaryInterceptor(interceptor.MyUnaryClientInterceptor),   // 设置客户端一元拦截器
		grpc.WithStreamInterceptor(interceptor.MyStreamClientInterceptor), // 设置客户端流拦截器
		opts[0],
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

	// 测试限流器是否生效
	for i := 0; i < 200; i++ {
		// 一元通信RPC调用
		if _, err := client.AddScoreByUserID(ctx, &pb.AddScoreByUserIDReq{
			UserID: 1,
		}); err != nil { // 如果超时了这里会收到error code = DeadlineExceeded
			code := status.Code(err)
			switch code {
			case codes.PermissionDenied:
				fmt.Println(err.Error())
				errorStatus := status.Convert(err)
				for _, detail := range errorStatus.Details() {
					fmt.Println("--", detail)
					switch t := detail.(type) {
					case *errdetails.BadRequest_FieldViolation:
						log.Printf("请求错误: %v \n", t)
					default:
						log.Printf(err.Error())
					}
				}
			case codes.ResourceExhausted:
				clientCounterCache.Increment(limitedKey, 1)
			default:
				log.Printf(err.Error())
			}
			clientCounterCache.Increment(failedKey, 1)
			continue
		}
		clientCounterCache.Increment(succeedKey, 1)
	}

	// 服务端流通信
	//service.CallStreamScoreListByUser(ctx, client, &pb.GetScoreListByUserIDReq{
	//	UserID: 1,
	//})

	// 客户端流通信
	//service.CallStreamAddScoreByUser(ctx, client)

	// 双向流通信
	//service.CallStreamBidirectional(ctx, client)

	sucee, _ := clientCounterCache.Get(succeedKey)
	fail, _ := clientCounterCache.Get(failedKey)
	limited, _ := clientCounterCache.Get(limitedKey)
	fmt.Printf("call succeed: %v, \n call failed: %v, \n call limited: %v", sucee, fail, limited)

}
