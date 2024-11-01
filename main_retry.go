package main

import (
	"fmt"
	"github.com/kontesthq/go-load-balancer/client"
	"github.com/kontesthq/go-load-balancer/loadbalancer"
	"github.com/kontesthq/go-load-balancer/server"
	"log/slog"
)

func main() {
	loadBalancerClient, err := client.NewConsulClientWithCustomRule("localhost", 5150, "KONTEST-API", loadbalancer.NewRetryRule(loadbalancer.NewRoundRobinRule(), 200))

	if err != nil {
		panic(err)
	}

	test(loadBalancerClient)

	slog.Info("Completed!")
}

func test(client *client.ConsulClient) {
	server, err := (*client).GetLoadBalancer().ChooseServer()
	if err != nil {
		slog.Error("No ConsulClient available")
		return
	}

	if server == nil {
		slog.Error("No server found")
		return
	}

	printServer(server)
}

func printServer(serverInstance server.Server) {
	//message := fmt.Sprintf("Kind: %s, ID: %s, Address: %s, Service: %s", server.Kind, server.ID, server.Address, server.Service)
	message := fmt.Sprintf("Server: %v\n", server.CommonServerString(serverInstance))
	slog.Info(message)
}
