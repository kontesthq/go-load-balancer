package go_load_balancer

import (
	"github.com/hashicorp/consul/api"
	"strconv"
	"sync"
)

import (
	error2 "github.com/ayushs-2k4/go-load-balancer/error"
	"math/rand"
	"time"
)

// LoadBalancer manages load balancing for a service.
type LoadBalancer struct {
	consulClient *api.Client
	mu           sync.RWMutex
	serviceName  string
}

// Load Balancer Registry to manage multiple load balancer instances.
var (
	loadBalancers = make(map[string]*LoadBalancer)
	mu            sync.Mutex
)

// newLoadBalancer creates a new LoadBalancer instance and returns it.
// This function is private to the package.
func newLoadBalancer(serviceName string, consulHost string, consulPort int) (*LoadBalancer, error) {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulHost + ":" + strconv.Itoa(consulPort)
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, err // Handle error if client creation fails
	}

	return &LoadBalancer{
		consulClient: consulClient,
		serviceName:  serviceName,
	}, nil
}

// GetLoadBalancer retrieves or creates a LoadBalancer instance for a service.
func GetLoadBalancer(serviceName, consulHost string, consulPort int) (*LoadBalancer, error) {
	mu.Lock()
	defer mu.Unlock()

	if lb, exists := loadBalancers[serviceName]; exists {
		return lb, nil
	}

	lb, err := newLoadBalancer(serviceName, consulHost, consulPort)
	if err != nil {
		return nil, err
	}

	loadBalancers[serviceName] = lb
	return lb, nil
}

// GetHealthyInstances retrieves healthy instances of the specified service.
func (lb *LoadBalancer) GetHealthyInstances() ([]*api.AgentService, error) {
	health, _, err := lb.consulClient.Health().Service(lb.serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}

	var instances []*api.AgentService
	for _, entry := range health {
		instances = append(instances, entry.Service)
	}

	return instances, nil
}

// ChooseInstance randomly selects one of the healthy instances.
func (lb *LoadBalancer) ChooseInstance() (*api.AgentService, error) {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	instances, err := lb.GetHealthyInstances()
	if err != nil {
		return nil, err
	}

	if len(instances) == 0 {
		return nil, &error2.NoHealthyInstanceAvailableError{ServiceName: lb.serviceName}
	}

	return instances[rand.Intn(len(instances))], nil
}

// Initialize the random number generator.
func init() {
	rand.Seed(time.Now().UnixNano())
}
