package qezap

import (
	"fmt"

	"google.golang.org/grpc/resolver"
)

// 实现一个本地地址解析器，用于 grpc 负载, 如果对应服务下线，不会变更地址。
const (
	LocalServiceName = "grpc_local_resolver_name"
	LocalScheme      = "local"
)

var (
	DialLocalServiceName = fmt.Sprintf("%s:///%s", LocalScheme, LocalServiceName)
)

type LocalResolverBuilder struct {
	addrs []string
}

func NewLocalResolverBuilder(addrs []string) *LocalResolverBuilder {
	return &LocalResolverBuilder{addrs: addrs}
}

func (lrb *LocalResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &localResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			LocalServiceName: lrb.addrs,
		},
	}
	r.start()

	return r, nil
}
func (*LocalResolverBuilder) Scheme() string { return LocalScheme }

type localResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *localResolver) start() {
	addrStrs := r.addrsStore[r.target.Endpoint]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}
func (*localResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*localResolver) Close()                                  {}
