package loadbalancer

// LoadBalancer defines the interface for different load balancing strategies.
type LoadBalancer interface {
	ChooseServer() (Server, error) // Chooses an instance based on the algorithm.
	GetServers() []Server
	GetServiceName() string
}
