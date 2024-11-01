package loadbalancer

import (
	"fmt"
	"github.com/kontesthq/go-load-balancer/server"
	"log/slog"
	"math/rand"
	"time"
)

type RandomRule struct {
}

func (r *RandomRule) ChooseServer(client Client) server.Server {
	if client == nil {
		return nil
	}

	var chosenServer server.Server = nil

	for chosenServer == nil {
		servers, err := client.GetHealthyInstances()

		if err != nil {
			slog.Error(fmt.Sprintf("Error in getting healthy instances: %v\n", err))
			return nil
		}

		if len(servers) == 0 {
			return nil
		}

		index := chooseRandomInt(len(servers))
		chosenServer = servers[index]

		if chosenServer == nil {
			/*
			 * The only time this should happen is if the chosenServer list were
			 * somehow trimmed. This is a transient condition. Retry after
			 * yielding.
			 */
			time.Sleep(1 * time.Millisecond)
			continue
		}
	}

	return chosenServer
}

func chooseRandomInt(serverCount int) int {
	return rand.Int() % serverCount
}
