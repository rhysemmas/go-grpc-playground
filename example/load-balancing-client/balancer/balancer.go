package balancer

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/resolver"
)

// Name is the name of round_robin balancer.
const Name = "my_round_robin"

var (
	// kube api to check for node IP/host?

	okStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "podlink_connection_status",
			Help: "Current temperature of the CPU.",
		},
		[]string{"pod_ip"},
	)

	conFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "podlink_connection_errors_total",
			Help: "Number of connection errors over time to pods.",
		},
		[]string{"pod_ip"},
	)
)

// ShimBalancer fuck off
type ShimBalancer struct {
	bal balancer.Balancer
	ip  string
}

// HandleSubConnStateChange fuck off
func (sb *ShimBalancer) HandleSubConnStateChange(sc balancer.SubConn, state connectivity.State) {
	sb.bal.HandleSubConnStateChange(sc, state)
}

// HandleResolvedAddrs fuck off
func (sb *ShimBalancer) HandleResolvedAddrs(addresses []resolver.Address, err error) {
	sb.bal.HandleResolvedAddrs(addresses, err)
}

// Close fuck off
func (sb *ShimBalancer) Close() {
	sb.bal.Close()
}

// UpdateClientConnState fuck off
func (sb *ShimBalancer) UpdateClientConnState(state balancer.ClientConnState) error {
	v2bal := sb.bal.(balancer.V2Balancer)
	return v2bal.UpdateClientConnState(state)
}

// ResolverError fuck off
func (sb *ShimBalancer) ResolverError(err error) {
	v2bal := sb.bal.(balancer.V2Balancer)
	v2bal.ResolverError(err)
}

// UpdateSubConnState fuck off
func (sb *ShimBalancer) UpdateSubConnState(sc balancer.SubConn, state balancer.SubConnState) {
	v2bal := sb.bal.(balancer.V2Balancer)

	if state.ConnectivityState == connectivity.TransientFailure {
		fmt.Println("woooooooooo", sc)
		r, _ := regexp.Compile("([0-9]+)(.)([0-9]+)(.)([0-9]+)(.)([0-9]+)(:?)([0-9]+)?")
		sb.ip = r.FindString(state.ConnectionError.Error())
		conFailures.With(prometheus.Labels{"pod_ip": sb.ip}).Inc()
		okStatus.With(prometheus.Labels{"pod_ip": sb.ip}).Set(1)
	}

	if state.ConnectivityState == connectivity.Ready {
		fmt.Println("fak")
		okStatus.With(prometheus.Labels{"pod_ip": sb.ip}).Set(0)
	}

	v2bal.UpdateSubConnState(sc, state)
}

// ShimBuilder fuck off
type ShimBuilder struct {
	builder balancer.Builder
}

// Build fuck off
func (b *ShimBuilder) Build(cc balancer.ClientConn, opts balancer.BuildOptions) balancer.Balancer {
	bal := b.builder.Build(cc, opts)

	return &ShimBalancer{bal: bal}
}

// Name fuck off
func (b *ShimBuilder) Name() string {
	return Name
}

// newBuilder creates a new roundrobin balancer builder.
func newBuilder() balancer.Builder {
	return &ShimBuilder{balancer.Get("round_robin")}
}

// PromHandler exported to be called by client making RPCs
func PromHandler() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		servErr := http.ListenAndServe(":50050", nil)
		if servErr != nil {
			log.Fatalf("metrics server failure")
		}
	}()
}

// fuck off
func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(conFailures)
	prometheus.MustRegister(okStatus)
	balancer.Register(newBuilder())
}
