package main

import (
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

	pb "../../proto"
)

var (
	addrs = []int{50051, 50052}
)

// main start a gRPC server and waits for connection
func createServer(addr int) {
	// create a listener on TCP port 7777
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", addr))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a server instance
	s := pb.Server{}

	// create a gRPC server object
	grpcServer := grpc.NewServer()

	// attach the Ping service to the server
	pb.RegisterPingServer(grpcServer, &s)

	// start the server
	log.Printf("starting to serve on %d", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func main() {
	// define a waitgroup
	var wg sync.WaitGroup

	// create a goroutine for each port
	for _, addr := range addrs {
		wg.Add(1)
		go func(addr int) {
			defer wg.Done()
			createServer(addr)
		}(addr)
	}

	wg.Wait()
}
