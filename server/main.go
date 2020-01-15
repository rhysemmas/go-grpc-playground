package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "../proto"
)

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

// main start a gRPC server and waits for connection
func main() {
	// create a listener on TCP port 7777
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 7777))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a server instance
	addr := getLocalIP()
	s := &pb.Server{Addr: addr}

	// create a gRPC server object
	grpcServer := grpc.NewServer()

	// attach the Ping service to the server
	pb.RegisterPingServer(grpcServer, s)

	// start the server
	log.Printf("starting to serve...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
