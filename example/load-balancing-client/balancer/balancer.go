package balancer

import (
	"log"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/resolver"
)

// Name is the name of round_robin balancer.
const Name = "my_round_robin"

type ShimBalancer struct {
	bal balancer.Balancer
}

func (sb *ShimBalancer) HandleSubConnStateChange(sc balancer.SubConn, state connectivity.State) {
	sb.bal.HandleSubConnStateChange(sc, state)
}

func (sb *ShimBalancer) HandleResolvedAddrs(addresses []resolver.Address, err error) {
	sb.bal.HandleResolvedAddrs(addresses, err)
}

func (sb *ShimBalancer) Close() {
	sb.bal.Close()
}

func (sb *ShimBalancer) UpdateClientConnState(state balancer.ClientConnState) error {
	v2bal := sb.bal.(balancer.V2Balancer)
	return v2bal.UpdateClientConnState(state)
}

func (sb *ShimBalancer) ResolverError(err error) {
	v2bal := sb.bal.(balancer.V2Balancer)
	v2bal.ResolverError(err)
}

func (sb *ShimBalancer) UpdateSubConnState(sc balancer.SubConn, state balancer.SubConnState) {
	v2bal := sb.bal.(balancer.V2Balancer)
	log.Println("wibble")
	v2bal.UpdateSubConnState(sc, state)
}

type ShimBuilder struct {
	builder balancer.Builder
}

func (b *ShimBuilder) Build(cc balancer.ClientConn, opts balancer.BuildOptions) balancer.Balancer {
	bal := b.builder.Build(cc, opts)

	return &ShimBalancer{bal}
}

func (b *ShimBuilder) Name() string {
	return Name
}

// newBuilder creates a new roundrobin balancer builder.
func newBuilder() balancer.Builder {
	return &ShimBuilder{balancer.Get("round_robin")}
}

func init() {
	balancer.Register(newBuilder())
}
