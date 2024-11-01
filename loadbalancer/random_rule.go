package loadbalancer

import (
	"github.com/kontesthq/go-load-balancer/server"
	"math/rand"
)

type RandomRule struct {
}

func (r *RandomRule) ChooseServer(lb LoadBalancer) server.Server {
	if lb == nil {
		return nil
	}

	var server server.Server = nil

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
