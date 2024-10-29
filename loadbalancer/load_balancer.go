package loadbalancer

import (
	"github.com/hashicorp/consul/api"
)

// LoadBalancer defines the interface for different load balancing strategies.
type LoadBalancer interface {
	ChooseServer() (*api.AgentService, error) // Chooses an instance based on the algorithm.
	GetServers() []*api.AgentService
	GetServiceName() string
}
