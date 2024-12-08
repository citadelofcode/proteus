package http

import (
	"strings"
)

// Handler to fetch static file and send the file contents as response back to the client.
var StaticFileHandler = func (request *HttpRequest, response *HttpResponse) {
	targetFilePath := request.staticFilePath
	targetFilePath = strings.TrimSpace(targetFilePath)
	if !request.isConditionalGet(targetFilePath) {
		response.Status(StatusOK)
		response.SendFile(targetFilePath, false)
	} else {
		response.Status(StatusNotModified)
		response.SendFile(targetFilePath, true)
	}
}

// Default error handler logic to be implemented for sending an error response back to client.
var ErrorHandler = func (request *HttpRequest, response *HttpResponse) {
	if response.StatusCode == int(StatusMethodNotAllowed) {
		response.AddHeader("Allow", getAllowedMethods(response.Version))
	} 

	statusCode := StatusCode(response.StatusCode)
	response.SendError(statusCode.GetErrorContent())
}