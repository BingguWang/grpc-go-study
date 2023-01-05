package etcdv3

import (
	"context"
	"fmt"
	mr "github.com/BingguWang/grpc-go-study/resolver"
	"net"
	"strings"
	"time"

	//"github.com/coreos/etcd/clientv3" // 不要用这个库，这个只能用1.26.0前版本的grpc
	clientv3 "go.etcd.io/etcd/client/v3"
)

/**
自定义封装的etcd服务注册的方法

注册应该是由服务端去调用的
*/

// Prefix should start and end with no slash
var Deregister = make(chan struct{})

// Register 注册服务到etcd, etcd，没有专门注册服务的函数，需要自己封装, 注册工作由服务端来完成
func Register(target, service, host, port string, interval time.Duration, ttl int) error {
	serviceValue := net.JoinHostPort(host, port)
	serviceKey := fmt.Sprintf("/%s/%s/%s", mr.ExampleScheme, service, serviceValue)

	// get endpoints for register dial address
	var err error
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: strings.Split(target, ","),
	})
	if err != nil {
		return fmt.Errorf("grpclb: create clientv3 client failed: %v", err)
	}
	// 设置租约
	resp, err := cli.Grant(context.TODO(), int64(ttl))
	if err != nil {
		return fmt.Errorf("grpclb: create clientv3 lease failed: %v", err)
	}

	// put key
	if _, err := cli.Put(context.TODO(), serviceKey, serviceValue, clientv3.WithLease(resp.ID)); err != nil {
		return fmt.Errorf("grpclb: set service '%s' with ttl to clientv3 failed: %s", service, err.Error())
	}

	if _, err := cli.KeepAlive(context.TODO(), resp.ID); err != nil {
		return fmt.Errorf("grpclb: refresh service '%s' with ttl to clientv3 failed: %s", service, err.Error())
	}

	// wait deregister then delete
	go func() {
		<-Deregister
		cli.Delete(context.Background(), serviceKey)
		Deregister <- struct{}{}
	}()

	return nil
}

// UnRegister delete registered service from etcd
func UnRegister() {
	Deregister <- struct{}{}
	<-Deregister
}
