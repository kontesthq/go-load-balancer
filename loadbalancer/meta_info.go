package loadbalancer

type MetaInfo interface {
	GetAppName() string
	GetServerGroup() string
	GetServiceIdForDiscovery() string
	GetInstanceId() string
}

// simpleMetaInfoImpl implements the MetaInfo interface with basic metadata.
type simpleMetaInfoImpl struct {
	appName     string
	serverGroup string
	serviceId   string
	instanceId  string
}

// NewSimpleMetaInfo creates a new simpleMetaInfoImpl with specified values.
func NewSimpleMetaInfo(appName, serverGroup, serviceId, instanceId string) *simpleMetaInfoImpl {
	return &simpleMetaInfoImpl{
		appName:     appName,
		serverGroup: serverGroup,
		serviceId:   serviceId,
		instanceId:  instanceId,
	}
}

// GetAppName returns the application name for the server.
func (m *simpleMetaInfoImpl) GetAppName() string {
	return m.appName
}

// GetServerGroup returns the server group (e.g., AWS auto-scaling group ID).
func (m *simpleMetaInfoImpl) GetServerGroup() string {
	return m.serverGroup
}

// GetServiceIdForDiscovery returns the service ID used for discovery.
func (m *simpleMetaInfoImpl) GetServiceIdForDiscovery() string {
	return m.serviceId
}

// GetInstanceId returns the instance ID of the server.
func (m *simpleMetaInfoImpl) GetInstanceId() string {
	return m.instanceId
}
