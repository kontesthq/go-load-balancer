package loadbalancer

import (
	"github.com/kontesthq/go-load-balancer/server"
)

type Client interface {
	GetLoadBalancer() LoadBalancer
	GetHealthyInstances() ([]server.Server, error)
	GetAllInstances() ([]server.Server, error)
}
