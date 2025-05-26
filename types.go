package proteus

import (
	"github.com/citadelofcode/proteus/lib/http"
)

// Represents a web server that can accept and process incoming HTTP requests.
type HttpServer = http.HttpServer

// Represents a HTTP request received by the web server.
type HttpRequest = http.HttpRequest

// Represents a HTTP response sent back in response to a HTTP request received by the web server.
type HttpResponse = http.HttpResponse

// A collection of HTTP request or response headers stored as key-value pairs.
type Headers = http.Headers

// A collection of parameters (either query or path) present in the HTTP request, stored as key-value pairs.
type Params = http.Params

// Represents the status code of the HTTP response sent back to the client.
type StatusCode = http.StatusCode

