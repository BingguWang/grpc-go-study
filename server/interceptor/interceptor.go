package interceptor

import (
	"context"
	lt "github.com/BingguWang/grpc-go-study/server/limiter"
	"github.com/BingguWang/grpc-go-study/server/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

/**
拦截器
	拦截器可以设在客户端也可以设在服务端

*/

// MyUnaryServerInterceptor 服务端一元拦截器
func MyUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	//通过检查传入的参数，获取关于当前RPC的信息
	// ....前置逻辑
	log.Println("======= [服务端一元拦截器]", utils.ToJsonString(info))
	// 新增限流器逻辑
	log.Println("...........获取限流器...........")
	limiter := lt.NewTimeRateLimiter()
	c, _ := context.WithTimeout(ctx, time.Millisecond*200) // 最多只会等待这个时间,可以设置时间,请求就会一直等到有令牌为止
	if err := limiter.Wait(c); err != nil {
		log.Println("被限流...")
		return nil, status.Errorf(codes.ResourceExhausted, "client exceeded rate limit")
	}

	// 调用handler完成一元RPC的正常执行
	log.Println("======= [调用服务执行handler]")
	m, err := handler(ctx, req) // 服务端调用的，有拦截器了，就走的是这里的逻辑（serviceName_method_Handler可以看到相关逻辑）
	// m就是响应给客户端的

	// ....后置逻辑，可以处理RPC响应
	log.Printf(" Post Proc Message : %s", m)
	return m, err
}

// MyStreamServerInterceptor 服务端流拦截器
// 参数grpc.ServerStream 就是实现方法的时候传入的参数stream，他有SendMsg 和RecvMsg 方法
// 可以在grpc.ServerStream的基础上封装自己的ServerStream
func MyStreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// 前置处理
	log.Println("======= [服务端流拦截器] ", utils.ToJsonString(info))

	// 调用服务
	log.Println("======= [调用服务执行handler]")
	err := handler(srv, ss)
	if err != nil {
		log.Printf("handler err  : %s \n", err.Error())
	}

	// 后置处理，如果响应信息不想在handler里发，可以放到这里的逻辑里发,也就是说可以不在你的实现方法里调用stream.SendMsg ，可以放到这里来调stream.SendMsg 发送响应
	log.Printf(" Post Proc Message : %s \n", utils.ToJsonString(info))
	return err
}
