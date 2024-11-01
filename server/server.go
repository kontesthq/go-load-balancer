package server

import (
	"fmt"
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
