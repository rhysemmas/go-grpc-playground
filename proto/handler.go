package proto

import (
	"fmt"
	"log"

	"golang.org/x/net/context"
)

// Server represents the gRPC server
type Server struct {
	Addr int
}

// SayHello generates response to a Ping request
func (s *Server) SayHello(ctx context.Context, in *PingMessage) (*PingMessage, error) {
	log.Printf("Receive message %s", in.Greeting)
	addrString := fmt.Sprint(s.Addr)
	return &PingMessage{Greeting: "bar", Server: addrString}, nil
}
