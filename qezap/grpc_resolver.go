package qezap

import (
	"fmt"

	"google.golang.org/grpc/resolver"
)

const (
	LocalServiceName = "grpc_local_resolver_name"
	LocalScheme      = "local"
)

var (
	DialLocalServiceName = fmt.Sprintf("%s:///%s", LocalScheme, LocalServiceName)
)

// LocalResolverBuilder impl local address resolver, used for GRPC load balancing
type LocalResolverBuilder struct {
	address []string
}

func NewLocalResolverBuilder(address []string) *LocalResolverBuilder {
	return &LocalResolverBuilder{address: address}
}

func (lrb *LocalResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &localResolver{
		target: target,
		cc:     cc,
		addressStore: map[string][]string{
			LocalServiceName: lrb.address,
		},
	}
	r.start()

	return r, nil
}
func (*LocalResolverBuilder) Scheme() string { return LocalScheme }

type localResolver struct {
	target       resolver.Target
	cc           resolver.ClientConn
	addressStore map[string][]string
}

func (r *localResolver) start() {
	addr := r.addressStore[r.target.Endpoint]
	address := make([]resolver.Address, len(addr))
	for i, s := range addr {
		address[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: address})
}
func (*localResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*localResolver) Close()                                  {}
