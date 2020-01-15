package myrrbalancer

import (
	"log"
	"math/rand"
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/grpclog"
)

// Name is the name of round_robin balancer.
const Name = "my_round_robin"

// newBuilder creates a new roundrobin balancer builder.
func newBuilder() balancer.Builder {
	return base.NewBalancerBuilderV2(Name, &rrPickerBuilder{}, base.Config{HealthCheck: true})
}

func init() {
	balancer.Register(newBuilder())
}

type rrPickerBuilder struct{}

func (*rrPickerBuilder) Build(info base.PickerBuildInfo) balancer.V2Picker {
	grpclog.Infof("roundrobinPicker: newPicker called with info: %v", info)
	if len(info.ReadySCs) == 0 {
		return base.NewErrPickerV2(balancer.ErrNoSubConnAvailable)
	}
	var scs []balancer.SubConn
	for sc := range info.ReadySCs {
		scs = append(scs, sc)
		log.Printf("ReadySCs: %v", scs)
	}
	log.Printf("SCs info: %v", info)
	return &rrPicker{
		SubConns: scs,
		// Start at a random index, as the same RR balancer rebuilds a new
		// picker when SubConn states change, and we don't want to apply excess
		// load to the first server in the list.
		next: rand.Intn(len(scs)),
	}
}

type rrPicker struct {
	// subConns is the snapshot of the roundrobin balancer when this picker was
	// created. The slice is immutable. Each Get() will do a round robin
	// selection from it and return the selected SubConn.

	//subConns []balancer.SubConn
	SubConns []balancer.SubConn

	mu   sync.Mutex
	next int
}

func (p *rrPicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	p.mu.Lock()
	sc := p.SubConns[p.next]
	p.next = (p.next + 1) % len(p.SubConns)
	p.mu.Unlock()
	return balancer.PickResult{SubConn: sc}, nil
}

// // ScStates exported
// type scStates struct {
// 	sc    []balancer.SubConn
// 	state balancer.SubConnState
// }
//
// // State func exported
// func (sc *scStates) State() {
// 	scs := sc.sc
// 	s := sc.state
//
// 	for i := range scs {
// 		log.Printf("SC state: %v, %v", i, s)
// 	}
// }
