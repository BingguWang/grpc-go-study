package utils

import (
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"io/ioutil"
	"log"
)

/**
封装一下TLS校验的代码
*/

// GetOneSideTlsServerOpts 单向TLS认证，服务端
func GetOneSideTlsServerOpts() (opts []grpc.ServerOption) {
	creds, err := credentials.NewServerTLSFromFile( // 单向TLS认证
		"/home/wangbing/grpc-test/key/server.pem",
		"/home/wangbing/grpc-test/key/server.key",
	)
	if err != nil {
		grpclog.Fatalf("Failed to load TLS credentials %v", err)
	}
	opts = append(opts, grpc.Creds(creds)) // 传入上面创建的启动TLS的证书，为所有传入的连接启用 TLS
	return opts
}

// GetOneSideTlsClientOpts 单向TLS认证，客户端
func GetOneSideTlsClientOpts() (opts []grpc.DialOption) {
	creds, err := credentials.NewClientTLSFromFile( // 单向TLS认证
		"/home/wangbing/grpc-test/key/server.pem",
		"x.binggu.example.com", // 证书里的commonName
	)
	if err != nil {
		grpclog.Fatalf("Failed to load TLS credentials %v", err)
	}
	opts = append(opts, grpc.WithTransportCredentials(creds))
	return opts
}

// GetBothSideTlsServerOpts 双向TLS认证，服务端
func GetBothSideTlsServerOpts() (opts []grpc.ServerOption) {
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
	// ca公钥加入池
	if ok := certPool.AppendCertsFromPEM(ca); !ok { // 尝试解析所传入的 PEM 编码的证书。如果解析成功会将其加到公钥池中
		log.Fatalf("certPool.AppendCertsFromPEM err")
	}
	// 创建credentials对象
	c := credentials.NewTLS(&tls.Config{ // 创建TSL连接
		Certificates: []tls.Certificate{cert},        // 服务端证书链，允许多个
		ClientAuth:   tls.RequireAndVerifyClientCert, // 需要并且验证客户端证书
		ClientCAs:    certPool,                       // 设置根证书的集合，校验方式使用 ClientAuth 中设定的模式
	})

	opts = append(opts, grpc.Creds(c))
	return opts
}

// GetBothSideTlsClientOpts 双向TLS认证，客户端
func GetBothSideTlsClientOpts() (opts []grpc.DialOption) {
	cert, err := tls.LoadX509KeyPair(
		"/home/wangbing/grpc-test/ce-client/client.pem",
		"/home/wangbing/grpc-test/ce-client/client.key",
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
	// ca公钥加入池
	if ok := certPool.AppendCertsFromPEM(ca); !ok { // 尝试解析所传入的 PEM 编码的证书。如果解析成功会将其加到公钥池中
		log.Fatalf("certPool.AppendCertsFromPEM err")
	}
	// 创建credentials对象
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert}, // 客户端证书
		ServerName:   "x.binggu.example.com",  // 证书里的commonName
		RootCAs:      certPool,
	})

	opts = append(opts, grpc.WithTransportCredentials(creds))
	return opts
}
