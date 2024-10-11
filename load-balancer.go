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

func NewLoadBalancer(serviceName string) *LoadBalancer {
	consulClient, _ := api.NewClient(api.DefaultConfig())
	return &LoadBalancer{
		consulClient: consulClient,
		serviceName:  serviceName,
	}
}

func (lb *LoadBalancer) GetHealthyInstances() ([]*api.AgentService, error) {
	services, err := lb.consulClient.Agent().Services()
	if err != nil {
		return nil, err
	}

	var instances []*api.AgentService
	for _, service := range services {
		if service.Service == lb.serviceName {
			instances = append(instances, service)
		}
	}

	return instances, nil
}

func (lb *LoadBalancer) ChooseInstance() (*api.AgentService, error) {
	instances, err := lb.GetHealthyInstances()
	if err != nil || len(instances) == 0 {
		return nil, err
	}

	// Randomly select an instance
	return instances[rand.Intn(len(instances))], nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

//func main() {
//	fmt.Println("Hello, World!")
//}
