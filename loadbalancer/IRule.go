package loadbalancer

import (
	"github.com/hashicorp/consul/api"
)

type IRule interface {
	ChooseServer(lb LoadBalancer) *api.AgentService
}
