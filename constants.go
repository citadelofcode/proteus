package proteus

import (
	"github.com/citadelofcode/proteus/internal"
)

// List of HTTP response codes supported by the server.
const (
	StatusOK StatusCode = internal.StatusOK
	StatusCreated StatusCode = internal.StatusCreated
	StatusAccepted StatusCode = internal.StatusAccepted
	StatusNonAuthoritative StatusCode = internal.StatusNonAuthoritative
	StatusNoContent StatusCode = internal.StatusNoContent
	StatusResetContent StatusCode = internal.StatusResetContent
	StatusPartialContent StatusCode = internal.StatusPartialContent
	StatusMultipleChoices StatusCode = internal.StatusMultipleChoices
	StatusMovedPermanently StatusCode = internal.StatusMovedPermanently
	StatusMovedTemporarily StatusCode = internal.StatusMovedTemporarily
	StatusSeeOther StatusCode = internal.StatusSeeOther
	StatusNotModified StatusCode = internal.StatusNotModified
	StatusUseProxy StatusCode = internal.StatusUseProxy
	StatusTemporaryRedirect StatusCode = internal.StatusTemporaryRedirect
	StatusBadRequest StatusCode = internal.StatusBadRequest
	StatusUnauthorized StatusCode = internal.StatusUnauthorized
	StatusPaymentRequired StatusCode = internal.StatusPaymentRequired
	StatusForbidden StatusCode = internal.StatusForbidden
	StatusNotFound StatusCode = internal.StatusNotFound
	StatusMethodNotAllowed StatusCode = internal.StatusMethodNotAllowed
	StatusNoneAcceptable StatusCode = internal.StatusNoneAcceptable
	StatusProxyAuth StatusCode = internal.StatusProxyAuth
	StatusRequestTimeout StatusCode = internal.StatusRequestTimeout
	StatusConflict StatusCode = internal.StatusConflict
	StatusGone StatusCode = internal.StatusGone
	StatusLengthMissing StatusCode = internal.StatusLengthMissing
	StatusInternalServerError StatusCode = internal.StatusInternalServerError
	StatusNotImplemented StatusCode = internal.StatusNotImplemented
	StatusBadGateway StatusCode = internal.StatusBadGateway
	StatusServiceUnavailable StatusCode = internal.StatusServiceUnavailable
	StatusGatewayTimeout StatusCode = internal.StatusGatewayTimeout
	StatusHTTPVersionNotSupported StatusCode = internal.StatusHTTPVersionNotSupported
)

// Logging levels available for server logs.
const (
	// Informational data logged to the terminal.
	INFO_LEVEL = internal.INFO_LEVEL
	// Error data logged to the terminal.
	ERROR_LEVEL = internal.ERROR_LEVEL
	// Warning(s) logged to the terminal.
	WARN_LEVEL = internal.WARN_LEVEL
)

// Format(s) in which the request processing logs can be generated.
const (
	// :remote-addr [:date[clf]] ":method :url HTTP/:http-version" :status :res[content-length]
	COMMON_LOGGER = internal.COMMON_LOGGER
	// :method :url :status :response-time ms - :res[content-length]
	DEV_LOGGER = internal.DEV_LOGGER
	// :method :url :status :res[content-length] - :response-time ms
	TINY_LOGGER = internal.TINY_LOGGER
	// :remote-addr :method :url HTTP/:http-version :status :res[content-length] - :response-time ms
	SHORT_LOGGER = internal.SHORT_LOGGER
)

// Exposes member functions to apply colors for texts before being logged to any ANSI-supported terminals.
var TextColor = internal.TextColor
