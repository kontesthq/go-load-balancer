package client

import (
	"github.com/kontesthq/go-load-balancer/loadbalancer"
)

type Client interface {
	GetLoadBalancer() loadbalancer.LoadBalancer
}
