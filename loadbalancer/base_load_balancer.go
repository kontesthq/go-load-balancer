package loadbalancer

import (
	"errors"
	"github.com/kontesthq/go-load-balancer/server"
	"sync"
)

// BaseLoadBalancer manages load balancing for a service.
type BaseLoadBalancer struct {
	mu          sync.RWMutex
	rule        IRule
	serviceName string
}

func (lb *BaseLoadBalancer) GetServiceName() string {
	return lb.serviceName
}

func NewBaseLoadBalancer(serviceName string) *BaseLoadBalancer {
	return &BaseLoadBalancer{
		rule:        &RandomRule{},
		serviceName: serviceName,
	}
}

func NewBaseLoadBalancerWithRule(rule IRule, serviceName string) *BaseLoadBalancer {
	return &BaseLoadBalancer{
		rule:        rule,
		serviceName: serviceName,
	}
}

func (lb *BaseLoadBalancer) ChooseServer(client Client) (server.Server, error) {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	chosenServer := lb.rule.ChooseServer(client)
	if chosenServer == nil {
		return nil, errors.New("no healthy chosenServer instance available")
	}

	return chosenServer, nil
}
