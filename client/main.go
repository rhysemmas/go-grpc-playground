package main

import (
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "../proto"
)

func main() {
	var conn *grpc.ClientConn

	// dial server
	conn, err := grpc.Dial(":7777", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	// defer closing the connection until function returns
	defer conn.Close()

	// create new client with connection
	c := pb.NewPingClient(conn)

	// ping the server
	for {
		response, err := c.SayHello(context.Background(), &pb.PingMessage{Greeting: "foo"})
		log.Printf("Response from server: %s", response.Greeting)
		if err != nil {
			log.Fatalf("Error when calling SayHello: %s", err)
			promMetric()
		}
		time.Sleep(1 * time.Second)
	}
}

// podlink_error_count
func promMetric() {}
