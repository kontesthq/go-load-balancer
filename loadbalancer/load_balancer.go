package loadbalancer

import (
	"github.com/kontesthq/go-load-balancer/server"
)

// LoadBalancer defines the interface for different load balancing strategies.
type LoadBalancer interface {
	ChooseServer(client Client) (server.Server, error) // Chooses an instance based on the algorithm.
	GetServiceName() string
}
