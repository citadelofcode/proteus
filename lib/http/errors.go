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
	return fmt.Sprintf("RequestParseError :: Section: (%s) :: Value: (%s) :: %s", rpe.Section, rpe.Value, rpe.Message)
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
	return fmt.Sprintf("RoutingError :: Route - [%s] :: %s", re.RoutePath, re.Message)
}

// Custom error to track errors raised when a HTTP response message is being formed
type ResponseError struct {
	// Refers to the part of the request which while being parsed raised the error - Header, Body, RespWrite, StatusLine are the possible values.
	Section string
	// The invalid value that caused the error.
	Value string
	// Refers to the actual error message raised.
	Message string
}

// Returns the error message associated with the instance of RequestParseError.
func (resErr ResponseError) Error() string {
	return fmt.Sprintf("ResponseError :: Section: (%s) :: Value: (%s) :: %s", resErr.Section, resErr.Value, resErr.Message)
}

// Custom error to track read timeout errors raised on incoming TCP Connections
type ReadTimeoutError struct {}

// Error message associated with the read timeout error raised for the underlying TCP connection.
func (rte *ReadTimeoutError) Error() string {
	return "Read timeout error occurred on the underlying TCP Connection."
}