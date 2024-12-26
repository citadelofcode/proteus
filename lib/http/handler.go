package http

import (
	"strings"
)

// Represents a handler function that is executed once any received request is parsed. You can define different handlers for different routes and HTTP methods.
type Handler func (*HttpRequest, *HttpResponse) error

// Handler to fetch static file and send the file contents as response back to the client.
var StaticFileHandler = func (request *HttpRequest, response *HttpResponse) error {
	targetFilePath := request.staticFilePath
	targetFilePath = strings.TrimSpace(targetFilePath)
	isCondGet, err := request.isConditionalGet(targetFilePath)
	if err != nil {
		return err
	}

	if !isCondGet {
		response.Status(StatusOK)
		return response.SendFile(targetFilePath, false)
	} else {
		response.Status(StatusNotModified)
		return response.SendFile(targetFilePath, true)
	}
}

// Default error handler logic to be implemented for sending an error response back to client.
var ErrorHandler = func (request *HttpRequest, response *HttpResponse) error {
	if response.StatusCode == int(StatusMethodNotAllowed) {
		response.Headers.Add("Allow", getAllowedMethods(response.Version))
	} 

	statusCode := StatusCode(response.StatusCode)
	return response.SendError(statusCode.GetErrorContent())
}