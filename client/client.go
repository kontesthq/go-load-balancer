package client

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/kontesthq/go-load-balancer/loadbalancer"
	"strconv"
)

type Client struct {
	consulClient *api.Client
	serviceName  string
	loadbalancer loadbalancer.LoadBalancer
}

// LoadBalancerType defines the types of load balancers available
type LoadBalancerType int

// Constants for different load balancer types
const (
	RoundRobin LoadBalancerType = iota
	Random
)

func NewClient(consulHost string, consulPort int, serviceName string, lbType LoadBalancerType) (*Client, error) {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulHost + ":" + strconv.Itoa(consulPort)
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}

	// Initialize the load balancer based on the specified type
	var lb *loadbalancer.BaseLoadBalancer
	switch lbType {
	case RoundRobin:
		lb = loadbalancer.NewBaseLoadBalancerWithRule(consulClient, loadbalancer.NewRoundRobinRule(), serviceName)
	case Random:
		lb = loadbalancer.NewBaseLoadBalancerWithRule(consulClient, &loadbalancer.RandomRule{}, serviceName)
	default:
		return nil, fmt.Errorf("unsupported load balancer type: %v", lbType)
	}

	return &Client{
		consulClient: consulClient,
		serviceName:  serviceName,
		loadbalancer: lb,
	}, nil
}

func NewClientWithCustomRule(consulHost string, consulPort int, serviceName string, rule loadbalancer.IRule) (*Client, error) {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulHost + ":" + strconv.Itoa(consulPort)
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}

	// Initialize the load balancer based on the specified type
	lb := loadbalancer.NewBaseLoadBalancerWithRule(consulClient, rule, serviceName)

	return &Client{
		consulClient: consulClient,
		serviceName:  serviceName,
		loadbalancer: lb,
	}, nil
}

func (c *Client) GetLoadBalancer() loadbalancer.LoadBalancer {
	return c.loadbalancer
}
