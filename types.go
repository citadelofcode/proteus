package proteus

import (
	"github.com/citadelofcode/proteus/internal"
)

// Represents a web server that can accept and process incoming HTTP requests.
type HttpServer = internal.HttpServer

// Represents a HTTP request received by the web server.
type HttpRequest = internal.HttpRequest

// Represents a HTTP response sent back in response to a HTTP request received by the web server.
type HttpResponse = internal.HttpResponse

// A collection of HTTP request or response headers stored as key-value pairs.
type Headers = internal.Headers

// A collection of parameters (either query or path) present in the HTTP request, stored as key-value pairs.
type Params = internal.Params

// Represents the status code of the HTTP response sent back to the client.
type StatusCode = internal.StatusCode

// Router instance to let users declare endpoints and associated handlers.
type Router = internal.Router

// Strucure to represent a single file in the local file system.
type File = internal.File
