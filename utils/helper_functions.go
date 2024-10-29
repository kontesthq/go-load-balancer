package utils

import "github.com/hashicorp/consul/api"

// GetHealthyInstancesOfAService retrieves healthy instances of the specified service.
func GetHealthyInstancesOfAService(consulClient *api.Client, serviceName string) ([]*api.AgentService, error) {
	health, _, err := consulClient.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}

	var instances []*api.AgentService
	for _, entry := range health {
		instances = append(instances, entry.Service)
	}

	return instances, nil
}

// GetAllInstances retrieves all instances of all services registered with Consul.
func GetAllInstances(consulClient *api.Client) ([]*api.AgentService, error) {
	// Get the list of all services
	services, _, err := consulClient.Catalog().Services(nil)
	if err != nil {
		return nil, err
	}

	// Prepare to hold all AgentServices
	var allInstances []*api.AgentService

	// Iterate over all services
	for serviceName := range services {
		// Get service instances from the catalog
		instances, _, err := consulClient.Catalog().Service(serviceName, "", nil)
		if err != nil {
			return nil, err
		}

		// Convert CatalogService instances to AgentService instances
		for _, instance := range instances {
			// Create a new AgentService
			agentService := &api.AgentService{
				ID:      instance.ID,
				Service: instance.ServiceID, // Assuming ServiceID is correct
				Address: instance.Address,
				Port:    instance.ServicePort, // Use ServicePort to match the structure
				Tags:    instance.ServiceTags, // Use ServiceTags to match the structure
			}
			allInstances = append(allInstances, agentService)
		}
	}

	return allInstances, nil
}
