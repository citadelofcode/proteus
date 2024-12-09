package http

import (
	"fmt"
)

// Custom error to track errors raised when a HTTP request received is being read and parsed.
type RequestParseError struct {
	// Refers to the part of the request which while being parsed raised the error - Header, Body, QueryParams are the possible values.
	Section string
	// The invalid value that caused the error.
	Value string
	// Refers to the actual error message raised.
	Message string
}

// Returns the error message associated with the instance of RequestParseError.
func (rpe *RequestParseError) Error() string {
	return fmt.Sprintf("Error while parsing request :: Section: (%s) :: Value: (%s) :: %s", rpe.Section, rpe.Value, rpe.Message)
}

// Custom error to track errors raised by the router associated with the web server.
type RoutingError struct {
	// The target route path which has caused the issue.
	RoutePath string
	// The actual error message raised
	Message string
}

// Returns the error message associated with the RoutingError instance.
func (re *RoutingError) Error() string {
	return fmt.Sprintf("Routing Error :: Route - [%s] :: %s", re.RoutePath, re.Message)
}