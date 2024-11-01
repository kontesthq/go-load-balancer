package loadbalancer

import (
	"fmt"
	"github.com/kontesthq/go-load-balancer/server"
	"log/slog"
	"sync/atomic"
)

type RoundRobinRule struct {
	nextServerCyclicCounter int32
}

func NewRoundRobinRule() *RoundRobinRule {
	return &RoundRobinRule{}
}

func (r *RoundRobinRule) ChooseServer(lb LoadBalancer) server.Server {
	if lb == nil {
		slog.Warn("LoadBalancer is nil")
		return nil
	}

	var server server.Server = nil
	var count int = 0

	for count < 10 {
		count++

		servers := lb.GetServers()

		if servers == nil || len(servers) == 0 {
			slog.Warn("No servers available")
			return nil
		}

		nextServerIndex := r.incrementAndGetModulo(int32(len(servers)))
		server := servers[nextServerIndex]

		if server == nil {
			/* Transient */
			continue
		}

		return server
	}

	slog.Warn(fmt.Sprintf("No available alive servers after 10 tries from load balancer: %v", lb))

	return server
}

func (r *RoundRobinRule) incrementAndGetModulo(modulo int32) int32 {
	for {
		current := atomic.LoadInt32(&r.nextServerCyclicCounter)
		next := (current + 1) % modulo

		if atomic.CompareAndSwapInt32(&r.nextServerCyclicCounter, current, next) {
			return next
		}
	}
}
