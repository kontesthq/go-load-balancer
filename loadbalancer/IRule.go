package loadbalancer

import (
	"github.com/kontesthq/go-load-balancer/server"
)

type IRule interface {
	ChooseServer(client Client) server.Server
}
