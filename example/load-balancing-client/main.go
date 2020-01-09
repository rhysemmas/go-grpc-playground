package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"

	pb "../../proto"
)

const (
	exampleScheme      = "test"
	exampleServiceName = "test.example.com"
)

var (
	addrs = []string{"localhost:50051", "localhost:50052"}
)

func makeRPCs(conn *grpc.ClientConn, n int) {
	// create new client with connection
	c := pb.NewPingClient(conn)
	// ping the server
	for {
		response, err := c.SayHello(context.Background(), &pb.PingMessage{Greeting: "foo"})
		log.Printf("Response from server: %s", response.Greeting)
		if err != nil {
			log.Fatalf("Error when calling SayHello: %s", err)
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	// Make a ClientConn with round_robin policy.
	roundrobinConn, err := grpc.Dial(
		// We set the address to connect to as the scheme/serviceName, which will be resolved by our example resolver below
		fmt.Sprintf("%s:///%s", exampleScheme, exampleServiceName),
		grpc.WithBalancerName(roundrobin.Name), // This sets the initial balancing policy.
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer roundrobinConn.Close()

	fmt.Println("--- calling helloworld.Greeter/SayHello with round_robin ---")
	makeRPCs(roundrobinConn, 10)
}

// Following is an example name resolver implementation based on: https://github.com/grpc/grpc-go/tree/master/examples/features/name_resolving
type exampleResolverBuilder struct{}

func (*exampleResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &exampleResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			exampleServiceName: addrs,
		},
	}
	r.start()
	return r, nil
}
func (*exampleResolverBuilder) Scheme() string { return exampleScheme }

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

func init() {
	resolver.Register(&exampleResolverBuilder{})
}
