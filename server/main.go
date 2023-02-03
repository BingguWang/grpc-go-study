package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/BingguWang/grpc-go-study/etcdv3"
	"github.com/BingguWang/grpc-go-study/server/interceptor"
	"github.com/BingguWang/grpc-go-study/server/limiter"
	"github.com/BingguWang/grpc-go-study/server/service"
	"github.com/BingguWang/grpc-go-study/server/utils"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
	addr := net.JoinHostPort(*host, *port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 服务端注册服务
	if err := etcdv3.Register(*reg, *serv, *host, *port, time.Second*10, 150000); err != nil {
		panic(err)
	}

	// 单向TLS校验, 不论是哪个客户端，只要有了公钥和服务器名的就都可以调用到服务
	opts := utils.GetOneSideTlsServerOpts()
	opts = append(opts,
		grpc.InTapHandle(limiter.RateLimiter),                         // 好像限流没有生效，那就去拦截器里用
		grpc.UnaryInterceptor(interceptor.MyUnaryServerInterceptor),   // 设置一个一元拦截器
		grpc.StreamInterceptor(interceptor.MyStreamServerInterceptor), // 设置一个流拦截器
		grpc.MaxRecvMsgSize(2<<10),                                    // 设置服务器端可接收的最大请求体长度为2KB
	)
	grpcServer := grpc.NewServer(
		opts...,
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

	// 网关相对于服务端相当于是客户,把http请求转为grpc请求后通过命名解析负载均衡的方式选择grpcServer
	r := etcdv3.NewResolver(*reg, *serv)
	resolver.Register(r)
	endpoint := r.Scheme() + "://authority/" + *serv // etcd的命名解析，格式要写对 scheme名称://authority/servicename

	pb.RegisterScoreServiceServer(grpcServer, service.GetServer())
	log.Printf("server listening at %v", lis.Addr())

	// 输出注册完的serviceInfo看下
	fmt.Println(utils.ToJsonString(grpcServer.GetServiceInfo()))

	// gw server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 因为gateway是要去调grpc server的
	//所以这里gateway相对于grpc server来说是grpc的  客户端
	dopts := utils.GetOneSideTlsClientOpts()
	dopts = append(dopts, grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`)) // 设置负载均衡策略
	gwmux := runtime.NewServeMux()
	// grpc网关通过指定的endpoint参数连接到grpcServer， 如果是命名解析方式，就可以实现http转为grpc请求后，负载均衡选择一个grpcServer来调用服务
	//if err := pb.RegisterScoreServiceHandlerFromEndpoint(ctx, gwmux, addr, dopts); err != nil {
	if err := pb.RegisterScoreServiceHandlerFromEndpoint(ctx, gwmux, endpoint, dopts); err != nil { // 转为grpc请求后用命名解析的方式选择grpcServer
		grpclog.Fatalf("Failed to register gw server: %v\n", err)
	}
	// http服务
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)
	/**
	  处理http请求的时候他其实是分为  2段  的
	  通信的整个方式其实是：
	  	http请求发送到addr --> 网关 --> 转http请求为grpc请求 --> 命名解析endpoint根据负载均衡选择一个grpcServer --> C grpcServer
	  																							 B grpcServer
	  																							 A grpcServer
	*/
	srv := &http.Server{
		Addr:      addr, // 注意！！！这里的addr是http请求的地址，其实和net.listen的lis的端口是同一个，也就是http请求发送的地址
		Handler:   grpcHandlerFunc(grpcServer, mux),
		TLSConfig: getTLSConfig(),
	}
	grpclog.Infof("gRPC and https listen on: %s\n", addr)
	log.Println("-------注册网关成功，可提供http服务, ", addr)

	if err := srv.Serve(tls.NewListener(lis, srv.TLSConfig)); err != nil {
		grpclog.Fatal("ListenAndServe: ", err)
	}

	//if err := s.Serve(lis); err != nil {
	//	log.Fatalf("failed to serve: %v", err)
	//}
}

// 用于判断请求来源于Rpc客户端还是Restful api的请求，根据不同的请求注册不同的ServerHTTP服务
func grpcHandlerFunc(gs *grpc.Server, otherHandler http.Handler) http.Handler {
	if otherHandler == nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gs.ServeHTTP(w, r)
		})
	}

	//return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	// 根据请求头判断是否是grpc调用
	//	if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") { // 请求基于HTTP/2且请求的是grpc服务
	//		gs.ServeHTTP(w, r)
	//	} else { // 请求的是http服务
	//		otherHandler.ServeHTTP(w, r)
	//	}
	//})

	// 使用官方的h2库,用这个库，网关服务用证书和不用证书都可以，就是同时支持http2和http
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 拦截了所有 h2c 流量，然后根据不同的请求流量类型将其劫持并重定向到相应的 Hander 中去处理。
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			gs.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})

	/**
	  curl -k --insecure -d '{"userMobile":"18956457845"}' https://localhost:50051/user/loginOrRegister

	*/
}

func getTLSConfig() *tls.Config {
	cert, _ := ioutil.ReadFile("/home/wangbing/grpc-test/key/server.pem")
	key, _ := ioutil.ReadFile("/home/wangbing/grpc-test/key/server.key")
	var demoKeyPair *tls.Certificate
	pair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		grpclog.Fatalf("TLS KeyPair err: %v\n", err)
	}
	demoKeyPair = &pair
	return &tls.Config{
		Certificates: []tls.Certificate{*demoKeyPair},
		NextProtos:   []string{http2.NextProtoTLS}, // HTTP2 TLS支持
	}
}
