package loadbalancer

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/kontesthq/go-load-balancer/server"
	"log/slog"
	"strconv"
	"time"
)

type ConsulClient struct {
	consulClient           *api.Client
	cachedHealthyInstances []server.Server
	cacheDuration          time.Duration
	lastUpdatedTime        time.Time
	serviceName            string
	loadBalancer           LoadBalancer
}

func (c *ConsulClient) GetAllInstances() ([]server.Server, error) {
	return c.GetHealthyInstances()
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
		loadBalancer: lb,
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
		loadBalancer: lb,
	}, nil
}

func (c *ConsulClient) GetHealthyInstances() ([]server.Server, error) {
	// Check if the cache is valid
	if time.Since(c.lastUpdatedTime) <= c.cacheDuration {
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
	c.lastUpdatedTime = time.Now()

	return c.cachedHealthyInstances, nil
}

func (c *ConsulClient) GetLoadBalancer() LoadBalancer {
	return c.loadBalancer
}

func (c *ConsulClient) SetCacheDuration(duration time.Duration) {
	c.cacheDuration = duration
}
