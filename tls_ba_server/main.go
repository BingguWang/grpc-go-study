package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	pb "github.com/BingguWang/grpc-go-study/server/proto"
	"github.com/BingguWang/grpc-go-study/server/service"
	"github.com/BingguWang/grpc-go-study/server/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
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

	cert, err := tls.LoadX509KeyPair(
		"/home/wangbing/grpc-test/ce-server/server.pem",
		"/home/wangbing/grpc-test/ce-server/server.key",
	) // 从证书相关文件中读取和解析信息，得到证书公钥-密钥对
	if err != nil {
		log.Fatalf("tls.LoadX509KeyPair err: %v", err)
	}
	// 建立公钥池
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("/home/wangbing/grpc-test/ca.pem")
	if err != nil {
		log.Fatalf("ioutil.ReadFile err: %v", err)
	}
	// ca公钥入池
	if ok := certPool.AppendCertsFromPEM(ca); !ok { // 尝试解析所传入的 PEM 编码的证书。如果解析成功会将其加到 CertPool 中
		log.Fatalf("certPool.AppendCertsFromPEM err")
	}
	c := credentials.NewTLS(&tls.Config{ // 创建TSL连接
		Certificates: []tls.Certificate{cert}, // 证书链，允许多个
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool, // 设置根证书的集合，校验方式使用 ClientAuth 中设定的模式
	})

	server := grpc.NewServer(grpc.Creds(c)) // 传入服务器
	pb.RegisterScoreServiceServer(server, service.GetServer())
	log.Printf("server listening at %v", listen.Addr())

	// 输出注册完的serviceInfo看下
	fmt.Println(utils.ToJsonString(server.GetServiceInfo()))

	if err := server.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
