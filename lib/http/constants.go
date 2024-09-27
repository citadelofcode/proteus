package http

const (
	ERROR_MSG_CONTENT_TYPE = "text/html"
	HEADER_LINE_SEPERATOR = "\r\n"
	REQUEST_LINE_SEPERATOR = " "
	HEADER_KEY_VALUE_SEPERATOR = ":"
	VALIDATE_ROUTE_PATTERN = "^[a-zA-z][a-zA-Z0-9_/:-]*$"
)

var ServerName string
var DateHeaders []string
var DefaultHostname string
var DefaultPortNumber int