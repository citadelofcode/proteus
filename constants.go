package proteus

import (
	"github.com/citadelofcode/proteus/internal"
)

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
