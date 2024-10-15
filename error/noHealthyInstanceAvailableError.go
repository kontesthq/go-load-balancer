package error

import (
	"fmt"
)

// NoHealthyInstanceAvailableError represents an error when no healthy service instance is available.
type NoHealthyInstanceAvailableError struct {
	ServiceName string // Unexported field, making direct struct creation outside this package impossible.
}

// Error implements the error interface, returning a formatted error message.
func (e *NoHealthyInstanceAvailableError) Error() string {
	return fmt.Sprintf("No healthy instance available for service: %s", e.ServiceName)
}
