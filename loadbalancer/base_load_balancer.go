package loadbalancer

import (
	"errors"
	"github.com/hashicorp/consul/api"
	"github.com/kontesthq/go-load-balancer/utils"
	"sync"
)

// BaseLoadBalancer manages load balancing for a service.
type BaseLoadBalancer struct {
	mu           sync.RWMutex
	rule         IRule
	consulClient *api.Client
	serviceName  string
}

func (lb *BaseLoadBalancer) GetServiceName() string {
	return lb.serviceName
}

func NewBaseLoadBalancer(consulClient *api.Client, serviceName string) *BaseLoadBalancer {
	return &BaseLoadBalancer{
		rule:         &RandomRule{},
		consulClient: consulClient,
		serviceName:  serviceName,
	}
}

func NewBaseLoadBalancerWithRule(consulClient *api.Client, rule IRule, serviceName string) *BaseLoadBalancer {
	return &BaseLoadBalancer{
		rule:         rule,
		consulClient: consulClient,
		serviceName:  serviceName,
	}
}

func (lb *BaseLoadBalancer) GetServers() []*api.AgentService {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	instances, err := utils.GetHealthyInstancesOfAService(lb.consulClient, lb.GetServiceName())
	if err != nil {
		return nil
	}

	return instances
}

func (lb *BaseLoadBalancer) ChooseServer() (*api.AgentService, error) {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	server := lb.rule.ChooseServer(lb)
	if server == nil {
		return nil, errors.New("no healthy server instance available")
	}

	return server, nil
}
