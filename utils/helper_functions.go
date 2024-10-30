package utils

import (
	"github.com/hashicorp/consul/api"
	"log/slog"
	"time"
)

var (
	cacheDuration          time.Duration = 30 * time.Second
	cachedHealthyInstances []*api.AgentService
	cachedServices         []*api.CatalogService
	lastUpdatedTime        time.Time
)

// SetCacheDuration allows users to set a custom cache duration.
func SetCacheDuration(duration time.Duration) {
	cacheDuration = duration
}

// GetHealthyInstancesOfAService retrieves healthy instances of the specified service.
func GetHealthyInstancesOfAService(consulClient *api.Client, serviceName string) ([]*api.AgentService, error) {
	// Check if the cache is valid
	if time.Since(lastUpdatedTime) <= cacheDuration {
		slog.Info("Using cached healthy instances")
		return cachedHealthyInstances, nil // Return cached healthy instances if valid
	}
	slog.Info("Fetching healthy instances from Consul")

	// Fetch healthy instances from Consul
	health, _, err := consulClient.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}

	// Clear previous cached healthy instances
	cachedHealthyInstances = nil

	// Collect healthy instances
	for _, entry := range health {
		cachedHealthyInstances = append(cachedHealthyInstances, entry.Service) // Append to cached healthy instances
	}

	// Update the last updated time
	lastUpdatedTime = time.Now()

	return cachedHealthyInstances, nil
}

// GetAllInstances retrieves all instances of all services registered with Consul.
func GetAllInstances(consulClient *api.Client) ([]*api.CatalogService, error) {
	if time.Since(lastUpdatedTime) > cacheDuration {
		cachedServices = nil

		// Get the list of all services
		services, _, err := consulClient.Catalog().Services(nil)
		if err != nil {
			return nil, err
		}

		// Iterate over all services
		for serviceName := range services {
			// Get service instances from the catalog
			instances, _, err := consulClient.Catalog().Service(serviceName, "", nil)
			if err != nil {
				return nil, err
			}

			cachedServices = append(cachedServices, instances...)
		}

		lastUpdatedTime = time.Now()
	}

	return cachedServices, nil
}
