package loadbalancer

import (
	"github.com/hashicorp/consul/api"
	"math/rand"
)

type RandomRule struct {
}

func (r *RandomRule) ChooseServer(lb LoadBalancer) *api.AgentService {
	if lb == nil {
		return nil
	}

	var server *api.AgentService = nil

	for server == nil {
		servers := (lb).GetServers()

		if len(servers) == 0 {
			return nil
		}

		index := chooseRandomInt(len(servers))
		server = servers[index]

		if server == nil {
			/*
			 * The only time this should happen is if the server list were
			 * somehow trimmed. This is a transient condition. Retry after
			 * yielding.
			 */
			continue
		}
	}

	return server
}

func chooseRandomInt(serverCount int) int {
	return rand.Int() % serverCount
}
