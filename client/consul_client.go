package client

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/kontesthq/go-load-balancer/loadbalancer"
	"strconv"
)

type ConsulClient struct {
	consulClient *api.Client
	serviceName  string
	loadbalancer loadbalancer.LoadBalancer
}

func NewConsulClient(consulHost string, consulPort int, serviceName string, lbType loadbalancer.LoadBalancerType) (*ConsulClient, error) {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulHost + ":" + strconv.Itoa(consulPort)
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}

	// Initialize the load balancer based on the specified type
	var lb *loadbalancer.BaseLoadBalancer
	switch lbType {
	case loadbalancer.RoundRobin:
		lb = loadbalancer.NewBaseLoadBalancerWithRule(consulClient, loadbalancer.NewRoundRobinRule(), serviceName)
	case loadbalancer.Random:
		lb = loadbalancer.NewBaseLoadBalancerWithRule(consulClient, &loadbalancer.RandomRule{}, serviceName)
	default:
		return nil, fmt.Errorf("unsupported load balancer type: %v", lbType)
	}

	return &ConsulClient{
		consulClient: consulClient,
		serviceName:  serviceName,
		loadbalancer: lb,
	}, nil
}

func NewConsulClientWithCustomRule(consulHost string, consulPort int, serviceName string, rule loadbalancer.IRule) (*ConsulClient, error) {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulHost + ":" + strconv.Itoa(consulPort)
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}

	// Initialize the load balancer based on the specified type
	lb := loadbalancer.NewBaseLoadBalancerWithRule(consulClient, rule, serviceName)

	return &ConsulClient{
		consulClient: consulClient,
		serviceName:  serviceName,
		loadbalancer: lb,
	}, nil
}

func (c *ConsulClient) GetLoadBalancer() loadbalancer.LoadBalancer {
	return c.loadbalancer
}
