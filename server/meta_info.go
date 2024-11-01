package server

// MetaInfo provides additional meta information about a server, including
// details about the application running on it and deployment-specific
// identification, such as server identifiers in environments like AWS.
type MetaInfo interface {

	// GetAppName returns the name of the application running on this server.
	// Returns an empty string if not available.
	GetAppName() string

	// GetServerGroup returns the group of the server, such as an auto-scaling group ID in AWS.
	// Returns an empty string if not available.
	GetServerGroup() string

	// GetServiceIdForDiscovery returns a virtual address used by the server to register with a discovery service.
	// Returns an empty string if not available.
	GetServiceIdForDiscovery() string

	// GetInstanceId returns the ID of the server instance.
	// Returns an empty string if not available.
	GetInstanceId() string
}

// SimpleMetaInfoImpl implements the MetaInfo interface with basic metadata.
type SimpleMetaInfoImpl struct {
	appName     string
	serverGroup string
	serviceId   string
	instanceId  string
}

// NewSimpleMetaInfo creates a new SimpleMetaInfoImpl with specified values.
func NewSimpleMetaInfo(appName, serverGroup, serviceId, instanceId string) *SimpleMetaInfoImpl {
	return &SimpleMetaInfoImpl{
		appName:     appName,
		serverGroup: serverGroup,
		serviceId:   serviceId,
		instanceId:  instanceId,
	}
}

// GetAppName returns the application name for the server.
func (m *SimpleMetaInfoImpl) GetAppName() string {
	return m.appName
}

// GetServerGroup returns the server group (e.g., AWS auto-scaling group ID).
func (m *SimpleMetaInfoImpl) GetServerGroup() string {
	return m.serverGroup
}

// GetServiceIdForDiscovery returns the service ID used for discovery.
func (m *SimpleMetaInfoImpl) GetServiceIdForDiscovery() string {
	return m.serviceId
}

// GetInstanceId returns the instance ID of the server.
func (m *SimpleMetaInfoImpl) GetInstanceId() string {
	return m.instanceId
}
