package main

import (
	"fmt"
	"github.com/kontesthq/go-load-balancer/client"
	"github.com/kontesthq/go-load-balancer/loadbalancer"
	"log/slog"
)

func main() {
	loadBalancerClient, err := client.NewClientWithCustomRule("localhost", 5150, "KONTEST-API", loadbalancer.NewRetryRuleWithDefaults())

	if err != nil {
		panic(err)
	}

	test(loadBalancerClient)

	slog.Info("Completed!")
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

func printServer(server loadbalancer.Server) {
	//message := fmt.Sprintf("Kind: %s, ID: %s, Address: %s, Service: %s", server.Kind, server.ID, server.Address, server.Service)
	message := fmt.Sprintf("Server: %v\n", loadbalancer.CommonServerString(server))
	slog.Info(message)
}
