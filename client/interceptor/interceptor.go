package interceptor

import (
	"context"
	"github.com/BingguWang/grpc-go-study/server/utils"
	"google.golang.org/grpc"
	"log"
)

/**
一元拦截器，客户端调用服务存根时，
由于执行的是 Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...CallOption) error
不需要返回流，执行完invoke就得到了响应
*/

func MyUnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// 前置处理, 可以在这个阶段通过检查传入的参数来访问关于当前 RPC 的信息，
	//比如 RPC 的上下文、方法字符串、要发送的请求以及 CallOption 配置
	log.Println("======= [客户端一元拦截器] ", utils.ToJsonString(method))

	// 调用rpc， invoker方法会执行send和recv操作，
	// 执行的其实就是ClientConn的方法(cc *ClientConn) Invoke(ctx context.Context, method string, args, reply interface{}, opts ...CallOption)
	log.Println("======= 调用RPC ")
	if err := invoker(ctx, method, req, reply, cc, opts...); err != nil {
		return err
	}

	// 后置处理，在后置处理阶段，可以访问 RPC 的响应结果或错误结果,因为调用了Invoker方法之后，已经是收到了服务端的响应结果了的。
	log.Println(utils.ToJsonString(reply))
	return nil
}

/**
流拦截器，
	客户端调用服务存根时，由于只是获取返回的grpc.ClientStream，这个流返回后客户端的service里还要继续进行逻辑处理
*/

func MyStreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	log.Println("======= [客户端流拦截器] ", utils.ToJsonString(method))
	// 请求和响应是Streamer 调用 SendMsg 和 RecvMsg 这两个方法获取的。

	log.Println("======= 调用RPC ")
	stream, err := streamer(ctx, desc, cc, method)
	// 我们其实可以对获取到stream进行自定义封装
	return stream, err
}
