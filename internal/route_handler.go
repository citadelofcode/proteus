package internal

import (
	"strings"
)

// Represents a handler function that is executed once any received request is parsed.
// You can define different handlers for different routes and HTTP methods.
type RouteHandler func (*HttpRequest, *HttpResponse)

// Handler to fetch static file and send the file contents as response back to the client.
var StaticFileHandler = func (request *HttpRequest, response *HttpResponse) {
	targetFilePath, ok := request.Locals["StaticFilePath"].(string)
	if !ok {
		request.Server.Log("Static file path not available in the request instance", ERROR_LEVEL)
		return
	}
	targetFilePath = strings.TrimSpace(targetFilePath)
	isCondGet, err := request.IsConditionalGet(targetFilePath)
	if err != nil {
		request.Server.Log(err.Error(), ERROR_LEVEL)
	}

	if !isCondGet {
		response.Status(StatusOK)
		err := response.SendFile(targetFilePath, false)
		if err != nil {
			request.Server.Log(err.Error(), ERROR_LEVEL)
		}
	} else {
		response.Status(StatusNotModified)
		err := response.SendFile(targetFilePath, true)
		if err != nil {
			request.Server.Log(err.Error(), ERROR_LEVEL)
		}
	}
}

// Default error handler logic to be implemented for sending an error response back to client.
var ErrorHandler = func (request *HttpRequest, response *HttpResponse) {
	if response.StatusCode < 400 {
		request.Server.Log("Response Status code should be 400 or above to invoke the default error handler", ERROR_LEVEL)
		return
	}

	if response.StatusCode == int(StatusMethodNotAllowed) {
		response.Headers.Add("Allow", GetAllowedMethods(response.Version))
	}

	statusCode := StatusCode(response.StatusCode)
	err := response.SendError(statusCode.GetErrorContent())
	if err != nil {
		request.Server.Log(err.Error(), ERROR_LEVEL)
	}
}
