package load_balancer

import (
	"math/rand"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
)

type LoadBalancer struct {
	consulClient *api.Client
	mu           sync.RWMutex
	serviceName  string
}

// NewLoadBalancer creates a new LoadBalancer instance and returns it
func NewLoadBalancer(serviceName string) (*LoadBalancer, error) {
	consulClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err // Handle error if client creation fails
	}
	return &LoadBalancer{
		consulClient: consulClient,
		serviceName:  serviceName,
	}, nil
}

// GetHealthyInstances retrieves healthy instances of the specified service
func (lb *LoadBalancer) GetHealthyInstances() ([]*api.AgentService, error) {
	// Fetch healthy services directly using the Health package
	health, _, err := lb.consulClient.Health().Service(lb.serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}

	// Convert []*api.ServiceEntry to []*api.AgentService
	var instances []*api.AgentService
	for _, entry := range health {
		instances = append(instances, entry.Service) // Append the Service field
	}

	return instances, nil // Return healthy instances
}

// ChooseInstance randomly selects one of the healthy instances
func (lb *LoadBalancer) ChooseInstance() (*api.AgentService, error) {
	lb.mu.RLock() // Read lock for safe concurrent access
	defer lb.mu.RUnlock()

	instances, err := lb.GetHealthyInstances()
	if err != nil || len(instances) == 0 {
		return nil, err
	}

	// Randomly select an instance
	return instances[rand.Intn(len(instances))], nil
}

// init initializes the random number generator
func init() {
	rand.Seed(time.Now().UnixNano())
}
