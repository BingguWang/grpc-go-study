package resolver

import (
	"google.golang.org/grpc/resolver"
)

/**
自定义一个命名解析器
根据服务名解析出地址列表

*/

var (
	ExampleScheme = "example"
	//exampleServiceName = "resolver.example.grpc.io"
	ExampleServiceName = "x.binggu.example.com"

	backendAddr  = "localhost:50051"
	backendAddr2 = "localhost:50052"
)

type ExampleResolverBuilder struct{}

func (*ExampleResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &exampleResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			// 这里同一个服务名配置两个地址, 客户端会根据服务名exampleServiceName解析出这个地址列表，然后根据负载均衡策略来选择用哪个地址!!
			// 这里使用的是rr策略，所以会看到是轮询的选择解析出的地址
			ExampleServiceName: {backendAddr, backendAddr2},
		},
	}
	r.start()
	return r, nil
}
func (*ExampleResolverBuilder) Scheme() string { return ExampleScheme }

// exampleResolver is a
// Resolver(https://godoc.org/google.golang.org/grpc/resolver#Resolver).
type exampleResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *exampleResolver) start() {
	addrStrs := r.addrsStore[r.target.Endpoint]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}
func (*exampleResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*exampleResolver) Close()                                  {}
