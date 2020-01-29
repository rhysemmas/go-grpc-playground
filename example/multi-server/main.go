package main

import (
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

	pb "github.com/rhysemmas/go-grpc-playground/proto"
)

var (
	//addrs = []string{"50051", "50052"}
	addrs = []string{"50051"}
)

// main start a gRPC server and waits for connection
func createServer(addr string) {
	// create a listener on TCP port 7777
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", addr))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a server instance
	localIP := getLocalIP()
	addr = localIP + ":" + addr
	s := &pb.Server{Addr: addr}

	// create a gRPC server object
	grpcServer := grpc.NewServer()

	// attach the Ping service to the server
	pb.RegisterPingServer(grpcServer, s)

	// start the server
	log.Printf("starting to serve on %s", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func main() {
	// define a waitgroup
	var wg sync.WaitGroup

	// create a goroutine for each port
	for _, addr := range addrs {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			createServer(addr)
		}(addr)
	}

	wg.Wait()
}
