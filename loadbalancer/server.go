package loadbalancer

import (
	"fmt"
	"strconv"
	"strings"
)

type Server interface {
	GetID() string
	GetHost() string
	GetPort() int
	GetScheme() string
	IsAlive() bool
	SetAlive(isAlive bool)
	GetHostPort() string
	GetMetaInfo() MetaInfo
	GetZone() string
	SetZone(zone string)
	IsReadyToServe() bool
	SetReadyToServe(ready bool)
}

// CommonServerString returns a string representation of any Server implementation.
func CommonServerString(s Server) string {
	return fmt.Sprintf(
		"Server{id='%s', host='%s', port=%d, scheme='%s', alive=%t, zone='%s', readyToServe=%t}",
		s.GetID(),
		s.GetHost(),
		s.GetPort(),
		s.GetScheme(),
		s.IsAlive(),
		s.GetZone(),
		s.IsReadyToServe(),
	)
}

const UNKNOWN_ZONE = "UNKNOWN"

type ServerImpl struct {
	host           string
	port           int
	scheme         string
	id             string
	isAlive        bool
	zone           string
	readyToServe   bool
	simpleMetaInfo MetaInfo
}

// NewServerImpl creates a new ServerImpl instance with host and port.
func NewServerImpl(host string, port int) *ServerImpl {
	return &ServerImpl{
		host:           host,
		port:           port,
		id:             fmt.Sprintf("%s:%d", host, port),
		zone:           UNKNOWN_ZONE,
		readyToServe:   true,
		simpleMetaInfo: &simpleMetaInfoImpl{},
	}
}

// NewServerWithScheme creates a new ServerImpl with a specific scheme, host, and port.
func NewServerWithScheme(scheme, host string, port int) *ServerImpl {
	server := NewServerImpl(host, port)
	server.scheme = scheme
	return server
}

// NewServerWithID creates a new ServerImpl using host:port format for the ID.
func NewServerWithID(id string) *ServerImpl {
	server := &ServerImpl{}
	server.SetID(id)
	server.zone = UNKNOWN_ZONE
	server.readyToServe = true
	return server
}

// SetAlive sets the server's alive status.
func (s *ServerImpl) SetAlive(isAlive bool) {
	s.isAlive = isAlive
}

// IsAlive returns the server's alive status.
func (s *ServerImpl) IsAlive() bool {
	return s.isAlive
}

// SetHostPort sets the host and port for the server.
func (s *ServerImpl) SetHostPort(hostPort string) {
	s.SetID(hostPort)
}

// NormalizeID normalizes a server ID to the form "host:port".
func NormalizeID(id string) string {
	host, port := GetHostPort(id)
	if host == "" {
		return ""
	}
	return fmt.Sprintf("%s:%d", host, port)
}

// GetScheme returns the scheme (http or https) based on the prefix in ID.
func GetScheme(id string) string {
	if strings.HasPrefix(strings.ToLower(id), "http://") {
		return "http"
	} else if strings.HasPrefix(strings.ToLower(id), "https://") {
		return "https"
	}
	return ""
}

// GetHostPort parses and returns the host and port from an ID.
func GetHostPort(id string) (string, int) {
	if id == "" {
		return "", 0
	}

	var host string
	port := 80

	if strings.HasPrefix(strings.ToLower(id), "http://") {
		id = id[7:]
		port = 80
	} else if strings.HasPrefix(strings.ToLower(id), "https://") {
		id = id[8:]
		port = 443
	}

	if slashIdx := strings.Index(id, "/"); slashIdx != -1 {
		id = id[:slashIdx]
	}

	colonIdx := strings.Index(id, ":")
	if colonIdx == -1 {
		host = id
	} else {
		host = id[:colonIdx]
		if parsedPort, err := strconv.Atoi(id[colonIdx+1:]); err == nil {
			port = parsedPort
		}
	}
	return host, port
}

// SetID sets the server's ID and updates host, port, and scheme accordingly.
func (s *ServerImpl) SetID(id string) {
	host, port := GetHostPort(id)
	if host != "" {
		s.id = fmt.Sprintf("%s:%d", host, port)
		s.host = host
		s.port = port
		s.scheme = GetScheme(id)
	} else {
		s.id = ""
	}
}

// SetScheme sets the server's scheme.
func (s *ServerImpl) SetScheme(scheme string) {
	s.scheme = scheme
}

// SetPort sets the server's port.
func (s *ServerImpl) SetPort(port int) {
	s.port = port
}

// SetHost sets the server's host.
func (s *ServerImpl) SetHost(host string) {
	if host != "" {
		s.host = host
	}
}

// GetID returns the server's ID.
func (s *ServerImpl) GetID() string {
	return s.id
}

// GetHost returns the server's host.
func (s *ServerImpl) GetHost() string {
	return s.host
}

// GetPort returns the server's port.
func (s *ServerImpl) GetPort() int {
	return s.port
}

// GetScheme returns the server's scheme.
func (s *ServerImpl) GetScheme() string {
	return s.scheme
}

// GetHostPort returns the server's host:port combination.
func (s *ServerImpl) GetHostPort() string {
	return fmt.Sprintf("%s:%d", s.host, s.port)
}

// GetMetaInfo returns the server's meta information.
func (s *ServerImpl) GetMetaInfo() MetaInfo {
	return s.simpleMetaInfo
}

// Equals checks if another object is equal to this server.
func (s *ServerImpl) Equals(other interface{}) bool {
	otherServer, ok := other.(*ServerImpl)
	return ok && otherServer.GetID() == s.GetID()
}

// HashCode returns a hash code for the server.
func (s *ServerImpl) HashCode() int {
	hash := 7
	hash = 31*hash + len(s.GetID())
	return hash
}

// GetZone returns the server's zone.
func (s *ServerImpl) GetZone() string {
	return s.zone
}

// SetZone sets the server's zone.
func (s *ServerImpl) SetZone(zone string) {
	s.zone = zone
}

// IsReadyToServe returns whether the server is ready to serve.
func (s *ServerImpl) IsReadyToServe() bool {
	return s.readyToServe
}

// SetReadyToServe sets whether the server is ready to serve.
func (s *ServerImpl) SetReadyToServe(readyToServe bool) {
	s.readyToServe = readyToServe
}
