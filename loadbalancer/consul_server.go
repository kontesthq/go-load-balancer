package loadbalancer

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"log/slog"
)

type ConsulServer struct {
	service *api.AgentService
	alive   bool
	ready   bool
}

func NewConsulServer(service *api.AgentService, alive, ready bool) *ConsulServer {
	return &ConsulServer{
		service: service,
		alive:   alive,
		ready:   ready,
	}
}

func (c *ConsulServer) GetID() string {
	return c.service.ID
}

func (c *ConsulServer) GetHost() string {
	return c.service.Address
}

func (c *ConsulServer) GetPort() int {
	return c.service.Port
}

func (c *ConsulServer) GetScheme() string {
	for _, tag := range c.service.Tags {
		if tag == "https" {
			return "https"
		}
	}
	return "http" // Default scheme if no https tag is found
}

func (c *ConsulServer) IsAlive() bool {
	return c.alive
}

func (c *ConsulServer) SetAlive(isAlive bool) {
	c.alive = isAlive
}

func (c *ConsulServer) GetHostPort() string {
	return fmt.Sprintf("%s:%d", c.GetHost(), c.GetPort())
}

func (c *ConsulServer) GetMetaInfo() MetaInfo {
	return &simpleMetaInfoImpl{
		appName:    c.service.Service,
		instanceId: c.service.ID,
	}
}

func (c *ConsulServer) GetZone() string {
	if c.service.Locality == nil {
		return UNKNOWN_ZONE
	}
	return c.service.Locality.Zone
}

func (c *ConsulServer) SetZone(zone string) {
	if c.service.Locality == nil {
		slog.Warn("Locality is nil")
	}
	c.service.Locality.Zone = zone
}

func (c *ConsulServer) IsReadyToServe() bool {
	return c.ready
}

func (c *ConsulServer) SetReadyToServe(ready bool) {
	c.ready = ready
}
