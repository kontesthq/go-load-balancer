package main

import (
	"fmt"
	"github.com/kontesthq/go-load-balancer/loadbalancer"
	"github.com/kontesthq/go-load-balancer/server"
	"log/slog"
	"time"
)

func main() {
	loadBalancerClient, err := loadbalancer.NewConsulClient("localhost", 5150, "KONTEST-API", loadbalancer.RoundRobin)

	if err != nil {
		panic(err)
	}

	test(loadBalancerClient)
	test(loadBalancerClient)

	time.Sleep(10 * time.Second)
	test(loadBalancerClient)

	time.Sleep(10 * time.Second)
	test(loadBalancerClient)

	time.Sleep(10 * time.Second)
	test(loadBalancerClient)

	time.Sleep(10 * time.Second)
	test(loadBalancerClient)

}

func test(client *loadbalancer.ConsulClient) {
	server, err := (*client).GetLoadBalancer().ChooseServer(client)
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
