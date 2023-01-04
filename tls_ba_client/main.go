package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/BingguWang/grpc-go-study/client/service"
	pb "github.com/BingguWang/grpc-go-study/server/proto"
	"github.com/BingguWang/grpc-go-study/server/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"time"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

/**
这里是TLS双向认证
	也就是客户端有了公钥和服务名后并不能随心所欲的调用服务，服务端对客户端也要进行筛选
*/
func main() {
	flag.Parse()
	fmt.Println(utils.ToJsonString(addr))
	cert, err := tls.LoadX509KeyPair(
		"/home/wangbing/grpc-test/ce-client/client.pem",
		"/home/wangbing/grpc-test/ce-client/client.key",
	)
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
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "x.binggu.example.com",
		RootCAs:      certPool,
	})

	conn, err := grpc.Dial(
		*addr,
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewScoreServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	service.CallStreamBidirectional(ctx, client)

}
