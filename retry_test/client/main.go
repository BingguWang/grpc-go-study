package main

import (
	"context"
	"flag"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc/codes"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/BingguWang/hystrix-study/retry_test/server/proto"
)

/**
测试grpc重试
*/
var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	// see https://github.com/grpc/grpc/blob/master/doc/service_config.md to know more about service config
	// [方式1] 重试策略配置
	retryPolicy = `{
		"methodConfig": [{
		  "name": [{"service": "proto.ScoreService"}],
		  "waitForReady": true,
		  "retryPolicy": {
			  "MaxAttempts": 4,
			  "InitialBackoff": ".01s",
			  "MaxBackoff": ".01s",
			  "BackoffMultiplier": 1.0,
			  "RetryableStatusCodes": [ "UNAVAILABLE" ]
		  }
		}]}`
	// name:["service":"服务名","method":"方法名"]
	// MaxAttempts 最大重试次数，包括首次调用本身，该字段的值必须大于1，超过5会当做5处理。
	// InitialBackoff 指数退避参数，首次重试会在random(0, initial_backoff)时间后触发
	// 之后第n次重试会在random(0, min(InitialBackoff * BackoffMultiplier**(n-1), max_backoff))时间后触发
	// 只有在返回的错误码是RetryableStatusCodes中罗列的才会重试

	// [方式2] 重试策略配置重试策略配置
	retryOpt = []grpc_retry.CallOption{
		grpc_retry.WithMax(3),                                                          // 最多重试3次
		grpc_retry.WithPerRetryTimeout(3 * time.Second),                                // 每次重试3秒超时
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(100 * time.Millisecond)),       // 线性退避重试，首次重试间隔100毫秒
		grpc_retry.WithCodes(codes.Unavailable, codes.Aborted, codes.DeadlineExceeded), // 仅当grpc响应状态码是列出来的才重试
	}
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(
		*addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		//grpc.WithDefaultServiceConfig(retryPolicy), // [方式1]
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpt...)),   // 重试拦截器 [方式2]
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(retryOpt...)), // 重试拦截器 [方式2]
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() {
		if e := conn.Close(); e != nil {
			log.Printf("failed to close connection: %s", e)
		}
	}()

	c := pb.NewScoreServiceClient(conn)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second).UTC())
	defer cancel()

	reply, err := c.AddScoreByUserID(ctx, &pb.AddScoreByUserIDReq{UserID: 1})
	if err != nil {
		log.Fatalf("AddScoreByUserID error: %v", err)
	}
	log.Printf("AddScoreByUserID reply: %v", reply)
}
