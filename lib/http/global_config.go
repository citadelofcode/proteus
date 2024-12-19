package http

import (
	"github.com/mkbworks/proteus/lib/config"
)

const (
	ERROR_MSG_CONTENT_TYPE = "text/html"
	HEADER_LINE_SEPERATOR = "\r\n"
	REQUEST_LINE_SEPERATOR = " "
	HEADER_KEY_VALUE_SEPERATOR = ":"
)

// Collection of headers supported by the server that has a date value.
var DateHeaders []string
// List of content types supported by the web server.
var AllowedContentTypes map[string]string
// A map containing all the default server configuration values.
var ServerDefaults map[string]string
// List of all versions of HTTP supported by the web server.
var Versions map[string][]string
// List of response status codes and their associated information.
var ResponseStatusCodes []respStatus

// Initializes the global variables used in the 'http' package.
func init() {
	ServerConfig, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	DateHeaders = make([]string, 0)
	DateHeaders = append(DateHeaders, ServerConfig.DateHeaders...)
	AllowedContentTypes = ServerConfig.AllowedContentTypes
	ServerDefaults = ServerConfig.ServerDefaults
	Versions = ServerConfig.GetVersionMap()
	ResponseStatusCodes = make([]respStatus, 0)
	for _, stat := range ServerConfig.ResponseStatus {
		newStat := respStatus{
			Code: StatusCode(stat.Code),
			Message: stat.Message,
			ErrorDescription: stat.ErrorDescription,
		}
		ResponseStatusCodes = append(ResponseStatusCodes, newStat)
	}
}