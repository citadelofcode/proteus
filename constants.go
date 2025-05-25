package proteus

import (
	"github.com/citadelofcode/proteus/lib/http"
)

const (
	StatusOK StatusCode = http.StatusOK
	StatusCreated StatusCode = http.StatusCreated
	StatusAccepted StatusCode = http.StatusAccepted
	StatusNonAuthoritative StatusCode = http.StatusNonAuthoritative
	StatusNoContent StatusCode = http.StatusNoContent
	StatusResetContent StatusCode = http.StatusResetContent
	StatusPartialContent StatusCode = http.StatusPartialContent
	StatusMultipleChoices StatusCode = http.StatusMultipleChoices
	StatusMovedPermanently StatusCode = http.StatusMovedPermanently
	StatusMovedTemporarily StatusCode = http.StatusMovedTemporarily
	StatusSeeOther StatusCode = http.StatusSeeOther
	StatusNotModified StatusCode = http.StatusNotModified
	StatusUseProxy StatusCode = http.StatusUseProxy
	StatusTemporaryRedirect StatusCode = http.StatusTemporaryRedirect
	StatusBadRequest StatusCode = http.StatusBadRequest
	StatusUnauthorized StatusCode = http.StatusUnauthorized
	StatusPaymentRequired StatusCode = http.StatusPaymentRequired
	StatusForbidden StatusCode = http.StatusForbidden
	StatusNotFound StatusCode = http.StatusNotFound
	StatusMethodNotAllowed StatusCode = http.StatusMethodNotAllowed
	StatusNoneAcceptable StatusCode = http.StatusNoneAcceptable
	StatusProxyAuth StatusCode = http.StatusProxyAuth
	StatusRequestTimeout StatusCode = http.StatusRequestTimeout
	StatusConflict StatusCode = http.StatusConflict
	StatusGone StatusCode = http.StatusGone
	StatusLengthMissing StatusCode = http.StatusLengthMissing
	StatusInternalServerError StatusCode = http.StatusInternalServerError
	StatusNotImplemented StatusCode = http.StatusNotImplemented
	StatusBadGateway StatusCode = http.StatusBadGateway
	StatusServiceUnavailable StatusCode = http.StatusServiceUnavailable
	StatusGatewayTimeout StatusCode = http.StatusGatewayTimeout
	StatusHTTPVersionNotSupported StatusCode = http.StatusHTTPVersionNotSupported
)

