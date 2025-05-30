package internal

import (
	"bytes"
	"html/template"
)

type StatusCode int

// Structure to represent a response status code and its associated information.
type HttpStatus struct {
	// HTTP response status code.
	Code StatusCode
	// Short message for the corresponding status code.
	Message string
	// Error description for error status codes (>=400).
	ErrorDescription string
}

const (
	StatusOK StatusCode = 200
	StatusCreated StatusCode = 201
	StatusAccepted StatusCode = 202
	StatusNonAuthoritative StatusCode = 203
	StatusNoContent StatusCode = 204
	StatusResetContent StatusCode = 205
	StatusPartialContent StatusCode = 206
	StatusMultipleChoices StatusCode = 300
	StatusMovedPermanently StatusCode = 301
	StatusMovedTemporarily StatusCode = 302
	StatusSeeOther StatusCode = 303
	StatusNotModified StatusCode = 304
	StatusUseProxy StatusCode = 305
	StatusTemporaryRedirect StatusCode = 307
	StatusBadRequest StatusCode = 400
	StatusUnauthorized StatusCode = 401
	StatusPaymentRequired StatusCode = 402
	StatusForbidden StatusCode = 403
	StatusNotFound StatusCode = 404
	StatusMethodNotAllowed StatusCode = 405
	StatusNoneAcceptable StatusCode = 406
	StatusProxyAuth StatusCode = 407
	StatusRequestTimeout StatusCode = 408
	StatusConflict StatusCode = 409
	StatusGone StatusCode = 410
	StatusLengthMissing StatusCode = 411
	StatusInternalServerError StatusCode = 500
	StatusNotImplemented StatusCode = 501
	StatusBadGateway StatusCode = 502
	StatusServiceUnavailable StatusCode = 503
	StatusGatewayTimeout StatusCode = 504
	StatusHTTPVersionNotSupported StatusCode = 505
)

// Gets the minified message assosciated with a HTTP status code.
func (code StatusCode) GetStatusMessage() string {
	for _, stat := range ResponseStatusCodes {
		if stat.Code == code {
			return stat.Message
		}
	}

	return ""
}

// Gets the default error content for a HTTP status code.
func (code StatusCode) GetErrorContent() string {
	htmlTemplate := `<html>
					<head>
					<title>{{printf "%s - Response" .Code}</title>
					</head>
					<body>
					<h1>{{printf "%d - %s" .Code .Message}}</h1>
					<p>{{.ErrorDescription}}</p>
					</body>
				</html>`

	for _, stat := range ResponseStatusCodes {
		if code == stat.Code {
			temp, err := template.New("errorResponse").Parse(htmlTemplate)
			if err != nil {
				break
			}

			var tmpBytes bytes.Buffer
			err = temp.Execute(&tmpBytes, stat)
			if err != nil {
				break
			}

			return tmpBytes.String()
		}
	}

	return ""
}
