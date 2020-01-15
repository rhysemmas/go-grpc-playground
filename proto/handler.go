package proto

import (
	"log"

	"golang.org/x/net/context"
)

// Server represents the gRPC server
type Server struct {
	Addr string
}

// SayHello generates response to a Ping request
func (s *Server) SayHello(ctx context.Context, in *PingMessage) (*PingMessage, error) {
	log.Printf("Receive message %s", in.Greeting)
	return &PingMessage{Greeting: "bar", Server: s.Addr}, nil
}
