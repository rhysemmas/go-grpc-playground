package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "../proto"
)

// prom metric
var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "podlink_error_count_total",
		Help: "The total number of processed failure events",
	})
)

//prom handler
func promHandler() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		servErr := http.ListenAndServe(":7778", nil)
		if servErr != nil {
			log.Printf("metrics server failure")
		}
	}()
}

func main() {
	var conn *grpc.ClientConn

	// start prom server
	promHandler()

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
		log.Printf("Response from '%s': %s", response.Server, response.Greeting)
		if err != nil {
			log.Printf("Error when calling SayHello: %s", err)
			opsProcessed.Inc()
		}
		time.Sleep(1 * time.Second)
	}
}
