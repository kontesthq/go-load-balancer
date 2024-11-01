package loadbalancer

// LoadBalancerType defines the types of load balancers available
type LoadBalancerType int

// Constants for different load balancer types
const (
	RoundRobin LoadBalancerType = iota
	Random
)
