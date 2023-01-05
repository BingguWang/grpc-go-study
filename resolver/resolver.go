package resolver

import (
	"google.golang.org/grpc/resolver"
)

/**
自定义一个命名解析器,
	没有注册中心的时候，客户端需要自己实现服务名-地址的映射，所以这里自定义了客户端的命名解析器，有了注册中心就不需要了
根据服务名解析出地址列表

到这为止， 还只是完成了客户端实现了 服务名--> 地址的解析而已，但实际我们不会把映射关系写死在代码里，而是会把这个解析的工作以及映射关系交给注册中心
我们还要实现注册中心，把服务名和地址的关系应该移到注册中心里去，让注册中心去解析，这才是最后的目的
*/

var (
	ExampleScheme      = "example"
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
