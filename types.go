package proteus

import (
	"github.com/citadelofcode/proteus/lib/http"
)

// A structure representing the web server that can accept and process incoming HTTP requests.
type HttpServer = http.HttpServer

// A structure to represent HTTP requests received by the web server.
type HttpRequest = http.HttpRequest

// A structure to represent HTTP response created for each request received by tht web server instance.
type HttpResponse = http.HttpResponse

// A structure to represent the collection of headers received with each request received or included in each response sent back by the server instance.
type Headers = http.Headers

// A structure to represent the collection of parameters (both query and path parameters) present in each request received by the server instance.
type Params = http.Params

// A structure to represent the status code returned with each response sent back by the server instance.
type StatusCode = http.StatusCode
