package etcdv3

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"

	//"github.com/coreos/etcd/clientv3" // 不要用这个库，这个只能用1.26.0前版本的grpc
	clientv3 "go.etcd.io/etcd/client/v3"

	"google.golang.org/grpc/resolver"
	"strings"
)

/*
 etcd作服务中心
	需要etcd实现命名解析的能力，所以etcd需要实现grpc.resolver
*/

var (
	Etcdv3Scheme = "example" // 可以看成是个前缀
)

// Resolver 自定义的resolver,实现了grpc.resolve.Builder
type Resolver struct {
	target  string // 注册中心etcd的地址
	service string
	cli     *clientv3.Client
	cc      resolver.ClientConn
}

// NewResolver 返回我们自定义的etcd解析器, 其实只是实现了grpc的resolver而已
// target example: "http://127.0.0.1:2379,http://127.0.0.1:12379,http://127.0.0.1:22379"
func NewResolver(target string, serviceName string) resolver.Builder {
	return &Resolver{target: target, service: serviceName}
}

// Scheme return etcdv3 schema
func (r *Resolver) Scheme() string {
	return Etcdv3Scheme
}
func (r *Resolver) ResolveNow(rn resolver.ResolveNowOptions) {
}
func (r *Resolver) Close() {
}

// Build 自定义etcd实现Build方法
func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var err error

	// 连接etcd
	r.cli, err = clientv3.New(clientv3.Config{
		Endpoints: strings.Split(r.target, ","),
	})
	if err != nil {
		return nil, fmt.Errorf("grpclb: create clientv3 client failed: %v", err)
	}

	r.cc = cc

	go r.watch(fmt.Sprintf("/%s/%s/", Etcdv3Scheme, r.service)) // 监听/example/score_service前缀的key

	return r, nil
}

func (r *Resolver) watch(prefix string) {
	addrDict := make(map[string]resolver.Address)

	update := func() {
		addrList := make([]resolver.Address, 0, len(addrDict))
		for _, v := range addrDict {
			addrList = append(addrList, v)
		}
		r.cc.UpdateState(resolver.State{Addresses: addrList})
	}

	resp, err := r.cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err == nil {
		for i := range resp.Kvs {
			addrDict[string(resp.Kvs[i].Value)] = resolver.Address{Addr: string(resp.Kvs[i].Value)}
		}
	}

	update()

	rch := r.cli.Watch(context.Background(), prefix, clientv3.WithPrefix(), clientv3.WithPrevKV())
	for n := range rch {
		for _, ev := range n.Events {
			switch ev.Type {
			case mvccpb.PUT:
				addrDict[string(ev.Kv.Key)] = resolver.Address{Addr: string(ev.Kv.Value)}
			case mvccpb.DELETE:
				delete(addrDict, string(ev.PrevKv.Key))
			}
		}
		update()
	}
}
