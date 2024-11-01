package loadbalancer

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/kontesthq/go-load-balancer/server"
	"log/slog"
	"strconv"
	"time"
)

var (
	cacheDuration   = 30 * time.Second
	lastUpdatedTime time.Time
)

type ConsulClient struct {
	consulClient           *api.Client
	cachedHealthyInstances []server.Server
	serviceName            string
	loadbalancer           LoadBalancer
}

func (c *ConsulClient) GetAllInstances() ([]server.Server, error) {
	return c.GetHealthyInstances()
}

func (c *ConsulClient) GetHealthyInstances() ([]server.Server, error) {
	// Check if the cache is valid
	if time.Since(lastUpdatedTime) <= cacheDuration {
		slog.Info("Using cached healthy instances")
		return c.cachedHealthyInstances, nil // Return cached healthy instances if valid
	}
	slog.Info("Fetching healthy instances from Consul")

	// Fetch healthy instances from Consul
	health, _, err := c.consulClient.Health().Service(c.serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}

	// Clear previous cached healthy instances
	c.cachedHealthyInstances = nil

	// Collect healthy instances
	for _, entry := range health {
		c.cachedHealthyInstances = append(c.cachedHealthyInstances, server.NewConsulServer(entry.Service, true, true)) // Append to cached healthy instances
	}

	// Update the last updated time
	lastUpdatedTime = time.Now()

	return c.cachedHealthyInstances, nil
}

func NewConsulClient(consulHost string, consulPort int, serviceName string, lbType LoadBalancerType) (*ConsulClient, error) {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulHost + ":" + strconv.Itoa(consulPort)
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}

	// Initialize the load balancer based on the specified type
	var lb *BaseLoadBalancer
	switch lbType {
	case RoundRobin:
		lb = NewBaseLoadBalancerWithRule(NewRoundRobinRule(), serviceName)
	case Random:
		lb = NewBaseLoadBalancerWithRule(&RandomRule{}, serviceName)
	default:
		return nil, fmt.Errorf("unsupported load balancer type: %v", lbType)
	}

	return &ConsulClient{
		consulClient: consulClient,
		serviceName:  serviceName,
		loadbalancer: lb,
	}, nil
}

func NewConsulClientWithCustomRule(consulHost string, consulPort int, serviceName string, rule IRule) (*ConsulClient, error) {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulHost + ":" + strconv.Itoa(consulPort)
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}

	// Initialize the load balancer based on the specified type
	lb := NewBaseLoadBalancerWithRule(rule, serviceName)

	return &ConsulClient{
		consulClient: consulClient,
		serviceName:  serviceName,
		loadbalancer: lb,
	}, nil
}

func (c *ConsulClient) GetLoadBalancer() LoadBalancer {
	return c.loadbalancer
}
