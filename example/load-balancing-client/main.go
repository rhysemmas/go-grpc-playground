package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"

	myrrbalancer "github.com/rhysemmas/go-grpc-playground/example/load-balancing-client/balancer"
	pb "github.com/rhysemmas/go-grpc-playground/proto"
)

const (
	exampleScheme      = "test"
	exampleServiceName = "test.example.com"
)

var (
	addrs = []string{"192.168.112.176:50051", "192.168.112.176:50052"}
)

func makeRPCs(conn *grpc.ClientConn, n int) {
	// create new client with connection
	c := pb.NewPingClient(conn)
	log.Printf("test1: %v", conn.Target())
	// ping the server
	for i := 0; i < n; i++ {
		log.Printf("test2: %v", conn.GetState())
		response, err := c.SayHello(context.Background(), &pb.PingMessage{Greeting: "foo"})
		log.Printf("Response from '%s': %s", response.Server, response.Greeting)
		if err != nil {
			log.Printf("Error when calling SayHello: %s", err)
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	// start prom server
	myrrbalancer.PromHandler()

	// Make a ClientConn with round_robin policy.
	roundrobinConn, err := grpc.Dial(
		// We set the address to connect to as the scheme/serviceName, which will be resolved by our example resolver below
		fmt.Sprintf("%s:///%s", exampleScheme, exampleServiceName),
		grpc.WithBalancerName(myrrbalancer.Name), // This sets the initial balancing policy.
		grpc.WithInsecure(),
		grpc.WithBlock(),
		// we are using a basic Unary RPC
		grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//myrrbalancer.State()
	defer roundrobinConn.Close()

	fmt.Println("--- calling podlink.Ping/SayHello with round_robin ---")
	makeRPCs(roundrobinConn, 1000)
}

// The below will need replacing with some Kube/DNS resolving stuff
//
// Following is an example name resolver implementation based on: https://github.com/grpc/grpc-go/tree/master/examples/features/name_resolving
//
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
	state      resolver.State
	addrsStore map[string][]string
}

func (r *exampleResolver) start() {
	addrStrs := r.addrsStore[r.target.Endpoint]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	log.Printf("Addrs: %v", addrs)
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}
func (*exampleResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*exampleResolver) Close()                                  {}

func init() {
	resolver.Register(&exampleResolverBuilder{})
}
