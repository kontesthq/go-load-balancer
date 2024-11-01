package loadbalancer

import (
	"github.com/kontesthq/go-load-balancer/server"
)

type IRule interface {
	ChooseServer(lb LoadBalancer) server.Server
}
