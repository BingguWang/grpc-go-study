package main

import (
	"flag"
	"fmt"
	"github.com/BingguWang/grpc-go-study/etcdv3"
	"github.com/BingguWang/grpc-go-study/server/interceptor"
	"github.com/BingguWang/grpc-go-study/server/service"
	"github.com/BingguWang/grpc-go-study/server/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/BingguWang/grpc-go-study/server/proto"
	"google.golang.org/grpc"
)

var (
	serv = flag.String("service", "score_service", "service name")              //服务名
	host = flag.String("host", "localhost", "listening host")                   // 服务的host
	port = flag.String("port", "50051", "The server port")                      // 服务的port
	reg  = flag.String("reg", "http://localhost:2379", "register etcd address") // 注册中心etcd的地址
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", net.JoinHostPort(*host, *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 服务端注册服务
	err = etcdv3.Register(*reg, *serv, *host, *port, time.Second*10, 150000)
	if err != nil {
		panic(err)
	}

	// 单向TLS校验, 不论是哪个客户端，只要有了公钥和服务器名的就都可以调用到服务
	cred, err := credentials.NewServerTLSFromFile(
		"/home/wangbing/grpc-test/key/server.pem",
		"/home/wangbing/grpc-test/key/server.key",
	) // 读取公钥-私钥对，返回启动TLS的证书
	if err != nil {
		panic(err)
	}
	/**
	NewServer()
	创建返回一个没有注册的服务，这个服务还没开始接收请求
	方法内核心的地方就是给server结构体的service成员初始化:
		services: make(map[string]*serviceInfo), // key 就是服务名service name
	可以看到只是初始化而已
	其中的serviceInfo,结构如下
		type serviceInfo struct {
			serviceImpl interface{} // 服务的方法的实现
			methods     map[string]*MethodDesc
			streams     map[string]*StreamDesc
			mdata       interface{}
		}
	*/
	s := grpc.NewServer(
		grpc.Creds(cred), // 传入上面创建的启动TLS的证书，，从而为所有传入的连接启用 TLS
		grpc.UnaryInterceptor(interceptor.MyUnaryServerInterceptor),   // 设置一个一元拦截器
		grpc.StreamInterceptor(interceptor.MyStreamServerInterceptor), // 设置一个流拦截器
	)

	// 服务down了之后，去删除etcd里注册的服务
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		fmt.Println("server down ... delete item in register")
		logrus.Infof("receive signal '%v'", s)
		etcdv3.UnRegister()
		os.Exit(1)
	}()

	/**
	RegisterScoreServiceServer(s grpc.ServiceRegistrar, srv ScoreServiceServer)
	注册服务
	实际上，真正最后是调用的(s *Server) RegisterService(sd *ServiceDesc, ss interface{})
		ServiceDesc是一个结构，它定义了RPC服务的规范
		ss就是你手动实现了server api的实现接口
	在这方法里会检查你是否实现了serviceDesc里的所有接口，是的就往之前server结构体的成员services里加入serviceInfo：
		info := &serviceInfo{
			serviceImpl: ss,
			methods:     make(map[string]*MethodDesc),
			streams:     make(map[string]*StreamDesc),
			mdata:       sd.Metadata,
		}
		methods和streams都封装到serviceInfo后，把info加入到services成员里
		for i := range sd.Methods {
			d := &sd.Methods[i]
			info.methods[d.MethodName] = d
		}
		for i := range sd.Streams {
			d := &sd.Streams[i]
			info.streams[d.StreamName] = d
		}
		s.services[sd.ServiceName] = info

	每个service name 只对应一个serviceInfo
	serviceDesc信息是由pb生成的，具体看ScoreService_ServiceDesc, 可以看到服务定义的具体的信息
	*/
	pb.RegisterScoreServiceServer(s, service.GetServer())
	log.Printf("server listening at %v", lis.Addr())

	// 输出注册完的serviceInfo看下
	fmt.Println(utils.ToJsonString(s.GetServiceInfo()))

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
