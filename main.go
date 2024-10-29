package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/kontesthq/go-load-balancer/client"
	"log/slog"
)

func main() {
	loadBalancerClient, err := client.NewClient("localhost", 5150, "KONTEST-API", client.RoundRobin)

	if err != nil {
		panic(err)
	}

	test(loadBalancerClient)
	test(loadBalancerClient)
	test(loadBalancerClient)
	test(loadBalancerClient)
	test(loadBalancerClient)
	test(loadBalancerClient)
	test(loadBalancerClient)
	test(loadBalancerClient)
	test(loadBalancerClient)
	test(loadBalancerClient)

}

func test(client *client.Client) {
	server, err := (*client).GetLoadBalancer().ChooseServer()
	if err != nil {
		slog.Error("No Client available")
		return
	}

	if server == nil {
		slog.Error("No server found")
		return
	}

	printServer(server)
}

func printServer(server *api.AgentService) {
	message := fmt.Sprintf("Kind: %s, ID: %s, Address: %s, Service: %s",
		server.Kind, server.ID, server.Address, server.Service)
	slog.Info(message)
}
