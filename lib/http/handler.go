package http

import (
	"strings"
)

// Handler to fetch static file and send the file contents as response back to the client.
var StaticFileHandler = func (request *HttpRequest, response *HttpResponse) {
	targetFilePath := request.staticFilePath
	targetFilePath = strings.TrimSpace(targetFilePath)
	isCondGet, err := request.isConditionalGet(targetFilePath)
	if err != nil {
		LogError(err.Error())
		response.Status(StatusInternalServerError)
		err = response.SendError(StatusInternalServerError.GetErrorContent())
		if err != nil {
			LogError(err.Error())
		}
		return
	}

	var errSend error
	if isCondGet {
		response.Status(StatusOK)
		errSend = response.SendFile(targetFilePath, false)
	} else {
		response.Status(StatusNotModified)
		errSend = response.SendFile(targetFilePath, true)
	}

	if errSend != nil {
		LogError(errSend.Error())
	}
}

// Default error handler logic to be implemented for sending an error response back to client.
var ErrorHandler = func (request *HttpRequest, response *HttpResponse) {
	if response.StatusCode == int(StatusMethodNotAllowed) {
		response.Headers.Add("Allow", getAllowedMethods(response.Version))
	} 

	statusCode := StatusCode(response.StatusCode)
	err := response.SendError(statusCode.GetErrorContent())
	if err != nil {
		LogError(err.Error())
	}
}